package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// @swagger:model User
type User struct {
    ID        uint       `json:"id" gorm:"primaryKey"`
    Username  string     `json:"username" gorm:"unique;not null"`
    Name      string     `json:"name"`
    Password  string     `json:"password,omitempty"` 
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty"`
    Tasks     []Task     `json:"tasks"`
}

func (u *User) SetPassword(password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashedPassword)
    return nil
}

func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
}
