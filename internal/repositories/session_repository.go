package repositories

import (
	"database/sql"
	"time"

	"github.com/gatehide/gatehide-api/internal/models"
)

// SessionRepositoryInterface defines the interface for session repository
type SessionRepositoryInterface interface {
	CreateSession(userID int, userType, sessionToken string, deviceInfo, ipAddress, userAgent *string, expiresAt time.Time) (*models.UserSession, error)
	GetSessionByToken(sessionToken string) (*models.UserSession, error)
	GetActiveSessionsByUserID(userID int, userType string) ([]models.UserSession, error)
	UpdateSessionActivity(sessionID int) error
	DeactivateSession(sessionID int) error
	DeactivateAllUserSessions(userID int, userType string) error
	DeactivateAllOtherUserSessions(userID int, userType string, currentSessionToken string) error
	CleanupExpiredSessions() error
	DeleteSession(sessionID int) error
}

// SessionRepository implements SessionRepositoryInterface
type SessionRepository struct {
	db *sql.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *sql.DB) SessionRepositoryInterface {
	return &SessionRepository{db: db}
}

// CreateSession creates a new user session
func (r *SessionRepository) CreateSession(userID int, userType, sessionToken string, deviceInfo, ipAddress, userAgent *string, expiresAt time.Time) (*models.UserSession, error) {
	query := `
		INSERT INTO user_sessions (user_id, user_type, session_token, device_info, ip_address, user_agent, expires_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query, userID, userType, sessionToken, deviceInfo, ipAddress, userAgent, expiresAt)
	if err != nil {
		return nil, err
	}

	sessionID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.UserSession{
		ID:           int(sessionID),
		UserID:       userID,
		UserType:     userType,
		SessionToken: sessionToken,
		DeviceInfo:   deviceInfo,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		IsActive:     true,
		ExpiresAt:    expiresAt,
	}, nil
}

// GetSessionByToken retrieves a session by its token
func (r *SessionRepository) GetSessionByToken(sessionToken string) (*models.UserSession, error) {
	query := `
		SELECT id, user_id, user_type, session_token, device_info, ip_address, user_agent, 
		       is_active, last_activity_at, created_at, expires_at
		FROM user_sessions 
		WHERE session_token = ?
	`

	var session models.UserSession
	err := r.db.QueryRow(query, sessionToken).Scan(
		&session.ID,
		&session.UserID,
		&session.UserType,
		&session.SessionToken,
		&session.DeviceInfo,
		&session.IPAddress,
		&session.UserAgent,
		&session.IsActive,
		&session.LastActivityAt,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}

// GetActiveSessionsByUserID retrieves all active sessions for a user
func (r *SessionRepository) GetActiveSessionsByUserID(userID int, userType string) ([]models.UserSession, error) {
	query := `
		SELECT id, user_id, user_type, session_token, device_info, ip_address, user_agent, 
		       is_active, last_activity_at, created_at, expires_at
		FROM user_sessions 
		WHERE user_id = ? AND user_type = ? AND is_active = TRUE AND expires_at > NOW()
		ORDER BY last_activity_at DESC
	`

	rows, err := r.db.Query(query, userID, userType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.UserSession
	for rows.Next() {
		var session models.UserSession
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.UserType,
			&session.SessionToken,
			&session.DeviceInfo,
			&session.IPAddress,
			&session.UserAgent,
			&session.IsActive,
			&session.LastActivityAt,
			&session.CreatedAt,
			&session.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// UpdateSessionActivity updates the last activity time for a session
func (r *SessionRepository) UpdateSessionActivity(sessionID int) error {
	query := `
		UPDATE user_sessions 
		SET last_activity_at = NOW() 
		WHERE id = ? AND is_active = TRUE
	`

	_, err := r.db.Exec(query, sessionID)
	return err
}

// DeactivateSession deactivates a specific session
func (r *SessionRepository) DeactivateSession(sessionID int) error {
	query := `
		UPDATE user_sessions 
		SET is_active = FALSE 
		WHERE id = ?
	`

	_, err := r.db.Exec(query, sessionID)
	return err
}

// DeactivateAllUserSessions deactivates all sessions for a user
func (r *SessionRepository) DeactivateAllUserSessions(userID int, userType string) error {
	query := `
		UPDATE user_sessions 
		SET is_active = FALSE 
		WHERE user_id = ? AND user_type = ?
	`

	_, err := r.db.Exec(query, userID, userType)
	return err
}

// DeactivateAllOtherUserSessions deactivates all sessions for a user except the current one
func (r *SessionRepository) DeactivateAllOtherUserSessions(userID int, userType string, currentSessionToken string) error {
	query := `
		UPDATE user_sessions 
		SET is_active = FALSE 
		WHERE user_id = ? AND user_type = ? AND session_token != ?
	`

	_, err := r.db.Exec(query, userID, userType, currentSessionToken)
	return err
}

// CleanupExpiredSessions removes expired sessions from the database
func (r *SessionRepository) CleanupExpiredSessions() error {
	query := `
		DELETE FROM user_sessions 
		WHERE expires_at < NOW()
	`

	_, err := r.db.Exec(query)
	return err
}

// DeleteSession permanently deletes a session
func (r *SessionRepository) DeleteSession(sessionID int) error {
	query := `DELETE FROM user_sessions WHERE id = ?`

	_, err := r.db.Exec(query, sessionID)
	return err
}
