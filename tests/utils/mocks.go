package utils

import (
	"testing"
	"time"

	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/utils"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateLastLogin(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockAdminRepository is a mock implementation of AdminRepository
type MockAdminRepository struct {
	mock.Mock
}

func (m *MockAdminRepository) GetByEmail(email string) (*models.Admin, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Admin), args.Error(1)
}

func (m *MockAdminRepository) UpdateLastLogin(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// CreateMockUser creates a mock user for testing
func CreateMockUser(id int, email, name string) *models.User {
	hashedPassword, _ := models.HashPassword("password123")
	return &models.User{
		ID:       id,
		Name:     name,
		Mobile:   "+1234567890",
		Email:    email,
		Password: hashedPassword,
	}
}

// CreateMockAdmin creates a mock admin for testing
func CreateMockAdmin(id int, email, name string) *models.Admin {
	hashedPassword, _ := models.HashPassword("admin123")
	return &models.Admin{
		ID:       id,
		Name:     name,
		Mobile:   "+1234567890",
		Email:    email,
		Password: hashedPassword,
	}
}

// CreateMockJWTClaims creates mock JWT claims for testing
func CreateMockJWTClaims(userID int, userType, email, name string) interface{} {
	return struct {
		UserID   int    `json:"user_id"`
		UserType string `json:"user_type"`
		Email    string `json:"email"`
		Name     string `json:"name"`
	}{
		UserID:   userID,
		UserType: userType,
		Email:    email,
		Name:     name,
	}
}

// CreateMockLoginResponse creates a mock login response for testing
func CreateMockLoginResponse(userType string, user interface{}) *models.LoginResponse {
	return &models.LoginResponse{
		Token:     "mock.jwt.token",
		UserType:  userType,
		User:      user,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
}

// SetupMockUserRepository sets up a mock user repository with common expectations
func SetupMockUserRepository(t *testing.T) *MockUserRepository {
	mockRepo := new(MockUserRepository)
	return mockRepo
}

// SetupMockAdminRepository sets up a mock admin repository with common expectations
func SetupMockAdminRepository(t *testing.T) *MockAdminRepository {
	mockRepo := new(MockAdminRepository)
	return mockRepo
}

// AssertUserRepositoryExpectations asserts all expectations on the mock user repository
func AssertUserRepositoryExpectations(t *testing.T, mockRepo *MockUserRepository) {
	mockRepo.AssertExpectations(t)
}

// AssertAdminRepositoryExpectations asserts all expectations on the mock admin repository
func AssertAdminRepositoryExpectations(t *testing.T, mockRepo *MockAdminRepository) {
	mockRepo.AssertExpectations(t)
}

// MockAuthService is a mock implementation of AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(email, password string, rememberMe bool) (*models.LoginResponse, error) {
	args := m.Called(email, password, rememberMe)
	return args.Get(0).(*models.LoginResponse), args.Error(1)
}

func (m *MockAuthService) ValidateToken(tokenString string) (*utils.JWTClaims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*utils.JWTClaims), args.Error(1)
}

func (m *MockAuthService) RefreshToken(tokenString string, rememberMe bool) (string, error) {
	args := m.Called(tokenString, rememberMe)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) GetUserFromToken(tokenString string) (*utils.JWTClaims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*utils.JWTClaims), args.Error(1)
}

func (m *MockAuthService) ForgotPassword(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockAuthService) ResetPassword(token, email, newPassword, confirmPassword string) error {
	args := m.Called(token, email, newPassword, confirmPassword)
	return args.Error(0)
}

func (m *MockAuthService) ValidateResetToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}
