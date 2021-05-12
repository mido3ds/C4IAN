package dzd

import (
	"fmt"
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/tables"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

type DZDController struct {
	reqMacConn                *MACLayerConn
	resMacConn                *MACLayerConn
	dzCahce                   *DZCache
	packetsBuffer             *PacketsBuffer
	topology                  *Topology
	myIP                      net.IP
	getUnicastNextHopCallback func(dst NodeID) (net.HardwareAddr, bool)
	sendUnicastCallback       func(packet []byte, dstIP net.IP)
}

func NewDZDController(ip net.IP, iface *net.Interface, topology *Topology) (*DZDController, error) {
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

	dzdController := &DZDController{
		reqMacConn: reqMacConn,
		resMacConn: resMacConn,
		dzCahce:    dzCahce,
		topology:   topology,
		myIP:       ip,
	}

	packetsBuffer := NewPacketsBuffer(dzdController.FindDstZone)
	dzdController.packetsBuffer = packetsBuffer

	return dzdController, nil
}

func (d *DZDController) SetForwarderCallbacks(
	getUnicastNextHopCallback func(dst NodeID) (net.HardwareAddr, bool),
	sendUnicastCallback func(packet []byte, dstIP net.IP)) {
	d.getUnicastNextHopCallback = getUnicastNextHopCallback
	d.sendUnicastCallback = sendUnicastCallback
}

func (d *DZDController) Start() {
	go d.receiveDZRequestPackets()
	go d.receiveDZResponsePackets()
}

func (d *DZDController) CachedDstZone(dstIP net.IP) (ZoneID, bool) {
	return d.dzCahce.Get(dstIP)
}

func (d *DZDController) BufferPacket(dstIP net.IP, packet []byte) {
	d.packetsBuffer.AppendPacket(dstIP, packet)
}

func (d *DZDController) FindDstZone(dstIP net.IP) {
	_, inSearch := d.packetsBuffer.Get(dstIP)
	if inSearch {
		return
	}

	neighborsZones := d.topology.GetNeighborsZones(ToNodeID(d.myIP))
	for _, zone := range neighborsZones {
		dzRequestPacket := d.createDZRequestPacket(d.myIP, MyZone().ID, zone, dstIP, []ZoneID{MyZone().ID})
		nextHopMAC, reachable := d.getUnicastNextHopCallback(zone)
		if !reachable {
			log.Panicln(zone, "is unreachable")
		}
		d.reqMacConn.Write(dzRequestPacket, nextHopMAC)
	}
}

func (d *DZDController) receiveDZRequestPackets() {
	for {
		packet := d.reqMacConn.Read()
		go d.handleDZRequestPackets(packet)
	}
}

func (d *DZDController) handleDZRequestPackets(packet []byte) {
	zidHeader, dzRequestHeader := d.unpackDZRequestPacket(packet)
	requiredDstZoneID, exist := d.CachedDstZone(dzRequestHeader.requiredDstIP)
	if exist {
		srcZone := &Zone{ID: dzRequestHeader.srcZone, Len: zidHeader.ZLen}
		dzResponsePacket := d.createDZResponsePacket(dzRequestHeader.requiredDstIP, requiredDstZoneID, dzRequestHeader.srcIP, dzRequestHeader.srcZone)
		var dstNodeID NodeID
		if MyZone().Equal(srcZone) {
			dstNodeID = ToNodeID(dzRequestHeader.srcIP)
		} else {
			dstNodeID = ToNodeID(srcZone.ID)
		}
		nextHopMAC, reachable := d.getUnicastNextHopCallback(dstNodeID)
		if !reachable {
			log.Panicln(dstNodeID, "is unreachable")
		}
		//	Forward tha packet as is
		d.resMacConn.Write(dzResponsePacket, nextHopMAC)
		return
	}

	dstZone := &Zone{ID: zidHeader.DstZID, Len: zidHeader.ZLen}
	if !MyZone().Equal(dstZone) {
		nextHopMAC, reachable := d.getUnicastNextHopCallback(ToNodeID(dstZone.ID))
		if !reachable {
			log.Panicln(ToNodeID(dstZone.ID), "is unreachable")
		}
		//	Forward tha packet as is
		d.reqMacConn.Write(packet, nextHopMAC)
		return
	} else {
		fmt.Println(d.myIP, "Received ", dzRequestHeader)
		// Check if the required dstIP exist in my zone
		_, inMyZone := d.getUnicastNextHopCallback(ToNodeID(dzRequestHeader.requiredDstIP))
		if inMyZone {
			srcZone := &Zone{ID: dzRequestHeader.srcZone, Len: zidHeader.ZLen}
			dzResponsePacket := d.createDZResponsePacket(dzRequestHeader.requiredDstIP, MyZone().ID, dzRequestHeader.srcIP, dzRequestHeader.srcZone)
			var dstNodeID NodeID
			if MyZone().Equal(srcZone) {
				dstNodeID = ToNodeID(dzRequestHeader.srcIP)
			} else {
				dstNodeID = ToNodeID(srcZone.ID)
			}
			nextHopMAC, reachable := d.getUnicastNextHopCallback(dstNodeID)
			if !reachable {
				log.Panicln("Neighbor :", dstNodeID, "is unreachable")
			}
			//	Forward tha packet as is
			d.resMacConn.Write(dzResponsePacket, nextHopMAC)
			return
		} else {
			neighborsZones := d.topology.GetNeighborsZones(ToNodeID(d.myIP))
			nextZones := discardVisitedZones(dzRequestHeader.visitedZones, neighborsZones)
			visitedZones := append(dzRequestHeader.visitedZones, MyZone().ID)
			for _, zone := range nextZones {
				dzRequestPacket := d.createDZRequestPacket(dzRequestHeader.srcIP, dzRequestHeader.srcZone, zone, dzRequestHeader.requiredDstIP, visitedZones)
				nextHopMAC, reachable := d.getUnicastNextHopCallback(zone)
				if !reachable {
					log.Panicln(zone, "is unreachable")
				}
				d.reqMacConn.Write(dzRequestPacket, nextHopMAC)
			}
			return
		}
	}

}

func (d *DZDController) receiveDZResponsePackets() {
	for {
		packet := d.resMacConn.Read()
		go d.handleDZResponsePackets(packet)
	}
}

func (d *DZDController) handleDZResponsePackets(packet []byte) {
	zidHeader, dzResponseHeader := d.unpackDZResponsePacket(packet)

	fmt.Println(d.myIP, "Received ", dzResponseHeader)

	d.dzCahce.Set(dzResponseHeader.requiredDstIP, dzResponseHeader.requiredDstZone)
	go d.sendBufferedMsgs(dzResponseHeader.requiredDstIP)

	if dzResponseHeader.dstIP.Equal(d.myIP) {
		return
	}

	dstZone := &Zone{ID: zidHeader.DstZID, Len: zidHeader.ZLen}
	var dstNodeID NodeID
	if MyZone().Equal(dstZone) {
		dstNodeID = ToNodeID(dzResponseHeader.dstIP)
	} else {
		dstNodeID = ToNodeID(dstZone.ID)
	}
	nextHopMAC, reachable := d.getUnicastNextHopCallback(dstNodeID)
	if !reachable {
		log.Panicln(dstNodeID, "is unreachable")
	}
	//	Forward tha packet as is
	d.resMacConn.Write(packet, nextHopMAC)
}

func (d *DZDController) createDZRequestPacket(srcIP net.IP, srcZone ZoneID, nextZone NodeID, dstIP net.IP, visitedZones []ZoneID) []byte {
	nextZoneID, valid := nextZone.ToZoneID()
	if !valid {
		log.Panicln("Invalid next zoneID in dzd request packet")
	}
	zidHeader := MyZIDHeader(nextZoneID)
	dzRequestHeader := &DZRequestHeader{srcIP: srcIP, srcZone: srcZone, requiredDstIP: dstIP, visitedZones: visitedZones}
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

func (d *DZDController) sendBufferedMsgs(dstIP net.IP) {
	bufferQueue, exist := d.packetsBuffer.Get(dstIP)
	if !exist {
		return
	}

	for _, packet := range bufferQueue {
		d.sendUnicastCallback(packet, dstIP)
	}

	d.packetsBuffer.Del(dstIP)
}

func (d *DZDController) Close() {
	d.resMacConn.Close()
	d.reqMacConn.Close()
}
