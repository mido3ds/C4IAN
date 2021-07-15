package models

const (
	MESSAGE_EVENT      = "msg"
	AUDIO_EVENT        = "audio"
	VIDEO_FRAME_EVENT  = "video-fragment"
	SENSORS_DATA_EVENT = "sensors-data"
)

type Type byte

const (
	MessageType Type = iota
	AudioType
	VideoFragmentType
	SensorDataType
)

type Event interface {
	EventType() string
}

// Time as int64 to store in db as Unix Time (SQLite3 does not support a time type)
type Message struct {
	Src       string `json:"src"`            // Only in SSEs between CMD Daemon & its UI
	SentByCMD bool   `json:"sent" db:"sent"` // Only between CMD Daemon & its UI
	Time      int64  `json:"time" db:"time"`
	Code      int    `json:"code" db:"code"`
}

type Audio struct {
	Src  string `json:"src"` // Only in SSEs between CMD Daemon & its UI
	Time int64  `json:"time" db:"time"`
	Body []byte `json:"body" db:"body"`
}

type VideoFragment struct {
	Src  string `json:"src"` // Only in SSEs between CMD Daemon & its UI
	ID   string `json:"id" db:"id"`
	Time int64  `json:"time" db:"time"`
	Body []byte `json:"body" db:"body"`
}

type Video struct {
	Time int64  `json:"time" db:"time"`
	Path string `json:"path" db:"path"`
}

type SensorData struct {
	Src       string `json:"src"` // Only in SSEs between CMD Daemon & its UI
	Time      int64  `json:"time" db:"time"`
	Heartbeat int    `json:"heartbeat" db:"heartbeat"`
	Loc_x     int    `json:"loc_x" db:"loc_x"`
	Loc_y     int    `json:"loc_y" db:"loc_y"`
}

func (msg *Message) EventType() string {
	return MESSAGE_EVENT
}

func (audio *Audio) EventType() string {
	return AUDIO_EVENT
}

func (frame *VideoFragment) EventType() string {
	return VIDEO_FRAME_EVENT
}

func (sensors *SensorData) EventType() string {
	return SENSORS_DATA_EVENT
}
