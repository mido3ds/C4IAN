package main

import (
	"bytes"
	"log"
	"net"

	"github.com/AkihiroSuda/go-netfilter-queue"
	"github.com/mdlayher/ethernet"
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
	uniForwTable   *UniForwardTable
	multiForwTable *MultiForwardTable
	neighborsTable *NeighborsTable

	// multicast controller callback
	mcGetMissingEntries func(grpIP net.IP) (*MultiForwardingEntry, bool)
}

func NewForwarder(iface *net.Interface, ip net.IP, msec *MSecLayer, zlen byte,
	neighborsTable *NeighborsTable,
	mcGetMissingEntries func(grpIP net.IP) (*MultiForwardingEntry, bool)) (*Forwarder, error) {
	// connect to mac layer for ZID packets
	zidMacConn, err := NewMACLayerConn(iface, ZIDEtherType)
	if err != nil {
		return nil, err
	}

	// connect to mac layer for multicast IP packets
	ipMacConn, err := NewMACLayerConn(iface, ethernet.EtherTypeIPv4)
	if err != nil {
		return nil, err
	}

	// connect to ip layer
	ipConn, err := NewIPLayerConn()
	if err != nil {
		return nil, err
	}

	uniForwTable := NewUniForwardTable()
	multiForwTable := NewMultiForwardTable()

	log.Println("initalized forwarder")

	return &Forwarder{
		iface:               iface,
		msec:                msec,
		ip:                  ip,
		zlen:                zlen,
		zidMacConn:          zidMacConn,
		ipMacConn:           ipMacConn,
		ipConn:              ipConn,
		uniForwTable:        uniForwTable,
		neighborsTable:      neighborsTable,
		multiForwTable:      multiForwTable,
		mcGetMissingEntries: mcGetMissingEntries,
	}, nil
}

func (f *Forwarder) Start(controllerChannel chan *UnicastControlPacket) {
	go f.forwardFromIPLayer()
	go f.forwardZIDFromMACLayer(controllerChannel)
	go f.forwardIPFromMACLayer()
}

func (f *Forwarder) broadcastDummy() {
	dummy := []byte("Dummy")
	zid := &ZIDHeader{ZLen: f.zlen, PacketType: LSRFloodPacket, SrcZID: f.zoneID, DstZID: f.zoneID}
	packet := append(zid.MarshalBinary(), dummy...)

	f.zidMacConn.Write(f.msec.Encrypt(packet), ethernet.Broadcast)

	log.Println("Broadcasting dummy control packet..")
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
		zid, valid := pd.DecryptAndVerifyZID()
		if !valid {
			log.Println("Received a packet with invalid ZID header")
			continue
		}

		if zid.isControlPacket() {
			packet, err := pd.DecryptAll()
			if err != nil {
				continue
			}

			controllerChannel <- &UnicastControlPacket{zidHeader: zid, payload: packet[ZIDHeaderLen:]}
			continue
		}

		// TODO: check if this is the destZoneID before decrypting IP header
		ip, valid := pd.DecryptAndVerifyIP()
		if !valid {
			continue
		}

		if imDestination(f.ip, ip.DestIP, zid.DstZID) { // i'm destination,
			packet, err := pd.DecryptAll()
			if err != nil {
				continue
			}

			ippacket := packet[ZIDHeaderLen:]

			// receive message by injecting it in loopback
			err = f.ipConn.Write(ippacket)
			if err != nil {
				log.Panic("failed to write to lo interface: ", err)
			}
		} else { // i'm a forwarder
			IPv4DecrementTTL(packet[ZIDHeaderLen:])

			e, ok := getNextHop(ip.DestIP, f.uniForwTable, f.neighborsTable, f.zoneID)
			if !ok {
				// TODO: call controller
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
		ip, valid := pd.DecryptAndVerifyIP()
		if !valid {
			continue
		}

		if imInMulticastGrp(ip.DestIP) { // i'm destination,
			packet, err := pd.DecryptAll()
			if err != nil {
				continue
			}

			// receive message by injecting it in loopback
			err = f.ipConn.Write(packet)
			if err != nil {
				log.Panic("failed to write to lo interface: ", err)
			}
		}

		// even if im destination, i may forward it
		IPv4DecrementTTL(packet)

		es, ok := f.multiForwTable.Get(ip.DestIP)
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

func (f *Forwarder) sendUnicast(packet []byte, destIP net.IP) {
	e, ok := getNextHop(destIP, f.uniForwTable, f.neighborsTable, f.zoneID)
	if !ok {
		// TODO: call controller
		return
	}

	zid := &ZIDHeader{ZLen: f.zlen, PacketType: DataPacket, SrcZID: f.zoneID, DstZID: e.DestZoneID}

	// build packet
	buffer := bytes.NewBuffer(make([]byte, 0, f.iface.MTU))
	buffer.Write(zid.MarshalBinary())
	buffer.Write(packet)

	// encrypt
	encryptedPacket := f.msec.Encrypt(buffer.Bytes())

	// write to device driver
	f.zidMacConn.Write(encryptedPacket, e.NextHopMAC)
}

func (f *Forwarder) sendMulticast(packet []byte, grpIP net.IP) {
	es, ok := f.multiForwTable.Get(grpIP)
	if !ok {
		es, ok = f.mcGetMissingEntries(grpIP)
		if !ok {
			return
		}
	}

	// encrypt
	encryptedPacket := f.msec.Encrypt(packet)

	// write to device driver
	for i := 0; i < len(es.NextHopMACs); i++ {
		f.zidMacConn.Write(encryptedPacket, es.NextHopMACs[i])
	}
}

func (f *Forwarder) sendBroadcast(packet []byte) {
	// write to device driver
	// TODO: for now ethernet broadcast
	f.zidMacConn.Write(f.msec.Encrypt(packet), ethernet.Broadcast)
}

func (f *Forwarder) OnZoneIDChanged(z ZoneID) {
	f.zoneID = z
}

func (f *Forwarder) Close() {
	f.zidMacConn.Close()
	f.ipConn.Close()
}

func imDestination(ip, destIP net.IP, destZoneID ZoneID) bool {
	// TODO: use destZID with the ip
	return destIP.Equal(ip) || destIP.IsLoopback()
}

func imInMulticastGrp(destGrpIP net.IP) bool {
	// TODO
	return false
}

func getNextHop(destIP net.IP, ft *UniForwardTable, nt *NeighborsTable, zoneID ZoneID) (*UniForwardingEntry, bool) {
	fe, ok := ft.Get(destIP)
	if !ok {
		ne, ok := nt.Get(destIP)
		if !ok {
			return nil, false
		}
		return &UniForwardingEntry{NextHopMAC: ne.MAC, DestZoneID: zoneID}, true
	}
	return fe, true
}
