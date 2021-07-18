package main

import (
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/mido3ds/C4IAN/src/unit/halapi"
)

func main() {
	defer log.Println("finished cleaning up, closing")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetOutput(os.Stdout)

	args, err := parseArgs()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Println(args)

	context := newContext(args)

	go context.sendAudioMsgs()
	go context.sendCodeMsgs()
	go context.sendSensorsData()
	go context.streamVideo()

	go context.receiveMsgs()

	log.Println("finished initalizing all")

	waitSIGINT()
}

type Context struct {
	Args
	videosFiles []string
	audiosFiles []string
	halConn     net.Conn
}

func newContext(args *Args) Context {
	videosFiles := listDir(args.videosDirPath)
	log.Println("videos:", videosFiles)

	audiosFiles := listDir(args.audiosDirPath)
	log.Println("audios:", audiosFiles)

	conn, err := net.Dial("unix", args.halSocketPath)
	if err != nil {
		log.Panic(err)
	}

	return Context{
		Args:        *args,
		videosFiles: videosFiles,
		audiosFiles: audiosFiles,
		halConn:     conn,
	}
}

func (c *Context) close() {
	c.halConn.Close()
}

func (c *Context) sendAudioMsgs() {
	// TODO
	// every rand(avg=2s, stdev=300ms): send(rand(audio msg))
}

func (c *Context) streamVideo() {
	// TODO
	// every 10s: start video streaming mode (which lasts for 10s)
	// in video streaming mode: send index.m3u8 with last fragment
}

func (c *Context) sendCodeMsgs() {
	// TODO
	// every rand(avg=3s, stdev=1s): send(rand(number for code msg))
}

func (c *Context) sendSensorsData() {
	// TODO
	// every 10s with probabliy=60%: send(location=rand(avg=(lon,lat), stdev=(.02,.03)),heartbeat=rand(avg=70,stdev=20))
}

func (c *Context) receiveMsgs() {
	for {
		// TODO
	}
}

func onReceivedCodeMsg(code int) {
	// TODO
	// print any code msg
}

func onRecievedAudioMsg(audio []byte) {
	// TODO
	// save any audio msg in getTmpFile() and print path
}
