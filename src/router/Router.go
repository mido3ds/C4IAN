package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Router struct {
	iface    *net.Interface
	ip       net.IP
	msec     *MSecLayer
	locAgent *LocationAgent
}

func NewRouter(ifaceName, passphrase, locSocket string) (*Router, error) {
	// get interface
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, fmt.Errorf("couldn't get interface %s, error: %s", ifaceName, err)
	}

	// get initial ip addresses
	ip, _, err := GetMyIPs(iface)
	if err != nil {
		return nil, fmt.Errorf("failed to get iface ips, err: %s", err)
	}
	log.Println("iface ipv4: ", ip)

	// location agent
	loc, err := NewLocationAgent(locSocket)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize location agent, err: %s", err)
	}

	return &Router{
		iface:    iface,
		msec:     NewMSecLayer(passphrase),
		ip:       ip,
		locAgent: loc,
	}, nil
}

func (r *Router) Start() error {
	// initial modules
	controller, err := NewController(r)
	if err != nil {
		return fmt.Errorf("failed to initialize controller, err: %s", err)
	}

	forwarder, err := NewForwarder(r)
	if err != nil {
		return fmt.Errorf("failed to initialize forwarder, err: %s", err)
	}

	// start modules
	go r.locAgent.Start()
	go controller.ListenForControlPackets()
	go controller.sARP()
	go forwarder.ForwardFromIPLayer()
	go forwarder.ForwardFromMACLayer(controller.inputChannel)

	time.Sleep(5 * time.Second)
	controller.floodDummy()
	
	return nil
}
