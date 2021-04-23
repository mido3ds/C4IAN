package main

import (
	"log"
	"os/exec"
)

// TODO: support parallelism and fan-out

func AddIPTablesRules() {
	// add rule
	exec.Command("iptables", "-t", "filter", "-D", "OUTPUT", "-j", "NFQUEUE", "-w").Run()
	cmd := exec.Command("iptables", "-t", "filter", "-A", "OUTPUT", "-j", "NFQUEUE", "-w", "--queue-num", "0")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic("couldn't add iptables rule, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("added NFQUEUE rule to OUTPUT chain in iptables")

	// register gateway
	exec.Command("route", "del", "default", "gw", "localhost").Run()
	cmd = exec.Command("route", "add", "default", "gw", "localhost")
	stdoutStderr, err = cmd.CombinedOutput()
	if err != nil {
		log.Panic("couldn't add default gateway, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("added default gateway")
}

func RemoveIPTablesRules() {
	// remove rule
	cmd := exec.Command("iptables", "-t", "filter", "-D", "OUTPUT", "-j", "NFQUEUE", "-w")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic("couldn't remove iptables rule, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("remove NFQUEUE rule to OUTPUT chain in iptables")

	// remove gateway
	cmd = exec.Command("route", "del", "default", "gw", "localhost")
	stdoutStderr, err = cmd.CombinedOutput()
	if err != nil {
		log.Panic("couldn't remove default gateway, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("remove default gateway")
}
