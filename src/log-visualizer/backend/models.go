package main

type Packet struct {
	StartTime int64  `json:"start" db:"start_time"`
	Dst       string `json:"dst" db:"dst"`
	Hash      string `json:"hash" db:"packet_hash"`
}

type Location struct {
	Lat float64 `json:"lat" db:"lat"`
	Lon float64 `json:"lon" db:"lon"`
}

type PacketData struct {
	Locations map[string]Location `json:"locations"`
	Path      []string            `json:"path"`
}
