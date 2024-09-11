package utils

import (
	"go-todo-list/config"
	"go-todo-list/models"
)

type Response struct {
    Name string
    Detail string
}

func CreateUser(user *models.User) error {
    return config.DB.Create(&user).Error
}

func GetUserExist(username string)  (existingUser *models.User,err error){
	result := config.DB.Where("username = ?", username).First(&existingUser)
	if result != nil {
		err = result.Error
	}
	return
}

func GetTasks(userId uint) (tasks []*models.Task,err error){
	result := config.DB.Where("user_id = ?", userId).Find(&tasks)
	if result !=nil {
		err = result.Error
	}
	return
}

func CreateTask(task *models.Task) error {
	return config.DB.Create(&task).Error
}

func GetTasksExist(taskId string, userId uint) (task *models.Task,err error) {
	result := config.DB.Where("id = ? AND user_id = ?", taskId, userId).First(&task)
	if result !=nil {
		err = result.Error
	}
	return
}

func UpdateTask(task *models.Task) error {
	return config.DB.Save(&task).Error
}

func DeleteTask(task *models.Task) error {
	return config.DB.Delete(&task).Error
}