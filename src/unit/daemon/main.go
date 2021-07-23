package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"database/sql"

	"github.com/mido3ds/C4IAN/src/models"
	"github.com/mido3ds/C4IAN/src/unit/halapi"
)

const VideoStreamingNoEndTimeout = 1 * time.Minute

func testPeriodicCMDMsgs(context *Context, code int) {
	time.AfterFunc(2*time.Second, func() {
		sendCode := (code % 3) + 2
		context.onCodeMsgReceivedFromCMD(&models.Message{Code: sendCode})
		testPeriodicCMDMsgs(context, code+1)
	})
}

func testPeriodicCMDAudios(context *Context) {
	audioMock := []uint8{1, 2, 3, 4}
	time.AfterFunc(2*time.Second, func() {
		context.onAudioMsgReceivedFromCMD(&models.Audio{Body: audioMock})
		testPeriodicCMDAudios(context)
	})
}

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

	context := newContext(args)
	defer context.close()

	go context.listenHAL()
	go context.listenCmdTcp()
	go context.listenCmdUdp()
	go context.api.start(args.UISocket)

	log.Println("finished initalizing all")

	// testPeriodicCMDMsgs(context, 0)
	// testPeriodicCMDAudios(context)

	waitSIGINT()
}

type Context struct {
	Args
	name    string
	storeDB *sql.DB

	halConn net.Conn

	halMutex              sync.Mutex
	_expectingVideoStream bool
	isConnectedToHAL      bool

	videoMutex       sync.Mutex
	_videoSeqno      uint64
	videoID          int
	videoStreamTimer *time.Timer

	// videoManager *VideoFilesManager
	api *API
}

func newContext(args *Args) *Context {
	context := &Context{
		name:                  strings.Split(args.iface, "-")[0],
		Args:                  *args,
		_videoSeqno:           0,
		videoID:               0,
		_expectingVideoStream: false,
		isConnectedToHAL:      false,
	}

	// db
	if len(context.storePath) > 0 {
		context.storeDB = openDB(context.storePath)
		context.createTables()
	}

	// videoManager := NewVideoFilesManager("/tmp/unitvideos/assets/media")
	// context.videoManager = videoManager

	context.api = newAPI(context)

	return context
}

func (c *Context) close() {
	if c.storeDB != nil {
		c.storeDB.Close()
	}
}

func (c *Context) expectingVideoStream() bool {
	c.halMutex.Lock()
	defer c.halMutex.Unlock()
	return c._expectingVideoStream
}

func (c *Context) setExpectingVideoStream(v bool) {
	c.halMutex.Lock()
	defer c.halMutex.Unlock()
	c._expectingVideoStream = v
}

func (c *Context) resetVideoSeqNo() {
	c.videoMutex.Lock()
	defer c.videoMutex.Unlock()
	c._videoSeqno = 0
}

func (c *Context) incrementAndGetVideoSeqNo() uint64 {
	c.videoMutex.Lock()
	defer c.videoMutex.Unlock()
	v := c._videoSeqno
	c._videoSeqno++
	return v
}

func (c *Context) listenCmdTcp() {
	cmdLisener, err := net.Listen("tcp", fmt.Sprint("0.0.0.0:", c.port))
	if err != nil {
		log.Panic(err)
	}
	defer cmdLisener.Close()
	log.Println("listening for cmd on tcp port:", c.port)

	// go simulateCMD(c)

	for {
		conn, err := cmdLisener.Accept()
		if err != nil {
			log.Panic("failed to accept, err:", err)
		}

		log.Println("connected to cmd")

		go func() {
			defer conn.Close()

			var packetType models.Type
			var msg models.Message
			var audio models.Audio
			decoder := gob.NewDecoder(conn)

			err := decoder.Decode(&packetType)
			if err != nil {
				log.Println("failed to decode type, err:", err)
				return
			}

			if packetType == models.MessageType {
				err := decoder.Decode(&msg)
				if err != nil {
					log.Println("error occurred, will close connection with cmd:", err)
					return
				}
				c.onCodeMsgReceivedFromCMD(&msg)
			} else if packetType == models.AudioType {
				err := decoder.Decode(&audio)
				if err != nil {
					log.Println("error occurred, will close connection with cmd:", err)
					return
				}
				c.onAudioMsgReceivedFromCMD(&audio)
			} else {
				log.Panic("received unrecognized msg type on tcp")
			}
		}()
	}
}

func simulateCMD(c *Context) {
	for {
		time.Sleep(5 * time.Second)

		c.onCodeMsgReceivedFromCMD(&models.Message{Code: 5})
		c.onAudioMsgReceivedFromCMD(&models.Audio{Body: []byte{45, 23, 45, 1, 2, 4}})
	}
}

func (c *Context) listenCmdUdp() {
	addr := net.UDPAddr{
		Port: c.port,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()
	log.Println("listening for cmd on udp port:", c.port)

	var packetType models.Type
	var msg models.Message
	var audio models.Audio

	for {
		buffer := make([]byte, 64*1024)
		length, err := conn.Read(buffer)
		if err != nil {
			log.Panic(err)
		}
		decoder := gob.NewDecoder(bytes.NewBuffer(buffer[:length]))

		for decoder.Decode(&packetType) == nil {
			if packetType == models.MessageType {
				err := decoder.Decode(&msg)
				if err != nil {
					log.Panic("error occurred, will close connection with cmd:", err)
				}
				c.onCodeMsgReceivedFromCMD(&msg)
			} else if packetType == models.AudioType {
				err := decoder.Decode(&audio)
				if err != nil {
					log.Panic("error occurred, will close connection with cmd:", err)
				}
				c.onAudioMsgReceivedFromCMD(&audio)
			} else {
				log.Panic("received unrecognized msg type on tcp")
			}
		}
	}
}

func (c *Context) onCodeMsgReceivedFromCMD(msg *models.Message) {
	switch msg.Code {
	case StartVideoStreamingCode:
		log.Println("start video streaming code")
		c.setExpectingVideoStream(true)
		c.resetVideoSeqNo()
		if c.videoStreamTimer != nil {
			c.videoStreamTimer.Stop()
			c.videoStreamTimer = nil
		}
		c.videoStreamTimer = time.AfterFunc(VideoStreamingNoEndTimeout, func() {
			log.Println("didn't receive end video streaming for 1 minute, closing video streaming")
			c.setExpectingVideoStream(false)
			c.videoID++
		})
		c.api.sendCodeMsgEvent(msg.Code)
		break
	case StopVideStreamingCode:
		log.Println("stop video streaming code")
		if c.videoStreamTimer != nil {
			c.videoStreamTimer.Stop()
			c.videoStreamTimer = nil
		}
		c.setExpectingVideoStream(false)
		c.videoID++
		c.api.sendCodeMsgEvent(msg.Code)
		break
	case Hello:
		// Ignore hello messages
		break
	default:
		log.Println("generic code msg: ", msg.Code)

		c.halMutex.Lock()
		defer c.halMutex.Unlock()

		if c.isConnectedToHAL {
			err := halapi.ShowCodeMsg{Code: msg.Code}.Write(c.halConn)
			if err != nil {
				log.Println("error in sending code msg to hal:", err)
			}
		}
		c.api.sendCodeMsgEvent(msg.Code)
		break
	}
}

func (c *Context) onAudioMsgReceivedFromCMD(audio *models.Audio) {
	log.Println("Received an audio with len= ", len(audio.Body))

	c.halMutex.Lock()
	defer c.halMutex.Unlock()

	if c.isConnectedToHAL {
		err := halapi.ShowAudioMsg{Audio: audio.Body}.Write(c.halConn)
		if err != nil {
			log.Println("error in sending show audio msg to HAL, err:", err)
		}
	} else {
		log.Println("received msg but coudln't connect to HAL to play it, dropping msg")
	}

	c.api.sendAudioMsgEvent(audio)
}

func (c *Context) sendVideoFragmentUDP(fragment, metadata []byte, filename string) error {
	conn, err := net.DialTimeout("tcp", c.cmdAddress, c.timeout)
	if err != nil {
		return fmt.Errorf("couldn't open tcp port, err: %v", err)
	}
	defer conn.Close()

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err = encoder.Encode(models.VideoFragmentType)
	if err != nil {
		return fmt.Errorf("failed to encode type, error: %v", err)
	}

	err = encoder.Encode(&models.VideoFragment{
		ID:       c.videoID,
		Time:     time.Now().UnixNano(),
		Body:     fragment,
		Metadata: metadata,
		FileName: filename,
	})

	if err != nil {
		return fmt.Errorf("failed to encode fragment, err: %v", err)
	}

	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send bytes, err: %v", err)
	}

	return nil
}

func (c *Context) sendSensorDataUDP(lon, lat float64, beatsPerMinute int) error {
	udpConn, err := net.DialTimeout("udp", c.cmdAddress, c.timeout)
	if err != nil {
		log.Panic("couldn't open udp port, err:", err)
	}
	defer udpConn.Close()

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err = encoder.Encode(models.SensorDataType)
	if err != nil {
		return fmt.Errorf("failed to encode type, error: %v", err)
	}

	err = encoder.Encode(&models.SensorData{
		Time:      time.Now().UnixNano(),
		Heartbeat: beatsPerMinute,
		Lat:       lat,
		Lon:       lon,
	})
	if err != nil {
		return fmt.Errorf("failed to encode sensor data, err: %v", err)
	}

	_, err = udpConn.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send bytes, err: %v", err)
	}

	return nil
}

func (c *Context) sendAudioMessageTCP(fragment []byte) error {
	conn, err := net.DialTimeout("tcp", c.cmdAddress, c.timeout)
	if err != nil {
		return fmt.Errorf("couldn't open tcp port, err: %v", err)
	}
	defer conn.Close()

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err = encoder.Encode(models.AudioType)
	if err != nil {
		return fmt.Errorf("failed to encode type, error: %v", err)
	}

	err = encoder.Encode(&models.Audio{
		Time: time.Now().UnixNano(),
		Body: fragment,
	})
	if err != nil {
		return fmt.Errorf("failed to encode fragment, err: %v", err)
	}

	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send bytes, err: %v", err)
	}

	return nil
}

func (c *Context) sendCodeMessageTCP(code int) error {
	conn, err := net.DialTimeout("tcp", c.cmdAddress, c.timeout)
	if err != nil {
		return fmt.Errorf("couldn't open tcp port, err: %v", err)
	}
	defer conn.Close()

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err = encoder.Encode(models.MessageType)
	if err != nil {
		return fmt.Errorf("failed to encode type, error: %v", err)
	}

	err = encoder.Encode(&models.Message{
		Code: code,
		Time: time.Now().UnixNano(),
	})
	if err != nil {
		return fmt.Errorf("failed to encode code message, err: %v", err)
	}

	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send bytes, err: %v", err)
	}

	return nil
}
