package main

import (
	"fmt"
	"math/rand"
)

// "Zone IDentification (ZID)" protocol structs and functions

const ZIDHeaderLen = 12

var (
	errZeroZlen            = fmt.Errorf("zone len must not be 0")
	errTooSmallSGZIPHeader = fmt.Errorf("ZID header is too small")
	errNegativeMTU         = fmt.Errorf("MTU can't be negative")
)

type ZIDHeader struct {
	Checksum        int16
	RandomSalt      uint16
	ZLen            byte
	DestZID, SrcZID int32
}

func UnpackZIDHeader(packet []byte) (*ZIDHeader, bool, error) {
	if len(packet) < ZIDHeaderLen {
		return nil, false, errTooSmallSGZIPHeader
	}

	// basicChecksum
	var csum int16 = int16(packet[0])<<8 | int16(packet[1])

	header := &ZIDHeader{
		Checksum:   csum,
		RandomSalt: uint16(packet[2])<<8 | uint16(packet[3]&0b11100000),
		ZLen:       packet[3] & 0b11111,
		DestZID:    int32(packet[4])<<24 | int32(packet[5])<<16 | int32(packet[6])<<8 | int32(packet[7]),
		SrcZID:     int32(packet[8])<<24 | int32(packet[9])<<16 | int32(packet[10])<<8 | int32(packet[11]),
	}

	return header, csum == basicChecksum(packet[2:ZIDHeaderLen]), nil
}

type ZIDPacketMarshaler struct {
	buffer []byte
}

func NewZIDPacketMarshaler(mtu int) (*ZIDPacketMarshaler, error) {
	if mtu <= 0 {
		return nil, errNegativeMTU
	}

	return &ZIDPacketMarshaler{make([]byte, mtu-ZIDHeaderLen)}, nil
}

func (m *ZIDPacketMarshaler) MarshalBinary(zlen byte, destZID, srcZID int32, payload []byte) ([]byte, error) {
	if zlen == 0 {
		return nil, errZeroZlen
	}

	// mix salt and zlen
	saltedZlen := uint16(rand.Uint32())
	saltedZlen <<= 5
	saltedZlen |= 0b11111 & uint16(zlen)

	// write to buffer
	m.buffer[2] = byte(saltedZlen >> 8)
	m.buffer[3] = byte(saltedZlen)

	m.buffer[4] = byte(destZID >> 24)
	m.buffer[5] = byte(destZID >> 16)
	m.buffer[6] = byte(destZID >> 8)
	m.buffer[7] = byte(destZID)

	m.buffer[8] = byte(srcZID >> 24)
	m.buffer[9] = byte(srcZID >> 16)
	m.buffer[10] = byte(srcZID >> 8)
	m.buffer[11] = byte(srcZID)

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
