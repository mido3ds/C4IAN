package zid

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

type ZoneIDAgent struct {
	conn      *net.UnixConn
	zlen      byte
	listeners []func(ZoneID)

	// don't read it, carries garbage at beginning, use AddListener to be notifies when it gets a correct value
	lastZoneID ZoneID
}

func NewZoneIDAgent(locSocket string, zlen byte) (*ZoneIDAgent, error) {
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

	log.Println("initailized ZoneIDAgent, sock=", locSocket)

	return &ZoneIDAgent{
		conn:      l,
		zlen:      zlen,
		listeners: make([]func(ZoneID), 0),
	}, nil
}

func (a *ZoneIDAgent) Start() {
	log.Println("started ZoneIDAgent")

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
func (a *ZoneIDAgent) AddListener(f func(ZoneID)) {
	a.listeners = append(a.listeners, f)
}
