package forward

import (
	"log"
	"net"

	"github.com/AkihiroSuda/go-netfilter-queue"
	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/tables"
	. "github.com/mido3ds/C4IAN/src/router/zhls/dzd"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

type Forwarder struct {
	iface          *net.Interface
	msec           *MSecLayer
	ip             net.IP
	zidMacConn     *MACLayerConn
	ipMacConn      *MACLayerConn
	ipConn         *IPLayerConn
	UniForwTable   *UniForwardTable
	MultiForwTable *MultiForwardTable
	neighborsTable *NeighborsTable
	dzdController  *DZDController

	// multicast controller callback
	mcGetMissingEntries func(grpIP net.IP) (*MultiForwardingEntry, bool)

	// Unicast controller callbacks
	updateUnicastForwardingTable func(ft *UniForwardTable)
}

func NewForwarder(iface *net.Interface, ip net.IP, msec *MSecLayer,
	neighborsTable *NeighborsTable, dzdController *DZDController,
	mcGetMissingEntries func(grpIP net.IP) (*MultiForwardingEntry, bool),
	updateUnicastForwardingTable func(ft *UniForwardTable)) (*Forwarder, error) {
	// connect to mac layer for ZID packets
	zidMacConn, err := NewMACLayerConn(iface, ZIDDataEtherType)
	if err != nil {
		return nil, err
	}

	// connect to mac layer for multicast IP packets
	ipMacConn, err := NewMACLayerConn(iface, IPv4EtherType)
	if err != nil {
		return nil, err
	}

	// connect to ip layer
	ipConn, err := NewIPLayerConn()
	if err != nil {
		return nil, err
	}

	UniForwTable := NewUniForwardTable()
	MultiForwTable := NewMultiForwardTable()
	log.Println("initalized forwarder")

	return &Forwarder{
		iface:                        iface,
		msec:                         msec,
		ip:                           ip,
		zidMacConn:                   zidMacConn,
		ipMacConn:                    ipMacConn,
		ipConn:                       ipConn,
		UniForwTable:                 UniForwTable,
		neighborsTable:               neighborsTable,
		dzdController:                dzdController,
		MultiForwTable:               MultiForwTable,
		mcGetMissingEntries:          mcGetMissingEntries,
		updateUnicastForwardingTable: updateUnicastForwardingTable,
	}, nil
}

func (f *Forwarder) Start() {
	go f.forwardFromIPLayer()
	go f.forwardZIDFromMACLayer()
	go f.forwardIPFromMACLayer()
}

// forwardZIDFromMACLayer continuously receives messages from the interface,
// then either repeats it over loopback (if this is destination), or forwards it for another node.
// The messages may be up to the interface's MTU in size.
func (f *Forwarder) forwardZIDFromMACLayer() {
	log.Println("started receiving from MAC layer")

	for {
		packet := f.zidMacConn.Read()
		// TODO: speed up by goroutine workers

		// decrypt and verify
		zid, valid := UnmarshalZIDHeader(f.msec.Decrypt(packet[:ZIDHeaderLen]))
		if !valid {
			log.Println("Received a packet with an invalid ZID header")
			continue
		}

		ipHdr := f.msec.Decrypt(packet[ZIDHeaderLen : ZIDHeaderLen+IPv4HeaderLen])
		ip, valid := UnmarshalIPHeader(ipHdr)
		if !valid {
			log.Println("Received a packet with an invalid IP header")
			continue
		}

		if imDestination(f.ip, ip.DestIP) {
			ipPayload := f.msec.Decrypt(packet[ZIDHeaderLen+IPv4HeaderLen:])
			ipPacket := append(ipHdr, ipPayload...)

			// receive message by injecting it in loopback
			err := f.ipConn.Write(ipPacket)
			if err != nil {
				log.Panic("failed to write to lo interface: ", err)
			}
		} else { // i'm a forwarder
			if valid := IPv4DecrementTTL(ipHdr); !valid {
				log.Println("ttl < 0, drop packet")
				continue
			}

			// re-encrypt ip hdr
			copy(packet[ZIDHeaderLen:ZIDHeaderLen+IPv4HeaderLen], f.msec.Encrypt(ipHdr))

			var nextHopMAC net.HardwareAddr
			var inMyZone, reachable bool
			if zid.DstZID == MyZone().ID {
				// The destination is in my zone, search in the forwarding table by its ip
				nextHopMAC, inMyZone = f.GetUnicastNextHop(ToNodeID(ip.DestIP))
				if !inMyZone {
					// If the IP is not found in the forwarding table (my zone)
					// although the src claims that it is
					// then the dest may have moved out of this zone
					// and the src have an old value for the dst zone
					// discover its new zone
					dstZoneID, cached := f.dzdController.CachedDstZone(ip.DestIP)
					if cached {
						nextHopMAC, reachable = f.GetUnicastNextHop(ToNodeID(dstZoneID))
						if !reachable {
							// TODO: Should we do anything else here?
							log.Println("Destination zone is unreachable:", zid.DstZID)
							continue
						}
					} else {
						// if dst zone isn't cached, search for it
						// and buffer this msg to be sent when dst zone response arrive
						f.dzdController.FindDstZone(ip.DestIP)
						f.dzdController.BufferPacket(ip.DestIP, packet)
						continue
					}
				}
			} else {
				// The dst is in a different zone,
				// search in the forwarding table by its zone
				nextHopMAC, reachable = f.GetUnicastNextHop(ToNodeID(zid.DstZID))
				if !reachable {
					// TODO: Should we do anything else here?
					log.Println("Destination zone is unreachable:", zid.DstZID)
					continue
				}
			}

			// hand it directly to the interface
			f.zidMacConn.Write(packet, nextHopMAC)
		}
	}
}

// forwardIPFromMACLayer continuously receives messages from the interface,
// then either repeats it over loopback (if this is destination), or forwards it for another node.
// The messages may be up to the interface's MTU in size.
func (f *Forwarder) forwardIPFromMACLayer() {
	log.Println("started receiving from MAC layer")

	for {
		packet := f.ipMacConn.Read()
		// TODO: speed up by goroutine workers

		// decrypt and verify
		ipHdr := f.msec.Decrypt(packet[:IPv4HeaderLen])
		ip, valid := UnmarshalIPHeader(ipHdr)
		if !valid {
			continue
		}

		if imInMulticastGrp(ip.DestIP) { // i'm destination,
			ipPayload := f.msec.Decrypt(packet[IPv4HeaderLen:])
			ipPacket := append(ipHdr, ipPayload...)

			// receive message by injecting it in loopback
			err := f.ipConn.Write(ipPacket)
			if err != nil {
				log.Panic("failed to write to lo interface: ", err)
			}
		}

		if valid := IPv4DecrementTTL(ipHdr); !valid {
			log.Println("ttl < 0, drop packet")
			continue
		}

		// re-encrypt ip hdr
		copy(packet[:IPv4HeaderLen], f.msec.Encrypt(ipHdr))

		// even if im destination, i may forward it
		es, ok := f.MultiForwTable.Get(ip.DestIP)
		if !ok {
			// TODO: call controller
			return
		}

		// write to device driver
		for i := 0; i < len(es.NextHopMACs); i++ {
			f.ipMacConn.Write(packet, es.NextHopMACs[i])
		}
	}
}

// forwardFromIPLayer periodically forwards packets from IP to MAC
// after encrypting them and determining their destination
func (f *Forwarder) forwardFromIPLayer() {
	log.Println("started receiving from IP layer")

	for {
		p := f.ipConn.Read()
		packet := p.Packet.Data()

		// TODO: speed up by goroutine workers
		// TODO: speed up by fanout netfilter feature

		ip, valid := UnmarshalIPHeader(packet)
		if !valid {
			log.Panic("ip header must have been valid from ip layer!")
		}

		if imDestination(f.ip, ip.DestIP) {
			p.SetVerdict(netfilter.NF_ACCEPT)
		} else { // to out
			if ip.DestIP.IsGlobalUnicast() {
				go f.SendUnicast(packet, ip.DestIP)
			} else if ip.DestIP.IsMulticast() {
				go f.sendMulticast(packet, ip.DestIP)
			} else {
				go f.sendBroadcast(packet)
			}

			// sender shall know the papcket is sent
			p.SetVerdict(netfilter.NF_DROP)
		}
	}
}

func (f *Forwarder) Close() {
	f.zidMacConn.Close()
	f.ipConn.Close()
}
