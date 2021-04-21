package main

import (
	"log"
)

type UnicastControlPacket struct {
	zidHeader *ZIDHeader
	payload   []byte
}

type UnicastController struct {
	router       *Router
	macConn      *MACLayerConn
	sARP         *SARP
	flooder      *Flooder
	lsr          *LSR
	inputChannel chan *UnicastControlPacket
}

func (c *UnicastController) floodDummy() {
	dummy := []byte("Dummy")
	c.flooder.Flood(dummy)
}

func NewUnicastController(router *Router, sARP *SARP) (*UnicastController, error) {
	macConn, err := NewMACLayerConn(router.iface)
	if err != nil {
		return nil, err
	}

	c := make(chan *UnicastControlPacket)

	flooder, err := NewFlooder(router)
	if err != nil {
		log.Fatal("failed to initiate flooder, err: ", err)
	}

	lsr := NewLSR()

	log.Println("initalized controller")

	return &UnicastController{
		router:       router,
		macConn:      macConn,
		inputChannel: c,
		sARP:         sARP,
		flooder:      flooder,
		lsr:          lsr,
	}, nil
}

func (c *UnicastController) ListenForControlPackets() {
	log.Println("UnicastController started listening for control packets from the forwarder")
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

func (c *UnicastController) runSARP() {
	onNeighborhoodChange := func() {
		c.lsr.topology.Update(c.router.ip, c.sARP.neighborsTable)
		c.lsr.SendLSRPacket(c.flooder, c.sARP.neighborsTable)
	}
	c.sARP.run(onNeighborhoodChange)
}
