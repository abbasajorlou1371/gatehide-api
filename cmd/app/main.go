package main

import (
	"fmt"
	"log"

	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Initialize Gin router
	router := gin.New()

	// Setup routes
	routes.SetupRoutes(router, cfg)

	// Server information
	log.Printf("🚀 Starting %s v%s", cfg.App.Name, cfg.App.Version)
	log.Printf("📡 Server running on port %s", cfg.Server.Port)
	log.Printf("🔧 Environment: %s", cfg.Server.GinMode)
	log.Printf("🏥 Health check available at: http://localhost:%s/health", cfg.Server.Port)

	// Start server
	address := fmt.Sprintf(":%s", cfg.Server.Port)
	if err := router.Run(address); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
