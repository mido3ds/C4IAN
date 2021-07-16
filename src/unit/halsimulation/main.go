package main

import (
	"fmt"
	"log"
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
	go context.send()
	go context.streamVideo()
	go context.showMsgs()

	log.Println("finished initalizing all")

	waitSIGINT()
}

type Context struct {
	Args
	videosFiles []string
	audiosFiles []string
}

func newContext(args *Args) Context {
	videosFiles := listDir(args.videosDirPath)
	log.Println("videos:", videosFiles)

	audiosFiles := listDir(args.audiosDirPath)
	log.Println("audios:", audiosFiles)

	return Context{
		Args:        *args,
		videosFiles: videosFiles,
		audiosFiles: audiosFiles,
	}
}

func (c *Context) send() {
	// TODO
	// open connection
	// every rand(avg=2s, stdev=300ms): send(rand(audio msg))
	// every rand(avg=3s, stdev=1s): send(rand(number for code msg))
	// every 10s with probabliy=60%: send(location=rand(avg=(lon,lat), stdev=(.02,.03)),heartbeat=rand(avg=70,stdev=20))
	// every 10s: toggle video streaming mode, video=rand(videos)
}

func (c *Context) streamVideo() {
	// TODO
	// in video streaming mode: send next fragment (what size?)
}

func (c *Context) showMsgs() {
	// TODO
	// print any code msg
	// save any audio msg in getTmpFile() and print path
}
