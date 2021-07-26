package forward

import (
	"bytes"
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/database_logger"
	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/tables"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

func (f *Forwarder) sendUnicast(packet []byte, dstIP net.IP) {
	// Get the next hop using the dstIP
	// return true only if the dst inside my zone
	nextHopMac, inMyZone := f.GetUnicastNextHop(ToNodeID(dstIP))

	//log.Println(f.ip, "want to send unicast message to", dstIP)

	var zid *ZIDHeader
	if inMyZone {
		//log.Println(dstIP, "is in", f.ip, "zone")
		zid = MyZIDHeader(MyZone().ID)
	} else {
		// Check if this dst zone is cached
		//log.Println(dstIP, "isn't in", f.ip, "zone")
		dstZoneID, cached := f.dzdController.CachedDstZone(dstIP)
		if cached {
			zid = MyZIDHeader(dstZoneID)
			var reachable bool
			//log.Println(dstIP, "zone is cached")
			nextHopMac, reachable = f.GetUnicastNextHop(ToNodeID(dstZoneID))
			if !reachable {
				// If dst zone is cached but unreachable, it may have moved to a reachable zone -> rediscover
				//log.Println(dstIP, "cached zone isn't reachable")
				f.dzdController.FindDstZone(dstIP)
				f.dzdController.BufferPacket(dstIP, packet, f.sendUnicast)
				return
			}
			//log.Println(dstIP, "cached zone is reachable")
		} else {
			// if dst zone isn't cached, search for it
			// and buffer this msg to be sent when dst zone response arrive
			//log.Println(dstIP, "zone isn't cached")
			f.dzdController.FindDstZone(dstIP)
			f.dzdController.BufferPacket(dstIP, packet, f.sendUnicast)
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
	DatabaseLogger.LogForwarding(buffer.Bytes()[ZIDHeaderLen+IPv4HeaderLen:], dstIP)
}

// sendMulticast takes a packet the router wants to send and the multicast group ip to send the packet to
func (f *Forwarder) sendMulticast(packet []byte, grpIP net.IP) {
	// log.Printf("Node IP:%#v, fwd table: %#v\n", f.ip.String(), f.MultiForwTable.String())
	_, ok := f.MultiForwTable.Get(grpIP)
	if !ok {
		// get missing entries from the multi forward table
		ok = f.mcGetMissingEntries(grpIP)
		if !ok {
			log.Println("error")
			return
		}
		// time.Sleep(1 * time.Second)
	}

	// build packet
	buffer := bytes.NewBuffer(make([]byte, 0, f.iface.MTU))
	buffer.Write(f.msec.Encrypt(packet[:IPv4HeaderLen])) // ip header
	buffer.Write(f.msec.Encrypt(packet[IPv4HeaderLen:])) // ip payload
	encryptedPacket := buffer.Bytes()

	// write to device driver
	es, exist := f.MultiForwTable.Get(grpIP)
	// log.Println(f.MultiForwTable.String())
	if exist {
		for item := range es.Items.Iter() {
			// log.Printf("Send packet to:%#v\n", item.Value.(*NextHopEntry).NextHop.String())
			f.ipMacConn.Write(encryptedPacket, item.Value.(*NextHopEntry).NextHop)
		}
	}
}

func (f *Forwarder) sendBroadcast(packet []byte) {
	zid := MyZIDHeader(ZoneID(0))

	// build packet
	buffer := bytes.NewBuffer(make([]byte, 0, f.iface.MTU))
	buffer.Write(f.msec.Encrypt(zid.MarshalBinary()))    // zid
	buffer.Write(f.msec.Encrypt(packet[:IPv4HeaderLen])) // ip header
	buffer.Write(f.msec.Encrypt(packet[IPv4HeaderLen:])) // ip payload

	// flood packet
	f.bcFlooder.Flood(buffer.Bytes())
}
