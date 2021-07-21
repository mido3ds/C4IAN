package kernel

import (
	"log"
	"os/exec"
)

// TODO (low priority): support parallelism and fan-out

func AddIPTablesRule(ifaceName string) {
	exec.Command("iptables", "-t", "filter", "-D", "OUTPUT", "-j", "NFQUEUE", "-w", "-o", ifaceName).Run()
	cmd := exec.Command("iptables", "-t", "filter", "-A", "OUTPUT", "-j", "NFQUEUE", "-w", "-o", ifaceName, "--queue-num", "0")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic("couldn't add iptables rule, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("added NFQUEUE rule to OUTPUT chain in iptables")
}

func DeleteIPTablesRule(ifaceName string) {
	cmd := exec.Command("iptables", "-t", "filter", "-D", "OUTPUT", "-j", "NFQUEUE", "-w", "-o", ifaceName)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic("couldn't remove iptables rule, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("remove NFQUEUE rule to OUTPUT chain in iptables")
}
