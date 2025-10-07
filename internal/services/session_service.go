package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/repositories"
	"github.com/gatehide/gatehide-api/internal/utils"
)

// SessionService implements SessionServiceInterface
type SessionService struct {
	sessionRepo repositories.SessionRepositoryInterface
	jwtManager  *utils.JWTManager
	cfg         *config.Config
}

// NewSessionService creates a new session service
func NewSessionService(sessionRepo repositories.SessionRepositoryInterface, cfg *config.Config) SessionServiceInterface {
	return &SessionService{
		sessionRepo: sessionRepo,
		jwtManager:  utils.NewJWTManager(cfg),
		cfg:         cfg,
	}
}

// CreateSession creates a new user session and returns both session and JWT token
func (s *SessionService) CreateSession(userID int, userType, deviceInfo, ipAddress, userAgent string, rememberMe bool) (*models.UserSession, string, error) {
	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(userID, userType, "", "", rememberMe)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Calculate expiration time
	expiration := time.Duration(s.cfg.Security.JWTExpiration) * time.Hour
	if rememberMe {
		expiration = expiration * 24 * 7 // 7 days for remember me
	}
	expiresAt := time.Now().Add(expiration)

	// Create session in database
	var deviceInfoPtr, ipAddressPtr, userAgentPtr *string
	if deviceInfo != "" {
		deviceInfoPtr = &deviceInfo
	}
	if ipAddress != "" {
		ipAddressPtr = &ipAddress
	}
	if userAgent != "" {
		userAgentPtr = &userAgent
	}

	session, err := s.sessionRepo.CreateSession(userID, userType, token, deviceInfoPtr, ipAddressPtr, userAgentPtr, expiresAt)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create session: %w", err)
	}

	return session, token, nil
}

// ValidateAndUpdateSession validates a session token and updates its activity
func (s *SessionService) ValidateAndUpdateSession(sessionToken string) (*models.UserSession, error) {
	// First validate JWT token
	_, err := s.jwtManager.ValidateToken(sessionToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Get session from database
	session, err := s.sessionRepo.GetSessionByToken(sessionToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		return nil, errors.New("session not found")
	}

	// Check if session is valid
	if !session.IsValid() {
		return nil, errors.New("session is not active or expired")
	}

	// Update session activity
	err = s.sessionRepo.UpdateSessionActivity(session.ID)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to update session activity: %v\n", err)
	}

	// Update the session's last activity time for the response
	session.LastActivityAt = time.Now()

	return session, nil
}

// GetActiveSessions retrieves all active sessions for a user
func (s *SessionService) GetActiveSessions(userID int, userType string, currentSessionToken string) ([]models.SessionResponse, error) {
	sessions, err := s.sessionRepo.GetActiveSessionsByUserID(userID, userType)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}

	var responses []models.SessionResponse
	for _, session := range sessions {
		response := session.ToResponse()
		// Mark current session
		response.IsCurrent = session.SessionToken == currentSessionToken
		responses = append(responses, response)
	}

	return responses, nil
}

// LogoutSession deactivates a specific session
func (s *SessionService) LogoutSession(sessionID int, userID int, userType string) error {
	// Verify the session belongs to the user
	sessions, err := s.sessionRepo.GetActiveSessionsByUserID(userID, userType)
	if err != nil {
		return fmt.Errorf("failed to verify session ownership: %w", err)
	}

	// Check if the session exists and belongs to the user
	found := false
	for _, session := range sessions {
		if session.ID == sessionID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("session not found or does not belong to user")
	}

	err = s.sessionRepo.DeactivateSession(sessionID)
	if err != nil {
		return fmt.Errorf("failed to logout session: %w", err)
	}

	return nil
}

// LogoutAllOtherSessions deactivates all sessions except the current one
func (s *SessionService) LogoutAllOtherSessions(userID int, userType string, currentSessionToken string) error {
	err := s.sessionRepo.DeactivateAllOtherUserSessions(userID, userType, currentSessionToken)
	if err != nil {
		return fmt.Errorf("failed to logout other sessions: %w", err)
	}

	return nil
}

// LogoutAllSessions deactivates all sessions for a user
func (s *SessionService) LogoutAllSessions(userID int, userType string) error {
	err := s.sessionRepo.DeactivateAllUserSessions(userID, userType)
	if err != nil {
		return fmt.Errorf("failed to logout all sessions: %w", err)
	}

	return nil
}

// CleanupExpiredSessions removes expired sessions from the database
func (s *SessionService) CleanupExpiredSessions() error {
	err := s.sessionRepo.CleanupExpiredSessions()
	if err != nil {
		return fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}

	return nil
}
