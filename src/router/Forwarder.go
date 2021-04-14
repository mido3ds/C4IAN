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
	table   *ForwardTable
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

	table := NewForwardTable()

	log.Println("initalized forwarder")

	return &Forwarder{
		router:  router,
		macConn: macConn,
		ipConn:  ipConn,
		nfq:     nfq,
		table:   table,
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

		// decrypt and verify
		pd, err := f.router.msec.NewPacketDecrypter(packet)
		if err != nil {
			log.Fatal("failed to build packet decrypter, err: ", err)
		}
		zid, verified := pd.DecryptAndVerifyZID()
		if !verified {
			continue
		} // TODO: check if this is the destZoneID
		ip, verified := pd.DecryptAndVerifyIP()
		if !verified {
			continue
		}

		if imDestination(f.router.ip, ip.DestIP, zid.DestZID) { // i'm destination,
			packet, err := pd.DecryptAll()
			if err != nil {
				continue
			}

			IPHeader := packet[ZIDHeaderLen:]

			// receive message by injecting it in loopback
			err = f.ipConn.Write(IPHeader)
			if err != nil {
				log.Fatal("failed to write to lo interface: ", err)
			}
		} else { // i'm a forwarder
			// determine next hop
			e, ok := f.table.Get(ip.DestIP)
			if !ok {
				// TODO: call controller

				// TODO: for now its all broadcast
				e = &ForwardingEntry{NextHopMAC: ethernet.Broadcast}
				f.table.Set(ip.DestIP, e)
			}

			// hand it directly to the interface
			err = f.macConn.Write(packet, e.NextHopMAC)
			if err != nil {
				log.Fatal("failed to write to the interface: ", err)
			}
		}
	}
}

// ForwardFromIPLayer periodically forwards packets from IP to MAC
// after encrypting them and determining their destination
func (f *Forwarder) ForwardFromIPLayer() {
	zid, err := NewZIDPacketMarshaler(f.router.iface.MTU)
	if err != nil {
		log.Fatal(err)
	}

	packets := f.nfq.GetPackets()

	log.Println("started receiving from IP layer")

	for {
		select {
		case p := <-packets:
			packet := p.Packet.Data()

			ip, err := ParseIPHeader(packet)
			if err != nil {
				log.Println("[error] failed to parse dest ip, drop it, err: ", err)
			} else if IsRawPacket(packet) || imDestination(f.router.ip, ip.DestIP, 0) {
				p.SetVerdict(netfilter.NF_ACCEPT)
			} else { // to out
				// steal packet
				p.SetVerdict(netfilter.NF_DROP)

				// determine next hop
				e, ok := f.table.Get(ip.DestIP)
				if !ok {
					// TODO: call controller

					// TODO: for now its all broadcast
					e = &ForwardingEntry{NextHopMAC: ethernet.Broadcast}
					f.table.Set(ip.DestIP, e)
				}

				// TODO: put this zone id, and zlen
				// wrapp with ZID header
				zidPacket, err := zid.MarshalBinary(&ZIDHeader{ZLen: 1, SrcZID: 2, DestZID: int32(e.DestZoneID)}, packet)
				if err != nil {
					log.Fatal(err)
				}

				// encrypt
				encryptedPacket, err := f.router.msec.Encrypt(zidPacket)
				if err != nil {
					log.Fatal("failed to encrypt packet, err: ", err)
				}

				// hand it directly to the interface
				err = f.macConn.Write(encryptedPacket, e.NextHopMAC)
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

func imDestination(ip, destIP net.IP, destZoneID int32) bool {
	// TODO: use destZID with the ip
	return destIP.Equal(ip) || destIP.IsLoopback()
}

func getNextHopHWAddr(destIP *net.IP) (net.HardwareAddr, error) {
	// TODO: lookup forwarding table for given ipaddr
	return ethernet.Broadcast, nil
}
