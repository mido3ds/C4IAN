package zhls

import (
	"log"
	"net"
	"time"

	. "github.com/mido3ds/C4IAN/src/router/constants"
	. "github.com/mido3ds/C4IAN/src/router/flood"
	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/tables"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
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

func newLSR(iface *net.Interface, msec *MSecLayer, myIP net.IP, neighborsTable *NeighborsTable, t *Topology) *LSRController {
	zoneFlooder := NewZoneFlooder(iface, myIP, msec)
	globalFlooder := NewGlobalFlooder(myIP, iface, InterzoneLSREtherType, msec)

	return &LSRController{
		myIP:           myIP,
		neighborsTable: neighborsTable,
		topology:       t,
		zoneFlooder:    zoneFlooder,
		globalFlooder:  globalFlooder,
	}
}

func (lsr *LSRController) Start() {
	go lsr.zoneFlooder.ListenForFloodedMsgs(lsr.handleIntrazoneLSRPacket)
	go lsr.globalFlooder.ListenForFloodedMsgs(lsr.handleInterzoneLSRPacket)
	go lsr.sendInterzoneLSR()
}

func (lsr *LSRController) OnZoneChange(newZoneID ZoneID) {
	lsr.topology.Clear()
}

func (lsr *LSRController) Close() {
	lsr.zoneFlooder.Close()
	lsr.globalFlooder.Close()
}

func (lsr *LSRController) sendInterzoneLSR() {
	for {
		time.Sleep(InterzoneLSRDelay)
		// Get a list of neighbor zones
		neighborZones, isMaxIP := lsr.topology.GetNeighborZones(ToNodeID(lsr.myIP))

		// Only the node with the maximum IP value in a zone floods the zone LSR
		if !isMaxIP {
			continue
		}

		// Create a zone neighbors table
		zoneNeighborsTable := NewNeighborsTable()
		for _, neighborZoneID := range neighborZones {
			zoneNeighborsTable.Set(neighborZoneID, &NeighborEntry{Cost: 65535}) // Cost = MAX_UINT16
		}

		zidHeader := MyZIDHeader(0)
		//log.Println("Sending Interzone LSR Packet from zone: ", zidHeader.SrcZID)
		lsr.globalFlooder.Flood(append(zidHeader.MarshalBinary(), zoneNeighborsTable.MarshalBinary()...))
	}
}

func (lsr *LSRController) sendIntrazoneLSR(isUpdated bool) {
	// if isUpdated {
	// 	log.Println(lsr.neighborsTable)
	// }
	lsr.topology.Update(ToNodeID(lsr.myIP), lsr.neighborsTable)
	lsr.zoneFlooder.Flood(lsr.neighborsTable.MarshalBinary())
}

func (lsr *LSRController) handleIntrazoneLSRPacket(srcIP net.IP, payload []byte) {
	srcNeighborsTable, valid := UnmarshalNeighborsTable(payload)
	if !valid {
		log.Panicln("Corrupted neighbors table in intrazone LSR packet received")
	}

	lsr.topology.Update(ToNodeID(srcIP), srcNeighborsTable)
	lsr.dirtyTopology = true
}

func (lsr *LSRController) handleInterzoneLSRPacket(payload []byte) []byte {
	zidHeader, valid := UnmarshalZIDHeader(payload[:ZIDHeaderLen])
	if !valid {
		log.Panicln("Corrupted ZID header in interzone LSR packet received")
	}

	zoneNeighborsTable, valid := UnmarshalNeighborsTable(payload[ZIDHeaderLen:])
	if !valid {
		log.Panicln("Corrupted neighbors table in interzone LSR packet received")
	}

	//log.Println("Received Interzone LSR Packet from zone:", zidHeader.SrcZID)
	//log.Println(zoneNeighborsTable)
	//lsr.displaySinkTreeParents(lsr.topology.CalculateSinkTree(ToNodeID(lsr.myIP)))

	lsr.topology.Update(ToNodeID(zidHeader.SrcZID.ToLen(MyZone().Len)), zoneNeighborsTable)
	lsr.dirtyTopology = true
	return payload
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
		// It will eventually be removed from the topology when its timer fires
		if parent == nil {
			log.Println(dst, "is unreachable (LSR)")
			continue
		}

		// Dst is a direct neighbor
		var nextHop goraph.ID
		if parent == ToNodeID(lsr.myIP) {
			nextHop = dst
		}

		// Iterate till reaching one of the direct neighbors
		// or one of the nodes that we have already known its nextHop
		// TODO (low priority):
		//		Optimize by collecting nodes along a path and adding next hop for all of them together,
		// 		then removing them from the map
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
				//log.Println(lsr.neighborsTable)
				//lsr.topology.DisplaySinkTreeParents(sinkTreeParents)
				log.Panicln("Attempting to make a next hop through a non-neighbor, dst: ", nextHop.(NodeID))
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
