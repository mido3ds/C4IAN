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

func (f *ZoneFlooder) Flood(msg []byte) {
	f.seqNumber++

	zidHdr := ZIDHeader{ZLen: f.zlen, PacketType: LSRFloodPacket, SrcZID: f.zoneID}
	encZidHdr := f.msec.Encrypt(zidHdr.MarshalBinary())

	fldHdr := FloodHeader{SrcIP: f.ip, SeqNum: f.seqNumber}
	encFloodedMsg := f.msec.Encrypt(append(fldHdr.MarshalBinary(), msg...))

	f.macConn.Write(append(encZidHdr, encFloodedMsg...), BroadcastMACAddr)
}

func (f *ZoneFlooder) ReceiveFloodedMsg(msg []byte, payloadProcessor func(net.IP, []byte)) {
	hdr, payload, ok := UnmarshalFloodedHeader(msg)
	if !ok {
		return
	}

	if net.IP.Equal(hdr.SrcIP, f.ip) {
		return
	}

	tableSeq, exist := f.fTable.Get(hdr.SrcIP)

	if exist && hdr.SeqNum <= tableSeq {
		return
	}

	f.fTable.Set(hdr.SrcIP, hdr.SeqNum)

	// Call the payload processor in a separate goroutine to avoid delays during flooding
	go payloadProcessor(hdr.SrcIP, payload)

	// re-wrap with zid header
	// TODO: should i use the ZID header that came with it?
	zidHdr := ZIDHeader{ZLen: f.zlen, PacketType: LSRFloodPacket, SrcZID: f.zoneID}
	encZidHdr := f.msec.Encrypt(zidHdr.MarshalBinary())

	encFloodedMsg := f.msec.Encrypt(msg)

	// reflood the msg
	f.macConn.Write(append(encZidHdr, encFloodedMsg...), BroadcastMACAddr)
}

func (f *ZoneFlooder) OnZoneIDChanged(z ZoneID) {
	f.zoneID = z
}

func (f *ZoneFlooder) Close() {
	f.macConn.Close()
}
