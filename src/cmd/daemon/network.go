package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/mido3ds/C4IAN/src/models"
)

const RetryTimeout = 2 * time.Second
const DialTimeout = 2 * time.Second

type onReceiveMsgCallback = func(models.Message)
type onReceiveAudioCallback = func(models.Audio)
type onReceiveVideoFragmentCallback = func(models.VideoFragment)
type onReceiveSensorsDataCallback = func(models.SensorData)

type NetworkManager struct {
	onReceiveMsg           onReceiveMsgCallback
	onReceiveAudio         onReceiveAudioCallback
	onReceiveVideoFragment onReceiveVideoFragmentCallback
	onReceiveSensorsData   onReceiveSensorsDataCallback
}

func NewNetworkManager(
	onReceiveMsg onReceiveMsgCallback,
	onReceiveAudio onReceiveAudioCallback,
	onReceiveVideoFragment onReceiveVideoFragmentCallback,
	onReceiveSensorsData onReceiveSensorsDataCallback,
) *NetworkManager {
	return &NetworkManager{
		onReceiveMsg:           onReceiveMsg,
		onReceiveAudio:         onReceiveAudio,
		onReceiveVideoFragment: onReceiveVideoFragment,
		onReceiveSensorsData:   onReceiveSensorsData,
	}
}

func (netManager *NetworkManager) Listen(port int) {
	go netManager.ListenTCP(port)
	go netManager.ListenUDP(port)
}

func (netManager *NetworkManager) SendTCP(dstAddrss string, dstPort int, payload interface{}) {
	// Connect to remote TCP socket
	conn, err := net.DialTimeout("tcp", dstAddrss+":"+strconv.Itoa(dstPort), DialTimeout)
	if err != nil {
		log.Println("Could not connect to unit: ", dstAddrss, " over TCP port: ", dstPort)
		log.Println("Retry in ", RetryTimeout)
		time.AfterFunc(RetryTimeout, func() { netManager.SendTCP(dstAddrss, dstPort, payload) })
		return
	}
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	// Get the payload type and send it
	if _, ok := payload.(models.Message); ok {
		encoder.Encode(models.MessageType)
	} else if _, ok := payload.(models.Audio); ok {
		encoder.Encode(models.AudioType)
	} else {
		log.Panic("Unknown payload type")
	}

	// Send the payload
	err = encoder.Encode(payload)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Sent TCP packet to: ", dstAddrss)
}

func (netManager *NetworkManager) SendGroupsHello(groupMembers map[string][]string, port int) {
	for {
		time.Sleep(2 * time.Second)
		for group := range groupMembers {
			netManager.SendUDP(group, port, models.Message{Code: -1})
		}
	}
}

func (netManager *NetworkManager) SendUDP(dstAddrss string, dstPort int, payload interface{}) {
	address, err := net.ResolveUDPAddr("udp", dstAddrss+":"+strconv.Itoa(dstPort))
	if err != nil {
		log.Panic(err)
	}

	conn, err := net.DialUDP("udp", nil, address)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	// Get the payload type and add it to the buffer
	if _, ok := payload.(models.Message); ok {
		encoder.Encode(models.MessageType)
	} else if _, ok := payload.(models.Audio); ok {
		encoder.Encode(models.AudioType)
	} else {
		log.Panic("Unknown payload type")
	}

	// Add the payload to the buffer
	err = encoder.Encode(payload)
	if err != nil {
		log.Panic(err)
	}

	// Send the buffer
	n, err := conn.Write(buffer.Bytes())
	if n != buffer.Len() {
		log.Panic("Could not write the whole buffer to UDP socket, buffer size: ",
			buffer.Len(), ", sent bytes: ", n)
	}
	if err != nil {
		log.Panic(err)
	}
	log.Println("Sent UDP (", n, "bytes) to: ", dstAddrss+":"+strconv.Itoa(dstPort))
}

func (netManager *NetworkManager) ListenTCP(port int) {
	// Get local TCP address
	address, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Panic(err)
	}

	// Listen for remote TCP connections
	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		log.Panic(err)
	}
	defer listener.Close()
	log.Println("listening on tcp port:", port)

	for {
		// Handle any incoming connections
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Could not accept TCP connection from: ", conn.RemoteAddr(), ", err: ", err)
			continue
		}
		log.Println("received tcp connection from unit, address:", conn.RemoteAddr()) // TODO: remove

		go netManager.handleTCPConnection(conn)
	}
}

func (netManager *NetworkManager) ListenUDP(port int) {
	// Get local UDP address
	address, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
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
		buffer := make([]byte, 64*1024*10)
		length, src, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Panic(err)
		}
		decoder := gob.NewDecoder(bytes.NewBuffer(buffer[:length]))
		// Decode any packets in the buffer by reading the type then the payload, then make appropriate callbacks
		for decoder.Decode(&packetType) == nil {
			if packetType == models.SensorDataType {
				var sensorsData models.SensorData
				err := decoder.Decode(&sensorsData)
				if err != nil {
					log.Panic(err)
				}
				sensorsData.Src = src.IP.String()
				log.Println("received sensor data:", sensorsData) // TODO: remove
				netManager.onReceiveSensorsData(sensorsData)
			} else {
				log.Panic("Unknow packet type received through UDP from: ", src.IP.String())
			}
		}
	}
}

func (netManager *NetworkManager) handleTCPConnection(conn net.Conn) {
	defer func() {
		log.Println("Closing TCP connection with ", conn.RemoteAddr())
		err := conn.Close()
		if err != nil {
			log.Println("Could not close connection with ", conn.RemoteAddr(), ", err: ", err)
		}
	}()

	// Extract the IP address of the source
	srcIP := strings.Split(conn.RemoteAddr().String(), ":")[0]

	// Decode the type of the packet
	decoder := gob.NewDecoder(conn)
	var packetType models.Type
	err := decoder.Decode(&packetType)
	if err != nil {
		log.Panic("failed to decode type, err:", err)
	}

	// Decode the payload of the packet and make appropriate callbacks
	switch packetType {
	case models.MessageType:
		var msg models.Message
		err := decoder.Decode(&msg)
		if err != nil {
			log.Panic("failed to decode code msg from unit, err:", err)
		}

		msg.Src = srcIP
		go netManager.onReceiveMsg(msg)

		log.Println("received code msg from unit:", msg) // TODO: remove
	case models.AudioType:
		var audio models.Audio
		err := decoder.Decode(&audio)
		if err != nil {
			log.Panic("failed to decode audio msg from unit, err:", err)
		}

		audio.Src = srcIP
		go netManager.onReceiveAudio(audio)

		log.Println("received audio msg from unit: len=", len(audio.Body)) // TODO: remove
	case models.VideoFragmentType:
		var videoFragment models.VideoFragment
		err := decoder.Decode(&videoFragment)
		if err != nil {
			log.Panic(err)
		}
		videoFragment.Src = srcIP
		log.Println(
			"received video fragment: filename=", videoFragment.FileName, " id=", videoFragment.ID,
		) // TODO: remove
		go netManager.onReceiveVideoFragment(videoFragment)
	default:
		log.Panic("Unknow packet type received through TCP from: ", srcIP)
	}
}
