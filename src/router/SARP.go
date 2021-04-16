package main

import (
	"bytes"
	"log"
	"net"
	"time"

	"github.com/mdlayher/ethernet"
	"golang.org/x/crypto/sha3"
)

const (
	sARPDelay     = 5 * time.Second
	hashLen       = 64 // bytes at the end
	SARPHeaderLen = 10 // excluding the hash at the end
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

func (s *SARP) run() {
	for {
		s.sendSARPReq()

		// TODO: Replace with scheduling if necessary
		time.Sleep(sARPDelay)
	}
}

func (s *SARP) OnSARPReq(payload []byte) {
	if !verifySARPHeader(payload) {
		log.Println("received malformed SARP header, ignore it")
	} else {
		ip := net.IP(payload[:4])
		mac := net.HardwareAddr(payload[4:10])
		s.neighborsTable.Set(ip, &NeighborEntry{MAC: mac})
		s.sendSARPRes(mac)
	}
}

func (s *SARP) OnSARPRes(payload []byte) {
	if !verifySARPHeader(payload) {
		log.Println("received malformed SARP header, ignore it")
	} else {
		ip := net.IP(payload[:4])
		mac := net.HardwareAddr(payload[4:10])
		s.neighborsTable.Set(ip, &NeighborEntry{MAC: mac})
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
	payload = append(payload, hash(payload)...)

	// add ZID Header
	zid := &ZIDHeader{zLen: 1, packetType: packetType, srcZID: 2, dstZID: 3}
	packet := append(zid.MarshalBinary(), payload...)

	encryptedPacket, err := s.router.msec.Encrypt(packet)
	if err != nil {
		log.Fatal("failed to encrypt packet, err: ", err)
	}

	s.macConn.Write(encryptedPacket, dst)
}

func hash(b []byte) []byte {
	h := sha3.New512()

	n, err := h.Write(b)
	if err != nil {
		log.Fatal("failed to hash, err: ", err)
	} else if n != len(b) {
		log.Fatal("failed to hash")
	}

	return h.Sum(nil)
}

func verifySARPHeader(b []byte) bool {
	if len(b) < SARPHeaderLen+hashLen {
		return false
	}

	h := hash(b[:SARPHeaderLen])
	h2 := b[SARPHeaderLen : SARPHeaderLen+hashLen]
	return bytes.Compare(h, h2) == 0
}
