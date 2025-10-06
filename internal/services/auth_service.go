package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/repositories"
	"github.com/gatehide/gatehide-api/internal/utils"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo            repositories.UserRepository
	adminRepo           repositories.AdminRepository
	passwordResetRepo   repositories.PasswordResetRepositoryInterface
	notificationService NotificationServiceInterface
	jwtManager          *utils.JWTManager
	config              *config.Config
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo repositories.UserRepository,
	adminRepo repositories.AdminRepository,
	passwordResetRepo repositories.PasswordResetRepositoryInterface,
	notificationService NotificationServiceInterface,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		userRepo:            userRepo,
		adminRepo:           adminRepo,
		passwordResetRepo:   passwordResetRepo,
		notificationService: notificationService,
		jwtManager:          utils.NewJWTManager(cfg),
		config:              cfg,
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

// GetJWTManager returns the JWT manager instance
func (s *AuthService) GetJWTManager() *utils.JWTManager {
	return s.jwtManager
}

// generateResetToken generates a secure random token for password reset
func (s *AuthService) generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// ForgotPassword initiates the password reset process
func (s *AuthService) ForgotPassword(email string) error {
	// First, try to find the user as a regular user
	user, userErr := s.userRepo.GetByEmail(email)
	if userErr == nil {
		// Invalidate any existing tokens for this user
		if err := s.passwordResetRepo.InvalidateUserTokens(user.ID, "user"); err != nil {
			fmt.Printf("Warning: failed to invalidate existing tokens for user %d: %v\n", user.ID, err)
		}

		// Generate new reset token
		token, err := s.generateResetToken()
		if err != nil {
			return fmt.Errorf("failed to generate reset token: %w", err)
		}

		// Set token expiration (15 minutes from now)
		expiresAt := time.Now().Add(15 * time.Minute)

		// Create the token in database
		if err := s.passwordResetRepo.CreateToken(user.ID, "user", token, expiresAt); err != nil {
			return fmt.Errorf("failed to create reset token: %w", err)
		}

		// Send password reset email
		if err := s.sendPasswordResetEmail(user.Email, user.Name, token); err != nil {
			fmt.Printf("Warning: failed to send password reset email to %s: %v\n", email, err)
			// Don't return error here, as the token was created successfully
		}

		return nil
	}

	// If user not found, try admin
	admin, adminErr := s.adminRepo.GetByEmail(email)
	if adminErr == nil {
		// Invalidate any existing tokens for this admin
		if err := s.passwordResetRepo.InvalidateUserTokens(admin.ID, "admin"); err != nil {
			fmt.Printf("Warning: failed to invalidate existing tokens for admin %d: %v\n", admin.ID, err)
		}

		// Generate new reset token
		token, err := s.generateResetToken()
		if err != nil {
			return fmt.Errorf("failed to generate reset token: %w", err)
		}

		// Set token expiration (15 minutes from now)
		expiresAt := time.Now().Add(15 * time.Minute)

		// Create the token in database
		if err := s.passwordResetRepo.CreateToken(admin.ID, "admin", token, expiresAt); err != nil {
			return fmt.Errorf("failed to create reset token: %w", err)
		}

		// Send password reset email
		if err := s.sendPasswordResetEmail(admin.Email, admin.Name, token); err != nil {
			fmt.Printf("Warning: failed to send password reset email to %s: %v\n", email, err)
			// Don't return error here, as the token was created successfully
		}

		return nil
	}

	// If neither user nor admin found, return error
	return fmt.Errorf("email not found")
}

// ResetPassword resets the password using a valid token
func (s *AuthService) ResetPassword(token, email, newPassword, confirmPassword string) error {
	// Validate passwords match
	if newPassword != confirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	// Validate password strength
	if len(newPassword) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}

	// Get the token from database
	resetToken, err := s.passwordResetRepo.GetTokenByToken(token)
	if err != nil {
		return fmt.Errorf("invalid or expired token")
	}

	// Check if token is valid (not expired and not used)
	if !resetToken.IsValid() {
		return fmt.Errorf("invalid or expired token")
	}

	// Validate that the email matches the token's user
	switch resetToken.UserType {
	case "user":
		user, err := s.userRepo.GetByEmail(email)
		if err != nil || user.ID != resetToken.UserID {
			return fmt.Errorf("invalid email for this token")
		}
	case "admin":
		admin, err := s.adminRepo.GetByEmail(email)
		if err != nil || admin.ID != resetToken.UserID {
			return fmt.Errorf("invalid email for this token")
		}
	}

	// Hash the new password
	hashedPassword, err := models.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password based on user type
	switch resetToken.UserType {
	case "user":
		if err := s.userRepo.UpdatePassword(resetToken.UserID, hashedPassword); err != nil {
			return fmt.Errorf("failed to update user password: %w", err)
		}
	case "admin":
		if err := s.adminRepo.UpdatePassword(resetToken.UserID, hashedPassword); err != nil {
			return fmt.Errorf("failed to update admin password: %w", err)
		}
	default:
		return fmt.Errorf("invalid user type")
	}

	// Mark token as used
	if err := s.passwordResetRepo.MarkTokenAsUsed(token); err != nil {
		fmt.Printf("Warning: failed to mark token as used: %v\n", err)
	}

	// Invalidate all other tokens for this user
	if err := s.passwordResetRepo.InvalidateUserTokens(resetToken.UserID, resetToken.UserType); err != nil {
		fmt.Printf("Warning: failed to invalidate other tokens: %v\n", err)
	}

	return nil
}

// ValidateResetToken validates a password reset token
func (s *AuthService) ValidateResetToken(token string) error {
	resetToken, err := s.passwordResetRepo.GetTokenByToken(token)
	if err != nil {
		return fmt.Errorf("invalid token")
	}

	if !resetToken.IsValid() {
		return fmt.Errorf("token is expired or already used")
	}

	return nil
}

// sendPasswordResetEmail sends a password reset email using the notification service
func (s *AuthService) sendPasswordResetEmail(email, name, token string) error {
	if s.notificationService == nil {
		return fmt.Errorf("notification service not available")
	}

	// Create reset link with email parameter
	resetLink := fmt.Sprintf("http://localhost:3000/reset-password?token=%s&email=%s", token, email)
	unsubscribeLink := "http://localhost:3000/unsubscribe?email=" + email
	supportLink := "http://localhost:3000/support"

	// Create notification request
	notification := &models.CreateNotificationRequest{
		Type:      models.NotificationTypeEmail,
		Priority:  models.NotificationPriorityHigh,
		Recipient: email,
		Subject:   fmt.Sprintf("بازنشانی رمز عبور - %s", s.config.App.Name),
		Content:   fmt.Sprintf("کاربر گرامی %s،\n\nدرخواست بازنشانی رمز عبور برای حساب کاربری شما در %s دریافت شده است.\n\nبرای تنظیم رمز عبور جدید، لطفاً روی لینک زیر کلیک کنید:\n%s\n\nاین لینک تا 0.25 ساعت معتبر است.\n\nاگر شما این درخواست را انجام نداده\u200cاید، لطفاً این ایمیل را نادیده بگیرید.\n\nبا احترام،\nتیم %s", name, s.config.App.Name, resetLink, s.config.App.Name),
		TemplateData: map[string]interface{}{
			"app_name":         s.config.App.Name,
			"user_name":        name,
			"reset_link":       resetLink,
			"expiry_hours":     "0.25", // 15 minutes
			"unsubscribe_link": unsubscribeLink,
			"support_link":     supportLink,
		},
	}

	// Send the notification
	ctx := context.Background()
	return s.notificationService.SendNotification(ctx, notification)
}
