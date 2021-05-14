package tables

import (
	"fmt"
	"log"
	"net"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/msec"
)

// NeighborsTable is lock-free thread-safe hash table
// optimized for fastest read access
// key: 4 bytes IPv4, value: *NeighborEntry
type NeighborsTable struct {
	m *hashmap.HashMap
}

// 8 bytes (NodeID) + 2 bytes (Cost)
// The MAC addresses are not transfered in LSR packets, so they are not marshalled
const neighborEntryLen = 10

type NeighborEntry struct {
	MAC  net.HardwareAddr
	Cost uint16
}

func NewNeighborsTable() *NeighborsTable {
	return &NeighborsTable{
		m: &hashmap.HashMap{},
	}
}

func UnmarshalNeighborsTable(payload []byte) (*NeighborsTable, bool) {
	// extract number of entries
	numberOfEntries := uint16(payload[0])<<8 | uint16(payload[1])
	payloadLen := numberOfEntries * neighborEntryLen

	// extract checksum
	csum := uint16(payload[2])<<8 | uint16(payload[3])
	if csum != BasicChecksum(payload[4:4+payloadLen]) {
		return nil, false
	}

	payload = payload[4 : 4+payloadLen]
	neighborsTable := &NeighborsTable{m: &hashmap.HashMap{}}

	start := 0
	for i := 0; i < int(numberOfEntries); i++ {
		nodeID := uint64(payload[start])<<56 |
			uint64(payload[start+1])<<48 |
			uint64(payload[start+2])<<40 |
			uint64(payload[start+3])<<32 |
			uint64(payload[start+4])<<24 |
			uint64(payload[start+5])<<16 |
			uint64(payload[start+6])<<8 |
			uint64(payload[start+7])

		cost := uint16(payload[start+8])<<8 | uint16(payload[start+9])
		neighborsTable.Set(NodeID(nodeID), &NeighborEntry{Cost: cost})
		start += neighborEntryLen
	}

	return neighborsTable, true
}

// Get returns value associated with the given key, and whether the key existed or not
func (n *NeighborsTable) Get(nodeID NodeID) (*NeighborEntry, bool) {
	v, ok := n.m.Get(uint64(nodeID))
	if !ok {
		return nil, false
	}
	return v.(*NeighborEntry), true
}

func (n *NeighborsTable) Set(nodeID NodeID, entry *NeighborEntry) {
	if entry == nil {
		log.Panic("you can't enter nil entry")
	}
	n.m.Set(uint64(nodeID), entry)
}

// Del silently fails if key doesn't exist
func (n *NeighborsTable) Del(nodeID NodeID) {
	n.m.Del(uint64(nodeID))
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
		s += fmt.Sprintf(" (%#v,mac=%#v,Cost=%d)", NodeID(item.Key.(uint64)).String(), item.Value.(*NeighborEntry).MAC.String(), item.Value.(*NeighborEntry).Cost)
	}
	s += " }"
	return s
}

func (n *NeighborsTable) MarshalBinary() []byte {
	payloadLen := n.Len() * neighborEntryLen
	payload := make([]byte, payloadLen+4)

	// 0:2 => number of entries
	payload[0] = byte(uint16(n.Len()) >> 8)
	payload[1] = byte(uint16(n.Len()))

	start := 4
	for item := range n.m.Iter() {
		// Insert Node ID: 8 bytes
		nodeID := item.Key.(uint64)
		payload[start] = byte(nodeID >> 56)
		payload[start+1] = byte(nodeID >> 48)
		payload[start+2] = byte(nodeID >> 40)
		payload[start+3] = byte(nodeID >> 32)
		payload[start+4] = byte(nodeID >> 24)
		payload[start+5] = byte(nodeID >> 16)
		payload[start+6] = byte(nodeID >> 8)
		payload[start+7] = byte(nodeID)

		// Insert cost: 2 bytes
		cost := item.Value.(*NeighborEntry).Cost
		payload[start+8] = byte(cost >> 8)
		payload[start+9] = byte(cost)

		start += neighborEntryLen
	}

	// add checksum
	csum := BasicChecksum(payload[4 : 4+payloadLen])
	payload[2] = byte(csum >> 8)
	payload[3] = byte(csum)

	return payload[:]
}

// The table hash depends on who the neighbors are, regardless of the costs and the MAC addresses
func (n *NeighborsTable) GetTableHash() []byte {
	// Create a list of node ids in the table
	nodeIDs := make([]uint64, 0, n.m.Len())
	for item := range n.m.Iter() {
		nodeIDs = append(nodeIDs, item.Key.(uint64))
	}

	// Convert to a bytes slice to hash
	b := make([]byte, len(nodeIDs)*8)
	start := 0
	for _, nodeID := range nodeIDs {
		b[start] = byte(nodeID >> 56)
		b[start+1] = byte(nodeID >> 48)
		b[start+2] = byte(nodeID >> 40)
		b[start+3] = byte(nodeID >> 32)
		b[start+4] = byte(nodeID >> 24)
		b[start+5] = byte(nodeID >> 16)
		b[start+6] = byte(nodeID >> 8)
		b[start+7] = byte(nodeID)
		start += 8
	}
	return HashSHA3(b)
}
