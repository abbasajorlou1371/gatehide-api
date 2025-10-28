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
	GetUserByID(userID int) (*models.User, error)
	GetAdminByID(adminID int) (*models.Admin, error)
	GetGamenetByID(gamenetID int) (*models.Gamenet, error)
	UpdateUserProfile(userID int, name, mobile, image string) (*models.UserResponse, error)
	UpdateAdminProfile(adminID int, name, mobile, image string) (*models.AdminResponse, error)
	UpdateGamenetProfile(gamenetID int, name, mobile, image string) (*models.GamenetResponse, error)
	UpdateUserEmail(userID int, newEmail string) (*models.UserResponse, error)
	UpdateAdminEmail(adminID int, newEmail string) (*models.AdminResponse, error)
	UpdateGamenetEmail(gamenetID int, newEmail string) (*models.GamenetResponse, error)
	ForgotPassword(email string) error
	ResetPassword(token, email, newPassword, confirmPassword string) error
	ValidateResetToken(token string) error
	ChangePassword(userID int, userType, currentPassword, newPassword, confirmPassword string) error
	SendEmailVerification(userID int, userType, newEmail string) (string, error)
	VerifyEmailCode(userID int, userType, email, code string) (bool, error)
	CheckEmailExists(email string) (bool, error)
	GetUserPermissions(userType string) ([]string, error)
	GetUserPermissionsByID(userID int, userType string) ([]string, error)
}
