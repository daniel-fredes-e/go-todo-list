package main

import (
	"go-todo-list/config"
	_ "go-todo-list/docs"
	"go-todo-list/middleware"
	"go-todo-list/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

    // @title Go Todo List
    // @version 1
    // @Description To Do List

    // @securityDefinitions.apikey
    // @in header
    // @name Authorization

    r := gin.Default()
    config.ConnectDatabase()

    // Rutas públicas
    r.POST("/login", routes.Login)
    r.POST("/register", routes.Register)  // Ruta para registro de usuarios

    // Ruta pública de Swagger
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Rutas protegidas
    authorized := r.Group("/")
    authorized.Use(middleware.JWTMiddleware())
    authorized.GET("/tasks", routes.GetTasks)
    authorized.POST("/tasks", routes.CreateTask)
    authorized.PATCH("/tasks/:id/resolve", routes.MarkTaskResolved)
    authorized.DELETE("/tasks/:id", routes.DeleteTask)

    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Next()
    })
    
    r.Run(":4000")
}
