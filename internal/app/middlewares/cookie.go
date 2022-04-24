package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	config "github.com/vvkosty/go_ya_final/internal/app/config"
	"github.com/vvkosty/go_ya_final/internal/app/helpers"
)

type Middleware struct {
	Config  *config.Config
	Encoder *helpers.Encoder
}

func (m *Middleware) NeedAuth(c *gin.Context) {
	authCookie, _ := c.Request.Cookie("user")
	if authCookie == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID, err := m.Encoder.Decrypt([]byte(authCookie.Value))
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
