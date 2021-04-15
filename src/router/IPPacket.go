package main

import (
	"net"
)

type IPHeader struct {
	Version byte
	DestIP  net.IP
}

const IPv4HeaderLen = 20

func UnpackIPHeader(buffer []byte) (*IPHeader, bool) {
	var ip net.IP
	version := byte(buffer[0]) >> 4

	valid := false
	if version == 4 {
		if len(buffer) < IPv4HeaderLen {
			return nil, false
		}
		ip = net.IPv4(buffer[16], buffer[17], buffer[18], buffer[19])
		valid = ipv4Checksum(buffer) == 0
	} else if version == 6 {
		if len(buffer) < 40 {
			return nil, false
		}
		ip = buffer[24:40]
		valid = true
	} else {
		return nil, false
	}

	return &IPHeader{
		Version: version,
		DestIP:  ip,
	}, valid
}

func ipv4Checksum(b []byte) uint16 {
	var sum uint32 = 0
	for i := 0; i < IPv4HeaderLen; i += 2 {
		sum += uint32(b[i])<<8 | uint32(b[i+1])
	}
	return ^uint16((sum >> 16) + sum)
}

func IPToUInt32(ip []byte) uint32 {
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}
