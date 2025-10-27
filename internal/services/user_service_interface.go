package services

import (
	"context"

	"github.com/gatehide/gatehide-api/internal/models"
)

// UserServiceInterface defines the interface for user business logic
type UserServiceInterface interface {
	GetAll(ctx context.Context) ([]models.UserResponse, error)
	GetAllByGamenet(ctx context.Context, gamenetID int) ([]models.UserResponse, error)
	GetByID(ctx context.Context, id int) (*models.UserResponse, error)
	GetByEmail(ctx context.Context, email string) (*models.UserResponse, error)
	GetByMobile(ctx context.Context, mobile string) (*models.UserResponse, error)
	Create(ctx context.Context, req *models.UserCreateRequest, gamenetID *int) (*models.UserResponse, error)
	Update(ctx context.Context, id int, req *models.UserUpdateRequest) (*models.UserResponse, error)
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, req *models.UserSearchRequest) (*models.UserSearchResponse, error)
	SearchByGamenet(ctx context.Context, req *models.UserSearchRequest, gamenetID int) (*models.UserSearchResponse, error)
	AttachToGamenet(ctx context.Context, userID, gamenetID int) error
	DetachFromGamenet(ctx context.Context, userID, gamenetID int) error
	CanModifyUser(ctx context.Context, userID, requesterID int, requesterType string) (bool, error)
	ResendCredentials(ctx context.Context, id int) error
}
