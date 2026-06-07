package store

import (
	"database/sql"
	"time"
)

type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	DateTime    time.Time `json:"dateTime" db:"date_time"`
	UserID      string    `json:"userId" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateEventRequest struct {
	Title       string    `json:"title" validate:"required,min=3,max=50"`
	Description string    `json:"description"`
	Location    string    `json:"location" validate:"required"`
	DateTime    time.Time `json:"dateTime" validate:"required"`
}

type PatchEventRequest struct {
	Title       *string    `json:"title" validate:"omitempty,min=3,max=50"`
	Description *string    `json:"description" validate:"omitempty"`
	Location    *string    `json:"location" validate:"omitempty"`
	DateTime    *time.Time `json:"dateTime" validate:"omitempty"`
}

type EventStore interface {
	GetAllEvents() ([]Event, error)
	GetEvent(id string) (*Event, error)
	CreateEvent(*Event) (*Event, error)
	UpdateEvent(*Event) (*Event, error)
	DeleteEvent(id string) error
	GetEventOwner(eventId string) (string, error)
}

type PostgresEventStore struct {
	db *sql.DB
}

func NewPostgresEventStore(db *sql.DB) *PostgresEventStore {
	return &PostgresEventStore{db: db}
}

func (pg *PostgresEventStore) GetAllEvents() ([]Event, error) {
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

func (pg *PostgresEventStore) GetEvent(id string) (*Event, error) {
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

func (pg *PostgresEventStore) CreateEvent(e *Event) (*Event, error) {
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

func (pg *PostgresEventStore) UpdateEvent(e *Event) (*Event, error) {
	tx, err := pg.db.Begin()

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	UPDATE events
	SET title = $1, description = $2, location = $3, updated_at = NOW()
	WHERE id = $4
	RETURNING updated_at
	`

	err = tx.QueryRow(query, e.Title, e.Description, e.Location, e.ID).Scan(&e.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return e, nil
}

func ApplyEventPatch(e *Event, p PatchEventRequest) error {
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
	return nil
}

func (pg *PostgresEventStore) DeleteEvent(id string) error {
	query := `DELETE FROM events WHERE id = $1`

	stmt, err := pg.db.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.Exec(id)

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (pg *PostgresEventStore) GetEventOwner(eventId string) (string, error) {
	var userID string

	query := `
  SELECT user_id
  FROM events
  WHERE id = $1
  `

	err := pg.db.QueryRow(query, eventId).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}
