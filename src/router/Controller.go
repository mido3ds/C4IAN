package main

import (
	"log"
	"net"
	"time"

	"github.com/mdlayher/ethernet"
)

const sARPDelay = 20 * time.Second

type ControlPacket struct {
	zidHeader *ZIDHeader
	payload   []byte
}

type Controller struct {
	router       *Router
	macConn      *MACLayerConn
	inputChannel chan *ControlPacket
	flooder		 *Flooder
	neighborsTable *NeighborsTable
}

func (c *Controller) floodDummy() {
	dummy:= []byte{0xAA, 0xBB}
	c.flooder.Flood(dummy)
}


func NewController(router *Router) (*Controller, error) {
	macConn, err := NewMACLayerConn(router.iface)
	if err != nil {
		return nil, err
	}

	c := make(chan *ControlPacket)
	flooder, err := NewFlooder(router)

	neighborsTable := NewNeighborsTable()

	log.Println("initalized controller")

	return &Controller{
		router:       router,
		macConn:      macConn,
		inputChannel: c,
		flooder:	  flooder,	 
		neighborsTable: neighborsTable,
	}, nil
}

func (c *Controller) ListenForControlPackets() {
	log.Println("Controller started listening for control packets from the forwarder")

	for {
		controlPacket := <-c.inputChannel

		switch controlPacket.zidHeader.packetType {
		case SARPReq:
			ip := net.IP(controlPacket.payload[:4])
			mac := net.HardwareAddr(controlPacket.payload[4:10])
			log.Println("Received sARP Request from: ", ip, mac)
			c.sendSARPRes(mac)
		case SARPRes:
			ip := net.IP(controlPacket.payload[:4])
			mac := net.HardwareAddr(controlPacket.payload[4:10])
			log.Println("Received sARP: ", ip, mac)
			c.sendSARP(mac)
		case FloodPacket:
			c.flooder.receiveFlood(controlPacket.payload)
		}
	}
}

func (c *Controller) sARP() {
	log.Println("Initiating sARP")
	for {
		c.sendSARPReq()

		// TODO: Replace with scheduling if necessary
		time.Sleep(sARPDelay)
	}
}

func (c *Controller) sendSARPReq() {
	log.Print("Sending sARP Request: ")
	c.sendSARP(SARPReq, ethernet.Broadcast)
}

func (c *Controller) sendSARPRes(dst net.HardwareAddr) {
	log.Print("Sending sARP Response: ")
	c.sendSARP(SARPRes, dst)
}

func (c *Controller) sendSARP(packetType PacketType, dst net.HardwareAddr) {
	payload := append(c.router.ip, (c.router.iface.HardwareAddr)...)
	log.Println(c.router.ip, c.router.iface.HardwareAddr)

	zid, err := NewZIDPacketMarshaler(c.router.iface.MTU)
	if err != nil {
		log.Fatal(err)
	}

	packet, err := zid.MarshalBinary(&ZIDHeader{zLen: 1, packetType: packetType}, payload)
	if err != nil {
		log.Fatal(err)
	}

	encryptedPacket, err := c.router.msec.Encrypt(packet)
	if err != nil {
		log.Fatal("failed to encrypt packet, err: ", err)
	}

	c.macConn.Write(encryptedPacket, dst)
}
