package main

import (
	"bytes"
	"encoding/binary"
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

	// read the rest
	header := &ZIDHeader{}
	err := binary.Read(bytes.NewBuffer(packet[2:]), binary.BigEndian, header)

	// remove random salt from zlen
	header.ZLen &= 0b11111

	return header, csum == basicChecksum(packet[2:ZIDHeaderLen]), err
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
	err := binary.Write(m.nonCSUMBuffer, binary.BigEndian, header)
	if err != nil {
		return nil, err
	}

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
