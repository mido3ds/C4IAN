package dzd

import (
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/tables"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

type DZDController struct {
	reqMacConn                *MACLayerConn
	resMacConn                *MACLayerConn
	dzCache                   *DZCache
	packetsBuffer             *PacketsBuffer
	topology                  *Topology
	myIP                      net.IP
	getUnicastNextHopCallback func(dst NodeID) (net.HardwareAddr, bool)
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

	dzCache := NewDZCache()

	log.Println("initalized dzd controller")

	dzdController := &DZDController{
		reqMacConn: reqMacConn,
		resMacConn: resMacConn,
		dzCache:    dzCache,
		topology:   topology,
		myIP:       ip,
	}

	packetsBuffer := NewPacketsBuffer(dzdController.forwardToNeighborZones)
	dzdController.packetsBuffer = packetsBuffer

	return dzdController, nil
}

func (d *DZDController) SetGetNextHopCallback(getUnicastNextHopCallback func(dst NodeID) (net.HardwareAddr, bool)) {
	d.getUnicastNextHopCallback = getUnicastNextHopCallback
}

func (d *DZDController) CachedDstZone(dstIP net.IP) (ZoneID, bool) {
	return d.dzCache.Get(dstIP)
}

func (d *DZDController) BufferPacket(dstIP net.IP, packet []byte, sendCallback SendPacketCallback) {
	d.packetsBuffer.AppendPacket(dstIP, packet, sendCallback)
}

func (d *DZDController) Start() {
	go d.receiveDZRequestPackets()
	go d.receiveDZResponsePackets()
}

func (d *DZDController) forwardToNeighborZones(dstIP net.IP) {
	// log.Println("Searching for", dstIP)
	neighborsZones, _ := d.topology.GetNeighborZones(ToNodeID(d.myIP))
	for _, zone := range neighborsZones {
		dzRequestPacket := d.createDZRequestPacket(d.myIP, MyZone(), zone, dstIP, []ZoneID{MyZone().ID})
		d.sendDZRequestPackets(dzRequestPacket, zone)
	}
}

func (d *DZDController) FindDstZone(dstIP net.IP) {
	_, inSearch := d.packetsBuffer.Get(dstIP)
	if inSearch {
		return
	}

	d.forwardToNeighborZones(dstIP)
}

func (d *DZDController) receiveDZRequestPackets() {
	for {
		packet := d.reqMacConn.Read()
		go d.handleDZRequestPackets(packet)
	}
}

func (d *DZDController) handleDZRequestPackets(packet []byte) {
	zidHeader, dzRequestHeader := d.unpackDZRequestPacket(packet)

	// Cache the src IP/Zone information
	d.dzCache.Set(dzRequestHeader.srcIP, zidHeader.SrcZID)
	go d.sendBufferedMsgs(dzRequestHeader.srcIP)

	// I'm the required destination
	if dzRequestHeader.requiredDstIP.Equal(d.myIP) {
		d.sendDZResponsePackets(zidHeader, dzRequestHeader, MyZone().ID)
		return
	}

	// The required destination is cached
	/*requiredDstZoneID, exist := d.CachedDstZone(dzRequestHeader.requiredDstIP)
	if exist {
		d.sendDZResponsePackets(zidHeader, dzRequestHeader, requiredDstZoneID)
		return
	}*/

	// Check if this msg is forwarded to my zone
	dstZone := &Zone{ID: zidHeader.DstZID, Len: zidHeader.ZLen}
	if MyZone().Equal(dstZone) {
		//log.Println(d.myIP, "Received ", dzRequestHeader)
		// The msg is forwarded to my zone
		// is the required dstIP exist in my zone ?
		_, inMyZone := d.getUnicastNextHopCallback(ToNodeID(dzRequestHeader.requiredDstIP))
		if inMyZone {
			// It's here, send the response
			d.sendDZResponsePackets(zidHeader, dzRequestHeader, MyZone().ID)
			return
		} else {
			// The required destination isn't here, forward to the neighbors zones
			// but discard the already visited zones
			neighborsZones, _ := d.topology.GetNeighborZones(ToNodeID(d.myIP))
			nextZones := discardVisitedZones(dzRequestHeader.visitedZones, neighborsZones)
			visitedZones := append(dzRequestHeader.visitedZones, MyZone().ID)
			for _, dstZoneID := range nextZones {
				srcZone := Zone{ID: zidHeader.SrcZID, Len: zidHeader.ZLen}
				dzRequestPacket := d.createDZRequestPacket(dzRequestHeader.srcIP, srcZone, dstZoneID, dzRequestHeader.requiredDstIP, visitedZones)
				d.sendDZRequestPackets(dzRequestPacket, dstZoneID)
			}
			return
		}
	} else {
		// Forward it to the right zone in the zidHeader
		d.sendDZRequestPackets(packet, ToNodeID(dstZone.ID))
		return
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

	//log.Println(d.myIP, "Received ", dzResponseHeader)

	d.dzCache.Set(dzResponseHeader.requiredDstIP, dzResponseHeader.requiredDstZone)
	go d.sendBufferedMsgs(dzResponseHeader.requiredDstIP)

	if dzResponseHeader.dstIP.Equal(d.myIP) {
		// log.Println("Discovered zone of:", dzResponseHeader.requiredDstIP, ", zone: ", dzResponseHeader.requiredDstZone)
		return
	}

	var dstNodeID NodeID
	myZone := MyZone()
	dstZID := zidHeader.DstZID.ToLen(myZone.Len)
	if myZone.ID == dstZID {
		dstNodeID = ToNodeID(dzResponseHeader.dstIP)
	} else {
		dstNodeID = ToNodeID(dstZID)
	}

	nextHopMAC, reachable := d.getUnicastNextHopCallback(dstNodeID)
	if !reachable {
		log.Println(dstNodeID, "is unreachable (DZD)")
		return
	}
	//	Forward tha packet as is
	d.resMacConn.Write(packet, nextHopMAC)
}

func (d *DZDController) createDZRequestPacket(srcIP net.IP, srcZone Zone, nextZone NodeID, dstIP net.IP, visitedZones []ZoneID) []byte {
	nextZoneID, valid := nextZone.ToZoneID()
	if !valid {
		log.Panicln("Invalid next zoneID in dzd request packet")
	}
	zidHeader := &ZIDHeader{SrcZID: srcZone.ID, ZLen: srcZone.Len, DstZID: nextZoneID}
	dzRequestHeader := &DZRequestHeader{srcIP: srcIP, requiredDstIP: dstIP, visitedZones: visitedZones}
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

func (d *DZDController) sendDZRequestPackets(dzRequestPacket []byte, zone NodeID) {
	nextHopMAC, reachable := d.getUnicastNextHopCallback(zone)
	if !reachable {
		log.Println(zone, "is unreachable (DZD)")
		return
	}
	d.reqMacConn.Write(dzRequestPacket, nextHopMAC)
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

func (d *DZDController) sendDZResponsePackets(zidHeader *ZIDHeader, dzRequestHeader *DZRequestHeader, requiredDstZoneID ZoneID) {
	srcZone := &Zone{ID: zidHeader.SrcZID, Len: zidHeader.ZLen}
	var srcNodeID NodeID
	if MyZone().Equal(srcZone) {
		srcNodeID = ToNodeID(dzRequestHeader.srcIP)
	} else {
		srcNodeID = ToNodeID(srcZone.ID)
	}
	dzResponsePacket := d.createDZResponsePacket(dzRequestHeader.requiredDstIP, requiredDstZoneID, dzRequestHeader.srcIP, srcZone.ID)
	nextHopMAC, reachable := d.getUnicastNextHopCallback(srcNodeID)
	if !reachable {
		log.Println(srcNodeID, "is unreachable (DZD)")
		return
	}
	d.resMacConn.Write(dzResponsePacket, nextHopMAC)
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
		packet.Send(dstIP)
	}

	d.packetsBuffer.Del(dstIP)
}

func (d *DZDController) Close() {
	d.resMacConn.Close()
	d.reqMacConn.Close()
}
