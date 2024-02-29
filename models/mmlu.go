package models

import (
	"time"

	"gorm.io/gorm"
)

type Mmlu struct {
	ID          uint            `gorm:"primaryKey;autoIncrement"`
	Name        string          `gorm:"type:varchar(100);default:''"`
	Description string          `gorm:"type:varchar(256);default:''"`
	PhotoURL    string          `gorm:"type:varchar(1000);default:''"`
	Model       string          `gorm:"type:varchar(100);default:''"`
	Provider    string          `gorm:"type:varchar(256);default:'';check:provider IN ('ollama', 'openai')"`
	OwnerId     uint            `gorm:"not null"`
	Owner       User            `gorm:"foreignKey:OwnerId"`
	CreationAt  time.Time       `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time       `gorm:"type:timestamptz"`
	DeletedAt   *gorm.DeletedAt `gorm:"type:timestamptz"`
}

func (Mmlu) ProviderCheck() string {
	return "provider IN ('ollama', 'openai')"
}

func (m Mmlu) TableName() string {
	return "mmlus"
}
