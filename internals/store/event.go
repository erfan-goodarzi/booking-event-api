package store

import (
	"database/sql"

	"github.com/erfan-goodarzi/booking-event-api/internals/models"
)

type EventStore interface {
	GetAllEvents() ([]models.Event, error)
	GetEvent(id string) (*models.Event, error)
	CreateEvent(*models.Event) (*models.Event, error)
	UpdateEvent(*models.Event) (*models.Event, error)
	DeleteEvent(id string) error
	GetEventOwner(eventId string) (string, error)
}

type PostgresEventStore struct {
	db *sql.DB
}

func NewPostgresEventStore(db *sql.DB) *PostgresEventStore {
	return &PostgresEventStore{db: db}
}

func (pg *PostgresEventStore) GetAllEvents() ([]models.Event, error) {
	var events []models.Event

	eventsQuery := `
		SELECT
    e.id,
    e.title,
    e.description,
    e.location,
    e.date_time,
    e.user_id,
    e.duration,
    e.created_at,
    e.updated_at,
    u.id as host_id,
    u.username as host_username,
    u.email as host_email
		FROM events e
		JOIN users u ON u.id = e.user_id
		ORDER BY e.created_at DESC
	`

	rows, err := pg.db.Query(eventsQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventsMap := make(map[string]*models.Event)
	var order []string

	for rows.Next() {
		var e models.Event

		err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.Description,
			&e.Location,
			&e.DateTime,
			&e.UserId,
			&e.Duration,
			&e.CreatedAt,
			&e.UpdatedAt,
			&e.Host.ID,
			&e.Host.Username,
			&e.Host.Email,
		)

		if err != nil {
			return nil, err
		}

		e.Tickets = make([]models.Ticket, 0)
		eventsMap[e.ID] = &e
		order = append(order, e.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	ticketsQuery := `
		SELECT
			id,
			user_id,
			event_id,
			type,
			price,
			quantity,
			created_at,
			updated_at
		FROM tickets
	`

	hRow, err := pg.db.Query(ticketsQuery)
	if err != nil {
		return nil, err
	}
	defer hRow.Close()

	for hRow.Next() {
		var t models.Ticket

		err := hRow.Scan(
			&t.ID, &t.UserId, &t.EventId, &t.Type, &t.Price, &t.Quantity, &t.CreatedAt, &t.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if event, exists := eventsMap[t.EventId]; exists {
			event.Tickets = append(event.Tickets, t)
		}
	}

	if err := hRow.Err(); err != nil {
		return nil, err
	}

	events = make([]models.Event, 0, len(order))
	for _, id := range order {
		events = append(events, *eventsMap[id])
	}

	return events, nil
}

func (pg *PostgresEventStore) GetEvent(id string) (*models.Event, error) {
	var event models.Event

	query := `
	SELECT
    e.id,
    e.title,
    e.description,
    e.location,
    e.date_time,
    e.user_id,
    e.duration,
    e.created_at,
    e.updated_at,
    u.id as host_id,
    u.username as host_username,
    u.email as host_email
		FROM events e
		JOIN users u ON u.id = e.user_id
		WHERE e.id = $1
	`

	row := pg.db.QueryRow(query, id)

	err := row.Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Location,
		&event.DateTime,
		&event.UserId,
		&event.Duration,
		&event.CreatedAt,
		&event.UpdatedAt,
		&event.Host.ID,
		&event.Host.Username,
		&event.Host.Email,
	)

	if err != nil {
		return nil, err
	}

	ticketsQuery := `
		SELECT
			id,
			user_id,
			event_id,
			type,
			price,
			quantity,
			created_at,
			updated_at
		FROM tickets
		WHERE event_id = $1
	`

	tRows, err := pg.db.Query(ticketsQuery, id)
	if err != nil {
		return nil, err
	}
	defer tRows.Close()

	for tRows.Next() {
		var t models.Ticket

		err := tRows.Scan(
			&t.ID, &t.UserId, &t.EventId, &t.Type, &t.Price, &t.Quantity, &t.CreatedAt, &t.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		event.Tickets = append(event.Tickets, t)
	}

	if err := tRows.Err(); err != nil {
		return nil, err
	}

	return &event, nil
}

func (pg *PostgresEventStore) CreateEvent(e *models.Event) (*models.Event, error) {
	tx, err := pg.db.Begin()

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO events(title, description, location, date_time, user_id, duration)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, updated_at`

	err = tx.QueryRow(query, e.Title, e.Description, e.Location, e.DateTime, e.UserId, e.Duration).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (pg *PostgresEventStore) UpdateEvent(e *models.Event) (*models.Event, error) {
	tx, err := pg.db.Begin()

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	UPDATE events
	SET title = $1, description = $2, location = $3, duration = $4, updated_at = NOW()
	WHERE id = $5
	RETURNING updated_at
	`

	err = tx.QueryRow(query, e.Title, e.Description, e.Location, e.Duration, e.ID).Scan(&e.UpdatedAt)

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

func ApplyEventPatch(e *models.Event, p models.PatchEventRequest) error {
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
	if p.Duration != nil {
		e.Duration = *p.Duration
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
