package main

import (
	"log"
	"net"
	"time"

	"github.com/mdlayher/ethernet"
)

type JoinType int

const (
	JoinQueryType JoinType = iota
	JoinReplyType
)

type MulticastControlPacket struct {
	packetType JoinType
	payload    []byte
}

type MulticastController struct {
	router       *Router
	macConn      *MACLayerConn
	flooder      *Flooder
	inputChannel chan *MulticastControlPacket
}

func (c *MulticastController) floodDummy() {
	dummy := []byte("Dummy")
	c.flooder.Flood(dummy)
}

func NewMulticastController(router *Router, neighborsTable *NeighborsTable, neighborhoodUpdateSignal chan bool) (*MulticastController, error) {
	macConn, err := NewMACLayerConn(router.iface, uint16(ethernet.EtherTypeIPv4))
	if err != nil {
		return nil, err
	}

	c := make(chan *MulticastControlPacket)

	flooder, err := NewFlooder(router)
	if err != nil {
		log.Fatal("failed to initiate flooder, err: ", err)
	}

	log.Println("initalized Multicast controller")

	return &MulticastController{
		router:       router,
		macConn:      macConn,
		inputChannel: c,
		flooder:      flooder,
	}, nil
}

func (c *MulticastController) Start(ft *UniForwardTable) {
	go c.ListenForControlPackets()

	time.AfterFunc(10*time.Second, func() {
		log.Println(ft)
	})
}

func (c *MulticastController) HandleMulticastControlPacket(srcIP net.IP, payload []byte) {
	jq, valid := UnmarshalJoinQuery(payload)
	if !valid {
		log.Panicln("Corrupted LSR packet received")
	}
	log.Println(jq)
}

func (c *MulticastController) ListenForControlPackets() {
	log.Println("MulticastController started listening for control packets from the forwarder")
	// TODO: receive encrypted packet and packet decrypter
	for {
		controlPacket := <-c.inputChannel

		switch controlPacket.packetType {
		case JoinQueryType:
			c.flooder.ReceiveFloodedMsg(controlPacket.payload, c.HandleMulticastControlPacket)
		}
	}
}
