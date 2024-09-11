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
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func Login(c *gin.Context) {
    var input UserLogin
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    username := input.Username
    password := input.Password

    user, err := utils.GetUserExist(username)

    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    if !user.CheckPassword(password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": username,
    })

    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Register maneja el registro de nuevos usuarios.
// @Summary Register
// @Description Registra nuevo usuario
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserRegister true "User Registration"
// @Success 201 {object} UserResponse
// @Failure 400 {object} map[string]string
// @Router /register [post]
func Register(c *gin.Context) {
    var user models.User
    var input UserRegister
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Verifica si el nombre de usuario ya existe
    
    if _, err := utils.GetUserExist(input.Username); err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
        return
    }

    user.Username = input.Username
    user.Name = input.Name
    user.Password = input.Password

    // Cifra la contraseña
    if err := user.SetPassword(user.Password); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
        return
    }

    // Guarda el nuevo usuario en la base de datos
    if err := utils.CreateUser(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    // Crear una respuesta sin la contraseña
    userResponse := UserResponse{
        ID:        user.ID,
        Username:  user.Username,
        Name:      user.Name,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }

    c.JSON(http.StatusCreated, userResponse)
}