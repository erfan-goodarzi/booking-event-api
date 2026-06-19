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

type DB interface {
	Ping() error
	Close() error
	SetMaxOpenConns(int)
	SetMaxIdleConns(int)
	SetConnMaxLifetime(time.Duration)
	Begin() (*sql.Tx, error)
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
}

type Tx interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
	Commit() error
	Rollback() error
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

func ConnectDB(cfg DatabaseConfig) (DB, error) {
	db, err := Open(cfg)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenDBConn)
	db.SetMaxIdleConns(maxDBIdleConn)
	db.SetConnMaxLifetime(maxDBLifeTime)

	return db, nil
}

func Open(cfg DatabaseConfig) (*sql.DB, error) {
	driver := cfg.Driver
	if driver == "" {
		driver = "pgx"
	}

	db, err := sql.Open(driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	return db, nil
}

func MigrateFs(db DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()

	sqlDB, ok := db.(*sql.DB)
	if !ok {
		return fmt.Errorf("migrate: unsupported db type %T", db)
	}

	return Migrate(sqlDB, dir)
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
