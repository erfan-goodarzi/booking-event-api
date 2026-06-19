package store

import (
	"database/sql"
	"errors"

	"github.com/erfan-goodarzi/booking-event-api/internal/db"
	"github.com/erfan-goodarzi/booking-event-api/internal/models"
	"github.com/jackc/pgconn"
)

type BookingStore interface {
	Create(ticketId string, t *models.Booking) (*models.Booking, error)
	UpdateStatus(id string, t *models.Booking) (*models.Booking, error)
}

type PostgresBookingStore struct {
	db db.DB
}

func NewPostgresBookingStore(db db.DB) *PostgresBookingStore {
	return &PostgresBookingStore{db: db}
}

func (pg *PostgresBookingStore) Create(ticketId string, b *models.Booking) (*models.Booking, error) {
	tx, err := pg.db.Begin()

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var remaining int

	err = tx.QueryRow(`
    UPDATE tickets
    SET quantity = quantity - 1
    WHERE id = $1
      AND quantity > 0
    RETURNING quantity
`, ticketId).Scan(&remaining)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("TICKET_SOLD_OUT")
		}
		return nil, err
	}

	query := `INSERT INTO bookings(user_id, ticket_id, status)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(query, b.UserId, ticketId, b.Status).Scan(
		&b.ID,
		&b.CreatedAt,
		&b.UpdatedAt,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return nil, errors.New("USER_ALREADY_REGISTERED_FOR_TICKET")
			}
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (pg *PostgresBookingStore) UpdateStatus(id string, b *models.Booking) (*models.Booking, error) {
	query := `UPDATE bookings
	SET status = $1
	WHERE id = $2
	RETURNING updated_at, id, user_id, ticket_id
	`

	err := pg.db.QueryRow(query, b.Status, id).Scan(
		&b.UpdatedAt,
		&b.ID,
		&b.UserId,
		&b.TicketId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return b, nil
}

func ApplyRegistrationPatch(b *models.Booking, p models.PatchBookingRequest) error {
	if p.Status != nil {
		b.Status = *p.Status
	}

	return nil
}
