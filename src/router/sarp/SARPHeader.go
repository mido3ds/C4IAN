package sarp

import (
	"bytes"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/msec"
)

type SARPHeader struct {
	IP       net.IP
	MAC      net.HardwareAddr
	sendTime int64
}

func UnmarshalSARPHeader(packet []byte) (*SARPHeader, bool) {
	ok := verifySARPHeader(packet)
	if !ok {
		return nil, false
	}
	// sendTime -> packet[10:18]
	sendTime := int64(packet[10])<<56 | int64(packet[11])<<48 | int64(packet[12])<<40 | int64(packet[13])<<32 |
		int64(packet[14])<<24 | int64(packet[15])<<16 | int64(packet[16])<<8 | int64(packet[17])
	return &SARPHeader{
		IP:       net.IP(packet[:4]),
		MAC:      net.HardwareAddr(packet[4:10]),
		sendTime: sendTime,
	}, true
}

func (s *SARPHeader) MarshalBinary() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, sARPTotalLen))

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
