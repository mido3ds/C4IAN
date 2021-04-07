package main

import (
	"log"
	"net"

	"github.com/AkihiroSuda/go-netfilter-queue"
	"github.com/mdlayher/ethernet"
)

type Forwarder struct {
	router  *Router
	macConn *MACLayerConn
	ipConn  *IPLayerConn
	nfq     *netfilter.NFQueue
}

func NewForwarder(router *Router) (*Forwarder, error) {
	// connect to mac layer
	macConn, err := NewMACLayerConn(router.iface)
	if err != nil {
		return nil, err
	}

	// connect to ip layer
	ipConn, err := NewIPLayerConn()
	if err != nil {
		return nil, err
	}

	// get packets from netfilter queue
	nfq, err := netfilter.NewNFQueue(0, 200, netfilter.NF_DEFAULT_PACKET_SIZE)
	if err != nil {
		return nil, err
	}

	log.Println("initalized forwarder")

	return &Forwarder{
		router:  router,
		macConn: macConn,
		ipConn:  ipConn,
		nfq:     nfq,
	}, nil
}

// ForwardFromMACLayer continuously receives messages from the interface,
// then either repeats it over loopback (if this is destination), or forwards it for another node.
// The messages may be up to the interface's MTU in size.
func (f *Forwarder) ForwardFromMACLayer() {
	log.Println("started receiving from MAC layer")

	for {
		packet, err := f.macConn.Read()
		if err != nil {
			log.Fatal("failed to read from interface, err: ", err)
		}

		pd, err := f.router.msec.NewPacketDecrypter(packet)
		if err != nil {
			log.Fatal("failed build packet decrypter, err: ", err)
		}

		if imDestination(f.router.ip4, f.router.ip6, pd.DestIP) { // i'm destination,
			packet, err := pd.DecryptAll()
			if err != nil {
				log.Fatal("failed to decrypt rest of the packet")
			}

			// receive message by injecting it in loopback
			err = f.ipConn.Write(packet, pd.Version)
			if err != nil {
				log.Fatal("failed to write to lo interface: ", err)
			}
		} else { // i'm a forwarder
			// determine next hop
			nextHopHWAddr, err := getNextHopHWAddr(&pd.DestIP)
			if err != nil {
				log.Fatal("failed to determine packets destination: ", err)
			}

			// hand it directly to the interface
			err = f.macConn.Write(packet, nextHopHWAddr)
			if err != nil {
				log.Fatal("failed to write to the interface: ", err)
			}
		}
	}
}

// ForwardFromIPLayer periodically forwards packets from IP to MAC
// after encrypting them and determining their destination
func (f *Forwarder) ForwardFromIPLayer() {
	packets := f.nfq.GetPackets()

	log.Println("started receiving from IP layer")

	for {
		select {
		case p := <-packets:
			packet := p.Packet.Data()

			ipPacket, err := ParseIPPacket(packet)
			if err != nil {
				log.Println("[error] failed to parse dest ip, drop it, err: ", err)
				continue
			}

			// TODO: should use imDestination(ip4, ip6, destIP), but it slows down alot
			// without it you can't send to yourself from non-loopback
			// (e.g `ping 10.0.0.1` when your ip is 10.0.0.1`)

			if IsRawPacket(packet) || ipPacket.destIP.IsLoopback() { // to me
				p.SetVerdict(netfilter.NF_ACCEPT)
				continue
			} else { // to out
				// steal packet
				p.SetVerdict(netfilter.NF_DROP)

				// determine next hop
				nextHopHWAddr, err := getNextHopHWAddr(&ipPacket.destIP)
				if err != nil {
					log.Fatal("failed to determine packets destination: ", err)
				}

				if ipPacket.version == 6 {
					// fix stupid issue with IPv6 headers srcIP = ::1
					// which result in the response not returning back to the original sender
					// TODO: find better solution
					copy(packet[8:24], f.router.ip6)
				}

				// encrypt
				encryptedPacket, err := f.router.msec.Encrypt(packet)
				if err != nil {
					log.Fatal("failed to encrypt packet, err: ", err)
				}

				// hand it directly to the interface
				err = f.macConn.Write(encryptedPacket, nextHopHWAddr)
				if err != nil {
					log.Fatal("failed to write to the interface: ", err)
				}
			}
		}
	}
}

func (f *Forwarder) Close() {
	f.macConn.Close()
	f.ipConn.Close()
	f.nfq.Close()
}

func imDestination(ip4, ip6, destIP net.IP) bool {
	return destIP.Equal(ip4) || destIP.Equal(ip6) || destIP.IsLoopback()
}

func getNextHopHWAddr(destIP *net.IP) (net.HardwareAddr, error) {
	// TODO: lookup forwarding table for given ipaddr
	return ethernet.Broadcast, nil
}
