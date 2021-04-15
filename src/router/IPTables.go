package main

import (
	"log"
	"os/exec"
)

// TODO: support parallelism and fan-out

func AddIPTablesRules() {
	// ipv4
	exec.Command("iptables", "-t", "filter", "-D", "OUTPUT", "-j", "NFQUEUE").Run()
	cmd := exec.Command("iptables", "-t", "filter", "-A", "OUTPUT", "-j", "NFQUEUE", "--queue-num", "0")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("couldn't add iptables rule, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("added NFQUEUE rule to OUTPUT chain in iptables")
}

func RemoveIPTablesRules() {
	// ipv4
	cmd := exec.Command("iptables", "-t", "filter", "-D", "OUTPUT", "-j", "NFQUEUE")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("couldn't remove iptables rule, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("remove NFQUEUE rule to OUTPUT chain in iptables")
}
