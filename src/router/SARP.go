package main

import (
	"log"
	"net"
	"time"

	"github.com/mdlayher/ethernet"
)

const sARPDelay = 5 * time.Second

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

func (s *SARP) onSRPReq(payload []byte) {
	ip := net.IP(payload[:4])
	mac := net.HardwareAddr(payload[4:10])
	s.neighborsTable.Set(ip, &NeighborEntry{MAC: mac})
	s.sendSARPRes(mac)
}

func (s *SARP) onSRPRes(payload []byte) {
	ip := net.IP(payload[:4])
	mac := net.HardwareAddr(payload[4:10])
	s.neighborsTable.Set(ip, &NeighborEntry{MAC: mac})
}

func (s *SARP) sendSARPReq() {
	s.sendSARP(SARPReq, ethernet.Broadcast)
}

func (s *SARP) sendSARPRes(dst net.HardwareAddr) {
	s.sendSARP(SARPRes, dst)
}

func (s *SARP) sendSARP(packetType PacketType, dst net.HardwareAddr) {
	payload := append(s.router.ip, (s.router.iface.HardwareAddr)...)

	zid, err := NewZIDPacketMarshaler(s.router.iface.MTU)
	if err != nil {
		log.Fatal(err)
	}

	packet, err := zid.MarshalBinary(&ZIDHeader{zLen: 1, packetType: packetType}, payload)
	if err != nil {
		log.Fatal(err)
	}

	encryptedPacket, err := s.router.msec.Encrypt(packet)
	if err != nil {
		log.Fatal("failed to encrypt packet, err: ", err)
	}

	s.macConn.Write(encryptedPacket, dst)
}
