package main

import (
	"fmt"
	"math/rand"
)

// "Zone IDentification (ZID)" protocol structs and functions

const (
	ZIDHeaderLen = 12

	// Make use of an unassigned EtherType to differentiate between ZID traffic and other traffic
	// https://www.iana.org/assignments/ieee-802-numbers/ieee-802-numbers.xhtml
	ZIDEtherType = 0x7031
)

type PacketType uint8

const (
	// TODO: Add actual data/control types
	DataPacket PacketType = iota
	LSRFloodPacket
	DummyControlPacket
)

var (
	errNegativeMTU = fmt.Errorf("MTU can't be negative")
)

type ZIDHeader struct {
	PacketType     PacketType // Most significant 4 bits of the 4th byte
	ZLen           uint8      // Least significant 4 bits of the 4th byte
	DstZID, SrcZID ZoneID
}

func (z *ZIDHeader) isControlPacket() bool {
	return z.PacketType == LSRFloodPacket ||
		z.PacketType == DummyControlPacket
}

func UnmarshalZIDHeader(packet []byte) (*ZIDHeader, bool) {
	if len(packet) < ZIDHeaderLen {
		return nil, false
	}

	// extract checksum
	csum := uint16(packet[0])<<8 | uint16(packet[1])
	if csum != BasicChecksum(packet[2:ZIDHeaderLen]) {
		return nil, false
	}

	return &ZIDHeader{
		PacketType: PacketType(packet[3] >> 4),
		ZLen:       uint8(packet[3] & 0b1111),
		DstZID:     ZoneID(uint32(packet[4])<<24 | uint32(packet[5])<<16 | uint32(packet[6])<<8 | uint32(packet[7])),
		SrcZID:     ZoneID(uint32(packet[8])<<24 | uint32(packet[9])<<16 | uint32(packet[10])<<8 | uint32(packet[11])),
	}, true
}

func (header *ZIDHeader) MarshalBinary() []byte {
	var buf [ZIDHeaderLen]byte

	// random salt
	buf[2] = byte(rand.Uint32())

	// packet type + zlen
	buf[3] = byte(header.PacketType)<<4 | (byte(header.ZLen) & 0b1111)

	// destZID
	buf[4] = byte(header.DstZID >> 24)
	buf[5] = byte(header.DstZID >> 16)
	buf[6] = byte(header.DstZID >> 8)
	buf[7] = byte(header.DstZID)

	// SrcZID
	buf[8] = byte(header.SrcZID >> 24)
	buf[9] = byte(header.SrcZID >> 16)
	buf[10] = byte(header.SrcZID >> 8)
	buf[11] = byte(header.SrcZID)

	// add checksum
	csum := BasicChecksum(buf[2:ZIDHeaderLen])
	buf[0] = byte(csum >> 8)
	buf[1] = byte(csum)

	return buf[:]
}
