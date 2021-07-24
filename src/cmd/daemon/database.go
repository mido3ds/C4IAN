package main

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mido3ds/C4IAN/src/models"
)

var mutex sync.RWMutex

const InactiveThresholdInNanoseconds = 60 * time.Second

type DatabaseManager struct {
	db *sqlx.DB
}

func NewDatabaseManager(dbPath string) *DatabaseManager {
	db := sqlx.MustOpen("sqlite3", dbPath)

	// Make sure foreign key constraints are enabled
	// db.MustExec("PRAGMA foreign_keys = ON")

	// Create database from schema script
	db.MustExec(schemaSQL)

	return &DatabaseManager{db: db}
}

func (dbManager *DatabaseManager) Initialize(units []models.Unit, groupMembers map[string][]string) {
	for _, unit := range units {
		dbManager.AddUnit(unit)
	}
	for group, members := range groupMembers {
		dbManager.AddGroup(group, members)
	}
}

func (dbManager *DatabaseManager) AddUnit(unit models.Unit) {
	mutex.Lock()
	defer mutex.Unlock()
	dbManager.db.MustExec("INSERT INTO units VALUES ($1, $2, $3)", unit.IP, unit.Name, 0)
}

func (dbManager *DatabaseManager) AddGroup(groupIP string, memberIPs []string) {
	mutex.Lock()
	defer mutex.Unlock()
	dbManager.db.MustExec("INSERT INTO groups VALUES ($1)", groupIP)
	for _, memberIP := range memberIPs {
		dbManager.db.MustExec("INSERT INTO members VALUES ($1, $2)", groupIP, memberIP)
	}
}

func (dbManager *DatabaseManager) AddSentMessage(dstIP string, msg *models.Message) {
	mutex.Lock()
	defer mutex.Unlock()
	dbManager.db.MustExec("INSERT INTO sent_msgs VALUES ($1, $2, $3)",
		msg.Time, dstIP, msg.Code)
}

func (dbManager *DatabaseManager) AddSentAudio(dstIP string, audio *models.Audio) {
	mutex.Lock()
	defer mutex.Unlock()
	dbManager.db.MustExec("INSERT INTO sent_audios VALUES ($1, $2, $3)",
		audio.Time, dstIP, audio.Body)
}

func (dbManager *DatabaseManager) AddReceivedMessage(msg *models.Message) {
	mutex.Lock()
	defer mutex.Unlock()
	dbManager.db.MustExec("INSERT INTO received_msgs VALUES ($1, $2, $3)",
		msg.Time, msg.Src, msg.Code)
}

func (dbManager *DatabaseManager) AddReceivedAudio(audio *models.Audio) {
	mutex.Lock()
	defer mutex.Unlock()
	dbManager.db.MustExec("INSERT INTO received_audios VALUES ($1, $2, $3)",
		audio.Time, audio.Src, audio.Body)
}

func (dbManager *DatabaseManager) AddReceivedSensorsData(data *models.SensorData) {
	mutex.Lock()
	defer mutex.Unlock()
	dbManager.db.MustExec("INSERT INTO received_sensors_data VALUES ($1, $2, $3, $4, $5)",
		data.Time, data.Src, data.Heartbeat, data.Lat, data.Lon)
}

func (dbManager *DatabaseManager) UpdateLastActivity(ip string, lastActivity int64) {
	mutex.Lock()
	defer mutex.Unlock()
	dbManager.db.MustExec("UPDATE units SET last_activity = $1 WHERE ip = $2",
		lastActivity, ip)
}

func (dbManager *DatabaseManager) GetUnitsNames() map[string]string {
	mutex.Lock()
	defer mutex.Unlock()
	units := make([]models.Unit, 0)
	err := dbManager.db.Select(&units, "SELECT ip, name FROM units")
	if err != nil {
		log.Panic(err)
	}
	unitsNames := make(map[string]string)
	for _, unit := range units {
		unitsNames[unit.IP] = unit.Name
	}
	return unitsNames
}

func (dbManager *DatabaseManager) GetUnits() []models.Unit {
	mutex.Lock()
	defer mutex.Unlock()
	units := make([]models.Unit, 0)
	err := dbManager.db.Select(
		&units,
		`
		SELECT name, ip, $1 - last_activity < $2 as active,
		ifnull(heartbeat, 1000) as heartbeat,
		ifnull(lon, 1000) as lon,
		ifnull(lat, 1000) as lat
		FROM units LEFT JOIN
		(SELECT * FROM received_sensors_data)
		ON ip = src AND last_activity = time
		`,
		time.Now().UnixNano(),
		InactiveThresholdInNanoseconds,
	)
	if err != nil {
		log.Panic(err)
	}
	return units
}

func (dbManager *DatabaseManager) GetGroups() []models.Group {
	mutex.Lock()
	defer mutex.Unlock()
	groups := make([]models.Group, 0)
	err := dbManager.db.Select(
		&groups,
		"SELECT * FROM groups",
	)
	if err != nil {
		log.Panic(err)
	}
	return groups
}

func (dbManager *DatabaseManager) GetMemberships() []models.Membership {
	mutex.Lock()
	defer mutex.Unlock()
	memberships := make([]models.Membership, 0)
	err := dbManager.db.Select(
		&memberships,
		"SELECT * FROM members",
	)
	if err != nil {
		log.Panic(err)
	}
	return memberships
}

func (dbManager *DatabaseManager) GetConversation(unitIP string) []models.Message {
	mutex.Lock()
	defer mutex.Unlock()
	msgs := make([]models.Message, 0)
	err := dbManager.db.Select(
		&msgs,
		`
		SELECT time, code, 0 as sent FROM received_msgs WHERE src = $1
			UNION
		SELECT time, code, 1 as sent FROM sent_msgs WHERE dst = $1
		ORDER BY time
		`,
		unitIP,
	)
	if err != nil {
		log.Panic(err)
	}
	return msgs
}

func (dbManager *DatabaseManager) GetAllReceivedMessages() []models.Message {
	mutex.Lock()
	defer mutex.Unlock()
	msgs := make([]models.Message, 0)
	err := dbManager.db.Select(&msgs, "SELECT *, 0 as sent FROM received_msgs ORDER BY time")
	if err != nil {
		log.Panic(err)
	}
	return msgs
}

func (dbManager *DatabaseManager) GetReceivedAudio(srcIP string) []models.Audio {
	mutex.Lock()
	defer mutex.Unlock()
	audios := make([]models.Audio, 0)
	err := dbManager.db.Select(
		&audios,
		"SELECT time, body FROM received_audios WHERE src = $1 ORDER BY time",
		srcIP,
	)
	if err != nil {
		log.Panic(err)
	}
	return audios
}

func (dbManager *DatabaseManager) GetReceivedSensorsData(srcIP string) []models.SensorData {
	mutex.Lock()
	defer mutex.Unlock()
	data := make([]models.SensorData, 0)
	err := dbManager.db.Select(
		&data,
		"SELECT time, heartbeat, lat, lon FROM received_sensors_data WHERE src = $1 ORDER BY time",
		srcIP,
	)
	if err != nil {
		log.Panic(err)
	}
	return data
}

func (dbManager *DatabaseManager) GetReceivedVideos(srcIP string) []models.Video {
	mutex.Lock()
	defer mutex.Unlock()
	videos := make([]models.Video, 0)
	err := dbManager.db.Select(
		&videos,
		"SELECT * FROM received_videos WHERE src = $1 ORDER BY time",
		srcIP,
	)
	if err != nil {
		log.Panic(err)
	}
	return videos
}

func (dbManager *DatabaseManager) AddVideoIfNew(frag *models.VideoFragment) bool {
	mutex.Lock()
	defer mutex.Unlock()
	// Check if the video already exists in the database
	exists := true
	row := dbManager.db.QueryRowx(
		"SELECT * FROM received_videos WHERE src = $1 AND id = $2",
		frag.Src, frag.ID,
	)

	var video models.Video
	err := row.StructScan(&video)
	if err == sql.ErrNoRows {
		exists = false
	} else if err != nil {
		log.Panic(err)
	}

	// Add video if it does not exist
	if !exists {
		dbManager.db.MustExec("INSERT INTO received_videos VALUES ($1, $2, $3)",
			frag.Time, frag.Src, frag.ID)
	}
	return exists
}
