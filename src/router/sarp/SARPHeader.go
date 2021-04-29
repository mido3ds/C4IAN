package sarp

import (
	"bytes"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/msec"
)

type SARPType uint8

const (
	hashLen       = 64 // bytes at the end
	sARPHeaderLen = 19 // excluding the hash at the end
	sARPTotalLen  = sARPHeaderLen + hashLen
)

const (
	SARPReq SARPType = iota
	SARPRes
)

type SARPHeader struct {
	Type     SARPType
	IP       net.IP
	MAC      net.HardwareAddr
	sendTime int64
}

func UnmarshalSARPHeader(packet []byte) (*SARPHeader, bool) {
	ok := verifySARPHeader(packet)
	if !ok {
		return nil, false
	}
	// sendTime -> packet[11:19]
	sendTime := int64(packet[11])<<56 | int64(packet[12])<<48 | int64(packet[13])<<40 | int64(packet[14])<<32 |
		int64(packet[15])<<24 | int64(packet[16])<<16 | int64(packet[17])<<8 | int64(packet[18])
	return &SARPHeader{
		Type:     SARPType(packet[0]),
		IP:       net.IP(packet[1:5]),
		MAC:      net.HardwareAddr(packet[5:11]),
		sendTime: sendTime,
	}, true
}

func (s *SARPHeader) MarshalBinary() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, sARPTotalLen))

	buf.WriteByte(byte(s.Type))
	buf.Write(s.IP.To4())
	buf.Write(s.MAC)
	for i := 56; i >= 0; i -= 8 {
		buf.WriteByte(byte(s.sendTime >> i))
	}
	buf.Write(HashSHA3(buf.Bytes()[:sARPHeaderLen]))

	return buf.Bytes()
}

func verifySARPHeader(b []byte) bool {
	if len(b) < sARPTotalLen {
		return false
	}

	return VerifySHA3Hash(b[:sARPHeaderLen], b[sARPHeaderLen:sARPTotalLen])
}
