package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"database/sql"

	"github.com/akamensky/argparse"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mido3ds/C4IAN/src/unit/daemon/halapi"
)

func main() {
	defer log.Println("finished cleaning up, closing")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetOutput(os.Stdout)

	args, err := parseArgs()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Println(args)

	context := newContext(args.StorePath)
	defer context.close()

	go context.listenHAL(args.HALSocketPath)
	go context.listenCmdTcp(args.Port)

	waitSIGINT()
}

func fileExists(path string) bool {
	_, err := os.Open(path)
	return err == nil
}

func openDB(storePath string) *sql.DB {
	if fileExists(storePath) {
		log.Println("db exists, won't recreate it")
	} else {
		file, err := os.Create(storePath)
		if err != nil {
			log.Panic(err.Error())
		}
		file.Close()

		log.Println("opened file for db")
	}

	sqliteDatabase, err := sql.Open("sqlite3", storePath)
	if err != nil {
		log.Panic(err)
	}
	return sqliteDatabase
}

type Context struct {
	expectingVideoStream bool
	storeDB              *sql.DB
}

func newContext(storePath string) Context {
	context := Context{
		expectingVideoStream: false,
	}

	if len(storePath) > 0 {
		context.storeDB = openDB(storePath)
		context.createTables()
	}

	return context
}

func (c *Context) close() {
	if c.storeDB != nil {
		c.storeDB.Close()
	}
}

func (c *Context) createTables() {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS HeartBeats (beatsPerMinute INT, time INT);`,
		`CREATE TABLE IF NOT EXISTS Locations (lon REAL, lat REAL, time INT);`,
		`CREATE TABLE IF NOT EXISTS VideoFragments (data BLOB, time INT);`,
	}

	for _, v := range tables {
		statement, err := c.storeDB.Prepare(v)
		if err != nil {
			log.Panic(err)
		}
		defer statement.Close()

		_, err = statement.Exec()
		if err != nil {
			log.Panic(err)
		}
		log.Println(v)
	}

	log.Println("created all tables")
}

func (c *Context) saveHeartbeat(beatsPerMinute int) error {
	if c.storeDB == nil {
		return nil
	}

	statement, err := c.storeDB.Prepare(`INSERT INTO HeartBeats (beatsPerMinute, time) VALUES(?, strftime('%s','now'));`)
	if err != nil {
		return fmt.Errorf("couldn't insert beatsPerMinute, err: %v", err)
	}
	defer statement.Close()

	_, err = statement.Exec(beatsPerMinute)
	if err != nil {
		return fmt.Errorf("couldn't insert beatsPerMinute, err: %v", err)
	}

	return nil
}

func (c *Context) saveLocation(lon, lat float64) error {
	if c.storeDB == nil {
		return nil
	}

	statement, err := c.storeDB.Prepare(`INSERT INTO Locations (lon, lat, time) VALUES(?, ?, strftime('%s','now'));`)
	if err != nil {
		return fmt.Errorf("couldn't insert location, err: %v", err)
	}
	defer statement.Close()

	_, err = statement.Exec(lon, lat)
	if err != nil {
		return fmt.Errorf("couldn't insert location, err: %v", err)
	}

	return nil
}

// TODO: append video fragments to one row
func (c *Context) saveVideoFragment(data []byte) error {
	if c.storeDB == nil {
		return nil
	}

	statement, err := c.storeDB.Prepare(`INSERT INTO VideoFragments (data, time) VALUES(?, strftime('%s','now'));`)
	if err != nil {
		return fmt.Errorf("couldn't insert heartbeat, err: %v", err)
	}
	defer statement.Close()

	_, err = statement.Exec(data)
	if err != nil {
		return fmt.Errorf("couldn't insert heartbeat, err: %v", err)
	}

	return nil
}

func (c *Context) listenCmdTcp(port int) {
	cmdLisener, err := net.Listen("tcp", fmt.Sprint("0.0.0.0:", port))
	if err != nil {
		log.Panic(err)
	}
	defer cmdLisener.Close()

	for {
		conn, err := cmdLisener.Accept()
		if err != nil {
			log.Println("failed to accept, err:", err)
		}
		defer conn.Close()

		c.serveCmdTcp(conn)
	}
}

func (c *Context) listenCmdUdp(port int) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()
	c.serveCmdUdp(conn)
}

func (c *Context) serveCmdTcp(conn net.Conn) {
	// TODO: receive msgs
}

func (c *Context) serveCmdUdp(conn *net.UDPConn) {
	// TODO: receive msgs/streams/reqeuestToStream
}

// TODO: remove
func simulateHALClient(HALSocketPath string) {
	conn, err := net.Dial("unix", HALSocketPath)
	if err != nil {
		log.Panic(err)
	}

	enc := gob.NewEncoder(conn)
	halapi.Location{Lon: 5, Lat: 2}.Send(enc)
	halapi.HeartBeat{BeatsPerMinut: 5}.Send(enc)
	halapi.VideoFragment{Video: []byte{5, 3}}.Send(enc)
}

func (c *Context) listenHAL(HALSocketPath string) {
	// remove loc socket file
	err := os.RemoveAll(HALSocketPath)
	if err != nil {
		log.Panic("failed to remove socket:", err)
	}

	halListener, err := net.Listen("unix", HALSocketPath)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer halListener.Close()

	for {
		go simulateHALClient(HALSocketPath)
		conn, err := halListener.Accept()
		defer conn.Close()

		if err != nil {
			log.Println("accept error:", err)
		} else {
			log.Println("HAL connected")
			c.serveHAL(conn)
		}
	}
}

func (context *Context) serveHAL(conn net.Conn) {
	dec := gob.NewDecoder(conn)

	var video halapi.VideoFragment
	var heartbeat halapi.HeartBeat
	var loc halapi.Location

	for {
		sentType, err := halapi.RecvFromHAL(dec, &video, &heartbeat, &loc)
		if err != nil {
			log.Println(err)
		} else {
			switch sentType {
			case halapi.VideoFragmentType:
				context.onVideoReceived(&video)
				break
			case halapi.HeartBeatType:
				context.onHeartBeatReceived(&heartbeat)
				break
			case halapi.LocationType:
				context.onLocationReceived(&loc)
				break
			default:
				log.Panicf("received unkown msg type = %v", sentType)
				break
			}
		}
	}
}

func (c *Context) onVideoReceived(v *halapi.VideoFragment) {
	log.Printf("VideoFragment: %v\n", *v)
	c.saveVideoFragment(v.Video)

	if !c.expectingVideoStream {
		log.Println("error, not expecting video stream, but received packet from HAL")
		return
	}

	// TODO: send packet to cmd
}

func (c *Context) onHeartBeatReceived(hb *halapi.HeartBeat) {
	log.Printf("HeartBeat: %v\n", *hb)
	c.saveHeartbeat(hb.BeatsPerMinut)

	// TODO: send heart beat to cmd
}

func (c *Context) onLocationReceived(loc *halapi.Location) {
	log.Printf("Location: %v\n", *loc)
	c.saveLocation(loc.Lon, loc.Lat)

	// TODO: send location to cmd
}

// Args store command line arguments
type Args struct {
	// null if no store path
	StorePath     string
	Port          int
	HALSocketPath string
}

func (a *Args) String() string {
	return fmt.Sprintf("&Args{StorePath: %v, Port: %v, HALSocketPath: %v}", a.StorePath, a.Port, a.HALSocketPath)
}

func parseArgs() (*Args, error) {
	parser := argparse.NewParser("unit-daemon", "Unit client daemon")

	storePath := parser.String("s", "store", &argparse.Options{Help: "Path to archive video/positions/heartbeats. If not provided, won't store them.", Default: nil})
	port := parser.Int("p", "port", &argparse.Options{Help: "Main port the client will bind to, to receive connections from other clients.", Default: 4070})
	ctrlSocketPath := parser.String("", "hal-socket-path", &argparse.Options{Help: "Path to unix socket file to communicate over with HAL.", Default: "/tmp/unit.hal.sock"})

	err := parser.Parse(os.Args)
	if err != nil {
		return nil, errors.New(parser.Usage(err))
	}

	return &Args{
		StorePath:     *storePath,
		Port:          *port,
		HALSocketPath: *ctrlSocketPath,
	}, nil
}

func waitSIGINT() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
