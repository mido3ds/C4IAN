package main

import (
	"bytes"
	"log"
	"net"
	"fmt"
	"github.com/mdlayher/ethernet"
)


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

func (f *FloodHeader) String() string {
	return fmt.Sprintf("received a msg flooded by:%#v, with seq=%#v", f.SrcIP.String(), f.SeqNum)
}

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

func (flooder *Flooder) Flood(msg []byte) {
	hdr := FloodHeader{SrcIP: flooder.router.ip, SeqNum: flooder.seqNumber}
	msg = append(hdr.MarshalBinary(), msg...)

	flooder.seqNumber++

	// add ZID Header
	zid := &ZIDHeader{zLen: 1, packetType: FloodPacket, srcZID: 2, dstZID: 3}
	msg = append(zid.MarshalBinary(), msg...)

	encryptedPacket, err := flooder.router.msec.Encrypt(msg)
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
		log.Println("received my flooded msg")
		return
	}

	tableSeq, exist := flooder.fTable.Get(hdr.SrcIP)

	//log.Println(hdr)
	//log.Println("before: ",  flooder.fTable)

	if exist && hdr.SeqNum <= tableSeq {
		//log.Println("this flooded msg is discarded")
		return
	}

	flooder.fTable.Set(hdr.SrcIP, hdr.SeqNum)

	//log.Println("this flooded msg is accepted")
	//log.Println("after: ",  flooder.fTable)

	log.Println(hdr)

	// add ZID Header
	zid := &ZIDHeader{zLen: 1, packetType: FloodPacket, srcZID: 2, dstZID: 3}
	msg = append(zid.MarshalBinary(), msg...)

	// encrypt the msg
	encryptedPacket, err := flooder.router.msec.Encrypt(msg)
	if err != nil {
		log.Fatal("failed to encrypt packet, err: ", err)
	}

	// reflood the msg
	err = flooder.macConn.Write(encryptedPacket, ethernet.Broadcast)
	if err != nil {
		log.Fatal("failed to write to the device driver: ", err)
	}
}

