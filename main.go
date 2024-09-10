package main

import (
	"go-todo-list/config"
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

    // @securityDefinitions.apikey barerToken
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

    r.Run(":4000")
}
