package database

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/pageza/chat-app/internal/config"
	"github.com/pageza/chat-app/internal/models"
)

var DB *gorm.DB

func InitializeDB() (*gorm.DB, error) {
	dsn := config.PostgreDSN
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func AutoMigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}

func GetDB() {
	var err error
	DB, err = InitializeDB()
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
		return
	}

	if err := AutoMigrateDB(DB); err != nil {
		logrus.Fatalf("Failed to auto-migrate User model: %v", err)
		return
	}

	logrus.Info("Successfully connected to the database and auto-migrated the User model.")
}
