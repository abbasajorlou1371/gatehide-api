package services

import (
	"context"
	"fmt"

	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/repositories"
	"github.com/gatehide/gatehide-api/internal/utils"
)

// gamenetService implements GamenetServiceInterface
type gamenetService struct {
	gamenetRepo  repositories.GamenetRepository
	smsService   *SMSService
	emailService *EmailService
}

// NewGamenetService creates a new gamenet service
func NewGamenetService(gamenetRepo repositories.GamenetRepository, smsService *SMSService, emailService *EmailService) GamenetServiceInterface {
	return &gamenetService{
		gamenetRepo:  gamenetRepo,
		smsService:   smsService,
		emailService: emailService,
	}
}

// GetAll retrieves all gamenets
func (s *gamenetService) GetAll(ctx context.Context) ([]models.GamenetResponse, error) {
	gamenets, err := s.gamenetRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get gamenets: %w", err)
	}

	var responses []models.GamenetResponse
	for _, gamenet := range gamenets {
		responses = append(responses, gamenet.ToResponse())
	}

	return responses, nil
}

// GetByID retrieves a gamenet by ID
func (s *gamenetService) GetByID(ctx context.Context, id int) (*models.GamenetResponse, error) {
	gamenet, err := s.gamenetRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get gamenet: %w", err)
	}

	response := gamenet.ToResponse()
	return &response, nil
}

// Create creates a new gamenet
func (s *gamenetService) Create(ctx context.Context, req *models.GamenetCreateRequest) (*models.GamenetResponse, error) {
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

	gamenet := &models.Gamenet{
		Name:              req.Name,
		OwnerName:         req.OwnerName,
		OwnerMobile:       req.OwnerMobile,
		Address:           req.Address,
		Email:             req.Email,
		Password:          hashedPassword,
		LicenseAttachment: req.LicenseAttachment,
	}

	err = s.gamenetRepo.Create(gamenet)
	if err != nil {
		return nil, fmt.Errorf("failed to create gamenet: %w", err)
	}

	// Send credentials via SMS using Kavenegar Verify Lookup
	if s.smsService != nil {
		err = s.smsService.SendGamenetCredentials(ctx, req.OwnerMobile, req.Email, randomPassword)
		if err != nil {
			// Log the error but don't fail the creation
			fmt.Printf("Warning: Failed to send credentials SMS to %s: %v\n", req.OwnerMobile, err)
		} else {
			fmt.Printf("Successfully sent credentials SMS to %s\n", req.OwnerMobile)
		}
	}

	response := gamenet.ToResponse()
	return &response, nil
}

// Update updates an existing gamenet
func (s *gamenetService) Update(ctx context.Context, id int, req *models.GamenetUpdateRequest) (*models.GamenetResponse, error) {
	// Check if gamenet exists
	_, err := s.gamenetRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("gamenet not found: %w", err)
	}

	err = s.gamenetRepo.Update(id, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update gamenet: %w", err)
	}

	// Get updated gamenet
	updatedGamenet, err := s.gamenetRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated gamenet: %w", err)
	}

	response := updatedGamenet.ToResponse()
	return &response, nil
}

// Delete deletes a gamenet
func (s *gamenetService) Delete(ctx context.Context, id int) error {
	// Check if gamenet exists
	_, err := s.gamenetRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("gamenet not found: %w", err)
	}

	err = s.gamenetRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete gamenet: %w", err)
	}

	return nil
}

// Search searches gamenets with pagination
func (s *gamenetService) Search(ctx context.Context, req *models.GamenetSearchRequest) (*models.GamenetSearchResponse, error) {
	result, err := s.gamenetRepo.Search(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search gamenets: %w", err)
	}

	return result, nil
}

// ResendCredentials generates new password and sends credentials via email
func (s *gamenetService) ResendCredentials(ctx context.Context, id int) error {
	// Get gamenet details
	gamenet, err := s.gamenetRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("gamenet not found: %w", err)
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
	updateReq := &models.GamenetUpdateRequest{
		Password: &hashedPassword,
	}
	err = s.gamenetRepo.Update(id, updateReq)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Send credentials via SMS
	if s.smsService != nil {
		err = s.smsService.SendGamenetCredentials(ctx, gamenet.OwnerMobile, gamenet.Email, newPassword)
		if err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Warning: Failed to send credentials SMS to %s: %v\n", gamenet.OwnerMobile, err)
			return fmt.Errorf("password updated but failed to send SMS: %w", err)
		}

		fmt.Printf("Successfully sent credentials SMS to %s\n", gamenet.OwnerMobile)
	} else {
		return fmt.Errorf("SMS service not configured")
	}

	return nil
}
