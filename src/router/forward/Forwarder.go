package forward

import (
	"log"
	"net"

	"github.com/AkihiroSuda/go-netfilter-queue"
	. "github.com/mido3ds/C4IAN/src/router/database_logger"
	. "github.com/mido3ds/C4IAN/src/router/flood"
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
	ucFlooder      *GlobalFlooder

	// multicast controller callback
	mcGetMissingEntries func(grpIP net.IP) bool
	isInMCastGroup      func(grpIP net.IP) bool

	// Unicast controller callbacks
	updateUnicastForwardingTable func(ft *UniForwardTable)

	// braodcast
	bcFlooder *GlobalFlooder
}

// Create new forwarder
func NewForwarder(iface *net.Interface, ip net.IP, msec *MSecLayer,
	neighborsTable *NeighborsTable, dzdController *DZDController,
	mcGetMissingEntries func(grpIP net.IP) bool,
	isInMCastGroup func(grpIP net.IP) bool,
	updateUnicastForwardingTable func(ft *UniForwardTable),
	timers *TimersQueue) (*Forwarder, error) {
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
	MultiForwTable := NewMultiForwardTable(timers)

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
		isInMCastGroup:               isInMCastGroup,
		bcFlooder:                    NewGlobalFlooder(ip, iface, ZIDBroadcastEtherType, msec),
		ucFlooder:                    NewGlobalFlooder(ip, iface, ZIDFloodEtherType, msec),
	}, nil
}

func (f *Forwarder) Start() {
	go f.forwardFromIPLayer()
	go f.forwardZIDFromMACLayer()
	go f.forwardIPFromMACLayer()
	go f.forwardFloodedBroadcastMessages()
	go f.forwardFloodedUnicastMessages()
}

// forwardZIDFromMACLayer continuously receives messages from the interface,
// then either repeats it over loopback (if this is destination), or forwards it for another node.
// The messages may be up to the interface's MTU in size.
func (f *Forwarder) forwardZIDFromMACLayer() {
	log.Println("started receiving from MAC layer")

	for {
		packet := f.zidMacConn.Read()

		go func() {
			// decrypt and verify
			zid, valid := UnmarshalZIDHeader(f.msec.Decrypt(packet[:ZIDHeaderLen]))
			if !valid {
				log.Println("Received a packet with an invalid ZID header")
				return
			}

			ipHdr := f.msec.Decrypt(packet[ZIDHeaderLen : ZIDHeaderLen+IPv4HeaderLen])
			ip, valid := UnmarshalIPHeader(ipHdr)
			if !valid {
				log.Println("Received a packet with an invalid IP header")
				return
			}

			if imDestination(f.ip, ip.DestIP) {
				ipPayload := f.msec.Decrypt(packet[ZIDHeaderLen+IPv4HeaderLen:])
				ipPacket := append(ipHdr, ipPayload...)

				// receive message by injecting it in loopback
				DatabaseLogger.LogForwarding(packet[ZIDHeaderLen+IPv4HeaderLen:], ip.DestIP)
				log.Println("Forwarding unicast packet to IP layer")
				f.ipConn.Write(ipPacket)
			} else { // i'm a forwarder
				//log.Println("Forward msg from: ", ip.SrcIP, "to: ", ip.DestIP, "ttl: ", ip.TTL)
				if valid := IPv4DecrementTTL(ipHdr); !valid {
					log.Println("ttl <= 0, drop packet")
					return
				}
				IPv4UpdateChecksum(ipHdr)

				// re-encrypt ip hdr
				copy(packet[ZIDHeaderLen:ZIDHeaderLen+IPv4HeaderLen], f.msec.Encrypt(ipHdr))

				myZone := MyZone()
				dstZone := Zone{ID: zid.DstZID, Len: zid.ZLen}

				if dstZone.Len == myZone.Len {
					f.forwardZIDToNextHop(packet, myZone.ID, dstZone.ID, ip.DestIP)
				} else if castedDstZID, intersects := dstZone.Intersection(myZone); !intersects {
					// not same area, but no intersection. safe to ignore difference in zlen
					f.forwardZIDToNextHop(packet, myZone.ID, castedDstZID, ip.DestIP)
				} else {
					// flood in biggest(dstzone, myzone)
					f.ucFlooder.Flood(packet)
				}
			}
		}()
	}
}

func (f *Forwarder) forwardZIDToNextHop(packet []byte, myZID, dstZID ZoneID, dstIP net.IP) {
	var nextHopMAC net.HardwareAddr
	var inMyZone, reachable bool

	if dstZID == myZID {
		// The destination is in my zone, search in the forwarding table by its ip
		nextHopMAC, inMyZone = f.GetUnicastNextHop(ToNodeID(dstIP))
		if !inMyZone {
			// If the IP is not found in the forwarding table (my zone)
			// although the src claims that it is
			// then the dest may have moved out of this zone
			// or the src have an old cached value for the dst zone
			// discover its new zone
			dstZID, cached := f.dzdController.CachedDstZone(dstIP)
			if cached {
				nextHopMAC, reachable = f.GetUnicastNextHop(ToNodeID(dstZID))
				if !reachable {
					// if dst zone is cached, but the cached zone is unreachable
					// try to search one more time
					// and buffer this msg to be sent when dst zone response arrive
					log.Println(dstZID, "is unreachable (Forwarder)")
					f.dzdController.FindDstZone(dstIP)
					f.dzdController.BufferPacket(dstIP, packet, f.forwardBufferedPacketDirectly)
					return
				}
			} else {
				// if dst zone isn't cached, search for it
				// and buffer this msg to be sent when dst zone response arrive
				f.dzdController.FindDstZone(dstIP)
				f.dzdController.BufferPacket(dstIP, packet, f.forwardBufferedPacketDirectly)
				return
			}
		}
	} else {
		// The dst is in a different zone,
		// search in the forwarding table by its zone
		nextHopMAC, reachable = f.GetUnicastNextHop(ToNodeID(dstZID))
		if !reachable {
			// TODO (low priority):
			//		Buffer messages to unreachable zones for a short while and send them
			//	    if the zone becomes reachable
			log.Println(dstZID, "is unreachable (Forwarder)")
			return
		}
	}

	// hand it directly to the interface
	f.zidMacConn.Write(packet, nextHopMAC)
	DatabaseLogger.LogForwarding(packet[ZIDHeaderLen+IPv4HeaderLen:], dstIP)
}

func (f *Forwarder) forwardBufferedPacketDirectly(packet []byte, dstIP net.IP) {
	dstZoneID, cached := f.dzdController.CachedDstZone(dstIP)
	if cached {
		nextHopMAC, reachable := f.GetUnicastNextHop(ToNodeID(dstZoneID))
		if !reachable {
			// TODO (low priority):
			// 		Here I think we have to terminate the search process
			// 		Search for A, A in zone x, x is unreachable(How? I succeded to know that A in x so x must be reachable)
			// 		Unless the nodes moves very quickly, so we can repeat the search process again but how many times?
			log.Println(dstZoneID, "is unreachable (Forwarder)")
			return
		}
		// Update destination zone ID in ZID Header
		zid, ok := UnmarshalZIDHeader(f.msec.Decrypt(packet[:ZIDHeaderLen]))
		if !ok {
			return
		}
		zid.DstZID = dstZoneID
		copy(packet[:ZIDHeaderLen], f.msec.Encrypt(zid.MarshalBinary()))
		// Send to the next hop for the original destination
		f.zidMacConn.Write(packet, nextHopMAC)
		DatabaseLogger.LogForwarding(packet[ZIDHeaderLen+IPv4HeaderLen:], dstIP)
	} else {
		log.Panicln("Dst Zone must be cached here")
	}
}

// forwardIPFromMACLayer continuously receives messages from the interface,
// then either repeats it over loopback (if this is destination), or forwards it for another node.
// The messages may be up to the interface's MTU in size.
func (f *Forwarder) forwardIPFromMACLayer() {
	log.Println("started receiving from MAC layer")

	for {
		packet := f.ipMacConn.Read()

		go func() {
			log.Printf("Node IP:%#v, fwd table: %#v\n", f.ip.String(), f.MultiForwTable.String())

			// decrypt and verify
			ipHdr := f.msec.Decrypt(packet[:IPv4HeaderLen])
			ip, valid := UnmarshalIPHeader(ipHdr)
			if !valid {
				log.Println("Received a packet with an invalid IP header")
				return
			}

			if f.isInMCastGroup(ip.DestIP) { // i'm destination,
				ipPayload := f.msec.Decrypt(packet[IPv4HeaderLen:])
				ipPacket := append(ipHdr, ipPayload...)

				// receive message by injecting it in loopback
				log.Println("Forwarding unicast packet to IP layer")
				f.ipConn.Write(ipPacket)
			}

			if valid := IPv4DecrementTTL(ipHdr); !valid {
				log.Println("ttl <= 0, drop packet")
				return
			}
			IPv4UpdateChecksum(ipHdr)

			// even if im destination, i may forward it
			es, exist := f.MultiForwTable.Get(ip.DestIP)

			if exist {
				// re-encrypt ip hdr
				copy(packet[:IPv4HeaderLen], f.msec.Encrypt(ipHdr))
				// write to device driver
				for item := range es.Items.Iter() {
					log.Printf("Forward packet to:%#v\n", item.Value.(*NextHopEntry).NextHop.String())
					f.ipMacConn.Write(packet, item.Value.(*NextHopEntry).NextHop)
				}
			}
		}()
	}
}

// forwardFromIPLayer periodically forwards packets from IP to MAC
// after encrypting them and determining their destination
func (f *Forwarder) forwardFromIPLayer() {
	log.Println("started receiving from IP layer")

	for {
		p := f.ipConn.Read()

		go func() {
			packet := p.Packet.Data()
			log.Println("Received IP packet")

			// TODO (low priority): speed up by fanout netfilter feature

			ip, valid := UnmarshalIPHeader(packet)

			if !valid {
				log.Println("ip6 is not supported, drop packet")
				p.SetVerdict(netfilter.NF_DROP)
			} else if imDestination(f.ip, ip.DestIP) || f.isInMCastGroup(ip.DestIP) {
				p.SetVerdict(netfilter.NF_ACCEPT)
			} else { // to out
				// sender shall know the papcket is sent
				p.SetVerdict(netfilter.NF_DROP)

				// reset ttl if ip layer, weirdly, gave low ttl
				// doesn't work for traceroute
				IPv4ResetTTL(packet)
				IPv4UpdateChecksum(packet)

				switch iptype := GetIPAddrType(ip.DestIP); iptype {
				case UnicastIPAddr:
					f.sendUnicast(packet, ip.DestIP)
				case MulticastIPAddr:
					log.Println("Sending multicast packet")
					f.sendMulticast(packet, ip.DestIP)
				case BroadcastIPAddr:
					log.Println("Sending broadcast packet")
					f.sendBroadcast(packet)
				default:
					log.Panic("got invalid ip address from ip layer")
				}
			}
		}()
	}
}

func (f *Forwarder) forwardFloodedBroadcastMessages() {
	f.bcFlooder.ListenForFloodedMsgs(func(encryptedPacket []byte) []byte {
		zidhdr := f.msec.Decrypt(encryptedPacket[:ZIDHeaderLen])
		zid, ok := UnmarshalZIDHeader(zidhdr)
		if !ok {
			// invalid zid header, stop here
			return nil
		}

		iphdr := f.msec.Decrypt(encryptedPacket[ZIDHeaderLen : ZIDHeaderLen+IPv4HeaderLen])
		ip, ok := UnmarshalIPHeader(iphdr)
		if !ok {
			// invalid ip header, stop here
			return nil
		}

		r1 := BroadcastRadius(ip.DestIP)
		r2 := MyZone().ID.DistTo(zid.SrcZID)
		if r2 > r1 {
			// out of zone broadcast, stop here
			return nil
		}

		go func() {
			// inject it into my ip layer
			payload := f.msec.Decrypt(encryptedPacket[ZIDHeaderLen+IPv4HeaderLen:])
			IPv4SetDest(iphdr, f.ip)
			f.ipConn.Write(append(iphdr, payload...))
		}()

		// continue flooding
		return encryptedPacket
	})
}

func (f *Forwarder) forwardFloodedUnicastMessages() {
	f.ucFlooder.ListenForFloodedMsgs(func(encryptedPacket []byte) []byte {
		zidhdr := f.msec.Decrypt(encryptedPacket[:ZIDHeaderLen])
		zid, ok := UnmarshalZIDHeader(zidhdr)
		if !ok {
			// invalid zid header, stop here
			return nil
		}

		myzone := MyZone()
		dstzone := Zone{ID: zid.DstZID, Len: zid.ZLen}
		if _, intersects := myzone.Intersection(dstzone); !intersects {
			// out of zone, stop here
			return nil
		}

		iphdr := f.msec.Decrypt(encryptedPacket[ZIDHeaderLen : ZIDHeaderLen+IPv4HeaderLen])
		ip, ok := UnmarshalIPHeader(iphdr)
		if !ok {
			// invalid ip header, stop here
			return nil
		}

		go func() {
			if ip.DestIP.Equal(f.ip) {
				// inject it into my ip layer
				payload := f.msec.Decrypt(encryptedPacket[ZIDHeaderLen+IPv4HeaderLen:])
				IPv4SetDest(iphdr, f.ip)
				f.ipConn.Write(append(iphdr, payload...))
			}
		}()

		// continue flooding
		return encryptedPacket
	})
}

func (f *Forwarder) Close() {
	f.zidMacConn.Close()
	f.ipMacConn.Close()
	f.ipConn.Close()

	f.ucFlooder.Close()
	f.bcFlooder.Close()
}
