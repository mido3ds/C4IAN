package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/mido3ds/C4IAN/src/models"
)

type DatabaseManager struct {
	db *sqlx.DB
}

func NewDatabaseManager(dbPath string) *DatabaseManager {
	db := sqlx.MustOpen("sqlite3", dbPath)

	// TODO: load any necessary config to the database (e.g. units ips)
	_, err := sqlx.LoadFile(db, "schema.sql")
	if err != nil {
		log.Panic(err.Error())
	}
	return &DatabaseManager{db: db}
}

func (dbManager *DatabaseManager) addUnit(IP string) {
	dbManager.db.MustExec("INSERT INTO units VALUES ($1)", IP)
}

func (dbManager *DatabaseManager) addGroup(groupIP string, memberIPs []string) {
	dbManager.db.MustExec("INSERT INTO groups VALUES ($1)", groupIP)
	for _, memberIP := range memberIPs {
		dbManager.db.MustExec("INSERT INTO members VALUES ($1, $2)", groupIP, memberIP)
	}
}

func (dbManager *DatabaseManager) addSentMessage(dstIP string, msg *models.Message) {
	dbManager.db.MustExec("INSERT INTO sent_msgs VALUES ($1, $2, $3)",
		msg.Time, dstIP, msg.Code)
}

func (dbManager *DatabaseManager) addSentAudio(dstIP string, audio *models.Audio) {
	// TODO: check if body is inserted correctly
	dbManager.db.MustExec("INSERT INTO sent_audios VALUES ($1, $2, $3)",
		audio.Time, dstIP, audio.Body)
}

func (dbManager *DatabaseManager) addReceivedMessage(srcIP string, msg *models.Message) {
	dbManager.db.MustExec("INSERT INTO received_msgs VALUES ($1, $2, $3)",
		msg.Time, srcIP, msg.Code)
}

func (dbManager *DatabaseManager) addReceivedAudio(srcIP string, audio *models.Audio) {
	dbManager.db.MustExec("INSERT INTO received_audios VALUES ($1, $2, $3)",
		audio.Time, srcIP, audio.Body)
}

func (dbManager *DatabaseManager) addReceivedSensorsData(srcIP string, data *models.SensorData) {
	dbManager.db.MustExec("INSERT INTO received_sensors_data VALUES ($1, $2, $3, $4, $5)",
		data.Time, srcIP, data.Heartbeat, data.Loc_x, data.Loc_y)
}

func (dbManager *DatabaseManager) getReceivedMessages(srcIP string) []models.Message {
	var msgs []models.Message
	err := dbManager.db.Select(
		&msgs,
		"SELECT time, code FROM received_msgs WHERE src = $1",
		srcIP,
	)
	if err != nil {
		log.Panic(err)
	}
	return msgs
}

func (dbManager *DatabaseManager) getReceivedAudio(srcIP string) []models.Audio {
	var audios []models.Audio
	err := dbManager.db.Select(
		&audios,
		"SELECT time, body FROM received_audios WHERE src = $1",
		srcIP,
	)
	if err != nil {
		log.Panic(err)
	}
	return audios
}

func (dbManager *DatabaseManager) getReceivedSensorsData(srcIP string) []models.SensorData {
	var data []models.SensorData
	err := dbManager.db.Select(
		&data,
		"SELECT time, heartbeat, loc_x, loc_y FROM received_sensors_data WHERE src = $1",
		srcIP,
	)
	if err != nil {
		log.Panic(err)
	}
	return data
}
