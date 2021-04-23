package kernel

import (
	"fmt"
	"log"
	"os/exec"
)

func RegisterGateway() error {
	exec.Command("route", "del", "default", "gw", "localhost").Run()
	cmd := exec.Command("route", "add", "default", "gw", "localhost")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("couldn't add default gateway, err: %#v, stderr: %#v", err, string(stdoutStderr))
	}
	log.Println("added default gateway")
	return nil
}

func UnregisterGateway() {
	cmd := exec.Command("route", "del", "default", "gw", "localhost")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic("couldn't remove default gateway, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("remove default gateway")
}
