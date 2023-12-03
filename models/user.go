package models

import (
	"time"

	"gorm.io/gorm"
)

const tableName = "users"

type User struct {
	ID             uint           `gorm:"primaryKey"`
	ValidationCode string         `gorm:"type:varchar(6)"`
	Verified       bool           `gorm:"default:false"`
	Name           string         `gorm:"type:varchar(100);default:'';nullable"`
	LastName       string         `gorm:"type:varchar(100);default:'';nullable"`
	Email          string         `gorm:"type:varchar(300);default:'';nullable"`
	Password       string         `gorm:"type:varchar(512);default:'';not null"`
	PhotoURL       string         `gorm:"type:varchar(1000);default:'';nullable"`
	Phone          string         `gorm:"type:varchar(15);default:'';unique"`
	Rol            string         `gorm:"type:varchar(15);default:''"`
	CreationAt     time.Time      `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time      `gorm:"type:timestamptz"`
	DeletedAt      gorm.DeletedAt `gorm:"type:timestamptz"`
}

func (u *User) TableName() string {
	return tableName
}
