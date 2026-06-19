package store

import (
	"database/sql"

	"github.com/erfan-goodarzi/booking-event-api/internal/models"
)

type TicketStore interface {
	Create(id string, t *models.Ticket) (*models.Ticket, error)
}

type PostgresTicketStore struct {
	db *sql.DB
}

func NewPostgresTicketStore(db *sql.DB) *PostgresTicketStore {
	return &PostgresTicketStore{db: db}
}

func (pg *PostgresTicketStore) Create(id string, t *models.Ticket) (*models.Ticket, error) {
	query := `INSERT INTO tickets(user_id, event_id, type, price, quantity)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`

	err := pg.db.QueryRow(query, t.UserId, id, t.Type, t.Price, t.Quantity).Scan(
		&t.ID,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return t, nil
}
