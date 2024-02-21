package server

import (
	"github.com/gin-gonic/gin"
	"github.com/juliotorresmoreno/tana-api/server/auth"
	"github.com/juliotorresmoreno/tana-api/server/connections"
	"github.com/juliotorresmoreno/tana-api/server/conversation"
	"github.com/juliotorresmoreno/tana-api/server/events"
	"github.com/juliotorresmoreno/tana-api/server/mmlu"
)

func SetupAPIRoutes(r *gin.RouterGroup) {
	auth.SetupAPIRoutes(r)
	mmlu.SetupAPIRoutes(r.Group("/mmlu"))
	connections.SetupAPIRoutes(r.Group("/connections"))
	conversation.SetupAPIRoutes(r.Group("/conversation"))
	events.SetupAPIRoutes(r.Group("/events"))
}
