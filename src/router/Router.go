package main

import (
	"fmt"
	"log"
	"net"
)

type Router struct {
	iface    *net.Interface
	ip       net.IP
	zlen     byte
	msec     *MSecLayer
	locAgent *LocationAgent
}

func NewRouter(ifaceName, passphrase, locSocket string, zlen byte) (*Router, error) {
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
	loc, err := NewLocationAgent(locSocket, zlen)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize location agent, err: %s", err)
	}

	return &Router{
		iface:    iface,
		msec:     NewMSecLayer(passphrase),
		ip:       ip,
		zlen:     zlen,
		locAgent: loc,
	}, nil
}

func (r *Router) Start(mgrpContent string) error {
	// initialize modules
	sARP, err := NewSARPController(r)
	if err != nil {
		return fmt.Errorf("failed to initiate sARP, err: %s", err)
	}

	unicCont, err := NewUnicastController(r, sARP.neighborsTable, sARP.neighborhoodUpdateSignal)
	if err != nil {
		return fmt.Errorf("failed to initialize unicast controller, err: %s", err)
	}

	multicCont, err := NewMulticastController(r, mgrpContent)
	if err != nil {
		return fmt.Errorf("failed to initialize unicast controller, err: %s", err)
	}

	forwarder, err := NewForwarder(r, sARP.neighborsTable)
	if err != nil {
		return fmt.Errorf("failed to initialize forwarder, err: %s", err)
	}

	r.locAgent.AddListener(forwarder.OnZoneIDChanged)
	r.locAgent.AddListener(unicCont.OnZoneIDChanged)

	// start modules
	go sARP.Start()
	go r.locAgent.Start()
	go unicCont.Start(forwarder.uniForwTable)
	go multicCont.Start(forwarder.multiForwTable)
	go forwarder.Start(unicCont.inputChannel)

	return nil
}
