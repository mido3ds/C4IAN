package main

import (
	"fmt"
	"log"
	"net"
)

type Router struct {
	faceName string
	iface    *net.Interface
	ip4, ip6 net.IP
	msec     *MSecLayer
}

func NewRouter(ifaceName, passphrase string) (*Router, error) {
	// get interface
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, fmt.Errorf("couldn't get interface %s, error: %s", ifaceName, err)
	}

	// get initial ip addresses
	ip4, ip6, err := GetMyIPs(iface)
	if err != nil {
		return nil, fmt.Errorf("failed to get iface ips, err: %s", err)
	}
	log.Println("iface ipv4: ", ip4)
	log.Println("iface ipv6: ", ip6)

	return &Router{
		faceName: ifaceName,
		iface:    iface,
		msec:     NewMSecLayer(passphrase),
		ip4:      ip4,
		ip6:      ip6,
	}, nil
}

func (r *Router) Start() error {
	// initial modules
	forwarder, err := NewForwarder(r)
	if err != nil {
		return fmt.Errorf("failed to initialize forwarder, err: %s", err)
	}
	// TODO: initialize controller

	// start modules
	go forwarder.ForwardFromIPLayer()
	go forwarder.ForwardFromMACLayer()
	// TODO: start controller

	return nil
}
