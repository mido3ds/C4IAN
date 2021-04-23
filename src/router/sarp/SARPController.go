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
	sARPDelay     = 5 * time.Second // Time between consequent sARP requests (neighborhood discoveries)
	hashLen       = 64              // bytes at the end
	sARPHeaderLen = 18              // excluding the hash at the end
	sARPTotalLen  = sARPHeaderLen + hashLen
)

type SARPController struct {
	msec                     *MSecLayer
	reqMacConn               *MACLayerConn
	resMacConn               *MACLayerConn
	NeighborsTable           *NeighborsTable
	NeighborhoodUpdateSignal chan bool
	myIP                     net.IP
	myMAC                    net.HardwareAddr
}

func NewSARPController(ip net.IP, iface *net.Interface, msec *MSecLayer) (*SARPController, error) {
	reqMacConn, err := NewMACLayerConn(iface, SARPReqEtherType)
	if err != nil {
		return nil, err
	}
	resMacConn, err := NewMACLayerConn(iface, SARPResEtherType)
	if err != nil {
		return nil, err
	}

	neighborsTable := NewNeighborsTable()

	log.Println("initalized sARP controller")

	return &SARPController{
		msec:                     msec,
		reqMacConn:               reqMacConn,
		resMacConn:               resMacConn,
		NeighborsTable:           neighborsTable,
		NeighborhoodUpdateSignal: make(chan bool),
		myIP:                     ip,
		myMAC:                    iface.HardwareAddr,
	}, nil
}

func (s *SARPController) Start() {
	go s.sendMsgs()
	go s.recvRequests()
	go s.recvResponses()
}

func (s *SARPController) sendMsgs() {
	tableHash := s.NeighborsTable.GetTableHash()
	for {
		// TODO: Neighbors table may be used while it is cleared, find another way to remove gone neighbors
		s.NeighborsTable.Clear()

		// broadcast request
		header := &SARPHeader{s.myIP, s.myMAC, time.Now().UnixNano()}
		s.reqMacConn.Write(s.msec.Encrypt(header.MarshalBinary()), BroadcastMACAddr)

		time.Sleep(sARPHoldTime)
		newTableHash := s.NeighborsTable.GetTableHash()

		if !bytes.Equal(tableHash, newTableHash) {
			s.NeighborhoodUpdateSignal <- true
		}

		// TODO: Replace with scheduling if necessary
		time.Sleep(sARPDelay - sARPHoldTime)
		tableHash = newTableHash
	}
}

func (s *SARPController) recvRequests() {
	for {
		packet := s.reqMacConn.Read()
		packet = s.msec.Decrypt(packet[:sARPTotalLen])

		if header, ok := UnmarshalSARPHeader(packet); ok {
			// store it
			delay := time.Since(time.Unix(0, header.sendTime))
			s.NeighborsTable.Set(header.IP, &NeighborEntry{MAC: header.MAC, Cost: uint16(delay.Microseconds())})

			// unicast response
			header := &SARPHeader{s.myIP, s.myMAC, time.Now().UnixNano()}
			s.resMacConn.Write(s.msec.Encrypt(header.MarshalBinary()), header.MAC)
		}
	}
}

func (s *SARPController) recvResponses() {
	for {
		packet := s.resMacConn.Read()
		packet = s.msec.Decrypt(packet[:sARPTotalLen])

		if header, ok := UnmarshalSARPHeader(packet); ok {
			// store it
			delay := time.Since(time.Unix(0, header.sendTime))
			s.NeighborsTable.Set(header.IP, &NeighborEntry{MAC: header.MAC, Cost: uint16(delay.Microseconds())})
		}
	}
}

func (s *SARPController) Close() {
	s.reqMacConn.Close()
	s.resMacConn.Close()
}
