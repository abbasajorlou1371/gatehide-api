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
	Balance     float64    `json:"balance" db:"balance"`
	Debt        float64    `json:"debt" db:"debt"`
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
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	RememberMe bool   `json:"remember_me"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token       string      `json:"token"`
	UserType    string      `json:"user_type"`
	User        interface{} `json:"user"`
	Permissions []string    `json:"permissions"`
	ExpiresAt   time.Time   `json:"expires_at"`
}

// UserResponse represents a user response without sensitive data
type UserResponse struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Mobile      string     `json:"mobile"`
	Email       string     `json:"email"`
	Image       *string    `json:"image"`
	Balance     float64    `json:"balance"`
	Debt        float64    `json:"debt"`
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

// ProfileResponse represents a profile response with permissions
type ProfileResponse struct {
	User        interface{} `json:"user"`
	UserType    string      `json:"user_type"`
	Permissions []string    `json:"permissions"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Name:        u.Name,
		Mobile:      u.Mobile,
		Email:       u.Email,
		Image:       u.Image,
		Balance:     u.Balance,
		Debt:        u.Debt,
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

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        int        `json:"id" db:"id"`
	UserID    int        `json:"user_id" db:"user_id"`
	UserType  string     `json:"user_type" db:"user_type"`
	Token     string     `json:"token" db:"token"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	UsedAt    *time.Time `json:"used_at" db:"used_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// ForgotPasswordRequest represents a forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents a reset password request
type ResetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6"`
}

// ChangePasswordRequest represents a change password request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6"`
}

// IsExpired checks if the token is expired
func (prt *PasswordResetToken) IsExpired() bool {
	return time.Now().After(prt.ExpiresAt)
}

// IsUsed checks if the token has been used
func (prt *PasswordResetToken) IsUsed() bool {
	return prt.UsedAt != nil
}

// IsValid checks if the token is valid (not expired and not used)
func (prt *PasswordResetToken) IsValid() bool {
	return !prt.IsExpired() && !prt.IsUsed()
}

// UserCreateRequest represents a request to create a new user
type UserCreateRequest struct {
	Name   string `json:"name" binding:"required,min=2"`
	Email  string `json:"email" binding:"required,email"`
	Mobile string `json:"mobile" binding:"required,min=11,max=11"`
}

// UserUpdateRequest represents a request to update a user
type UserUpdateRequest struct {
	Name   *string `json:"name,omitempty"`
	Email  *string `json:"email,omitempty"`
	Mobile *string `json:"mobile,omitempty"`
	Image  *string `json:"image,omitempty"`
}

// UserSearchRequest represents a search request for users
type UserSearchRequest struct {
	Query    string `json:"query"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

// UserSearchResponse represents a search response for users
type UserSearchResponse struct {
	Data       []UserResponse `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}
