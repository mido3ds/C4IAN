package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

type LocationAgent struct {
	conn     *net.UnixConn
	location Location
}

func NewLocationAgent(locSocket string) (*LocationAgent, error) {
	// remove loc socket file
	err := os.RemoveAll(locSocket)
	if err != nil {
		return nil, err
	}

	adr, err := net.ResolveUnixAddr("unixgram", locSocket)
	if err != nil {
		return nil, err
	}

	// create loc socket file
	l, err := net.ListenUnixgram("unixgram", adr)
	if err != nil {
		return nil, err
	}

	log.Println("initailized location agent, sock=", locSocket)

	return &LocationAgent{conn: l}, nil
}

func (a *LocationAgent) Start() {
	log.Println("started location agent")

	d := json.NewDecoder(a.conn)

	for {
		var loc Location
		err := d.Decode(&loc)
		if err != nil {
			continue
		}

		a.location = loc
		// TODO: calculate zone id
		// TODO: send signal of zone id change, if changed
	}
}

// Location is gps position
// where Lat and Lon are in degrees
type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
