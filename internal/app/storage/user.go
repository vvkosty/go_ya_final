package app

import (
	"database/sql"
	"strings"
)

type UserStorage struct {
	*sql.DB
}

type UserModel struct {
	ID       int
	Login    string
	Password string
}

var user *UserModel

func (u *UserStorage) Create(login string, password string) (*UserModel, error) {
	user = &UserModel{}
	row := u.DB.QueryRow(
		"INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id, login, password",
		login,
		password,
	)

	if row.Err() != nil {
		if strings.Contains(row.Err().Error(), "23505") {
			return nil, NewUniqueViolatesError(row.Err())
		}
		return nil, row.Err()
	}

	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, NewEntityNotFoundError(err)
		}

		return nil, err
	}

	return user, nil
}

func (u *UserStorage) Find(login string, password string) (*UserModel, error) {
	user = &UserModel{}

	row := u.DB.QueryRow("SELECT * FROM users WHERE login = $1 and password = $2", login, password)
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, NewEntityNotFoundError(err)
		}

		return nil, err
	}

	return user, nil
}
