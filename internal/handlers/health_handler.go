package handlers

import (
	"net/http"
	"time"

	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	config *config.Config
}

// NewHealthHandler creates a new health handler instance
func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		config: cfg,
	}
}

// Check handles the health check endpoint
// @Summary Health Check
// @Description Check if the API is running and healthy
// @Tags Health
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Router /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	response := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   h.config.App.Name,
		Version:   h.config.App.Version,
	}

	c.JSON(http.StatusOK, response)
}
