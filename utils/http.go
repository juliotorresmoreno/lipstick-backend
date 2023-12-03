package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juliotorresmoreno/tana-api/logger"
)

var log = logger.SetupLogger()

func GetToken(c *gin.Context) {
	token, _ := c.Cookie("token")
	if token == "" {
		token = c.Request.URL.Query().Get("token")
	}

	if token == "" {
		token = c.Request.Header.Get("authorization")
	}
	if len(token) > 7 && strings.ToLower(token[:6]) == "bearer" {
		token = token[7:]
	}

	if token == "" {
		log.Error("StatusUnauthorized")
		Response(c, StatusUnauthorized)
		return
	}
}
