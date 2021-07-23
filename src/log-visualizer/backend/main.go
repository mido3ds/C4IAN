package main

import (
	"errors"
	"log"
	"os"

	"github.com/akamensky/argparse"
)

const port = 5000

func main() {
	args := parseArgs()
	log.Println(args.Path)
	dbManager := NewDatabaseManager(args.Path)
	api := NewAPI(dbManager)
	api.Start(port)
}

type Args struct {
	Path string
}

func parseArgs() *Args {
	parser := argparse.NewParser("cmd-daemon", "Command-Center client daemon")

	path := parser.String("p", "path", &argparse.Options{Default: "/var/log/caian/log.db", Help: "Log path."})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Panic(errors.New(parser.Usage(err)))
	}

	return &Args{
		Path: *path,
	}
}
