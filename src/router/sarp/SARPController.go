package sarp

import (
	"bytes"
	"log"
	"net"
	"time"

	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/tables"
)

const (
	sARPHoldTime  = time.Second     // Time allowed for sARP responses to arrive and neighborhood table to be updated
	sARPDelay     = 3 * time.Second // Time between consequent sARP requests (neighborhood discoveries)
	hashLen       = 64              // bytes at the end
	sARPHeaderLen = 19              // excluding the hash at the end
	sARPTotalLen  = sARPHeaderLen + hashLen
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
		header := &SARPHeader{SARPReq, s.myIP, s.myMAC, time.Now().UnixNano()}
		s.macConn.Write(s.msec.Encrypt(header.MarshalBinary()), BroadcastMACAddr)

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
		}
		tableHash = newTableHash
	}
}

func (s *SARPController) receiveMsgs() {
	for {
		packet := s.macConn.Read()
		packet = s.msec.Decrypt(packet[:sARPTotalLen])

		header, ok := UnmarshalSARPHeader(packet)
		if !ok {
			log.Panicln("Invalid sARP packet received")
		}
		switch header.Type {
		case SARPReq:
			delay := time.Since(time.Unix(0, header.sendTime))
			s.NeighborsTable.Set(header.IP, &NeighborEntry{MAC: header.MAC, Cost: uint16(delay.Microseconds())})

			myHeader := &SARPHeader{SARPRes, s.myIP, s.myMAC, time.Now().UnixNano()}
			s.macConn.Write(s.msec.Encrypt(myHeader.MarshalBinary()), header.MAC)
		case SARPRes:
			delay := time.Since(time.Unix(0, header.sendTime))
			s.dirtyNeighborsTable.Set(header.IP, &NeighborEntry{MAC: header.MAC, Cost: uint16(delay.Microseconds())})
		}
	}
}

func (s *SARPController) Close() {
	s.macConn.Close()
}
