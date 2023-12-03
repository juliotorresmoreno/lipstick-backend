package cache

import (
	"os"

	"github.com/juliotorresmoreno/tana-api/logger"
	"github.com/redis/go-redis/v9"
)

var log = logger.SetupLogger()
var DefaultClient *redis.Client

func Init() {
	var err error
	DefaultClient, err = NewCache()
	if err == nil {
		log.Info("Connected to redis")
	} else {
		log.Panic("Failed conection to redis")
	}
}

func NewCache() (*redis.Client, error) {
	opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		return &redis.Client{}, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
		PoolSize: 10,
	})
	return rdb, nil
}
