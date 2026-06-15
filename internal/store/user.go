package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/erfan-goodarzi/booking-event-api/internals/models"
	"github.com/erfan-goodarzi/booking-event-api/pkg/apiUtils"
	"github.com/jackc/pgconn"
)

type UserStore interface {
	Create(u *models.User) error
	ValidateCredential(u *models.User) error
	SaveRefreshToken(userID string, token string, expiresAt time.Time) error
	DeleteRefreshToken(token string) error
	GetUserByRefreshToken(token string) (*models.User, error)
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

func (pg *PostgresUserStore) Create(u *models.User) error {
	tx, err := pg.db.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO users(email, password, username) 
	VALUES($1, $2, $3)
	RETURNING id
	`

	hashesPassword, err := apiUtils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	err = tx.QueryRow(query, u.Email, hashesPassword, u.Username).Scan(&u.ID)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				switch pgErr.ConstraintName {
				case "users_email_key":
					return errors.New("EMAIL_ALREADY_EXISTS")
				case "users_username_key":
					return errors.New("USERNAME_ALREADY_EXISTS")
				}
			}
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresUserStore) ValidateCredential(u *models.User) error {
	query := "SELECT id, password FROM users WHERE email = $1"

	row := pg.db.QueryRow(query, u.Email)

	var password string
	err := row.Scan(&u.ID, &password)

	if err != nil {
		return errors.New("INVALID_CREDENTIAL")
	}

	isValidPassword := apiUtils.CheckPassword(password, u.Password)

	if !isValidPassword {
		return errors.New("INVALID_CREDENTIAL")
	}
	return nil
}

func (pg *PostgresUserStore) SaveRefreshToken(userID string, token string, expiresAt time.Time) error {
	query := `
	INSERT INTO refresh_tokens(user_id, token, expires_at) 
	VALUES($1, $2, $3)
	`

	_, err := pg.db.Exec(query, userID, token, expiresAt)

	return err
}

func (pg *PostgresUserStore) DeleteRefreshToken(token string) error {
	query := "DELETE FROM refresh_tokens WHERE token = $1"

	_, err := pg.db.Exec(query, token)

	return err
}

func (pg *PostgresUserStore) GetUserByRefreshToken(token string) (*models.User, error) {
	query := `
	SELECT u.id, u.email, u.username
	FROM users u
	INNER JOIN refresh_tokens rt ON u.id = rt.user_id
	WHERE rt.token = $1 AND rt.expires_at > NOW()
	`

	row := pg.db.QueryRow(query, token)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("INVALID_TOKEN")
		}
		return nil, err
	}

	return &user, nil
}
