package tables

import (
	"fmt"
	"net"
	"time"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/constants"
	. "github.com/mido3ds/C4IAN/src/router/ip"
)

type SendPacketCallback = func([]byte, net.IP)
type PacketsQueue = []*PacketEntry

// PacketsBugger is lock-free thread-safe hash table
// optimized for fastest read access
// key: 4 bytes dst IPv4, value: queue for msgs
type PacketsBuffer struct {
	m                   hashmap.HashMap
	findDstZoneCallback func(dstIP net.IP)
}

type PacketEntry struct {
	payload  []byte
	callback SendPacketCallback
}

func (p *PacketEntry) Send(dstIP net.IP) {
	p.callback(p.payload, dstIP)
}

type BufferEntry struct {
	packetsQueue     PacketsQueue
	ageTimer         *time.Timer
	numOfDstSearches uint8
}

func NewPacketsBuffer(findDstZoneCallback func(dstIP net.IP)) *PacketsBuffer {
	return &PacketsBuffer{
		findDstZoneCallback: findDstZoneCallback,
	}
}

// Get returns value associated with the given key
// and whether the key existed or not
func (p *PacketsBuffer) Get(dstIP net.IP) (PacketsQueue, bool) {
	v, ok := p.m.Get(IPv4ToUInt32(dstIP))
	if !ok {
		return nil, false
	}

	return v.(*BufferEntry).packetsQueue, true
}

func (p *PacketsBuffer) AppendPacket(dstIP net.IP, packet []byte, sendCallback SendPacketCallback) {
	v, ok := p.m.Get(IPv4ToUInt32(dstIP))
	var queue PacketsQueue
	if ok {
		// enqueue the new packet
		queue = v.(*BufferEntry).packetsQueue
		queue = append(queue, &PacketEntry{payload: packet, callback: sendCallback})
		timer := v.(*BufferEntry).ageTimer
		numOfDstSearches := v.(*BufferEntry).numOfDstSearches

		p.m.Set(IPv4ToUInt32(dstIP), &BufferEntry{packetsQueue: queue, ageTimer: timer, numOfDstSearches: numOfDstSearches})

	} else {
		// Start new Timer
		fireFunc := bufferFireTimer(dstIP, p)
		newTimer := time.AfterFunc(DZDRetryTimeout, fireFunc)

		// make new queue for upcoming messages to the same destination
		queue = make([]*PacketEntry, 0)
		queue = append(queue, &PacketEntry{payload: packet, callback: sendCallback})
		p.m.Set(IPv4ToUInt32(dstIP), &BufferEntry{packetsQueue: queue, ageTimer: newTimer, numOfDstSearches: 0})
	}
}

// Del silently fails if key doesn't exist
func (p *PacketsBuffer) Del(dstIP net.IP) {
	v, ok := p.m.Get(IPv4ToUInt32(dstIP))
	if !ok {
		return
	}

	v.(*BufferEntry).ageTimer.Stop()
	p.m.Del(IPv4ToUInt32(dstIP))
}

func (p *PacketsBuffer) Len() int {
	return p.m.Len()
}

func (p *PacketsBuffer) String() string {
	s := "&PacketsBuffer{\n"
	for item := range p.m.Iter() {
		s += fmt.Sprintf("IP=%#v\n", UInt32ToIPv4(item.Key.(uint32)).String())
		s += "Buffer: "
		for _, value := range item.Value.(*BufferEntry).packetsQueue {
			s += fmt.Sprint(value)
		}
		s += "\n"
	}
	s += " }"
	return s
}

func bufferFireTimerHelper(dstIP net.IP, p *PacketsBuffer) {
	v, ok := p.m.Get(IPv4ToUInt32(dstIP))
	if !ok {
		return
	}
	numOfDstSearches := v.(*BufferEntry).numOfDstSearches
	if numOfDstSearches < DZDMaxRetry {
		p.findDstZoneCallback(dstIP)
		// Start new Timer
		fireFunc := bufferFireTimer(dstIP, p)
		newTimer := time.AfterFunc(DZDRetryTimeout, fireFunc)
		queue := v.(*BufferEntry).packetsQueue

		p.m.Set(IPv4ToUInt32(dstIP), &BufferEntry{packetsQueue: queue, ageTimer: newTimer, numOfDstSearches: numOfDstSearches + 1})
	} else {
		p.Del(dstIP)
	}
}

func bufferFireTimer(dstIP net.IP, p *PacketsBuffer) func() {
	return func() {
		bufferFireTimerHelper(dstIP, p)
	}
}
