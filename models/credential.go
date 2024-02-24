package models

import (
	"time"

	"gorm.io/gorm"
)

type Credential struct {
	ID         uint           `gorm:"primaryKey"`
	ApiKey     string         `gorm:"type:varchar(100);default:'';nullable"`
	ApiSecret  string         `gorm:"type:varchar(100);default:'';nullable"`
	UserId     uint           `gorm:"not null"`
	User       User           `gorm:"foreignKey:UserId"`
	LastUsed   *time.Time     `gorm:"type:timestamptz"`
	CreationAt time.Time      `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time      `gorm:"type:timestamptz"`
	DeletedAt  gorm.DeletedAt `gorm:"type:timestamptz"`
}

func (u Credential) TableName() string {
	return "credentials"
}
