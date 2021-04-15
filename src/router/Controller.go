package main

import (
	"log"
)

type Controller struct {
	macConn      *MACLayerConn
	inputChannel chan []byte
}

func NewController(router *Router) (*Controller, error) {
	macConn, err := NewMACLayerConn(router.iface)
	if err != nil {
		return nil, err
	}

	c := make(chan []byte)

	log.Println("initalized controller")

	return &Controller{
		macConn:      macConn,
		inputChannel: c,
	}, nil
}

func (controller *Controller) ListenForControlPackets() {
	log.Println("Controller started listening for control packets from the forwarder")

	for {
		<-controller.inputChannel
		log.Println("Controller received a control packet")

		// TODO: Check for control packet type and handle accordingly
	}
}
