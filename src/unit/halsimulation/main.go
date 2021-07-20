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
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mido3ds/C4IAN/src/unit/halapi"
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
	videoPath        string
	audiosFiles      []string
	halConn          net.Conn
	videoStreamingOn bool
	tempDir          string
	lastTSIndex      int
}

func newContext(args *Args) Context {
	videoPath := args.videoPath
	log.Println("videoPath:", videoPath)

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

	return Context{
		Args:             *args,
		videoPath:        videoPath,
		audiosFiles:      audiosFiles,
		halConn:          conn,
		videoStreamingOn: false,
		tempDir:          tempdir,
		lastTSIndex:      0,
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
		time.Sleep(time.Duration(normal(2, 0.3)) * time.Second)

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
	go runFFmpeg(c.ffmpegPath, c.videoPath, m3u8Path, c.tempDir, c.fragmentDurSecs)

	// every 10s: start video streaming mode (which lasts for 10s)
	for {
		time.Sleep(10 * time.Second)
		c.videoStreamingOn = !c.videoStreamingOn
	}
}

func (c *Context) watchM3U8(m3u8path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panic(err)
	}
	defer watcher.Close()

	err = watcher.Add(m3u8path)
	if err != nil {
		log.Panic(err)
	}

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			if c.videoStreamingOn {
				c.sendM3U8(m3u8path)
			}
		case event, _ := <-watcher.Events:
			if (event.Op&fsnotify.Write) == fsnotify.Write && c.videoStreamingOn {
				c.sendM3U8(m3u8path)
			}
		case err, _ := <-watcher.Errors:
			log.Panic(err)
		}
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
	log.Println("sending m3u8")

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

	// send video fragment(s)
	for i := 0; i < len(tsfiles); i++ {
		err := halapi.VideoFragment{
			Video:    tsfiles[i],
			Metadata: m3u8,
			Filename: filenames[i],
		}.Write(c.halConn)
		if err != nil {
			log.Panic("failed to send video fragment, err:", err)
		}
	}

	// increment counter to latest ts file(s)
	c.lastTSIndex += numTSToSend

	log.Println("sent", numTSToSend, "TSs, filenames:", filenames) // TODO: remove
}

func runFFmpeg(ffmpegPath, videoPath, m3u8Path, outdir string, fragmentDurSecs int) {
	args := []string{
		`-i`, videoPath,
		`-framerate`, `60`,
		`-s`, `480x360`,
		`-level`, `3.0`,
		`-fs`, `6500`,
		`-start_number`, `0`,
		`-f`, `hls`,
		`-hls_time`, fmt.Sprint(fragmentDurSecs),
		`-hls_playlist_type`, `event`,
		`-hls_flags`, `independent_segments`,
		`-hls_flags`, `split_by_time`,
		`-hls_segment_type`, `mpegts`,
		`-hls_list_size`, `5`,
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
		time.Sleep(time.Duration(normal(3, 1)) * time.Second)

		err := halapi.CodeMsg{
			Code: rand.Intn(400),
		}.Write(c.halConn)
		if err != nil {
			log.Panic(err)
		}
	}
}

func (c *Context) sendSensorsData() {
	// every 10s with probabliy=60%: send(location=rand(avg=(lon,lat), stdev=(.02,.03)),heartbeat=rand(avg=70,stdev=20))
	for {
		time.Sleep(10 * time.Second)

		if rand.Intn(100) < 60 {
			lon := normal(32.4, .02)
			lat := normal(43.098, .03)
			hb := int(normal(70, 20))

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
		}
	}
}

func (c *Context) receiveMsgs() {
	var svs halapi.StartVideoStream
	var evs halapi.EndVideoStream
	var sam halapi.ShowAudioMsg
	var scm halapi.ShowCodeMsg

	for {
		receivedType, err := halapi.ReadFromUnit(c.halConn, &svs, &evs, &sam, &scm)
		if err != nil {
			log.Panic(err)
		}

		switch receivedType {
		case halapi.StartVideoStreamType:
			log.Println("started vidoe streaming")
			c.videoStreamingOn = true
			break
		case halapi.EndVideoStreamType:
			log.Println("ended video streaming")
			c.videoStreamingOn = false
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
