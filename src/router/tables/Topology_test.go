package tables

import (
	"fmt"
	"net"
	"testing"

	. "github.com/mido3ds/C4IAN/src/router/ip"
)

// Example : https://www.geeksforgeeks.org/dijkstras-shortest-path-algorithm-greedy-algo-7/
func TestTopology(t *testing.T) {
	ip0 := net.IP([]byte{0x00, 0x00, 0x00, 0x00})
	ip1 := net.IP([]byte{0x00, 0x00, 0x00, 0x01})
	ip2 := net.IP([]byte{0x00, 0x00, 0x00, 0x02})
	ip3 := net.IP([]byte{0x00, 0x00, 0x00, 0x03})
	ip4 := net.IP([]byte{0x00, 0x00, 0x00, 0x04})
	ip5 := net.IP([]byte{0x00, 0x00, 0x00, 0x05})
	ip6 := net.IP([]byte{0x00, 0x00, 0x00, 0x06})
	ip7 := net.IP([]byte{0x00, 0x00, 0x00, 0x07})
	ip8 := net.IP([]byte{0x00, 0x00, 0x00, 0x08})

	topology := NewTopology()

	neighborsTable0 := NewNeighborsTable()
	neighborsTable0.Set(ToNodeID(ip1), &NeighborEntry{Cost: 4})
	neighborsTable0.Set(ToNodeID(ip7), &NeighborEntry{Cost: 8})
	topology.Update(ToNodeID(ip0), neighborsTable0)

	neighborsTable1 := NewNeighborsTable()
	neighborsTable1.Set(ToNodeID(ip2), &NeighborEntry{Cost: 8})
	neighborsTable1.Set(ToNodeID(ip7), &NeighborEntry{Cost: 11})
	topology.Update(ToNodeID(ip1), neighborsTable1)

	neighborsTable2 := NewNeighborsTable()
	neighborsTable2.Set(ToNodeID(ip1), &NeighborEntry{Cost: 8})
	neighborsTable2.Set(ToNodeID(ip3), &NeighborEntry{Cost: 7})
	neighborsTable2.Set(ToNodeID(ip8), &NeighborEntry{Cost: 2})
	topology.Update(ToNodeID(ip2), neighborsTable2)

	neighborsTable3 := NewNeighborsTable()
	neighborsTable3.Set(ToNodeID(ip2), &NeighborEntry{Cost: 7})
	neighborsTable3.Set(ToNodeID(ip4), &NeighborEntry{Cost: 9})
	neighborsTable3.Set(ToNodeID(ip5), &NeighborEntry{Cost: 14})
	topology.Update(ToNodeID(ip3), neighborsTable3)

	neighborsTable4 := NewNeighborsTable()
	neighborsTable4.Set(ToNodeID(ip3), &NeighborEntry{Cost: 9})
	neighborsTable4.Set(ToNodeID(ip5), &NeighborEntry{Cost: 10})
	topology.Update(ToNodeID(ip4), neighborsTable4)

	neighborsTable5 := NewNeighborsTable()
	neighborsTable5.Set(ToNodeID(ip2), &NeighborEntry{Cost: 4})
	neighborsTable5.Set(ToNodeID(ip3), &NeighborEntry{Cost: 14})
	neighborsTable5.Set(ToNodeID(ip4), &NeighborEntry{Cost: 10})
	neighborsTable5.Set(ToNodeID(ip6), &NeighborEntry{Cost: 2})
	topology.Update(ToNodeID(ip5), neighborsTable5)

	neighborsTable6 := NewNeighborsTable()
	neighborsTable6.Set(ToNodeID(ip5), &NeighborEntry{Cost: 2})
	neighborsTable6.Set(ToNodeID(ip7), &NeighborEntry{Cost: 1})
	neighborsTable6.Set(ToNodeID(ip8), &NeighborEntry{Cost: 6})
	topology.Update(ToNodeID(ip6), neighborsTable6)

	neighborsTable7 := NewNeighborsTable()
	neighborsTable7.Set(ToNodeID(ip0), &NeighborEntry{Cost: 8})
	neighborsTable7.Set(ToNodeID(ip1), &NeighborEntry{Cost: 11})
	neighborsTable7.Set(ToNodeID(ip6), &NeighborEntry{Cost: 1})
	neighborsTable7.Set(ToNodeID(ip8), &NeighborEntry{Cost: 7})
	topology.Update(ToNodeID(ip7), neighborsTable7)

	neighborsTable8 := NewNeighborsTable()
	neighborsTable8.Set(ToNodeID(ip2), &NeighborEntry{Cost: 2})
	neighborsTable8.Set(ToNodeID(ip6), &NeighborEntry{Cost: 6})
	neighborsTable8.Set(ToNodeID(ip7), &NeighborEntry{Cost: 7})
	topology.Update(ToNodeID(ip8), neighborsTable8)

	parents := topology.CalculateSinkTree(ToNodeID(ip0))

	for key, value := range parents {
		fmt.Println("dst:", key, "prev:", value)
	}

	fmt.Println("========= Try after removing some vertex ==============")
	topology.g.DeleteVertex(IPv4ToUInt32(ip2))

	parents = topology.CalculateSinkTree(ToNodeID(ip0))

	for key, value := range parents {
		fmt.Println("dst:", key, "prev:", value)
	}

	fmt.Println("========= Try adding the same vertex but with no edges ==============")
	neighborsTable_ := NewNeighborsTable()
	topology.Update(ToNodeID(ip2), neighborsTable_)

	parents = topology.CalculateSinkTree(ToNodeID(ip0))

	for key, value := range parents {
		fmt.Println("dst:", key, "prev:", value)
	}

	fmt.Println("========= Try make this node unreachable from 0 ==============")
	neighborsTable1_ := NewNeighborsTable()
	neighborsTable1.Set(ToNodeID(ip7), &NeighborEntry{Cost: 11})
	topology.Update(ToNodeID(ip1), neighborsTable1_)

	parents = topology.CalculateSinkTree(ToNodeID(ip0))

	for key, value := range parents {
		fmt.Println("dst:", key, "prev:", value)
	}
	//fmt.Println(topology.g.GetVertex(IPv4ToUInt32(ip0)))
}

func TestTopology2(t *testing.T) {
	ip0 := net.IP([]byte{0x00, 0x00, 0x00, 0x00})
	ip1 := net.IP([]byte{0x00, 0x00, 0x00, 0x01})

	topology := NewTopology()

	neighborsTable0 := NewNeighborsTable()
	neighborsTable0.Set(ToNodeID(ip1), &NeighborEntry{Cost: 1})
	topology.Update(ToNodeID(ip0), neighborsTable0)

	neighborsTable1 := NewNeighborsTable()
	neighborsTable1.Set(ToNodeID(ip0), &NeighborEntry{Cost: 1})
	topology.Update(ToNodeID(ip1), neighborsTable1)

	parents0 := topology.CalculateSinkTree(ToNodeID(ip0))
	fmt.Println("Src0")
	for key, value := range parents0 {
		fmt.Println("dst:", key, "prev:", value)
	}

	parents1 := topology.CalculateSinkTree(ToNodeID(ip1))
	fmt.Println("Src1")
	for key, value := range parents1 {
		fmt.Println("dst:", key, "prev:", value)
	}
}
