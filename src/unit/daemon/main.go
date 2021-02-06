package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

func main() {
	args := parseArgs()

	// TODO: open store db
	// TODO: wrap writing to db
	// TODO: read keys file
	// TODO: open port
	// TODO: define interface for virt
	fmt.Println(args)
}

// Args store command line arguments
type Args struct {
	// null if no store path
	StorePath string
	KeysPath  string
	Port      int
	IsVirt    bool
	UIPort    int
}

func parseArgs() Args {
	parser := argparse.NewParser("unit-daemon", "Unit client daemon")

	storePath := parser.String("s", "store", &argparse.Options{Help: "Path to archive video/positions/heartbeats. If not provided, won't store them.", Default: nil})
	keysPath := parser.String("k", "keys", &argparse.Options{Help: "Path to keys file. To get commands IPs and their pub-keys.", Required: true})
	port := parser.Int("p", "port", &argparse.Options{Help: "Main port the client will bind to, to receive connections from other clients.", Default: 4070})

	virt := parser.NewCommand("virt", "Run in virtual mode")
	uiPort := virt.Int("", "ui-port", &argparse.Options{Default: 3070, Help: "UI port the client will bind to, to connect with its UI."})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	return Args{
		StorePath: *storePath,
		KeysPath:  *keysPath,
		Port:      *port,
		IsVirt:    virt.Happened(),
		UIPort:    *uiPort,
	}
}
