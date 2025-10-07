package handlers

import (
	"net/http"
	"strconv"

	"github.com/gatehide/gatehide-api/internal/middlewares"
	"github.com/gatehide/gatehide-api/internal/services"
	"github.com/gin-gonic/gin"
)

// SessionHandler handles session-related HTTP requests
type SessionHandler struct {
	sessionService services.SessionServiceInterface
}

// NewSessionHandler creates a new session handler
func NewSessionHandler(sessionService services.SessionServiceInterface) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
	}
}

// GetActiveSessions retrieves all active sessions for the current user
// @Summary Get active sessions
// @Description Get all active sessions for the current user
// @Tags sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Active sessions retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /sessions [get]
func (h *SessionHandler) GetActiveSessions(c *gin.Context) {
	// Get current user from context
	claims, exists := middlewares.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found in context",
		})
		return
	}

	// Get current session token
	currentToken, err := middlewares.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to extract token",
		})
		return
	}

	// Get active sessions
	sessions, err := h.sessionService.GetActiveSessions(claims.UserID, claims.UserType, currentToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get active sessions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Active sessions retrieved successfully",
		"sessions": sessions,
		"count":    len(sessions),
	})
}

// LogoutSession deactivates a specific session
// @Summary Logout a session
// @Description Deactivate a specific session by ID
// @Tags sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id path int true "Session ID"
// @Success 200 {object} map[string]interface{} "Session logged out successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Session not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /sessions/{session_id}/logout [post]
func (h *SessionHandler) LogoutSession(c *gin.Context) {
	// Get current user from context
	claims, exists := middlewares.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found in context",
		})
		return
	}

	// Get session ID from URL parameter
	sessionID, err := parseSessionID(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid session ID",
		})
		return
	}

	// Logout the session
	err = h.sessionService.LogoutSession(sessionID, claims.UserID, claims.UserType)
	if err != nil {
		if err.Error() == "session not found or does not belong to user" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Session not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to logout session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Session logged out successfully",
	})
}

// LogoutAllOtherSessions deactivates all sessions except the current one
// @Summary Logout all other sessions
// @Description Deactivate all sessions except the current one
// @Tags sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "All other sessions logged out successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /sessions/logout-others [post]
func (h *SessionHandler) LogoutAllOtherSessions(c *gin.Context) {
	// Get current user from context
	claims, exists := middlewares.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found in context",
		})
		return
	}

	// Get current session token
	currentToken, err := middlewares.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to extract token",
		})
		return
	}

	// Logout all other sessions
	err = h.sessionService.LogoutAllOtherSessions(claims.UserID, claims.UserType, currentToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to logout other sessions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All other sessions logged out successfully",
	})
}

// LogoutAllSessions deactivates all sessions for the user
// @Summary Logout all sessions
// @Description Deactivate all sessions for the current user
// @Tags sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "All sessions logged out successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /sessions/logout-all [post]
func (h *SessionHandler) LogoutAllSessions(c *gin.Context) {
	// Get current user from context
	claims, exists := middlewares.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found in context",
		})
		return
	}

	// Logout all sessions
	err := h.sessionService.LogoutAllSessions(claims.UserID, claims.UserType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to logout all sessions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All sessions logged out successfully",
	})
}

// parseSessionID parses session ID from string to int
func parseSessionID(sessionIDStr string) (int, error) {
	return strconv.Atoi(sessionIDStr)
}
