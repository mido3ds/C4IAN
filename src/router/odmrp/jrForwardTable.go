package odmrp

import (
	"fmt"
	"log"
	"net"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/ip"
)

// jrForwardTable is lock-free thread-safe hash table
// for multicast forwarding
// key: 4 bytes src IPv4, value: *jrForwardEntry
type jrForwardTable struct {
	m *hashmap.HashMap
}

type jrForwardEntry struct {
	seqNum  uint64
	nextHop net.HardwareAddr
	// TODO: timer
}

func newJRForwardTable() *jrForwardTable {
	return &jrForwardTable{
		m: &hashmap.HashMap{},
	}
}

// Get returns value associated with the given key, and whether the key existed or not
func (f *jrForwardTable) Get(src net.IP) (*jrForwardEntry, bool) {
	v, ok := f.m.Get(IPv4ToUInt32(src))
	if !ok {
		return nil, false
	}

	return v.(*jrForwardEntry), true
}

func (f *jrForwardTable) Set(src net.IP, entry *jrForwardEntry) {
	if entry == nil {
		log.Panic("you can't enter nil entry")
	}
	f.m.Set(IPv4ToUInt32(src), entry)
}

// Del silently fails if key doesn't exist
func (f *jrForwardTable) Del(src net.IP) {
	f.m.Del(IPv4ToUInt32(src))
}

func (f *jrForwardTable) Len() int {
	return f.m.Len()
}

// Clear jrForwardTable
func (f *jrForwardTable) Clear() {
	f.m = &hashmap.HashMap{}
}

func (f *jrForwardTable) String() string {
	s := "&jrForwardTable{"
	for item := range f.m.Iter() {
		v := item.Value.(*jrForwardEntry)
		s += fmt.Sprintf(" (srcIP=%#v,seq=%d,nexhop=%#v)", UInt32ToIPv4(item.Key.(uint32)).String(), v.seqNum, v.nextHop)
	}
	s += " }"

	return s
}
