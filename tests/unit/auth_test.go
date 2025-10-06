package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gatehide/gatehide-api/internal/middlewares"
	"github.com/gatehide/gatehide-api/internal/utils"
	testutils "github.com/gatehide/gatehide-api/tests/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		mockSetup      func(*testutils.MockAuthService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:       "valid token",
			authHeader: "Bearer valid.jwt.token",
			mockSetup: func(m *testutils.MockAuthService) {
				claims := &utils.JWTClaims{
					UserID:   1,
					UserType: "user",
					Email:    "user@example.com",
					Name:     "Test User",
				}
				m.On("ValidateToken", "valid.jwt.token").Return(claims, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:       "invalid token",
			authHeader: "Bearer invalid.jwt.token",
			mockSetup: func(m *testutils.MockAuthService) {
				m.On("ValidateToken", "invalid.jwt.token").Return((*utils.JWTClaims)(nil), assert.AnError)
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
				m.On("ValidateToken", "InvalidFormat valid.jwt.token").Return((*utils.JWTClaims)(nil), assert.AnError)
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

			// Setup router with middleware
			router := gin.New()
			router.Use(middlewares.AuthMiddleware(mockService))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Setup request
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Setup response recorder
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError {
				assert.Contains(t, w.Body.String(), "error")
			} else {
				assert.Contains(t, w.Body.String(), "success")
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthMiddleware_ContextValues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup mock
	mockService := new(testutils.MockAuthService)
	claims := &utils.JWTClaims{
		UserID:   123,
		UserType: "admin",
		Email:    "admin@example.com",
		Name:     "Test Admin",
	}
	mockService.On("ValidateToken", "valid.jwt.token").Return(claims, nil)

	// Setup router with middleware
	router := gin.New()
	router.Use(middlewares.AuthMiddleware(mockService))
	router.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		assert.True(t, exists)
		assert.Equal(t, 123, userID)

		userType, exists := c.Get("user_type")
		assert.True(t, exists)
		assert.Equal(t, "admin", userType)

		email, exists := c.Get("user_email")
		assert.True(t, exists)
		assert.Equal(t, "admin@example.com", email)

		name, exists := c.Get("user_name")
		assert.True(t, exists)
		assert.Equal(t, "Test Admin", name)

		user, exists := c.Get("user")
		assert.True(t, exists)
		assert.Equal(t, claims, user)

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Setup request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")

	// Setup response recorder
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")

	mockService.AssertExpectations(t)
}

func TestAdminMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userType       string
		setupUser      func(*gin.Context)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:     "valid admin user",
			userType: "admin",
			setupUser: func(c *gin.Context) {
				c.Set("user_type", "admin")
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:     "invalid user type",
			userType: "user",
			setupUser: func(c *gin.Context) {
				c.Set("user_type", "user")
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  true,
		},
		{
			name:     "no user type set",
			userType: "",
			setupUser: func(c *gin.Context) {
				// Don't set user_type
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup router with middleware
			router := gin.New()
			router.Use(func(c *gin.Context) {
				tt.setupUser(c)
				c.Next()
			})
			router.Use(middlewares.AdminMiddleware())
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
			})

			// Setup request
			req := httptest.NewRequest("GET", "/test", nil)

			// Setup response recorder
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError {
				assert.Contains(t, w.Body.String(), "error")
			} else {
				assert.Contains(t, w.Body.String(), "admin access granted")
			}
		})
	}
}

func TestUserMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userType       string
		setupUser      func(*gin.Context)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:     "valid user",
			userType: "user",
			setupUser: func(c *gin.Context) {
				c.Set("user_type", "user")
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:     "invalid user type (admin)",
			userType: "admin",
			setupUser: func(c *gin.Context) {
				c.Set("user_type", "admin")
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  true,
		},
		{
			name:     "no user type set",
			userType: "",
			setupUser: func(c *gin.Context) {
				// Don't set user_type
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup router with middleware
			router := gin.New()
			router.Use(func(c *gin.Context) {
				tt.setupUser(c)
				c.Next()
			})
			router.Use(middlewares.UserMiddleware())
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "user access granted"})
			})

			// Setup request
			req := httptest.NewRequest("GET", "/test", nil)

			// Setup response recorder
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError {
				assert.Contains(t, w.Body.String(), "error")
			} else {
				assert.Contains(t, w.Body.String(), "user access granted")
			}
		})
	}
}

func TestOptionalAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		mockSetup      func(*testutils.MockAuthService)
		expectedStatus int
		expectUserSet  bool
	}{
		{
			name:       "valid token provided",
			authHeader: "Bearer valid.jwt.token",
			mockSetup: func(m *testutils.MockAuthService) {
				claims := &utils.JWTClaims{
					UserID:   1,
					UserType: "user",
					Email:    "user@example.com",
					Name:     "Test User",
				}
				m.On("ValidateToken", "valid.jwt.token").Return(claims, nil)
			},
			expectedStatus: http.StatusOK,
			expectUserSet:  true,
		},
		{
			name:       "invalid token provided",
			authHeader: "Bearer invalid.jwt.token",
			mockSetup: func(m *testutils.MockAuthService) {
				m.On("ValidateToken", "invalid.jwt.token").Return((*utils.JWTClaims)(nil), assert.AnError)
			},
			expectedStatus: http.StatusOK,
			expectUserSet:  false,
		},
		{
			name:           "no token provided",
			authHeader:     "",
			mockSetup:      func(m *testutils.MockAuthService) {},
			expectedStatus: http.StatusOK,
			expectUserSet:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockService := new(testutils.MockAuthService)
			tt.mockSetup(mockService)

			// Setup router with middleware
			router := gin.New()
			router.Use(middlewares.OptionalAuthMiddleware(mockService))
			router.GET("/test", func(c *gin.Context) {
				_, userExists := c.Get("user")
				if tt.expectUserSet {
					assert.True(t, userExists, "User should be set in context")
				} else {
					assert.False(t, userExists, "User should not be set in context")
				}
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Setup request
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Setup response recorder
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), "success")

			mockService.AssertExpectations(t)
		})
	}
}

func TestRequireAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userType       string
		setupUser      func(*gin.Context)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:     "valid user",
			userType: "user",
			setupUser: func(c *gin.Context) {
				c.Set("user_type", "user")
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:     "valid admin",
			userType: "admin",
			setupUser: func(c *gin.Context) {
				c.Set("user_type", "admin")
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:     "invalid user type",
			userType: "guest",
			setupUser: func(c *gin.Context) {
				c.Set("user_type", "guest")
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:     "no user type set",
			userType: "",
			setupUser: func(c *gin.Context) {
				// Don't set user_type
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup router with middleware
			router := gin.New()
			router.Use(func(c *gin.Context) {
				tt.setupUser(c)
				c.Next()
			})
			router.Use(middlewares.RequireAuthMiddleware())
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "access granted"})
			})

			// Setup request
			req := httptest.NewRequest("GET", "/test", nil)

			// Setup response recorder
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError {
				assert.Contains(t, w.Body.String(), "error")
			} else {
				assert.Contains(t, w.Body.String(), "access granted")
			}
		})
	}
}

func TestExtractTokenFromHeader(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		wantToken string
		wantErr   bool
	}{
		{
			name:      "valid bearer token",
			header:    "Bearer valid.jwt.token",
			wantToken: "valid.jwt.token",
			wantErr:   false,
		},
		{
			name:      "empty header",
			header:    "",
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "invalid format",
			header:    "InvalidFormat valid.jwt.token",
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "no bearer prefix",
			header:    "valid.jwt.token",
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tt.header)

			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = req

			token, err := middlewares.ExtractTokenFromHeader(c)
			if (err != nil) != tt.wantErr {
				t.Errorf("middlewares.ExtractTokenFromHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if token != tt.wantToken {
				t.Errorf("middlewares.ExtractTokenFromHeader() token = %v, want %v", token, tt.wantToken)
			}
		})
	}
}

func TestGetCurrentUser(t *testing.T) {
	tests := []struct {
		name      string
		setupUser func(*gin.Context)
		wantUser  bool
	}{
		{
			name: "user set in context",
			setupUser: func(c *gin.Context) {
				claims := &utils.JWTClaims{
					UserID:   1,
					UserType: "user",
					Email:    "user@example.com",
					Name:     "Test User",
				}
				c.Set("user", claims)
			},
			wantUser: true,
		},
		{
			name: "no user in context",
			setupUser: func(c *gin.Context) {
				// Don't set user
			},
			wantUser: false,
		},
		{
			name: "invalid user type in context",
			setupUser: func(c *gin.Context) {
				c.Set("user", "invalid-type")
			},
			wantUser: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			tt.setupUser(c)

			user, exists := middlewares.GetCurrentUser(c)
			if tt.wantUser {
				assert.True(t, exists)
				assert.NotNil(t, user)
			} else {
				assert.False(t, exists)
				assert.Nil(t, user)
			}
		})
	}
}
