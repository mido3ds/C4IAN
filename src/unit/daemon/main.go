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

	KeysPath string
	IsVirt   bool
	Port     int
}

func parseArgs() Args {
	parser := argparse.NewParser("unit-daemon", "Unit client daemon")
	storePath := parser.String("s", "store", &argparse.Options{Help: "Path to archive video/positions/heartbeats. If not provided, won't store them.", Default: nil})
	keysPath := parser.String("k", "keys", &argparse.Options{Help: "Path to keys file. To get commands IPs and their pub-keys", Required: true})

	virt := parser.NewCommand("virt", "Run in virtual mode")
	port := virt.Int("p", "port", &argparse.Options{Default: 3070, Help: "Port to bind to"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	return Args{
		StorePath: *storePath,
		KeysPath:  *keysPath,
		IsVirt:    virt.Happened(),
		Port:      *port,
	}
}
