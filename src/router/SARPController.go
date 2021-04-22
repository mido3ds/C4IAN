package main

import (
	"bytes"
	"log"
	"net"
	"time"

	"github.com/mdlayher/ethernet"
)

const (
	sARPHoldTime  = time.Second     // Time allowed for sARP responses to arrive and neighborhood table to be updated
	sARPDelay     = 5 * time.Second // Time between consequent sARP requests (neighborhood discoveries)
	hashLen       = 64              // bytes at the end
	sARPHeaderLen = 10              // excluding the hash at the end
	sARPTotalLen  = sARPHeaderLen + hashLen

	// Make use of an unassigned EtherType to differentiate between SARP traffic and other traffic
	// https://www.iana.org/assignments/ieee-802-numbers/ieee-802-numbers.xhtml
	sARPReqEtherType = 0x0809
	sARPResEtherType = 0x080A
)

type SARPController struct {
	msec                     *MSecLayer
	reqMacConn               *MACLayerConn
	resMacConn               *MACLayerConn
	neighborsTable           *NeighborsTable
	neighborhoodUpdateSignal chan bool
	encryptedHdr             []byte
}

func NewSARPController(router *Router) (*SARPController, error) {
	reqMacConn, err := NewMACLayerConn(router.iface, sARPReqEtherType)
	if err != nil {
		return nil, err
	}
	resMacConn, err := NewMACLayerConn(router.iface, sARPResEtherType)
	if err != nil {
		return nil, err
	}

	neighborsTable := NewNeighborsTable()

	header := &SARPHeader{router.ip, router.iface.HardwareAddr}
	encryptedHdr, err := router.msec.Encrypt(header.MarshalBinary())
	if err != nil {
		log.Fatal("failed to encrypt packet, err: ", err)
	}

	log.Println("initalized sARP controller")

	return &SARPController{
		msec:                     router.msec,
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
		log.Println("Sending sARP request")
		s.neighborsTable.Clear()

		// broadcast request
		if err := s.reqMacConn.Write(s.encryptedHdr, ethernet.Broadcast); err != nil {
			log.Fatal("failed to write to device driver, err: ", err)
		}

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
		packet, err := s.reqMacConn.Read()
		if err != nil {
			log.Fatal("couldn't read from device driver, err: ", err)
		}

		packet, err = s.msec.Decrypt(packet[:sARPTotalLen])
		if err != nil {
			log.Fatal("couldn't decrypt msg, err: ", err)
		}

		if header, ok := UnmarshalSARPHeader(packet); ok {
			// store it
			s.neighborsTable.Set(header.ip, &NeighborEntry{MAC: header.mac, cost: 1})

			// unicast response
			if err := s.resMacConn.Write(s.encryptedHdr, header.mac); err != nil {
				log.Fatal("failed to write to device driver, err: ", err)
			}
		}
	}
}

func (s *SARPController) recvResponses() {
	for {
		packet, err := s.resMacConn.Read()
		if err != nil {
			log.Fatal("couldn't read from device driver, err: ", err)
		}

		packet, err = s.msec.Decrypt(packet[:sARPTotalLen])
		if err != nil {
			log.Fatal("couldn't decrypt msg, err: ", err)
		}

		if header, ok := UnmarshalSARPHeader(packet); ok {
			// store it
			s.neighborsTable.Set(header.ip, &NeighborEntry{MAC: header.mac, cost: 1})
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
	buf.Write(Hash_SHA3(buf.Bytes()[:sARPHeaderLen]))

	return buf.Bytes()
}

func verifySARPHeader(b []byte) bool {
	if len(b) < sARPTotalLen {
		return false
	}

	return verifyHash_SHA3(b[:sARPHeaderLen], b[sARPHeaderLen:sARPTotalLen])
}
