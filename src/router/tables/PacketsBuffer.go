package tables

import (
	"fmt"
	"net"
	"time"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/ip"
)

const BufferAge = 10
const MaxNumOfSearches = 3

// PacketsBugger is lock-free thread-safe hash table
// optimized for fastest read access
// key: 4 bytes dst IPv4, value: queue for msgs
type PacketsBuffer struct {
	m                   hashmap.HashMap
	findDstZoneCallback func(dstIP net.IP)
}

type BufferEntry struct {
	packetsQueue     [][]byte
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
func (p *PacketsBuffer) Get(dstIP net.IP) ([][]byte, bool) {
	v, ok := p.m.Get(IPv4ToUInt32(dstIP))
	if !ok {
		return nil, false
	}

	// Stop the timer
	timer := v.(*BufferEntry).ageTimer
	timer.Stop()

	return v.(*BufferEntry).packetsQueue, true
}

func (p *PacketsBuffer) AppendPacket(dstIP net.IP, packet []byte) {
	v, ok := p.m.Get(IPv4ToUInt32(dstIP))
	var queue [][]byte
	if ok {
		// enqueue the new packet
		queue = v.(*BufferEntry).packetsQueue
		queue = append(queue, packet)
		timer := v.(*BufferEntry).ageTimer
		numOfDstSearches := v.(*BufferEntry).numOfDstSearches

		p.m.Set(IPv4ToUInt32(dstIP), &BufferEntry{packetsQueue: queue, ageTimer: timer, numOfDstSearches: numOfDstSearches})

	} else {
		// Start new Timer
		fireFunc := bufferFireTimer(dstIP, p)
		newTimer := time.AfterFunc(BufferAge*time.Second, fireFunc)

		// make new queue for upcoming messages to the same destination
		queue = make([][]byte, 0)
		queue = append(queue, packet)
		p.m.Set(IPv4ToUInt32(dstIP), &BufferEntry{packetsQueue: queue, ageTimer: newTimer, numOfDstSearches: 0})
	}
}

// Del silently fails if key doesn't exist
func (p *PacketsBuffer) Del(dstIP net.IP) {
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
	if numOfDstSearches < MaxNumOfSearches {
		p.findDstZoneCallback(dstIP)
		// Start new Timer
		fireFunc := bufferFireTimer(dstIP, p)
		newTimer := time.AfterFunc(BufferAge*time.Second, fireFunc)
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
