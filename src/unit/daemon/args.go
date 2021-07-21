package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/akamensky/argparse"
)

// Args store command line arguments
type Args struct {
	// null if no store path
	storePath           string
	port                int
	halSocketPath       string
	cmdAddress          string
	timeout             time.Duration
	retryOrCloseTimeout time.Duration
	UISocket            string
	iface               string
}

func (a *Args) String() string {
	return fmt.Sprintf("&Args{StorePath: %v, Port: %v, HALSocketPath: %v, CMDAddress: %v, Timeout: %v, RetryOrCloseTimeout: %v}", a.storePath, a.port, a.halSocketPath, a.cmdAddress, a.timeout, a.retryOrCloseTimeout)
}

func parseArgs() (*Args, error) {
	parser := argparse.NewParser("unit-daemon", "Unit client daemon")

	storePath := parser.String("s", "store", &argparse.Options{Help: "Path to archive video/positions/heartbeats. If not provided, won't store them.", Default: nil})
	port := parser.Int("p", "port", &argparse.Options{Help: "Main port the client will bind to, to receive connections from other clients.", Default: 4070})
	uiSocket := parser.String("", "ui-socket", &argparse.Options{Default: "/tmp/unit.sock", Help: "Unix socket file that the client will listen on, to connect with its UI."})
	timeout := parser.Int("", "timeout-millis", &argparse.Options{Help: "Timeout of writing to CMD in millis, if timeouted will reopen socket and send.", Default: 1000})
	retryTimeout := parser.Int("", "retry-millis", &argparse.Options{Help: "Timeout of resending same packet in millis, if failed again will reopen socket.", Default: 3000})
	cmdAddress := parser.String("", "cmd-addr", &argparse.Options{Help: "Address+Port of cmd center to communicate with.", Default: "127.0.0.1:4170"})
	ctrlSocketPath := parser.String("", "hal-socket-path", &argparse.Options{Help: "Path to unix socket file to communicate over with HAL.", Default: "/tmp/unit.hal.sock"})

	iface := parser.String("", "iface", &argparse.Options{Help: "Name of this interface. Default is to list the ifaces with /proc/net/route.", Default: ""})

	err := parser.Parse(os.Args)
	if err != nil {
		return nil, errors.New(parser.Usage(err))
	}

	return &Args{
		storePath:           *storePath,
		port:                *port,
		halSocketPath:       *ctrlSocketPath,
		cmdAddress:          *cmdAddress,
		timeout:             time.Duration(*timeout) * time.Millisecond,
		retryOrCloseTimeout: time.Duration(*retryTimeout) * time.Millisecond,
		UISocket:            *uiSocket,
		iface:               *iface,
	}, nil
}
