package services

import (
	"context"

	"github.com/gatehide/gatehide-api/internal/models"
)

// GamenetServiceInterface defines the contract for gamenet services
type GamenetServiceInterface interface {
	// GetAll retrieves all gamenets
	GetAll(ctx context.Context) ([]models.GamenetResponse, error)

	// GetByID retrieves a gamenet by ID
	GetByID(ctx context.Context, id int) (*models.GamenetResponse, error)

	// Create creates a new gamenet
	Create(ctx context.Context, req *models.GamenetCreateRequest) (*models.GamenetResponse, error)

	// Update updates an existing gamenet
	Update(ctx context.Context, id int, req *models.GamenetUpdateRequest) (*models.GamenetResponse, error)

	// Delete deletes a gamenet
	Delete(ctx context.Context, id int) error

	// Search searches gamenets with pagination
	Search(ctx context.Context, req *models.GamenetSearchRequest) (*models.GamenetSearchResponse, error)

	// ResendCredentials generates new password and sends credentials via email
	ResendCredentials(ctx context.Context, id int) error
}
