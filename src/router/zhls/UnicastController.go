package zhls

import (
	"fmt"
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/flood"
	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/tables"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

type UnicastControlPacket struct {
	ZIDHeader *ZIDHeader
	Payload   []byte
}

type UnicastController struct {
	ip                       net.IP
	macConn                  *MACLayerConn
	flooder                  *ZoneFlooder
	lsr                      *LSR
	InputChannel             chan *UnicastControlPacket
	neighborhoodUpdateSignal chan bool
	neighborsTable           *NeighborsTable
}

func NewUnicastController(iface *net.Interface, ip net.IP, neighborsTable *NeighborsTable, neighborhoodUpdateSignal chan bool, msec *MSecLayer, zlen byte) (*UnicastController, error) {
	macConn, err := NewMACLayerConn(iface, ZIDEtherType)
	if err != nil {
		return nil, err
	}

	c := make(chan *UnicastControlPacket)

	flooder, err := NewZoneFlooder(iface, ip, msec, zlen)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate flooder, err: %#v", err)
	}

	lsr := NewLSR()

	log.Println("initalized controller")

	return &UnicastController{
		macConn:                  macConn,
		ip:                       ip,
		InputChannel:             c,
		flooder:                  flooder,
		lsr:                      lsr,
		neighborhoodUpdateSignal: neighborhoodUpdateSignal,
		neighborsTable:           neighborsTable,
	}, nil
}

func (c *UnicastController) Start(ft *UniForwardTable) {
	go c.listenForControlPackets()
	go c.listenNeighChanges()

	// time.AfterFunc(10*time.Second, func() {
	// 		c.lsr.UpdateForwardingTable(c.ip, ft, c.neighborsTable)
	// 		log.Println(ft)
	// })
}

func (c *UnicastController) listenForControlPackets() {
	log.Println("UnicastController started listening for control packets from the forwarder")
	// TODO: receive encrypted packet and packet decrypter
	for {
		controlPacket := <-c.InputChannel

		switch controlPacket.ZIDHeader.PacketType {
		case LSRFloodPacket:
			c.flooder.ReceiveFloodedMsg(controlPacket.Payload, c.lsr.HandleLSRPacket)
		}
	}
}

func (c *UnicastController) listenNeighChanges() {
	for {
		<-c.neighborhoodUpdateSignal
		c.lsr.topology.Update(c.ip, c.neighborsTable)
		c.lsr.SendLSRPacket(c.flooder, c.neighborsTable)
	}
}

func (c *UnicastController) OnZoneIDChanged(z ZoneID) {
	c.flooder.OnZoneIDChanged(z)
}

func (c *UnicastController) Close() {
	c.flooder.Close()
	c.macConn.Close()
}
