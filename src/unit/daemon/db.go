package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

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

func (c *Context) createTables() {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS HeartBeats (beatsPerMinute INT, time INT);`,
		`CREATE TABLE IF NOT EXISTS Locations (lon REAL, lat REAL, time INT);`,
		`CREATE TABLE IF NOT EXISTS VideoFragments (data BLOB, metadata BLOB, filename TEXT, time INT);`,
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
func (c *Context) saveVideoFragment(data, metadata []byte, filename string) error {
	if c.storeDB == nil {
		return nil
	}

	statement, err := c.storeDB.Prepare(`INSERT INTO VideoFragments (data, metadata, filename, time) VALUES(?, ?, ?, strftime('%s','now'));`)
	if err != nil {
		return fmt.Errorf("couldn't insert heartbeat, err: %v", err)
	}
	defer statement.Close()

	_, err = statement.Exec(data, metadata, filename)
	if err != nil {
		return fmt.Errorf("couldn't insert heartbeat, err: %v", err)
	}

	return nil
}
