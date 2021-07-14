package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/akamensky/argparse"
	_ "github.com/mattn/go-sqlite3"
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

	// TODO: read config
	dbManager := NewDatabaseManager(args.StorePath)
	api := NewAPI(dbManager)
	go api.Start(args.UIPort)
	// TODO: wrap writing to db
	// TODO: open port
	// TODO: define interface for ui
	fmt.Println(args)
	// id := 0
	// for {
	// 	api.SendEvent(&models.Message{Code: id})
	// 	api.SendEvent(&models.Audio{Body: []byte(strconv.Itoa(id))})
	// 	api.SendEvent(&models.VideoFrame{Body: []byte(strconv.Itoa(id + 100))})
	// 	api.SendEvent(&models.SensorData{Heartbeat: id * 2, Loc_x: 5, Loc_y: 73})
	// 	id++
	// 	time.Sleep(time.Second)
	// }
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
