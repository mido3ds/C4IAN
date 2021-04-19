package main

import (
	"fmt"
	"net"

	"github.com/cornelk/hashmap"
)

// NeighborsTable is lock-free thread-safe hash table
// optimized for fastest read access
// key: 4 bytes IPv4, value: *NeighborEntry
type NeighborsTable struct {
	m *hashmap.HashMap
}

type NeighborEntry struct {
	MAC  net.HardwareAddr
	cost uint32
}

func NewNeighborsTable() *NeighborsTable {
	return &NeighborsTable{
		m: &hashmap.HashMap{},
	}
}

// Get returns value associated with the given key, and whether the key existed or not
func (n *NeighborsTable) Get(ip net.IP) (*NeighborEntry, bool) {
	v, ok := n.m.Get(IPv4ToUInt32(ip))
	if !ok {
		return nil, false
	}
	return v.(*NeighborEntry), true
}

func (n *NeighborsTable) Set(ip net.IP, entry *NeighborEntry) {
	if entry == nil {
		panic(fmt.Errorf("you can't enter nil entry"))
	}
	n.m.Set(IPv4ToUInt32(ip), entry)
}

// Del silently fails if key doesn't exist
func (n *NeighborsTable) Del(ip net.IP) {
	n.m.Del(IPv4ToUInt32(ip))
}

func (n *NeighborsTable) Len() int {
	return n.m.Len()
}

func (n *NeighborsTable) Clear() {
	// Create a new hashmap as the underlying hashmap lacks a clear function :)
	n.m = &hashmap.HashMap{}
}

func (n *NeighborsTable) String() string {
	s := "&NeighborsTable{"
	for item := range n.m.Iter() {
		s += fmt.Sprintf(" (ip=%#v,mac=%#v)", UInt32ToIPv4(item.Key.(uint32)).String(), item.Value.(*NeighborEntry).MAC.String())

	}
	s += " }"
	return s
}

func (n *NeighborsTable) getTableHash() []byte {
	return Hash_SHA3([]byte(n.String()))
}
