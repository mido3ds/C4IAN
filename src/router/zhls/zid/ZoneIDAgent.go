package zid

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

type ZoneIDAgent struct {
	conn                *net.UnixConn
	zlen                byte
	locSocket           string
	zoneChangeCallbacks []func(ZoneID)
}

func NewZoneIDAgent(locSocket string, zlen byte) (*ZoneIDAgent, error) {
	myZone.Len = zlen

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

	myZoneMutex.Lock()
	myZone.Len = zlen
	myZoneMutex.Unlock()

	return &ZoneIDAgent{
		conn:                l,
		zlen:                zlen,
		locSocket:           locSocket,
		zoneChangeCallbacks: make([]func(ZoneID), 0),
	}, nil
}

func (a *ZoneIDAgent) Start() {
	log.Println("started ZoneIDAgent")

	d := json.NewDecoder(a.conn)

	for {
		var loc gpsLocation
		err := d.Decode(&loc)
		if err != nil {
			continue
		}

		id := newZoneID(loc, a.zlen)
		myZoneMutex.Lock()
		if id != myZone.ID {
			myZone.ID = id
			for _, cb := range a.zoneChangeCallbacks {
				cb(id)
			}
			log.Println("New Zone =", &myZone)
		}
		myZoneMutex.Unlock()
	}
}

func (a *ZoneIDAgent) Close() {
	a.conn.Close()
	os.RemoveAll(a.locSocket)
}

func (a *ZoneIDAgent) AddZoneChangeCallback(cb func(ZoneID)) {
	a.zoneChangeCallbacks = append(a.zoneChangeCallbacks, cb)
}
