package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID         uint           `gorm:"primaryKey"`
	Content    string         `gorm:"type:text;default:'';nullable"`
	OwnerId    uint           `gorm:"not null"`
	Owner      User           `gorm:"foreignKey:OwnerId"`
	MmluId     uint           `gorm:"not null"`
	Mmlu       Mmlu           `gorm:"foreignKey:MmluId"`
	Role       string         `gorm:"type:varchar(255);default:'';not null"`
	CreationAt time.Time      `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time      `gorm:"type:timestamptz"`
	DeletedAt  gorm.DeletedAt `gorm:"type:timestamptz"`
}

func (u Message) TableName() string {
	return "messages"
}
