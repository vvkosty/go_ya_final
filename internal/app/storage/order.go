package app

import (
	"database/sql"
	"log"
	"strings"
	"time"
)

type OrderStorage struct {
	*sql.DB
}

type OrderModel struct {
	ID         string
	UserID     int
	Accrual    float64
	Status     string
	UploadedAt time.Time
}

func (o *OrderStorage) Create(order *OrderModel) error {
	_, err := o.DB.Exec(
		"INSERT INTO orders (id, user_id, accrual, status, uploaded_at) VALUES ($1, $2, $3, $4, $5)",
		order.ID,
		order.UserID,
		order.Accrual,
		order.Status,
		time.Now(),
	)

	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			return NewUniqueViolatesError(err)
		}
		return err
	}

	return nil
}

func (o *OrderStorage) Find(id string) (*OrderModel, error) {
	row := o.DB.QueryRow("SELECT id, user_id, accrual, status FROM orders WHERE id = $1", id)

	order := &OrderModel{}
	err := row.Scan(&order.ID, &order.UserID, &order.Accrual, &order.Status)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, NewEntityNotFoundError(err)
		}

		return nil, err
	}

	return order, nil
}

func (o *OrderStorage) List(userID int) ([]*OrderModel, error) {
	var list []*OrderModel

	rows, err := o.DB.Query(
		"SELECT id, user_id, accrual, status, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC",
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
		orderModel := &OrderModel{}
		err := rows.Scan(
			&orderModel.ID,
			&orderModel.UserID,
			&orderModel.Accrual,
			&orderModel.Status,
			&orderModel.UploadedAt,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, orderModel)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, nil
}
