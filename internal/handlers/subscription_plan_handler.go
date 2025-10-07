package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/services"
	"github.com/gin-gonic/gin"
)

// SubscriptionPlanHandler handles subscription plan HTTP requests
type SubscriptionPlanHandler struct {
	service services.SubscriptionPlanServiceInterface
}

// NewSubscriptionPlanHandler creates a new subscription plan handler
func NewSubscriptionPlanHandler(service services.SubscriptionPlanServiceInterface) *SubscriptionPlanHandler {
	return &SubscriptionPlanHandler{service: service}
}

// CreatePlan handles plan creation requests
func (h *SubscriptionPlanHandler) CreatePlan(c *gin.Context) {
	var req models.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	plan, err := h.service.CreatePlan(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create plan",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Plan created successfully",
		"data":    plan,
	})
}

// GetPlan handles plan retrieval requests
func (h *SubscriptionPlanHandler) GetPlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid plan ID",
		})
		return
	}

	plan, err := h.service.GetPlan(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Plan not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": plan,
	})
}

// GetAllPlans handles plan listing requests
func (h *SubscriptionPlanHandler) GetAllPlans(c *gin.Context) {
	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	isActiveStr := c.Query("is_active")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	var isActive *bool
	if isActiveStr != "" {
		active, err := strconv.ParseBool(isActiveStr)
		if err == nil {
			isActive = &active
		}
	}

	plans, total, err := h.service.GetAllPlans(limit, offset, isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get plans",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": plans,
		"pagination": gin.H{
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// UpdatePlan handles plan update requests
func (h *SubscriptionPlanHandler) UpdatePlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid plan ID",
		})
		return
	}

	var req models.UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	plan, err := h.service.UpdatePlan(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update plan",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plan updated successfully",
		"data":    plan,
	})
}

// DeletePlan handles plan deletion requests
func (h *SubscriptionPlanHandler) DeletePlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid plan ID",
		})
		return
	}

	err = h.service.DeletePlan(id)
	if err != nil {
		// Check if it's a security error (active subscriptions)
		if err.Error() == "cannot delete plan: plan has active subscriptions" {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Cannot delete plan",
				"details": "This plan has active subscriptions and cannot be deleted. Please cancel all active subscriptions first.",
			})
			return
		}

		// Check if it's a not found error
		if err.Error() == "subscription plan not found" || strings.Contains(err.Error(), "subscription plan not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Plan not found",
				"details": err.Error(),
			})
			return
		}

		// Other errors
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete plan",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plan deleted successfully",
	})
}
