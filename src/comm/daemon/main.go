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
	Port      int
}

func parseArgs() Args {
	parser := argparse.NewParser("unit-daemon", "Unit client daemon")

	storePath := parser.String("s", "store", &argparse.Options{Help: "Path to archive data.", Required: true})
	keysPath := parser.String("k", "keys", &argparse.Options{Help: "Path to keys file. To get commands IPs and their pub-keys", Required: true})
	port := parser.Int("p", "port", &argparse.Options{Default: 3170, Help: "Port to bind to"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	return Args{
		StorePath: *storePath,
		KeysPath:  *keysPath,
		Port:      *port,
	}
}
