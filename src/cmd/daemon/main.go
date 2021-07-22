package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	args, err := parseArgs()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Println(args)
	defer log.Println("finished cleaning up, closing")

	var iface string
	if len(args.Iface) > 0 {
		iface = args.Iface
	} else {
		iface = getDefaultInterface()
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetOutput(os.Stdout)
	log.SetPrefix("[" + iface + "] ")

	units, groupMembers := parseConfig(args.UnitsPath, args.GroupsPath)
	dbManager := NewDatabaseManager(args.StorePath)
	dbManager.Initialize(units, groupMembers)
	api := NewAPI()
	videoFilesManager := NewVideoFilesManager(args.VideosPath)
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
			videoFilesManager.AddFragment(&frag)
			exists := dbManager.AddVideoIfNew(&frag)
			if !exists {
				api.SendEvent(&models.Video{Src: frag.Src, ID: frag.ID, Time: frag.Time})
			}
		},
		// onReceiveSensorsData
		func(data models.SensorData) {
			api.SendEvent(&data)
			dbManager.AddReceivedSensorsData(&data)
			dbManager.UpdateLastActivity(data.Src, data.Time)
		},
	)
	go api.Start(args.UISocket, args.UnitsPort, args.VideosPath, dbManager, netManager)
	go netManager.SendGroupsHello(groupMembers, args.UnitsPort)
	netManager.Listen(args.Port)
	log.Println("finished initalizing all")
	waitSIGINT()
}

// Args store command line arguments
type Args struct {
	StorePath  string
	VideosPath string
	UnitsPath  string
	GroupsPath string
	Iface      string
	UISocket   string
	Port       int
	UnitsPort  int
}

func parseArgs() (*Args, error) {
	parser := argparse.NewParser("cmd-daemon", "Command-Center client daemon")

	storePath := parser.String("s", "store", &argparse.Options{Help: "Path to archive data.",
		Default: time.Now().Format(time.RFC3339) + storePathSuffix})

	videosPath := parser.String("v", "videos-path", &argparse.Options{Default: "videos", Help: "Path to store videos received from units."})
	unitsPath := parser.String("u", "units-path", &argparse.Options{Default: "../../units.json", Help: "Path to units.json."})
	groupsPath := parser.String("g", "groups-path", &argparse.Options{Default: "../../groups.json", Help: "Path to groups.json."})

	port := parser.Int("p", "port", &argparse.Options{Default: 4170, Help: "Main port the client will bind to, to receive connections from other clients."})
	uiSocket := parser.String("", "ui-socket", &argparse.Options{Default: "/tmp/cmd.sock", Help: "Unix socket file that the client will listen on, to connect with its UI."})
	unitsPort := parser.Int("", "units-port", &argparse.Options{Default: 4070, Help: "Main port used in units."})

	iface := parser.String("", "iface", &argparse.Options{Help: "Name of this interface. Default is to list the ifaces with /proc/net/route.", Default: ""})

	err := parser.Parse(os.Args)
	if err != nil {
		return nil, errors.New(parser.Usage(err))
	}

	// Enforce .db extension to storePath
	if !strings.HasSuffix(*storePath, storePathSuffix) {
		*storePath = *storePath + storePathSuffix
	}

	return &Args{
		StorePath:  *storePath,
		VideosPath: *videosPath,
		UnitsPath:  *unitsPath,
		GroupsPath: *groupsPath,
		Port:       *port,
		UISocket:   *uiSocket,
		UnitsPort:  *unitsPort,
		Iface:      *iface,
	}, nil
}

func parseConfig(unitsPath string, groupsPath string) (units []models.Unit, groups map[string][]string) {
	data, err := ioutil.ReadFile(unitsPath)
	if err != nil {
		log.Println("failed to read units file, err:", err)
	}
	var result map[string][]models.Unit
	json.Unmarshal(data, &result)
	units = result["units"]

	data, err = ioutil.ReadFile(groupsPath)
	if err != nil {
		log.Println("failed to read groups file, err:", err)
	}
	json.Unmarshal(data, &groups)
	return
}

func waitSIGINT() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func getDefaultInterface() string {
	file, err := os.Open("/proc/net/route")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		const line = 1 // line containing the gateway addr. (first line: 0)
		// jump to line containing the agteway address
		for i := 0; i < line; i++ {
			scanner.Scan()
		}

		// get field containing gateway address
		tokens := strings.Split(scanner.Text(), "\t")
		iface := tokens[0]
		return iface
	}

	log.Panic("no default interface found")
	return "unreachable"
}
