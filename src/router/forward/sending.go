package forward

import (
	"bytes"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

func (f *Forwarder) sendUnicast(packet []byte, destIP net.IP) {
	e, ok := getNextHop(destIP, f.UniForwTable, f.neighborsTable, f.zoneID)
	if !ok {
		// TODO: call controller
		return
	}

	zid := &ZIDHeader{ZLen: f.zlen, PacketType: DataPacket, SrcZID: f.zoneID, DstZID: ZoneID(e.DestZoneID)}

	// build packet
	buffer := bytes.NewBuffer(make([]byte, 0, f.iface.MTU))
	buffer.Write(zid.MarshalBinary())
	buffer.Write(packet)

	// encrypt
	encryptedPacket := f.msec.Encrypt(buffer.Bytes())

	// write to device driver
	f.zidMacConn.Write(encryptedPacket, e.NextHopMAC)
}

func (f *Forwarder) sendMulticast(packet []byte, grpIP net.IP) {
	es, ok := f.MultiForwTable.Get(grpIP)
	if !ok {
		es, ok = f.mcGetMissingEntries(grpIP)
		if !ok {
			return
		}
	}

	// encrypt
	encryptedPacket := f.msec.Encrypt(packet)

	// write to device driver
	for i := 0; i < len(es.NextHopMACs); i++ {
		f.zidMacConn.Write(encryptedPacket, es.NextHopMACs[i])
	}
}

func (f *Forwarder) sendBroadcast(packet []byte) {
	// write to device driver
	// TODO: for now ethernet broadcast
	f.zidMacConn.Write(f.msec.Encrypt(packet), BroadcastMACAddr)
}
