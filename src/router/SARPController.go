package main

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
	sARPHeaderLen = 10              // excluding the hash at the end
	sARPTotalLen  = sARPHeaderLen + hashLen
)

type SARPController struct {
	msec                     *MSecLayer
	reqMacConn               *MACLayerConn
	resMacConn               *MACLayerConn
	neighborsTable           *NeighborsTable
	neighborhoodUpdateSignal chan bool
	encryptedHdr             []byte
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

	header := &SARPHeader{ip, iface.HardwareAddr}
	encryptedHdr := msec.Encrypt(header.MarshalBinary())

	log.Println("initalized sARP controller")

	return &SARPController{
		msec:                     msec,
		reqMacConn:               reqMacConn,
		resMacConn:               resMacConn,
		neighborsTable:           neighborsTable,
		neighborhoodUpdateSignal: make(chan bool),
		encryptedHdr:             encryptedHdr,
	}, nil
}

func (s *SARPController) Start() {
	go s.sendMsgs()
	go s.recvRequests()
	go s.recvResponses()
}

func (s *SARPController) sendMsgs() {
	tableHash := s.neighborsTable.GetTableHash()
	for {
		s.neighborsTable.Clear()

		// broadcast request
		s.reqMacConn.Write(s.encryptedHdr, BroadcastMACAddr)

		time.Sleep(sARPHoldTime)
		newTableHash := s.neighborsTable.GetTableHash()

		if !bytes.Equal(tableHash, newTableHash) {
			s.neighborhoodUpdateSignal <- true
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
			s.neighborsTable.Set(header.ip, &NeighborEntry{MAC: header.mac, Cost: 1})

			// unicast response
			s.resMacConn.Write(s.encryptedHdr, header.mac)
		}
	}
}

func (s *SARPController) recvResponses() {
	for {
		packet := s.resMacConn.Read()
		packet = s.msec.Decrypt(packet[:sARPTotalLen])

		if header, ok := UnmarshalSARPHeader(packet); ok {
			// store it
			s.neighborsTable.Set(header.ip, &NeighborEntry{MAC: header.mac, Cost: 1})
		}
	}
}

type SARPHeader struct {
	ip  net.IP
	mac net.HardwareAddr
}

func UnmarshalSARPHeader(packet []byte) (*SARPHeader, bool) {
	ok := verifySARPHeader(packet)
	if !ok {
		return nil, false
	}

	return &SARPHeader{
		ip:  net.IP(packet[:4]),
		mac: net.HardwareAddr(packet[4:10]),
	}, true
}

func (s *SARPHeader) MarshalBinary() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, sARPTotalLen))

	buf.Write(s.ip[:net.IPv4len])
	buf.Write(s.mac)
	buf.Write(HashSHA3(buf.Bytes()[:sARPHeaderLen]))

	return buf.Bytes()
}

func verifySARPHeader(b []byte) bool {
	if len(b) < sARPTotalLen {
		return false
	}

	return VerifySHA3Hash(b[:sARPHeaderLen], b[sARPHeaderLen:sARPTotalLen])
}
