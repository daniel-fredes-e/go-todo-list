package routes

import (
	"go-todo-list/models"
	"go-todo-list/utils"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type UserResponse struct {
    ID        uint      `json:"id"`
    Username  string    `json:"username"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type UserRegister struct {
    Username string `json:"username"`
    Name     string `json:"name"`
    Password string `json:"password"`
}

type UserLogin struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

var jwtKey = []byte("your_secret_key")

// Login
// @Summary Login
// @Description Login de usuario, devuelve Token JWT
// @Tags users
// @Accept  json
// @Produce  json
// @Param login body UserLogin true "Login Input"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /login [post]
func Login(c *gin.Context) {
    var input UserLogin
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, utils.Response{Name: "error", Detail: err.Error()})
        return
    }
    username := input.Username
    password := input.Password

    user, err := utils.GetUserExist(username)

    if err != nil {
        c.JSON(http.StatusUnauthorized, utils.Response{Name: "error", Detail: "Invalid username or password"})
        return
    }

    if !user.CheckPassword(password) {
        c.JSON(http.StatusUnauthorized, utils.Response{Name: "error", Detail: "Invalid username or password"})
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": username,
    })

    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.Response{Name: "error", Detail: "Could not generate token"})
        return
    }

    c.JSON(http.StatusOK, utils.Response{Name: "token", Detail: tokenString})
}

// Register maneja el registro de nuevos usuarios.
// @Summary Register
// @Description Registra nuevo usuario
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserRegister true "User Registration"
// @Success 201 {object} UserResponse
// @Failure 400 {object} utils.Response
// @Router /register [post]
func Register(c *gin.Context) {
    var user models.User
    var input UserRegister
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, utils.Response{Name: "error", Detail: err.Error()})
        return
    }

    if _, err := utils.GetUserExist(input.Username); err == nil {
        c.JSON(http.StatusBadRequest, utils.Response{Name: "error", Detail: "Username already taken"})
        return
    }

    user.Username = input.Username
    user.Name = input.Name
    user.Password = input.Password

    if err := user.SetPassword(user.Password); err != nil {
        c.JSON(http.StatusInternalServerError, utils.Response{Name: "error", Detail: "Failed to encrypt password"})
        return
    }

    if err := utils.CreateUser(&user); err != nil {
        c.JSON(http.StatusInternalServerError, utils.Response{Name: "error", Detail: "Failed to create user"})
        return
    }

    userResponse := UserResponse{
        ID:        user.ID,
        Username:  user.Username,
        Name:      user.Name,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }

    c.JSON(http.StatusCreated, userResponse)
}