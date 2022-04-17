package app

import "database/sql"

type Repository interface {
	Close() error
	Instance() *sql.DB
}
