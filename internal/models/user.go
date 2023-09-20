// Package models defines the data structures used in the application.
// This file specifically includes the User model and its validation logic.

package models

import (
	"errors"
	"regexp"
	"time"
)

// User represents a user in the system. It includes fields for the user's ID, username, email, and password.
// It also includes timestamps for when the user was created and last updated.
type User struct {
	ID        uint      `gorm:"primaryKey"`      // Primary key for the user
	Username  string    `gorm:"unique;not null"` // Unique username, cannot be null
	Email     string    `gorm:"unique;not null"` // Unique email, cannot be null
	Password  string    `gorm:"not null"`        // Password, cannot be null
	CreatedAt time.Time // Timestamp for when the user was created
	UpdatedAt time.Time // Timestamp for when the user was last updated
}

// Validate checks if the User fields are valid.
// It validates the length and format of the username, email, and password.
func (u *User) Validate() error {
	// Validate username length
	if len(u.Username) < 3 || len(u.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	// Validate email format using regex
	var emailRe = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRe.MatchString(u.Email) {
		return errors.New("invalid email format")
	}

	// Validate password length
	if len(u.Password) < 8 || len(u.Password) > 50 {
		return errors.New("password must be between 8 and 50 characters")
	}

	// Validate password complexity using regex
	var (
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString     // At least one uppercase letter
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString     // At least one lowercase letter
		hasDigit   = regexp.MustCompile(`\d`).MatchString        // At least one digit
		hasSpecial = regexp.MustCompile(`[@$!%*?&]`).MatchString // At least one special character
	)

	if !hasUpper(u.Password) || !hasLower(u.Password) || !hasDigit(u.Password) || !hasSpecial(u.Password) {
		return errors.New("password must include at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	return nil
}
