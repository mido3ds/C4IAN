package tables

import (
	"fmt"
	"net"
	"testing"
)

func TestMarshalBinary(t *testing.T) {
	ip0 := net.IP([]byte{0x01, 0x02, 0x03, 0x04})
	ip1 := net.IP([]byte{0x05, 0x06, 0x07, 0x08})
	ip2 := net.IP([]byte{0x09, 0x0A, 0x0B, 0x0C})
	ip4 := net.IP([]byte{0x0D, 0x0E, 0x0F, 0x10})

	neighborsTable := NewNeighborsTable()

	neighborsTable.Set(ToNodeID(ip0), &NeighborEntry{Cost: 1})
	neighborsTable.Set(ToNodeID(ip1), &NeighborEntry{Cost: 2})
	neighborsTable.Set(ToNodeID(ip2), &NeighborEntry{Cost: 3})
	neighborsTable.Set(ToNodeID(ip4), &NeighborEntry{Cost: 4})

	payload := neighborsTable.MarshalBinary()

	newNeighborsTable, ok := UnmarshalNeighborsTable(payload)
	fmt.Println(ok)
	fmt.Println(newNeighborsTable)
}
