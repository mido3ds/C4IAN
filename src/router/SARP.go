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
	SARPHeaderLen = 10              // excluding the hash at the end
)

type SARP struct {
	router         *Router
	macConn        *MACLayerConn
	neighborsTable *NeighborsTable
}

func NewSARP(router *Router) (*SARP, error) {
	macConn, err := NewMACLayerConn(router.iface)
	if err != nil {
		return nil, err
	}

	neighborsTable := NewNeighborsTable()

	log.Println("initalized sARP")

	return &SARP{
		router:         router,
		macConn:        macConn,
		neighborsTable: neighborsTable,
	}, nil
}

func (s *SARP) run(onNeighborhoodUpdate func()) {
	tableHash := s.neighborsTable.GetTableHash()
	for {
		//log.Println("Sending sARP request")
		s.neighborsTable.Clear()
		s.sendSARPReq()
		time.Sleep(sARPHoldTime)
		newTableHash := s.neighborsTable.GetTableHash()

		if !bytes.Equal(tableHash, newTableHash) {
			onNeighborhoodUpdate()
		}

		// TODO: Replace with scheduling if necessary
		time.Sleep(sARPDelay - sARPHoldTime)
		tableHash = newTableHash
	}
}

func (s *SARP) OnSARPReq(payload []byte) {
	if !verifySARPHeader(payload) {
		log.Println("received malformed SARP header, ignore it")
	} else {
		ip := net.IP(payload[:4])
		mac := net.HardwareAddr(payload[4:10])
		s.neighborsTable.Set(ip, &NeighborEntry{MAC: mac, cost: 1})
		s.sendSARPRes(mac)
	}
}

func (s *SARP) OnSARPRes(payload []byte) {
	if !verifySARPHeader(payload) {
		log.Println("received malformed SARP header, ignore it")
	} else {
		ip := net.IP(payload[:4])
		mac := net.HardwareAddr(payload[4:10])
		s.neighborsTable.Set(ip, &NeighborEntry{MAC: mac, cost: 1})
	}
}

func (s *SARP) sendSARPReq() {
	s.sendSARP(SARPReq, ethernet.Broadcast)
}

func (s *SARP) sendSARPRes(dst net.HardwareAddr) {
	s.sendSARP(SARPRes, dst)
}

func (s *SARP) sendSARP(packetType PacketType, dst net.HardwareAddr) {
	payload := append(s.router.ip, (s.router.iface.HardwareAddr)...)
	payload = append(payload, Hash_SHA3(payload)...)

	// Add ZID Header
	zid := &ZIDHeader{zLen: 1, packetType: packetType, srcZID: 2, dstZID: 3}
	packet := append(zid.MarshalBinary(), payload...)

	encryptedPacket, err := s.router.msec.Encrypt(packet)
	if err != nil {
		log.Fatal("failed to encrypt packet, err: ", err)
	}

	s.macConn.Write(encryptedPacket, dst)
}

func verifySARPHeader(b []byte) bool {
	if len(b) < SARPHeaderLen+hashLen {
		return false
	}

	return verifyHash_SHA3(b[:SARPHeaderLen], b[SARPHeaderLen:SARPHeaderLen+hashLen])
}
