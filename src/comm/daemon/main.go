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
	// TODO: define interface for ui
	fmt.Println(args)
}

// Args store command line arguments
type Args struct {
	StorePath string
	KeysPath  string
	Adr       string
	UIAdr     string
}

func parseArgs() Args {
	parser := argparse.NewParser("unit-daemon", "Unit client daemon")

	storePath := parser.String("s", "store", &argparse.Options{Help: "Path to archive data.", Required: true})
	keysPath := parser.String("k", "keys", &argparse.Options{Help: "Path to keys file. To get commands IPs and their pub-keys", Required: true})

	adr := parser.String("a", "adr", &argparse.Options{Help: "Main address+port the client will bind to receive connections from other clients.", Default: "0.0.0.0:4170"})
	uiAdr := parser.String("", "ui-adr", &argparse.Options{Default: "0.0.0.0:3170", Help: "UI address+port the client will bind to connect with its UI."})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	return Args{
		StorePath: *storePath,
		KeysPath:  *keysPath,
		Adr:       *adr,
		UIAdr:     *uiAdr,
	}
}
