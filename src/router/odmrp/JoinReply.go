package odmrp

import (
	"fmt"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/msec"
)

type JoinReply struct {
	SeqNo    uint64
	GrpIP    net.IP
	SrcIPs   []net.IP
	NextHops []net.HardwareAddr
}

func (jr *JoinReply) MarshalBinary() []byte {
	extraBytes := 5
	seqNoSize := 8

	totalSize := seqNoSize + net.IPv4len + net.IPv4len*len(jr.SrcIPs) + hwAddrLen*len(jr.NextHops)
	payload := make([]byte, totalSize+extraBytes)
	// 0:2 => number of NextHops
	payload[0] = byte(uint8(len(jr.SrcIPs)))
	payload[1] = byte(uint16(len(jr.NextHops)) >> bitsInByte)
	payload[2] = byte(uint16(len(jr.NextHops)))

	start := extraBytes
	for shift := seqNoSize*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
		payload[start] = byte(jr.SeqNo >> shift)
		start++
	}

	for shift := net.IPv4len*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
		payload[start] = byte(IPv4ToUInt32(jr.GrpIP) >> shift)
		start++
	}

	for i := 0; i < len(jr.SrcIPs); i++ {
		for shift := net.IPv4len*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
			payload[start] = byte(IPv4ToUInt32(jr.SrcIPs[i]) >> shift)
			start++
		}
	}

	for i := 0; i < len(jr.NextHops); i++ {
		for shift := hwAddrLen*bitsInByte - bitsInByte; shift >= 0; shift -= bitsInByte {
			payload[start] = byte(hwAddrToUInt64(jr.NextHops[i]) >> shift)
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

	nextHopsLength := uint16(b[1])<<bitsInByte | uint16(b[2])
	jr.NextHops = make([]net.HardwareAddr, nextHopsLength)

	totalSize := seqNoSize + net.IPv4len + net.IPv4len*len(jr.SrcIPs) + hwAddrLen*len(jr.NextHops)

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

	jr.GrpIP = net.IP(b[start : start+net.IPv4len])
	start += net.IPv4len

	for i := 0; i < len(jr.SrcIPs); i++ {
		jr.SrcIPs[i] = net.IP(b[start : start+net.IPv4len])
		start += net.IPv4len
	}

	for i := 0; i < len(jr.NextHops); i++ {
		jr.NextHops[i] = net.HardwareAddr(b[start : start+hwAddrLen])
		start += hwAddrLen
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

func prettyMacIPs(ips []net.HardwareAddr) string {
	s := "[]net.HardwareAddr{"
	for i := 0; i < len(ips); i++ {
		s += fmt.Sprintf("%v", ips[i].String())
		if i != len(ips)-1 {
			s += ", "
		}
	}
	return s + "}"
}

func (j *JoinReply) String() string {
	return fmt.Sprintf("JoinReply { SeqNo: %d, SrcIPs: %v, GrpIP: %v, NextHops: %v }", j.SeqNo, prettyIPs(j.SrcIPs), j.GrpIP.String(), prettyMacIPs(j.NextHops))
}
