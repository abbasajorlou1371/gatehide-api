package unit

import (
	"context"
	"errors"
	"testing"

	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAll() ([]models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetAllByGamenet(gamenetID int) ([]models.User, error) {
	args := m.Called(gamenetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByMobile(mobile string) (*models.User, error) {
	args := m.Called(mobile)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(id int, user *models.UserUpdateRequest) error {
	args := m.Called(id, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) Search(req *models.UserSearchRequest) (*models.UserSearchResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserSearchResponse), args.Error(1)
}

func (m *MockUserRepository) SearchByGamenet(req *models.UserSearchRequest, gamenetID int) (*models.UserSearchResponse, error) {
	args := m.Called(req, gamenetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserSearchResponse), args.Error(1)
}

func (m *MockUserRepository) LinkToGamenet(userID, gamenetID int) error {
	args := m.Called(userID, gamenetID)
	return args.Error(0)
}

func (m *MockUserRepository) UnlinkFromGamenet(userID, gamenetID int) error {
	args := m.Called(userID, gamenetID)
	return args.Error(0)
}

func (m *MockUserRepository) GetGamenetIDByUser(userID int) (*int, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockUserRepository) UpdateLastLogin(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePassword(id int, hashedPassword string) error {
	args := m.Called(id, hashedPassword)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateProfile(id int, name, mobile, image string) error {
	args := m.Called(id, name, mobile, image)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateEmail(id int, email string) error {
	args := m.Called(id, email)
	return args.Error(0)
}

// MockPermissionRepository is a mock implementation of PermissionRepositoryInterface
type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) GetPermissionsByRole(roleType string) ([]models.Permission, error) {
	args := m.Called(roleType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Permission), args.Error(1)
}

func (m *MockPermissionRepository) HasPermission(roleType, resource, action string) (bool, error) {
	args := m.Called(roleType, resource, action)
	return args.Bool(0), args.Error(1)
}

func (m *MockPermissionRepository) GetRoleWithPermissions(roleType string) (*models.RoleWithPermissions, error) {
	args := m.Called(roleType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RoleWithPermissions), args.Error(1)
}

func (m *MockPermissionRepository) GetRoleByName(roleName string) (*models.Role, error) {
	args := m.Called(roleName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockPermissionRepository) GetAllRoles() ([]models.Role, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockPermissionRepository) GetAllPermissions() ([]models.Permission, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Permission), args.Error(1)
}

func (m *MockPermissionRepository) AssignRoleToUser(userID int, userType string, roleName string) error {
	args := m.Called(userID, userType, roleName)
	return args.Error(0)
}

func (m *MockPermissionRepository) GetUserRoles(userID int, userType string) ([]models.Role, error) {
	args := m.Called(userID, userType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockPermissionRepository) GetUserPermissions(userID int, userType string) ([]models.Permission, error) {
	args := m.Called(userID, userType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Permission), args.Error(1)
}

func (m *MockPermissionRepository) RemoveRoleFromUser(userID int, userType string, roleName string) error {
	args := m.Called(userID, userType, roleName)
	return args.Error(0)
}

func (m *MockPermissionRepository) HasUserRole(userID int, userType string, roleName string) (bool, error) {
	args := m.Called(userID, userType, roleName)
	return args.Bool(0), args.Error(1)
}

// MockSMSService is a mock implementation of SMSService
type MockSMSService struct {
	mock.Mock
}

func (m *MockSMSService) SendSMS(ctx context.Context, sms *models.SMSNotification) error {
	args := m.Called(ctx, sms)
	return args.Error(0)
}

func (m *MockSMSService) SendUserCredentials(ctx context.Context, mobile, email, password string) error {
	args := m.Called(ctx, mobile, email, password)
	return args.Error(0)
}

func TestUserService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		userService := services.NewUserService(mockRepo, mockPermissionRepo, nil, nil)

		req := &models.UserCreateRequest{
			Name:   "Test User",
			Email:  "test@example.com",
			Mobile: "09123456789",
		}

		// Mock GetByEmail to return not found
		mockRepo.On("GetByEmail", req.Email).Return(nil, errors.New("user not found"))
		// Mock GetByMobile to return not found
		mockRepo.On("GetByMobile", req.Mobile).Return(nil, errors.New("user not found"))
		// Mock Create to succeed
		mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)
		// Mock AssignRoleToUser to succeed
		mockPermissionRepo.On("AssignRoleToUser", mock.AnythingOfType("int"), "user", "user").Return(nil)

		user, err := userService.Create(ctx, req, nil)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Name, user.Name)
		assert.Equal(t, req.Email, user.Email)
		assert.Equal(t, req.Mobile, user.Mobile)
		mockRepo.AssertExpectations(t)
		mockPermissionRepo.AssertExpectations(t)
	})

	t.Run("Email Already Exists", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		userService := services.NewUserService(mockRepo, mockPermissionRepo, nil, nil)

		req := &models.UserCreateRequest{
			Name:   "Test User",
			Email:  "existing@example.com",
			Mobile: "09123456789",
		}

		existingUser := &models.User{
			ID:    1,
			Email: req.Email,
		}

		// Mock GetByEmail to return existing user
		mockRepo.On("GetByEmail", req.Email).Return(existingUser, nil)

		user, err := userService.Create(ctx, req, nil)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "email already exists")
		mockRepo.AssertExpectations(t)
		mockPermissionRepo.AssertExpectations(t)
	})

	t.Run("Mobile Already Exists", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		userService := services.NewUserService(mockRepo, mockPermissionRepo, nil, nil)

		req := &models.UserCreateRequest{
			Name:   "Test User",
			Email:  "test@example.com",
			Mobile: "09123456789",
		}

		existingUser := &models.User{
			ID:     1,
			Mobile: req.Mobile,
		}

		// Mock GetByEmail to return not found
		mockRepo.On("GetByEmail", req.Email).Return(nil, errors.New("user not found"))
		// Mock GetByMobile to return existing user
		mockRepo.On("GetByMobile", req.Mobile).Return(existingUser, nil)

		user, err := userService.Create(ctx, req, nil)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "mobile number already exists")
		mockRepo.AssertExpectations(t)
		mockPermissionRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		userService := services.NewUserService(mockRepo, mockPermissionRepo, nil, nil)

		expectedUser := &models.User{
			ID:     1,
			Name:   "Test User",
			Email:  "test@example.com",
			Mobile: "09123456789",
		}

		mockRepo.On("GetByID", 1).Return(expectedUser, nil)

		user, err := userService.GetByID(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Name, user.Name)
		mockRepo.AssertExpectations(t)
		mockPermissionRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		userService := services.NewUserService(mockRepo, mockPermissionRepo, nil, nil)

		mockRepo.On("GetByID", 999).Return(nil, errors.New("user not found"))

		user, err := userService.GetByID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
		mockPermissionRepo.AssertExpectations(t)
	})
}

func TestUserService_Update(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		userService := services.NewUserService(mockRepo, mockPermissionRepo, nil, nil)

		existingUser := &models.User{
			ID:     1,
			Name:   "Old Name",
			Email:  "old@example.com",
			Mobile: "09123456789",
		}

		newName := "New Name"
		req := &models.UserUpdateRequest{
			Name: &newName,
		}

		mockRepo.On("GetByID", 1).Return(existingUser, nil).Once()
		mockRepo.On("Update", 1, req).Return(nil)

		// Mock for getting updated user
		updatedUser := &models.User{
			ID:     1,
			Name:   newName,
			Email:  existingUser.Email,
			Mobile: existingUser.Mobile,
		}
		mockRepo.On("GetByID", 1).Return(updatedUser, nil).Once()

		user, err := userService.Update(ctx, 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, newName, user.Name)
		mockRepo.AssertExpectations(t)
		mockPermissionRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		userService := services.NewUserService(mockRepo, mockPermissionRepo, nil, nil)

		newName := "New Name"
		req := &models.UserUpdateRequest{
			Name: &newName,
		}

		mockRepo.On("GetByID", 999).Return(nil, errors.New("user not found"))

		user, err := userService.Update(ctx, 999, req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "not found")
		mockRepo.AssertExpectations(t)
		mockPermissionRepo.AssertExpectations(t)
	})
}

func TestUserService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		userService := services.NewUserService(mockRepo, mockPermissionRepo, nil, nil)

		existingUser := &models.User{
			ID:   1,
			Name: "Test User",
		}

		mockRepo.On("GetByID", 1).Return(existingUser, nil)
		mockRepo.On("Delete", 1).Return(nil)

		err := userService.Delete(ctx, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockPermissionRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		userService := services.NewUserService(mockRepo, mockPermissionRepo, nil, nil)

		mockRepo.On("GetByID", 999).Return(nil, errors.New("user not found"))

		err := userService.Delete(ctx, 999)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		mockRepo.AssertExpectations(t)
		mockPermissionRepo.AssertExpectations(t)
	})
}

func TestUserService_GetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		userService := services.NewUserService(mockRepo, mockPermissionRepo, nil, nil)

		expectedUsers := []models.User{
			{ID: 1, Name: "User 1", Email: "user1@example.com"},
			{ID: 2, Name: "User 2", Email: "user2@example.com"},
		}

		mockRepo.On("GetAll").Return(expectedUsers, nil)

		users, err := userService.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Len(t, users, 2)
		mockRepo.AssertExpectations(t)
		mockPermissionRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		userService := services.NewUserService(mockRepo, mockPermissionRepo, nil, nil)

		mockRepo.On("GetAll").Return(nil, errors.New("database error"))

		users, err := userService.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, users)
		mockRepo.AssertExpectations(t)
		mockPermissionRepo.AssertExpectations(t)
	})
}

func TestUserService_Search(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		userService := services.NewUserService(mockRepo, mockPermissionRepo, nil, nil)

		searchReq := &models.UserSearchRequest{
			Query:    "test",
			Page:     1,
			PageSize: 10,
		}

		expectedResponse := &models.UserSearchResponse{
			Data: []models.UserResponse{
				{ID: 1, Name: "Test User", Email: "test@example.com"},
			},
			Pagination: models.PaginationInfo{
				CurrentPage: 1,
				PageSize:    10,
				TotalItems:  1,
				TotalPages:  1,
				HasNext:     false,
				HasPrev:     false,
			},
		}

		mockRepo.On("Search", searchReq).Return(expectedResponse, nil)

		result, err := userService.Search(ctx, searchReq)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Data, 1)
		mockRepo.AssertExpectations(t)
		mockPermissionRepo.AssertExpectations(t)
	})
}
