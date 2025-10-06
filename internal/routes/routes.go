package routes

import (
	"database/sql"

	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/internal/handlers"
	"github.com/gatehide/gatehide-api/internal/middlewares"
	"github.com/gatehide/gatehide-api/internal/repositories"
	"github.com/gatehide/gatehide-api/internal/services"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, cfg *config.Config, db *sql.DB) {
	// Apply global middlewares
	router.Use(middlewares.Logger())
	router.Use(middlewares.CORS())
	router.Use(middlewares.SecurityHeaders())
	router.Use(gin.Recovery())

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	passwordResetRepo := repositories.NewPasswordResetRepository(db)
	notificationRepo := repositories.NewMySQLNotificationRepository(db)

	// Initialize services
	emailService := services.NewEmailService(&cfg.Notification.Email)
	notificationService := services.NewNotificationService(
		emailService, nil, nil, nil, notificationRepo, cfg)
	authService := services.NewAuthService(userRepo, adminRepo, passwordResetRepo, notificationService, cfg)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(cfg)
	authHandler := handlers.NewAuthHandler(authService)
	notificationHandler := handlers.NewNotificationHandler(
		notificationService, nil, nil, authService.GetJWTManager())

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		public := v1.Group("/")
		{
			// Health check endpoint
			public.GET("/health", healthHandler.Check)

			// Authentication routes
			auth := public.Group("/auth")
			{
				// Unified login endpoint (automatically determines user type)
				auth.POST("/login", authHandler.Login)
				auth.POST("/refresh", authHandler.RefreshToken)
				auth.POST("/logout", authHandler.Logout)

				// Password reset routes
				auth.POST("/forgot-password", authHandler.ForgotPassword)
				auth.POST("/reset-password", authHandler.ResetPassword)
				auth.GET("/validate-reset-token", authHandler.ValidateResetToken)
			}
		}

		// Protected routes (authentication required)
		protected := v1.Group("/")
		protected.Use(middlewares.AuthMiddleware(authService))
		{
			// User profile routes (accessible by both users and admins)
			protected.GET("/profile", authHandler.GetProfile)

			// Notification routes
			notifications := protected.Group("/notifications")
			{
				notifications.GET("/", notificationHandler.GetNotifications)
				notifications.GET("/:id", notificationHandler.GetNotification)
			}

			// Admin-only routes
			admin := protected.Group("/admin")
			admin.Use(middlewares.AdminMiddleware())
			{
				// Add admin-specific routes here
				admin.GET("/dashboard", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "Admin dashboard", "user": c.GetString("user_name")})
				})
			}

			// User-only routes
			user := protected.Group("/user")
			user.Use(middlewares.UserMiddleware())
			{
				// Add user-specific routes here
				user.GET("/dashboard", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "User dashboard", "user": c.GetString("user_name")})
				})
			}
		}
	}

	// Root health endpoint (for load balancers)
	router.GET("/health", healthHandler.Check)
}
