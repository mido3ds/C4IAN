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
	// TODO: get key
	// TODO: open port
	// TODO: define interface for ui
	fmt.Println(args)
}

// Args store command line arguments
type Args struct {
	StorePath  string
	Passphrase string
	Port       int
	UIPort     string
}

func parseArgs() (*Args, error) {
	parser := argparse.NewParser("unit-daemon", "Unit client daemon")

	storePath := parser.String("s", "store", &argparse.Options{Help: "Path to archive data.", Required: true})
	passphrase := parser.String("", "pass", &argparse.Options{Help: "Passphrase.", Required: true})

	port := parser.Int("p", "port", &argparse.Options{Help: "Main port the client will bind to, to receive connections from other clients.", Default: 4170})
	uiPort := parser.String("", "ui-port", &argparse.Options{Default: 3170, Help: "UI port the client will bind to, to connect with its UI."})

	err := parser.Parse(os.Args)
	if err != nil {
		return nil, errors.New(parser.Usage(err))
	}

	return &Args{
		StorePath:  *storePath,
		Passphrase: *passphrase,
		Port:       *port,
		UIPort:     *uiPort,
	}, nil
}
