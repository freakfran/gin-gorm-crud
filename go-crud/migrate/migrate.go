package main

import (
	"github.com/gookit/slog"
	"go-crud/initializers"
	"go-crud/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	err := initializers.DB.AutoMigrate(&models.Post{})
	if err != nil {
		slog.Fatalf("Failed to migrate database:%v", err)
	}
}
