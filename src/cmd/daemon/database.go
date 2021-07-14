package main

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type DatabaseManager struct {
	db *sqlx.DB
}

func NewDatabaseManager(dbPath string) *DatabaseManager {
	db := sqlx.MustOpen("sqlite3", dbPath)

	// TODO: load any necessary config to the database (e.g. units ips)
	_, err := sqlx.LoadFile(db, "schema.sql")
	if err != nil {
		log.Fatalln(err.Error())
	}
	return &DatabaseManager{db: db}
}
