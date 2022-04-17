package app

import "database/sql"

type Repository interface {
	Find(id string) (string, error)
	Save(url string, userID string) (string, error)
	List(userID string) map[string]string
	Close() error
	Instance() *sql.DB
}
