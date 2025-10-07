package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/services"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService services.AuthServiceInterface
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService services.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RefreshToken handles token refresh requests
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization header required",
		})
		return
	}

	// Extract token from "Bearer <token>" format
	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	// Get remember me preference from request body (optional)
	var req struct {
		RememberMe bool `json:"remember_me"`
	}
	c.ShouldBindJSON(&req) // Ignore errors, default to false

	newToken, err := h.authService.RefreshToken(tokenString, req.RememberMe)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or expired token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"data": gin.H{
			"token": newToken,
		},
	})
}

// Logout handles logout requests with token validation and logging
func (h *AuthHandler) Logout(c *gin.Context) {
	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Authorization header required",
		})
		return
	}

	// Extract token from "Bearer <token>" format
	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	// Validate token to get user information for logging
	claims, err := h.authService.ValidateToken(tokenString)
	if err != nil {
		// Even if token is invalid, we should still allow logout
		// This handles cases where token expired but user wants to logout
		c.JSON(http.StatusOK, gin.H{
			"message": "Logout successful",
		})
		return
	}

	// Log the logout event for security auditing
	fmt.Printf("User logout: ID=%d, Email=%s, Type=%s, Time=%s\n",
		claims.UserID, claims.Email, claims.UserType, time.Now().Format(time.RFC3339))

	// Since we're using stateless JWT tokens, logout is handled client-side
	// by removing the token from storage. This endpoint confirms the logout.
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
		"data": gin.H{
			"user_id":   claims.UserID,
			"user_type": claims.UserType,
		},
	})
}

// Login handles unified login requests (automatically determines user type)
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Extract device information from request headers
	deviceInfo := c.GetHeader("X-Device-Info")
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Use LoginWithSession to create a session during login
	response, err := h.authService.LoginWithSession(req.Email, req.Password, req.RememberMe, deviceInfo, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    response,
	})
}

// GetProfile returns the current user's profile information
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user info from context (set by middleware)
	userInfo, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User information not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"data":    userInfo,
	})
}

// ForgotPassword handles forgot password requests
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	err := h.authService.ForgotPassword(req.Email)
	if err != nil {
		if err.Error() == "email not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Email not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process password reset request",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset email sent successfully",
	})
}

// ResetPassword handles password reset requests
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	err := h.authService.ResetPassword(req.Token, req.Email, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		if err.Error() == "invalid or expired token" || err.Error() == "token is expired or already used" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}
		if err.Error() == "passwords do not match" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Passwords do not match",
			})
			return
		}
		if err.Error() == "password must be at least 6 characters long" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Password must be at least 6 characters long",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to reset password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}

// ValidateResetToken validates a password reset token
func (h *AuthHandler) ValidateResetToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token is required",
		})
		return
	}

	err := h.authService.ValidateResetToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid or expired token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token is valid",
	})
}

// ChangePassword handles change password requests for authenticated users
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "داده‌های درخواست نامعتبر است",
			"details": err.Error(),
		})
		return
	}

	// Extract user information from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "شناسه کاربر یافت نشد",
		})
		return
	}

	userType, exists := c.Get("user_type")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "نوع کاربر یافت نشد",
		})
		return
	}

	// Call the service to change password
	err := h.authService.ChangePassword(
		userID.(int),
		userType.(string),
		req.CurrentPassword,
		req.NewPassword,
		req.ConfirmPassword,
	)

	if err != nil {
		if err.Error() == "رمز عبور فعلی اشتباه است" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "رمز عبور فعلی اشتباه است",
			})
			return
		}
		if err.Error() == "رمز عبور جدید و تأیید رمز عبور مطابقت ندارند" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "رمز عبور جدید و تأیید رمز عبور مطابقت ندارند",
			})
			return
		}
		if err.Error() == "رمز عبور باید حداقل 6 کاراکتر باشد" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "رمز عبور باید حداقل 6 کاراکتر باشد",
			})
			return
		}
		if err.Error() == "کاربر یافت نشد" || err.Error() == "مدیر یافت نشد" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "کاربر یافت نشد",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "خطا در تغییر رمز عبور",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "رمز عبور با موفقیت تغییر یافت",
		"data": gin.H{
			"user_id":   userID,
			"user_type": userType,
		},
	})
}
