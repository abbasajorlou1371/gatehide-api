package utils

import (
	"context"
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

func (m *MockUserRepository) GetByID(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateLastLogin(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePassword(id int, hashedPassword string) error {
	args := m.Called(id, hashedPassword)
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

func (m *MockAdminRepository) GetByID(id int) (*models.Admin, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Admin), args.Error(1)
}

func (m *MockAdminRepository) UpdateLastLogin(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAdminRepository) UpdatePassword(id int, hashedPassword string) error {
	args := m.Called(id, hashedPassword)
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

func (m *MockAuthService) ChangePassword(userID int, userType, currentPassword, newPassword, confirmPassword string) error {
	args := m.Called(userID, userType, currentPassword, newPassword, confirmPassword)
	return args.Error(0)
}

func (m *MockAuthService) LoginWithSession(email, password string, rememberMe bool, deviceInfo, ipAddress, userAgent string) (*models.LoginResponse, error) {
	args := m.Called(email, password, rememberMe, deviceInfo, ipAddress, userAgent)
	return args.Get(0).(*models.LoginResponse), args.Error(1)
}

func (m *MockAuthService) GetUserByID(userID int) (*models.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) GetAdminByID(adminID int) (*models.Admin, error) {
	args := m.Called(adminID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Admin), args.Error(1)
}

func (m *MockAuthService) GetGamenetByID(gamenetID int) (*models.Gamenet, error) {
	args := m.Called(gamenetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Gamenet), args.Error(1)
}

func (m *MockAuthService) UpdateUserProfile(userID int, name, mobile, image string) (*models.UserResponse, error) {
	args := m.Called(userID, name, mobile, image)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockAuthService) UpdateAdminProfile(adminID int, name, mobile, image string) (*models.AdminResponse, error) {
	args := m.Called(adminID, name, mobile, image)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AdminResponse), args.Error(1)
}

func (m *MockAuthService) UpdateGamenetProfile(gamenetID int, name, mobile, image string) (*models.GamenetResponse, error) {
	args := m.Called(gamenetID, name, mobile, image)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GamenetResponse), args.Error(1)
}

func (m *MockAuthService) UpdateUserEmail(userID int, newEmail string) (*models.UserResponse, error) {
	args := m.Called(userID, newEmail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockAuthService) UpdateAdminEmail(adminID int, newEmail string) (*models.AdminResponse, error) {
	args := m.Called(adminID, newEmail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AdminResponse), args.Error(1)
}

func (m *MockAuthService) UpdateGamenetEmail(gamenetID int, newEmail string) (*models.GamenetResponse, error) {
	args := m.Called(gamenetID, newEmail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GamenetResponse), args.Error(1)
}

func (m *MockAuthService) SendEmailVerification(userID int, userType, newEmail string) (string, error) {
	args := m.Called(userID, userType, newEmail)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) VerifyEmailCode(userID int, userType, email, code string) (bool, error) {
	args := m.Called(userID, userType, email, code)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthService) CheckEmailExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

// MockSessionRepository is a mock implementation of SessionRepositoryInterface
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) CreateSession(session *models.UserSession) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *MockSessionRepository) GetSessionByToken(token string) (*models.UserSession, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserSession), args.Error(1)
}

func (m *MockSessionRepository) GetActiveSessionsByUserID(userID int, userType string) ([]models.UserSession, error) {
	args := m.Called(userID, userType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserSession), args.Error(1)
}

func (m *MockSessionRepository) DeactivateSession(sessionID int) error {
	args := m.Called(sessionID)
	return args.Error(0)
}

func (m *MockSessionRepository) DeactivateAllOtherUserSessions(userID int, userType string, currentSessionToken string) error {
	args := m.Called(userID, userType, currentSessionToken)
	return args.Error(0)
}

func (m *MockSessionRepository) DeactivateAllUserSessions(userID int, userType string) error {
	args := m.Called(userID, userType)
	return args.Error(0)
}

func (m *MockSessionRepository) UpdateSessionActivity(sessionID int) error {
	args := m.Called(sessionID)
	return args.Error(0)
}

func (m *MockSessionRepository) CleanupExpiredSessions() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteSession(sessionID int) error {
	args := m.Called(sessionID)
	return args.Error(0)
}

// MockNotificationService is a mock implementation of NotificationServiceInterface
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) SendNotification(ctx context.Context, notification *models.CreateNotificationRequest) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

func (m *MockNotificationService) SendEmail(ctx context.Context, email *models.SendEmailRequest) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockNotificationService) SendSMS(ctx context.Context, sms *models.SendSMSRequest) error {
	args := m.Called(ctx, sms)
	return args.Error(0)
}

func (m *MockNotificationService) SendDatabaseNotification(ctx context.Context, dbNotification *models.DatabaseNotification) error {
	args := m.Called(ctx, dbNotification)
	return args.Error(0)
}

func (m *MockNotificationService) GetNotification(ctx context.Context, id int) (*models.Notification, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Notification), args.Error(1)
}

func (m *MockNotificationService) GetNotifications(ctx context.Context, filters map[string]interface{}) ([]*models.Notification, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Notification), args.Error(1)
}

func (m *MockNotificationService) UpdateNotificationStatus(ctx context.Context, id int, status models.NotificationStatus, errorMsg *string) error {
	args := m.Called(ctx, id, status, errorMsg)
	return args.Error(0)
}

func (m *MockNotificationService) RetryFailedNotification(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockFileUploader is a mock implementation of FileUploader
type MockFileUploader struct {
	mock.Mock
}

func (m *MockFileUploader) UploadFile(file interface{}, subfolder string) (interface{}, error) {
	args := m.Called(file, subfolder)
	return args.Get(0), args.Error(1)
}

func (m *MockFileUploader) DeleteFile(filePath string) error {
	args := m.Called(filePath)
	return args.Error(0)
}

func (m *MockFileUploader) GetFileInfo(filePath string) (interface{}, error) {
	args := m.Called(filePath)
	return args.Get(0), args.Error(1)
}

// MockSubscriptionPlanRepository is a mock implementation of SubscriptionPlanRepositoryInterface
type MockSubscriptionPlanRepository struct {
	mock.Mock
}

func (m *MockSubscriptionPlanRepository) Create(plan *models.SubscriptionPlan) error {
	args := m.Called(plan)
	return args.Error(0)
}

func (m *MockSubscriptionPlanRepository) GetByID(id int) (*models.SubscriptionPlan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionPlanRepository) GetAll(limit, offset int, isActive *bool) ([]*models.SubscriptionPlan, error) {
	args := m.Called(limit, offset, isActive)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionPlanRepository) Update(id int, plan *models.SubscriptionPlan) error {
	args := m.Called(id, plan)
	return args.Error(0)
}

func (m *MockSubscriptionPlanRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSubscriptionPlanRepository) Count(isActive *bool) (int, error) {
	args := m.Called(isActive)
	return args.Int(0), args.Error(1)
}

func (m *MockSubscriptionPlanRepository) HasActiveSubscriptions(planID int) (bool, error) {
	args := m.Called(planID)
	return args.Bool(0), args.Error(1)
}

// CreateMockSubscriptionPlan creates a mock subscription plan for testing
func CreateMockSubscriptionPlan(id int, name, planType string, price float64) *models.SubscriptionPlan {
	now := time.Now()
	return &models.SubscriptionPlan{
		ID:                       id,
		Name:                     name,
		PlanType:                 planType,
		Price:                    price,
		AnnualDiscountPercentage: nil,
		TrialDurationDays:        nil,
		IsActive:                 true,
		SubscriptionCount:        0,
		CreatedAt:                now,
		UpdatedAt:                now,
	}
}

// CreateMockPlanResponse creates a mock plan response for testing
func CreateMockPlanResponse(id int, name, planType string, price float64) *models.PlanResponse {
	now := time.Now()
	return &models.PlanResponse{
		ID:                       id,
		Name:                     name,
		PlanType:                 planType,
		Price:                    price,
		AnnualDiscountPercentage: nil,
		TrialDurationDays:        nil,
		IsActive:                 true,
		CreatedAt:                now,
		UpdatedAt:                now,
	}
}

// SetupMockSubscriptionPlanRepository sets up a mock subscription plan repository
func SetupMockSubscriptionPlanRepository(t *testing.T) *MockSubscriptionPlanRepository {
	mockRepo := new(MockSubscriptionPlanRepository)
	return mockRepo
}

// AssertSubscriptionPlanRepositoryExpectations asserts all expectations on the mock repository
func AssertSubscriptionPlanRepositoryExpectations(t *testing.T, mockRepo *MockSubscriptionPlanRepository) {
	mockRepo.AssertExpectations(t)
}
