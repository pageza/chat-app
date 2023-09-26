package database

import (
	"fmt"

	"github.com/pageza/chat-app/internal/config"
	"github.com/pageza/chat-app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	InitializeDB() (*gorm.DB, error)
	AutoMigrateDB() error
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	// Additional methods as needed
	UpdateLastLoginTime(user *models.User) error      // To update the last login time of the user
	HandleFailedLoginAttempt(user *models.User) error // To handle failed login attempts, if needed
	// Any other methods related to login/logout that interact with the database
}

type GormDatabase struct {
	DB *gorm.DB
}

func NewGormDatabase() (*GormDatabase, error) {
	db, err := gorm.Open(postgres.Open(config.PostgreDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate User model: %w", err)
	}
	return &GormDatabase{DB: db}, nil
}

func (g *GormDatabase) InitializeDB() (*gorm.DB, error) {
	// Perform any needed initialization logic here
	return g.DB, nil
}

func (g *GormDatabase) AutoMigrateDB() error {
	return g.DB.AutoMigrate(&models.User{})
}

func (g *GormDatabase) CreateUser(user *models.User) error {
	return g.DB.Create(user).Error
}

func (g *GormDatabase) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := g.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
