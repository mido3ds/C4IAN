package main

import (
	"net"
)

type JoinReply struct {
	SeqNo      uint64
	SrcIPs     []net.IP
	GrpIPs     []net.IP
	Forwarders []net.IP
}

func (j *JoinReply) MarshalBinary() []byte {
	extraBytes := 5
	seqNoSize := 8

	totalSize := seqNoSize + netIPSize*len(j.SrcIPs) + netIPSize*len(j.SrcIPs) + netIPSize*len(j.SrcIPs)
	payload := make([]byte, totalSize+extraBytes)
	// 0:2 => number of Forwarders
	payload[0] = byte(uint8(len(j.SrcIPs)))
	payload[1] = byte(uint8(len(j.GrpIPs)))
	payload[2] = byte(uint8(len(j.Forwarders)))

	start := extraBytes
	for shift := seqNoSize*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
		payload[start] = byte(j.SeqNo >> shift)
		start++
	}
	for i := 0; i < len(j.SrcIPs); i++ {
		for shift := netIPSize*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
			payload[start] = byte(IPv4ToUInt32(j.SrcIPs[i]) >> shift)
			start++
		}
	}
	for i := 0; i < len(j.GrpIPs); i++ {
		for shift := netIPSize*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
			payload[start] = byte(IPv4ToUInt32(j.GrpIPs[i]) >> shift)
			start++
		}
	}
	for i := 0; i < len(j.Forwarders); i++ {
		for shift := netIPSize*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
			payload[start] = byte(IPv4ToUInt32(j.Forwarders[i]) >> shift)
			start++
		}
	}

	// add checksum
	csum := BasicChecksum(payload[extraBytes:])
	payload[3] = byte(csum >> bitsInByte)
	payload[4] = byte(csum)

	return payload[:]
}

// UnmarshalJoinReply returns false if packet is invalid
func UnmarshalJoinReply(b []byte) (*JoinReply, bool) {
	extraBytes := 5
	seqNoSize := 8

	var jr JoinReply
	jr.SrcIPs = make([]net.IP, uint8(b[0]))
	jr.GrpIPs = make([]net.IP, uint8(b[1]))
	jr.Forwarders = make([]net.IP, uint8(b[2]))

	totalSize := seqNoSize + netIPSize*len(jr.SrcIPs) + netIPSize*len(jr.GrpIPs) + netIPSize*len(jr.Forwarders)

	// extract checksum
	csum := uint16(b[3])<<bitsInByte | uint16(b[4])
	// if invalid packet
	if csum != BasicChecksum(b[extraBytes:totalSize+extraBytes]) {
		return nil, false
	}

	start := extraBytes
	for shift := seqNoSize*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
		jr.SeqNo |= (uint64(b[start]) << shift)
		start++
	}

	for i := 0; i < len(jr.SrcIPs); i++ {
		jr.SrcIPs[i] = net.IP(b[start : start+netIPSize])
		start += netIPSize
	}

	for i := 0; i < len(jr.GrpIPs); i++ {
		jr.GrpIPs[i] = net.IP(b[start : start+netIPSize])
		start += netIPSize
	}

	for i := 0; i < len(jr.Forwarders); i++ {
		jr.Forwarders[i] = net.IP(b[start : start+netIPSize])
		start += netIPSize
	}

	return &jr, true
}
