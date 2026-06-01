package db

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func Open() (*sql.DB, error) {
	dbSrc := fmt.Sprintf(
		"host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	db, err := sql.Open("pgx", dbSrc)

	if err != nil {
		panic("Could not connect to DB." + err.Error())
	}

	if err := db.Ping(); err != nil {
		fmt.Printf("db: open %v", err)
	}

	db.SetMaxOpenConns(10)

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
