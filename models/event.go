package models

import (
	"time"

	"example.com/booking-event/db"
)

type Event struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	DateTime    time.Time `json:"dateTime" db:"date_time"`
	UserID      int       `json:"userId" db:"user_id"`
}

func GetAllEvents() ([]Event, error) {
	var events []Event

	query := `
	SELECT
		id,
		title,
		description,
		location,
		date_time,
		user_id
	FROM events
	`

	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var event Event

		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Location,
			&event.DateTime,
			&event.UserID,
		)

		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func GetEvent(id int64) (*Event, error) {
	var event Event

	query := `
	SELECT
		id,
		title,
		description,
		location,
		date_time,
		user_id
	FROM events
	WHERE id = ?
	`

	row := db.DB.QueryRow(query, id)

	err := row.Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Location,
		&event.DateTime,
		&event.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (e *Event) Create() error {
	query := `
	INSERT INTO events(title, description, location, date_time, user_id)
	VALUES (?,?,?,?,?)`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(e.Title, e.Description, e.Location, e.DateTime, e.UserID)

	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	e.ID = id
	return err
}
