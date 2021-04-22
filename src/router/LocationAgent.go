package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

type LocationAgent struct {
	conn      *net.UnixConn
	zlen      byte
	listeners []func(ZoneID)

	// don't read it, carries garbage at beginning, use AddListener to be notifies when it gets a correct value
	lastZoneID ZoneID
}

func NewLocationAgent(locSocket string, zlen byte) (*LocationAgent, error) {
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

	return &LocationAgent{
		conn:      l,
		zlen:      zlen,
		listeners: make([]func(ZoneID), 0),
	}, nil
}

func (a *LocationAgent) Start() {
	log.Println("started location agent")

	d := json.NewDecoder(a.conn)

	for {
		var loc GPSLocation
		err := d.Decode(&loc)
		if err != nil {
			continue
		}

		id := NewZoneID(loc, a.zlen)
		if id != a.lastZoneID {
			a.lastZoneID = id

			go func() {
				// call listeners
				for _, f := range a.listeners {
					f(id)
				}
			}()
		}
	}
}

// AddListener appends a function to the list of functions
// that will be called when zoneid changes
func (a *LocationAgent) AddListener(f func(ZoneID)) {
	a.listeners = append(a.listeners, f)
}
