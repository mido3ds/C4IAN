package dzd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

const (
	DZRequestChecksumLen = 2
	DZRequestIPsLen      = 8 // src IP and dst IP
	DzRequestSrcZoneLen  = 4
	numOfVisitedZonesLen = 4
	DZRequestHeaderLen   = DZRequestChecksumLen + DZRequestIPsLen + DzRequestSrcZoneLen + numOfVisitedZonesLen
)

type DZRequestHeader struct {
	srcIP         net.IP
	srcZone       ZoneID
	requiredDstIP net.IP
	visitedZones  []ZoneID
}

func UnmarshalDZRequestHeader(packet []byte) (*DZRequestHeader, bool) {
	numOfVisistedZones := binary.BigEndian.Uint32(packet[:numOfVisitedZonesLen])
	DZDTotalLen := DZRequestHeaderLen + 4*numOfVisistedZones

	// extract checksum
	csum := binary.BigEndian.Uint16(packet[DZDTotalLen-DZRequestChecksumLen:])
	if csum != BasicChecksum(packet[numOfVisitedZonesLen:DZDTotalLen-DZRequestChecksumLen]) {
		return nil, false
	}

	visitedZones := make([]ZoneID, numOfVisistedZones)
	start := DZRequestIPsLen + numOfVisitedZonesLen + DzRequestSrcZoneLen
	for i := uint32(0); i < numOfVisistedZones; i++ {
		visitedZones[i] = ZoneID(binary.BigEndian.Uint32(packet[start : start+4]))
		start += 4
	}

	return &DZRequestHeader{
		srcZone:       ZoneID(binary.BigEndian.Uint32(packet[4:8])),
		srcIP:         net.IP(packet[8:12]),
		requiredDstIP: net.IP(packet[12:16]),
		visitedZones:  visitedZones,
	}, true
}

func (d *DZRequestHeader) MarshalBinary() []byte {
	numOfVisistedZones := uint32(len(d.visitedZones))
	DZRequestTotalLen := DZRequestHeaderLen + 4*numOfVisistedZones

	buffer := bytes.NewBuffer(make([]byte, 0, DZRequestTotalLen))

	for i := 24; i >= 0; i -= 8 {
		buffer.WriteByte(byte(numOfVisistedZones >> i))
	}

	for i := 24; i >= 0; i -= 8 {
		buffer.WriteByte(byte(d.srcZone >> i))
	}

	buffer.Write(d.srcIP.To4())
	buffer.Write(d.requiredDstIP.To4())

	for _, zoneID := range d.visitedZones {
		for i := 24; i >= 0; i -= 8 {
			buffer.WriteByte(byte(zoneID >> i))
		}
	}

	// add checksum
	csum := BasicChecksum(buffer.Bytes()[numOfVisitedZonesLen : DZRequestTotalLen-DZRequestChecksumLen])
	for i := 8; i >= 0; i -= 8 {
		buffer.WriteByte(byte(csum >> i))
	}

	return buffer.Bytes()
}

func (d DZRequestHeader) String() string {
	s := "DZRequestMsg: "
	s += "srcIP=" + d.srcIP.String()
	s += ", srcZone=" + d.srcZone.String()
	s += ", requiredDstIP=" + d.requiredDstIP.String() + "\n"
	s += fmt.Sprint(d.visitedZones) + "\n"
	return s
}
