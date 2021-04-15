package main

import (
	"log"
	"net"
	"time"

	"github.com/mdlayher/ethernet"
)

const sARPDelay = 5 * time.Second

type ControlPacket struct {
	zidHeader *ZIDHeader
	payload   []byte
}

type Controller struct {
	router       *Router
	macConn      *MACLayerConn
	inputChannel chan *ControlPacket
}

func NewController(router *Router) (*Controller, error) {
	macConn, err := NewMACLayerConn(router.iface)
	if err != nil {
		return nil, err
	}

	c := make(chan *ControlPacket)

	log.Println("initalized controller")

	return &Controller{
		router:       router,
		macConn:      macConn,
		inputChannel: c,
	}, nil
}

func (c *Controller) ListenForControlPackets() {
	log.Println("Controller started listening for control packets from the forwarder")

	for {
		controlPacket := <-c.inputChannel
		log.Println("Received a control packet")

		switch controlPacket.zidHeader.packetType {
		case SARP:
			ip := net.IP(controlPacket.payload[:4])
			mac := net.HardwareAddr(controlPacket.payload[4:10])
			log.Println("Received sARP: ", ip, mac)
			c.sendSARP(mac)
		}

	}
}

func (c *Controller) sARP() {
	log.Println("Initiating sARP")
	for {
		c.sendSARP(ethernet.Broadcast)

		// TODO: Replace with scheduling if necessary
		time.Sleep(sARPDelay)
	}
}

func (c *Controller) sendSARP(dst net.HardwareAddr) {
	payload := append([]byte(c.router.ip.To4()), []byte(c.router.iface.HardwareAddr)...)
	log.Println("Sending sARP: ", c.router.ip.To4(), c.router.iface.HardwareAddr)

	zid, err := NewZIDPacketMarshaler(c.router.iface.MTU)
	if err != nil {
		log.Fatal(err)
	}

	packet, err := zid.MarshalBinary(&ZIDHeader{zLen: 1, packetType: SARP}, payload)
	if err != nil {
		log.Fatal(err)
	}

	encryptedPacket, err := c.router.msec.Encrypt(packet)
	if err != nil {
		log.Fatal("failed to encrypt packet, err: ", err)
	}

	c.macConn.Write(encryptedPacket, dst)
}
