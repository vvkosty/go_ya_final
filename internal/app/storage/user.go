package app

import (
	"database/sql"
	"strings"
)

type UserStorage struct {
	*sql.DB
}

type UserModel struct {
	Id       int
	Login    string
	Password string
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

func (u *UserStorage) Find(login string, password string) (*UserModel, error) {
	var user UserModel
	row := u.DB.QueryRow("SELECT * FROM users WHERE login = $1 and password = $2", login, password)
	err := row.Scan(&user.Id, &user.Login, &user.Password)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, NewUserNotFoundError(err)
		}

		return nil, err
	}

	return &user, nil
}
