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

var (
	errZeroZlen          = fmt.Errorf("zone len must not be 0")
	errTooSmallZIDHeader = fmt.Errorf("ZID header is too small")
	errNegativeMTU       = fmt.Errorf("MTU can't be negative")
)

type ZIDHeader struct {
	ZLen            uint16
	DestZID, SrcZID int32
}

func UnpackZIDHeader(packet []byte) (*ZIDHeader, bool, error) {
	if len(packet) < ZIDHeaderLen {
		return nil, false, errTooSmallZIDHeader
	}

	// extract checksum
	var csum int16 = int16(packet[0])<<8 | int16(packet[1])

	header := &ZIDHeader{
		ZLen:    uint16(packet[3]) & 0b11111,
		DestZID: int32(packet[4])<<24 | int32(packet[5])<<16 | int32(packet[6])<<8 | int32(packet[7]),
		SrcZID:  int32(packet[8])<<24 | int32(packet[9])<<16 | int32(packet[10])<<8 | int32(packet[11]),
	}

	return header, csum == basicChecksum(packet[2:ZIDHeaderLen]), nil
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
	if header.ZLen == 0 {
		return nil, errZeroZlen
	}

	// mix salt and header.ZLen
	header.ZLen |= uint16(rand.Uint32()) << 5

	// write to buffer
	m.buffer[2] = byte(header.ZLen >> 8)
	m.buffer[3] = byte(header.ZLen)

	m.buffer[4] = byte(header.DestZID >> 24)
	m.buffer[5] = byte(header.DestZID >> 16)
	m.buffer[6] = byte(header.DestZID >> 8)
	m.buffer[7] = byte(header.DestZID)

	m.buffer[8] = byte(header.SrcZID >> 24)
	m.buffer[9] = byte(header.SrcZID >> 16)
	m.buffer[10] = byte(header.SrcZID >> 8)
	m.buffer[11] = byte(header.SrcZID)

	// basicChecksum
	csum := basicChecksum(m.buffer[2:ZIDHeaderLen])
	m.buffer[0] = byte(csum >> 8)
	m.buffer[1] = byte(csum)

	// copy payload
	copy(m.buffer[ZIDHeaderLen:ZIDHeaderLen+len(payload)], payload)

	return m.buffer, nil
}

func basicChecksum(buf []byte) int16 {
	var sum int16 = 0
	for i := 0; i < len(buf); i++ {
		sum += int16(buf[i])
	}
	return sum
}
