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
	defer db.Close()
	defer testutils.CleanupTestDB(t, db)

	// Ensure database is clean before starting
	testutils.CleanupTestDB(t, db)

	userRepo := repositories.NewUserRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	cfg := testutils.TestConfig()
	authService := services.NewAuthService(userRepo, adminRepo, cfg)

	// Create a test user with unique email
	_ = testutils.CreateTestUser(t, db, "user1@example.com", "password123", "Test User 1")

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "valid user login",
			email:    "user1@example.com",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "invalid password",
			email:    "user1@example.com",
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
			email:    "user1@example.com",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.Login(tt.email, tt.password, false)
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
	defer db.Close()
	defer testutils.CleanupTestDB(t, db)

	userRepo := repositories.NewUserRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	cfg := testutils.TestConfig()
	authService := services.NewAuthService(userRepo, adminRepo, cfg)

	// Create a test admin with unique email
	_ = testutils.CreateTestAdmin(t, db, "admin1@example.com", "admin123", "Test Admin 1")

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "valid admin login",
			email:    "admin1@example.com",
			password: "admin123",
			wantErr:  false,
		},
		{
			name:     "invalid password",
			email:    "admin1@example.com",
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
			email:    "admin1@example.com",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.Login(tt.email, tt.password, false)
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
	defer db.Close()
	defer testutils.CleanupTestDB(t, db)

	userRepo := repositories.NewUserRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	cfg := testutils.TestConfig()
	authService := services.NewAuthService(userRepo, adminRepo, cfg)

	// Create test user and admin with unique emails
	user := testutils.CreateTestUser(t, db, "user2@example.com", "password123", "Test User 2")
	admin := testutils.CreateTestAdmin(t, db, "admin2@example.com", "admin123", "Test Admin 2")

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
			email:    "user2@example.com",
			password: "admin123",
			wantErr:  true,
		},
		{
			name:     "admin email with user password",
			email:    "admin2@example.com",
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
			response, err := authService.Login(tt.email, tt.password, false)
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
	defer db.Close()
	defer testutils.CleanupTestDB(t, db)

	userRepo := repositories.NewUserRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	cfg := testutils.TestConfig()
	authService := services.NewAuthService(userRepo, adminRepo, cfg)

	// Create a test user and get a valid token
	testUser := testutils.CreateTestUser(t, db, "user3@example.com", "password123", "Test User 3")
	loginResponse, err := authService.Login(testUser.Email, "password123", false)
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
	defer db.Close()
	defer testutils.CleanupTestDB(t, db)

	userRepo := repositories.NewUserRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	cfg := testutils.TestConfig()
	authService := services.NewAuthService(userRepo, adminRepo, cfg)

	// Create a test user and get a valid token
	testUser := testutils.CreateTestUser(t, db, "user4@example.com", "password123", "Test User 4")
	loginResponse, err := authService.Login(testUser.Email, "password123", false)
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
			newToken, err := authService.RefreshToken(tt.token, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if newToken == "" {
					t.Error("AuthService.RefreshToken() returned empty token")
				}

				// Check that the refreshed token is valid by trying to validate it
				// This is more realistic than checking if tokens are different
				claims, err := authService.ValidateToken(newToken)
				if err != nil {
					t.Errorf("AuthService.RefreshToken() returned invalid token: %v", err)
				}
				if claims == nil {
					t.Error("AuthService.RefreshToken() returned token with no claims")
				}
			}
		})
	}
}

func TestAuthService_GetUserFromToken(t *testing.T) {
	testutils.SkipIfNoDB(t)

	db := testutils.SetupTestDB(t)
	defer db.Close()
	defer testutils.CleanupTestDB(t, db)

	userRepo := repositories.NewUserRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	cfg := testutils.TestConfig()
	authService := services.NewAuthService(userRepo, adminRepo, cfg)

	// Create a test user and get a valid token
	testUser := testutils.CreateTestUser(t, db, "user5@example.com", "password123", "Test User 5")
	loginResponse, err := authService.Login(testUser.Email, "password123", false)
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
	defer db.Close()
	defer testutils.CleanupTestDB(t, db)

	userRepo := repositories.NewUserRepository(db)
	adminRepo := repositories.NewAdminRepository(db)
	cfg := testutils.TestConfig()
	authService := services.NewAuthService(userRepo, adminRepo, cfg)

	// Create test user and admin with unique emails
	_ = testutils.CreateTestUser(t, db, "user6@example.com", "password123", "Test User 6")
	_ = testutils.CreateTestAdmin(t, db, "admin3@example.com", "admin123", "Test Admin 3")

	tests := []struct {
		name       string
		email      string
		password   string
		expectType string
		wantErr    bool
	}{
		{
			name:       "user login returns user type",
			email:      "user6@example.com",
			password:   "password123",
			expectType: "user",
			wantErr:    false,
		},
		{
			name:       "admin login returns admin type",
			email:      "admin3@example.com",
			password:   "admin123",
			expectType: "admin",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.Login(tt.email, tt.password, false)
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
