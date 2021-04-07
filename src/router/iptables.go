package main

import (
	"log"
	"os/exec"
)

func AddIptablesRules() {
	// ipv4
	exec.Command("iptables", "-t", "filter", "-D", "OUTPUT", "-j", "NFQUEUE").Run()
	cmd := exec.Command("iptables", "-t", "filter", "-A", "OUTPUT", "-j", "NFQUEUE", "--queue-num", "0")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("couldn't add iptables rule, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("added NFQUEUE rule to OUTPUT chain in iptables")

	// ipv6
	exec.Command("ip6tables", "-t", "filter", "-D", "OUTPUT", "-j", "NFQUEUE").Run()
	cmd = exec.Command("ip6tables", "-t", "filter", "-A", "OUTPUT", "-j", "NFQUEUE", "--queue-num", "0")
	stdoutStderr, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatal("couldn't add ip6tables rule, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("added NFQUEUE rule to OUTPUT chain in ip6tables")
}

func RemoveIptablesRules() {
	// ipv4
	cmd := exec.Command("iptables", "-t", "filter", "-D", "OUTPUT", "-j", "NFQUEUE")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("couldn't remove iptables rule, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("remove NFQUEUE rule to OUTPUT chain in iptables")

	// ipv6
	cmd = exec.Command("ip6tables", "-t", "filter", "-D", "OUTPUT", "-j", "NFQUEUE")
	stdoutStderr, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatal("couldn't remove ip6tables rule, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("remove NFQUEUE rule to OUTPUT chain in ip6tables")
}
