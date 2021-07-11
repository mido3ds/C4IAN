package odmrp

import (
	"fmt"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/msec"
)

const (
	odmrpDefaultTTL = 100
	ttlSize         = 1
	bitsInByte      = 8
	hwAddrLen       = 6
)

type joinQuery struct {
	seqNum  uint64
	ttl     int8
	srcIP   net.IP
	prevHop net.HardwareAddr
	grpIP   net.IP
	dests   []net.IP
}

func (j *joinQuery) marshalBinary() []byte {
	extraBytes := 4
	seqNoSize := 8

	destsSize := net.IPv4len * len(j.dests)
	totalSize := seqNoSize + ttlSize + net.IPv4len + hwAddrLen + net.IPv4len + destsSize
	payload := make([]byte, totalSize+extraBytes)
	// 0:2 => number of Dests
	payload[0] = byte(uint16(len(j.dests)) >> bitsInByte)
	payload[1] = byte(uint16(len(j.dests)))
	start := 4
	for shift := seqNoSize*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
		payload[start] = byte(j.seqNum >> shift)
		start++
	}
	for shift := ttlSize*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
		payload[start] = byte(j.ttl >> shift)
		start++
	}
	for shift := net.IPv4len*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
		payload[start] = byte(IPv4ToUInt32(j.srcIP) >> shift)
		start++
	}
	for shift := hwAddrLen*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
		payload[start] = byte(HwAddrToUInt64(j.prevHop) >> shift)
		start++
	}
	for shift := net.IPv4len*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
		payload[start] = byte(IPv4ToUInt32(j.grpIP) >> shift)
		start++
	}
	for i := 0; i < len(j.dests); i++ {
		for shift := net.IPv4len*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
			payload[start] = byte(IPv4ToUInt32(j.dests[i]) >> shift)
			start++
		}
	}

	// add checksum
	csum := BasicChecksum(payload[extraBytes:])
	payload[2] = byte(csum >> bitsInByte)
	payload[3] = byte(csum)

	return payload[:]
}

func unmarshalJoinQuery(b []byte) (*joinQuery, bool) {
	extraBytes := int64(4)
	seqNoSize := int64(8)
	lenOfDests := uint16(b[0])<<bitsInByte | uint16(b[1])
	destsSize := net.IPv4len * int64(lenOfDests)
	totalSize := seqNoSize + ttlSize + net.IPv4len + hwAddrLen + net.IPv4len + destsSize

	var jq joinQuery
	start := extraBytes
	for shift := seqNoSize*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
		jq.seqNum |= (uint64(b[start]) << shift)
		start++
	}
	jq.ttl = int8(b[start])
	start += ttlSize

	// extract checksum
	csum := uint16(b[2])<<bitsInByte | uint16(b[3])
	// if invalid packet
	if csum != BasicChecksum(b[extraBytes:totalSize+extraBytes]) {
		return nil, false
	}
	jq.srcIP = net.IP(b[start : start+net.IPv4len])
	start += net.IPv4len
	jq.prevHop = net.HardwareAddr(b[start : start+hwAddrLen])
	start += hwAddrLen
	jq.grpIP = net.IP(b[start : start+net.IPv4len])
	start += net.IPv4len

	jq.dests = make([]net.IP, lenOfDests)
	for i := 0; i < len(jq.dests); i++ {
		jq.dests[i] = net.IP(b[start : start+net.IPv4len])
		start += net.IPv4len
	}

	return &jq, true
}

func (j *joinQuery) String() string {
	return fmt.Sprintf("JoinQuery { SeqNo: %d, TTL: %#v, SrcIP: %v, PrevHop: %v, GrpIP: %v, Dests: %v }", j.seqNum, j.ttl, j.srcIP.String(), j.prevHop.String(), j.grpIP.String(), j.dests)
}
