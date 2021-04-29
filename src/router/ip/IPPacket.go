package ip

import (
	"net"
)

const (
	IPv4HeaderLen = 20
)

type IPHeader struct {
	Version byte
	DestIP  net.IP
}

func UnmarshalIPHeader(buffer []byte) (*IPHeader, bool) {
	if len(buffer) < IPv4HeaderLen {
		return nil, false
	}

	// Only IPv4 is supported
	version := byte(buffer[0]) >> 4
	if version != 4 {
		return nil, false
	}

	// checksum and ttl
	ttl := int8(buffer[8])
	valid := ipv4Checksum(buffer) == 0 && ttl >= 0
	if !valid {
		return nil, false
	}

	return &IPHeader{
		Version: version,
		DestIP:  net.IPv4(buffer[16], buffer[17], buffer[18], buffer[19]).To4(),
	}, true
}

// IPv4DecrementTTL decrements TTL and returns whether
// packet is valid or not
func IPv4DecrementTTL(packet []byte) bool {
	ttl := int8(packet[8])
	if ttl == 0 {
		return false
	}

	packet[8] = byte(ttl) - 1

	// update checksum
	packet[10] = 0 // MSB
	packet[11] = 0 // LSB
	csum := ipv4Checksum(packet)
	packet[10] = byte(csum >> 8) // MSB
	packet[11] = byte(csum)      // LSB

	return true
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
	return net.IPv4(byte(i>>24), byte(i>>16), byte(i>>8), byte(i)).To4()
}
