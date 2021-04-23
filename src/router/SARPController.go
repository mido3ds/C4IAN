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
	sARPHeaderLen = 18              // excluding the hash at the end
	sARPTotalLen  = sARPHeaderLen + hashLen
)

type SARPController struct {
	msec                     *MSecLayer
	reqMacConn               *MACLayerConn
	resMacConn               *MACLayerConn
	neighborsTable           *NeighborsTable
	neighborhoodUpdateSignal chan bool
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
		neighborsTable:           neighborsTable,
		neighborhoodUpdateSignal: make(chan bool),
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
	tableHash := s.neighborsTable.GetTableHash()
	for {
		s.neighborsTable.Clear()

		// broadcast request
		header := &SARPHeader{s.myIP, s.myMAC, time.Now().UnixNano()}
		s.reqMacConn.Write(s.msec.Encrypt(header.MarshalBinary()), BroadcastMACAddr)

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
			delay := time.Since(time.Unix(0, header.sendTime))
			s.neighborsTable.Set(header.IP, &NeighborEntry{MAC: header.MAC, Cost: uint16(delay.Microseconds())})

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
			s.neighborsTable.Set(header.IP, &NeighborEntry{MAC: header.MAC, Cost: uint16(delay.Microseconds())})
		}
	}
}

func (s *SARPController) Close() {
	s.reqMacConn.Close()
	s.resMacConn.Close()
}

type SARPHeader struct {
	IP       net.IP
	MAC      net.HardwareAddr
	sendTime int64
}

func UnmarshalSARPHeader(packet []byte) (*SARPHeader, bool) {
	ok := verifySARPHeader(packet)
	if !ok {
		return nil, false
	}
	// sendTime -> packet[10:18]
	sendTime := int64(packet[10])<<56 | int64(packet[11])<<48 | int64(packet[12])<<40 | int64(packet[13])<<32 |
		int64(packet[14])<<24 | int64(packet[15])<<16 | int64(packet[16])<<8 | int64(packet[17])
	return &SARPHeader{
		IP:       net.IP(packet[:4]),
		MAC:      net.HardwareAddr(packet[4:10]),
		sendTime: sendTime,
	}, true
}

func (s *SARPHeader) MarshalBinary() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, sARPTotalLen))

	buf.Write(s.IP.To4())
	buf.Write(s.MAC)
	for i := 56; i >= 0; i -= 8 {
		buf.WriteByte(byte(s.sendTime >> i))
	}
	buf.Write(HashSHA3(buf.Bytes()[:sARPHeaderLen]))

	return buf.Bytes()
}

func verifySARPHeader(b []byte) bool {
	if len(b) < sARPTotalLen {
		return false
	}

	return VerifySHA3Hash(b[:sARPHeaderLen], b[sARPHeaderLen:sARPTotalLen])
}
