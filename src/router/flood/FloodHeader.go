package flood

import (
	"fmt"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/msec"
)

type FloodHeader struct {
	// [0:2] checksum here
	SrcIP  net.IP // [2:6]
	SeqNum uint32 // [6:10]
}

const floodHeaderLen = 2 + 2*4

func UnmarshalFloodedHeader(b []byte) (*FloodHeader, bool) {
	if len(b) < floodHeaderLen {
		return nil, false
	}

	// extract checksum
	csum := uint16(b[0])<<8 | uint16(b[1])
	if csum != BasicChecksum(b[2:floodHeaderLen]) {
		return nil, false
	}

	return &FloodHeader{
		SrcIP:  b[2:6],
		SeqNum: uint32(b[6])<<24 | uint32(b[7])<<16 | uint32(b[8])<<8 | uint32(b[9]),
	}, true
}

func (f *FloodHeader) MarshalBinary() []byte {
	var header [floodHeaderLen]byte

	// ip
	copy(header[2:6], f.SrcIP[:])

	// seqnum
	header[6] = byte(f.SeqNum >> 24)
	header[7] = byte(f.SeqNum >> 16)
	header[8] = byte(f.SeqNum >> 8)
	header[9] = byte(f.SeqNum)

	// add checksum
	csum := BasicChecksum(header[2:floodHeaderLen])
	header[0] = byte(csum >> 8)
	header[1] = byte(csum)

	return header[:]
}

func (f *FloodHeader) String() string {
	return fmt.Sprintf("received a msg flooded by:%#v, with seq=%#v", f.SrcIP.String(), f.SeqNum)
}
