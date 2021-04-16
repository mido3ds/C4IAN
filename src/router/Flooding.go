package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"unsafe"

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
	// append srcIP
	srcIP := flooder.router.ip
	msg = append(srcIP, msg...)

	// append sequence number
	seqBytes := (*[4]byte)(unsafe.Pointer(&flooder.seqNumber))[:]
	msg = append(seqBytes, msg...)
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

func (flooder *Flooder) receiveFlood(msg []byte) {
	srcIP := msg[4:8]
	myIP := flooder.router.ip

	if bytes.Equal(srcIP, myIP) {
		log.Println("My flooded msg returned to me")
		return
	}

	log.Println("I received a msg from ", net.IP(srcIP))

	seq := binary.LittleEndian.Uint32(msg[:4])
	tableSeq, exist := flooder.fTable.Get(srcIP)

	log.Println("Seq Number: ", seq)
	log.Println("Exist: ", exist)
	log.Println("Table seq: ", tableSeq)

	if exist && tableSeq <= seq {
		return
	}

	flooder.fTable.Set(srcIP, seq)

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
