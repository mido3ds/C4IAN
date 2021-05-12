package zhls

import (
	"fmt"
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/tables"
)

type UnicastController struct {
	ip                           net.IP
	lsr                          *LSRController
	neighborhoodUpdateSignal     chan bool
	neighborsTable               *NeighborsTable
	UpdateUnicastForwardingTable func(ft *UniForwardTable)
}

func NewUnicastController(iface *net.Interface, ip net.IP, neighborsTable *NeighborsTable, neighborhoodUpdateSignal chan bool, msec *MSecLayer, topology *Topology) (*UnicastController, error) {
	lsr, err := newLSR(iface, msec, ip, neighborsTable, topology)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate LSR Controller, err: %#v", err)
	}

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
		<-c.neighborhoodUpdateSignal
		c.lsr.onNeighborhoodUpdate()
	}
}

func (c *UnicastController) Close() {
	c.lsr.Close()
}
