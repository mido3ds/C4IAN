package main

import (
	"fmt"
	"net"
	"testing"
)

func Benchmark_MarshalBinary(b *testing.B) {
	ip0 := net.IP([]byte{0x01, 0x02, 0x03, 0x04})
	ip1 := net.IP([]byte{0x05, 0x06, 0x07, 0x08})
	ip2 := net.IP([]byte{0x09, 0x0A, 0x0B, 0x0C})
	ip4 := net.IP([]byte{0x0D, 0x0E, 0x0F, 0x10})

	neighborsTable := NewNeighborsTable()

	neighborsTable.Set(ip0, &NeighborEntry{cost: 1})
	neighborsTable.Set(ip1, &NeighborEntry{cost: 2})
	neighborsTable.Set(ip2, &NeighborEntry{cost: 3})
	neighborsTable.Set(ip4, &NeighborEntry{cost: 4})

	payload := neighborsTable.MarshalBinary()

	newNeighborsTable, err := UnmarshalNeighborsTable(payload)
	fmt.Println(err)
	fmt.Println(newNeighborsTable)
}
