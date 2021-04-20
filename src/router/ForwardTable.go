package main

import (
	"fmt"
	"net"

	"github.com/cornelk/hashmap"
)

// ForwardTable is lock-free thread-safe hash table
// optimized for fastest read access
// key: 4 bytes IPv4, value: *ForwardingEntry
type ForwardTable struct {
	m *hashmap.HashMap
}

type ForwardingEntry struct {
	NextHopMAC net.HardwareAddr
	DestZoneID uint32
}

func NewForwardTable() *ForwardTable {
	return &ForwardTable{
		m: &hashmap.HashMap{},
	}
}

// Get returns value associated with the given key, and whether the key existed or not
func (f *ForwardTable) Get(destIP net.IP) (*ForwardingEntry, bool) {
	v, ok := f.m.Get(IPv4ToUInt32(destIP))
	if !ok {
		return nil, false
	}
	return v.(*ForwardingEntry), true
}

func (f *ForwardTable) Set(destIP net.IP, entry *ForwardingEntry) {
	if entry == nil {
		panic(fmt.Errorf("you can't enter nil entry"))
	}
	f.m.Set(IPv4ToUInt32(destIP), entry)
}

// Del silently fails if key doesn't exist
func (f *ForwardTable) Del(destIP net.IP) {
	f.m.Del(IPv4ToUInt32(destIP))
}

func (f *ForwardTable) Len() int {
	return f.m.Len()
}

func (f *ForwardTable) Clear() {
	// Create a new hashmap as the underlying hashmap lacks a clear function :)
	f.m = &hashmap.HashMap{}
}

func (f *ForwardTable) String() string {
	s := "&ForwardTable{"
	for item := range f.m.Iter() {
		s += fmt.Sprintf(" (ip=%#v,mac=%#v)", UInt32ToIPv4(item.Key.(uint32)).String(), item.Value.(*ForwardingEntry).NextHopMAC.String())

	}
	s += " }"
	return s
}
