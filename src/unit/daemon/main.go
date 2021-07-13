package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/akamensky/argparse"
	"github.com/mido3ds/C4IAN/src/unit/daemon/halapi"
)

type Context struct {
	expectingVideoStream bool
	args                 Args
}

func main() {
	defer log.Println("finished cleaning up, closing")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetOutput(os.Stdout)

	args, err := parseArgs()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	// TODO: open store db
	// TODO: wrap writing to db
	// TODO: open port

	// TODO: remove
	fmt.Println(args)

	context := Context{expectingVideoStream: false, args: *args}

	// remove loc socket file
	err = os.RemoveAll(args.HALSocketPath)
	if err != nil {
		log.Panic("failed to remove socket:", err)
	}

	l, err := net.Listen("unix", args.HALSocketPath)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer l.Close()

	for {
		go simulateClient(args.HALSocketPath)
		conn, err := l.Accept()
		defer conn.Close()

		if err != nil {
			log.Println("accept error:", err)
		} else {
			log.Println("HAL connected")
			serve(&context, conn)
		}
	}
}

// TODO: remove
func simulateClient(HALSocketPath string) {
	conn, err := net.Dial("unix", HALSocketPath)
	if err != nil {
		log.Panic(err)
	}

	enc := gob.NewEncoder(conn)
	halapi.Location{Lon: 5, Lat: 2}.Send(enc)
	halapi.HeartBeat{BeatsPerMinut: 5}.Send(enc)
	halapi.VideoPart{Video: []byte{5, 3}}.Send(enc)
}

func serve(context *Context, conn net.Conn) {
	dec := gob.NewDecoder(conn)

	var video halapi.VideoPart
	var heartbeat halapi.HeartBeat
	var loc halapi.Location

	for {
		sentType, err := halapi.RecvFromHAL(dec, &video, &heartbeat, &loc)
		if err != nil {
			log.Println(err)
		} else {
			switch sentType {
			case halapi.VideoPartType:
				context.onVideoReceived(&video)
				break
			case halapi.HeartBeatType:
				context.onHeartBeatReceived(&heartbeat)
				break
			case halapi.LocationType:
				context.onLocationReceived(&loc)
				break
			default:
				log.Panicf("received unkown msg type = %v", sentType)
				break
			}
		}
	}
}

func (c *Context) onVideoReceived(v *halapi.VideoPart) {
	fmt.Printf("VideoPart: %v\n", *v)

	if !c.expectingVideoStream {
		log.Println("error, not expecting video stream, but received packet from HAL")
		return
	}

	// TODO: send packet to cmd
}

func (c *Context) onHeartBeatReceived(hb *halapi.HeartBeat) {
	fmt.Printf("HeartBeat: %v\n", *hb)
	// TODO: send heart beat to cmd
}

func (c *Context) onLocationReceived(loc *halapi.Location) {
	fmt.Printf("Location: %v\n", *loc)

	// TODO: send location to cmd
}

// Args store command line arguments
type Args struct {
	// null if no store path
	StorePath     string
	Port          int
	HALSocketPath string
}

func parseArgs() (*Args, error) {
	parser := argparse.NewParser("unit-daemon", "Unit client daemon")

	storePath := parser.String("s", "store", &argparse.Options{Help: "Path to archive video/positions/heartbeats. If not provided, won't store them.", Default: nil})
	port := parser.Int("p", "port", &argparse.Options{Help: "Main port the client will bind to, to receive connections from other clients.", Default: 4070})
	ctrlSocketPath := parser.String("", "hal-socket-path", &argparse.Options{Help: "Path to unix socket file to communicate over with HAL.", Default: "/tmp/unit.hal.sock"})

	err := parser.Parse(os.Args)
	if err != nil {
		return nil, errors.New(parser.Usage(err))
	}

	return &Args{
		StorePath:     *storePath,
		Port:          *port,
		HALSocketPath: *ctrlSocketPath,
	}, nil
}
