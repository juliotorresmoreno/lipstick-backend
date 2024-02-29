package models

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type Model struct {
	Code     string `json:"code"`
	Provider string `json:"provider"`
}

func SetupAPIRoutes(g *gin.RouterGroup) {
	g.GET("", func(ctx *gin.Context) {
		models := strings.Split(os.Getenv("OLLAMA_MODELS"), " ")
		result := []*Model{}
		for _, model := range models {
			result = append(result, &Model{
				Code:     model,
				Provider: "ollama",
			})
		}
		ctx.JSON(200, result)
	})
}
