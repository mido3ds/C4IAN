package forward

import (
	"bytes"
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/tables"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

func (f *Forwarder) SendUnicast(packet []byte, dstIP net.IP) {
	// Get the next hop using the dstIP
	// return true only if the dst inside my zone
	nextHopMac, inMyZone := f.GetUnicastNextHop(ToNodeID(dstIP))

	var zid *ZIDHeader
	if inMyZone {
		zid = MyZIDHeader(MyZone().ID)
	} else {
		// Check if this dst zone is cached
		dstZoneID, cached := f.dzdController.CachedDstZone(dstIP)
		if cached {
			zid = MyZIDHeader(dstZoneID)
			nextHopMac, _ = f.GetUnicastNextHop(ToNodeID(dstZoneID))
		} else {
			// if dst zone isn't cached, search for it
			// and buffer this msg to be sent when dst zone response arrive
			f.dzdController.FindDstZone(dstIP)
			f.dzdController.BufferPacket(dstIP, packet)
			return
		}
	}

	// build packet
	buffer := bytes.NewBuffer(make([]byte, 0, f.iface.MTU))
	buffer.Write(f.msec.Encrypt(zid.MarshalBinary()))    // zid
	buffer.Write(f.msec.Encrypt(packet[:IPv4HeaderLen])) // ip header
	buffer.Write(f.msec.Encrypt(packet[IPv4HeaderLen:])) // ip payload

	// write to device driver
	log.Println("Sending to: ", nextHopMac)
	f.zidMacConn.Write(buffer.Bytes(), nextHopMac)
}

func (f *Forwarder) sendMulticast(packet []byte, grpIP net.IP) {
	es, ok := f.MultiForwTable.Get(grpIP)
	if !ok {
		es, ok = f.mcGetMissingEntries(grpIP)
		if !ok {
			return
		}
	}

	// build packet
	buffer := bytes.NewBuffer(make([]byte, 0, f.iface.MTU))
	buffer.Write(f.msec.Encrypt(packet[:IPv4HeaderLen])) // ip header
	buffer.Write(f.msec.Encrypt(packet[IPv4HeaderLen:])) // ip payload
	encryptedPacket := buffer.Bytes()

	// write to device driver
	for i := 0; i < len(es.NextHopMACs); i++ {
		f.ipMacConn.Write(encryptedPacket, es.NextHopMACs[i])
	}
}

func (f *Forwarder) sendBroadcast(packet []byte) {
	// build packet
	buffer := bytes.NewBuffer(make([]byte, 0, f.iface.MTU))
	buffer.Write(f.msec.Encrypt(packet[:IPv4HeaderLen])) // ip header
	buffer.Write(f.msec.Encrypt(packet[IPv4HeaderLen:])) // ip payload

	// write to device driver
	// TODO: for now ethernet broadcast
	f.zidMacConn.Write(buffer.Bytes(), BroadcastMACAddr)
}
