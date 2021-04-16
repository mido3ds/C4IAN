package main

import (
	"bytes"
	"log"
	"net"

	"github.com/mdlayher/ethernet"
)

type Flooder struct {
	seqNumber uint32
	fTable    *FloodingTable
	router    *Router
	macConn   *MACLayerConn
}

func NewFlooder(router *Router) (*Flooder, error) {
	// connect to mac layer
	macConn, err := NewMACLayerConn(router.iface)
	if err != nil {
		return nil, err
	}

	fTable := NewFloodingTable()

	log.Println("initalized flooder")

	return &Flooder{
		seqNumber: 0,
		fTable:    fTable,
		router:    router,
		macConn:   macConn,
	}, nil
}

// seq (4 Bytes) + srcIP (4 Bytes) + packet
func (flooder *Flooder) Flood(msg []byte) {
	hdr := FloodHeader{SrcIP: flooder.router.ip, SeqNum: flooder.seqNumber}
	msg = append(hdr.MarshalBinary(), msg...)

	flooder.seqNumber++

	// ADD ZID Header
	zid, err := NewZIDPacketMarshaler(flooder.router.iface.MTU)
	if err != nil {
		log.Fatal(err)
	}

	packet, err := zid.MarshalBinary(&ZIDHeader{zLen: 1, packetType: FloodPacket, srcZID: 2, dstZID: 3}, msg)
	if err != nil {
		log.Fatal(err)
	}

	encryptedPacket, err := flooder.router.msec.Encrypt(packet)
	if err != nil {
		log.Fatal("failed to encrypt packet, err: ", err)
	}

	err = flooder.macConn.Write(encryptedPacket, ethernet.Broadcast)
	if err != nil {
		log.Fatal("failed to write to the device driver: ", err)
	}
}

func (flooder *Flooder) ReceiveFloodedMsg(msg []byte) {
	hdr, ok := UnpackFloodHeader(msg)
	if !ok {
		return
	}

	if bytes.Equal(hdr.SrcIP, flooder.router.ip) {
		log.Println("My flooded msg returned to me")
		return
	}

	log.Println("I received a msg from ", net.IP(hdr.SrcIP))

	tableSeq, exist := flooder.fTable.Get(hdr.SrcIP)

	log.Println("Seq Number: ", hdr.SeqNum)
	log.Println("Exist: ", exist)
	log.Println("Table seq: ", tableSeq)

	if exist && tableSeq <= hdr.SeqNum {
		return
	}

	flooder.fTable.Set(hdr.SrcIP, hdr.SeqNum)

	// ADD ZID Header
	zid, err := NewZIDPacketMarshaler(flooder.router.iface.MTU)
	if err != nil {
		log.Fatal(err)
	}

	packet, err := zid.MarshalBinary(&ZIDHeader{zLen: 1, packetType: FloodPacket, srcZID: 2, dstZID: 3}, msg)
	if err != nil {
		log.Fatal(err)
	}

	// encrypt the msg
	encryptedPacket, err := flooder.router.msec.Encrypt(packet)
	if err != nil {
		log.Fatal("failed to encrypt packet, err: ", err)
	}

	// reflood the msg
	err = flooder.macConn.Write(encryptedPacket, ethernet.Broadcast)
	if err != nil {
		log.Fatal("failed to write to the device driver: ", err)
	}
}

type FloodHeader struct {
	// [0:2] checksum here
	SrcIP  net.IP // [2:6]
	SeqNum uint32 // [6:10]
}

const FloodHeaderLen = 2 + 2*4

func UnpackFloodHeader(b []byte) (*FloodHeader, bool) {
	if len(b) < FloodHeaderLen {
		return nil, false
	}

	// extract checksum
	csum := uint16(b[0])<<8 | uint16(b[1])
	if csum != BasicChecksum(b[2:FloodHeaderLen]) {
		return nil, false
	}

	return &FloodHeader{
		SrcIP:  b[2:6],
		SeqNum: uint32(b[6])<<24 | uint32(b[7])<<16 | uint32(b[8])<<8 | uint32(b[9]),
	}, true
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
