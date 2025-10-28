package services

import (
	"context"
	"fmt"

	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/repositories"
	"github.com/gatehide/gatehide-api/internal/utils"
)

// userService implements UserServiceInterface
type userService struct {
	userRepo       repositories.UserRepository
	permissionRepo repositories.PermissionRepositoryInterface
	smsService     *SMSService
	emailService   *EmailService
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepository, permissionRepo repositories.PermissionRepositoryInterface, smsService *SMSService, emailService *EmailService) UserServiceInterface {
	return &userService{
		userRepo:       userRepo,
		permissionRepo: permissionRepo,
		smsService:     smsService,
		emailService:   emailService,
	}
}

// GetAll retrieves all users
func (s *userService) GetAll(ctx context.Context) ([]models.UserResponse, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	return responses, nil
}

// GetAllByGamenet retrieves all users for a specific gamenet
func (s *userService) GetAllByGamenet(ctx context.Context, gamenetID int) ([]models.UserResponse, error) {
	users, err := s.userRepo.GetAllByGamenet(gamenetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	return responses, nil
}

// GetByID retrieves a user by ID
func (s *userService) GetByID(ctx context.Context, id int) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

// GetByEmail retrieves a user by email
func (s *userService) GetByEmail(ctx context.Context, email string) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

// GetByMobile retrieves a user by mobile
func (s *userService) GetByMobile(ctx context.Context, mobile string) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByMobile(mobile)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

// Create creates a new user
func (s *userService) Create(ctx context.Context, req *models.UserCreateRequest, gamenetID *int) (*models.UserResponse, error) {
	// Check if user with email already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	// Check if user with mobile already exists
	existingUser, err = s.userRepo.GetByMobile(req.Mobile)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with this mobile number already exists")
	}

	// Generate random 8-digit password
	randomPassword, err := utils.GenerateRandomPassword()
	if err != nil {
		return nil, fmt.Errorf("failed to generate password: %w", err)
	}

	// Hash the password
	hashedPassword, err := models.HashPassword(randomPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Password: hashedPassword,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Assign user role to the newly created user
	err = s.permissionRepo.AssignRoleToUser(user.ID, "user", "user")
	if err != nil {
		// Log error but don't fail creation
		fmt.Printf("Warning: Failed to assign user role to user %d: %v\n", user.ID, err)
	}

	// Link user to gamenet if gamenetID is provided
	if gamenetID != nil && *gamenetID > 0 {
		err = s.userRepo.LinkToGamenet(user.ID, *gamenetID)
		if err != nil {
			// Log error but don't fail creation
			fmt.Printf("Warning: Failed to link user to gamenet: %v\n", err)
		}
	}

	// Send credentials via SMS using Kavenegar Verify Lookup
	if s.smsService != nil {
		err = s.smsService.SendUserCredentials(ctx, req.Mobile, req.Email, randomPassword)
		if err != nil {
			// Log the error but don't fail the creation
			fmt.Printf("Warning: Failed to send credentials SMS to %s: %v\n", req.Mobile, err)
		} else {
			fmt.Printf("Successfully sent credentials SMS to %s\n", req.Mobile)
		}
	}

	response := user.ToResponse()
	return &response, nil
}

// Update updates an existing user
func (s *userService) Update(ctx context.Context, id int, req *models.UserUpdateRequest) (*models.UserResponse, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// If email is being updated, check if it's already taken by another user
	if req.Email != nil {
		existingUser, err := s.userRepo.GetByEmail(*req.Email)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, fmt.Errorf("user with this email already exists")
		}
	}

	// If mobile is being updated, check if it's already taken by another user
	if req.Mobile != nil {
		existingUser, err := s.userRepo.GetByMobile(*req.Mobile)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, fmt.Errorf("user with this mobile number already exists")
		}
	}

	err = s.userRepo.Update(id, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Get updated user
	updatedUser, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated user: %w", err)
	}

	response := updatedUser.ToResponse()
	return &response, nil
}

// Delete deletes a user
func (s *userService) Delete(ctx context.Context, id int) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	err = s.userRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// Search searches users with pagination
func (s *userService) Search(ctx context.Context, req *models.UserSearchRequest) (*models.UserSearchResponse, error) {
	result, err := s.userRepo.Search(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return result, nil
}

// SearchByGamenet searches users for a specific gamenet with pagination
func (s *userService) SearchByGamenet(ctx context.Context, req *models.UserSearchRequest, gamenetID int) (*models.UserSearchResponse, error) {
	result, err := s.userRepo.SearchByGamenet(req, gamenetID)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return result, nil
}

// AttachToGamenet attaches a user to a gamenet
func (s *userService) AttachToGamenet(ctx context.Context, userID, gamenetID int) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	err = s.userRepo.LinkToGamenet(userID, gamenetID)
	if err != nil {
		return fmt.Errorf("failed to attach user to gamenet: %w", err)
	}

	return nil
}

// DetachFromGamenet detaches a user from a gamenet
func (s *userService) DetachFromGamenet(ctx context.Context, userID, gamenetID int) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	err = s.userRepo.UnlinkFromGamenet(userID, gamenetID)
	if err != nil {
		return fmt.Errorf("failed to detach user from gamenet: %w", err)
	}

	return nil
}

// CanModifyUser checks if a requester can modify a user
func (s *userService) CanModifyUser(ctx context.Context, userID, requesterID int, requesterType string) (bool, error) {
	// Admins can modify any user
	if requesterType == "admin" {
		return true, nil
	}

	// Gamenets can only modify users they created (first gamenet to link the user)
	if requesterType == "gamenet" {
		creatorGamenetID, err := s.userRepo.GetGamenetIDByUser(userID)
		if err != nil {
			return false, err
		}

		if creatorGamenetID == nil {
			return false, nil
		}

		return *creatorGamenetID == requesterID, nil
	}

	return false, nil
}

// ResendCredentials generates new password and sends credentials via SMS
func (s *userService) ResendCredentials(ctx context.Context, id int) error {
	// Get user details
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Generate new random 8-digit password
	newPassword, err := utils.GenerateRandomPassword()
	if err != nil {
		return fmt.Errorf("failed to generate password: %w", err)
	}

	// Hash the new password
	hashedPassword, err := models.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password in database
	err = s.userRepo.UpdatePassword(id, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Send credentials via SMS
	if s.smsService != nil {
		err = s.smsService.SendUserCredentials(ctx, user.Mobile, user.Email, newPassword)
		if err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Warning: Failed to send credentials SMS to %s: %v\n", user.Mobile, err)
			return fmt.Errorf("password updated but failed to send SMS: %w", err)
		}

		fmt.Printf("Successfully sent credentials SMS to %s\n", user.Mobile)
	} else {
		return fmt.Errorf("SMS service not configured")
	}

	return nil
}
