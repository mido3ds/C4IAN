package forward

import (
	"bytes"
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

func (f *Forwarder) sendUnicast(packet []byte, destIP net.IP) {
	e, reachable := getUnicastNextHop(destIP, f)

	if !reachable {
		// TODO: Should we do anything else here?
		log.Println("Destination unreachable:", destIP)
		return
	}

	zid := &ZIDHeader{ZLen: f.zlen, SrcZID: f.zoneID, DstZID: ZoneID(e.DestZoneID)}

	// build packet
	buffer := bytes.NewBuffer(make([]byte, 0, f.iface.MTU))
	buffer.Write(f.msec.Encrypt(zid.MarshalBinary()))    // zid
	buffer.Write(f.msec.Encrypt(packet[:IPv4HeaderLen])) // ip header
	buffer.Write(f.msec.Encrypt(packet[IPv4HeaderLen:])) // ip payload

	// write to device driver
	f.zidMacConn.Write(buffer.Bytes(), e.NextHopMAC)
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
		f.zidMacConn.Write(encryptedPacket, es.NextHopMACs[i])
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
