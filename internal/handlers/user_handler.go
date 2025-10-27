package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/services"
	"github.com/gin-gonic/gin"
)

// UserHandler handles user HTTP requests
type UserHandler struct {
	userService services.UserServiceInterface
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetAllUsers handles GET /users
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	// Check if search parameters are provided
	query := c.Query("query")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	// Check if user is a gamenet
	userType, _ := c.Get("user_type")
	userID, _ := c.Get("user_id")

	var gamenetID *int
	if userType == "gamenet" {
		if id, ok := userID.(int); ok {
			gamenetID = &id
		}
	}

	// Debug logging
	fmt.Printf("DEBUG: query=%s, pageStr=%s, pageSizeStr=%s, userType=%v, gamenetID=%v\n", query, pageStr, pageSizeStr, userType, gamenetID)

	// If search parameters are provided, use search endpoint
	if query != "" || pageStr != "" || pageSizeStr != "" {
		fmt.Printf("DEBUG: Using search path\n")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}

		searchReq := &models.UserSearchRequest{
			Query:    query,
			Page:     page,
			PageSize: pageSize,
		}

		var result *models.UserSearchResponse
		if gamenetID != nil {
			// Gamenet users only see their own users
			result, err = h.userService.SearchByGamenet(c.Request.Context(), searchReq, *gamenetID)
		} else {
			// Admins see all users
			result, err = h.userService.Search(c.Request.Context(), searchReq)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to search users",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Users retrieved successfully",
			"data":       result.Data,
			"pagination": result.Pagination,
		})
		return
	}

	// Default behavior - get all users
	var users []models.UserResponse
	var err error

	if gamenetID != nil {
		// Gamenet users only see their own users
		users, err = h.userService.GetAllByGamenet(c.Request.Context(), *gamenetID)
	} else {
		// Admins see all users
		users, err = h.userService.GetAll(c.Request.Context())
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"data":    users,
	})
}

// GetUserByID handles GET /users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User retrieved successfully",
		"data":    user,
	})
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if user is a gamenet
	userType, _ := c.Get("user_type")
	userID, _ := c.Get("user_id")

	var gamenetID *int
	if userType == "gamenet" {
		if id, ok := userID.(int); ok {
			gamenetID = &id
		}
	}

	user, err := h.userService.Create(c.Request.Context(), &req, gamenetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    user,
	})
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Check permissions
	userType, _ := c.Get("user_type")
	userID, _ := c.Get("user_id")

	requesterType := ""
	requesterID := 0

	if ut, ok := userType.(string); ok {
		requesterType = ut
	}
	if uid, ok := userID.(int); ok {
		requesterID = uid
	}

	// Check if requester can modify this user
	canModify, err := h.userService.CanModifyUser(c.Request.Context(), id, requesterID, requesterType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check permissions",
		})
		return
	}

	if !canModify {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to modify this user",
		})
		return
	}

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.userService.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"data":    user,
	})
}

// DeleteUser handles DELETE /users/:id (Admin only)
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Check if user is admin
	userType, _ := c.Get("user_type")
	if userType != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only admins can delete users",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	err = h.userService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// ResendCredentials handles POST /users/:id/resend-credentials
func (h *UserHandler) ResendCredentials(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	err = h.userService.ResendCredentials(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Credentials sent successfully via SMS",
	})
}

// SearchUserByIdentifier handles GET /users/search-by-identifier?q=email_or_mobile
func (h *UserHandler) SearchUserByIdentifier(c *gin.Context) {
	identifier := c.Query("q")
	if identifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Search query is required",
		})
		return
	}

	// Try to find user by email first
	user, err := h.userService.GetByEmail(c.Request.Context(), identifier)
	if err == nil && user != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "User found",
			"data":    user,
			"found":   true,
		})
		return
	}

	// Try to find user by mobile
	user, err = h.userService.GetByMobile(c.Request.Context(), identifier)
	if err == nil && user != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "User found",
			"data":    user,
			"found":   true,
		})
		return
	}

	// User not found
	c.JSON(http.StatusOK, gin.H{
		"message": "User not found",
		"found":   false,
	})
}

// AttachUserToGamenet handles POST /users/:id/attach
func (h *UserHandler) AttachUserToGamenet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Check if user is a gamenet
	userType, _ := c.Get("user_type")
	userID, _ := c.Get("user_id")

	if userType != "gamenet" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only gamenets can attach users",
		})
		return
	}

	gamenetID, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid gamenet ID",
		})
		return
	}

	err = h.userService.AttachToGamenet(c.Request.Context(), id, gamenetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User attached to gamenet successfully",
	})
}

// DetachUserFromGamenet handles POST /users/:id/detach
func (h *UserHandler) DetachUserFromGamenet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Check if user is a gamenet
	userType, _ := c.Get("user_type")
	userID, _ := c.Get("user_id")

	if userType != "gamenet" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only gamenets can detach users",
		})
		return
	}

	gamenetID, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid gamenet ID",
		})
		return
	}

	err = h.userService.DetachFromGamenet(c.Request.Context(), id, gamenetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User detached from gamenet successfully",
	})
}
