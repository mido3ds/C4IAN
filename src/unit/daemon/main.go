package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
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
	expectingVideoStream   bool
	storeDB                *sql.DB
	cmdUDPConn, cmdTCPConn net.Conn
	videoSeqno             uint64
	videoID                int
	isConnectedToCMD       bool
	halConn                net.Conn
	isConnectedToHAL       bool
	videoStreamTimer       *time.Timer
}

func newContext(args *Args) Context {
	context := Context{
		Args:                 *args,
		expectingVideoStream: false,
		videoSeqno:           0,
		videoID:              0,
		isConnectedToCMD:     false,
		isConnectedToHAL:     false,
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
	c.closeConnectionWithCMD()
}

func (c *Context) closeConnectionWithCMD() {
	if c.isConnectedToCMD {
		c.isConnectedToCMD = false
		c.cmdTCPConn.Close()
	}
}

func (c *Context) tryConnectWithCMD() {
	// tcp
	tcpConn, err := net.DialTimeout("tcp", c.cmdAddress, c.timeout)
	if err != nil {
		log.Println("couldn't open tcp port, err:", err)
		return
	}
	c.cmdTCPConn = tcpConn
	c.isConnectedToCMD = true
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
			log.Println("failed to accept, err:", err)
		}
		defer conn.Close()

		if !c.isConnectedToCMD {
			c.isConnectedToCMD = true
			c.cmdTCPConn = conn
		}

		c.serveCmd(conn)
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
	log.Println("listenign for cmd on udp port:", c.port)

	c.serveCmd(conn)
}

func (c *Context) serveCmd(conn net.Conn) {
	var packetType models.Type
	var msg models.Message
	var audio models.Audio
	decoder := gob.NewDecoder(conn)

	for {
		err := decoder.Decode(&packetType)
		if err != nil {
			log.Panic(err)
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
	}
}

func (c *Context) onCodeMsgReceivedFromCMD(msg *models.Message) {
	log.Println("Received a msg: ", msg)
	switch msg.Code {
	case StartVideoStreamingCode:
		log.Println("start video streaming code")
		c.expectingVideoStream = true
		c.videoID++
		c.videoSeqno = 0
		c.videoStreamTimer = time.AfterFunc(VideoStreamingNoEndTimeout, func() {
			log.Println("didn't receive end video streaming for 1 minute, closing video streaming")
			c.expectingVideoStream = false
		})
		break
	case StopVideStreamingCode:
		log.Println("stop video streaming code")
		if c.videoStreamTimer != nil {
			c.videoStreamTimer.Stop()
			c.videoStreamTimer = nil
		}
		c.expectingVideoStream = false
		break
	default:
		log.Println("generic code msg")
		break
	}
}

func (c *Context) onAudioMsgReceivedFromCMD(audio *models.Audio) {
	log.Println("Received an audio with len= ", len(audio.Body))
	if c.isConnectedToHAL {
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
		SeqNo: c.videoSeqno,
	})
	c.videoSeqno++

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
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(models.AudioType)
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

	n, err := c.cmdTCPConn.Write(buffer.Bytes())
	if n != len(buffer.Bytes()) {
		return fmt.Errorf("failed to send all bytes")
	} else if err != nil {
		return fmt.Errorf("failed to send bytes, err: %v", err)
	}

	return nil
}

func (c *Context) sendCodeMessageTCP(code int) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(models.MessageType)
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

	n, err := c.cmdTCPConn.Write(buffer.Bytes())
	if n != len(buffer.Bytes()) {
		return fmt.Errorf("failed to send all bytes")
	} else if err != nil {
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

	n, err := c.cmdUDPConn.Write(buffer.Bytes())
	if n != len(buffer.Bytes()) {
		return fmt.Errorf("failed to send all bytes")
	} else if err != nil {
		return fmt.Errorf("failed to send bytes, err: %v", err)
	}

	return nil
}
