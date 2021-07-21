package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

type Location struct {
	Lat float64
	Lon float64
}

type LocAgent struct {
	conn      *net.UnixConn
	decoder   *json.Decoder
	locSocket string
	location  Location
}

func newLocAgent(locSocket string) (*LocAgent, error) {
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

	d := json.NewDecoder(l)

	log.Println("initailized LocAgent, sock=", locSocket)

	locAgent := LocAgent{
		conn:      l,
		locSocket: locSocket,
		decoder:   d,
		location:  Location{Lon: 32.4, Lat: 43.098},
	}

	return &locAgent, nil
}

func (a *LocAgent) start() {
	log.Println("started LocAgent")

	for {
		err := a.decoder.Decode(&a.location)
		if err != nil {
			log.Println("err in loc decoding, err:", err)
			return
		}
	}
}

func (a *LocAgent) close() {
	a.conn.Close()
	os.RemoveAll(a.locSocket)
}
