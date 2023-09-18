package database

import (
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pageza/chat-app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	dsn := os.Getenv("POSTGRE_DSN")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DB.AutoMigrate(&models.User{})

}
