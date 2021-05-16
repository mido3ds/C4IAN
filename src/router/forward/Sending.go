package forward

import (
	"bytes"
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/tables"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

func (f *Forwarder) SendUnicast(packet []byte, dstIP net.IP) {
	// Get the next hop using the dstIP
	// return true only if the dst inside my zone
	nextHopMac, inMyZone := f.GetUnicastNextHop(ToNodeID(dstIP))

	//log.Println(f.ip , "want to send unicast message to", dstIP)

	var zid *ZIDHeader
	if inMyZone {
		zid = MyZIDHeader(MyZone().ID)
	} else {
		// Check if this dst zone is cached
		dstZoneID, cached := f.dzdController.CachedDstZone(dstIP)
		if cached {
			zid = MyZIDHeader(dstZoneID)
			var reachable bool
			nextHopMac, reachable = f.GetUnicastNextHop(ToNodeID(dstZoneID))
			if !reachable {
				// If dst zone is cached but unreachable, it may have moved to a reachable zone -> rediscover
				f.dzdController.FindDstZone(dstIP)
				f.dzdController.BufferPacket(dstIP, packet, f.SendUnicast)
				return
			}
			log.Println("Sending a msg to: ", dstIP, " in zone: ", dstZoneID, "through: ", nextHopMac)
		} else {
			// if dst zone isn't cached, search for it
			// and buffer this msg to be sent when dst zone response arrive
			f.dzdController.FindDstZone(dstIP)
			f.dzdController.BufferPacket(dstIP, packet, f.SendUnicast)
			return
		}
	}

	// build packet
	buffer := bytes.NewBuffer(make([]byte, 0, f.iface.MTU))
	buffer.Write(f.msec.Encrypt(zid.MarshalBinary()))    // zid
	buffer.Write(f.msec.Encrypt(packet[:IPv4HeaderLen])) // ip header
	buffer.Write(f.msec.Encrypt(packet[IPv4HeaderLen:])) // ip payload

	// write to device driver
	f.zidMacConn.Write(buffer.Bytes(), nextHopMac)
}

func (f *Forwarder) SendMulticast(packet []byte, grpIP net.IP) {
	log.Printf("Node IP:%#v, fwd table: %#v\n", f.ip.String(), f.MultiForwTable.String())
	_, ok := f.MultiForwTable.Get(grpIP)
	if !ok {
		ok = f.mcGetMissingEntries(grpIP)
		if !ok {
			log.Println("error")
			return
		}
	}

	// build packet
	buffer := bytes.NewBuffer(make([]byte, 0, f.iface.MTU))
	buffer.Write(f.msec.Encrypt(packet[:IPv4HeaderLen])) // ip header
	buffer.Write(f.msec.Encrypt(packet[IPv4HeaderLen:])) // ip payload
	encryptedPacket := buffer.Bytes()

	// write to device driver
	es, exist := f.MultiForwTable.Get(grpIP)
	log.Println(f.MultiForwTable.String())
	if exist {
		for item := range es.Items.Iter() {
			log.Printf("Send packet to:%#v\n", item.Value.(*NextHopEntry).NextHop.String())
			f.ipMacConn.Write(encryptedPacket, item.Value.(*NextHopEntry).NextHop)
		}
	}
}

func (f *Forwarder) SendBroadcast(packet []byte) {
	zid := MyZIDHeader(ZoneID(0))

	// build packet
	buffer := bytes.NewBuffer(make([]byte, 0, f.iface.MTU))
	buffer.Write(f.msec.Encrypt(zid.MarshalBinary()))    // zid
	buffer.Write(f.msec.Encrypt(packet[:IPv4HeaderLen])) // ip header
	buffer.Write(f.msec.Encrypt(packet[IPv4HeaderLen:])) // ip payload

	// flood packet
	f.bcFlooder.Flood(buffer.Bytes())
}
