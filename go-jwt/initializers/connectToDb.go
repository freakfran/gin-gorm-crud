package initializers

import (
	"github.com/gookit/slog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func ConnectToDb() {
	var err error
	// 连接数据库
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Fatal("failed to connect database:", err)
	}
	slog.Info("Connected to database")
}
