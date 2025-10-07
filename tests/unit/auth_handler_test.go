package unit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gatehide/gatehide-api/internal/handlers"
	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/utils"
	testutils "github.com/gatehide/gatehide-api/tests/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*testutils.MockAuthService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "valid user login",
			requestBody: models.LoginRequest{
				Email:    "user@example.com",
				Password: "password123",
			},
			mockSetup: func(m *testutils.MockAuthService) {
				response := &models.LoginResponse{
					Token:     "valid.jwt.token",
					UserType:  "user",
					User:      models.UserResponse{ID: 1, Email: "user@example.com", Name: "Test User"},
					ExpiresAt: time.Now().Add(24 * time.Hour),
				}
				m.On("LoginWithSession", "user@example.com", "password123", false, "", "192.0.2.1", "").Return(response, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "valid admin login",
			requestBody: models.LoginRequest{
				Email:    "admin@example.com",
				Password: "admin123",
			},
			mockSetup: func(m *testutils.MockAuthService) {
				response := &models.LoginResponse{
					Token:     "valid.jwt.token",
					UserType:  "admin",
					User:      models.AdminResponse{ID: 1, Email: "admin@example.com", Name: "Test Admin"},
					ExpiresAt: time.Now().Add(24 * time.Hour),
				}
				m.On("LoginWithSession", "admin@example.com", "admin123", false, "", "192.0.2.1", "").Return(response, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "invalid credentials",
			requestBody: models.LoginRequest{
				Email:    "user@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(m *testutils.MockAuthService) {
				m.On("LoginWithSession", "user@example.com", "wrongpassword", false, "", "192.0.2.1", "").Return((*models.LoginResponse)(nil), assert.AnError)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockService := new(testutils.MockAuthService)
			tt.mockSetup(mockService)

			// Setup file uploader
			cfg := testutils.TestConfig()
			fileUploader := utils.NewFileUploader(&cfg.FileStorage)

			// Setup handler
			handler := handlers.NewAuthHandler(mockService, fileUploader)

			// Setup request
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Setup response recorder
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Execute
			handler.Login(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			} else {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
				assert.Contains(t, response, "data")
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		mockSetup      func(*testutils.MockAuthService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:       "valid token refresh",
			authHeader: "Bearer valid.jwt.token",
			mockSetup: func(m *testutils.MockAuthService) {
				m.On("RefreshToken", "valid.jwt.token", false).Return("new.jwt.token", nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:       "invalid token refresh",
			authHeader: "Bearer invalid.jwt.token",
			mockSetup: func(m *testutils.MockAuthService) {
				m.On("RefreshToken", "invalid.jwt.token", false).Return("", assert.AnError)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			mockSetup:      func(m *testutils.MockAuthService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:       "invalid authorization header format",
			authHeader: "InvalidFormat valid.jwt.token",
			mockSetup: func(m *testutils.MockAuthService) {
				m.On("RefreshToken", "InvalidFormat valid.jwt.token", false).Return("", assert.AnError)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockService := new(testutils.MockAuthService)
			tt.mockSetup(mockService)

			// Setup file uploader
			cfg := testutils.TestConfig()
			fileUploader := utils.NewFileUploader(&cfg.FileStorage)

			// Setup handler
			handler := handlers.NewAuthHandler(mockService, fileUploader)

			// Setup request
			req := httptest.NewRequest("POST", "/auth/refresh", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Setup response recorder
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Execute
			handler.RefreshToken(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			} else {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
				assert.Contains(t, response, "data")
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup handler
	mockService := new(testutils.MockAuthService)
	// Setup mock expectation for ValidateToken
	mockService.On("ValidateToken", "valid.jwt.token").Return(&utils.JWTClaims{
		UserID:   1,
		UserType: "user",
		Email:    "test@example.com",
		Name:     "Test User",
	}, nil)
	cfg := testutils.TestConfig()
	fileUploader := utils.NewFileUploader(&cfg.FileStorage)
	handler := handlers.NewAuthHandler(mockService, fileUploader)

	// Setup request
	req := httptest.NewRequest("POST", "/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")

	// Setup response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.Logout(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "message")
	assert.Equal(t, "Logout successful", response["message"])
}

func TestAuthHandler_GetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "valid user in context",
			setupContext: func(c *gin.Context) {
				user := &utils.JWTClaims{
					UserID:   1,
					UserType: "user",
					Email:    "user@example.com",
					Name:     "Test User",
				}
				c.Set("user", user)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "no user in context",
			setupContext: func(c *gin.Context) {
				// Don't set any user in context
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup handler
			mockService := new(testutils.MockAuthService)

			// Setup mock expectations for valid user case
			if tt.name == "valid user in context" {
				mockUser := &models.User{
					ID:    1,
					Name:  "Test User",
					Email: "user@example.com",
				}
				mockService.On("GetUserByID", 1).Return(mockUser, nil)
			}

			cfg := testutils.TestConfig()
			fileUploader := utils.NewFileUploader(&cfg.FileStorage)
			handler := handlers.NewAuthHandler(mockService, fileUploader)

			// Setup request
			req := httptest.NewRequest("GET", "/profile", nil)

			// Setup response recorder
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Setup context
			tt.setupContext(c)

			// Execute
			handler.GetProfile(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			} else {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
				assert.Contains(t, response, "data")
			}
		})
	}
}
