package models

import (
	"database/sql"
	"time"
)

type Event struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	DateTime    time.Time `json:"dateTime" db:"date_time"`
	UserID      int64     `json:"userId" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type PatchEventRequest struct {
	Title       *string
	Description *string
	Location    *string
	DateTime    *time.Time
}

type EventStore interface {
	GetAllEvents() ([]Event, error)
	GetEvent(id int64) (*Event, error)
	CreateEvent(*Event) (*Event, error)
	UpdateEvent(*Event) (*Event, error)
	DeleteEvent(id int64) error
	GetEventOwner(id int64) (int, error)
	ApplyPatch(e *Event, p PatchEventRequest)
}

type PostgresEventModel struct {
	db *sql.DB
}

func NewPostgresEventStore(db *sql.DB) *PostgresEventModel {
	return &PostgresEventModel{db: db}
}

func (pg *PostgresEventModel) GetAllEvents() ([]Event, error) {
	var events []Event

	query := `
	SELECT
		id,
		title,
		description,
		location,
		date_time,
		user_id,
		created_at,
		updated_at
	FROM events
	`

	rows, err := pg.db.Query(query)

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
			&event.CreatedAt,
			&event.UpdatedAt,
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

func (pg *PostgresEventModel) GetEvent(id int64) (*Event, error) {
	var event Event

	query := `
	SELECT
		id,
		title,
		description,
		location,
		date_time,
		user_id,
		created_at,
		updated_at
		FROM events
		WHERE id = $1 
	`

	row := pg.db.QueryRow(query, id)

	err := row.Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Location,
		&event.DateTime,
		&event.UserID,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (pg *PostgresEventModel) CreateEvent(e *Event) (*Event, error) {
	tx, err := pg.db.Begin()

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO events(title, description, location, date_time, user_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at`

	err = tx.QueryRow(query, e.Title, e.Description, e.Location, e.DateTime, e.UserID).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (pg *PostgresEventModel) UpdateEvent(e *Event) (*Event, error) {
	tx, err := pg.db.Begin()

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	UPDATE events
	SET title = $1, description = $2, location = $3
	WHERE id = $4
	`

	stmt, err := tx.Prepare(query)

	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(e.Title, e.Description, e.Location, e.ID)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (pg *PostgresEventModel) ApplyPatch(e *Event, p PatchEventRequest) {
	if p.Title != nil {
		e.Title = *p.Title
	}
	if p.Description != nil {
		e.Description = *p.Description
	}
	if p.Location != nil {
		e.Location = *p.Location
	}
	if p.DateTime != nil {
		e.DateTime = *p.DateTime
	}
}

func (pg *PostgresEventModel) DeleteEvent(eventId int64) error {
	query := `DELETE FROM events WHERE id = $1`

	stmt, err := pg.db.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.Exec(eventId)

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (pg *PostgresEventModel) GetEventOwner(eventId int64) (int, error) {
	var userID int

	query := `
  SELECT user_id
  FROM events
  WHERE id = $1
  `

	err := pg.db.QueryRow(query, eventId).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
