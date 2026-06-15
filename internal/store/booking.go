package store

import (
	"database/sql"

	"github.com/erfan-goodarzi/booking-event-api/internals/models"
)

type BookingStore interface {
	CreateBooking(ticketId string, t *models.Booking) (*models.Booking, error)
	UpdateBookingStatus(id string, t *models.Booking) (*models.Booking, error)
}

type PostgresBookingStore struct {
	db *sql.DB
}

func NewPostgresBookingStore(db *sql.DB) *PostgresBookingStore {
	return &PostgresBookingStore{db: db}
}

func (pg *PostgresBookingStore) CreateBooking(ticketId string, b *models.Booking) (*models.Booking, error) {
	query := `INSERT INTO bookings(user_id, ticket_id, status)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`

	err := pg.db.QueryRow(query, b.UserId, ticketId, b.Status).Scan(
		&b.ID,
		&b.CreatedAt,
		&b.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func (pg *PostgresBookingStore) UpdateBookingStatus(id string, b *models.Booking) (*models.Booking, error) {
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
