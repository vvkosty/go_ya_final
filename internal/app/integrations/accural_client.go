package integrations

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type AccrualAPIClient struct {
	BaseAddress string
}

type AccrualOrder struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func (a AccrualOrder) String() string {
	return fmt.Sprintf("Number: %s, Status: %s, Accrual: %f", a.Number, a.Status, a.Accrual)
}

func (as AccrualAPIClient) GetAccrualOrderData(orderID string) (*AccrualOrder, error) {
	var ao AccrualOrder

	res, err := http.Get(fmt.Sprintf("%s/api/orders/%s", as.BaseAddress, orderID))
	if err != nil {
		return nil, errors.New("accrual response error")
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("accrual internal error: " + strconv.Itoa(res.StatusCode))
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&ao)
	if err != nil {
		return nil, err
	}

	return &ao, nil
}
