package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./api.db")

	if err != nil {
		panic("Could not connect to DB." + err.Error())
	}

	DB.SetMaxOpenConns(10)

	createTables()
}

func createTables() {
	createEventsTables := `
		CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
		  title TEXT NOT NULL,
			description TEXT NOT NULL,
			location TEXT NOT NULL,
			date_time DATETIME NOT NULL,
			user_id INTEGER
		)
	`

	_, err := DB.Exec(createEventsTables)

	if err != nil {
		panic("Could not create event table: " + err.Error())
	}
}
