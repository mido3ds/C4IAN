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
	m hashmap.HashMap
}

type ForwardingEntry struct {
	NextHopMAC net.HardwareAddr
	DestZoneID uint32
}

func NewForwardTable() *ForwardTable {
	return &ForwardTable{}
}

// Get returns value associated with the given key, and whether the key existed or not
func (f *ForwardTable) Get(destIP []byte) (*ForwardingEntry, bool) {
	v, ok := f.m.Get(ipToUInt32(destIP))
	if v == nil {
		return nil, false
	}
	return v.(*ForwardingEntry), ok
}

func (f *ForwardTable) Set(destIP []byte, entry *ForwardingEntry) {
	if entry == nil {
		panic(fmt.Errorf("you can't enter nil entry"))
	}
	f.m.Set(ipToUInt32(destIP), entry)
}

// Del silently fails if key doesn't exist
func (f *ForwardTable) Del(destIP []byte) {
	f.m.Del(ipToUInt32(destIP))
}

func (f *ForwardTable) Len() int {
	return f.m.Len()
}

func ipToUInt32(ip []byte) uint32 {
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}
