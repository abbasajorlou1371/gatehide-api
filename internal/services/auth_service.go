package services

import (
	"fmt"
	"time"

	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/repositories"
	"github.com/gatehide/gatehide-api/internal/utils"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo   repositories.UserRepository
	adminRepo  repositories.AdminRepository
	jwtManager *utils.JWTManager
	config     *config.Config
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo repositories.UserRepository,
	adminRepo repositories.AdminRepository,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		adminRepo:  adminRepo,
		jwtManager: utils.NewJWTManager(cfg),
		config:     cfg,
	}
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*utils.JWTClaims, error) {
	return s.jwtManager.ValidateToken(tokenString)
}

// RefreshToken generates a new token with extended expiration
func (s *AuthService) RefreshToken(tokenString string, rememberMe bool) (string, error) {
	return s.jwtManager.RefreshToken(tokenString, rememberMe)
}

// Login unified authentication that determines user type by email
func (s *AuthService) Login(email, password string, rememberMe bool) (*models.LoginResponse, error) {
	// First, try to find the user as a regular user
	user, userErr := s.userRepo.GetByEmail(email)
	if userErr == nil {
		// Verify password for user
		if models.CheckPassword(password, user.Password) {
			// Generate JWT token for user
			token, err := s.jwtManager.GenerateToken(user.ID, "user", user.Email, user.Name, rememberMe)
			if err != nil {
				return nil, fmt.Errorf("failed to generate token: %w", err)
			}

			// Update last login
			if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
				fmt.Printf("Warning: failed to update last login for user %d: %v\n", user.ID, err)
			}

			// Calculate token expiration
			expiresAt := time.Now().Add(time.Duration(s.config.Security.JWTExpiration) * time.Hour)

			return &models.LoginResponse{
				Token:     token,
				UserType:  "user",
				User:      user.ToResponse(),
				ExpiresAt: expiresAt,
			}, nil
		}
	}

	// If user login failed, try admin login
	admin, adminErr := s.adminRepo.GetByEmail(email)
	if adminErr == nil {
		// Verify password for admin
		if models.CheckPassword(password, admin.Password) {
			// Generate JWT token for admin
			token, err := s.jwtManager.GenerateToken(admin.ID, "admin", admin.Email, admin.Name, rememberMe)
			if err != nil {
				return nil, fmt.Errorf("failed to generate token: %w", err)
			}

			// Update last login
			if err := s.adminRepo.UpdateLastLogin(admin.ID); err != nil {
				fmt.Printf("Warning: failed to update last login for admin %d: %v\n", admin.ID, err)
			}

			// Calculate token expiration
			expiresAt := time.Now().Add(time.Duration(s.config.Security.JWTExpiration) * time.Hour)

			return &models.LoginResponse{
				Token:     token,
				UserType:  "admin",
				User:      admin.ToResponse(),
				ExpiresAt: expiresAt,
			}, nil
		}
	}

	// If both failed, return invalid credentials error
	return nil, fmt.Errorf("invalid credentials")
}

// GetUserFromToken extracts user information from a JWT token
func (s *AuthService) GetUserFromToken(tokenString string) (*utils.JWTClaims, error) {
	claims, err := s.jwtManager.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}
