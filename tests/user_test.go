package tests

import (
	"go-todo-list/config"
	"go-todo-list/models"
	"go-todo-list/routes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRouter() *gin.Engine {
    r := gin.Default()
    r.POST("/register", routes.Register)
    r.POST("/login", routes.Login)
    r.POST("/tasks", routes.CreateTask)
    r.GET("/tasks", routes.GetTasks)
    r.PATCH("/tasks/:id/resolve", routes.MarkTaskResolved)
    r.DELETE("/tasks/:id", routes.DeleteTask)
    return r
}

func setupDatabase() *gorm.DB {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&models.User{})
    return db
}

func TestRegister(t *testing.T) {
    router := setupRouter()
    db := setupDatabase()
    config.DB = db

    //Register

    newUser := `{"username":"testuser","name": "New User","password":"testpass"}`

    req, _ := http.NewRequest("POST", "/register", strings.NewReader(newUser))
    req.Header.Set("Content-Type", "application/json")

    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusCreated, recorder.Code)
}

func LoginUser(t *testing.T) {
    
    TestRegister(t)
    router := setupRouter()

    //Login

    logUser := `{"username":"testuser","password":"testpass"}`

    req, _ := http.NewRequest("POST", "/login", strings.NewReader(logUser))
    req.Header.Set("Content-Type", "application/json")

    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)
}

func CreateTask(t *testing.T) {
    router := setupRouter()
    LoginUser(t)
    
    //Create Task

    task := `{"name": "New Task","description": "This is a new task"}`

    req, _ := http.NewRequest("POST", "/tasks", strings.NewReader(task))
    req.Header.Set("Content-Type", "application/json")

    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusCreated, recorder.Code)

}

func GetTasks(t *testing.T) {
    router := setupRouter()
    CreateTask(t)
    
    //Get Tasks

    req, _ := http.NewRequest("GET", "/tasks", nil)
    req.Header.Set("Content-Type", "application/json")

    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)

}

func MarkTaskResolved(t *testing.T) {
    router := setupRouter()
    CreateTask(t)
    
    //Mark Tasks Resolved

    req, _ := http.NewRequest("PATCH", "/tasks/1/resolve", nil)
    req.Header.Set("Content-Type", "application/json")

    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)

}

func DeleteTask(t *testing.T) {
    router := setupRouter()
    CreateTask(t)
    
    //Delete Tasks

    req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
    req.Header.Set("Content-Type", "application/json")

    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)

}