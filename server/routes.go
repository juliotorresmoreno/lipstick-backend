package server

import (
	"github.com/gin-gonic/gin"
	"github.com/juliotorresmoreno/tana-api/server/auth"
	"github.com/juliotorresmoreno/tana-api/server/connections"
	"github.com/juliotorresmoreno/tana-api/server/conversation"
	"github.com/juliotorresmoreno/tana-api/server/credentials"
	"github.com/juliotorresmoreno/tana-api/server/events"
	"github.com/juliotorresmoreno/tana-api/server/mmlu"
	"github.com/juliotorresmoreno/tana-api/server/models"
	"github.com/juliotorresmoreno/tana-api/server/users"
)

func SetupAPIRoutes(r *gin.RouterGroup) {
	auth.SetupAPIRoutes(r)
	mmlu.SetupAPIRoutes(r.Group("/mmlu"))
	users.SetupAPIRoutes(r.Group("/users"))
	events.SetupAPIRoutes(r.Group("/events"))
	models.SetupAPIRoutes(r.Group("/models"))
	connections.SetupAPIRoutes(r.Group("/connections"))
	credentials.SetupAPIRoutes(r.Group("/credentials"))
	conversation.SetupAPIRoutes(r.Group("/conversation"))
}
