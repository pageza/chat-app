package database

import (
	_ "github.com/lib/pq"
	"github.com/pageza/chat-app/internal/config"
	"github.com/pageza/chat-app/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitializeDB sets up and returns a database instance
func InitializeDB() (*gorm.DB, error) {
	dsn := config.PostgreDSN // Get DSN from config package
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// AutoMigrateDB performs auto-migration for the given database instance
func AutoMigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}

// GetDB initializes the database and performs auto-migration
func GetDB() {
	var err error
	DB, err = InitializeDB()
	if err != nil {
		logrus.Fatalf("Failed to initialize database, got error: %v", err)
		return
	}

	if err := AutoMigrateDB(DB); err != nil {
		logrus.Fatalf("Failed to auto-migrate User model, got error: %v", err)
		return
	}

	logrus.Info("Successfully connected to the database and auto-migrated the User model.")
}
