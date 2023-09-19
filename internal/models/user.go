package models

import (
	"errors"
	"regexp"
	"time"
)

// User represents a user in the system.
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Validate checks if the User fields are valid.
func (u *User) Validate() error {
	if len(u.Username) < 3 || len(u.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	// Validate email format using regex
	var emailRe = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRe.MatchString(u.Email) {
		return errors.New("invalid email format")
	}

	// Validate password complexity
	if len(u.Password) < 8 || len(u.Password) > 50 {
		return errors.New("password must be between 8 and 50 characters")
	}

	var (
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString
		hasDigit   = regexp.MustCompile(`\d`).MatchString
		hasSpecial = regexp.MustCompile(`[@$!%*?&]`).MatchString
	)

	if !hasUpper(u.Password) || !hasLower(u.Password) || !hasDigit(u.Password) || !hasSpecial(u.Password) {
		return errors.New("password must include at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	return nil
}
