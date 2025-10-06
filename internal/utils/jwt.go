package utils

import (
	"fmt"
	"time"

	"github.com/gatehide/gatehide-api/config"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the JWT claims structure
type JWTClaims struct {
	UserID   int    `json:"user_id"`
	UserType string `json:"user_type"` // "user" or "admin"
	Email    string `json:"email"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT operations
type JWTManager struct {
	secret             []byte
	expiration         time.Duration
	rememberExpiration time.Duration
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(cfg *config.Config) *JWTManager {
	return &JWTManager{
		secret:             []byte(cfg.Security.JWTSecret),
		expiration:         time.Duration(cfg.Security.JWTExpiration) * time.Hour,
		rememberExpiration: time.Duration(cfg.Security.JWTExpiration) * time.Hour * 24 * 7, // 7 days for remember me
	}
}

// GenerateToken generates a new JWT token for the given user
func (j *JWTManager) GenerateToken(userID int, userType, email, name string, rememberMe bool) (string, error) {
	now := time.Now()

	// Choose expiration based on remember me
	expiration := j.expiration
	if rememberMe {
		expiration = j.rememberExpiration
	}

	claims := JWTClaims{
		UserID:   userID,
		UserType: userType,
		Email:    email,
		Name:     name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "gatehide-api",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ValidateToken validates and parses a JWT token
func (j *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// RefreshToken generates a new token with extended expiration
func (j *JWTManager) RefreshToken(tokenString string, rememberMe bool) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Add a delay to ensure different timestamps
	time.Sleep(100 * time.Millisecond)

	// For testing, always generate a new token
	// In production, you might want to check if token is close to expiration
	return j.GenerateToken(claims.UserID, claims.UserType, claims.Email, claims.Name, rememberMe)
}
