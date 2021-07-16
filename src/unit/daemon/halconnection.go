package main

import (
	"encoding/gob"
	"log"
	"net"
	"os"
	"time"

	"github.com/mido3ds/C4IAN/src/unit/daemon/halapi"
)

// TODO: remove
func simulateHALClient(HALSocketPath string) {
	conn, err := net.Dial("unix", HALSocketPath)
	if err != nil {
		log.Panic(err)
	}

	enc := gob.NewEncoder(conn)

	for {
		time.Sleep(time.Second * 2)
		log.Println("hal sending data")

		halapi.SensorData{
			HeartBeat: halapi.HeartBeat{BeatsPerMinut: 5},
			Location:  halapi.Location{Lon: 5, Lat: 2},
		}.Send(enc)
		halapi.VideoFragment{Video: []byte{5, 3}}.Send(enc)
		halapi.AudioMsg{Audio: []byte{10, 2}}.Send(enc)
		halapi.CodeMsg{Code: 6}.Send(enc)
	}
}

func (c *Context) listenHAL() {
	// remove loc socket file
	err := os.RemoveAll(c.halSocketPath)
	if err != nil {
		log.Panic("failed to remove socket:", err)
	}

	halListener, err := net.Listen("unix", c.halSocketPath)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer halListener.Close()
	log.Println("listening for HAL connection over unix socket:", c.halSocketPath)

	for {
		go simulateHALClient(c.halSocketPath)
		conn, err := halListener.Accept()
		c.halConn = conn
		defer conn.Close()

		if err != nil {
			log.Println("accept error:", err)
		} else {
			c.setIsConnectedToHAL(true)
			defer c.setIsConnectedToHAL(false)

			log.Println("HAL connected")
			c.serveHAL(conn)
		}
	}
}

func (context *Context) serveHAL(conn net.Conn) {
	dec := gob.NewDecoder(conn)

	var video halapi.VideoFragment
	var sensorData halapi.SensorData
	var audiomsg halapi.AudioMsg
	var codemsg halapi.CodeMsg

	for {
		sentType, err := halapi.RecvFromHAL(dec, &video, &sensorData, &audiomsg, &codemsg)
		if err != nil {
			log.Println(err)
		} else {
			switch sentType {
			case halapi.VideoFragmentType:
				context.onVideoReceivedFromHAL(&video)
				break
			case halapi.SensorDataType:
				context.onSensorDataReceivedFromHAL(&sensorData)
				break
			case halapi.AudioMsgType:
				context.onAudioMsgReceivedFromHAL(&audiomsg)
				break
			case halapi.CodeMsgType:
				context.onCodeMsgReceivedFromHAL(&codemsg)
				break
			default:
				log.Panicf("received unkown msg type = %v", sentType)
				break
			}
		}
	}
}

func (c *Context) onCodeMsgReceivedFromHAL(cm *halapi.CodeMsg) {
	log.Printf("From HAL:: CodeMsg= %v\n", *cm)

	err := c.sendCodeMessageTCP(cm.Code)
	if err != nil {
		log.Println("failed to send code message, err:", err)
		log.Println("will try again")
		time.AfterFunc(c.retryOrCloseTimeout, func() {
			err := c.sendCodeMessageTCP(cm.Code)
			if err != nil {
				log.Println("failed to send code message, err:", err)
			}
		})
	}
}

func (c *Context) onAudioMsgReceivedFromHAL(a *halapi.AudioMsg) {
	log.Printf("From HAL:: AudioMsg= %v\n", *a)

	err := c.sendAudioMessageTCP(a.Audio)
	if err != nil {
		log.Println("failed to send audio message, err:", err)
		log.Println("will try again")
		time.AfterFunc(c.retryOrCloseTimeout, func() {
			err := c.sendAudioMessageTCP(a.Audio)
			if err != nil {
				log.Println("failed to send audio message, err:", err)
			}
		})
	}
}

func (c *Context) onVideoReceivedFromHAL(v *halapi.VideoFragment) {
	log.Printf("From HAL:: VideoFragment= %v\n", *v)
	c.saveVideoFragment(v.Video)

	if !c.expectingVideoStream() {
		log.Println("received video fragment from HAL, but CMD didn't ask for it, dropping it")
		return
	}

	err := c.sendVideoFragmentUDP(v.Video)
	if err != nil {
		log.Println("error in sending video frag:", err)
	}
}

func (c *Context) onSensorDataReceivedFromHAL(s *halapi.SensorData) {
	hb := s.HeartBeat
	loc := s.Location

	log.Printf("From HAL:: HeartBeat= %v\n", hb)
	log.Printf("From HAL:: Location= %v\n", loc)

	c.saveHeartbeat(hb.BeatsPerMinut)
	c.saveLocation(loc.Lon, loc.Lat)

	err := c.sendSensorDataUDP(loc.Lon, loc.Lat, hb.BeatsPerMinut)
	if err != nil {
		log.Println("error in sending sensor data to cmd:", err)
	}
}
