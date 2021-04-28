package odmrp

import (
	"fmt"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/msec"
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

	totalSize := seqNoSize + net.IPv4len*len(j.SrcIPs) + net.IPv4len*len(j.SrcIPs) + net.IPv4len*len(j.SrcIPs)
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
		for shift := net.IPv4len*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
			payload[start] = byte(IPv4ToUInt32(j.SrcIPs[i]) >> shift)
			start++
		}
	}
	for i := 0; i < len(j.GrpIPs); i++ {
		for shift := net.IPv4len*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
			payload[start] = byte(IPv4ToUInt32(j.GrpIPs[i]) >> shift)
			start++
		}
	}
	for i := 0; i < len(j.Forwarders); i++ {
		for shift := net.IPv4len*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
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

	totalSize := seqNoSize + net.IPv4len*len(jr.SrcIPs) + net.IPv4len*len(jr.GrpIPs) + net.IPv4len*len(jr.Forwarders)

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
		jr.SrcIPs[i] = net.IP(b[start : start+net.IPv4len])
		start += net.IPv4len
	}

	for i := 0; i < len(jr.GrpIPs); i++ {
		jr.GrpIPs[i] = net.IP(b[start : start+net.IPv4len])
		start += net.IPv4len
	}

	for i := 0; i < len(jr.Forwarders); i++ {
		jr.Forwarders[i] = net.IP(b[start : start+net.IPv4len])
		start += net.IPv4len
	}

	return &jr, true
}

func prettyIPs(ips []net.IP) string {
	s := "[]net.IP{"
	for i := 0; i < len(ips); i++ {
		s += fmt.Sprintf("%v", ips[i].String())
		if i != len(ips)-1 {
			s += ", "
		}
	}
	return s + "}"
}

func (j *JoinReply) String() string {
	return fmt.Sprintf("JoinReply { SeqNo: %d, SrcIPs: %v, GrpIPs: %v, Forwarders: %v }", j.SeqNo, prettyIPs(j.SrcIPs), prettyIPs(j.GrpIPs), prettyIPs(j.Forwarders))
}
