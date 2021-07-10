package flood

import (
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

type ZoneFlooder struct {
	seqNumber   uint32
	floodingTbl *floodingTable
	macConn     *MACLayerConn
	ip          net.IP
	msec        *MSecLayer
}

func NewZoneFlooder(iface *net.Interface, ip net.IP, msec *MSecLayer) *ZoneFlooder {
	// connect to mac layer
	macConn, err := NewMACLayerConn(iface, ZoneFloodEtherType)
	if err != nil {
		log.Panic("failed to create device connection, err: ", err)
	}

	log.Println("initalized zone flooder")

	return &ZoneFlooder{
		seqNumber:   0,
		floodingTbl: newFloodingTable(),
		macConn:     macConn,
		ip:          ip,
		msec:        msec,
	}
}

func (f *ZoneFlooder) Flood(msg []byte) {
	f.seqNumber++

	zidHeader := MyZIDHeader(0)
	encryptedZIDHeader := f.msec.Encrypt(zidHeader.MarshalBinary())

	floodHeader := FloodHeader{SrcIP: f.ip, SeqNum: f.seqNumber}
	encryptedFloodHeader := f.msec.Encrypt(floodHeader.MarshalBinary())

	encryptedPayload := f.msec.Encrypt(msg)

	f.macConn.Write(append(encryptedZIDHeader, append(encryptedFloodHeader, encryptedPayload...)...), BroadcastMACAddr)
}

func (f *ZoneFlooder) ListenForFloodedMsgs(payloadProcessor func(net.IP, []byte)) {
	for {
		msg := f.macConn.Read()
		go f.handleFloodedMsg(msg, payloadProcessor)
	}
}

func (f *ZoneFlooder) handleFloodedMsg(msg []byte, payloadProcessor func(net.IP, []byte)) {
	zidHeader, ok := UnmarshalZIDHeader(f.msec.Decrypt(msg[:ZIDHeaderLen]))
	if !ok {
		return
	}

	myZone := MyZone()
	srcZone := &Zone{ID: zidHeader.SrcZID, Len: zidHeader.ZLen}

	if !myZone.Equal(srcZone) {
		return
	}

	floodHeader, ok := UnmarshalFloodedHeader(f.msec.Decrypt(msg[ZIDHeaderLen : ZIDHeaderLen+floodHeaderLen]))
	if !ok {
		return
	}

	if net.IP.Equal(floodHeader.SrcIP, f.ip) {
		return
	}

	if f.floodingTbl.isHighestSeqNum(floodHeader.SrcIP, floodHeader.SeqNum) {
		f.floodingTbl.set(floodHeader.SrcIP, floodHeader.SeqNum)

		// Call the payload processor in a separate goroutine to avoid delays during flooding
		go payloadProcessor(floodHeader.SrcIP, f.msec.Decrypt(msg[ZIDHeaderLen+floodHeaderLen:]))

		// reflood the msg
		f.macConn.Write(msg, BroadcastMACAddr)
	}
}

func (f *ZoneFlooder) Close() {
	f.macConn.Close()
}
