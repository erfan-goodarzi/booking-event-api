package store

import (
	"database/sql"
	"errors"

	"github.com/erfan-goodarzi/booking-event-api/utils"
)

type User struct {
	ID       int64
	Username string `binding:"omitempty"`
	Email    string `binding:"required"`
	Password string `binding:"required"`
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

	hashesPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	err = tx.QueryRow(query, u.Email, hashesPassword, u.Username).Scan(&u.ID)

	if err != nil {
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
		return errors.New("Invalid Credential")
	}

	isValidPassword := utils.CheckPassword(password, u.Password)

	if !isValidPassword {
		return errors.New("Invalid Credential")
	}
	return nil
}
