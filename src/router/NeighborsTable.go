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
	cost uint8
}

func NewNeighborsTable() *NeighborsTable {
	return &NeighborsTable{
		m: &hashmap.HashMap{},
	}
}

func UnmarshalNeighborsTable(payload []byte) (*NeighborsTable, bool) {
	// extract number of entries
	numberOfEntries := uint16(payload[0])<<8 | uint16(payload[1])
	payloadLen := numberOfEntries * 5

	// extract checksum
	csum := uint16(payload[2])<<8 | uint16(payload[3])
	if csum != BasicChecksum(payload[4:4+payloadLen]) {
		return nil, false
	}

	payload = payload[4 : 4+payloadLen]
	neighborsTable := &NeighborsTable{m: &hashmap.HashMap{}}

	start := 0
	for i := 0; i < int(numberOfEntries); i++ {
		IP := net.IP(payload[start : start+4])
		cost := uint8(payload[start+4])
		neighborsTable.Set(IP, &NeighborEntry{cost: cost})
		start += 5
	}

	return neighborsTable, true
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
	s := "&NeighborsTable{\n"
	for item := range n.m.Iter() {
		s += fmt.Sprintf(" (ip=%#v,mac=%#v,cost=%d)\n", UInt32ToIPv4(item.Key.(uint32)).String(), item.Value.(*NeighborEntry).MAC.String(), item.Value.(*NeighborEntry).cost)
	}
	s += " }"
	return s
}

func (n *NeighborsTable) MarshalBinary() []byte {
	payloadLen := n.Len() * 5
	payload := make([]byte, payloadLen+4)

	// 0:2 => number of entries
	payload[0] = byte(uint16(n.Len()) >> 8)
	payload[1] = byte(uint16(n.Len()))

	start := 4
	for item := range n.m.Iter() {
		// Insert IP: 4 Bytes
		IP := item.Key.(uint32)
		payload[start] = byte(IP >> 24)
		payload[start+1] = byte(IP >> 16)
		payload[start+2] = byte(IP >> 8)
		payload[start+3] = byte(IP)

		// Insert cost: 1 byte
		payload[start+4] = byte(item.Value.(*NeighborEntry).cost)
		start += 5
	}

	// add checksum
	csum := BasicChecksum(payload[4 : 4+payloadLen])
	payload[2] = byte(csum >> 8)
	payload[3] = byte(csum)

	return payload[:]
}

func (n *NeighborsTable) GetTableHash() []byte {
	return Hash_SHA3([]byte(n.String()))
}
