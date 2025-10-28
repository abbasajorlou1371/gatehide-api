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
	permissionRepo := repositories.NewPermissionRepository(db)

	// Initialize services
	emailService := services.NewEmailService(&cfg.Notification.Email)
	smsService := services.NewSMSService(&cfg.Notification.SMS)
	notificationService := services.NewNotificationService(
		emailService, smsService, nil, nil, notificationRepo, cfg)
	permissionService := services.NewPermissionService(permissionRepo, db)
	authService := services.NewAuthService(userRepo, adminRepo, gamenetRepo, passwordResetRepo, sessionRepo, emailVerificationRepo, notificationService, permissionService, cfg)
	sessionService := services.NewSessionService(sessionRepo, cfg)
	gamenetService := services.NewGamenetService(gamenetRepo, permissionRepo, smsService, emailService)
	userService := services.NewUserService(userRepo, permissionRepo, smsService, emailService)
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
	userHandler := handlers.NewUserHandler(userService)
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

			// Gamenet routes (admin only)
			gamenets := protected.Group("/gamenets")
			gamenets.Use(middlewares.RequirePermission(permissionService, "gamenets", "read"))
			{
				gamenets.GET("/", gamenetHandler.GetAllGamenets)
				gamenets.POST("/", middlewares.RequirePermission(permissionService, "gamenets", "create"), gamenetHandler.CreateGamenet)
				gamenets.GET("/:id", gamenetHandler.GetGamenetByID)
				gamenets.PUT("/:id", middlewares.RequirePermission(permissionService, "gamenets", "update"), gamenetHandler.UpdateGamenet)
				gamenets.DELETE("/:id", middlewares.RequirePermission(permissionService, "gamenets", "delete"), gamenetHandler.DeleteGamenet)
				gamenets.POST("/:id/resend-credentials", middlewares.RequirePermission(permissionService, "gamenets", "update"), gamenetHandler.ResendCredentials)
			}

			// User routes (gamenets can manage their users, admins can manage all)
			users := protected.Group("/users")
			users.Use(middlewares.RequirePermission(permissionService, "users", "read"))
			{
				users.GET("/", userHandler.GetAllUsers)
				users.GET("/search-by-identifier", userHandler.SearchUserByIdentifier)
				users.POST("/", middlewares.RequirePermission(permissionService, "users", "create"), userHandler.CreateUser)
				users.GET("/:id", middlewares.RequireResourceOwnership(permissionService, "users"), userHandler.GetUserByID)
				users.PUT("/:id", middlewares.RequirePermissionAndOwnership(permissionService, "users", "update"), userHandler.UpdateUser)
				users.DELETE("/:id", middlewares.RequirePermissionAndOwnership(permissionService, "users", "delete"), userHandler.DeleteUser)
				users.POST("/:id/resend-credentials", middlewares.RequirePermissionAndOwnership(permissionService, "users", "update"), userHandler.ResendCredentials)
				users.POST("/:id/attach", middlewares.RequirePermission(permissionService, "users", "update"), userHandler.AttachUserToGamenet)
				users.POST("/:id/detach", middlewares.RequirePermission(permissionService, "users", "update"), userHandler.DetachUserFromGamenet)
			}

			// Subscription Plan routes (admin only)
			plans := protected.Group("/subscription-plans")
			plans.Use(middlewares.RequirePermission(permissionService, "subscription_plans", "read"))
			{
				plans.GET("/", subscriptionPlanHandler.GetAllPlans)
				plans.POST("/", middlewares.RequirePermission(permissionService, "subscription_plans", "create"), subscriptionPlanHandler.CreatePlan)
				plans.GET("/:id", subscriptionPlanHandler.GetPlan)
				plans.PUT("/:id", middlewares.RequirePermission(permissionService, "subscription_plans", "update"), subscriptionPlanHandler.UpdatePlan)
				plans.DELETE("/:id", middlewares.RequirePermission(permissionService, "subscription_plans", "delete"), subscriptionPlanHandler.DeletePlan)
			}

			// Dashboard routes with permission checks
			admin := protected.Group("/admin")
			admin.Use(middlewares.RequirePermission(permissionService, "dashboard", "view"))
			{
				admin.GET("/dashboard", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "Admin dashboard", "user": c.GetString("user_name")})
				})
			}

			// User dashboard routes
			user := protected.Group("/user")
			user.Use(middlewares.RequirePermission(permissionService, "dashboard", "view"))
			{
				user.GET("/dashboard", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "User dashboard", "user": c.GetString("user_name")})
				})
			}

			// Gamenet dashboard routes
			gamenet := protected.Group("/gamenet")
			gamenet.Use(middlewares.RequirePermission(permissionService, "dashboard", "view"))
			{
				gamenet.GET("/dashboard", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "Gamenet dashboard", "gamenet": c.GetString("user_name")})
				})
			}
		}
	}

	// Root health endpoint (for load balancers)
	router.GET("/health", healthHandler.Check)
}
