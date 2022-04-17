package app

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	config "github.com/vvkosty/go_ya_final/internal/app/config"
	storage "github.com/vvkosty/go_ya_final/internal/app/storage"
)

type (
	Handler struct {
		Storage storage.Repository
		Config  *config.ServerConfig
	}

	UserLogin struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
)

var err error

//
//func (h *Handler) GetFullLink(c *gin.Context) {
//	urlID := c.Param("id")
//	originalURL, err := h.Storage.Find(urlID)
//	if err != nil {
//		log.Println(err)
//		c.Status(http.StatusBadRequest)
//		return
//	}
//
//	if len(originalURL) <= 0 {
//		c.Status(http.StatusBadRequest)
//		return
//	}
//
//	c.Header(`Location`, originalURL)
//	c.Status(http.StatusTemporaryRedirect)
//}
//
//func (h *Handler) Register(c *gin.Context) {
//	body, _ := io.ReadAll(c.Request.Body)
//	defer c.Request.Body.Close()
//
//	urlToEncode, err := url.ParseRequestURI(string(body))
//	if err != nil {
//		fmt.Println(err)
//		c.Status(http.StatusBadRequest)
//		return
//	}
//
//	checksum := helpers.GenerateChecksum(urlToEncode.String())
//	entity, _ := h.Storage.Find(checksum)
//
//	if entity != "" {
//		c.Status(http.StatusConflict)
//	} else {
//		userID, _ := c.Get("userId")
//		checksum, err = h.Storage.Save(urlToEncode.String(), userID.(string))
//		if err != nil {
//			log.Println(err)
//			c.Status(http.StatusBadRequest)
//			return
//		}
//		c.Status(http.StatusCreated)
//	}
//
//	c.Header(`Content-Type`, `plain/text`)
//	responseBody := fmt.Sprintf("%s/%s", h.Config.BaseURL, checksum)
//	c.Writer.Write([]byte(responseBody))
//}
//
//func (h *Handler) CreateJSONShortLink(c *gin.Context) {
//	body, _ := io.ReadAll(c.Request.Body)
//	defer c.Request.Body.Close()
//
//	requestURL := requestURL{}
//	if err := json.Unmarshal(body, &requestURL); err != nil {
//		log.Println(err)
//		c.Status(http.StatusBadRequest)
//		return
//	}
//
//	checksum := helpers.GenerateChecksum(requestURL.URL)
//	entity, _ := h.Storage.Find(checksum)
//
//	c.Header(`Content-Type`, gin.MIMEJSON)
//	if entity != "" {
//		c.Status(http.StatusConflict)
//	} else {
//		userID, _ := c.Get("userId")
//		checksum, err = h.Storage.Save(requestURL.URL, userID.(string))
//		if err != nil {
//			log.Println(err)
//			c.Status(http.StatusBadRequest)
//			return
//		}
//		c.Status(http.StatusCreated)
//	}
//
//	response := responseURL{
//		Result: fmt.Sprintf("%s/%s", h.Config.BaseURL, checksum),
//	}
//
//	encodedResponse, err := json.Marshal(&response)
//	if err != nil {
//		log.Println(err)
//		c.Status(http.StatusBadRequest)
//		return
//	}
//
//	c.Writer.Write(encodedResponse)
//}
//
//func (h *Handler) GetAllLinks(c *gin.Context) {
//	var response []listURL
//	userID, _ := c.Get("userId")
//
//	for checksum, originalURL := range h.Storage.List(userID.(string)) {
//		response = append(response, listURL{
//			ShortURL:    fmt.Sprintf("%s/%s", h.Config.BaseURL, checksum),
//			OriginalURL: originalURL,
//		})
//	}
//
//	if len(response) == 0 {
//		c.Status(http.StatusNoContent)
//		return
//	}
//
//	encodedResponse, err := json.Marshal(&response)
//	if err != nil {
//		log.Println(err)
//		c.Status(http.StatusBadRequest)
//		return
//	}
//
//	c.Header(`Content-Type`, gin.MIMEJSON)
//	c.Writer.Write(encodedResponse)
//}
//
//func (h *Handler) Ping(c *gin.Context) {
//	var ctx context.Context
//	db, err := sql.Open("pgx", h.Config.DatabaseDsn)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
//		os.Exit(1)
//	}
//	defer db.Close()
//
//	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
//	defer cancel()
//	if err := db.PingContext(ctx); err != nil {
//		panic(err)
//	}
//}
//
//func (h *Handler) CreateBatchLinks(c *gin.Context) {
//	var requestBatchURLs []requestBatchURL
//	var responseBatchURLs []responseBatchURL
//	var uniqueViolatesError *storage.UniqueViolatesError
//
//	body, _ := io.ReadAll(c.Request.Body)
//	defer c.Request.Body.Close()
//
//	if err := json.Unmarshal(body, &requestBatchURLs); err != nil {
//		log.Println(err)
//		c.Status(http.StatusBadRequest)
//		return
//	}
//
//	c.Header(`Content-Type`, gin.MIMEJSON)
//
//	for _, requestURL := range requestBatchURLs {
//		userID, _ := c.Get("userId")
//		checksum, err := h.Storage.Save(requestURL.OriginalURL, userID.(string))
//		if err != nil {
//			log.Println(err)
//			if errors.As(err, &uniqueViolatesError) {
//				c.Status(http.StatusConflict)
//				return
//			}
//			c.Status(http.StatusBadRequest)
//			return
//		}
//
//		responseBatchURLs = append(responseBatchURLs, responseBatchURL{
//			ShortURL:      fmt.Sprintf("%s/%s", h.Config.BaseURL, checksum),
//			CorrelationID: requestURL.CorrelationID,
//		})
//	}
//
//	encodedResponse, err := json.Marshal(&responseBatchURLs)
//	if err != nil {
//		log.Println(err)
//		c.Status(http.StatusBadRequest)
//		return
//	}
//
//	c.Status(http.StatusCreated)
//	c.Writer.Write(encodedResponse)
//}

// GetOrders Получение списка загруженных пользователем номеров заказов,
// статусов их обработки и информации о начислениях
func (h *Handler) GetOrders(c *gin.Context) {

}

// GetUserBalance Получение текущего баланса счёта баллов лояльности пользователя
func (h *Handler) GetUserBalance(c *gin.Context) {

}

// GetUserBalanceWithdrawals Получение информации о выводе средств с накопительного счёта пользователем
func (h *Handler) GetUserBalanceWithdrawals(c *gin.Context) {

}

// RegisterUser Регистрация пользователя
func (h *Handler) RegisterUser(c *gin.Context) {

}

// LoginUser Аутентификация пользователя
func (h *Handler) LoginUser(c *gin.Context) {

}

// SaveOrders Загрузка пользователем номера заказа для расчёта
func (h *Handler) SaveOrders(c *gin.Context) {

}

// Withdraw Списание баллов с накопительного счёта в счёт оплаты нового заказа
func (h *Handler) Withdraw(c *gin.Context) {

}
