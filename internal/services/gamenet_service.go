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
	gamenetRepo repositories.GamenetRepository
}

// NewGamenetService creates a new gamenet service
func NewGamenetService(gamenetRepo repositories.GamenetRepository) GamenetServiceInterface {
	return &gamenetService{
		gamenetRepo: gamenetRepo,
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
