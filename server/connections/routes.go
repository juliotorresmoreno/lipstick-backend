package connections

import (
	"github.com/gin-gonic/gin"
	"github.com/juliotorresmoreno/tana-api/server/mmlu"
	"github.com/juliotorresmoreno/tana-api/utils"
)

type ConnectionsRouter struct {
}

type Connection struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Feeling     string `json:"feeling"`
	PhotoURL    string `json:"photo_url"`
	Type        string `json:"type"`
}

type Connections []*Connection

func SetupAPIRoutes(r *gin.RouterGroup) {
	connections := &ConnectionsRouter{}
	r.GET("", connections.find)
}

func (h *ConnectionsRouter) find(c *gin.Context) {
	mmlus := &Connections{}
	err := mmlu.Find(mmlus)

	if err != nil {
		utils.Response(c, err)
		return
	}
	for _, item := range *mmlus {
		item.Type = "bot"
	}
	c.JSON(200, mmlus)
}
