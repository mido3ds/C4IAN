package main

import (
	"net"
)

type IPHeader struct {
	Version byte
	DestIP  net.IP
	TTL     int8
}

const IPv4HeaderLen = 20

func UnmarshalIPHeader(buffer []byte) (*IPHeader, bool) {
	var ip net.IP
	version := byte(buffer[0]) >> 4
	var ttl int8

	valid := false
	if version == 4 {
		if len(buffer) < IPv4HeaderLen {
			return nil, false
		}

		// DONT call net.IPv4 here, as it puts the 4 bytes
		// at the end of 16 byte array
		// while we expect to have only 4 bytes
		ip = []byte{buffer[16], buffer[17], buffer[18], buffer[19]}

		ttl = int8(buffer[8])
		valid = ipv4Checksum(buffer) == 0 && ttl > 0
	} else if version == 6 {
		if len(buffer) < 40 {
			return nil, false
		}
		ip = buffer[24:40]
		ttl = int8(buffer[7])
		valid = ttl > 0
	} else {
		return nil, false
	}

	// actually it's ttl>0, but there is an edge case
	// when this is the destination and it's 0, it should be processed
	// this is not spec complient to make it easy to write the function
	valid = valid && ttl >= 0

	return &IPHeader{
		Version: version,
		DestIP:  ip,
		TTL:     ttl,
	}, valid
}

func IPv4DecrementTTL(packet []byte) {
	ttl := int8(packet[8])
	packet[8] = byte(ttl) - 1
}

func ipv4Checksum(b []byte) uint16 {
	var sum uint32 = 0
	for i := 0; i < IPv4HeaderLen; i += 2 {
		sum += uint32(b[i])<<8 | uint32(b[i+1])
	}
	return ^uint16((sum >> 16) + sum)
}

func IPv4ToUInt32(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func UInt32ToIPv4(i uint32) net.IP {
	return net.IPv4(byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
}
