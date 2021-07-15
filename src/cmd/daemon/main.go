package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/akamensky/argparse"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mido3ds/C4IAN/src/models"
)

const storePathSuffix = ".db"

func main() {
	defer log.Println("finished cleaning up, closing")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetOutput(os.Stdout)

	args, err := parseArgs()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	// TODO: read config, add units and groups to database
	dbManager := NewDatabaseManager(args.StorePath)
	api := NewAPI()
	netManager := NewNetworkManager(
		// onReceiveMessage
		func(msg models.Message) {
			api.SendEvent(&msg)
			dbManager.AddReceivedMessage(&msg)
		},
		// onReceiveAudio
		func(audio models.Audio) {
			api.SendEvent(&audio)
			dbManager.AddReceivedAudio(&audio)
		},
		// onReceiveVideoFragment
		func(frag models.VideoFragment) {
			api.SendEvent(&frag)
			// TODO: Handle received video fragment (create/append to file)
		},
		// onReceiveSensorsData
		func(data models.SensorData) {
			api.SendEvent(&data)
			dbManager.AddReceivedSensorsData(&data)
		},
	)
	go api.Start(args.UIPort, dbManager, netManager)
	netManager.Listen(args.Port)
	waitSIGINT()
}

// Args store command line arguments
type Args struct {
	StorePath string
	Port      int
	UIPort    int
}

func parseArgs() (*Args, error) {
	parser := argparse.NewParser("cmd-daemon", "Command-Center client daemon")

	storePath := parser.String("s", "store", &argparse.Options{Help: "Path to archive data.",
		Default: time.Now().Format(time.RFC3339) + storePathSuffix})

	port := parser.Int("p", "port", &argparse.Options{Default: 4170, Help: "Main port the client will bind to, to receive connections from other clients."})
	uiPort := parser.Int("", "ui-port", &argparse.Options{Default: 3170, Help: "UI port the client will bind to, to connect with its UI."})

	err := parser.Parse(os.Args)
	if err != nil {
		return nil, errors.New(parser.Usage(err))
	}

	// Enforce .db extension to storePath
	if !strings.HasSuffix(*storePath, storePathSuffix) {
		*storePath = *storePath + storePathSuffix
	}

	return &Args{
		StorePath: *storePath,
		Port:      *port,
		UIPort:    *uiPort,
	}, nil
}

func waitSIGINT() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
