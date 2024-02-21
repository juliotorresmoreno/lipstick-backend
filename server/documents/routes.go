package documents

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/juliotorresmoreno/tana-api/logger"
	"github.com/juliotorresmoreno/tana-api/utils"
)

var log = logger.SetupLogger()

type DocumentsRouter struct {
}

func SetupAPIRoutes(r *gin.RouterGroup) {
	documents := &DocumentsRouter{}
	r.GET("/:id", documents.findOne)
}

func (h *DocumentsRouter) findOne(c *gin.Context) {
	token, err := utils.GetToken(c)
	if err != nil {
		log.Error("Error getting token", err)
		utils.Response(c, err)
		return
	}
	session, err := utils.ValidateSession(token)
	if err != nil {
		log.Error("Error validating session", err)
		utils.Response(c, err)
		return
	}

	fmt.Sprintln(session)

	c.JSON(200, gin.H{})
}
