package models

import "time"

// EmailVerificationCode represents an email verification code
type EmailVerificationCode struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	UserType  string    `json:"user_type" db:"user_type"`
	Email     string    `json:"email" db:"email"`
	Code      string    `json:"code" db:"code"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateEmailVerificationCodeRequest represents a request to create a verification code
type CreateEmailVerificationCodeRequest struct {
	UserID   int    `json:"user_id" binding:"required"`
	UserType string `json:"user_type" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Code     string `json:"code" binding:"required"`
}

// VerifyEmailCodeRequest represents a request to verify an email code
type VerifyEmailCodeRequest struct {
	UserID   int    `json:"user_id" binding:"required"`
	UserType string `json:"user_type" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Code     string `json:"code" binding:"required"`
}
