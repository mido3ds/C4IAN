package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/mido3ds/C4IAN/src/models"
)

var cmdPort int = 4170
var unitPort int = 4070

func TestUnit(*testing.T) {
	go Listen(unitPort)
	i := 0
	for {
		SendMessage(i, cmdPort)
		SendAudio(i, cmdPort)
		SendSensorsData(i, cmdPort)
		SendVideoFragment(i, cmdPort)
		i++
		time.Sleep(time.Second)
	}
}

func Listen(port int) {
	address, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Panic(err)
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		log.Panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go func() {
			decoder := gob.NewDecoder(conn)
			var packetType models.Type
			err := decoder.Decode(&packetType)
			if err != nil {
				log.Panic(err)
			}

			// Decode the payload of the packet and make appropriate callbacks
			if packetType == models.MessageType {
				var msg models.Message
				err := decoder.Decode(&msg)
				if err != nil {
					log.Panic(err)
				}
				log.Println("Test: Received a msg: ", msg)
			} else {
				var audio models.Audio
				err := decoder.Decode(&audio)
				if err != nil {
					log.Panic(err)
				}
				log.Println("Test: Received an audio: ", audio)
			}
			conn.Close()
		}()
	}
}

func SendSensorsData(i int, port int) {
	address, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Panic(err)
	}

	conn, err := net.DialUDP("udp", nil, address)
	if err != nil {
		log.Panic(err)
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(models.SensorDataType)
	encoder.Encode(&models.SensorData{
		Time:      time.Now().Unix(),
		Heartbeat: i,
		Loc_x:     i,
		Loc_y:     i,
	})
	conn.Write(buffer.Bytes())
	conn.Close()
}

func SendVideoFragment(id int, port int) {
	address, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Panic(err)
	}

	conn, err := net.DialUDP("udp", nil, address)
	if err != nil {
		log.Panic(err)
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(models.VideoFragmentType)
	encoder.Encode(&models.VideoFragment{
		Time: time.Now().Unix(),
		Body: []byte("video fragment" + strconv.Itoa(id)),
	})
	conn.Write(buffer.Bytes())
	conn.Close()
}

func SendMessage(code int, port int) {
	address, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, address)
	if err != nil {
		log.Panic(err)
	}

	encoder := gob.NewEncoder(conn)
	encoder.Encode(models.MessageType)
	encoder.Encode(&models.Message{Time: time.Now().Unix(), Code: code})
	conn.Close()
}

func SendAudio(id int, port int) {
	address, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, address)
	if err != nil {
		log.Panic(err)
	}

	encoder := gob.NewEncoder(conn)
	encoder.Encode(models.AudioType)
	encoder.Encode(&models.Audio{
		Time: time.Now().Unix(),
		Body: []byte("audio" + strconv.Itoa(id)),
	})
	conn.Close()
}
