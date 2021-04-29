package zid

import (
	"fmt"
	"math/rand"

	. "github.com/mido3ds/C4IAN/src/router/msec"
)

// "Zone IDentification (ZID)" protocol structs and functions

const (
	ZIDHeaderLen = 12
)

var (
	errNegativeMTU = fmt.Errorf("MTU can't be negative")
)

type ZIDHeader struct {
	ZLen           uint8
	DstZID, SrcZID ZoneID
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
		ZLen:   uint8(packet[3]),
		DstZID: ZoneID(uint32(packet[4])<<24 | uint32(packet[5])<<16 | uint32(packet[6])<<8 | uint32(packet[7])),
		SrcZID: ZoneID(uint32(packet[8])<<24 | uint32(packet[9])<<16 | uint32(packet[10])<<8 | uint32(packet[11])),
	}, true
}

func (header *ZIDHeader) MarshalBinary() []byte {
	var buf [ZIDHeaderLen]byte

	// Random salt
	buf[2] = byte(rand.Uint32())

	// Zlen
	buf[3] = byte(header.ZLen)

	// DstZID
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
