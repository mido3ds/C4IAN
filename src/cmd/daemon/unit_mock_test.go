package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/mido3ds/C4IAN/src/models"
)

const (
	unitPort         = 4070
	cmdPort          = 4170
	multicastGroupIP = "224.0.0.251"
	unitIP           = "172.0.0.2"
	videoID          = 50
)

func TestUnit(*testing.T) {
	go ListenTCP(unitPort)
	go ListenUDPMulticast(multicastGroupIP, unitPort)

	for i := 0; true; i++ {
		//SendVideoFragment(i/10, i%10, cmdPort)
		SendMessage(i, cmdPort)
		time.Sleep(time.Second)
	}
}

func ListenUDPMulticast(groupIP string, port int) {
	// Get local UDP address
	address, err := net.ResolveUDPAddr("udp", groupIP+":"+strconv.Itoa(port))
	if err != nil {
		log.Panic(err)
	}

	// Listen for any remote UDP packets
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	var packetType models.Type
	for {
		// Read any incoming UDP packet
		buffer := make([]byte, 1024)
		length, err := conn.Read(buffer)
		if err != nil {
			log.Panic(err)
		}
		decoder := gob.NewDecoder(bytes.NewBuffer(buffer[:length]))
		// Decode any packets in the buffer by reading the type then the payload, then make appropriate callbacks
		for decoder.Decode(&packetType) == nil {
			if packetType == models.AudioType {
				var audio models.Audio
				err := decoder.Decode(&audio)
				if err != nil {
					log.Panic(err)
				}
				log.Println("Test: Received multicast audio:", audio)
			} else {
				log.Panic("Invalid packet type received through multicast")
			}
		}
	}
}

func ListenTCP(port int) {
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

				saveAudio(audio.Body)

				//log.Println("Test: Received an audio: ", audio)
			}
			conn.Close()
		}()
	}
}

func SendSensorsData(i float64, port int) {
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
		Src:       "10.0.0.3",
		Time:      time.Now().UnixNano(),
		Heartbeat: 10.0,
		Lat:       41.4568 + i*0.2,
		Lon:       -79.0512 + i*0.3,
	})
	conn.Write(buffer.Bytes())
	conn.Close()
}

func SendVideoFragment(id int, index int, port int) {
	address, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, address)
	if err != nil {
		log.Panic(err)
	}

	encoder := gob.NewEncoder(conn)
	encoder.Encode(models.VideoFragmentType)
	encoder.Encode(&models.VideoFragment{
		ID:       id,
		Time:     time.Now().UnixNano(),
		FileName: "frag" + strconv.Itoa(index) + ".ts",
		Body:     []byte("video fragment" + strconv.Itoa(index) + " "),
		Metadata: []byte("metadata sent with frag id: " + strconv.Itoa(id) + ", index: " + strconv.Itoa(index)),
	})
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
	encoder.Encode(&models.Message{Src: "10.0.0.3", Time: time.Now().UnixNano(), Code: code})
	conn.Close()
}

func SendAudio(i int, port int) {
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
		Time: time.Now().UnixNano(),
		Body: []byte("audio" + strconv.Itoa(i)),
	})
	conn.Close()
}

func saveAudio(audio []byte) {
	f, err := os.Create("a.mp3")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.Write(audio)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("done")
}
