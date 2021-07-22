package kernel

import (
	"log"
	"net"
	"os/exec"
	"strconv"
)

func SetMaxMSS(ifaceName string, ip net.IP, mss int) {
	cmd := exec.Command(
		"ip", "route", "change", "10.0.0.0/8", "dev", ifaceName,
		"proto", "kernel", "scope", "link", "src", ip.String(), "advmss", strconv.Itoa(mss),
	)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic("Couldn't change MSS, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("Changed MSS to: ", mss)
}
