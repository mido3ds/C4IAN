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
	decoder             *json.Decoder
}

func NewZoneIDAgent(locSocket string, zlen byte) (*ZoneIDAgent, error) {
	myZoneMutex.Lock()
	myZone.ID = 0
	myZone.Len = zlen
	myZoneMutex.Unlock()

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

	log.Println("initailized ZoneIDAgent, sock=", locSocket)

	zidAgent := ZoneIDAgent{
		conn:                l,
		zlen:                zlen,
		locSocket:           locSocket,
		zoneChangeCallbacks: make([]func(ZoneID), 0),
		decoder:             d,
	}

	// Make sure ZID in initialised correctly
	zidAgent.updateZone()

	return &zidAgent, nil
}

func (a *ZoneIDAgent) Start() {
	log.Println("started ZoneIDAgent")

	for {
		a.updateZone()
	}
}

func (a *ZoneIDAgent) updateZone() {
	var loc gpsLocation
	err := a.decoder.Decode(&loc)
	if err != nil {
		log.Println("err in loc decoding")
		return
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

func (a *ZoneIDAgent) Close() {
	a.conn.Close()
	os.RemoveAll(a.locSocket)
}

func (a *ZoneIDAgent) AddZoneChangeCallback(cb func(ZoneID)) {
	a.zoneChangeCallbacks = append(a.zoneChangeCallbacks, cb)
}
