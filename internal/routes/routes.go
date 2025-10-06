package routes

import (
	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/internal/handlers"
	"github.com/gatehide/gatehide-api/internal/middlewares"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// Apply global middlewares
	router.Use(middlewares.Logger())
	router.Use(middlewares.CORS())
	router.Use(middlewares.SecurityHeaders())
	router.Use(gin.Recovery())

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(cfg)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Health check endpoint
		v1.GET("/health", healthHandler.Check)
	}

	// Root health endpoint (for load balancers)
	router.GET("/health", healthHandler.Check)
}
