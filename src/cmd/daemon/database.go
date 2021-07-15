package main

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/mido3ds/C4IAN/src/models"
)

type DatabaseManager struct {
	db *sqlx.DB
}

func NewDatabaseManager(dbPath string) *DatabaseManager {
	db := sqlx.MustOpen("sqlite3", dbPath)

	// Make sure foreign key constraints are enabled
	db.MustExec("PRAGMA foreign_keys = ON")

	// Create database from schema script
	_, err := sqlx.LoadFile(db, "schema.sql")
	if err != nil {
		log.Panic(err.Error())
	}

	return &DatabaseManager{db: db}
}

func (dbManager *DatabaseManager) Initialize(units []string, groups map[string][]string) {
	for _, unit := range units {
		dbManager.AddUnit(unit)
	}
	for group, members := range groups {
		dbManager.AddGroup(group, members)
	}
}

func (dbManager *DatabaseManager) AddUnit(IP string) {
	dbManager.db.MustExec("INSERT INTO units VALUES ($1)", IP)
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
		data.Time, data.Src, data.Heartbeat, data.Loc_x, data.Loc_y)
}

func (dbManager *DatabaseManager) AddReceivedVideo(src string, video *models.Video) {
	dbManager.db.MustExec("INSERT INTO received_videos VALUES ($1, $2, $3, $4)",
		video.Time, src, video.ID, video.Path)
}

func (dbManager *DatabaseManager) GetConversation(unitIP string) []models.Message {
	var msgs []models.Message
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
	var audios []models.Audio
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
	var data []models.SensorData
	err := dbManager.db.Select(
		&data,
		"SELECT time, heartbeat, loc_x, loc_y FROM received_sensors_data WHERE src = $1 ORDER BY time",
		srcIP,
	)
	if err != nil {
		log.Panic(err)
	}
	return data
}

func (dbManager *DatabaseManager) GetReceivedVideos(srcIP string) []models.Video {
	var videos []models.Video
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
