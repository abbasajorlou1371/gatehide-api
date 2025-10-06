package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gatehide/gatehide-api/internal/services"
	"github.com/gatehide/gatehide-api/internal/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens and sets user information in context
func AuthMiddleware(authService services.AuthServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		// Validate token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_type", claims.UserType)
		c.Set("user_email", claims.Email)
		c.Set("user_name", claims.Name)
		c.Set("user", claims)

		c.Next()
	}
}

// AdminMiddleware ensures the user is an admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User type not found in context",
			})
			c.Abort()
			return
		}

		if userType != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserMiddleware ensures the user is a regular user
func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User type not found in context",
			})
			c.Abort()
			return
		}

		if userType != "user" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "User access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens if present but doesn't require them
func OptionalAuthMiddleware(authService services.AuthServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Extract token from "Bearer <token>" format
		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		// Validate token if present
		claims, err := authService.ValidateToken(tokenString)
		if err == nil {
			// Set user information in context if token is valid
			c.Set("user_id", claims.UserID)
			c.Set("user_type", claims.UserType)
			c.Set("user_email", claims.Email)
			c.Set("user_name", claims.Name)
			c.Set("user", claims)
		}

		c.Next()
	}
}

// RequireAuthMiddleware ensures user is authenticated (either admin or user)
func RequireAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		if userType != "admin" && userType != "user" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user type",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ExtractTokenFromHeader extracts JWT token from Authorization header
func ExtractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header not found")
	}

	// Check for Bearer token format
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return authHeader[7:], nil
}

// GetCurrentUser extracts current user information from context
func GetCurrentUser(c *gin.Context) (*utils.JWTClaims, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	claims, ok := user.(*utils.JWTClaims)
	return claims, ok
}
