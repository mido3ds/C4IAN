package zhls

import (
	"fmt"
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/flood"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/tables"
	"github.com/starwander/goraph"
)

type LSRController struct {
	myIP           net.IP
	neighborsTable *NeighborsTable
	topology       *Topology
	zoneFlooder    *ZoneFlooder
	globalFlooder  *GlobalFlooder
	dirtyTopology  bool
}

func newLSR(iface *net.Interface, msec *MSecLayer, myIP net.IP, neighborsTable *NeighborsTable, t *Topology) (*LSRController, error) {
	zoneFlooder, err := NewZoneFlooder(iface, myIP, msec)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate flooder, err: %#v", err)
	}

	return &LSRController{myIP: myIP, neighborsTable: neighborsTable, topology: t, zoneFlooder: zoneFlooder}, nil
}

func (lsr *LSRController) Start() {
	go lsr.zoneFlooder.ListenForFloodedMsgs(lsr.handleLSRPacket)
}

func (lsr *LSRController) Close() {
	lsr.zoneFlooder.Close()
}

func (lsr *LSRController) onNeighborhoodUpdate() {
	lsr.topology.Update(ToNodeID(lsr.myIP), lsr.neighborsTable)
	lsr.zoneFlooder.Flood(lsr.neighborsTable.MarshalBinary())
}

func (lsr *LSRController) handleLSRPacket(srcIP net.IP, payload []byte) {
	srcNeighborsTable, valid := UnmarshalNeighborsTable(payload)

	if !valid {
		log.Panicln("Corrupted LSR packet received")
	}
	lsr.topology.Update(ToNodeID(srcIP), srcNeighborsTable)
	lsr.dirtyTopology = true
}

func (lsr *LSRController) updateForwardingTable(forwardingTable *UniForwardTable) {
	if !lsr.dirtyTopology {
		return
	}

	dirtyForwardingTable := NewUniForwardTable()
	sinkTreeParents := lsr.topology.CalculateSinkTree(ToNodeID(lsr.myIP))
	//lsr.displaySinkTreeParents(sinkTreeParents)

	for dst, parent := range sinkTreeParents {

		// Dst is the same as the src node
		if dst == ToNodeID(lsr.myIP) {
			continue
		}

		// Dst in unreachable
		// TODO : to be handled
		if parent == nil {
			log.Println(dst, "is unreachable")
			continue
		}

		// Dst is a direct neighbor
		var nextHop goraph.ID
		if parent == ToNodeID(lsr.myIP) {
			nextHop = dst
		}

		// Iterate till reaching one of the direct neighbors
		// or one of the nodes that we have already known its nextHop
		// TODO: Optimize by collecting nodes along a path and adding next hop for all of them together,
		// 		 then removing them from the map
		for parent != ToNodeID(lsr.myIP) {
			// check if the dst parent shortest path is calculated before
			forwardingEntry, exists := dirtyForwardingTable.Get(parent.(NodeID))
			if exists {
				dirtyForwardingTable.Set(dst.(NodeID), forwardingEntry)
				break
			}
			// return through the dst shortest path till reach one of the neighbors
			nextHop = parent
			parent = sinkTreeParents[parent]
		}

		// We iterated through the path until we reached a direct neighbor
		if parent == ToNodeID(lsr.myIP) {
			// Get the neighbor MAC using the neighbors table and construct its forwarding entry
			neighborEntry, exists := lsr.neighborsTable.Get(nextHop.(NodeID))
			if !exists {
				log.Panicln("Attempting to make a next hop through a non-neighbor")
			}
			dirtyForwardingTable.Set(dst.(NodeID), neighborEntry.MAC)
		}
	}
	// Shallow copy the forwarding table, this will make the hashmap pointer in forwardingTable
	// point to the new hashmap inside dityForwardingTable. The old hashmap in forwardingTable
	// will be deleted by the garbage collector
	*forwardingTable = *dirtyForwardingTable
	lsr.dirtyTopology = false
}

func (lsr *LSRController) displaySinkTreeParents(sinkTreeParents map[goraph.ID]goraph.ID) {
	log.Println("----------- Sink Tree -------------")
	for dst, parent := range sinkTreeParents {
		if dst == nil {
			log.Println("Dst: ", dst, "Parent: ", parent.(NodeID))
			continue
		}
		if parent == nil {
			log.Println("Dst: ", dst.(NodeID), "Parent: ", parent)
			lsr.topology.DisplayVertex(dst.(NodeID))
			continue
		}
		log.Println("Dst: ", dst.(NodeID), "Parent: ", parent.(NodeID))
		lsr.topology.DisplayVertex(dst.(NodeID))
	}
	log.Println("-----------------------------------")
}
