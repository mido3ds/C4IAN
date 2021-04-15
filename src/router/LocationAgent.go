package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"sync"
)

type LocationAgent struct {
	locMutex sync.Mutex
	listener net.Listener
	location Location
}

func NewLocationAgent(locSocket string) (*LocationAgent, error) {
	// remove loc socket file
	err := os.RemoveAll(locSocket)
	if err != nil {
		return nil, err
	}

	// create loc socket file
	l, err := net.Listen("unix", locSocket)
	if err != nil {
		return nil, err
	}

	return &LocationAgent{listener: l}, nil
}

func (a *LocationAgent) Start() {
	for {
		conn, err := a.listener.Accept()
		if err != nil {
			log.Fatal("can't accept connections to unix socket, err: ", err)
		}
		log.Println("Location provider connected")

		d := json.NewDecoder(conn)

		for {
			var loc Location
			err = d.Decode(&loc)
			if err != nil {
				log.Println("Location provider disconnected")
				break
			}

			a.locMutex.Lock()
			defer a.locMutex.Unlock()
			a.location = loc
		}
	}
}

func (a *LocationAgent) Location() Location {
	a.locMutex.Lock()
	defer a.locMutex.Unlock()
	return a.location
}

// Location is gps position
type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Z   float64 `json:"z"`
}
