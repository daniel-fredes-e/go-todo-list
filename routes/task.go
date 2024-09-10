package routes

import (
	"go-todo-list/config"
	"go-todo-list/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// GetTasks Obtiene las tareas de un usuario autorizado.
// @Summary Get Task
// @Description Get tasks for the logged-in user
// @Produce json
// @Success 200 {object} models.Task
// @Failure 401 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /tasks [post]
func GetTasks(c *gin.Context) {
    tokenString := c.GetHeader("Authorization")
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
        return
    }

    username := claims["username"].(string)
    var user models.User
    if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }

    var tasks []models.Task
    config.DB.Where("user_id = ?", user.ID).Find(&tasks)
    c.JSON(http.StatusOK, tasks)
}

// CreateTask maneja la creación de nuevas tareas.
// @Summary Create Task
// @Description Create a new task for the logged-in user
// @Accept json
// @Produce json
// @Param task body models.Task true "Task"
// @Success 201 {object} models.Task
// @Failure 401 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /tasks [post]
func CreateTask(c *gin.Context) {
    tokenString := c.GetHeader("Authorization")
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
        return
    }

    username := claims["username"].(string)
    var user models.User
    if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }

    var task models.Task
    if err := c.ShouldBindJSON(&task); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    task.CreatedAt = time.Now()
    task.UpdatedAt = time.Now()

    if err := config.DB.Create(&task).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
        return
    }

    c.JSON(http.StatusCreated, task)
}