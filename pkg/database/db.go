package database

import (
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pageza/chat-app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// TODO: You might want to add logging for successful database connection.

func init() {
	var err error
	dsn := os.Getenv("POSTGRE_DSN")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DB.AutoMigrate(&models.User{})

}
