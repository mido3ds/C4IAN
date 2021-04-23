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
	router          *Router
	macConn         *MACLayerConn
	grpMembersTable *GroupMembersTable

	// TODO: use global flooder
	flooder *ZoneFlooder

	// TODO: remove
	inputChannel chan *MulticastControlPacket
}

func (c *MulticastController) floodDummy() {
	dummy := []byte("Dummy")
	c.flooder.Flood(dummy)
}

func NewMulticastController(router *Router, mgroupContent string) (*MulticastController, error) {
	// TODO: create ether type for odmrp control messages
	macConn, err := NewMACLayerConn(router.iface, ethernet.EtherTypeIPv4)
	if err != nil {
		return nil, err
	}

	c := make(chan *MulticastControlPacket)

	flooder, err := NewZoneFlooder(router)
	if err != nil {
		log.Panic("failed to initiate flooder, err: ", err)
	}

	log.Println("initalized Multicast controller")

	return &MulticastController{
		router:          router,
		macConn:         macConn,
		grpMembersTable: NewGroupMembersTable(mgroupContent),
		inputChannel:    c,
		flooder:         flooder,
	}, nil
}

// TODO: create function for forwarder to call when forwarding table has no entry

func (c *MulticastController) Start(ft *MultiForwardTable) {
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

// TODO: remove
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
