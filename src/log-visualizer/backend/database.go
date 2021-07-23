package main

import (
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseManager struct {
	db *sqlx.DB
}

func NewDatabaseManager(dbPath string) *DatabaseManager {
	db := sqlx.MustOpen("sqlite3", dbPath)
	return &DatabaseManager{db: db}
}

func (dbManager *DatabaseManager) GetPackets() []Packet {
	packets := make([]Packet, 0)
	err := dbManager.db.Select(
		&packets,
		"SELECT packet_hash, MIN(time) as start_time, dst FROM forwarding GROUP BY packet_hash",
	)
	if err != nil {
		log.Panic(err)
	}
	return packets
}

func (dbManager *DatabaseManager) GetPacketData(hash string) *PacketData {
	// Get path
	row := dbManager.db.QueryRowx(
		`
		SELECT GROUP_CONCAT(ip) as path, MIN(time) as start_time FROM 
		(SELECT MIN(time) as time, ip FROM forwarding 
		 WHERE packet_hash = $1 GROUP BY ip ORDER BY time)
		`,
		hash,
	)
	var path string
	var startTime int64
	err := row.Scan(&path, &startTime)
	if err != nil {
		log.Println(hash)
		log.Println(err)
		return nil
	}

	// Get most recent locations of nodes before forwarding the packet
	rows, err := dbManager.db.Queryx(
		`
		SELECT ip, lat, lon FROM locations LEFT JOIN
		(SELECT ip as unit_ip, MAX(time) as latest_time FROM locations WHERE time <= $1 GROUP BY ip)
		ON ip = unit_ip AND time = latest_time
		`,
		startTime,
	)
	if err != nil {
		log.Println(err)
		return nil
	}
	locations := make(map[string]Location)
	for rows.Next() {
		var ip string
		var location Location
		rows.Scan(&ip, &location.Lat, &location.Lon)
		locations[ip] = location
	}
	return &PacketData{
		Path:      strings.Split(path, ","),
		Locations: locations,
	}
}
