package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gatehide/gatehide-api/internal/services"
	"github.com/gin-gonic/gin"
)

// RequirePermission checks if the authenticated user has a specific permission
func RequirePermission(permissionService services.PermissionServiceInterface, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User type not found in context",
			})
			c.Abort()
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID not found in context",
			})
			c.Abort()
			return
		}

		userTypeStr, ok := userType.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user type",
			})
			c.Abort()
			return
		}

		userIDInt, ok := userID.(int)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user ID",
			})
			c.Abort()
			return
		}

		// Check permission using user ID and type
		err := permissionService.CheckUserPermission(userIDInt, userTypeStr, resource, action)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Permission denied",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireResourceOwnership checks if the authenticated user owns the resource they're trying to access
func RequireResourceOwnership(permissionService services.PermissionServiceInterface, resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User type not found in context",
			})
			c.Abort()
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID not found in context",
			})
			c.Abort()
			return
		}

		userTypeStr, ok := userType.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user type",
			})
			c.Abort()
			return
		}

		userIDInt, ok := userID.(int)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user ID",
			})
			c.Abort()
			return
		}

		// Extract resource ID from URL parameter
		resourceIDStr := c.Param("id")
		if resourceIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Resource ID not provided",
			})
			c.Abort()
			return
		}

		resourceID, err := strconv.Atoi(resourceIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid resource ID",
			})
			c.Abort()
			return
		}

		canAccess, err := permissionService.CanAccessResource(userTypeStr, resourceType, resourceID, userIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check resource access",
			})
			c.Abort()
			return
		}

		if !canAccess {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied: insufficient ownership",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermissionAndOwnership combines permission check with resource ownership
func RequirePermissionAndOwnership(permissionService services.PermissionServiceInterface, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User type not found in context",
			})
			c.Abort()
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID not found in context",
			})
			c.Abort()
			return
		}

		userTypeStr, ok := userType.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user type",
			})
			c.Abort()
			return
		}

		userIDInt, ok := userID.(int)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user ID",
			})
			c.Abort()
			return
		}

		// Check permission first
		err := permissionService.CheckUserPermission(userIDInt, userTypeStr, resource, action)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Permission denied",
			})
			c.Abort()
			return
		}

		// Extract resource ID from URL parameter
		resourceIDStr := c.Param("id")
		if resourceIDStr != "" {
			resourceID, err := strconv.Atoi(resourceIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid resource ID",
				})
				c.Abort()
				return
			}

			// Check resource ownership
			canAccess, err := permissionService.CanAccessResource(userTypeStr, resource, resourceID, userIDInt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to check resource access",
				})
				c.Abort()
				return
			}

			if !canAccess {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "Access denied: insufficient ownership",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// RequireAdminOnly ensures only administrators can access
func RequireAdminOnly(permissionService services.PermissionServiceInterface) gin.HandlerFunc {
	return RequirePermission(permissionService, "admin", "access")
}

// RequireGamenetOnly ensures only gamenets can access
func RequireGamenetOnly(permissionService services.PermissionServiceInterface) gin.HandlerFunc {
	return RequirePermission(permissionService, "gamenet", "access")
}

// RequireUserOnly ensures only users can access
func RequireUserOnly(permissionService services.PermissionServiceInterface) gin.HandlerFunc {
	return RequirePermission(permissionService, "user", "access")
}
