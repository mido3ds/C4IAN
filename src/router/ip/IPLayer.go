package ip

import (
	"fmt"
	"net"
	"strings"
	"syscall"

	"github.com/AkihiroSuda/go-netfilter-queue"
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
	fd4     int
	nfq     *netfilter.NFQueue
	packets <-chan netfilter.NFPacket
}

func NewIPLayerConn() (*IPLayerConn, error) {
	fd4, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		return nil, err
	}

	// get packets from netfilter queue
	nfq, err := netfilter.NewNFQueue(0, 200, netfilter.NF_DEFAULT_PACKET_SIZE)
	if err != nil {
		return nil, err
	}
	packets := nfq.GetPackets()

	return &IPLayerConn{fd4: fd4,
		nfq:     nfq,
		packets: packets,
	}, nil
}

func (c *IPLayerConn) Write(packet []byte) error {
	markPacketAsInjected(packet)
	return syscall.Sendto(c.fd4, packet, 0, &loopbackRawAddrIPv4)
}

func (c *IPLayerConn) Read() netfilter.NFPacket {
	return <-c.packets
}

func (c *IPLayerConn) Close() {
	syscall.Close(c.fd4)
	c.nfq.Close()
}

// markPacketAsInjected puts a mark in a reserved bit in IPv4 header
// so we can notice it when it comes back in netfilter-queue
func markPacketAsInjected(b []byte) {
	b[1] |= 1
}

func IsInjectedPacket(packet []byte) bool {
	return packet[1]&1 == 1
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
			return nil, nil, fmt.Errorf("ip is not ip4 or ip6")
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
