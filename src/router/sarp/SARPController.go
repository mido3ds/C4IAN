package sarp

import (
	"bytes"
	"log"
	"net"
	"time"

	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/tables"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

const (
	sARPHoldTime = time.Second     // Time allowed for sARP responses to arrive and neighborhood table to be updated
	sARPDelay    = 3 * time.Second // Time between consequent sARP requests (neighborhood discoveries)
)

type SARPController struct {
	msec                     *MSecLayer
	macConn                  *MACLayerConn
	NeighborsTable           *NeighborsTable
	dirtyNeighborsTable      *NeighborsTable
	NeighborhoodUpdateSignal chan bool
	myIP                     net.IP
	myMAC                    net.HardwareAddr
}

func NewSARPController(ip net.IP, iface *net.Interface, msec *MSecLayer) (*SARPController, error) {
	macConn, err := NewMACLayerConn(iface, SARPEtherType)
	if err != nil {
		return nil, err
	}

	dirtyNeighborsTable := NewNeighborsTable()
	neighborsTable := NewNeighborsTable()

	log.Println("initalized sARP controller")

	return &SARPController{
		msec:                     msec,
		macConn:                  macConn,
		NeighborsTable:           neighborsTable,
		dirtyNeighborsTable:      dirtyNeighborsTable,
		NeighborhoodUpdateSignal: make(chan bool),
		myIP:                     ip,
		myMAC:                    iface.HardwareAddr,
	}, nil
}

func (s *SARPController) Start() {
	go s.sendMsgs()
	go s.receiveMsgs()
}

func (s *SARPController) sendMsgs() {
	tableHash := s.NeighborsTable.GetTableHash()
	for {
		// TODO: Replace with scheduling if necessary
		time.Sleep(sARPDelay - sARPHoldTime)

		// Create a new table to collect sARP responses
		s.dirtyNeighborsTable = NewNeighborsTable()

		// Broadcast sARP request
		s.macConn.Write(s.CreateSARPPacket(SARPReq), BroadcastMACAddr)

		// Wait for sARP responses (collected in dirtyNeighborsTable)
		time.Sleep(sARPHoldTime)

		// Update NeighborsTable
		// Shallow copy the forwarding table, this will make the hashmap pointer in s.NeighborsTable
		// point to the new hashmap inside s.dirtyNeighborsTable. The old hashmap in s.NeighborsTable
		// will be deleted by the garbage collector
		*s.NeighborsTable = *s.dirtyNeighborsTable

		// Check if the new table contains new data
		newTableHash := s.NeighborsTable.GetTableHash()
		if !bytes.Equal(tableHash, newTableHash) {
			s.NeighborhoodUpdateSignal <- true
			log.Println(s.NeighborsTable)
		}
		tableHash = newTableHash
	}
}

func (s *SARPController) receiveMsgs() {
	for {
		packet := s.macConn.Read()

		zidHeader, ok := UnmarshalZIDHeader(s.msec.Decrypt(packet[:ZIDHeaderLen]))
		if !ok {
			log.Panicln("Received sARP Packet with invalid ZID header")
		}

		sarpHeader, ok := UnmarshalSARPHeader(s.msec.Decrypt(packet[ZIDHeaderLen : ZIDHeaderLen+sARPTotalLen]))
		if !ok {
			log.Panicln("Received sARP Packet with invalid sARP header")
		}

		// TODO: Handle zones of different sizes
		// Construct NodeID based on whether the neighbor is in the same zone or not
		var nodeID NodeID
		srcZone := &Zone{ID: zidHeader.SrcZID, Len: zidHeader.ZLen}
		if MyZone().Equal(srcZone) {
			nodeID = ToNodeID(sarpHeader.IP)
		} else {
			nodeID = ToNodeID(srcZone.ID)
		}

		// Calculate the delay, which is the link cost in the topology
		delay := time.Since(time.Unix(0, sarpHeader.sendTime))
		switch sarpHeader.Type {
		case SARPReq:
			// Update neighbors table
			s.NeighborsTable.Set(nodeID, &NeighborEntry{MAC: sarpHeader.MAC, Cost: uint16(delay.Microseconds())})
			// Send sARP response to the request sender
			s.macConn.Write(s.CreateSARPPacket(SARPRes), sarpHeader.MAC)
		case SARPRes:
			// Update dirty neighbors table
			s.dirtyNeighborsTable.Set(nodeID, &NeighborEntry{MAC: sarpHeader.MAC, Cost: uint16(delay.Microseconds())})
		}
	}
}

func (s *SARPController) CreateSARPPacket(packetType SARPType) []byte {
	zidHeader := MyZIDHeader(0)
	encryptedZIDHeader := s.msec.Encrypt(zidHeader.MarshalBinary())
	mySarpHeader := &SARPHeader{packetType, s.myIP, s.myMAC, time.Now().UnixNano()}
	encryptedSARPHeader := s.msec.Encrypt(mySarpHeader.MarshalBinary())
	return append(encryptedZIDHeader, encryptedSARPHeader...)
}

func (s *SARPController) Close() {
	s.macConn.Close()
}
