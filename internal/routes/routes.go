package routes

import (
	"database/sql"

	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/internal/handlers"
	"github.com/gatehide/gatehide-api/internal/middlewares"
	"github.com/gatehide/gatehide-api/internal/repositories"
	"github.com/gatehide/gatehide-api/internal/services"
	"github.com/gatehide/gatehide-api/internal/utils"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, cfg *config.Config, db *sql.DB) {
	// Apply global middlewares
	router.Use(middlewares.Logger())
	router.Use(middlewares.CORS())
	router.Use(middlewares.SecurityHeaders())
	router.Use(gin.Recovery())

	// Serve uploaded files
	router.Static("/uploads", cfg.FileStorage.UploadPath)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	passwordResetRepo := repositories.NewPasswordResetRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	emailVerificationRepo := repositories.NewEmailVerificationRepository(db)
	notificationRepo := repositories.NewMySQLNotificationRepository(db)
	gamenetRepo := repositories.NewGamenetRepository(db)
	subscriptionPlanRepo := repositories.NewSubscriptionPlanRepository(db)

	// Initialize services
	emailService := services.NewEmailService(&cfg.Notification.Email)
	smsService := services.NewSMSService(&cfg.Notification.SMS)
	notificationService := services.NewNotificationService(
		emailService, smsService, nil, nil, notificationRepo, cfg)
	authService := services.NewAuthService(userRepo, adminRepo, passwordResetRepo, sessionRepo, emailVerificationRepo, notificationService, cfg)
	sessionService := services.NewSessionService(sessionRepo, cfg)
	gamenetService := services.NewGamenetService(gamenetRepo, smsService, emailService)
	subscriptionPlanService := services.NewSubscriptionPlanService(subscriptionPlanRepo)

	// Initialize file uploader
	fileUploader := utils.NewFileUploader(&cfg.FileStorage)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(cfg)
	authHandler := handlers.NewAuthHandler(authService, fileUploader)
	sessionHandler := handlers.NewSessionHandler(sessionService)
	notificationHandler := handlers.NewNotificationHandler(
		notificationService, nil, nil, authService.GetJWTManager())
	gamenetHandler := handlers.NewGamenetHandler(gamenetService, fileUploader)
	subscriptionPlanHandler := handlers.NewSubscriptionPlanHandler(subscriptionPlanService)

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
		protected.Use(middlewares.AuthMiddlewareWithSession(authService, sessionService))
		{
			// User profile routes (accessible by both users and admins)
			protected.GET("/profile", authHandler.GetProfile)
			protected.PUT("/profile", authHandler.UpdateProfile)
			protected.POST("/profile/upload-image", authHandler.UploadProfileImage)
			protected.POST("/change-password", authHandler.ChangePassword)
			protected.POST("/send-email-verification", authHandler.SendEmailVerification)
			protected.POST("/verify-email-code", authHandler.VerifyEmailCode)

			// Session management routes
			sessions := protected.Group("/sessions")
			{
				sessions.GET("/", sessionHandler.GetActiveSessions)
				sessions.POST("/:session_id/logout", sessionHandler.LogoutSession)
				sessions.POST("/logout-others", sessionHandler.LogoutAllOtherSessions)
				sessions.POST("/logout-all", sessionHandler.LogoutAllSessions)
			}

			// Notification routes
			notifications := protected.Group("/notifications")
			{
				notifications.GET("/", notificationHandler.GetNotifications)
				notifications.GET("/:id", notificationHandler.GetNotification)
			}

			// Gamenet routes
			gamenets := protected.Group("/gamenets")
			{
				gamenets.GET("/", gamenetHandler.GetAllGamenets)
				gamenets.POST("/", gamenetHandler.CreateGamenet)
				gamenets.GET("/:id", gamenetHandler.GetGamenetByID)
				gamenets.PUT("/:id", gamenetHandler.UpdateGamenet)
				gamenets.DELETE("/:id", gamenetHandler.DeleteGamenet)
				gamenets.POST("/:id/resend-credentials", gamenetHandler.ResendCredentials)
			}

			// Subscription Plan routes
			plans := protected.Group("/subscription-plans")
			{
				plans.GET("/", subscriptionPlanHandler.GetAllPlans)
				plans.POST("/", subscriptionPlanHandler.CreatePlan)
				plans.GET("/:id", subscriptionPlanHandler.GetPlan)
				plans.PUT("/:id", subscriptionPlanHandler.UpdatePlan)
				plans.DELETE("/:id", subscriptionPlanHandler.DeletePlan)
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
