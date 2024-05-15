package initializers

import (
	"github.com/gookit/slog"
	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		slog.Fatal("Error loading .env file")
	}
}
