package services

import (
	"github.com/gatehide/gatehide-api/internal/models"
)

// SessionServiceInterface defines the interface for session service
type SessionServiceInterface interface {
	CreateSession(userID int, userType, deviceInfo, ipAddress, userAgent string, rememberMe bool) (*models.UserSession, string, error)
	ValidateAndUpdateSession(sessionToken string) (*models.UserSession, error)
	GetActiveSessions(userID int, userType string, currentSessionToken string) ([]models.SessionResponse, error)
	LogoutSession(sessionID int, userID int, userType string) error
	LogoutAllOtherSessions(userID int, userType string, currentSessionToken string) error
	LogoutAllSessions(userID int, userType string) error
	CleanupExpiredSessions() error
}
