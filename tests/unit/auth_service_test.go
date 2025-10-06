package unit

import (
	"testing"

	"github.com/gatehide/gatehide-api/internal/repositories"
	"github.com/gatehide/gatehide-api/internal/services"
	testutils "github.com/gatehide/gatehide-api/tests/utils"
)

func TestAuthService_LoginUser(t *testing.T) {
	testutils.SkipIfNoDB(t)

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	cfg := testutils.TestConfig()
	authService := services.NewAuthService(userRepo, adminRepo, cfg)

	// Create a test user
	_ = testutils.CreateTestUser(t, db, "user@example.com", "password123", "Test User")

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "valid user login",
			email:    "user@example.com",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "invalid password",
			email:    "user@example.com",
			password: "wrongpassword",
			wantErr:  true,
		},
		{
			name:     "non-existing user",
			email:    "nonexistent@example.com",
			password: "password123",
			wantErr:  true,
		},
		{
			name:     "empty email",
			email:    "",
			password: "password123",
			wantErr:  true,
		},
		{
			name:     "empty password",
			email:    "user@example.com",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.Login(tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if response == nil {
					t.Error("AuthService.Login() returned nil response")
					return
				}

				if response.Token == "" {
					t.Error("AuthService.Login() returned empty token")
				}

				if response.UserType != "user" {
					t.Errorf("AuthService.Login() userType = %v, want %v", response.UserType, "user")
				}

				if response.User == nil {
					t.Error("AuthService.Login() returned nil user")
				}
			}
		})
	}
}

func TestAuthService_LoginAdmin(t *testing.T) {
	testutils.SkipIfNoDB(t)

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	cfg := testutils.TestConfig()
	authService := services.NewAuthService(userRepo, adminRepo, cfg)

	// Create a test admin
	_ = testutils.CreateTestAdmin(t, db, "admin@example.com", "admin123", "Test Admin")

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "valid admin login",
			email:    "admin@example.com",
			password: "admin123",
			wantErr:  false,
		},
		{
			name:     "invalid password",
			email:    "admin@example.com",
			password: "wrongpassword",
			wantErr:  true,
		},
		{
			name:     "non-existing admin",
			email:    "nonexistent@example.com",
			password: "admin123",
			wantErr:  true,
		},
		{
			name:     "empty email",
			email:    "",
			password: "admin123",
			wantErr:  true,
		},
		{
			name:     "empty password",
			email:    "admin@example.com",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.Login(tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if response == nil {
					t.Error("AuthService.Login() returned nil response")
					return
				}

				if response.Token == "" {
					t.Error("AuthService.Login() returned empty token")
				}

				if response.UserType != "admin" {
					t.Errorf("AuthService.Login() userType = %v, want %v", response.UserType, "admin")
				}

				if response.User == nil {
					t.Error("AuthService.Login() returned nil user")
				}
			}
		})
	}
}

func TestAuthService_Login_Unified(t *testing.T) {
	testutils.SkipIfNoDB(t)

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	cfg := testutils.TestConfig()
	authService := services.NewAuthService(userRepo, adminRepo, cfg)

	// Create test user and admin
	user := testutils.CreateTestUser(t, db, "user@example.com", "password123", "Test User")
	admin := testutils.CreateTestAdmin(t, db, "admin@example.com", "admin123", "Test Admin")

	tests := []struct {
		name       string
		email      string
		password   string
		wantErr    bool
		expectType string
	}{
		{
			name:       "valid user login",
			email:      user.Email,
			password:   "password123",
			wantErr:    false,
			expectType: "user",
		},
		{
			name:       "valid admin login",
			email:      admin.Email,
			password:   "admin123",
			wantErr:    false,
			expectType: "admin",
		},
		{
			name:     "user email with admin password",
			email:    "user@example.com",
			password: "admin123",
			wantErr:  true,
		},
		{
			name:     "admin email with user password",
			email:    "admin@example.com",
			password: "password123",
			wantErr:  true,
		},
		{
			name:     "non-existing email",
			email:    "nonexistent@example.com",
			password: "anypassword",
			wantErr:  true,
		},
		{
			name:     "empty credentials",
			email:    "",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.Login(tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if response == nil {
					t.Error("AuthService.Login() returned nil response")
					return
				}

				if response.Token == "" {
					t.Error("AuthService.Login() returned empty token")
				}

				if response.UserType != tt.expectType {
					t.Errorf("AuthService.Login() userType = %v, want %v", response.UserType, tt.expectType)
				}

				if response.User == nil {
					t.Error("AuthService.Login() returned nil user")
				}
			}
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	testutils.SkipIfNoDB(t)

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	cfg := testutils.TestConfig()
	authService := services.NewAuthService(userRepo, adminRepo, cfg)

	// Create a test user and get a valid token
	testUser := testutils.CreateTestUser(t, db, "user@example.com", "password123", "Test User")
	loginResponse, err := authService.Login(testUser.Email, "password123")
	if err != nil {
		t.Fatalf("Failed to login user for token validation test: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   loginResponse.Token,
			wantErr: false,
		},
		{
			name:    "invalid token",
			token:   "invalid.token.here",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := authService.ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if claims == nil {
					t.Error("AuthService.ValidateToken() returned nil claims")
					return
				}

				if claims.UserID != testUser.ID {
					t.Errorf("AuthService.ValidateToken() claims.UserID = %v, want %v", claims.UserID, testUser.ID)
				}

				if claims.Email != testUser.Email {
					t.Errorf("AuthService.ValidateToken() claims.Email = %v, want %v", claims.Email, testUser.Email)
				}
			}
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	testutils.SkipIfNoDB(t)

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	cfg := testutils.TestConfig()
	authService := services.NewAuthService(userRepo, adminRepo, cfg)

	// Create a test user and get a valid token
	testUser := testutils.CreateTestUser(t, db, "user@example.com", "password123", "Test User")
	loginResponse, err := authService.Login(testUser.Email, "password123")
	if err != nil {
		t.Fatalf("Failed to login user for token refresh test: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token refresh",
			token:   loginResponse.Token,
			wantErr: false,
		},
		{
			name:    "invalid token refresh",
			token:   "invalid.token.here",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newToken, err := authService.RefreshToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if newToken == "" {
					t.Error("AuthService.RefreshToken() returned empty token")
				}

				if newToken == tt.token {
					t.Error("AuthService.RefreshToken() returned the same token")
				}
			}
		})
	}
}

func TestAuthService_GetUserFromToken(t *testing.T) {
	testutils.SkipIfNoDB(t)

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	cfg := testutils.TestConfig()
	authService := services.NewAuthService(userRepo, adminRepo, cfg)

	// Create a test user and get a valid token
	testUser := testutils.CreateTestUser(t, db, "user@example.com", "password123", "Test User")
	loginResponse, err := authService.Login(testUser.Email, "password123")
	if err != nil {
		t.Fatalf("Failed to login user for GetUserFromToken test: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   loginResponse.Token,
			wantErr: false,
		},
		{
			name:    "invalid token",
			token:   "invalid.token.here",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := authService.GetUserFromToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.GetUserFromToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if claims == nil {
					t.Error("AuthService.GetUserFromToken() returned nil claims")
					return
				}

				if claims.UserID != testUser.ID {
					t.Errorf("AuthService.GetUserFromToken() claims.UserID = %v, want %v", claims.UserID, testUser.ID)
				}

				if claims.Email != testUser.Email {
					t.Errorf("AuthService.GetUserFromToken() claims.Email = %v, want %v", claims.Email, testUser.Email)
				}
			}
		})
	}
}

func TestAuthService_UserTypeDetection(t *testing.T) {
	testutils.SkipIfNoDB(t)

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	cfg := testutils.TestConfig()
	authService := services.NewAuthService(userRepo, adminRepo, cfg)

	// Create test user and admin with same email pattern but different domains
	_ = testutils.CreateTestUser(t, db, "user@example.com", "password123", "Test User")
	_ = testutils.CreateTestAdmin(t, db, "admin@example.com", "admin123", "Test Admin")

	tests := []struct {
		name       string
		email      string
		password   string
		expectType string
		wantErr    bool
	}{
		{
			name:       "user login returns user type",
			email:      "user@example.com",
			password:   "password123",
			expectType: "user",
			wantErr:    false,
		},
		{
			name:       "admin login returns admin type",
			email:      "admin@example.com",
			password:   "admin123",
			expectType: "admin",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.Login(tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if response.UserType != tt.expectType {
					t.Errorf("AuthService.Login() userType = %v, want %v", response.UserType, tt.expectType)
				}
			}
		})
	}
}
