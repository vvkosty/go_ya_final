package app

import (
	"database/sql"
	"strings"
)

type UserStorage struct {
	*sql.DB
}

func (u *UserStorage) Create(login string, password string) error {
	_, err := u.DB.Exec(
		"INSERT INTO users (login, password) VALUES ($1, $2)",
		login,
		password,
	)

	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			return NewUniqueViolatesError(err)
		}
		return err
	}

	return nil
}
