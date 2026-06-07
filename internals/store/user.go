package store

import (
	"database/sql"
	"errors"

	"github.com/erfan-goodarzi/booking-event-api/apiUtils"
	"github.com/jackc/pgconn"
)

type User struct {
	ID       string
	Username string `binding:"omitempty" json:"username" validate:"omitempty,min=3,max=20"`
	Email    string `binding:"required" json:"email" validate:"required,email"`
	Password string `binding:"required" json:"password" validate:"required,min=8"`
}

type UserStore interface {
	Create(u *User) error
	ValidateCredential(u *User) error
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

func (pg *PostgresUserStore) Create(u *User) error {
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

func (pg *PostgresUserStore) ValidateCredential(u *User) error {
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
