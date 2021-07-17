package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mido3ds/C4IAN/src/models"
)

const InactiveThresholdInSeconds = 60 * 2

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
	dbManager.db.MustExec("INSERT INTO units VALUES ($1, $2, $3)", unit.IP, unit.Name, 0)
}

func (dbManager *DatabaseManager) AddGroup(groupIP string, memberIPs []string) {
	dbManager.db.MustExec("INSERT INTO groups VALUES ($1)", groupIP)
	for _, memberIP := range memberIPs {
		dbManager.db.MustExec("INSERT INTO members VALUES ($1, $2)", groupIP, memberIP)
	}
}

func (dbManager *DatabaseManager) AddSentMessage(dstIP string, msg *models.Message) {
	dbManager.db.MustExec("INSERT INTO sent_msgs VALUES ($1, $2, $3)",
		msg.Time, dstIP, msg.Code)
}

func (dbManager *DatabaseManager) AddSentAudio(dstIP string, audio *models.Audio) {
	dbManager.db.MustExec("INSERT INTO sent_audios VALUES ($1, $2, $3)",
		audio.Time, dstIP, audio.Body)
}

func (dbManager *DatabaseManager) AddReceivedMessage(msg *models.Message) {
	dbManager.db.MustExec("INSERT INTO received_msgs VALUES ($1, $2, $3)",
		msg.Time, msg.Src, msg.Code)
}

func (dbManager *DatabaseManager) AddReceivedAudio(audio *models.Audio) {
	dbManager.db.MustExec("INSERT INTO received_audios VALUES ($1, $2, $3)",
		audio.Time, audio.Src, audio.Body)
}

func (dbManager *DatabaseManager) AddReceivedSensorsData(data *models.SensorData) {
	dbManager.db.MustExec("INSERT INTO received_sensors_data VALUES ($1, $2, $3, $4, $5)",
		data.Time, data.Src, data.Heartbeat, data.Lat, data.Lon)
}

func (dbManager *DatabaseManager) AddReceivedVideo(src string, video *models.Video) {
	dbManager.db.MustExec("INSERT INTO received_videos VALUES ($1, $2, $3, $4)",
		video.Time, src, video.ID, video.Path)
}

func (dbManager *DatabaseManager) UpdateLastActivity(ip string) {
	dbManager.db.MustExec("UPDATE units SET last_activity = $1 WHERE ip = $2",
		time.Now().Unix(), ip)
}

func (dbManager *DatabaseManager) GetUnits() []models.Unit {
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
		time.Now().Unix(),
		InactiveThresholdInSeconds,
	)
	if err != nil {
		log.Panic(err)
	}
	return units
}

func (dbManager *DatabaseManager) GetGroups() []models.Group {
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

func (dbManager *DatabaseManager) GetReceivedAudio(srcIP string) []models.Audio {
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
	videos := make([]models.Video, 0)
	err := dbManager.db.Select(
		&videos,
		"SELECT time, path, id FROM received_videos WHERE src = $1 ORDER BY time",
		srcIP,
	)
	if err != nil {
		log.Panic(err)
	}
	return videos
}

func (dbManager *DatabaseManager) GetReceivedVideo(srcIP string, id int) *models.Video {
	var video models.Video
	row := dbManager.db.QueryRowx(
		"SELECT time, path, id FROM received_videos WHERE src = $1 AND id = $2",
		srcIP, id,
	)
	err := row.StructScan(&video)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		log.Panic(err)
	}
	return &video
}
