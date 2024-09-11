package utils

import (
	"go-todo-list/config"
	"go-todo-list/models"
)

func CreateUser(user *models.User) error {
    return config.DB.Create(user).Error
}

func GetUserExist(username string)  (existingUser *models.User,err error){
	result := config.DB.Where("username = ?", username).First(&existingUser)
	if result != nil {
		err = result.Error
	}
	return
}