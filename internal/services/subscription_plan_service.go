package services

import (
	"fmt"

	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/repositories"
)

// SubscriptionPlanServiceInterface defines the interface for subscription plan operations
type SubscriptionPlanServiceInterface interface {
	CreatePlan(req *models.CreatePlanRequest) (*models.PlanResponse, error)
	GetPlan(id int) (*models.PlanResponse, error)
	GetAllPlans(limit, offset int, isActive *bool) ([]*models.PlanResponse, int, error)
	UpdatePlan(id int, req *models.UpdatePlanRequest) (*models.PlanResponse, error)
	DeletePlan(id int) error
}

// SubscriptionPlanService handles subscription plan business logic
type SubscriptionPlanService struct {
	repo repositories.SubscriptionPlanRepositoryInterface
}

// NewSubscriptionPlanService creates a new subscription plan service
func NewSubscriptionPlanService(repo repositories.SubscriptionPlanRepositoryInterface) *SubscriptionPlanService {
	return &SubscriptionPlanService{repo: repo}
}

// CreatePlan creates a new subscription plan
func (s *SubscriptionPlanService) CreatePlan(req *models.CreatePlanRequest) (*models.PlanResponse, error) {
	// Validate plan type specific requirements
	if err := s.validatePlanRequest(req); err != nil {
		return nil, err
	}

	plan := &models.SubscriptionPlan{
		Name:                     req.Name,
		PlanType:                 req.PlanType,
		Price:                    req.Price,
		AnnualDiscountPercentage: req.AnnualDiscountPercentage,
		TrialDurationDays:        req.TrialDurationDays,
		IsActive:                 req.IsActive,
	}

	if err := s.repo.Create(plan); err != nil {
		return nil, fmt.Errorf("failed to create plan: %w", err)
	}

	response := plan.ToResponse()
	return &response, nil
}

// GetPlan retrieves a subscription plan by ID
func (s *SubscriptionPlanService) GetPlan(id int) (*models.PlanResponse, error) {
	plan, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}

	response := plan.ToResponse()
	return &response, nil
}

// GetAllPlans retrieves all subscription plans with pagination
func (s *SubscriptionPlanService) GetAllPlans(limit, offset int, isActive *bool) ([]*models.PlanResponse, int, error) {
	plans, err := s.repo.GetAll(limit, offset, isActive)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get plans: %w", err)
	}

	total, err := s.repo.Count(isActive)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count plans: %w", err)
	}

	var responses []*models.PlanResponse
	for _, plan := range plans {
		response := plan.ToResponse()
		responses = append(responses, &response)
	}

	return responses, total, nil
}

// UpdatePlan updates an existing subscription plan
func (s *SubscriptionPlanService) UpdatePlan(id int, req *models.UpdatePlanRequest) (*models.PlanResponse, error) {
	// Get existing plan
	existingPlan, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing plan: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		existingPlan.Name = *req.Name
	}
	if req.PlanType != nil {
		existingPlan.PlanType = *req.PlanType
	}
	if req.Price != nil {
		existingPlan.Price = *req.Price
	}
	if req.AnnualDiscountPercentage != nil {
		existingPlan.AnnualDiscountPercentage = req.AnnualDiscountPercentage
	}
	if req.TrialDurationDays != nil {
		existingPlan.TrialDurationDays = req.TrialDurationDays
	}
	if req.IsActive != nil {
		existingPlan.IsActive = *req.IsActive
	}

	// Validate updated plan
	if err := s.validatePlanUpdate(existingPlan); err != nil {
		return nil, err
	}

	if err := s.repo.Update(id, existingPlan); err != nil {
		return nil, fmt.Errorf("failed to update plan: %w", err)
	}

	response := existingPlan.ToResponse()
	return &response, nil
}

// DeletePlan deletes a subscription plan
func (s *SubscriptionPlanService) DeletePlan(id int) error {
	// Check if plan exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get plan: %w", err)
	}

	// Check if plan has active subscriptions
	hasActiveSubscriptions, err := s.repo.HasActiveSubscriptions(id)
	if err != nil {
		return fmt.Errorf("failed to check active subscriptions: %w", err)
	}

	if hasActiveSubscriptions {
		return fmt.Errorf("cannot delete plan: plan has active subscriptions")
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete plan: %w", err)
	}

	return nil
}

// validatePlanRequest validates plan creation request
func (s *SubscriptionPlanService) validatePlanRequest(req *models.CreatePlanRequest) error {
	// Trial plans must have trial duration
	if req.PlanType == "trial" && (req.TrialDurationDays == nil || *req.TrialDurationDays <= 0) {
		return fmt.Errorf("trial plans must have a valid trial duration")
	}

	// Non-trial plans should have a price
	if req.PlanType != "trial" && req.Price <= 0 {
		return fmt.Errorf("non-trial plans must have a positive price")
	}

	// Annual plans can have discount
	if req.PlanType == "annual" && req.AnnualDiscountPercentage != nil {
		if *req.AnnualDiscountPercentage < 0 || *req.AnnualDiscountPercentage > 100 {
			return fmt.Errorf("annual discount percentage must be between 0 and 100")
		}
	}

	// Non-annual plans should not have discount
	if req.PlanType != "annual" && req.AnnualDiscountPercentage != nil && *req.AnnualDiscountPercentage > 0 {
		return fmt.Errorf("only annual plans can have discount percentage")
	}

	return nil
}

// validatePlanUpdate validates plan update request
func (s *SubscriptionPlanService) validatePlanUpdate(plan *models.SubscriptionPlan) error {
	// Trial plans must have trial duration
	if plan.PlanType == "trial" && (plan.TrialDurationDays == nil || *plan.TrialDurationDays <= 0) {
		return fmt.Errorf("trial plans must have a valid trial duration")
	}

	// Non-trial plans should have a price
	if plan.PlanType != "trial" && plan.Price <= 0 {
		return fmt.Errorf("non-trial plans must have a positive price")
	}

	// Annual plans can have discount
	if plan.PlanType == "annual" && plan.AnnualDiscountPercentage != nil {
		if *plan.AnnualDiscountPercentage < 0 || *plan.AnnualDiscountPercentage > 100 {
			return fmt.Errorf("annual discount percentage must be between 0 and 100")
		}
	}

	// Non-annual plans should not have discount
	if plan.PlanType != "annual" && plan.AnnualDiscountPercentage != nil && *plan.AnnualDiscountPercentage > 0 {
		return fmt.Errorf("only annual plans can have discount percentage")
	}

	return nil
}
