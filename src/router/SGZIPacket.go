package main

import (
	"fmt"
	"math/rand"
)

// "Secure Global Zone Identification Protocol (SG-ZIP)" headers and functions

const SGZIPHeaderLen = 12

var (
	errZeroZlen            = fmt.Errorf("zone len must not be 0")
	errTooSmallSGZIPHeader = fmt.Errorf("SGZIP header is too small")
	errNegativeMTU         = fmt.Errorf("MTU can't be negative")
)

type SGZIPHeader struct {
	Checksum        int16
	RandomSalt      uint16
	ZLen            byte
	DestZID, SrcZID int32
}

func UnpackSGZIPHeader(packet []byte) (*SGZIPHeader, bool, error) {
	if len(packet) < SGZIPHeaderLen {
		return nil, false, errTooSmallSGZIPHeader
	}

	// sgzipChecksum
	var csum int16 = int16(packet[0])<<8 | int16(packet[1])

	header := &SGZIPHeader{
		Checksum:   csum,
		RandomSalt: uint16(packet[2])<<8 | uint16(packet[3]&0b11100000),
		ZLen:       packet[3] & 0b11111,
		DestZID:    int32(packet[4])<<24 | int32(packet[5])<<16 | int32(packet[6])<<8 | int32(packet[7]),
		SrcZID:     int32(packet[8])<<24 | int32(packet[9])<<16 | int32(packet[10])<<8 | int32(packet[11]),
	}

	return header, csum == sgzipChecksum(packet[2:SGZIPHeaderLen]), nil
}

type SGZIPacketMarshaler struct {
	buffer []byte
}

func NewSGZIPacketMarshaler(mtu int) (*SGZIPacketMarshaler, error) {
	if mtu <= 0 {
		return nil, errNegativeMTU
	}

	return &SGZIPacketMarshaler{make([]byte, mtu-SGZIPHeaderLen)}, nil
}

func (m *SGZIPacketMarshaler) MarshalBinary(zlen byte, destZID, srcZID int32, payload []byte) ([]byte, error) {
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

	// sgzipChecksum
	csum := sgzipChecksum(m.buffer[2:SGZIPHeaderLen])
	m.buffer[0] = byte(csum >> 8)
	m.buffer[1] = byte(csum)

	// copy payload
	copy(m.buffer[SGZIPHeaderLen:SGZIPHeaderLen+len(payload)], payload)

	return m.buffer, nil
}

func sgzipChecksum(buf []byte) int16 {
	var sum int16 = 0
	for i := 0; i < len(buf); i++ {
		sum += int16(buf[i])
	}
	return sum
}
