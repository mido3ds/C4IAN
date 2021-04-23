package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/akamensky/argparse"
)

const DefaultZLen = 12

func main() {
	args, err := parseArgs()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer log.Println("finished cleaning up, closing")

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetOutput(os.Stdout)
	log.SetPrefix("[" + args.ifaceName + "] ")

	router, err := NewRouter(args.ifaceName, args.passphrase, args.locSocket, args.zlen)
	if err != nil {
		log.Panic(err)
	}

	AddIPTablesRules()
	defer RemoveIPTablesRules()

	mgrpContent := ReadJsonFile(args.mgrpFilePath)
	err = router.Start(mgrpContent)
	if err != nil {
		log.Panic(err)
	}

	// wait for SIGINT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Println("received SIGINT, started cleaning up")
}

type Args struct {
	ifaceName    string
	passphrase   string
	locSocket    string
	mgrpFilePath string
	zlen         byte
}

func parseArgs() (*Args, error) {
	parser := argparse.NewParser("router", "Sets forwarding table in linux to route packets in adhoc-network.")
	ifaceName := parser.String("i", "iface", &argparse.Options{Required: true, Help: "Interface name."})
	passphrase := parser.String("p", "pass", &argparse.Options{Required: true, Help: "Passphrase for MSec (en/de)cryption."})
	locSocket := parser.String("l", "location-socket", &argparse.Options{Required: true, Help: "Path to unix domain socket to listen for location stream."})
	mgrpFile := parser.String("g", "mgroups-file", &argparse.Options{Required: false, Help: "Path to mutlicast group member table file."})
	zlen := parser.Int("", "zlen", &argparse.Options{Required: false, Default: DefaultZLen, Help: "ZLen value to determine zone area, must be between 0 and 16 inclusive."})

	err := parser.Parse(os.Args)
	if err != nil {
		return nil, errors.New(parser.Usage(err))
	}

	if *zlen < 0 || *zlen > 16 {
		return nil, errors.New(parser.Usage("ZLen must be between 0 and 16 inclusive"))
	}

	if *mgrpFile != "" && !fileExists(*mgrpFile) {
		return nil, errors.New(parser.Usage("mgroups-file doesn't exist"))
	}

	return &Args{
		ifaceName:    *ifaceName,
		passphrase:   *passphrase,
		locSocket:    *locSocket,
		mgrpFilePath: *mgrpFile,
		zlen:         byte(*zlen),
	}, nil
}

func ReadJsonFile(path string) string {
	if path != "" {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Panic(err)
		}
		return string(content)
	}
	return "{}"
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
