package main

const port = 5000
const logDatabasePath = "/var/log/caian/log.db"

func main() {
	dbManager := NewDatabaseManager(logDatabasePath)
	api := NewAPI(dbManager)
	api.Start(port)
}
