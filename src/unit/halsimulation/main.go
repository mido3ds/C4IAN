package main

import (
	"log"

	"github.com/mido3ds/C4IAN/src/unit/halapi"
)

type Args struct {
	halSocketPath string
	videosDirPath string
	audiosDirPath string
}

func main() {
	log.Println("hello simulation")
	log.Println(halapi.VideoFragment{Video: []byte{1, 2, 3}})

	// read halsocketpath
	// read video file path
	// read audio file paths

	// open video
	// open audio

	// goroutine:
	// open connection
	// every rand(avg=2s, stdev=300ms): send(rand(audio msg))
	// every rand(avg=3s, stdev=1s): send(rand(number for code msg))
	// every 10s with probabliy=60%: send(location=rand(avg=(lon,lat), stdev=(.02,.03)),heartbeat=rand(avg=70,stdev=20))
	// every 10s: toggle video streaming mode, video=rand(videos)

	/// goroutine:
	// in video streaming mode: send next fragment (what size?)

	/// goroutine:
	// print any code msg
	// save any audio msg in getTmpFile() and print path
}
