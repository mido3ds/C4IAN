package odmrp

import (
	"fmt"
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/msec"
)

const CostSize = 16

type joinReply struct {
	seqNum   uint64
	destIP   net.IP
	grpIP    net.IP
	prevHop  net.HardwareAddr
	srcIPs   []net.IP
	nextHops []net.HardwareAddr
	cost     uint16
}

func (jr *joinReply) marshalBinary() []byte {
	extraBytes := 4
	seqNoSize := 8

	if len(jr.srcIPs) != len(jr.nextHops) {
		log.Panic("JoinReply SrcIPs should have the same length as NextHops, " + jr.String()) // TODO remove
	}

	totalSize := seqNoSize + net.IPv4len + net.IPv4len + MACAddrLen + net.IPv4len*len(jr.srcIPs) + MACAddrLen*len(jr.nextHops) + CostSize
	payload := make([]byte, totalSize+extraBytes)
	// 0:1 => number of NextHops & SrcIPs
	payload[0] = byte(uint16(len(jr.nextHops)) >> BitsInByte)
	payload[1] = byte(uint16(len(jr.nextHops)))

	start := extraBytes
	for shift := seqNoSize*BitsInByte - BitsInByte; shift >= 0; shift -= BitsInByte {
		payload[start] = byte(jr.seqNum >> shift)
		start++
	}

	for shift := net.IPv4len*BitsInByte - BitsInByte; shift >= 0; shift -= BitsInByte {
		payload[start] = byte(IPv4ToUInt32(jr.destIP) >> shift)
		start++
	}

	for shift := net.IPv4len*BitsInByte - BitsInByte; shift >= 0; shift -= BitsInByte {
		payload[start] = byte(IPv4ToUInt32(jr.grpIP) >> shift)
		start++
	}

	for shift := MACAddrLen*BitsInByte - BitsInByte; shift >= 0; shift -= BitsInByte {
		payload[start] = byte(HwAddrToUInt64(jr.prevHop) >> shift)
		start++
	}

	for i := 0; i < len(jr.srcIPs); i++ {
		for shift := net.IPv4len*BitsInByte - BitsInByte; shift >= 0; shift -= BitsInByte {
			payload[start] = byte(IPv4ToUInt32(jr.srcIPs[i]) >> shift)
			start++
		}
	}

	for i := 0; i < len(jr.nextHops); i++ {
		for shift := MACAddrLen*BitsInByte - BitsInByte; shift >= 0; shift -= BitsInByte {
			payload[start] = byte(HwAddrToUInt64(jr.nextHops[i]) >> shift)
			start++
		}
	}

	for shift := CostSize*BitsInByte - BitsInByte; shift >= 0; shift -= BitsInByte {
		payload[start] = byte(jr.cost >> shift)
		start++
	}

	// add checksum
	csum := BasicChecksum(payload[extraBytes:])
	payload[2] = byte(csum >> BitsInByte)
	payload[3] = byte(csum)

	return payload[:]
}

// unmarshalJoinReply returns false if packet is invalid
func unmarshalJoinReply(b []byte) (*joinReply, bool) {
	extraBytes := 4
	seqNoSize := 8

	var jr joinReply

	count := uint16(b[0])<<BitsInByte | uint16(b[1])
	jr.srcIPs = make([]net.IP, count)
	jr.nextHops = make([]net.HardwareAddr, count)

	totalSize := seqNoSize + net.IPv4len + net.IPv4len + MACAddrLen + net.IPv4len*len(jr.srcIPs) + MACAddrLen*len(jr.nextHops) + CostSize

	// extract checksum
	csum := uint16(b[2])<<BitsInByte | uint16(b[3])
	// if invalid packet
	if csum != BasicChecksum(b[extraBytes:totalSize+extraBytes]) {
		return nil, false
	}

	start := extraBytes
	for shift := seqNoSize*BitsInByte - BitsInByte; shift >= 0; shift -= BitsInByte {
		jr.seqNum |= (uint64(b[start]) << shift)
		start++
	}

	jr.destIP = net.IP(b[start : start+net.IPv4len])
	start += net.IPv4len

	jr.grpIP = net.IP(b[start : start+net.IPv4len])
	start += net.IPv4len

	jr.prevHop = net.HardwareAddr(b[start : start+MACAddrLen])
	start += MACAddrLen

	for i := 0; i < len(jr.srcIPs); i++ {
		jr.srcIPs[i] = net.IP(b[start : start+net.IPv4len])
		start += net.IPv4len
	}

	for i := 0; i < len(jr.nextHops); i++ {
		jr.nextHops[i] = net.HardwareAddr(b[start : start+MACAddrLen])
		start += MACAddrLen
	}

	for shift := CostSize*BitsInByte - BitsInByte; shift >= 0; shift -= BitsInByte {
		jr.cost |= (uint16(b[start]) << shift)
		start++
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

func (j *joinReply) String() string {
	return fmt.Sprintf("JoinReply { SeqNo: %d, DestIP: %#v, GrpIP: %v, SrcIPs: %v, NextHops: %v, Cost:%d }", j.seqNum, j.destIP.String(), j.grpIP.String(), prettyIPs(j.srcIPs), prettyMacIPs(j.nextHops), j.cost)
}
