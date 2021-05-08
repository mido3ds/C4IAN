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

func (f *Forwarder) sendUnicast(packet []byte, destIP net.IP) {
	e, reachable := getUnicastNextHop(destIP, f)

	if !reachable {
		// TODO: Should we do anything else here?
		log.Println("Destination unreachable:", destIP)
		return
	}

	zid := MyZIDHeader(ZoneID(e.DestZoneID))

	// build packet
	buffer := bytes.NewBuffer(make([]byte, 0, f.iface.MTU))
	buffer.Write(f.msec.Encrypt(zid.MarshalBinary()))    // zid
	buffer.Write(f.msec.Encrypt(packet[:IPv4HeaderLen])) // ip header
	buffer.Write(f.msec.Encrypt(packet[IPv4HeaderLen:])) // ip payload

	// write to device driver
	f.zidMacConn.Write(buffer.Bytes(), e.NextHopMAC)
}

func (f *Forwarder) sendMulticast(packet []byte, grpIP net.IP) {
	log.Printf("Node IP:%#v, fwd table: %#v\n", f.ip.String(), f.MultiForwTable.String())
	es, ok := f.MultiForwTable.Get(grpIP)
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
	es, ok = f.MultiForwTable.Get(grpIP)
	log.Println(f.MultiForwTable.String())
	if ok {
		log.Println("Sending....")
		for item := range es.Items.Iter() {
			log.Printf("Send packet to:%#v\n", item.Value.(*NextHopEntry).NextHop.String())
			f.ipMacConn.Write(encryptedPacket, item.Value.(*NextHopEntry).NextHop)
		}
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
