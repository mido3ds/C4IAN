package main

import (
	"bytes"
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
	ControlPacket
	FloodPacket
	DummyControlPacket
	SARP
)

var (
	errZeroZlen    = fmt.Errorf("zone len must not be 0")
	errNegativeMTU = fmt.Errorf("MTU can't be negative")
)

type ZIDHeader struct {
	packetType     PacketType // Left 4 bits of the 4th byte
	zLen           uint8      // Right 4 bits of the 4th byte
	dstZID, srcZID int32
}

func (z *ZIDHeader) isControlPacket() bool {
	return z.packetType == ControlPacket || 
		   z.packetType == FloodPacket ||
		   z.packetType == SARP ||
		   z.packetType == DummyControlPacket 
}

func UnpackZIDHeader(packet []byte) (*ZIDHeader, bool) {
	if len(packet) < ZIDHeaderLen {
		return nil, false
	}

	// extract checksum
	csum := uint16(packet[0])<<8 | uint16(packet[1])

	header := &ZIDHeader{
		packetType: PacketType(packet[3] >> 4),
		zLen:       uint8(packet[3] & 0b1111),
		dstZID:     int32(packet[4])<<24 | int32(packet[5])<<16 | int32(packet[6])<<8 | int32(packet[7]),
		srcZID:     int32(packet[8])<<24 | int32(packet[9])<<16 | int32(packet[10])<<8 | int32(packet[11]),
	}

	return header, csum == basicChecksum(packet[2:ZIDHeaderLen])
}

type ZIDPacketMarshaler struct {
	buffer        []byte
	nonCSUMBuffer *bytes.Buffer
}

func NewZIDPacketMarshaler(mtu int) (*ZIDPacketMarshaler, error) {
	if mtu <= 0 {
		return nil, errNegativeMTU
	}

	buffer := make([]byte, mtu-ZIDHeaderLen)
	nonCSUMBuffer := bytes.NewBuffer(buffer[:])

	return &ZIDPacketMarshaler{buffer: buffer, nonCSUMBuffer: nonCSUMBuffer}, nil
}

func (m *ZIDPacketMarshaler) MarshalBinary(header *ZIDHeader, payload []byte) ([]byte, error) {
	if header.zLen == 0 {
		return nil, errZeroZlen
	}

	// write to buffer
	m.buffer[2] = byte(rand.Uint32()) // Random salt
	m.buffer[3] = byte(header.packetType)<<4 | (byte(header.zLen) & 0b1111)

	m.buffer[4] = byte(header.dstZID >> 24)
	m.buffer[5] = byte(header.dstZID >> 16)
	m.buffer[6] = byte(header.dstZID >> 8)
	m.buffer[7] = byte(header.dstZID)

	m.buffer[8] = byte(header.srcZID >> 24)
	m.buffer[9] = byte(header.srcZID >> 16)
	m.buffer[10] = byte(header.srcZID >> 8)
	m.buffer[11] = byte(header.srcZID)

	// basicChecksum
	csum := basicChecksum(m.buffer[2:ZIDHeaderLen])
	m.buffer[0] = byte(csum >> 8)
	m.buffer[1] = byte(csum)

	// copy payload
	copy(m.buffer[ZIDHeaderLen:ZIDHeaderLen+len(payload)], payload)

	return m.buffer, nil
}

func basicChecksum(buf []byte) uint16 {
	var sum uint16 = 0
	for i := 0; i < len(buf); i++ {
		sum += uint16(buf[i])
	}
	return sum
}
