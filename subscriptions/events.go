package subscriptions

import (
	"context"
	"encoding/json"

	"github.com/juliotorresmoreno/tana-api/db"
	"github.com/juliotorresmoreno/tana-api/logger"
	"github.com/juliotorresmoreno/tana-api/models"
	"github.com/juliotorresmoreno/tana-api/server/events"
	"github.com/redis/go-redis/v9"
)

var log = logger.SetupLogger()

func Setup() {
	rdb, err := db.NewRedisClient()
	if err != nil {
		log.Fatal(err)
	}
	sub := rdb.Subscribe(context.Background(), "events")
	go handleEvents(sub)
}

func handleEvents(sub *redis.PubSub) {
	for {
		msg, err := sub.ReceiveMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		evt := &models.Event{}
		if err = json.Unmarshal([]byte(msg.Payload), evt); err != nil {
			continue
		}
		events.Manager.Event <- evt
	}
}
