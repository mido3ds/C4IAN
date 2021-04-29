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
	locSocket string
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

	myZlen = zlen

	return &ZoneIDAgent{
		conn:      l,
		zlen:      zlen,
		locSocket: locSocket,
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
		myZoneMutex.Lock()
		if id != myZone.ID {
			myZone.ID = id
			log.Println("New Zone =", myZone)
		}
		myZoneMutex.Unlock()
	}
}

func (a *ZoneIDAgent) Close() {
	a.conn.Close()
	os.RemoveAll(a.locSocket)
}
