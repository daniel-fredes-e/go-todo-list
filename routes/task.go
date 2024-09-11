package routes

import (
	"go-todo-list/models"
	"go-todo-list/utils"
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
// @Failure 401 {object} utils.Response
// @Router /tasks [get]
func GetTasks(c *gin.Context) {
    tokenString := c.GetHeader("Authorization")
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        c.JSON(http.StatusUnauthorized, utils.Response{Name: "error", Detail: "Invalid token"})
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        c.JSON(http.StatusUnauthorized, utils.Response{Name: "error", Detail: "Invalid token claims"})
        return
    }

    username := claims["username"].(string)
    user, err := utils.GetUserExist(username)
    if err != nil {
        c.JSON(http.StatusUnauthorized, utils.Response{Name: "error", Detail: "User not found"})
        return
    }

    tasks, _ := utils.GetTasks(user.ID)
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
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /tasks [post]
func CreateTask(c *gin.Context) {
    tokenString := c.GetHeader("Authorization")
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        c.JSON(http.StatusUnauthorized, utils.Response{Name: "error", Detail: "Invalid token"})
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        c.JSON(http.StatusUnauthorized, utils.Response{Name: "error", Detail: "Invalid token claims"})
        return
    }

    username := claims["username"].(string)
    user, err := utils.GetUserExist(username)
    if err != nil {
        c.JSON(http.StatusUnauthorized, utils.Response{Name: "error", Detail: "User not found"})
        return
    }

    var task models.Task
    var input TaskRegister
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, utils.Response{Name: "error", Detail: err.Error()})
        return
    }

    task.UserID = user.ID
    task.Status = models.Unresolved
    task.CreatedAt = time.Now()
    task.UpdatedAt = time.Now()
    task.Name = input.Name
    task.Description = input.Description

    err = utils.CreateTask(&task)

    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.Response{Name: "error", Detail: "Failed to create task"})
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
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /tasks/{id}/resolve [patch]
func MarkTaskResolved(c *gin.Context) {
    taskID := c.Param("id")

    tokenString := c.GetHeader("Authorization")
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        c.JSON(http.StatusUnauthorized, utils.Response{Name: "error", Detail: "Invalid token"})
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        c.JSON(http.StatusUnauthorized, utils.Response{Name: "error", Detail: "Invalid token claims"})
        return
    }

    username := claims["username"].(string)
    user, err := utils.GetUserExist(username)
    if err != nil {
        c.JSON(http.StatusUnauthorized, utils.Response{Name: "error", Detail: "User not found"})
        return
    }

    task, err := utils.GetTasksExist(taskID, user.ID)
    if err != nil {
        c.JSON(http.StatusNotFound, utils.Response{Name: "error", Detail: "Task not found"})
        return
    }

    task.Status = models.Resolved
    task.UpdatedAt = time.Now()

    err = utils.UpdateTask(task)

    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.Response{Name: "error", Detail: "Failed to mark task as resolved"})
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
// @Success 200 {object} utils.Response "message": "Tarea Eliminada"
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /tasks/{id} [delete]
func DeleteTask(c *gin.Context) {
    taskID := c.Param("id")              // Obtiene el ID de la tarea desde la URL

    tokenString := c.GetHeader("Authorization")
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        c.JSON(http.StatusUnauthorized, utils.Response{Name: "error", Detail: "Invalid token"})
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        c.JSON(http.StatusUnauthorized, utils.Response{Name: "error", Detail: "Invalid token claims"})
        return
    }

    username := claims["username"].(string)
    user, err := utils.GetUserExist(username)
    if err != nil {
        c.JSON(http.StatusUnauthorized, utils.Response{Name: "error", Detail: "User not found"})
        return
    }

    task, err := utils.GetTasksExist(taskID, user.ID)
    if err != nil {
        c.JSON(http.StatusNotFound, utils.Response{Name: "error", Detail: "Task not found"})
        return
    }

    err = utils.DeleteTask(task)

    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.Response{Name: "error", Detail: "Failed to delete task"})
        return
    }

    c.JSON(http.StatusOK, utils.Response{Name: "message", Detail: "Task deleted successfully"})
}