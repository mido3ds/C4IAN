package main

import (
	"log"
	"net"
)

type LSR struct {
	topology *Topology
}

func NewLSR() *LSR {
	t := NewTopology()
	return &LSR{topology: t}
}

func (lsr *LSR) SendLSRPacket(flooder *Flooder, neighborsTable *NeighborsTable) {
	flooder.Flood(neighborsTable.MarshalBinary())
}

func (lsr *LSR) HandleLSRPacket(srcIP net.IP, payload []byte) {
	srcNeighborsTable, valid := UnmarshalNeighborsTable(payload)
	if !valid {
		log.Println("Corrupted LSR packet received")
		return
	}
	lsr.topology.Update(srcIP, srcNeighborsTable)
}

func (lsr *LSR) UpdateForwardingTable(myIP net.IP,
	forwardingTable *ForwardTable,
	neighborsTable *NeighborsTable) {

	forwardingTable.Clear()
	sinkTreeParents := lsr.topology.CalculateSinkTree(myIP)

	for dst, parent := range sinkTreeParents {
		// dst is the same as the src node
		if dst == IPv4ToUInt32(myIP) {
			continue
		}

		dstIP := UInt32ToIPv4(dst.(uint32))

		// dst is one of the src neighbors
		if parent == IPv4ToUInt32(myIP) {
			neighborEntry, exists := neighborsTable.Get(dstIP)
			if !exists {
				log.Panicln("Attempting to make a next hop through a non-neighbor")
			}
			dstMAC := neighborEntry.MAC
			forwardingTable.Set(dstIP, &ForwardingEntry{
				NextHopMAC: dstMAC,
			})
			continue
		}

		// iterate till reaching one of the src neighbors
		// or one of the nodes that we have already known its nextHop

		// TODO: Optimize by collecting nodes along a path
		// and adding next hop for all of them together,
		// then removing them from the map
		var NextHopEntry *ForwardingEntry = nil
		var nextHopEntry *NeighborEntry = nil
		var parentIP, nextHopIP net.IP = nil, nil
		var exist bool = false
		nextHop := parent
		for parent != IPv4ToUInt32(myIP) {
			// check if the dst parent shortest path is calculated before
			parentIP = UInt32ToIPv4(parent.(uint32))
			NextHopEntry, exist = forwardingTable.Get(parentIP)
			if exist {
				break
			}
			// return through the dst shortest path till reach one of the neighbors
			nextHop = parent
			parent = sinkTreeParents[parent]
		}

		if parent == IPv4ToUInt32(myIP) {
			// Get the neighbor MAC using the neighbors table and construct its forwarding entry
			nextHopIP = UInt32ToIPv4(nextHop.(uint32))
			nextHopEntry, exist = neighborsTable.Get(nextHopIP)
			if exist {
				log.Panicln("Attempting to make a next hop through a non-neighbor")
			}
			NextHopEntry = &ForwardingEntry{NextHopMAC: nextHopEntry.MAC}
		}

		forwardingTable.Set(dstIP, NextHopEntry)
	}
}
