package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/internal/handlers"
	"github.com/gatehide/gatehide-api/internal/middlewares"
	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/repositories"
	"github.com/gatehide/gatehide-api/internal/services"
	"github.com/gatehide/gatehide-api/internal/utils"
	testutils "github.com/gatehide/gatehide-api/tests/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticationIntegration_UserLogin(t *testing.T) {
	testutils.SkipIfNoDB(t)

	db := testutils.SetupTestDB(t)
	defer db.Close()
	defer testutils.CleanupTestDBForce(t, db)

	// Clean up any existing test data before starting
	testutils.CleanupTestDBForce(t, db)

	// Setup test data with unique email
	testUser := testutils.CreateTestUser(t, db, "user1@example.com", "password123", "Test User 1")

	// Setup application
	cfg := testutils.TestConfig()
	router := setupTestRouter(cfg, db)

	tests := []struct {
		name           string
		requestBody    models.LoginRequest
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "valid user login",
			requestBody: models.LoginRequest{
				Email:    "user1@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name: "invalid password",
			requestBody: models.LoginRequest{
				Email:    "user1@example.com",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
		{
			name: "non-existing user",
			requestBody: models.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test unified login endpoint
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectToken {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "data")

				data := response["data"].(map[string]interface{})
				assert.Contains(t, data, "token")
				assert.Contains(t, data, "user_type")
				assert.Contains(t, data, "user")
				assert.Contains(t, data, "expires_at")

				assert.Equal(t, "user", data["user_type"])

				user := data["user"].(map[string]interface{})
				assert.Equal(t, testUser.ID, int(user["id"].(float64)))
				assert.Equal(t, testUser.Email, user["email"])
				assert.Equal(t, testUser.Name, user["name"])
			} else {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			}
		})
	}
}

func TestAuthenticationIntegration_AdminLogin(t *testing.T) {
	testutils.SkipIfNoDB(t)

	db := testutils.SetupTestDB(t)
	defer db.Close()
	defer testutils.CleanupTestDBForce(t, db)

	// Setup test data with unique email
	testAdmin := testutils.CreateTestAdmin(t, db, "admin1@example.com", "admin123", "Test Admin 1")

	// Setup application
	cfg := testutils.TestConfig()
	router := setupTestRouter(cfg, db)

	tests := []struct {
		name           string
		requestBody    models.LoginRequest
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "valid admin login",
			requestBody: models.LoginRequest{
				Email:    "admin1@example.com",
				Password: "admin123",
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name: "invalid password",
			requestBody: models.LoginRequest{
				Email:    "admin1@example.com",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test unified login endpoint
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectToken {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "data")

				data := response["data"].(map[string]interface{})
				assert.Contains(t, data, "token")
				assert.Contains(t, data, "user_type")
				assert.Contains(t, data, "user")
				assert.Contains(t, data, "expires_at")

				assert.Equal(t, "admin", data["user_type"])

				user := data["user"].(map[string]interface{})
				assert.Equal(t, testAdmin.ID, int(user["id"].(float64)))
				assert.Equal(t, testAdmin.Email, user["email"])
				assert.Equal(t, testAdmin.Name, user["name"])
			} else {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			}
		})
	}
}

func TestAuthenticationIntegration_ProtectedRoutes(t *testing.T) {
	testutils.SkipIfNoDB(t)

	db := testutils.SetupTestDB(t)
	defer db.Close()
	defer testutils.CleanupTestDBForce(t, db)

	// Setup test data with unique emails
	_ = testutils.CreateTestUser(t, db, "user2@example.com", "password123", "Test User 2")
	_ = testutils.CreateTestAdmin(t, db, "admin2@example.com", "admin123", "Test Admin 2")

	// Setup application
	cfg := testutils.TestConfig()
	router := setupTestRouter(cfg, db)

	// Get tokens for testing
	userToken := getAuthToken(t, router, "user2@example.com", "password123")
	adminToken := getAuthToken(t, router, "admin2@example.com", "admin123")

	tests := []struct {
		name           string
		endpoint       string
		method         string
		authToken      string
		expectedStatus int
	}{
		{
			name:           "user profile with user token",
			endpoint:       "/api/v1/profile",
			method:         "GET",
			authToken:      userToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user profile with admin token",
			endpoint:       "/api/v1/profile",
			method:         "GET",
			authToken:      adminToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user profile without token",
			endpoint:       "/api/v1/profile",
			method:         "GET",
			authToken:      "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "user dashboard with user token",
			endpoint:       "/api/v1/user/dashboard",
			method:         "GET",
			authToken:      userToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user dashboard with admin token",
			endpoint:       "/api/v1/user/dashboard",
			method:         "GET",
			authToken:      adminToken,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "admin dashboard with admin token",
			endpoint:       "/api/v1/admin/dashboard",
			method:         "GET",
			authToken:      adminToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "admin dashboard with user token",
			endpoint:       "/api/v1/admin/dashboard",
			method:         "GET",
			authToken:      userToken,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.endpoint, nil)
			if tt.authToken != "" {
				req.Header.Set("Authorization", "Bearer "+tt.authToken)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAuthenticationIntegration_TokenRefresh(t *testing.T) {
	testutils.SkipIfNoDB(t)

	db := testutils.SetupTestDB(t)
	defer db.Close()
	defer testutils.CleanupTestDBForce(t, db)

	// Setup test data with unique email
	_ = testutils.CreateTestUser(t, db, "user3@example.com", "password123", "Test User 3")

	// Setup application
	cfg := testutils.TestConfig()
	router := setupTestRouter(cfg, db)

	// Get initial token
	initialToken := getAuthToken(t, router, "user3@example.com", "password123")

	tests := []struct {
		name           string
		authToken      string
		expectedStatus int
		expectNewToken bool
	}{
		{
			name:           "valid token refresh",
			authToken:      initialToken,
			expectedStatus: http.StatusOK,
			expectNewToken: true,
		},
		{
			name:           "invalid token refresh",
			authToken:      "invalid.jwt.token",
			expectedStatus: http.StatusUnauthorized,
			expectNewToken: false,
		},
		{
			name:           "missing token",
			authToken:      "",
			expectedStatus: http.StatusUnauthorized,
			expectNewToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
			if tt.authToken != "" {
				req.Header.Set("Authorization", "Bearer "+tt.authToken)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectNewToken {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "data")

				data := response["data"].(map[string]interface{})
				assert.Contains(t, data, "token")

				// Check that we got a valid token (it might be the same if generated within the same second)
				refreshedToken := data["token"].(string)
				assert.NotEmpty(t, refreshedToken)

				// Verify the refreshed token is valid by trying to use it
				// This is more realistic than checking if tokens are different
				req2 := httptest.NewRequest("GET", "/api/v1/profile", nil)
				req2.Header.Set("Authorization", "Bearer "+refreshedToken)
				w2 := httptest.NewRecorder()
				router.ServeHTTP(w2, req2)
				assert.Equal(t, http.StatusOK, w2.Code)
			} else {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			}
		})
	}
}

func TestAuthenticationIntegration_Logout(t *testing.T) {
	testutils.SkipIfNoDB(t)

	cfg := testutils.TestConfig()
	db := testutils.SetupTestDB(t)
	defer db.Close()
	defer testutils.CleanupTestDBForce(t, db)

	// Setup test data
	_ = testutils.CreateTestUser(t, db, "logout@example.com", "password123", "Logout Test User")

	router := setupTestRouter(cfg, db)

	// First, login to get a valid token
	loginReq := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{
		"email": "logout@example.com",
		"password": "password123"
	}`))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReq)

	assert.Equal(t, http.StatusOK, loginW.Code)

	var loginResponse map[string]interface{}
	err := json.Unmarshal(loginW.Body.Bytes(), &loginResponse)
	assert.NoError(t, err)
	assert.Contains(t, loginResponse, "data")

	data := loginResponse["data"].(map[string]interface{})
	token := data["token"].(string)

	// Now test logout with the token
	req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "message")
	assert.Equal(t, "Logout successful", response["message"])
}

// Helper functions

func setupTestRouter(cfg *config.Config, db *sql.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	passwordResetRepo := repositories.NewPasswordResetRepository(db)
	gamenetRepo := repositories.NewGamenetRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	emailVerificationRepo := repositories.NewEmailVerificationRepository(db)
	notificationService := &testutils.MockNotificationService{}

	// Initialize file uploader
	fileUploader := utils.NewFileUploader(&cfg.FileStorage)

	// Initialize services
	authService := services.NewAuthService(userRepo, adminRepo, gamenetRepo, passwordResetRepo, sessionRepo, emailVerificationRepo, notificationService, cfg)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(cfg)
	authHandler := handlers.NewAuthHandler(authService, fileUploader)

	// Setup routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		public := v1.Group("/")
		{
			public.GET("/health", healthHandler.Check)

			auth := public.Group("/auth")
			{
				auth.POST("/login", authHandler.Login)
				auth.POST("/refresh", authHandler.RefreshToken)
				auth.POST("/logout", authHandler.Logout)
			}
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middlewares.AuthMiddleware(authService))
		{
			protected.GET("/profile", authHandler.GetProfile)

			admin := protected.Group("/admin")
			admin.Use(middlewares.AdminMiddleware())
			{
				admin.GET("/dashboard", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "Admin dashboard", "user": c.GetString("user_name")})
				})
			}

			user := protected.Group("/user")
			user.Use(middlewares.UserMiddleware())
			{
				user.GET("/dashboard", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "User dashboard", "user": c.GetString("user_name")})
				})
			}
		}
	}

	return router
}

func getAuthToken(t *testing.T, router *gin.Engine, email, password string) string {
	loginReq := models.LoginRequest{
		Email:    email,
		Password: password,
	}

	jsonBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	token := data["token"].(string)
	assert.NotEmpty(t, token)

	return token
}
