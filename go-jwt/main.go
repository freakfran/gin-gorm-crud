package main

import (
	"github.com/gin-gonic/gin"
	"go-jwt/controllers"
	"go-jwt/initializers"
	"go-jwt/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	route := gin.Default()

	route.POST("/signUp", controllers.SingUp)
	route.POST("/login", controllers.Login)
	route.GET("/validate", middleware.RequireAuth, controllers.Validate)
	route.Run()
}
