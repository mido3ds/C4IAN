package main

import (
	"log"
)

type UnicastControlPacket struct {
	zidHeader *ZIDHeader
	payload   []byte
}

type UnicastController struct {
	router                   *Router
	macConn                  *MACLayerConn
	flooder                  *Flooder
	lsr                      *LSR
	inputChannel             chan *UnicastControlPacket
	neighborhoodUpdateSignal chan bool
	neighborsTable           *NeighborsTable
}

func (c *UnicastController) floodDummy() {
	dummy := []byte("Dummy")
	c.flooder.Flood(dummy)
}

func NewUnicastController(router *Router, neighborsTable *NeighborsTable, neighborhoodUpdateSignal chan bool) (*UnicastController, error) {
	macConn, err := NewMACLayerConn(router.iface, ZIDEtherType)
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
		router:                   router,
		macConn:                  macConn,
		inputChannel:             c,
		flooder:                  flooder,
		lsr:                      lsr,
		neighborhoodUpdateSignal: neighborhoodUpdateSignal,
		neighborsTable:           neighborsTable,
	}, nil
}

func (c *UnicastController) Start() {
	go c.ListenForControlPackets()
	go c.listenNeighChanges()
}

func (c *UnicastController) ListenForControlPackets() {
	log.Println("UnicastController started listening for control packets from the forwarder")
	// TODO: receive encrypted packet and packet decrypter
	for {
		controlPacket := <-c.inputChannel

		switch controlPacket.zidHeader.packetType {
		case LSRFloodPacket:
			c.flooder.ReceiveFloodedMsg(controlPacket.payload, c.lsr.HandleLSRPacket)
		}
	}
}

func (c *UnicastController) listenNeighChanges() {
	for {
		<-c.neighborhoodUpdateSignal
		c.lsr.topology.Update(c.router.ip, c.neighborsTable)
		c.lsr.SendLSRPacket(c.flooder, c.neighborsTable)
	}
}
