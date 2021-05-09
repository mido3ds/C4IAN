package tables

import (
	"log"
	"net"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/ip"
)

// MultiForwardTable is lock-free thread-safe hash table
// for multicast forwarding
// key: 4 bytes IPv4, value: *MultiForwardingEntry
/*
------------------------------------
grpIP | Set of nextHops
------------------------------------
*/
type MultiForwardTable struct {
	m      *hashmap.HashMap
	timers *TimersQueue
}

func NewMultiForwardTable(timers *TimersQueue) *MultiForwardTable {
	return &MultiForwardTable{
		m:      &hashmap.HashMap{},
		timers: timers,
	}
}

// Get returns value associated with the given key, and whether the key existed or not
func (f *MultiForwardTable) Get(grpIP net.IP) (*MultiForwardEntrySet, bool) {
	v, ok := f.m.Get(IPv4ToUInt32(grpIP))
	if !ok {
		return nil, false
	}

	return v.(*MultiForwardEntrySet), true
}

func (f *MultiForwardTable) Set(grpIP net.IP, nextHop net.HardwareAddr) {
	if !grpIP.IsMulticast() {
		log.Panic("Group IP Is Not Multicast IP")
	}
	grpIPkey := IPv4ToUInt32(grpIP)
	v, ok := f.m.Get(grpIPkey)
	var entry *MultiForwardEntrySet
	if ok {
		entry = v.(*MultiForwardEntrySet)
	} else {
		entry = NewMultiForwardEntrySet(f.timers)
	}
	entry.Set(nextHop)
	f.m.Set(grpIPkey, entry)
}

// Del silently fails if key doesn't exist
func (f *MultiForwardTable) Del(grpIP net.IP) {
	f.m.Del(IPv4ToUInt32(grpIP))
}

func (f *MultiForwardTable) Len() int {
	return f.m.Len()
}

// Clear MultiForwardTable
func (f *MultiForwardTable) Clear() {
	f.m = &hashmap.HashMap{}
}

func (f *MultiForwardTable) String() string {
	s := "&MultiForwardTable{"
	for item := range f.m.Iter() {
		s += item.Value.(*MultiForwardEntrySet).String()
	}
	s += " }"

	return s
}
