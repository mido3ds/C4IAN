package main

import (
	"log"
	"time"
)

type UnicastControlPacket struct {
	zidHeader *ZIDHeader
	payload   []byte
}

type UnicastController struct {
	router                   *Router
	macConn                  *MACLayerConn
	flooder                  *ZoneFlooder
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

	flooder, err := NewZoneFlooder(router)
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

func (c *UnicastController) Start(ft *UniForwardTable) {
	go c.ListenForControlPackets()
	go c.listenNeighChanges()

	time.AfterFunc(10*time.Second, func() {
		c.lsr.UpdateForwardingTable(c.router.ip, ft, c.neighborsTable)
		log.Println(ft)
	})
}

func (c *UnicastController) ListenForControlPackets() {
	log.Println("UnicastController started listening for control packets from the forwarder")
	// TODO: receive encrypted packet and packet decrypter
	for {
		controlPacket := <-c.inputChannel

		switch controlPacket.zidHeader.PacketType {
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

func (c *UnicastController) OnZoneIDChanged(z ZoneID) {
	c.flooder.OnZoneIDChanged(z)
}
