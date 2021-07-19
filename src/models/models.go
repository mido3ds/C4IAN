package models

const (
	MESSAGE_EVENT      = "msg"
	AUDIO_EVENT        = "audio"
	VIDEO_EVENT        = "video"
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

type Unit struct {
	Name      string  `json:"name" db:"name"`
	IP        string  `json:"ip" db:"ip"`
	Active    bool    `json:"active" db:"active"`
	Heartbeat int     `json:"heartbeat" db:"heartbeat"`
	Lat       float64 `json:"lat" db:"lat"`
	Lon       float64 `json:"lon" db:"lon"`
}

type Group struct {
	IP string `json:"ip" db:"ip"`
}

type Membership struct {
	GroupIP string `json:"groupIP" db:"group_ip"`
	UnitIP  string `json:"unitIP" db:"unit_ip"`
}

// Time as int64 to store in db as Unix Time (SQLite3 does not support a time type)
type Message struct {
	Src       string `json:"src" db:"src"`
	SentByCMD bool   `json:"sent" db:"sent"` // Only between CMD Daemon & its UI
	Time      int64  `json:"time" db:"time"`
	Code      int    `json:"code" db:"code"`
}

type Audio struct {
	Src  string  `json:"src"`
	Time int64   `json:"time" db:"time"`
	Body []uint8 `json:"body" db:"body"`
}

// Only in SSEs between CMD Daemon & its UI
type VideoFragment struct {
	Src      string `json:"src"`
	ID       int    `json:"id"`
	FileName string `json:"file"`
	Time     int64  `json:"time"`
	Metadata []byte `json:"metadata"`
	Body     []byte `json:"body"`
}

type Video struct {
	Src  string `json:"src" db:"src"`
	ID   int    `json:"id" db:"id"`
	Time int64  `json:"time" db:"time"`
}

type SensorData struct {
	Src       string  `json:"src" db:"src"`
	Time      int64   `json:"time" db:"time"`
	Heartbeat int     `json:"heartbeat" db:"heartbeat"`
	Lat       float64 `json:"lat" db:"lat"`
	Lon       float64 `json:"lon" db:"lon"`
}

func (msg *Message) EventType() string {
	return MESSAGE_EVENT
}

func (audio *Audio) EventType() string {
	return AUDIO_EVENT
}

func (video *Video) EventType() string {
	return VIDEO_EVENT
}

func (frame *VideoFragment) EventType() string {
	return VIDEO_FRAME_EVENT
}

func (sensors *SensorData) EventType() string {
	return SENSORS_DATA_EVENT
}
