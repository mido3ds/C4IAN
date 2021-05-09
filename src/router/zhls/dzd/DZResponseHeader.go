package dzd

import (
	"bytes"
	"encoding/binary"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

const (
	DZResponseChecksumLen = 2
	DZResponseDstIPLen    = 4
	DZResponseDstZoneLen  = 4
	DZResponseTotalLen    = DZResponseChecksumLen + 2*DZResponseDstIPLen + DZResponseDstZoneLen
)

type DZResponseHeader struct {
	dstIP           net.IP
	requiredDstIP   net.IP
	requiredDstZone ZoneID
}

func UnmarshalDZResponseHeader(packet []byte) (*DZResponseHeader, bool) {
	// extract checksum
	csum := binary.BigEndian.Uint16(packet[DZResponseTotalLen-DZRequestChecksumLen:])
	if csum != BasicChecksum(packet[:DZResponseTotalLen-DZRequestChecksumLen]) {
		return nil, false
	}

	dstZone := ZoneID(binary.BigEndian.Uint32(packet[2*DZResponseDstIPLen : 2*DZResponseDstIPLen+DZResponseDstZoneLen]))

	return &DZResponseHeader{
		dstIP:           net.IP(packet[0:DZResponseDstIPLen]),
		requiredDstIP:   net.IP(packet[DZResponseDstIPLen : 2*DZResponseDstIPLen]),
		requiredDstZone: dstZone,
	}, true
}

func (d *DZResponseHeader) MarshalBinary() []byte {
	buffer := bytes.NewBuffer(make([]byte, 0, DZResponseTotalLen))

	buffer.Write(d.dstIP.To4())

	buffer.Write(d.requiredDstIP.To4())
	for i := 24; i >= 0; i -= 8 {
		buffer.WriteByte(byte(d.requiredDstZone >> i))
	}

	// add checksum
	csum := BasicChecksum(buffer.Bytes()[:DZResponseTotalLen-DZResponseChecksumLen])
	for i := 8; i >= 0; i -= 8 {
		buffer.WriteByte(byte(csum >> i))
	}

	return buffer.Bytes()
}

func (d DZResponseHeader) String() string {
	s := "DZResponseMsg: "
	s += "dstIP=" + d.dstIP.String()
	s += ", requiredDstIP=" + d.requiredDstIP.String()
	s += ", requiredDstZone=" + d.requiredDstZone.String() + "\n"
	return s
}
