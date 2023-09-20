// Package database contains functions for initializing and managing the database connection.

package database

import (
	_ "github.com/lib/pq"                        // Import Postgres driver
	"github.com/pageza/chat-app/internal/config" // Import configuration package
	"github.com/pageza/chat-app/internal/models" // Import models package
	"github.com/sirupsen/logrus"                 // Import logging package
	"gorm.io/driver/postgres"                    // Import GORM Postgres driver
	"gorm.io/gorm"                               // Import GORM package
)

// DB is a global variable that holds the database instance
var DB *gorm.DB

// InitializeDB initializes a new database connection using the DSN from the config package.
// It returns the database instance and any error encountered.
func InitializeDB() (*gorm.DB, error) {
	dsn := config.PostgreDSN // Get the DSN from the config package
	// Open a new database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err // Return nil and the error if something goes wrong
	}
	return db, nil // Return the database instance and nil as the error
}

// AutoMigrateDB performs auto-migration for the given database instance.
// It automatically creates the User table if it doesn't exist.
func AutoMigrateDB(db *gorm.DB) error {
	// Perform auto-migration for the User model
	return db.AutoMigrate(&models.User{})
}

// GetDB initializes the database and performs auto-migration.
// It sets the global DB variable.
func GetDB() {
	var err error
	// Initialize the database
	DB, err = InitializeDB()
	if err != nil {
		// Log fatal error and return if database initialization fails
		logrus.Fatalf("Failed to initialize database, got error: %v", err)
		return
	}

	// Perform auto-migration
	if err := AutoMigrateDB(DB); err != nil {
		// Log fatal error and return if auto-migration fails
		logrus.Fatalf("Failed to auto-migrate User model, got error: %v", err)
		return
	}

	// Log successful database connection and auto-migration
	logrus.Info("Successfully connected to the database and auto-migrated the User model.")
}
