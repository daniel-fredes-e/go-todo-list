package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
    gorm.Model
    Name        string
    Status      string `gorm:"default:'not resolved'"`
    Description string
    CreatedAt   time.Time
    UpdatedAt   time.Time
    UserID      uint
}
