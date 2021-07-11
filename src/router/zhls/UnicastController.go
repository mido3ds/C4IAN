package zhls

import (
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/tables"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

type UnicastController struct {
	ip                           net.IP
	lsr                          *LSRController
	neighborhoodUpdateSignal     chan bool
	neighborsTable               *NeighborsTable
	UpdateUnicastForwardingTable func(ft *UniForwardTable)
}

func NewUnicastController(iface *net.Interface, ip net.IP, neighborsTable *NeighborsTable, neighborhoodUpdateSignal chan bool, msec *MSecLayer, topology *Topology) (*UnicastController, error) {
	lsr := newLSR(iface, msec, ip, neighborsTable, topology)

	log.Println("initalized controller")

	return &UnicastController{
		ip:                           ip,
		lsr:                          lsr,
		neighborhoodUpdateSignal:     neighborhoodUpdateSignal,
		neighborsTable:               neighborsTable,
		UpdateUnicastForwardingTable: lsr.updateForwardingTable,
	}, nil
}

func (c *UnicastController) Start() {
	go c.lsr.Start()
	go c.listenForNeighborhoodChanges()
}

func (c *UnicastController) listenForNeighborhoodChanges() {
	for {
		isUpdated := <-c.neighborhoodUpdateSignal
		c.lsr.sendIntrazoneLSR(isUpdated)
	}
}

func (c *UnicastController) OnZoneChange(newZoneID ZoneID) {
	c.lsr.OnZoneChange(newZoneID)
}

func (c *UnicastController) Close() {
	c.lsr.Close()
}
