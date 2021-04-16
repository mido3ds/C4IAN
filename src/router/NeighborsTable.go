package main

import (
	"fmt"
	"log"
	"net"

	"github.com/cornelk/hashmap"
)

// NeighborsTable is lock-free thread-safe hash table
// optimized for fastest read access
// key: 4 bytes IPv4, value: *NeighborEntry
type NeighborsTable struct {
	m hashmap.HashMap
}

type NeighborEntry struct {
	MAC net.HardwareAddr
}

func NewNeighborsTable() *NeighborsTable {
	return &NeighborsTable{}
}

// Get returns value associated with the given key, and whether the key existed or not
func (f *NeighborsTable) Get(ip net.IP) (*NeighborEntry, bool) {
	v, ok := f.m.Get(IPv4ToUInt32(ip))
	if !ok {
		return nil, false
	}
	return v.(*NeighborEntry), true
}

func (f *NeighborsTable) Set(ip net.IP, entry *NeighborEntry) {
	if entry == nil {
		panic(fmt.Errorf("you can't enter nil entry"))
	}
	f.m.Set(IPv4ToUInt32(ip), entry)
}

// Del silently fails if key doesn't exist
func (f *NeighborsTable) Del(ip net.IP) {
	f.m.Del(IPv4ToUInt32(ip))
}

func (f *NeighborsTable) Len() int {
	return f.m.Len()
}

func (f *NeighborsTable) Display() {
	for item := range f.m.Iter() {
		log.Println(item.Key, item.Value)
	}
}
