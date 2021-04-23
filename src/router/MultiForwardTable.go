package main

import (
	"fmt"
	"log"
	"net"

	"github.com/cornelk/hashmap"
)

// MultiForwardTable is lock-free thread-safe hash table
// for multicast forwarding
// key: 4 bytes IPv4, value: *MultiForwardingEntry
type MultiForwardTable struct {
	m *hashmap.HashMap
}

type MultiForwardingEntry struct {
	NextHopMACs []net.HardwareAddr
}

func NewMultiForwardTable() *MultiForwardTable {
	return &MultiForwardTable{
		m: &hashmap.HashMap{},
	}
}

// Get returns value associated with the given key, and whether the key existed or not
func (f *MultiForwardTable) Get(grpIP net.IP) (*MultiForwardingEntry, bool) {
	v, ok := f.m.Get(IPv4ToUInt32(grpIP))
	if !ok {
		return nil, false
	}

	return v.(*MultiForwardingEntry), true
}

func (f *MultiForwardTable) Set(grpIP net.IP, entry *MultiForwardingEntry) {
	if entry == nil {
		log.Panic("you can't enter nil entry")
	}
	if !grpIP.IsMulticast() {
		log.Panic("Group IP Is Not Multicast IP")
	}
	f.m.Set(IPv4ToUInt32(grpIP), entry)
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
		s += fmt.Sprintf(" (ip=%#v,mac=%#v)", UInt32ToIPv4(item.Key.(uint32)).String(), item.Value.(*MultiForwardingEntry).NextHopMACs)
	}
	s += " }"

	return s
}
