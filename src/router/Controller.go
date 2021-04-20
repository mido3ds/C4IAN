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
	flooder      *Flooder
	lsr          *LSR
	inputChannel chan *ControlPacket
}

func (c *Controller) floodDummy() {
	dummy := []byte("Dummy")
	c.flooder.Flood(dummy)
}

func NewController(router *Router) (*Controller, error) {
	macConn, err := NewMACLayerConn(router.iface)
	if err != nil {
		return nil, err
	}

	c := make(chan *ControlPacket)

	flooder, err := NewFlooder(router)
	if err != nil {
		log.Fatal("failed to initiate flooder, err: ", err)
	}

	sARP, err := NewSARP(router)
	if err != nil {
		log.Fatal("failed to initiate sARP, err: ", err)
	}

	lsr := NewLSR()

	log.Println("initalized controller")

	return &Controller{
		router:       router,
		macConn:      macConn,
		inputChannel: c,
		sARP:         sARP,
		flooder:      flooder,
		lsr:          lsr,
	}, nil
}

func (c *Controller) ListenForControlPackets() {
	log.Println("Controller started listening for control packets from the forwarder")
	// TODO: receive encrypted packet and packet decrypter
	for {
		controlPacket := <-c.inputChannel

		switch controlPacket.zidHeader.packetType {
		case SARPReq:
			c.sARP.OnSARPReq(controlPacket.payload)
		case SARPRes:
			c.sARP.OnSARPRes(controlPacket.payload)
		case LSRFloodPacket:
			c.flooder.ReceiveFloodedMsg(controlPacket.payload, c.lsr.HandleLSRPacket)
		}
	}
}

func (c *Controller) runSARP() {
	onNeighborhoodChange := func() {
		c.lsr.SendLSRPacket(c.flooder, c.sARP.neighborsTable)
	}
	c.sARP.run(onNeighborhoodChange)
}
