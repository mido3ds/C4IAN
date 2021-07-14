package models

import "time"

type Message struct {
	Time time.Time `json:"time" db:"time"`
	Code int       `json:"code" db:"int"`
}

type Audio struct {
	Time time.Time `json:"time" db:"time"`
	Body []byte    `json:"body" db:"body"`
}

type VideoFrame struct {
	Time time.Time `json:"time" db:"time"`
	Body []byte    `json:"body" db:"body"`
}

type Video struct {
	Time time.Time `json:"time" db:"time"`
	Path string    `json:"path" db:"path"`
}

type SensorData struct {
	Time      time.Time `json:"time" db:"time"`
	Heartbeat int       `json:"heartbeat" db:"heartbeat"`
	Loc_x     int       `json:"loc_x" db:"loc_x"`
	Loc_y     int       `json:"loc_y" db:"loc_y"`
}
