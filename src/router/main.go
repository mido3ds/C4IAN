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

func main() {
	defer log.Println("finished cleaning up, closing")

	ifaceName, pass, locSocket, _, err := parseArgs()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetOutput(os.Stdout)
	log.SetPrefix("[" + ifaceName + "] ")

	router, err := NewRouter(ifaceName, pass, locSocket)
	if err != nil {
		log.Fatal(err)
	}

	AddIPTablesRules()
	defer RemoveIPTablesRules()

	err = router.Start()
	if err != nil {
		panic(err)
	}

	// wait for SIGINT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Println("received SIGINT, started cleaning up")
}

func parseArgs() (string, string, string, string, error) {
	parser := argparse.NewParser("router", "Sets forwarding table in linux to route packets in adhoc-network.")
	ifaceName := parser.String("i", "iface", &argparse.Options{Required: true, Help: "Interface name."})
	passphrase := parser.String("p", "pass", &argparse.Options{Required: true, Help: "Passphrase for MSec (en/de)cryption."})
	locSocket := parser.String("l", "location-socket", &argparse.Options{Required: true, Help: "Path to unix domain socket to listen for location stream."})
	mgrpFile := parser.String("mg", "mgroups-file", &argparse.Options{Required: false, Help: "Path to mutlicast group member table file."})

	err := parser.Parse(os.Args)
	if err != nil {
		return "", "", "", "", errors.New(parser.Usage(err))
	}

	return *ifaceName, *passphrase, *locSocket, *mgrpFile, nil
}

func ReadJsonFile(path string) string {
	if path != "" {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		return string(content)
	}
	return "{}"
}
