package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gatehide/gatehide-api/internal/models"
)

// PasswordResetRepositoryInterface defines the interface for password reset operations
type PasswordResetRepositoryInterface interface {
	CreateToken(userID int, userType, token string, expiresAt time.Time) error
	GetTokenByToken(token string) (*models.PasswordResetToken, error)
	MarkTokenAsUsed(token string) error
	CleanupExpiredTokens() error
	GetActiveTokensForUser(userID int, userType string) ([]*models.PasswordResetToken, error)
	InvalidateUserTokens(userID int, userType string) error
}

// PasswordResetRepository handles password reset token operations
type PasswordResetRepository struct {
	db *sql.DB
}

// NewPasswordResetRepository creates a new password reset repository
func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
	return &PasswordResetRepository{
		db: db,
	}
}

// CreateToken creates a new password reset token
func (r *PasswordResetRepository) CreateToken(userID int, userType, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO password_reset_tokens (user_id, user_type, token, expires_at)
		VALUES (?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, userID, userType, token, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to create password reset token: %w", err)
	}

	return nil
}

// GetTokenByToken retrieves a password reset token by token string
func (r *PasswordResetRepository) GetTokenByToken(token string) (*models.PasswordResetToken, error) {
	query := `
		SELECT id, user_id, user_type, token, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE token = ?
	`

	var prt models.PasswordResetToken
	err := r.db.QueryRow(query, token).Scan(
		&prt.ID,
		&prt.UserID,
		&prt.UserType,
		&prt.Token,
		&prt.ExpiresAt,
		&prt.UsedAt,
		&prt.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("password reset token not found")
		}
		return nil, fmt.Errorf("failed to get password reset token: %w", err)
	}

	return &prt, nil
}

// MarkTokenAsUsed marks a password reset token as used
func (r *PasswordResetRepository) MarkTokenAsUsed(token string) error {
	query := `
		UPDATE password_reset_tokens
		SET used_at = ?
		WHERE token = ? AND used_at IS NULL
	`

	result, err := r.db.Exec(query, time.Now(), token)
	if err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("token not found or already used")
	}

	return nil
}

// CleanupExpiredTokens removes expired tokens from the database
func (r *PasswordResetRepository) CleanupExpiredTokens() error {
	query := `
		DELETE FROM password_reset_tokens
		WHERE expires_at < ?
	`

	_, err := r.db.Exec(query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}

	return nil
}

// GetActiveTokensForUser gets all active (non-expired, non-used) tokens for a user
func (r *PasswordResetRepository) GetActiveTokensForUser(userID int, userType string) ([]*models.PasswordResetToken, error) {
	query := `
		SELECT id, user_id, user_type, token, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE user_id = ? AND user_type = ? AND expires_at > ? AND used_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID, userType, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to get active tokens: %w", err)
	}
	defer rows.Close()

	var tokens []*models.PasswordResetToken
	for rows.Next() {
		var prt models.PasswordResetToken
		err := rows.Scan(
			&prt.ID,
			&prt.UserID,
			&prt.UserType,
			&prt.Token,
			&prt.ExpiresAt,
			&prt.UsedAt,
			&prt.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan token: %w", err)
		}
		tokens = append(tokens, &prt)
	}

	return tokens, nil
}

// InvalidateUserTokens invalidates all active tokens for a user
func (r *PasswordResetRepository) InvalidateUserTokens(userID int, userType string) error {
	query := `
		UPDATE password_reset_tokens
		SET used_at = ?
		WHERE user_id = ? AND user_type = ? AND used_at IS NULL
	`

	_, err := r.db.Exec(query, time.Now(), userID, userType)
	if err != nil {
		return fmt.Errorf("failed to invalidate user tokens: %w", err)
	}

	return nil
}
