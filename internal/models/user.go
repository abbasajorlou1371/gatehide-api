package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID          int        `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Mobile      string     `json:"mobile" db:"mobile"`
	Email       string     `json:"email" db:"email"`
	Password    string     `json:"-" db:"password"` // Hidden from JSON
	Image       *string    `json:"image" db:"image"`
	LastLoginAt *time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// Admin represents an admin in the system
type Admin struct {
	ID          int        `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Mobile      string     `json:"mobile" db:"mobile"`
	Email       string     `json:"email" db:"email"`
	Password    string     `json:"-" db:"password"` // Hidden from JSON
	Image       *string    `json:"image" db:"image"`
	LastLoginAt *time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token     string      `json:"token"`
	UserType  string      `json:"user_type"`
	User      interface{} `json:"user"`
	ExpiresAt time.Time   `json:"expires_at"`
}

// UserResponse represents a user response without sensitive data
type UserResponse struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Mobile      string     `json:"mobile"`
	Email       string     `json:"email"`
	Image       *string    `json:"image"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// AdminResponse represents an admin response without sensitive data
type AdminResponse struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Mobile      string     `json:"mobile"`
	Email       string     `json:"email"`
	Image       *string    `json:"image"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Name:        u.Name,
		Mobile:      u.Mobile,
		Email:       u.Email,
		Image:       u.Image,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

// ToResponse converts Admin to AdminResponse
func (a *Admin) ToResponse() AdminResponse {
	return AdminResponse{
		ID:          a.ID,
		Name:        a.Name,
		Mobile:      a.Mobile,
		Email:       a.Email,
		Image:       a.Image,
		LastLoginAt: a.LastLoginAt,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword checks if the provided password matches the hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
