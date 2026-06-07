package db

import (
	"database/sql"
	"fmt"
	"io/fs"
	"time"

	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const maxOpenDBConn = 10
const maxDBIdleConn = 5
const maxDBLifeTime = 5 * time.Minute

func ConnectDB(distro string) (*sql.DB, error) {
	db, err := Open(distro)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(maxOpenDBConn)
	db.SetMaxIdleConns(maxDBIdleConn)
	db.SetConnMaxLifetime(maxDBLifeTime)

	return db, nil
}

func Open(constr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", constr)

	if err != nil {
		panic("Could not connect to DB." + err.Error())
	}

	if err := db.Ping(); err != nil {
		fmt.Printf("db: open %v", err)
	}

	return db, nil
}

func MigrateFs(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect(string(goose.DialectPostgres))

	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, dir)

	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}
