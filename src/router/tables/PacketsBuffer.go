package tables

import (
	"fmt"
	"net"
	"time"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/ip"
)

const Age = 5

// PacketsBugger is lock-free thread-safe hash table
// optimized for fastest read access
// key: 4 bytes dst IPv4, value: queue for msgs
type PacketsBuffer struct {
	m hashmap.HashMap
}

type BufferEntry struct {
	packetsQueue [][]byte
	ageTimer     *time.Timer
}

func NewPacketsBuffer() *PacketsBuffer {
	return &PacketsBuffer{}
}

// Get returns value associated with the given key
// and whether the key existed or not
func (p *PacketsBuffer) Get(dstIP net.IP) ([][]byte, bool) {
	v, ok := p.m.Get(IPv4ToUInt32(dstIP))
	if !ok {
		return nil, false
	}
	return v.(*BufferEntry).packetsQueue, true
}

// Set the srcIP to a new sequence number
// Restart the timer attached to that src
func (p *PacketsBuffer) AppendPacket(dstIP net.IP, packet []byte) {
	v, ok := p.m.Get(IPv4ToUInt32(dstIP))
	var queue [][]byte
	if ok {
		// Stop the previous timer if it wasn't fired
		timer := v.(*BufferEntry).ageTimer
		timer.Stop()

		// enqueue the new packet
		queue = v.(*BufferEntry).packetsQueue
		queue = append(queue, packet)

	} else {
		// make new queue for upcoming messages to the same destination
		queue = make([][]byte, 0)
		queue = append(queue, packet)
	}

	// Start new Timer
	fireFunc := bufferFireTimer(dstIP, p)
	newTimer := time.AfterFunc(Age*time.Second, fireFunc)

	p.m.Set(IPv4ToUInt32(dstIP), &BufferEntry{packetsQueue: queue, ageTimer: newTimer})

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
	p.Del(dstIP)
}

func bufferFireTimer(dstIP net.IP, p *PacketsBuffer) func() {
	return func() {
		bufferFireTimerHelper(dstIP, p)
	}
}
