package models

import (
	"time"
)

// UserSession represents an active user session
type UserSession struct {
	ID             int       `json:"id" db:"id"`
	UserID         int       `json:"user_id" db:"user_id"`
	UserType       string    `json:"user_type" db:"user_type"`
	SessionToken   string    `json:"session_token" db:"session_token"`
	DeviceInfo     *string   `json:"device_info" db:"device_info"`
	IPAddress      *string   `json:"ip_address" db:"ip_address"`
	UserAgent      *string   `json:"user_agent" db:"user_agent"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	LastActivityAt time.Time `json:"last_activity_at" db:"last_activity_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	ExpiresAt      time.Time `json:"expires_at" db:"expires_at"`
}

// SessionResponse represents a session response without sensitive data
type SessionResponse struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	UserType       string    `json:"user_type"`
	DeviceInfo     *string   `json:"device_info"`
	IPAddress      *string   `json:"ip_address"`
	UserAgent      *string   `json:"user_agent"`
	IsActive       bool      `json:"is_active"`
	LastActivityAt time.Time `json:"last_activity_at"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
	IsCurrent      bool      `json:"is_current"` // This will be set by the service
}

// ToResponse converts UserSession to SessionResponse
func (s *UserSession) ToResponse() SessionResponse {
	return SessionResponse{
		ID:             s.ID,
		UserID:         s.UserID,
		UserType:       s.UserType,
		DeviceInfo:     s.DeviceInfo,
		IPAddress:      s.IPAddress,
		UserAgent:      s.UserAgent,
		IsActive:       s.IsActive,
		LastActivityAt: s.LastActivityAt,
		CreatedAt:      s.CreatedAt,
		ExpiresAt:      s.ExpiresAt,
	}
}

// IsExpired checks if the session is expired
func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsValid checks if the session is valid (active and not expired)
func (s *UserSession) IsValid() bool {
	return s.IsActive && !s.IsExpired()
}

// GetDeviceDisplayName extracts a readable device name from device info
func (s *UserSession) GetDeviceDisplayName() string {
	if s.DeviceInfo == nil {
		return "نامشخص"
	}

	// This is a simple implementation - in production you might want to parse the device info
	// to extract browser, OS, etc.
	return *s.DeviceInfo
}

// GetLocationFromIP extracts location info from IP address
func (s *UserSession) GetLocationFromIP() string {
	if s.IPAddress == nil {
		return "نامشخص"
	}

	// This is a placeholder - in production you might want to use an IP geolocation service
	return "ایران"
}
