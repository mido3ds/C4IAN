package zhls

import (
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/flood"
	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/tables"
	"github.com/starwander/goraph"
)

type LSR struct {
	topology *Topology
}

func NewLSR() *LSR {
	t := NewTopology()
	return &LSR{topology: t}
}

func (lsr *LSR) SendLSRPacket(flooder *ZoneFlooder, neighborsTable *NeighborsTable) {
	// log.Println("&&&&&&&&&&&&& Sending LSR")
	// log.Println(neighborsTable)
	flooder.Flood(neighborsTable.MarshalBinary())
}

func (lsr *LSR) HandleLSRPacket(srcIP net.IP, payload []byte) {
	srcNeighborsTable, valid := UnmarshalNeighborsTable(payload)

	if !valid {
		log.Panicln("Corrupted LSR packet received")
	}
	// log.Println("//////////////Received LSR packet from: ", srcIP)
	// log.Println(srcNeighborsTable)
	lsr.topology.Update(srcIP, srcNeighborsTable)
}

func (lsr *LSR) UpdateForwardingTable(myIP net.IP,
	forwardingTable *UniForwardTable,
	neighborsTable *NeighborsTable) {

	dirtyForwardingTable := NewUniForwardTable()
	sinkTreeParents := lsr.topology.CalculateSinkTree(myIP)

	// log.Println(neighborsTable)
	// lsr.displaySinkTreeParents(sinkTreeParents)

	for dst, parent := range sinkTreeParents {
		dstIP := UInt32ToIPv4(dst.(uint32))

		// Dst is the same as the src node
		if dstIP.Equal(myIP) {
			continue
		}

		// Dst in unreachable
		// TODO : to be handled
		if parent == nil {
			log.Println(UInt32ToIPv4(dst.(uint32)), "is unreachable")
			continue
		}

		// Dst is a direct neighbor
		var nextHop goraph.ID
		if parent == IPv4ToUInt32(myIP) {
			nextHop = dst
		}

		// TODO: Optimize by collecting nodes along a path
		// and adding next hop for all of them together,
		// then removing them from the map

		// Iterate till reaching one of the direct neighbors
		// or one of the nodes that we have already known its nextHop
		for parent != IPv4ToUInt32(myIP) {
			// check if the dst parent shortest path is calculated before
			parentIP := UInt32ToIPv4(parent.(uint32))
			forwardingEntry, exists := dirtyForwardingTable.Get(parentIP)
			if exists {
				dirtyForwardingTable.Set(dstIP, forwardingEntry)
				break
			}
			// return through the dst shortest path till reach one of the neighbors
			nextHop = parent
			parent = sinkTreeParents[parent]
		}

		// We iterated through the path until we reached a direct neighbor
		if parent == IPv4ToUInt32(myIP) {
			// Get the neighbor MAC using the neighbors table and construct its forwarding entry
			nextHopIP := UInt32ToIPv4(nextHop.(uint32))
			neighborEntry, exists := neighborsTable.Get(nextHopIP)
			if !exists {
				log.Panicln("Attempting to make a next hop through a non-neighbor")
			}
			dirtyForwardingTable.Set(dstIP, &UniForwardingEntry{NextHopMAC: neighborEntry.MAC})
		}
		// Shallow copy the forwarding table, this will make the hashmap pointer in forwardingTable
		// point to the new hashmap inside dityForwardingTable. The old hashmap in forwardingTable
		// will be deleted by the garbage collector
		*forwardingTable = *dirtyForwardingTable
	}
}

func (lsr *LSR) displaySinkTreeParents(sinkTreeParents map[goraph.ID]goraph.ID) {
	log.Println("----------- Sink Tree -------------")
	for dst, parent := range sinkTreeParents {
		if dst == nil {
			log.Println("Dst: ", dst, "Parent: ", UInt32ToIPv4(parent.(uint32)))
			continue
		}
		if parent == nil {
			log.Println("Dst: ", UInt32ToIPv4(dst.(uint32)), "Parent: ", parent)
			lsr.topology.DisplayVertex(dst.(uint32))
			continue
		}
		log.Println("Dst: ", UInt32ToIPv4(dst.(uint32)), "Parent: ", UInt32ToIPv4(parent.(uint32)))
		lsr.topology.DisplayVertex(dst.(uint32))
	}
	log.Println("-----------------------------------")
}
