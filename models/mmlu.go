package models

import (
	"time"

	"gorm.io/gorm"
)

type Mmlu struct {
	ID          uint           `gorm:"primaryKey;autoIncrement"`
	Verified    bool           `gorm:"default:false"`
	Name        string         `gorm:"type:varchar(100);default:''"`
	Description string         `gorm:"type:varchar(256);default:''"`
	Feeling     string         `gorm:"type:varchar(100);default:''"`
	PhotoURL    string         `gorm:"type:varchar(1000);default:''"`
	Type        string         `gorm:"type:varchar(10);default:''"`
	Path        string         `gorm:"type:varchar(256);default:''"`
	CreationAt  time.Time      `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `gorm:"type:timestamptz"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamptz"`
}

func (m Mmlu) TableName() string {
	return "mmlus"
}
