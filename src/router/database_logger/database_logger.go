package database_logger

import (
	"net"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

type DatabaseManager struct {
	myIP net.IP
	db   *sqlx.DB
}

const path = "/var/log/caian/log.db"

var DatabaseLogger DatabaseManager

func (dbManager *DatabaseManager) Initialize(myIP net.IP) {
	dbManager.db = sqlx.MustOpen("sqlite3", path)
	dbManager.myIP = myIP
}

func (dbManager *DatabaseManager) LogForwarding(payload []byte, dst net.IP) {
	location := MyLocation()
	hash := HashSHA3(payload)
	dbManager.db.MustExec(
		"INSERT INTO forwarding VALUES ($1, $2, $3, $4, $5, $6)",
		time.Now().UnixNano(), dbManager.myIP.String(), dst.String(), hash, location.Lat, location.Lon,
	)
}
