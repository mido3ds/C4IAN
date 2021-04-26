package odmrp

import (
	"fmt"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/msec"
)

const (
	ODMRPDefaultTTL = 100
	ttlSize         = 1
	netIPSize       = 4
	bitsInByte      = 8
)

type JoinQuery struct {
	SeqNo uint64
	TTL   int8
	SrcIP net.IP
	GrpIP net.IP
	Dests []net.IP
}

func (j *JoinQuery) MarshalBinary() []byte {
	extraBytes := 4
	seqNoSize := 8

	destsSize := netIPSize * len(j.Dests)
	totalSize := seqNoSize + ttlSize + netIPSize + netIPSize + destsSize
	payload := make([]byte, totalSize+extraBytes)
	// 0:2 => number of Dests
	payload[0] = byte(uint16(len(j.Dests)) >> bitsInByte)
	payload[1] = byte(uint16(len(j.Dests)))
	start := 4
	for shift := seqNoSize*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
		payload[start] = byte(j.SeqNo >> shift)
		start++
	}
	for shift := ttlSize*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
		payload[start] = byte(j.TTL >> shift)
		start++
	}
	for shift := netIPSize*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
		payload[start] = byte(IPv4ToUInt32(j.SrcIP) >> shift)
		start++
	}
	for shift := netIPSize*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
		payload[start] = byte(IPv4ToUInt32(j.GrpIP) >> shift)
		start++
	}
	for i := 0; i < len(j.Dests); i++ {
		for shift := netIPSize*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
			payload[start] = byte(IPv4ToUInt32(j.Dests[i]) >> shift)
			start++
		}
	}

	// add checksum
	csum := BasicChecksum(payload[extraBytes:])
	payload[2] = byte(csum >> bitsInByte)
	payload[3] = byte(csum)

	return payload[:]
}

// UnmarshalJoinQuery returns false if packet is invalid or TTL < 0
func UnmarshalJoinQuery(b []byte) (*JoinQuery, bool) {
	extraBytes := int64(4)
	seqNoSize := int64(8)
	lenOfDests := uint16(b[0])<<bitsInByte | uint16(b[1])
	destsSize := netIPSize * int64(lenOfDests)
	totalSize := seqNoSize + ttlSize + netIPSize + netIPSize + destsSize

	var jq JoinQuery
	start := extraBytes
	for shift := seqNoSize*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
		jq.SeqNo |= (uint64(b[start]) << shift)
		start++
	}
	jq.TTL = int8(b[start])
	start += ttlSize
	if jq.TTL < 0 {
		return nil, false
	}
	// extract checksum
	csum := uint16(b[2])<<bitsInByte | uint16(b[3])
	// if invalid packet
	if csum != BasicChecksum(b[extraBytes:totalSize+extraBytes]) {
		return nil, false
	}
	jq.SrcIP = net.IP(b[start : start+netIPSize])
	start += netIPSize
	jq.GrpIP = net.IP(b[start : start+netIPSize])
	start += netIPSize

	jq.Dests = make([]net.IP, lenOfDests)
	for i := 0; i < len(jq.Dests); i++ {
		jq.Dests[i] = net.IP(b[start : start+netIPSize])
		start += netIPSize
	}

	return &jq, true
}

func (j *JoinQuery) String() string {
	return fmt.Sprintf("JoinQuery:seq=%#v, TTL=%#v, SrcIP=%#v, GrpIP=%#v", j.SeqNo, j.TTL, j.SrcIP.String(), j.GrpIP.String())
}
