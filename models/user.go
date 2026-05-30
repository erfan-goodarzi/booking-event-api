package models

import (
	"errors"

	"github.com/erfan-goodarzi/booking-event-api/db"
	"github.com/erfan-goodarzi/booking-event-api/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Create() error {
	query := "INSERT INTO users(email, password) VALUES(?,?)"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	hashesPassword, err := utils.HashPassword(u.Password)

	if err != nil {
		return err
	}

	res, err := stmt.Exec(u.Email, hashesPassword)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	u.ID = id
	return err
}

func (u *User) ValidateCredential() error {
	query := "SELECT id, password FROM users WHERE email = ?"

	row := db.DB.QueryRow(query, u.Email)

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
