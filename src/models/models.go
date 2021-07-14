package models

// Time as integer to store in db as Unix Time (SQLite3 does not support a time type)
type Message struct {
	Time int `json:"time" db:"time"`
	Code int `json:"code" db:"code"`
}

type Audio struct {
	Time int    `json:"time" db:"time"`
	Body []byte `json:"body" db:"body"`
}

type VideoFrame struct {
	Time int    `json:"time" db:"time"`
	Body []byte `json:"body" db:"body"`
}

type Video struct {
	Time int    `json:"time" db:"time"`
	Path string `json:"path" db:"path"`
}

type SensorData struct {
	Time      int `json:"time" db:"time"`
	Heartbeat int `json:"heartbeat" db:"heartbeat"`
	Loc_x     int `json:"loc_x" db:"loc_x"`
	Loc_y     int `json:"loc_y" db:"loc_y"`
}
