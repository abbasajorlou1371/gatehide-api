package services

import (
	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/utils"
)

// AuthServiceInterface defines the interface for authentication services
type AuthServiceInterface interface {
	Login(email, password string, rememberMe bool) (*models.LoginResponse, error)
	ValidateToken(tokenString string) (*utils.JWTClaims, error)
	RefreshToken(tokenString string, rememberMe bool) (string, error)
	GetUserFromToken(tokenString string) (*utils.JWTClaims, error)
}
