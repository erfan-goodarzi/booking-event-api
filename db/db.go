package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	DB, err := sql.Open("sqlite3", "api.db")

	if err != nil {
		panic("Could not connect to DB.")
	}

	DB.SetMaxOpenConns(10)

	createTables()
}

func createTables() {
	createEventsTables := `
		CREATE TABLE IF NOT EXIST events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			location TEXT NOT NULL,
			dateTime DATETIME NOT NULL,
			userId INTEGER,
		)
	`

	_, err := DB.Exec(createEventsTables)

	if err != nil {
		panic("Could not create event table")
	}
}
