package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
)

func main() {
	defer log.Println("finished cleaning up, closing")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetOutput(os.Stdout)

	args, err := parseArgs()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	// TODO: open store db
	// TODO: wrap writing to db
	// TODO: open port
	// TODO: open ctrl socket
	fmt.Println(args)
}

// Args store command line arguments
type Args struct {
	// null if no store path
	StorePath         string
	Port              int
	ControlSocketPath string
}

func parseArgs() (*Args, error) {
	parser := argparse.NewParser("unit-daemon", "Unit client daemon")

	storePath := parser.String("s", "store", &argparse.Options{Help: "Path to archive video/positions/heartbeats. If not provided, won't store them.", Default: nil})
	port := parser.Int("p", "port", &argparse.Options{Help: "Main port the client will bind to, to receive connections from other clients.", Default: 4070})
	ctrlSocketPath := parser.String("", "control-socket-path", &argparse.Options{Help: "Path to unix socket file to communicate over with controller.", Default: "/tmp/unit.ctrl.sock"})

	err := parser.Parse(os.Args)
	if err != nil {
		return nil, errors.New(parser.Usage(err))
	}

	return &Args{
		StorePath:         *storePath,
		Port:              *port,
		ControlSocketPath: *ctrlSocketPath,
	}, nil
}
