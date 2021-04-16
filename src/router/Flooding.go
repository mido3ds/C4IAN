
package main

import (
	"log"
	"unsafe"
	"encoding/binary"
	"github.com/mdlayher/ethernet"
)

type Flooder struct {
	seqNumber uint32
	fTable *FloodingTable
	router  *Router
	macConn  *MACLayerConn
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
		fTable: fTable,
		router:  router,
		macConn: macConn,
	}, nil
}

// seq (4 Bytes) + srcIP (4 Bytes) + packet
func (flooder *Flooder) Flood(msg []byte) {
	// append srcIP
	srcIP := flooder.router.ip
	msg = append(srcIP, msg...)

	// append sequence number
	seqBytes := (*[4]byte)(unsafe.Pointer(&flooder.seqNumber,))[:]
	msg = append(seqBytes, msg...)
	flooder.seqNumber++
    
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

func (flooder *Flooder) receiveFlood(packet []byte) {
	seq := binary.LittleEndian.Uint32(packet[:4])
	srcIP := packet[4:8]
	tableSeq, exist := flooder.fTable.Get(srcIP)
	if !exist || seq > tableSeq {
		flooder.fTable.Set(srcIP, seq)
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
}


