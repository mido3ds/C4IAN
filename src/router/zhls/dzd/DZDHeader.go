package dzd

import (
	"bytes"
	"encoding/binary"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

const (
	hashLen              = 64 // bytes at the end
	IPsLen               = 8  // src IP and dst IP
	numOfVisitedZonesLen = 4
)

type DZDHeader struct {
	srcIP        net.IP
	dstIP        net.IP
	visitedZones []ZoneID
}

func UnmarshalDZDPHeader(packet []byte) (*DZDHeader, bool) {
	ok := verifyDZDHeader(packet)
	if !ok {
		return nil, false
	}

	numOfVisistedZones := binary.BigEndian.Uint32(packet[IPsLen : IPsLen+numOfVisitedZonesLen])

	visitedZones := make([]ZoneID, numOfVisistedZones)
	start := IPsLen + numOfVisitedZonesLen
	for i := uint32(0); i < numOfVisistedZones; i++ {
		visitedZones[i] = ZoneID(binary.BigEndian.Uint32(packet[start : start+4]))
		start += 4
	}

	return &DZDHeader{
		srcIP:        net.IP(packet[:4]),
		dstIP:        net.IP(packet[4:8]),
		visitedZones: visitedZones,
	}, true
}

func (d *DZDHeader) MarshalBinary() []byte {
	numOfVisistedZones := uint32(len(d.visitedZones))
	DZDHeaderLen := IPsLen + numOfVisitedZonesLen + 4*numOfVisistedZones
	DZDTotalLen := DZDHeaderLen + hashLen

	buffer := bytes.NewBuffer(make([]byte, 0, DZDTotalLen))

	buffer.Write(d.srcIP.To4())
	buffer.Write(d.dstIP.To4())

	for i := 24; i >= 0; i -= 8 {
		buffer.WriteByte(byte(numOfVisistedZones >> i))
	}

	for _, zoneID := range d.visitedZones {
		for i := 24; i >= 0; i -= 8 {
			buffer.WriteByte(byte(zoneID >> i))
		}
	}

	buffer.Write(HashSHA3(buffer.Bytes()[:DZDHeaderLen]))

	return buffer.Bytes()
}

func verifyDZDHeader(dzdPacket []byte) bool {
	numOfVisistedZones := binary.BigEndian.Uint32(dzdPacket[IPsLen : IPsLen+numOfVisitedZonesLen])
	DZDHeaderLen := IPsLen + numOfVisitedZonesLen + 4*numOfVisistedZones
	DZDTotalLen := DZDHeaderLen + hashLen

	if uint32(len(dzdPacket)) < DZDTotalLen {
		return false
	}

	return VerifySHA3Hash(dzdPacket[:DZDHeaderLen], dzdPacket[DZDHeaderLen:DZDTotalLen])
}
