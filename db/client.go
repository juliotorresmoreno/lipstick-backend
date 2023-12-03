package db

import (
	"errors"
	"log"
	"os"
	"time"

	logger2 "github.com/juliotorresmoreno/tana-api/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var applog = logger2.SetupLogger()
var DefaultClient *gorm.DB

func Init() {
	var err error
	DefaultClient, err = NewClient()
	if err == nil {
		applog.Info("Connected to database")
	} else {
		applog.Panic("Failed conection to database")
	}
}

func NewClient() (*gorm.DB, error) {
	driver := os.Getenv("DATABASE_DRIVER")
	url := os.Getenv("DATABASE_URL")
	switch driver {
	case "postgres":
		return newPostgresClient(url, 10)
	}
	return &gorm.DB{}, errors.New("postgres isn't valid")
}

func newPostgresClient(dsn string, poolSize int) (*gorm.DB, error) {
	config := &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(poolSize)

	return db, nil
}
