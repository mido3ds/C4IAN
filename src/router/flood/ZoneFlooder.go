package flood

import (
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/tables"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

type ZoneFlooder struct {
	seqNumber uint32
	fTable    *FloodingTable
	macConn   *MACLayerConn
	ip        net.IP
	zoneID    ZoneID
	msec      *MSecLayer
	zlen      byte
}

func NewZoneFlooder(iface *net.Interface, ip net.IP, msec *MSecLayer, zlen byte) (*ZoneFlooder, error) {
	// connect to mac layer
	macConn, err := NewMACLayerConn(iface, ZIDEtherType)
	if err != nil {
		return nil, err
	}

	fTable := NewFloodingTable()

	log.Println("initalized zone flooder")

	return &ZoneFlooder{
		seqNumber: 0,
		fTable:    fTable,
		macConn:   macConn,
		ip:        ip,
		msec:      msec,
		zlen:      zlen,
	}, nil
}

func (flooder *ZoneFlooder) Flood(msg []byte) {
	hdr := FloodHeader{SrcIP: flooder.ip, SeqNum: flooder.seqNumber}
	msg = append(hdr.MarshalBinary(), msg...)

	flooder.seqNumber++

	// add ZID Header
	zid := &ZIDHeader{ZLen: flooder.zlen, PacketType: LSRFloodPacket, SrcZID: flooder.zoneID}
	msg = append(zid.MarshalBinary(), msg...)

	flooder.macConn.Write(flooder.msec.Encrypt(msg), BroadcastMACAddr)
}

func (flooder *ZoneFlooder) ReceiveFloodedMsg(msgZidHeader *ZIDHeader, msg []byte, payloadProcessor func(net.IP, []byte)) {
	myZone := &Zone{ID: flooder.zoneID, Len: flooder.zlen}
	srcZone := &Zone{ID: msgZidHeader.SrcZID, Len: msgZidHeader.ZLen}

	if !myZone.InZone(srcZone) {
		return
	}
	
	hdr, payload, ok := UnmarshalFloodedPacket(msg)
	if !ok {
		return
	}

	if net.IP.Equal(hdr.SrcIP, flooder.ip) {
		return
	}

	tableSeq, exist := flooder.fTable.Get(hdr.SrcIP)

	if exist && hdr.SeqNum <= tableSeq {
		return
	}

	flooder.fTable.Set(hdr.SrcIP, hdr.SeqNum)

	// Call the payload processor in a separate goroutine to avoid delays during flooding
	go payloadProcessor(hdr.SrcIP, payload)

	log.Println(hdr) // TODO: remove

	// add ZID Header
	zid := &ZIDHeader{ZLen: flooder.zlen, PacketType: LSRFloodPacket, SrcZID: flooder.zoneID}
	msg = append(zid.MarshalBinary(), msg...)

	// reflood the msg
	flooder.macConn.Write(flooder.msec.Encrypt(msg), BroadcastMACAddr)
}

func (f *ZoneFlooder) OnZoneIDChanged(z ZoneID) {
	f.zoneID = z
}
