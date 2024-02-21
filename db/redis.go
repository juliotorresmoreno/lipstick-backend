package db

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	conf, _ := redis.ParseURL(os.Getenv("REDIS_URL"))
	return redis.NewClient(conf)
}
