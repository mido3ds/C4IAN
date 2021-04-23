package zid

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

type ZoneIDController struct {
	conn      *net.UnixConn
	zlen      byte
	listeners []func(ZoneID)

	// don't read it, carries garbage at beginning, use AddListener to be notifies when it gets a correct value
	lastZoneID ZoneID
}

func NewZoneIDController(locSocket string, zlen byte) (*ZoneIDController, error) {
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

	log.Println("initailized ZoneIDController, sock=", locSocket)

	return &ZoneIDController{
		conn:      l,
		zlen:      zlen,
		listeners: make([]func(ZoneID), 0),
	}, nil
}

func (a *ZoneIDController) Start() {
	log.Println("started ZoneIDController")

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
func (a *ZoneIDController) AddListener(f func(ZoneID)) {
	a.listeners = append(a.listeners, f)
}
