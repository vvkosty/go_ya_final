package middlewares

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	config "github.com/vvkosty/go_ya_final/internal/app/config"
	"github.com/vvkosty/go_ya_final/internal/app/helpers"
)

const secretKey = "123e4567-e89b-12d3-a456-42661417"

type Middleware struct {
	Config *config.ServerConfig
}

func (m *Middleware) NeedAuth(c *gin.Context) {
	authCookie, _ := c.Request.Cookie("user")
	if authCookie == nil {
		log.Printf("Не авторизован")
		c.Status(http.StatusUnauthorized)

		return
	}

	userID, err := m.decrypt([]byte(authCookie.Value))
	if err != nil {
		log.Println(err)
		return
	}

	c.SetCookie(
		"user",
		authCookie.Value,
		3600,
		"/",
		m.Config.Host,
		false,
		false,
	)

	c.Set("userId", userID)

	c.Next()
}

func (m *Middleware) encrypt(value string) (string, error) {
	// получаем cipher.Block
	aesblock, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	// создаём вектор инициализации
	nonce, err := helpers.GenerateRandom(aesgcm.NonceSize())
	if err != nil {
		return "", err
	}

	// зашифровываем
	dst := aesgcm.Seal(nonce, nonce, []byte(value), nil)

	return hex.EncodeToString(dst), nil
}

func (m *Middleware) decrypt(value []byte) (string, error) {
	var decodedValue []byte
	decodedValue, _ = hex.DecodeString(string(value))

	// получаем cipher.Block
	aesblock, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	// создаём вектор инициализации
	nonce, cipherText := decodedValue[:aesgcm.NonceSize()], decodedValue[aesgcm.NonceSize():]

	// расшифровываем
	userID, err := aesgcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(userID), nil
}
