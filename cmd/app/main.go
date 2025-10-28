package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/internal/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Connect to database
	db, err := sql.Open(cfg.Database.Driver, cfg.GetDSN())
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}
	log.Printf("✅ Database connection established")

	// Initialize Gin router
	router := gin.New()

	// Setup routes
	routes.SetupRoutes(router, cfg, db)

	// Server information
	log.Printf("🚀 Starting %s v%s", cfg.App.Name, cfg.App.Version)
	log.Printf("📡 Server running on port %s", cfg.Server.Port)
	log.Printf("🔧 Environment: %s", cfg.Server.GinMode)
	log.Printf("🏥 Health check available at: http://localhost:%s/health", cfg.Server.Port)

	// Start server
	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := router.Run(address); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
