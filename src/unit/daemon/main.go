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
	fmt.Println("-----------------------")

	// TODO: open store db
	// TODO: wrap writing to db
	// TODO: get key
	// TODO: open port
	// TODO: define interface for virt
	fmt.Println(args)
}

// Args store command line arguments
type Args struct {
	// null if no store path
	StorePath  string
	KeysPath   string
	PassPhrase string
	Port       int
	IsVirt     bool
	UIPort     int
}

func parseArgs() (*Args, error) {
	parser := argparse.NewParser("unit-daemon", "Unit client daemon")

	storePath := parser.String("s", "store", &argparse.Options{Help: "Path to archive video/positions/heartbeats. If not provided, won't store them.", Default: nil})
	passphrase := parser.String("", "pass", &argparse.Options{Help: "Passphrase.", Required: true})
	port := parser.Int("p", "port", &argparse.Options{Help: "Main port the client will bind to, to receive connections from other clients.", Default: 4070})

	virt := parser.NewCommand("virt", "Run in virtual mode")
	uiPort := virt.Int("", "ui-port", &argparse.Options{Default: 3070, Help: "UI port the client will bind to, to connect with its UI."})

	err := parser.Parse(os.Args)
	if err != nil {
		return nil, errors.New(parser.Usage(err))
	}

	return &Args{
		StorePath:  *storePath,
		PassPhrase: *passphrase,
		Port:       *port,
		IsVirt:     virt.Happened(),
		UIPort:     *uiPort,
	}, nil
}
