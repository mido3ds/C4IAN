package main

import (
	"fmt"
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/database_logger"
	. "github.com/mido3ds/C4IAN/src/router/forward"
	. "github.com/mido3ds/C4IAN/src/router/ip"
	"github.com/mido3ds/C4IAN/src/router/kernel"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/odmrp"
	. "github.com/mido3ds/C4IAN/src/router/sarp"
	. "github.com/mido3ds/C4IAN/src/router/tables"
	. "github.com/mido3ds/C4IAN/src/router/zhls"
	. "github.com/mido3ds/C4IAN/src/router/zhls/dzd"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

type Router struct {
	iface *net.Interface
	ip    net.IP
	zlen  byte

	msec      *MSecLayer
	forwarder *Forwarder
	timers    *TimersQueue

	// controllers
	zidAgent *ZoneIDAgent
	unicCont *UnicastController
	multCont *MulticastController
	sarpCont *SARPController
	dzdCont  *DZDController
}

func NewRouter(ifaceName, passphrase, locSocket string, zlen byte, mgrpFilePath string) (*Router, error) {
	// tell linux im a router
	kernel.AddIPTablesRule()
	if err := kernel.RegisterGateway(); err != nil {
		kernel.DeleteIPTablesRule()
		return nil, err
	}

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

	// tell linux to reduce MSS to accomodate our added headers
	kernel.SetMaxMSS(ifaceName, ip, 1400)

	// Initialize database logger
	DatabaseLogger.Initialize(ip)

	msec := NewMSecLayer(passphrase)

	topology := NewTopology(ip)
	timers := NewTimersQueue()

	log.Println("ZLen =", zlen, ", Zone Max Area =", ZLenToAreaKMs(zlen), "km^2")
	zidAgent, err := NewZoneIDAgent(locSocket, zlen)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize location agent, err: %s", err)
	}

	sarpCont, err := NewSARPController(ip, iface, msec)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate sARP, err: %s", err)
	}

	unicCont, err := NewUnicastController(iface, ip, sarpCont.NeighborsTable, sarpCont.NeighborhoodUpdateSignal, msec, topology)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize unicast controller, err: %s", err)
	}

	multCont, err := NewMulticastController(iface, ip, iface.HardwareAddr, msec, mgrpFilePath, timers)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize unicast controller, err: %s", err)
	}

	dzdCont, err := NewDZDController(ip, iface, topology)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize dzd controller, err: %s", err)
	}

	forwarder, err := NewForwarder(iface, ip, msec, sarpCont.NeighborsTable, dzdCont, multCont.GetMissingEntries,
		multCont.IsDest, unicCont.UpdateUnicastForwardingTable, timers)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize forwarder, err: %s", err)
	}

	zidAgent.AddZoneChangeCallback(unicCont.OnZoneChange)
	dzdCont.SetGetNextHopCallback(forwarder.GetUnicastNextHop)

	return &Router{
		iface:     iface,
		msec:      msec,
		ip:        ip,
		zlen:      zlen,
		forwarder: forwarder,
		timers:    timers,
		zidAgent:  zidAgent,
		unicCont:  unicCont,
		multCont:  multCont,
		sarpCont:  sarpCont,
		dzdCont:   dzdCont,
	}, nil
}

func (r *Router) Start() {
	go r.timers.Start()

	// zid agent
	go r.zidAgent.Start()

	// start location logging
	go DatabaseLogger.StartLocationLogging()

	// start controllers
	go r.sarpCont.Start()
	go r.unicCont.Start()
	go r.multCont.Start(r.forwarder.MultiForwTable)
	go r.forwarder.Start()
	go r.dzdCont.Start()

	log.Println(r.ip, ":", r.iface.HardwareAddr)
}

func (r *Router) Close() {
	r.forwarder.Close()
	r.multCont.Close()
	r.unicCont.Close()
	r.sarpCont.Close()
	r.dzdCont.Close()

	r.zidAgent.Close()

	kernel.UnregisterGateway()
	kernel.DeleteIPTablesRule()
}
