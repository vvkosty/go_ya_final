package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type (
	PostgresDatabase struct {
		db *sql.DB
	}
)

type UniqueViolatesError struct{ Err error }

func (uve *UniqueViolatesError) Error() string {
	return fmt.Sprintf("UniqueViolatesError: %v", uve.Err)
}

func NewUniqueViolatesError(err error) error {
	return &UniqueViolatesError{
		Err: err,
	}
}

type EntityNotFoundError struct{ Err error }

func (unf *EntityNotFoundError) Error() string {
	return fmt.Sprintf("EntityNotFoundError: %v", unf.Err)
}

func NewEntityNotFoundError(err error) error {
	return &EntityNotFoundError{
		Err: err,
	}
}

func NewPostgresDatabase(dsn string) *PostgresDatabase {
	var md PostgresDatabase
	var err error

	md.db, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return &md
}

func (m *PostgresDatabase) Close() error {
	return m.db.Close()
}

func (m *PostgresDatabase) Instance() *sql.DB {
	return m.db
}
