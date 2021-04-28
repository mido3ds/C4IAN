package odmrp

import (
	"fmt"
	"log"
	"net"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/ip"
)

// RoutingTable is lock-free thread-safe hash table
// for multicast forwarding
// key: 4 bytes src IPv4, value: *jrForwardEntry
type RoutingTable struct {
	m *hashmap.HashMap
}

type jrForwardEntry struct {
	nextHop net.HardwareAddr
	// TODO: timer
}

func newRoutingTable() *RoutingTable {
	return &RoutingTable{
		m: &hashmap.HashMap{},
	}
}

// Get returns value associated with the given key, and whether the key existed or not
func (f *RoutingTable) Get(src net.IP) (*jrForwardEntry, bool) {
	v, ok := f.m.Get(IPv4ToUInt32(src))
	if !ok {
		return nil, false
	}

	return v.(*jrForwardEntry), true
}

func (f *RoutingTable) Set(src net.IP, entry *jrForwardEntry) {
	if entry == nil {
		log.Panic("you can't enter nil entry")
	}
	f.m.Set(IPv4ToUInt32(src), entry)
}

// Del silently fails if key doesn't exist
func (f *RoutingTable) Del(src net.IP) {
	f.m.Del(IPv4ToUInt32(src))
}

func (f *RoutingTable) Len() int {
	return f.m.Len()
}

// Clear RoutingTable
func (f *RoutingTable) Clear() {
	f.m = &hashmap.HashMap{}
}

func (f *RoutingTable) String() string {
	s := "&RoutingTable{"
	for item := range f.m.Iter() {
		v := item.Value.(*jrForwardEntry)
		s += fmt.Sprintf(" (srcIP=%#v, nexhop=%#v)", UInt32ToIPv4(item.Key.(uint32)).String(), v.nextHop)
	}
	s += " }"

	return s
}
