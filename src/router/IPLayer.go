package main

import (
	"fmt"
	"net"
	"strings"
	"syscall"
)

var (
	loopbackRawAddrIPv4 = syscall.SockaddrInet4{
		Port: 0,
		Addr: [4]byte{127, 0, 0, 1},
	}

	loopbackRawAddrIPv6 = syscall.SockaddrInet6{
		Port: 0,
		Addr: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	}
)

type IPLayerConn struct {
	fd4 int
}

func NewIPLayerConn() (*IPLayerConn, error) {
	fd4, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		return nil, err
	}

	return &IPLayerConn{fd4: fd4}, nil
}

func (c *IPLayerConn) Write(packet []byte) error {
	// mark as raw
	packet[1] |= byte(1)

	return syscall.Sendto(c.fd4, packet, 0, &loopbackRawAddrIPv4)
}

func (c *IPLayerConn) Close() {
	syscall.Close(c.fd4)
}

func IsRawPacket(packet []byte) bool {
	return (packet[1] & byte(1)) == byte(1)
}

func GetMyIPs(iface *net.Interface) (net.IP, net.IP, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, nil, err
	}

	var ip4, ip6 net.IP
	for _, addr := range addrs {
		var ip net.IP

		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if isIPv4(ip.String()) {
			ip4 = ip
		} else if isIPv6(ip.String()) {
			ip6 = ip
		} else {
			return nil, nil, fmt.Errorf("ip is not ip4 or ip6!")
		}
	}

	return ip4, ip6, nil
}

func isIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

func isIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}
