package connections

import (
	"github.com/gin-gonic/gin"
	"github.com/juliotorresmoreno/tana-api/models"
	"github.com/juliotorresmoreno/tana-api/server/mmlu"
	"github.com/juliotorresmoreno/tana-api/utils"
)

type ConnectionsRouter struct {
}

func SetupAPIRoutes(r *gin.RouterGroup) {
	connections := &ConnectionsRouter{}
	r.GET("", connections.find)
}

func (h *ConnectionsRouter) find(c *gin.Context) {
	connections := &models.Connections{}
	err := mmlu.Find(connections)

	if err != nil {
		utils.Response(c, err)
		return
	}
	for _, item := range *connections {
		item.Type = "bot"
	}
	c.JSON(200, connections)
}
