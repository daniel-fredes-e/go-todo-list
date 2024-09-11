package routes

import (
	"go-todo-list/config"
	"go-todo-list/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type TaskRegister struct {
    Name string `json:"name"`
    Description string `json:"description"`
}

// GetTasks Obtiene las tareas de un usuario autorizado.
// @Summary Get Tasks
// @Description Obtiene las tareas de un usuario autorizado
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.Task "Lista de tareas del usuario"
// @Failure 401 {object} map[string]string "error": "Unauthorized"
// @Router /tasks [get]
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
// @Description creación de nuevas tareas para usuario autorizado
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param task body TaskRegister true "Task"
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
    var input TaskRegister
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    task.UserID = user.ID
    task.Status = models.Unresolved // Estado inicial no resuelto
    task.CreatedAt = time.Now()
    task.UpdatedAt = time.Now()
    task.Name = input.Name
    task.Description = input.Description

    if err := config.DB.Create(&task).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
        return
    }

    c.JSON(http.StatusCreated, task)
}

// MarkTaskResolved marca una tarea como resuelta
// @Summary Mark Task Resolved
// @Description Actualiza el estado de una tarea del usuario autenticado a "RESOLVED".
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Security ApiKeyAuth
// @Success 200 {object} models.Task "Tarea Resuelta"
// @Failure 404 {object} map[string]string "error": "Tarea no encontrada"
// @Failure 500 {object} map[string]string "error": "Falló al actualizar la tarea"
// @Router /tasks/{id}/resolve [patch]
func MarkTaskResolved(c *gin.Context) {
    taskID := c.Param("id")              // Obtiene el ID de la tarea desde la URL

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
    // Verifica si la tarea existe y pertenece al usuario
    if err := config.DB.Where("id = ? AND user_id = ?", taskID, user.ID).First(&task).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }

    task.Status = models.Resolved
    task.UpdatedAt = time.Now()

    if err := config.DB.Save(&task).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark task as resolved"})
        return
    }

    c.JSON(http.StatusOK, task)
}

// DeleteTask elimina una tarea del usuario autenticado
// @Summary Delete Task
// @Description Elimina una tarea específica del usuario autenticado dado el ID de la tarea.
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string "message": "Tarea Eliminada"
// @Failure 404 {object} map[string]string "error": "Tarea no encontrada"
// @Failure 500 {object} map[string]string "error": "Falló al eliminar la tarea"
// @Router /tasks/{id} [delete]
func DeleteTask(c *gin.Context) {
    taskID := c.Param("id")              // Obtiene el ID de la tarea desde la URL

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
    // Verifica si la tarea existe y pertenece al usuario
    if err := config.DB.Where("id = ? AND user_id = ?", taskID, user.ID).First(&task).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }

    if err := config.DB.Delete(&task).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}