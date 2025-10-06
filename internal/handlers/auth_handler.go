package handlers

import (
	"net/http"

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

	newToken, err := h.authService.RefreshToken(tokenString)
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

// Logout handles logout requests (client-side token removal)
func (h *AuthHandler) Logout(c *gin.Context) {
	// Since we're using stateless JWT tokens, logout is handled client-side
	// by removing the token from storage. This endpoint just confirms the logout.
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
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

	response, err := h.authService.Login(req.Email, req.Password)
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
