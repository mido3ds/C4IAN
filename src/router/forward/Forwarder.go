package forward

import (
	"log"
	"net"

	"github.com/AkihiroSuda/go-netfilter-queue"
	. "github.com/mido3ds/C4IAN/src/router/flood"
	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/tables"
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

	// multicast controller callback
	mcGetMissingEntries func(grpIP net.IP) bool
	isInMCastGroup      func(grpIP net.IP) bool

	// Unicast controller callbacks
	updateUnicastForwardingTable func(ft *UniForwardTable)

	// braodcast
	bcFlooder *GlobalFlooder
}

func NewForwarder(iface *net.Interface, ip net.IP, msec *MSecLayer,
	neighborsTable *NeighborsTable,
	mcGetMissingEntries func(grpIP net.IP) bool,
	isInMCastGroup func(grpIP net.IP) bool,
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
		MultiForwTable:               MultiForwTable,
		mcGetMissingEntries:          mcGetMissingEntries,
		updateUnicastForwardingTable: updateUnicastForwardingTable,
		isInMCastGroup:               isInMCastGroup,
		bcFlooder:                    NewGlobalFlooder(ip, iface, ZIDBroadcastEtherType, msec),
	}, nil
}

func (f *Forwarder) Start() {
	go f.forwardFromIPLayer()
	go f.forwardZIDFromMACLayer()
	go f.forwardIPFromMACLayer()
	go f.forwardBroadcastMessages()
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

		// TODO: check if this is the destZoneID before decrypting IP header
		ipHdr := f.msec.Decrypt(packet[ZIDHeaderLen : ZIDHeaderLen+IPv4HeaderLen])
		ip, valid := UnmarshalIPHeader(ipHdr)
		if !valid {
			log.Println("Received a packet with an invalid IP header")
			continue
		}

		if imDestination(f.ip, ip.DestIP, zid.DstZID) {
			ipPayload := f.msec.Decrypt(packet[ZIDHeaderLen+IPv4HeaderLen:])
			ipPacket := append(ipHdr, ipPayload...)

			// receive message by injecting it in loopback
			err := f.ipConn.Write(ipPacket)
			if err != nil {
				log.Panic("failed to write to lo interface: ", err)
			}
		} else { // i'm a forwarder
			if valid := IPv4DecrementTTL(ipHdr); !valid {
				log.Println("ttl <= 0, drop packet")
				continue
			}
			IPv4UpdateChecksum(ipHdr)

			// re-encrypt ip hdr
			copy(packet[ZIDHeaderLen:ZIDHeaderLen+IPv4HeaderLen], f.msec.Encrypt(ipHdr))

			e, reachable := getUnicastNextHop(ip.DestIP, f)
			if !reachable {
				// TODO: Should we do anything else here?
				log.Println("Destination unreachable:", ip.DestIP)
				continue
			}
			// hand it directly to the interface
			f.zidMacConn.Write(packet, e.NextHopMAC)
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
		log.Printf("Node IP:%#v, fwd table: %#v\n", f.ip.String(), f.MultiForwTable.String())
		// TODO: speed up by goroutine workers

		// decrypt and verify
		ipHdr := f.msec.Decrypt(packet[:IPv4HeaderLen])
		ip, valid := UnmarshalIPHeader(ipHdr)
		if !valid {
			log.Printf("NOT valid header") // TODO delete
			continue
		}
		log.Printf("valid header") // TODO delete

		if f.isInMCastGroup(ip.DestIP) { // i'm destination,
			ipPayload := f.msec.Decrypt(packet[IPv4HeaderLen:])
			ipPacket := append(ipHdr, ipPayload...)

			// receive message by injecting it in loopback
			err := f.ipConn.Write(ipPacket)
			log.Println("recieved and destination") // TODO remove this
			if err != nil {
				log.Panic("failed to write to lo interface: ", err)
			}
		} else { // TODO remove
			log.Println("recieved but not destination")
		}

		if valid := IPv4DecrementTTL(ipHdr); !valid {
			log.Println("ttl <= 0, drop packet")
			continue
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
			log.Println("ip6 is not supported, drop packet")
			p.SetVerdict(netfilter.NF_DROP)
		} else if imDestination(f.ip, ip.DestIP, 0) || f.isInMCastGroup(ip.DestIP) {
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
				go f.sendUnicast(packet, ip.DestIP)
			case MulticastIPAddr:
				go f.sendMulticast(packet, ip.DestIP)
			case BroadcastIPAddr:
				go f.sendBroadcast(packet)
			default:
				log.Panic("got invalid ip address from ip layer")
			}
		}
	}
}

func (f *Forwarder) forwardBroadcastMessages() {
	f.bcFlooder.ListenForFloodedMsgs(f.onReceiveBroadcastMessage)
}

func (f *Forwarder) onReceiveBroadcastMessage(hdr *FloodHeader, encryptedPacket []byte) []byte {
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

	// inject it into my ip layer
	payload := f.msec.Decrypt(encryptedPacket[ZIDHeaderLen+IPv4HeaderLen:])
	IPv4SetDest(iphdr, f.ip)
	f.ipConn.Write(append(iphdr, payload...))

	// continue flooding
	return encryptedPacket
}

func (f *Forwarder) Close() {
	f.zidMacConn.Close()
	f.ipConn.Close()
}
