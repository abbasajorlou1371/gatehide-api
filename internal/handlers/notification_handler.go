package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/services"
	"github.com/gatehide/gatehide-api/internal/utils"
	"github.com/gin-gonic/gin"
)

// NotificationHandler handles notification-related HTTP requests
type NotificationHandler struct {
	notificationService services.NotificationServiceInterface
	templateService     services.TemplateServiceInterface
	emailService        services.EmailServiceInterface
	jwtManager          *utils.JWTManager
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(
	notificationService services.NotificationServiceInterface,
	templateService services.TemplateServiceInterface,
	emailService services.EmailServiceInterface,
	jwtManager *utils.JWTManager,
) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		templateService:     templateService,
		emailService:        emailService,
		jwtManager:          jwtManager,
	}
}

// GetNotification handles GET /api/notifications/:id
func (h *NotificationHandler) GetNotification(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	// Get user from token for authorization
	user, err := h.getUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	notification, err := h.notificationService.GetNotification(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Check if user has permission to view this notification
	if !h.canUserViewNotification(user, notification) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Mark notification as read if it's a database notification
	if notification.Type == models.NotificationTypeDatabase {
		// Add metadata to track read status
		if notification.Metadata == nil {
			notification.Metadata = make(map[string]interface{})
		}
		notification.Metadata["read_by"] = user.UserID
		notification.Metadata["read_at"] = time.Now()

		// Update notification status
		if err := h.notificationService.UpdateNotificationStatus(c.Request.Context(), id, models.NotificationStatusSent, nil); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Warning: failed to mark notification as read: %v\n", err)
		}
	}

	c.JSON(http.StatusOK, notification.ToResponse())
}

// GetNotifications handles GET /api/notifications
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	// Get user from token for authorization
	user, err := h.getUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Build filters from query parameters
	filters := make(map[string]interface{})

	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	if notificationType := c.Query("type"); notificationType != "" {
		filters["type"] = notificationType
	}

	if recipient := c.Query("recipient"); recipient != "" {
		filters["recipient"] = recipient
	}

	if priority := c.Query("priority"); priority != "" {
		filters["priority"] = priority
	}

	if limit := c.Query("limit"); limit != "" {
		if limitInt, err := strconv.Atoi(limit); err == nil {
			filters["limit"] = limitInt
		}
	}

	// Add user-specific filter for non-admin users
	if user.UserType != "admin" {
		filters["recipient"] = user.Email
	}

	notifications, err := h.notificationService.GetNotifications(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	var responses []models.NotificationResponse
	for _, notification := range notifications {
		responses = append(responses, notification.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{"notifications": responses})
}

// getUserFromToken extracts user information from JWT token
func (h *NotificationHandler) getUserFromToken(c *gin.Context) (*utils.JWTClaims, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		return nil, fmt.Errorf("no authorization header")
	}

	// Remove "Bearer " prefix
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	return h.jwtManager.ValidateToken(token)
}

// canUserViewNotification checks if user can view a notification
func (h *NotificationHandler) canUserViewNotification(user *utils.JWTClaims, notification *models.Notification) bool {
	// Admins can view all notifications
	if user.UserType == "admin" {
		return true
	}

	// Users can only view their own notifications
	return notification.Recipient == user.Email
}
