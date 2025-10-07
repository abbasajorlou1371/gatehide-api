package repositories

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
)

// EmailVerificationRepository handles email verification code operations
type EmailVerificationRepository struct {
	db *sql.DB
}

// NewEmailVerificationRepository creates a new email verification repository
func NewEmailVerificationRepository(db *sql.DB) *EmailVerificationRepository {
	return &EmailVerificationRepository{db: db}
}

// hashCode hashes a verification code using SHA-256
func (r *EmailVerificationRepository) hashCode(code string) string {
	hash := sha256.Sum256([]byte(code))
	return hex.EncodeToString(hash[:])
}

// StoreCode stores an email verification code with expiration
func (r *EmailVerificationRepository) StoreCode(userID int, userType, email, code string, expiresAt time.Time) error {
	// First, delete any existing codes for this user and email
	if err := r.DeleteUserCodes(userID, userType, email); err != nil {
		return fmt.Errorf("failed to delete existing codes: %w", err)
	}

	// Hash the code before storing
	hashedCode := r.hashCode(code)

	// Insert new hashed code
	query := `
		INSERT INTO email_verification_codes (user_id, user_type, email, code, expires_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, userID, userType, email, hashedCode, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to store verification code: %w", err)
	}

	return nil
}

// VerifyCode verifies an email verification code
func (r *EmailVerificationRepository) VerifyCode(userID int, userType, email, code string) (bool, error) {
	// Hash the input code for comparison
	hashedCode := r.hashCode(code)

	query := `
		SELECT id, expires_at 
		FROM email_verification_codes 
		WHERE user_id = ? AND user_type = ? AND email = ? AND code = ?
		ORDER BY created_at DESC 
		LIMIT 1
	`

	var id int
	var expiresAt time.Time

	err := r.db.QueryRow(query, userID, userType, email, hashedCode).Scan(&id, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Code not found
		}
		return false, fmt.Errorf("failed to query verification code: %w", err)
	}

	// Check if code has expired
	if time.Now().After(expiresAt) {
		// Delete expired code
		r.DeleteCode(id)
		return false, nil
	}

	// Code is valid, delete it (one-time use)
	if err := r.DeleteCode(id); err != nil {
		return false, fmt.Errorf("failed to delete used code: %w", err)
	}

	return true, nil
}

// DeleteCode deletes a specific verification code by ID
func (r *EmailVerificationRepository) DeleteCode(id int) error {
	query := `DELETE FROM email_verification_codes WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete verification code: %w", err)
	}
	return nil
}

// DeleteUserCodes deletes all verification codes for a specific user and email
func (r *EmailVerificationRepository) DeleteUserCodes(userID int, userType, email string) error {
	query := `DELETE FROM email_verification_codes WHERE user_id = ? AND user_type = ? AND email = ?`
	_, err := r.db.Exec(query, userID, userType, email)
	if err != nil {
		return fmt.Errorf("failed to delete user verification codes: %w", err)
	}
	return nil
}

// CleanupExpiredCodes removes all expired verification codes
func (r *EmailVerificationRepository) CleanupExpiredCodes() error {
	query := `DELETE FROM email_verification_codes WHERE expires_at < NOW()`
	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired codes: %w", err)
	}
	return nil
}
