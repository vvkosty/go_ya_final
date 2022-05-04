package app

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/theplant/luhn"
	config "github.com/vvkosty/go_ya_final/internal/app/config"
	"github.com/vvkosty/go_ya_final/internal/app/enums"
	"github.com/vvkosty/go_ya_final/internal/app/helpers"
	"github.com/vvkosty/go_ya_final/internal/app/integrations"
	storage "github.com/vvkosty/go_ya_final/internal/app/storage"
)

type (
	Handler struct {
		Config  *config.Config
		Encoder *helpers.Encoder

		UserStorage            *storage.UserStorage
		OrderStorage           *storage.OrderStorage
		UserBalanceStorage     *storage.UserBalanceStorage
		WithdrawHistoryStorage *storage.WithdrawHistoryStorage

		AccrualAPIClient *integrations.AccrualAPIClient
	}

	userLoginDto struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	withdrawDto struct {
		OrderID string  `json:"order"`
		Sum     float64 `json:"sum"`
	}

	getOrdersResponseDto struct {
		Number     string  `json:"number"`
		Status     string  `json:"status"`
		Accrual    float64 `json:"accrual"`
		UploadedAt string  `json:"uploaded_at"`
	}

	getUserBalanceResponseDto struct {
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}

	getUserWithdrawalsResponseDto struct {
		OrderID     string  `json:"order"`
		Sum         float64 `json:"sum"`
		ProcessedAt string  `json:"processed_at"`
	}
)

// GetOrders Получение списка загруженных пользователем номеров заказов,
// статусов их обработки и информации о начислениях
func (h *Handler) GetOrders(c *gin.Context) {
	var orders []*storage.OrderModel
	var response []*getOrdersResponseDto

	userID, exist := h.getUserID(c)
	if !exist {
		return
	}

	orders, err := h.OrderStorage.List(userID)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	for _, order := range orders {
		response = append(response, &getOrdersResponseDto{
			Number:     order.ID,
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: order.UploadedAt.Format(time.RFC3339),
		})
	}

	encodedResponse, err := json.Marshal(&response)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Header(`Content-Type`, gin.MIMEJSON)
	_, err = c.Writer.Write(encodedResponse)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
}

// GetUserBalance Получение текущего баланса счёта баллов лояльности пользователя
func (h *Handler) GetUserBalance(c *gin.Context) {
	userID, exist := h.getUserID(c)
	if !exist {
		return
	}

	userBalance, err := h.UserBalanceStorage.GetBalance(userID)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	response := &getUserBalanceResponseDto{
		Current:   userBalance.Balance,
		Withdrawn: userBalance.Withdraw,
	}

	encodedResponse, err := json.Marshal(&response)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Header(`Content-Type`, gin.MIMEJSON)
	_, err = c.Writer.Write(encodedResponse)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
}

// GetUserWithdrawals Получение информации о выводе средств с накопительного счёта пользователем
func (h *Handler) GetUserWithdrawals(c *gin.Context) {
	var withdraws []*storage.WithdrawHistoryModel
	var response []*getUserWithdrawalsResponseDto

	userID, exist := h.getUserID(c)
	if !exist {
		return
	}

	withdraws, err := h.WithdrawHistoryStorage.List(userID)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if len(withdraws) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	for _, withdrawHistory := range withdraws {
		response = append(response, &getUserWithdrawalsResponseDto{
			OrderID:     withdrawHistory.OrderID,
			Sum:         withdrawHistory.Withdraw,
			ProcessedAt: withdrawHistory.CreatedAt.Format(time.RFC3339),
		})
	}

	encodedResponse, err := json.Marshal(&response)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Header(`Content-Type`, gin.MIMEJSON)
	_, err = c.Writer.Write(encodedResponse)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
}

// RegisterUser Регистрация пользователя
func (h *Handler) RegisterUser(c *gin.Context) {
	var uniqueViolatesError *storage.UniqueViolatesError
	var userLogin *userLoginDto

	body, _ := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	if err := json.Unmarshal(body, &userLogin); err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := h.UserStorage.Create(userLogin.Login, userLogin.Password)
	if err != nil {
		log.Println(err)
		if errors.As(err, &uniqueViolatesError) {
			c.Status(http.StatusConflict)
			return
		}
		c.Status(http.StatusBadRequest)
		return
	}

	err = h.UserBalanceStorage.Init(user.ID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	authCookie, err := h.Encoder.Encrypt(strconv.Itoa(user.ID))
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.SetCookie(
		"user",
		authCookie,
		3600,
		"/",
		h.Config.Host,
		false,
		false,
	)
}

// LoginUser Аутентификация пользователя
func (h *Handler) LoginUser(c *gin.Context) {
	var userNotFoundError *storage.EntityNotFoundError
	var userLogin *userLoginDto

	body, _ := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	if err := json.Unmarshal(body, &userLogin); err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := h.UserStorage.Find(userLogin.Login, userLogin.Password)
	if err != nil {
		log.Println(err)
		if errors.As(err, &userNotFoundError) {
			c.Status(http.StatusUnauthorized)
			return
		}
		c.Status(http.StatusBadRequest)
		return
	}

	authCookie, err := h.Encoder.Encrypt(strconv.Itoa(user.ID))
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.SetCookie(
		"user",
		authCookie,
		3600,
		"/",
		h.Config.Host,
		false,
		false,
	)
}

// SaveOrder Загрузка пользователем номера заказа для расчёта
func (h *Handler) SaveOrder(c *gin.Context) {
	var accrualOrder *integrations.AccrualOrder
	var order *storage.OrderModel
	var entityNotFoundError *storage.EntityNotFoundError

	body, _ := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	orderID, err := strconv.Atoi(string(body))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if !luhn.Valid(orderID) {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	order, err = h.OrderStorage.Find(strconv.Itoa(orderID))
	if err != nil {
		if !errors.As(err, &entityNotFoundError) {
			c.Status(http.StatusBadRequest)
			return
		}
	}

	userID, exist := h.getUserID(c)
	if !exist {
		return
	}

	if order != nil {
		if order.UserID == userID {
			c.Status(http.StatusOK)
			return
		} else {
			c.Status(http.StatusConflict)
			return
		}
	}

	accrualOrder, err = h.AccrualAPIClient.GetAccrualOrderData(strconv.Itoa(orderID))
	if err != nil {
		log.Println(err)
	}

	accrual := 0.0
	status := enums.AccrualOrderStatusNew.String()
	if accrualOrder != nil {
		accrual = accrualOrder.Accrual
		status = accrualOrder.Status
	}

	orderModel := &storage.OrderModel{
		ID:      strconv.Itoa(orderID),
		UserID:  userID,
		Accrual: accrual,
		Status:  status,
	}
	err = h.OrderStorage.Create(orderModel)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	if accrual > 0 {
		userBalance, err := h.UserBalanceStorage.GetBalance(userID)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		userBalance.Balance += accrual
		err = h.UserBalanceStorage.Update(userBalance)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}
	}

	c.Status(http.StatusAccepted)
}

// Withdraw Списание баллов с накопительного счёта в счёт оплаты нового заказа
func (h *Handler) Withdraw(c *gin.Context) {
	var withdraw *withdrawDto

	body, _ := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	if err := json.Unmarshal(body, &withdraw); err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	orderID, err := strconv.Atoi(withdraw.OrderID)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if !luhn.Valid(orderID) {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	userID, exist := h.getUserID(c)
	if !exist {
		return
	}

	userBalance, err := h.UserBalanceStorage.GetBalance(userID)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if userBalance.Balance < withdraw.Sum {
		c.Status(http.StatusPaymentRequired)
		return
	}

	userBalance.Balance -= withdraw.Sum
	userBalance.Withdraw += withdraw.Sum
	err = h.UserBalanceStorage.Update(userBalance)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	_, err = h.WithdrawHistoryStorage.Create(userID, withdraw.Sum)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}
}

func (h *Handler) getUserID(c *gin.Context) (int, bool) {
	val, exist := c.Get("userId")
	if !exist {
		log.Println("userId not exist")
		c.Status(http.StatusInternalServerError)
		return 0, false
	}
	userID, _ := strconv.Atoi(val.(string))
	return userID, true
}
