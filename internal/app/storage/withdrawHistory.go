package app

import (
	"database/sql"
	"log"
	"strings"
	"time"
)

type WithdrawHistoryStorage struct {
	*sql.DB
}

type WithdrawHistoryModel struct {
	ID        int
	UserID    int
	OrderID   string
	Withdraw  float64
	CreatedAt time.Time
}

func (wh *WithdrawHistoryStorage) Create(userID int, withdraw float64) (*WithdrawHistoryModel, error) {
	withdrawModel := &WithdrawHistoryModel{}
	row := wh.DB.QueryRow(
		`INSERT INTO withdraw_history (user_id, withdraw, created_at) 
				VALUES ($1, $2, $3) 
				RETURNING id, user_id, withdraw, created_at`,
		userID, withdraw, time.Now(),
	)

	if row.Err() != nil {
		if strings.Contains(row.Err().Error(), "23505") {
			return nil, NewUniqueViolatesError(row.Err())
		}
		return nil, row.Err()
	}

	err := row.Scan(&withdrawModel.ID, &withdrawModel.UserID, &withdrawModel.Withdraw, &withdrawModel.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, NewEntityNotFoundError(err)
		}

		return nil, err
	}

	return withdrawModel, nil
}

func (wh *WithdrawHistoryStorage) List(userID int) ([]*WithdrawHistoryModel, error) {
	var list []*WithdrawHistoryModel

	rows, err := wh.DB.Query(
		`SELECT id, user_id, order_id, withdraw, created_at 
				FROM withdraw_history 
				WHERE user_id = $1 
				ORDER BY created_at DESC`,
		userID,
	)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		withdrawModel := &WithdrawHistoryModel{}
		err := rows.Scan(
			&withdrawModel.ID,
			&withdrawModel.UserID,
			&withdrawModel.OrderID,
			&withdrawModel.Withdraw,
			&withdrawModel.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, withdrawModel)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, nil
}
