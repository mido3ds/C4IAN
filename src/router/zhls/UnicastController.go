package zhls

import (
	"fmt"
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/flood"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/tables"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

type UnicastController struct {
	ip                           net.IP
	flooder                      *ZoneFlooder
	lsr                          *LSR
	neighborhoodUpdateSignal     chan bool
	neighborsTable               *NeighborsTable
	UpdateUnicastForwardingTable func(ft *UniForwardTable)
}

func NewUnicastController(iface *net.Interface, ip net.IP, neighborsTable *NeighborsTable, neighborhoodUpdateSignal chan bool, msec *MSecLayer, zlen byte) (*UnicastController, error) {
	flooder, err := NewZoneFlooder(iface, ip, msec, zlen)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate flooder, err: %#v", err)
	}

	lsr := NewLSR(ip, neighborsTable)

	log.Println("initalized controller")

	return &UnicastController{
		ip:                           ip,
		flooder:                      flooder,
		lsr:                          lsr,
		neighborhoodUpdateSignal:     neighborhoodUpdateSignal,
		neighborsTable:               neighborsTable,
		UpdateUnicastForwardingTable: lsr.UpdateForwardingTable,
	}, nil
}

func (c *UnicastController) Start() {
	go c.listenForNeighborhoodChanges()
	go c.flooder.ListenForFloodedMsgs(c.lsr.HandleLSRPacket)
}

func (c *UnicastController) listenForNeighborhoodChanges() {
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
}
