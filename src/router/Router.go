package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"

	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/odmrp"
	. "github.com/mido3ds/C4IAN/src/router/zhls"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

type Router struct {
	iface *net.Interface
	ip    net.IP
	zlen  byte

	msec      *MSecLayer
	forwarder *Forwarder

	// controllers
	zidAgent *ZoneIDAgent
	unicCont *UnicastController
	multCont *MulticastController
	sarpCont *SARPController
}

func NewRouter(ifaceName, passphrase, locSocket string, zlen byte, mgrpFilePath string) (*Router, error) {
	// tell linux im a router
	addIPTablesRule()
	if err := registerGateway(); err != nil {
		deleteIPTablesRule()
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

	msec := NewMSecLayer(passphrase)

	zidAgent, err := NewZoneIDAgent(locSocket, zlen)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize location agent, err: %s", err)
	}

	sarpCont, err := NewSARPController(ip, iface, msec)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate sARP, err: %s", err)
	}

	unicCont, err := NewUnicastController(iface, ip, sarpCont.neighborsTable, sarpCont.neighborhoodUpdateSignal, msec, zlen)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize unicast controller, err: %s", err)
	}

	multCont, err := NewMulticastController(iface, ip, msec, mgrpFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize unicast controller, err: %s", err)
	}

	forwarder, err := NewForwarder(iface, ip, msec, zlen, sarpCont.neighborsTable, multCont.GetMissingEntries)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize forwarder, err: %s", err)
	}

	return &Router{
		iface:     iface,
		msec:      msec,
		ip:        ip,
		zlen:      zlen,
		forwarder: forwarder,
		zidAgent:  zidAgent,
		unicCont:  unicCont,
		multCont:  multCont,
		sarpCont:  sarpCont,
	}, nil
}

func (r *Router) Start() {
	// zid agent
	r.zidAgent.AddListener(r.forwarder.OnZoneIDChanged)
	r.zidAgent.AddListener(r.unicCont.OnZoneIDChanged)
	go r.zidAgent.Start()

	// start controllers
	go r.sarpCont.Start()
	go r.unicCont.Start(r.forwarder.uniForwTable)
	go r.multCont.Start(r.forwarder.multiForwTable)
	go r.forwarder.Start(r.unicCont.InputChannel)
}

func (r *Router) Close() {
	deleteIPTablesRule()
	unregisterGateway()
}

// TODO: support parallelism and fan-out

func addIPTablesRule() {
	exec.Command("iptables", "-t", "filter", "-D", "OUTPUT", "-j", "NFQUEUE", "-w").Run()
	cmd := exec.Command("iptables", "-t", "filter", "-A", "OUTPUT", "-j", "NFQUEUE", "-w", "--queue-num", "0")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic("couldn't add iptables rule, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("added NFQUEUE rule to OUTPUT chain in iptables")
}

func deleteIPTablesRule() {
	cmd := exec.Command("iptables", "-t", "filter", "-D", "OUTPUT", "-j", "NFQUEUE", "-w")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic("couldn't remove iptables rule, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("remove NFQUEUE rule to OUTPUT chain in iptables")
}

func registerGateway() error {
	exec.Command("route", "del", "default", "gw", "localhost").Run()
	cmd := exec.Command("route", "add", "default", "gw", "localhost")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("couldn't add default gateway, err: %#v, stderr: %#v", err, string(stdoutStderr))
	}
	log.Println("added default gateway")
	return nil
}

func unregisterGateway() {
	cmd := exec.Command("route", "del", "default", "gw", "localhost")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic("couldn't remove default gateway, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("remove default gateway")
}
