package services

import (
	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/utils"
)

// AuthServiceInterface defines the interface for authentication services
type AuthServiceInterface interface {
	Login(email, password string, rememberMe bool) (*models.LoginResponse, error)
	LoginWithSession(email, password string, rememberMe bool, deviceInfo, ipAddress, userAgent string) (*models.LoginResponse, error)
	ValidateToken(tokenString string) (*utils.JWTClaims, error)
	RefreshToken(tokenString string, rememberMe bool) (string, error)
	GetUserFromToken(tokenString string) (*utils.JWTClaims, error)
	ForgotPassword(email string) error
	ResetPassword(token, email, newPassword, confirmPassword string) error
	ValidateResetToken(token string) error
	ChangePassword(userID int, userType, currentPassword, newPassword, confirmPassword string) error
}
