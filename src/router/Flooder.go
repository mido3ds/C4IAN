package main

import (
	"fmt"
	"log"
	"net"

	"github.com/mdlayher/ethernet"
	. "github.com/mido3ds/C4IAN/src/router/msec"
)

type FloodHeader struct {
	// [0:2] checksum here
	SrcIP  net.IP // [2:6]
	SeqNum uint32 // [6:10]
}

const FloodHeaderLen = 2 + 2*4

func UnmarshalFloodedPacket(b []byte) (*FloodHeader, []byte, bool) {
	if len(b) < FloodHeaderLen {
		return nil, nil, false
	}

	// extract checksum
	csum := uint16(b[0])<<8 | uint16(b[1])
	if csum != BasicChecksum(b[2:FloodHeaderLen]) {
		return nil, nil, false
	}

	return &FloodHeader{
		SrcIP:  b[2:6],
		SeqNum: uint32(b[6])<<24 | uint32(b[7])<<16 | uint32(b[8])<<8 | uint32(b[9]),
	}, b[FloodHeaderLen:], true
}

func (f *FloodHeader) MarshalBinary() []byte {
	var header [FloodHeaderLen]byte

	// ip
	copy(header[2:6], f.SrcIP[:])

	// seqnum
	header[6] = byte(f.SeqNum >> 24)
	header[7] = byte(f.SeqNum >> 16)
	header[8] = byte(f.SeqNum >> 8)
	header[9] = byte(f.SeqNum)

	// add checksum
	csum := BasicChecksum(header[2:FloodHeaderLen])
	header[0] = byte(csum >> 8)
	header[1] = byte(csum)

	return header[:]
}

func (f *FloodHeader) String() string {
	return fmt.Sprintf("received a msg flooded by:%#v, with seq=%#v", f.SrcIP.String(), f.SeqNum)
}

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
	// TODO: what should be the destZID?
	zid := &ZIDHeader{ZLen: flooder.zlen, PacketType: LSRFloodPacket, SrcZID: flooder.zoneID, DstZID: flooder.zoneID}
	msg = append(zid.MarshalBinary(), msg...)

	flooder.macConn.Write(flooder.msec.Encrypt(msg), ethernet.Broadcast)
}

func (flooder *ZoneFlooder) ReceiveFloodedMsg(msg []byte, payloadProcessor func(net.IP, []byte)) {
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
	// TODO: what should be the destZID?
	zid := &ZIDHeader{ZLen: flooder.zlen, PacketType: LSRFloodPacket, SrcZID: flooder.zoneID, DstZID: flooder.zoneID}
	msg = append(zid.MarshalBinary(), msg...)

	// reflood the msg
	flooder.macConn.Write(flooder.msec.Encrypt(msg), ethernet.Broadcast)
}

func (f *ZoneFlooder) OnZoneIDChanged(z ZoneID) {
	f.zoneID = z
}

type GlobalFlooder struct {
	seqNumber uint32
	fTable    *FloodingTable
	macConn   *MACLayerConn
	ip        net.IP
	msec      *MSecLayer
}

func NewGlobalFlooder(ip net.IP, iface *net.Interface, etherType ethernet.EtherType, msec *MSecLayer) (*GlobalFlooder, error) {
	// connect to mac layer
	macConn, err := NewMACLayerConn(iface, etherType)
	if err != nil {
		return nil, err
	}

	fTable := NewFloodingTable()

	log.Println("initalized global flooder")

	return &GlobalFlooder{
		seqNumber: 0,
		fTable:    fTable,
		macConn:   macConn,
		ip:        ip,
		msec:      msec,
	}, nil
}

func (f *GlobalFlooder) Flood(msg []byte) {
	hdr := FloodHeader{SrcIP: f.ip, SeqNum: f.seqNumber}
	msg = append(hdr.MarshalBinary(), msg...)

	f.seqNumber++

	f.macConn.Write(f.msec.Encrypt(msg), ethernet.Broadcast)
}

// ReceiveFloodedMsgs inf loop that receives any flooded msgs
// calls `payloadProcessor` when it receives the message, it gives it the header and the payload
// and returns whether to continue flooding or not
func (f *GlobalFlooder) ReceiveFloodedMsgs(payloadProcessor func(*FloodHeader, []byte) bool) {
	for {
		msg := f.macConn.Read()

		hdr, payload, ok := UnmarshalFloodedPacket(msg)
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

		go func() {
			if !payloadProcessor(hdr, payload) {
				return
			}

			f.macConn.Write(f.msec.Encrypt(msg), ethernet.Broadcast)
		}()
	}
}
