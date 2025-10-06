package services

import (
	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/utils"
)

// AuthServiceInterface defines the interface for authentication services
type AuthServiceInterface interface {
	Login(email, password string) (*models.LoginResponse, error)
	ValidateToken(tokenString string) (*utils.JWTClaims, error)
	RefreshToken(tokenString string) (string, error)
	GetUserFromToken(tokenString string) (*utils.JWTClaims, error)
}
