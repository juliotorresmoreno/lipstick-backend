package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/juliotorresmoreno/tana-api/db"
	"github.com/juliotorresmoreno/tana-api/logger"
	"github.com/juliotorresmoreno/tana-api/middlewares"
	"github.com/juliotorresmoreno/tana-api/server"
	"github.com/juliotorresmoreno/tana-api/subscriptions"
)

func main() {
	time.Local = time.UTC // default to UTC for all time values
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	logger.SetupLogrus()
	db.Setup()
	subscriptions.Setup()

	r := gin.Default()
	r.Use(middlewares.AuthMiddleware())
	server.SetupAPIRoutes(r.Group("api"))
	r.Run(os.Getenv("ADDR"))
}
