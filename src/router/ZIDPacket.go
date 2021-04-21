package main

import (
	"fmt"
	"math/rand"
)

// "Zone IDentification (ZID)" protocol structs and functions

// TODO: GPS location -> zoneID
// TODO: translate zondIDs from zlen to another

const ZIDHeaderLen = 12

type PacketType uint8

const (
	// TODO: Add actual data/control types
	DataPacket PacketType = iota
	LSRFloodPacket
	DummyControlPacket
	SARPReq
	SARPRes
)

var (
	errZeroZlen    = fmt.Errorf("zone len must not be 0")
	errNegativeMTU = fmt.Errorf("MTU can't be negative")
)

type ZIDHeader struct {
	packetType     PacketType // Most significant 4 bits of the 4th byte
	zLen           uint8      // Least significant 4 bits of the 4th byte
	dstZID, srcZID int32
}

func (z *ZIDHeader) isControlPacket() bool {
	return z.packetType == LSRFloodPacket ||
		z.packetType == SARPReq ||
		z.packetType == SARPRes ||
		z.packetType == DummyControlPacket
}

func UnpackZIDHeader(packet []byte) (*ZIDHeader, bool) {
	if len(packet) < ZIDHeaderLen {
		return nil, false
	}

	// extract checksum
	csum := uint16(packet[0])<<8 | uint16(packet[1])
	if csum != BasicChecksum(packet[2:ZIDHeaderLen]) {
		return nil, false
	}

	return &ZIDHeader{
		packetType: PacketType(packet[3] >> 4),
		zLen:       uint8(packet[3] & 0b1111),
		dstZID:     int32(packet[4])<<24 | int32(packet[5])<<16 | int32(packet[6])<<8 | int32(packet[7]),
		srcZID:     int32(packet[8])<<24 | int32(packet[9])<<16 | int32(packet[10])<<8 | int32(packet[11]),
	}, true
}

func (header *ZIDHeader) MarshalBinary() []byte {
	var buf [ZIDHeaderLen]byte

	// random salt
	buf[2] = byte(rand.Uint32())

	// packet type + zlen
	buf[3] = byte(header.packetType)<<4 | (byte(header.zLen) & 0b1111)

	// destZID
	buf[4] = byte(header.dstZID >> 24)
	buf[5] = byte(header.dstZID >> 16)
	buf[6] = byte(header.dstZID >> 8)
	buf[7] = byte(header.dstZID)

	// srcZID
	buf[8] = byte(header.srcZID >> 24)
	buf[9] = byte(header.srcZID >> 16)
	buf[10] = byte(header.srcZID >> 8)
	buf[11] = byte(header.srcZID)

	// add checksum
	csum := BasicChecksum(buf[2:ZIDHeaderLen])
	buf[0] = byte(csum >> 8)
	buf[1] = byte(csum)

	return buf[:]
}
