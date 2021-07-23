package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/mido3ds/C4IAN/src/unit/halapi"
)

func main() {
	args, err := parseArgs()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Println(args)
	defer log.Println("finished cleaning up, closing")

	var iface string
	if len(args.iface) > 0 {
		iface = args.iface
	} else {
		iface = getDefaultInterface()
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetOutput(os.Stdout)
	log.SetPrefix("[" + iface + "] ")

	// Create a random seed for rand based on the interface name
	var seed int64 = 0
	for i, char := range iface {
		seed += int64(i * int(char))
	}
	rand.Seed(seed)

	context := newContext(args)

	// go context.sendAudioMsgs()
	// go context.sendCodeMsgs()

	if context.sensors {
		go context.sendSensorsData()
	}

	go context.streamVideo()

	go context.receiveMsgs()

	log.Println("finished initalizing all")

	waitSIGINT()
}

type Context struct {
	Args
	videoPath   string
	audiosFiles []string
	halConn     net.Conn
	tempDir     string
	lastTSIndex int
	locAgent    *LocAgent
}

func newContext(args *Args) Context {
	videoPath := args.videoPath
	if len(videoPath) > 0 {
		log.Println("videoPath:", videoPath)
	} else {
		log.Println("no video to stream")
	}

	audiosFiles := []string{}
	if len(args.audiosDirPath) > 0 {
		audiosFiles = listDir(args.audiosDirPath)
		log.Println("audios:", audiosFiles)
	} else {
		log.Println("no audio dir provided")
	}

	tempdir, err := ioutil.TempDir("", "halsimulation.")
	if err != nil {
		log.Panic(err)
	}
	shouldCloseDir := true
	defer func() {
		if shouldCloseDir {
			os.RemoveAll(tempdir)
		}
	}()

	conn, err := net.Dial("unix", args.halSocketPath)
	if err != nil {
		log.Panic(err)
	}
	shouldCloseDir = false

	var locAgent *LocAgent = nil
	if len(args.locationSocket) > 0 {
		locAgent, err = newLocAgent(args.locationSocket)
		if err != nil {
			log.Panic(err)
		}
	}

	return Context{
		Args:        *args,
		videoPath:   videoPath,
		audiosFiles: audiosFiles,
		halConn:     conn,
		tempDir:     tempdir,
		lastTSIndex: 0,
		locAgent:    locAgent,
	}
}

func (c *Context) close() {
	c.halConn.Close()
	os.RemoveAll(c.tempDir)
}

func (c *Context) sendAudioMsgs() {
	if len(c.audiosDirPath) == 0 {
		log.Println("no audio dir, won't send audio messages")
		return
	}

	// every rand(avg=2s, stdev=300ms): send(rand(audio msg))
	for {
		time.Sleep(time.Duration(normal(5, 2)) * time.Second)

		audioBuffer, err := ioutil.ReadFile(c.audiosFiles[rand.Intn(len(c.audiosFiles))])
		if err != nil {
			log.Panic(err)
		}

		err = halapi.AudioMsg{
			Audio: audioBuffer,
		}.Write(c.halConn)
		if err != nil {
			log.Panic(err)
		}
	}
}

func (c *Context) streamVideo() {
	if len(c.videoPath) == 0 {
		log.Println("no video path, won't stream videos")
		return
	}

	// in video streaming mode: send index.m3u8 with last fragment
	m3u8Path := path.Join(c.tempDir, "index.m3u8")
	f, err := os.Create(m3u8Path)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	go c.watchM3U8(m3u8Path)
	runFFmpeg(c.ffmpegPath, c.videoPath, m3u8Path, c.tempDir, c.fragmentDurSecs)
}

func (c *Context) watchM3U8(m3u8path string) {
	for {
		time.Sleep(5 * time.Millisecond)
		c.sendM3U8(m3u8path)
	}
}

func getNumTS(m3u8 []byte) int {
	reg := regexp.MustCompile(`index(\d)+\.ts`)
	b := reg.FindAll(m3u8, -1)
	size := 0
	if b != nil {
		size = len(b)
	}
	return size
}

func (c *Context) sendM3U8(m3u8path string) {
	// load m3u8 file
	m3u8, err := ioutil.ReadFile(m3u8path)
	if err != nil {
		log.Panic(err)
	}

	// get latest ts file(s)
	numTS := getNumTS(m3u8)
	numTSToSend := numTS - c.lastTSIndex

	tsfiles := make([][]byte, 0)
	filenames := make([]string, 0)
	for i := c.lastTSIndex; i < numTS; i++ {
		name := fmt.Sprintf("index%d.ts", i)
		path := path.Join(c.tempDir, name)
		filenames = append(filenames, name)

		bts, err := ioutil.ReadFile(path)
		if err != nil {
			log.Panic(err)
		}

		tsfiles = append(tsfiles, bts)
	}

	strm3u8 := string(m3u8)

	// send video fragment(s)
	for i := 0; i < len(tsfiles); i++ {
		splittedFile := strings.SplitAfter(strm3u8, filenames[i])[0]

		err := halapi.VideoFragment{
			Video:    tsfiles[i],
			Metadata: []byte(splittedFile),
			Filename: filenames[i],
		}.Write(c.halConn)
		if err != nil {
			log.Panic("failed to send video fragment, err:", err)
		}
	}

	// increment counter to latest ts file(s)
	c.lastTSIndex += numTSToSend

	//log.Println("numTs", numTS)
	if numTSToSend > 0 {
		log.Println("sent", numTSToSend, "TSs, filenames:", filenames) // TODO: remove
	}
}

func runFFmpeg(ffmpegPath, videoPath, m3u8Path, outdir string, fragmentDurSecs int) {
	args := []string{
		`-i`, videoPath,
		`-framerate`, `20`,
		`-s`, `480x360`,
		`-level`, `3.0`,
		`-fs`, `6500`,
		`-start_number`, `0`,
		`-f`, `hls`,
		`-hls_time`, fmt.Sprint(fragmentDurSecs),
		`-hls_playlist_type`, `event`,
		`-hls_flags`, `independent_segments`,
		`-hls_segment_type`, `mpegts`,
		`-hls_list_size`, `0`,
		m3u8Path,
	}
	cmd := exec.Command(ffmpegPath, args...)
	log.Println("executing: ", cmd)

	stderrpath := path.Join(outdir, "stderr")
	stderr, err := os.Create(stderrpath)
	if err != nil {
		log.Panic(err)
	}
	defer stderr.Close()
	cmd.Stderr = stderr
	log.Println("ffmpeg stderr path:", stderrpath)

	stdoutpath := path.Join(outdir, "stdout")
	stdout, err := os.Create(stdoutpath)
	if err != nil {
		log.Panic(err)
	}
	defer stdout.Close()
	cmd.Stdout = stdout
	log.Println("ffmpeg stdout path:", stdoutpath)

	err = cmd.Run()
	if err != nil {
		stderr, _ := ioutil.ReadFile(stderrpath)
		stdout, _ := ioutil.ReadFile(stdoutpath)
		log.Panic("error:", err, ", stderr:", string(stderr), ", stdout:", string(stdout))
	}
}

func (c *Context) sendCodeMsgs() {
	// every rand(avg=3s, stdev=1s): send(rand(number for code msg))
	for {
		time.Sleep(time.Duration(normal(6, 3)) * time.Second)

		err := halapi.CodeMsg{
			Code: rand.Intn(400),
		}.Write(c.halConn)
		if err != nil {
			log.Panic(err)
		}
	}
}

func (c *Context) sendSensorsData() {
	var lon float64
	var lat float64

	if c.locAgent != nil {
		go c.locAgent.start()
	}

	// every 3s to 7s: send(location=rand(avg=(lon,lat), stdev=(.02,.03)),heartbeat=rand(avg=70,stdev=30))
	for {
		time.Sleep(time.Duration(uniform(4, 8)) * time.Second)

		hb := int(normal(80, 15))

		if c.locAgent != nil {
			loc := c.locAgent.location
			lon = loc.Lon
			lat = loc.Lat
		} else {
			lon = normal(32.4, .02)
			lat = normal(43.098, .03)
		}

		err := halapi.SensorData{
			Location: halapi.Location{
				Lon: lon,
				Lat: lat,
			},
			HeartBeat: halapi.HeartBeat{
				BeatsPerMinut: hb,
			},
		}.Write(c.halConn)

		if err != nil {
			log.Panic(err)
		}
		//}
	}
}

func (c *Context) receiveMsgs() {
	var svs halapi.StartVideoStream
	var evs halapi.EndVideoStream
	var sam halapi.ShowAudioMsg
	var scm halapi.ShowCodeMsg
	var timer *time.Timer

	for {
		receivedType, err := halapi.ReadFromUnit(c.halConn, &svs, &evs, &sam, &scm)
		if err != nil {
			log.Panic(err)
		}

		switch receivedType {
		case halapi.StartVideoStreamType:
			log.Println("started vidoe streaming")
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(time.Minute, func() {
				log.Println("timeouted in video streaming")

				if !c.live {
					log.Println("resettign lastTSIndex")
					c.lastTSIndex = 0
				}
			})
			break
		case halapi.EndVideoStreamType:
			log.Println("ended video streaming")
			if !c.live {
				log.Println("resettign lastTSIndex")
				c.lastTSIndex = 0
			}
			break
		case halapi.ShowAudioMsgType:
			onRecievedAudioMsg(sam.Audio, c.tempDir)
			break
		case halapi.ShowCodeMsgType:
			onReceivedCodeMsg(scm.Code)
			break
		}
	}
}

func onReceivedCodeMsg(code int) {
	log.Println("Received code msg:", code)
}

func onRecievedAudioMsg(audio []byte, tmpDir string) {
	// save any audio msg in getTmpFile() and print path
	path, err := ioutil.TempFile(tmpDir, "audio.")
	if err != nil {
		log.Panic(path)
	}
	defer path.Close()

	n, err := path.Write(audio)
	if n != len(audio) {
		log.Panic("didn't save all bytes of audio msg")
	}
	if err != nil {
		log.Panic(err)
	}

	log.Println("saved audio msg to:", path.Name())
}
