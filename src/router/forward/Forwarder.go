package forward

import (
	"log"
	"net"

	"github.com/AkihiroSuda/go-netfilter-queue"
	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/tables"
	. "github.com/mido3ds/C4IAN/src/router/zhls"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

type Forwarder struct {
	iface          *net.Interface
	msec           *MSecLayer
	ip             net.IP
	zlen           byte
	zoneID         ZoneID
	zidMacConn     *MACLayerConn
	ipMacConn      *MACLayerConn
	ipConn         *IPLayerConn
	UniForwTable   *UniForwardTable
	MultiForwTable *MultiForwardTable
	neighborsTable *NeighborsTable

	// multicast controller callback
	mcGetMissingEntries func(grpIP net.IP) (*MultiForwardingEntry, bool)

	// Unicast controller callbacks
	updateUnicastForwardingTable func(ft *UniForwardTable)
}

func NewForwarder(iface *net.Interface, ip net.IP, msec *MSecLayer, zlen byte,
	neighborsTable *NeighborsTable,
	mcGetMissingEntries func(grpIP net.IP) (*MultiForwardingEntry, bool),
	updateUnicastForwardingTable func(ft *UniForwardTable)) (*Forwarder, error) {
	// connect to mac layer for ZID packets
	zidMacConn, err := NewMACLayerConn(iface, ZIDEtherType)
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
		zlen:                         zlen,
		zidMacConn:                   zidMacConn,
		ipMacConn:                    ipMacConn,
		ipConn:                       ipConn,
		UniForwTable:                 UniForwTable,
		neighborsTable:               neighborsTable,
		MultiForwTable:               MultiForwTable,
		mcGetMissingEntries:          mcGetMissingEntries,
		updateUnicastForwardingTable: updateUnicastForwardingTable,
	}, nil
}

func (f *Forwarder) Start(controllerChannel chan *UnicastControlPacket) {
	go f.forwardFromIPLayer()
	go f.forwardZIDFromMACLayer(controllerChannel)
	go f.forwardIPFromMACLayer()
}

// forwardZIDFromMACLayer continuously receives messages from the interface,
// then either repeats it over loopback (if this is destination), or forwards it for another node.
// The messages may be up to the interface's MTU in size.
func (f *Forwarder) forwardZIDFromMACLayer(controllerChannel chan *UnicastControlPacket) {
	log.Println("started receiving from MAC layer")

	for {
		packet := f.zidMacConn.Read()
		// TODO: speed up by goroutine workers

		// decrypt and verify
		pd := f.msec.NewPacketDecrypter(packet)
		zid, valid := UnmarshalZIDHeader(pd.DecryptN(ZIDHeaderLen))
		if !valid {
			log.Println("Received a packet with an invalid ZID header")
			continue
		}

		if zid.IsControlPacket() {
			decryptedPacket := pd.DecryptAll()
			controllerChannel <- &UnicastControlPacket{ZIDHeader: zid, Payload: decryptedPacket[ZIDHeaderLen:]}
			continue
		}

		// TODO: check if this is the destZoneID before decrypting IP header
		ip, valid := UnmarshalIPHeader(pd.DecryptN(IPv4HeaderLen))
		if !valid {
			log.Println("Received a packet with an invalid IP header")
			continue
		}

		if imDestination(f.ip, ip.DestIP, zid.DstZID) { // i'm destination,
			ippacket := pd.DecryptAll()[ZIDHeaderLen:]

			// receive message by injecting it in loopback
			err := f.ipConn.Write(ippacket)
			if err != nil {
				log.Panic("failed to write to lo interface: ", err)
			}
		} else { // i'm a forwarder
			IPv4DecrementTTL(packet[ZIDHeaderLen:])

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
		// TODO: speed up by goroutine workers

		// decrypt and verify
		pd := f.msec.NewPacketDecrypter(packet)
		ip, valid := UnmarshalIPHeader(pd.DecryptN(IPv4HeaderLen))
		if !valid {
			continue
		}

		if imInMulticastGrp(ip.DestIP) { // i'm destination,
			// receive message by injecting it in loopback
			err := f.ipConn.Write(pd.DecryptAll())
			if err != nil {
				log.Panic("failed to write to lo interface: ", err)
			}
		}

		// even if im destination, i may forward it
		IPv4DecrementTTL(packet)

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

		if IsInjectedPacket(packet) || imDestination(f.ip, ip.DestIP, 0) {
			p.SetVerdict(netfilter.NF_ACCEPT)
		} else { // to out
			if ip.DestIP.IsGlobalUnicast() {
				go f.sendUnicast(packet, ip.DestIP)
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

func (f *Forwarder) OnZoneIDChanged(z ZoneID) {
	f.zoneID = z
}

func (f *Forwarder) Close() {
	f.zidMacConn.Close()
	f.ipConn.Close()
}
