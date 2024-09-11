package middleware

import (
	"go-todo-list/utils"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("your_secret_key")

// JWTMiddleware verifica el token JWT en las solicitudes entrantes.
func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, utils.Response{Name: "error", Detail: "Token is missing"})
            c.Abort()
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, utils.Response{Name: "error", Detail: "Invalid token"})
            c.Abort()
            return
        }

        c.Next()
    }
}
