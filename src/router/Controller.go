package main

import (
	"log"
)

type ControlPacket struct {
	zidHeader *ZIDHeader
	payload   []byte
}

type Controller struct {
	router       *Router
	macConn      *MACLayerConn
	sARP         *SARP
	inputChannel chan *ControlPacket
}

func NewController(router *Router) (*Controller, error) {
	macConn, err := NewMACLayerConn(router.iface)
	if err != nil {
		return nil, err
	}

	c := make(chan *ControlPacket)

	sARP, err := NewSARP(router)
	if err != nil {
		log.Fatal("failed to build initiate sARP, err: ", err)
	}

	log.Println("initalized controller")

	return &Controller{
		router:       router,
		macConn:      macConn,
		inputChannel: c,
		sARP:         sARP,
	}, nil
}

func (c *Controller) ListenForControlPackets() {
	log.Println("Controller started listening for control packets from the forwarder")

	for {
		controlPacket := <-c.inputChannel

		switch controlPacket.zidHeader.packetType {
		case SARPReq:
			c.sARP.onSRPReq(controlPacket.payload)
		case SARPRes:
			c.sARP.onSRPRes(controlPacket.payload)
		}
	}
}

func (c *Controller) runSARP() {
	c.sARP.run()
}
