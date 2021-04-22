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
	ip = ip.To4()
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
	neighborhoodUpdateSignal := make(chan bool)
	sARP, err := NewSARPController(r, neighborhoodUpdateSignal)
	if err != nil {
		return fmt.Errorf("failed to initiate sARP, err: %s", err)
	}

	unicCont, err := NewUnicastController(r, sARP.neighborsTable, neighborhoodUpdateSignal)
	if err != nil {
		return fmt.Errorf("failed to initialize unicast controller, err: %s", err)
	}

	forwarder, err := NewForwarder(r, sARP.neighborsTable)
	if err != nil {
		return fmt.Errorf("failed to initialize forwarder, err: %s", err)
	}

	// start modules
	go sARP.Start()
	go r.locAgent.Start()
	go unicCont.Start()
	go forwarder.Start(unicCont.inputChannel)

	time.AfterFunc(6*time.Second, func() {
		log.Println(controller.sARP.neighborsTable)
	})

	time.AfterFunc(10*time.Second, func() {
		unicCont.lsr.UpdateForwardingTable(r.ip, forwarder.uniForwTable, sARP.neighborsTable)
		log.Println(forwarder.uniForwTable)
	})

	return nil
}
