package main

import (
	"log"
)

type Controller struct {
	router       *Router
	macConn      *MACLayerConn
	inputChannel chan []byte
	flooder		 *Flooder
}

func (c *Controller) floodDummy() {
	dummy:= []byte{0xAA, 0xBB}
	c.flooder.flood(dummy)
}


func NewController(router *Router) (*Controller, error) {
	macConn, err := NewMACLayerConn(router.iface)
	if err != nil {
		return nil, err
	}

	c := make(chan []byte)
	flooder := NewFlooder(router)

	log.Println("initalized controller")

	return &Controller{
		router:       router,
		macConn:      macConn,
		inputChannel: c,
		flooder:	  flooder	 
	}, nil
}

func (controller *Controller) ListenForControlPackets() {
	log.Println("Controller started listening for control packets from the forwarder")

	for {
		packet <-controller.inputChannel
		log.Println("Controller received a control packet")

		// TODO: Check for control packet type and handle accordingly
	}
}
