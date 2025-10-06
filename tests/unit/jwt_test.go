package unit

import (
	"testing"
	"time"

	"github.com/gatehide/gatehide-api/internal/utils"
	testutils "github.com/gatehide/gatehide-api/tests/utils"
)

func TestJWTManager_GenerateToken(t *testing.T) {
	cfg := testutils.TestConfig()
	jwtManager := utils.NewJWTManager(cfg)

	tests := []struct {
		name     string
		userID   int
		userType string
		email    string
		userName string
		wantErr  bool
	}{
		{
			name:     "valid user token",
			userID:   1,
			userType: "user",
			email:    "user@example.com",
			userName: "Test User",
			wantErr:  false,
		},
		{
			name:     "valid admin token",
			userID:   1,
			userType: "admin",
			email:    "admin@example.com",
			userName: "Test Admin",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := jwtManager.GenerateToken(tt.userID, tt.userType, tt.email, tt.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("JWTManager.GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && token == "" {
				t.Error("JWTManager.GenerateToken() returned empty token")
			}
		})
	}
}

func TestJWTManager_ValidateToken(t *testing.T) {
	cfg := testutils.TestConfig()
	jwtManager := utils.NewJWTManager(cfg)

	// Generate a valid token
	token, err := jwtManager.GenerateToken(1, "user", "test@example.com", "Test User")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   token,
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
		{
			name:    "malformed token",
			token:   "not.a.token",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := jwtManager.ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("JWTManager.ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if claims == nil {
					t.Error("JWTManager.ValidateToken() returned nil claims for valid token")
					return
				}

				if claims.UserID != 1 {
					t.Errorf("JWTManager.ValidateToken() claims.UserID = %v, want %v", claims.UserID, 1)
				}

				if claims.UserType != "user" {
					t.Errorf("JWTManager.ValidateToken() claims.UserType = %v, want %v", claims.UserType, "user")
				}

				if claims.Email != "test@example.com" {
					t.Errorf("JWTManager.ValidateToken() claims.Email = %v, want %v", claims.Email, "test@example.com")
				}

				if claims.Name != "Test User" {
					t.Errorf("JWTManager.ValidateToken() claims.Name = %v, want %v", claims.Name, "Test User")
				}
			}
		})
	}
}

func TestJWTManager_RefreshToken(t *testing.T) {
	cfg := testutils.TestConfig()
	jwtManager := utils.NewJWTManager(cfg)

	// Generate a token that's close to expiration (within 1 hour)
	token, err := jwtManager.GenerateToken(1, "user", "test@example.com", "Test User")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token refresh (not close to expiration)",
			token:   token,
			wantErr: false, // Token is valid but not close to expiration, so it will return a new token
		},
		{
			name:    "invalid token refresh",
			token:   "invalid.token.here",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newToken, err := jwtManager.RefreshToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("JWTManager.RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && newToken == "" {
				t.Error("JWTManager.RefreshToken() returned empty token")
			}

			if !tt.wantErr && newToken == tt.token {
				// For valid tokens that are not close to expiration, they might return the same token
				// This is expected behavior as per the refresh logic
				t.Log("JWTManager.RefreshToken() returned the same token (expected for tokens not close to expiration)")
			}
		})
	}
}

func TestJWTManager_TokenExpiration(t *testing.T) {
	cfg := testutils.TestConfig()
	jwtManager := utils.NewJWTManager(cfg)

	token, err := jwtManager.GenerateToken(1, "user", "test@example.com", "Test User")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := jwtManager.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Check that expiration is set to approximately 1 hour from now
	expectedExpiration := time.Now().Add(time.Duration(cfg.Security.JWTExpiration) * time.Hour)
	timeDiff := expectedExpiration.Sub(claims.ExpiresAt.Time)

	// Allow 5 minutes tolerance
	if timeDiff > 5*time.Minute || timeDiff < -5*time.Minute {
		t.Errorf("Token expiration time is incorrect. Expected around %v, got %v", expectedExpiration, claims.ExpiresAt.Time)
	}
}

func TestJWTManager_DifferentUserTypes(t *testing.T) {
	cfg := testutils.TestConfig()
	jwtManager := utils.NewJWTManager(cfg)

	userTypes := []string{"user", "admin"}

	for _, userType := range userTypes {
		t.Run("user_type_"+userType, func(t *testing.T) {
			token, err := jwtManager.GenerateToken(1, userType, "test@example.com", "Test User")
			if err != nil {
				t.Fatalf("Failed to generate token for user type %s: %v", userType, err)
			}

			claims, err := jwtManager.ValidateToken(token)
			if err != nil {
				t.Fatalf("Failed to validate token for user type %s: %v", userType, err)
			}

			if claims.UserType != userType {
				t.Errorf("Expected user type %s, got %s", userType, claims.UserType)
			}
		})
	}
}

func TestJWTManager_TokenStructure(t *testing.T) {
	cfg := testutils.TestConfig()
	jwtManager := utils.NewJWTManager(cfg)

	userID := 123
	userType := "user"
	email := "test@example.com"
	userName := "Test User"

	token, err := jwtManager.GenerateToken(userID, userType, email, userName)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := jwtManager.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Verify all claims are correctly set
	if claims.UserID != userID {
		t.Errorf("UserID mismatch: expected %d, got %d", userID, claims.UserID)
	}

	if claims.UserType != userType {
		t.Errorf("UserType mismatch: expected %s, got %s", userType, claims.UserType)
	}

	if claims.Email != email {
		t.Errorf("Email mismatch: expected %s, got %s", email, claims.Email)
	}

	if claims.Name != userName {
		t.Errorf("Name mismatch: expected %s, got %s", userName, claims.Name)
	}

	if claims.Issuer != "gatehide-api" {
		t.Errorf("Issuer mismatch: expected gatehide-api, got %s", claims.Issuer)
	}

	if claims.Subject != "123" {
		t.Errorf("Subject mismatch: expected 123, got %s", claims.Subject)
	}
}
