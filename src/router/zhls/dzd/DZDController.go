package dzd

import (
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/forward"
	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/tables"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

type DZDController struct {
	reqMacConn *MACLayerConn
	resMacConn *MACLayerConn
	dzCahce    *DZCache
	topology   *Topology
	forwarder  *Forwarder
	myIP       net.IP
}

func NewDZDController(ip net.IP, iface *net.Interface, topology *Topology, forwarder *Forwarder) (*DZDController, error) {
	reqMacConn, err := NewMACLayerConn(iface, DZRequestEtherType)
	if err != nil {
		return nil, err
	}

	resMacConn, err := NewMACLayerConn(iface, DZResponseEtherType)
	if err != nil {
		return nil, err
	}

	dzCahce := NewDZCache()

	log.Println("initalized dzd controller")

	return &DZDController{
		reqMacConn: reqMacConn,
		resMacConn: resMacConn,
		dzCahce:    dzCahce,
		topology:   topology,
		forwarder:  forwarder,
		myIP:       ip,
	}, nil
}

func (d *DZDController) Start() {
	go d.receiveDZRequestPackets()
	go d.receiveDZResponsePackets()
}

func (d *DZDController) CachedDstZone(dstIP net.IP) (ZoneID, bool) {
	return d.dzCahce.Get(dstIP)
}

func (d *DZDController) FindDstZone(dstIP net.IP) {
	neighborsZones := d.topology.GetNeighborsZones(ToNodeID(d.myIP))
	for _, zone := range neighborsZones {
		dzRequestPacket := d.createDZRequestPacket(zone, dstIP, []ZoneID{MyZone().ID})
		nextHopMAC, reachable := d.forwarder.GetUnicastNextHop(zone)
		if !reachable {
			log.Panicln("Neighbor Zone:", zone, "must be reachable")
		}
		d.reqMacConn.Write(dzRequestPacket, nextHopMAC)
	}
}

func (d *DZDController) receiveDZRequestPackets() {
	for {
		packet := d.reqMacConn.Read()
		zidHeader, dzRequestHeader := d.unpackDZRequestPacket(packet)

		requiredDstZoneID, exist := d.CachedDstZone(dzRequestHeader.requiredDstIP)
		if exist {
			srcZone := &Zone{ID: zidHeader.SrcZID, Len: zidHeader.ZLen}
			dzResponsePacket := d.createDZResponsePacket(dzRequestHeader.requiredDstIP, requiredDstZoneID, dzRequestHeader.srcIP, srcZone.ID)
			var dstNodeID NodeID
			if MyZone().Equal(srcZone) {
				dstNodeID = ToNodeID(dzRequestHeader.srcIP)
			} else {
				dstNodeID = ToNodeID(srcZone.ID)
			}
			nextHopMAC, reachable := d.forwarder.GetUnicastNextHop(dstNodeID)
			if !reachable {
				log.Panicln("Neighbor :", dstNodeID, "must be reachable")
			}
			//	Forward tha packet as is
			d.resMacConn.Write(dzResponsePacket, nextHopMAC)
		}

		dstZone := &Zone{ID: zidHeader.DstZID, Len: zidHeader.ZLen}
		if !MyZone().Equal(dstZone) {
			nextHopMAC, reachable := d.forwarder.GetUnicastNextHop(ToNodeID(dstZone.ID))
			if !reachable {
				log.Panicln("Neighbor Zone:", ToNodeID(dstZone.ID), "must be reachable")
			}
			//	Forward tha packet as is
			d.reqMacConn.Write(packet, nextHopMAC)
		} else {
			// Check if the required dstIP exist in my zone
			_, inMyZone := d.forwarder.GetUnicastNextHop(ToNodeID(dzRequestHeader.requiredDstIP))
			if inMyZone {
				srcZone := &Zone{ID: zidHeader.SrcZID, Len: zidHeader.ZLen}
				dzResponsePacket := d.createDZResponsePacket(dzRequestHeader.requiredDstIP, MyZone().ID, dzRequestHeader.srcIP, srcZone.ID)
				var dstNodeID NodeID
				if MyZone().Equal(srcZone) {
					dstNodeID = ToNodeID(dzRequestHeader.srcIP)
				} else {
					dstNodeID = ToNodeID(srcZone.ID)
				}
				nextHopMAC, reachable := d.forwarder.GetUnicastNextHop(dstNodeID)
				if !reachable {
					log.Panicln("Neighbor :", dstNodeID, "must be reachable")
				}
				//	Forward tha packet as is
				d.resMacConn.Write(dzResponsePacket, nextHopMAC)
			} else {
				neighborsZones := d.topology.GetNeighborsZones(ToNodeID(d.myIP))
				nextZones := discardVisitedZones(dzRequestHeader.visitedZones, neighborsZones)
				visitedZones := append(dzRequestHeader.visitedZones, MyZone().ID)
				for _, zone := range nextZones {
					dzRequestPacket := d.createDZRequestPacket(zone, dzRequestHeader.requiredDstIP, visitedZones)
					nextHopMAC, reachable := d.forwarder.GetUnicastNextHop(zone)
					if !reachable {
						log.Panicln("Neighbor Zone:", zone, "must be reachable")
					}
					d.reqMacConn.Write(dzRequestPacket, nextHopMAC)
				}
			}
		}
	}
}

func (d *DZDController) receiveDZResponsePackets() {
	for {
		packet := d.resMacConn.Read()
		zidHeader, dzResponseHeader := d.unpackDZResponsePacket(packet)

		d.dzCahce.Set(dzResponseHeader.requiredDstIP, dzResponseHeader.requiredDstZone)

		dstZone := &Zone{ID: zidHeader.DstZID, Len: zidHeader.ZLen}
		var dstNodeID NodeID
		if MyZone().Equal(dstZone) {
			dstNodeID = ToNodeID(dzResponseHeader.dstIP)
		} else {
			dstNodeID = ToNodeID(dstZone.ID)
		}
		nextHopMAC, reachable := d.forwarder.GetUnicastNextHop(dstNodeID)
		if !reachable {
			log.Panicln("Neighbor :", dstNodeID, "must be reachable")
		}
		//	Forward tha packet as is
		d.resMacConn.Write(packet, nextHopMAC)
	}
}

func (d *DZDController) createDZRequestPacket(nextZone NodeID, dstIP net.IP, visitedZones []ZoneID) []byte {
	nextZoneID, valid := nextZone.ToZoneID()
	if !valid {
		log.Panicln("Invalid next zoneID in dzd request packet")
	}
	zidHeader := MyZIDHeader(nextZoneID)
	dzRequestHeader := &DZRequestHeader{srcIP: d.myIP, requiredDstIP: dstIP, visitedZones: visitedZones}
	return append(zidHeader.MarshalBinary(), dzRequestHeader.MarshalBinary()...)
}

func (d *DZDController) unpackDZRequestPacket(dzRequestPacket []byte) (*ZIDHeader, *DZRequestHeader) {
	zidHeader, valid := UnmarshalZIDHeader(dzRequestPacket[:ZIDHeaderLen])
	if !valid {
		log.Panicln("Received dzd request Packet with invalid ZID header")
	}

	dzRequestHeader, valid := UnmarshalDZRequestHeader(dzRequestPacket[ZIDHeaderLen:])
	if !valid {
		log.Panicln("Received dzd request Packet with invalid dzd request header")
	}

	return zidHeader, dzRequestHeader
}

func (d *DZDController) createDZResponsePacket(requiredDstIP net.IP, requiredDstZone ZoneID, dstIP net.IP, dstZone ZoneID) []byte {
	zidHeader := MyZIDHeader(dstZone)
	dzResponseHeader := &DZResponseHeader{dstIP: dstIP, requiredDstIP: requiredDstIP, requiredDstZone: requiredDstZone}
	return append(zidHeader.MarshalBinary(), dzResponseHeader.MarshalBinary()...)
}

func (d *DZDController) unpackDZResponsePacket(dzResponsePacket []byte) (*ZIDHeader, *DZResponseHeader) {
	zidHeader, valid := UnmarshalZIDHeader(dzResponsePacket[:ZIDHeaderLen])
	if !valid {
		log.Panicln("Received dzd response Packet with invalid ZID header")
	}

	dzResponseHeader, valid := UnmarshalDZResponseHeader(dzResponsePacket[ZIDHeaderLen:])
	if !valid {
		log.Panicln("Received dzd response Packet with invalid dzd response header")
	}

	return zidHeader, dzResponseHeader
}

func discardVisitedZones(visitedZones []ZoneID, neighborsZones []NodeID) (nextZones []NodeID) {
	for _, neighborZone := range neighborsZones {
		visited := false
		neighborZoneID, valid := neighborZone.ToZoneID()
		if !valid {
			log.Panicln("Invalid next zoneID in dzd request packet")
		}
		for _, visitedZoneID := range visitedZones {
			if neighborZoneID == visitedZoneID {
				visited = true
				break
			}
		}
		if !visited {
			nextZones = append(nextZones, neighborZone)
		}
	}
	return
}

func (d *DZDController) Close() {
	d.resMacConn.Close()
	d.reqMacConn.Close()
}
