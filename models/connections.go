package models

import (
	"time"

	"gorm.io/gorm"
)

type Connection struct {
	ID          uint           `gorm:"primaryKey;autoIncrement"`
	Name        string         `gorm:"type:varchar(100);default:''"`
	Description string         `gorm:"type:varchar(256);default:''"`
	PhotoURL    string         `gorm:"type:varchar(1000);default:''"`
	OwnerId     uint           `gorm:"not null"`
	Owner       User           `gorm:"foreignKey:OwnerId"`
	MmluId      uint           `gorm:"not null"`
	Mmlu        Mmlu           `gorm:"foreignKey:MmluId"`
	CreationAt  time.Time      `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `gorm:"type:timestamptz"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamptz"`
}

func (c Connection) TableName() string {
	return "connections"
}
