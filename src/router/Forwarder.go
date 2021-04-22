package main

import (
	"bytes"
	"log"
	"net"

	"github.com/AkihiroSuda/go-netfilter-queue"
	"github.com/mdlayher/ethernet"
)

type Forwarder struct {
	router         *Router
	zidMacConn     *MACLayerConn
	ipMacConn      *MACLayerConn
	ipConn         *IPLayerConn
	uniForwTable   *UniForwardTable
	multiForwTable *MultiForwardTable
	neighborsTable *NeighborsTable
}

func NewForwarder(router *Router, neighborsTable *NeighborsTable) (*Forwarder, error) {
	// connect to mac layer for ZID packets
	zidMacConn, err := NewMACLayerConn(router.iface, ZIDEtherType)
	if err != nil {
		return nil, err
	}

	// connect to mac layer for multicast IP packets
	ipMacConn, err := NewMACLayerConn(router.iface, IPv4EtherType)
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
		router:         router,
		zidMacConn:     zidMacConn,
		ipMacConn:      ipMacConn,
		ipConn:         ipConn,
		uniForwTable:   uniForwTable,
		neighborsTable: neighborsTable,
		multiForwTable: multiForwTable,
	}, nil
}

func (f *Forwarder) Start(inputChannel chan *UnicastControlPacket) {
	go f.forwardFromIPLayer()
	go f.forwardZIDFromMACLayer(inputChannel)
	go f.forwardIPFromMACLayer()
}

func (f *Forwarder) broadcastDummy() {
	dummy := []byte("Dummy")
	zid := &ZIDHeader{zLen: 1, packetType: LSRFloodPacket, srcZID: 2, dstZID: 3}
	packet := append(zid.MarshalBinary(), dummy...)

	encryptedPacket, err := f.router.msec.Encrypt(packet)
	if err != nil {
		log.Fatal("failed to encrypt packet, err: ", err)
	}

	err = f.zidMacConn.Write(encryptedPacket, ethernet.Broadcast)
	if err != nil {
		log.Fatal("failed to write to the device driver: ", err)
	}

	log.Println("Broadcasting dummy control packet..")
}

// forwardZIDFromMACLayer continuously receives messages from the interface,
// then either repeats it over loopback (if this is destination), or forwards it for another node.
// The messages may be up to the interface's MTU in size.
func (f *Forwarder) forwardZIDFromMACLayer(controllerChannel chan *UnicastControlPacket) {
	log.Println("started receiving from MAC layer")

	for {
		packet, err := f.zidMacConn.Read()
		if err != nil {
			log.Fatal("failed to read from interface, err: ", err)
		}
		// TODO: speed up by goroutine workers

		// decrypt and verify
		pd, err := f.router.msec.NewPacketDecrypter(packet)
		if err != nil {
			log.Fatal("failed to build packet decrypter, err: ", err)
		}
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

		if imDestination(f.router.ip, ip.DestIP, zid.dstZID) { // i'm destination,
			packet, err := pd.DecryptAll()
			if err != nil {
				continue
			}

			ippacket := packet[ZIDHeaderLen:]

			// receive message by injecting it in loopback
			err = f.ipConn.Write(ippacket)
			if err != nil {
				log.Fatal("failed to write to lo interface: ", err)
			}
		} else { // i'm a forwarder
			IPv4DecrementTTL(packet[ZIDHeaderLen:])

			e, ok := getNextHop(ip.DestIP, f.uniForwTable, f.neighborsTable)
			if !ok {
				// TODO: call controller
				continue
			}

			// hand it directly to the interface
			err = f.zidMacConn.Write(packet, e.NextHopMAC)
			if err != nil {
				log.Fatal("failed to write to the interface: ", err)
			}
		}
	}
}

// forwardIPFromMACLayer continuously receives messages from the interface,
// then either repeats it over loopback (if this is destination), or forwards it for another node.
// The messages may be up to the interface's MTU in size.
func (f *Forwarder) forwardIPFromMACLayer() {
	log.Println("started receiving from MAC layer")

	for {
		packet, err := f.ipMacConn.Read()
		if err != nil {
			log.Fatal("failed to read from interface, err: ", err)
		}
		// TODO: speed up by goroutine workers

		// decrypt and verify
		pd, err := f.router.msec.NewPacketDecrypter(packet)
		if err != nil {
			log.Fatal("failed to build packet decrypter, err: ", err)
		}
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
				log.Fatal("failed to write to lo interface: ", err)
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
			err = f.zidMacConn.Write(packet, es.NextHopMACs[i])
			if err != nil {
				log.Fatal("failed to write to the device driver: ", err)
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
			log.Fatal("ip header must have been valid from ip layer!")
		}

		if IsInjectedPacket(packet) || imDestination(f.router.ip, ip.DestIP, 0) {
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
	e, ok := getNextHop(destIP, f.uniForwTable, f.neighborsTable)
	if !ok {
		// TODO: call controller
		return
	}

	// TODO: put this zone id, and zlen
	zid := &ZIDHeader{zLen: 1, packetType: DataPacket, srcZID: 2, dstZID: int32(e.DestZoneID)}

	// build packet
	buffer := bytes.NewBuffer(make([]byte, 0, f.router.iface.MTU))
	buffer.Write(zid.MarshalBinary())
	buffer.Write(packet)

	// encrypt
	encryptedPacket, err := f.router.msec.Encrypt(buffer.Bytes())
	if err != nil {
		log.Fatal("failed to encrypt packet, err: ", err)
	}

	// write to device driver
	err = f.zidMacConn.Write(encryptedPacket, e.NextHopMAC)
	if err != nil {
		log.Fatal("failed to write to the device driver: ", err)
	}
}

func (f *Forwarder) sendMulticast(packet []byte, grpIP net.IP) {
	es, ok := f.multiForwTable.Get(grpIP)
	if !ok {
		// TODO: call controller
		return
	}

	// encrypt
	encryptedPacket, err := f.router.msec.Encrypt(packet)
	if err != nil {
		log.Fatal("failed to encrypt packet, err: ", err)
	}

	// write to device driver
	for i := 0; i < len(es.NextHopMACs); i++ {
		err = f.zidMacConn.Write(encryptedPacket, es.NextHopMACs[i])
		if err != nil {
			log.Fatal("failed to write to the device driver: ", err)
		}
	}
}

func (f *Forwarder) sendBroadcast(packet []byte) {
	// encrypt
	encryptedPacket, err := f.router.msec.Encrypt(packet)
	if err != nil {
		log.Fatal("failed to encrypt packet, err: ", err)
	}

	// write to device driver
	// TODO: for now ethernet broadcast
	err = f.zidMacConn.Write(encryptedPacket, ethernet.Broadcast)
	if err != nil {
		log.Fatal("failed to write to the device driver: ", err)
	}
}

func (f *Forwarder) Close() {
	f.zidMacConn.Close()
	f.ipConn.Close()
}

func imDestination(ip, destIP net.IP, destZoneID int32) bool {
	// TODO: use destZID with the ip
	return destIP.Equal(ip) || destIP.IsLoopback()
}

func imInMulticastGrp(destGrpIP net.IP) bool {
	// TODO
	return false
}

func getNextHop(destIP net.IP, ft *UniForwardTable, nt *NeighborsTable) (*UniForwardingEntry, bool) {
	fe, ok := ft.Get(destIP)
	if !ok {
		ne, ok := nt.Get(destIP)
		if !ok {
			return nil, false
		}
		return &UniForwardingEntry{NextHopMAC: ne.MAC, DestZoneID: 0 /*TODO: replace with this zoneid*/}, ok
	}
	return fe, true
}
