package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

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
	// TODO: define interface for ui
	fmt.Println(args)
}

// Args store command line arguments
type Args struct {
	StorePath string
	Port      int
	UIPort    int
}

func parseArgs() (*Args, error) {
	parser := argparse.NewParser("cmd-daemon", "Command-Center client daemon")

	storePath := parser.String("s", "store", &argparse.Options{Help: "Path to archive data.", Default: time.Now().Format(time.RFC3339) + ".db"})

	port := parser.Int("p", "port", &argparse.Options{Help: "Main port the client will bind to, to receive connections from other clients.", Default: 4170})
	uiPort := parser.Int("", "ui-port", &argparse.Options{Help: "UI port the client will bind to, to connect with its UI.", Default: 3170})

	err := parser.Parse(os.Args)
	if err != nil {
		return nil, errors.New(parser.Usage(err))
	}

	return &Args{
		StorePath: *storePath,
		Port:      *port,
		UIPort:    *uiPort,
	}, nil
}
