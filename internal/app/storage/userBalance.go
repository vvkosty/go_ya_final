package app

import (
	"database/sql"
	"strings"
	"time"
)

type UserBalanceStorage struct {
	*sql.DB
}

type UserBalanceModel struct {
	ID        string
	UserID    int
	Balance   float64
	Withdraw  float64
	UpdatedAt time.Time
}

func (ub *UserBalanceStorage) Init(userID int) error {
	_, err := ub.DB.Exec(
		`INSERT INTO user_balance (user_id, balance, withdraw, updated_at) 
				VALUES ($1, $2, $3, $4) 
				RETURNING id, user_id, balance, withdraw, updated_at`,
		userID, 0, 0, time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (ub *UserBalanceStorage) GetBalance(userID int) (*UserBalanceModel, error) {
	row := ub.DB.QueryRow(
		"SELECT id, user_id, balance, withdraw, updated_at FROM user_balance WHERE user_id = $1",
		userID,
	)

	userBalance := &UserBalanceModel{}
	err := row.Scan(
		&userBalance.ID,
		&userBalance.UserID,
		&userBalance.Balance,
		&userBalance.Withdraw,
		&userBalance.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, NewEntityNotFoundError(err)
		}

		return nil, err
	}

	return userBalance, nil
}

func (ub *UserBalanceStorage) Update(model *UserBalanceModel) error {
	_, err := ub.DB.Exec(
		`UPDATE user_balance SET user_id = $2, balance = $3, withdraw = $4, updated_at = $5 WHERE id = $1;`,
		model.ID, model.UserID, model.Balance, model.Withdraw, time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}
