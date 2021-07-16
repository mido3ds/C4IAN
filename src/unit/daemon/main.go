package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"database/sql"

	"github.com/mido3ds/C4IAN/src/models"
	"github.com/mido3ds/C4IAN/src/unit/daemon/halapi"
)

const VideoStreamingNoEndTimeout = 1 * time.Minute

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
	defer context.close()

	go context.listenHAL()
	go context.listenCmdTcp()
	go context.listenCmdUdp()

	log.Println("finished initalizing all")

	waitSIGINT()
}

type Context struct {
	Args
	storeDB *sql.DB

	cmdUDPConn net.Conn
	halConn    net.Conn

	halMutex              sync.Mutex
	_expectingVideoStream bool
	_isConnectedToHAL     bool

	videoMutex       sync.Mutex
	_videoSeqno      uint64
	videoID          int
	videoStreamTimer *time.Timer
}

func newContext(args *Args) Context {
	context := Context{
		Args:                  *args,
		_videoSeqno:           0,
		videoID:               0,
		_expectingVideoStream: false,
		_isConnectedToHAL:     false,
	}

	// db
	if len(context.storePath) > 0 {
		context.storeDB = openDB(context.storePath)
		context.createTables()
	}

	udpConn, err := net.DialTimeout("udp", context.cmdAddress, context.timeout)
	if err != nil {
		log.Panic("couldn't open udp port, err:", err)
	}
	context.cmdUDPConn = udpConn

	return context
}

func (c *Context) close() {
	if c.storeDB != nil {
		c.storeDB.Close()
	}
	c.cmdUDPConn.Close()
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

func (c *Context) isConnectedToHAL() bool {
	c.halMutex.Lock()
	defer c.halMutex.Unlock()
	return c._isConnectedToHAL
}

func (c *Context) setIsConnectedToHAL(v bool) {
	c.halMutex.Lock()
	defer c.halMutex.Unlock()
	c._isConnectedToHAL = v
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
	decoder := gob.NewDecoder(conn)

	for {
		err := decoder.Decode(&packetType)
		if err != nil {
			log.Panic("failed to decode type, err:", err)
		}

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

func (c *Context) onCodeMsgReceivedFromCMD(msg *models.Message) {
	log.Println("Received a msg: ", msg)
	switch msg.Code {
	case StartVideoStreamingCode:
		log.Println("start video streaming code")
		c.setExpectingVideoStream(true)
		c.videoID++
		c.resetVideoSeqNo()
		c.videoStreamTimer = time.AfterFunc(VideoStreamingNoEndTimeout, func() {
			log.Println("didn't receive end video streaming for 1 minute, closing video streaming")
			c.setExpectingVideoStream(false)
		})
		break
	case StopVideStreamingCode:
		log.Println("stop video streaming code")
		if c.videoStreamTimer != nil {
			c.videoStreamTimer.Stop()
			c.videoStreamTimer = nil
		}
		c.setExpectingVideoStream(false)
		break
	default:
		log.Println("generic code msg")
		break
	}
}

func (c *Context) onAudioMsgReceivedFromCMD(audio *models.Audio) {
	log.Println("Received an audio with len= ", len(audio.Body))
	if c.isConnectedToHAL() {
		enc := gob.NewEncoder(c.halConn)
		err := halapi.ShowAudioMsg{Audio: audio.Body}.Send(enc)
		if err != nil {
			log.Println("error in sending show audio msg to HAL, err:", err)
		}
	} else {
		log.Println("received msg but coudln't connect to HAL to play it, dropping msg")
	}
}

func (c *Context) sendVideoFragmentUDP(fragment []byte) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(models.VideoFragmentType)
	if err != nil {
		return fmt.Errorf("failed to encode type, error: %v", err)
	}

	err = encoder.Encode(&models.VideoFragment{
		ID:    c.videoID,
		Time:  time.Now().Unix(),
		Body:  fragment,
		SeqNo: c.incrementAndGetVideoSeqNo(),
	})

	if err != nil {
		return fmt.Errorf("failed to encode fragment, err: %v", err)
	}

	n, err := c.cmdUDPConn.Write(buffer.Bytes())
	if n != len(buffer.Bytes()) {
		return fmt.Errorf("failed to send all bytes")
	} else if err != nil {
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
		Time: time.Now().Unix(),
		Body: fragment,
	})
	if err != nil {
		return fmt.Errorf("failed to encode fragment, err: %v", err)
	}

	n, err := conn.Write(buffer.Bytes())
	if n != len(buffer.Bytes()) {
		return fmt.Errorf("failed to send all bytes")
	} else if err != nil {
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
		Time: time.Now().Unix(),
	})
	if err != nil {
		return fmt.Errorf("failed to encode code message, err: %v", err)
	}

	_, err = c.cmdUDPConn.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send bytes, err: %v", err)
	}

	return nil
}

func (c *Context) sendSensorDataUDP(lon, lat float64, beatsPerMinute int) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(models.SensorDataType)
	if err != nil {
		return fmt.Errorf("failed to encode type, error: %v", err)
	}

	err = encoder.Encode(&models.SensorData{
		Time:      time.Now().Unix(),
		Heartbeat: beatsPerMinute,
		Lat:       lat,
		Lon:       lon,
	})
	if err != nil {
		return fmt.Errorf("failed to encode sensor data, err: %v", err)
	}

	_, err = c.cmdUDPConn.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send bytes, err: %v", err)
	}

	return nil
}
