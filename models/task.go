package models

import (
	"time"
)

type TaskStatus string

const (
    Resolved   TaskStatus = "RESOLVED"
    Unresolved TaskStatus = "UNRESOLVED"
)

// @swagger:model Task
type Task struct {
    ID          uint       `json:"id" gorm:"primaryKey"`
    Name        string     `json:"name"`
    Status      TaskStatus `json:"status"`
    Description string     `json:"description"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty"`  // Campo de soft delete
    UserID      uint       `json:"user_id"`
    User        User       `json:"-"`
}
