package unit

import (
	"errors"
	"testing"

	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/services"
	"github.com/gatehide/gatehide-api/tests/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubscriptionPlanService_CreatePlan(t *testing.T) {
	tests := []struct {
		name           string
		request        *models.CreatePlanRequest
		mockSetup      func(*utils.MockSubscriptionPlanRepository)
		expectedError  string
		expectedResult *models.PlanResponse
	}{
		{
			name: "successful monthly plan creation",
			request: &models.CreatePlanRequest{
				Name:     "Basic Monthly",
				PlanType: "monthly",
				Price:    29.99,
				IsActive: true,
			},
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				mockRepo.On("Create", mock.AnythingOfType("*models.SubscriptionPlan")).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "successful annual plan creation with discount",
			request: &models.CreatePlanRequest{
				Name:                     "Premium Annual",
				PlanType:                 "annual",
				Price:                    299.99,
				AnnualDiscountPercentage: func() *float64 { v := 20.0; return &v }(),
				IsActive:                 true,
			},
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				mockRepo.On("Create", mock.AnythingOfType("*models.SubscriptionPlan")).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "successful trial plan creation",
			request: &models.CreatePlanRequest{
				Name:              "Free Trial",
				PlanType:          "trial",
				Price:             0.0,
				TrialDurationDays: func() *int { v := 14; return &v }(),
				IsActive:          true,
			},
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				mockRepo.On("Create", mock.AnythingOfType("*models.SubscriptionPlan")).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "trial plan without duration fails",
			request: &models.CreatePlanRequest{
				Name:     "Invalid Trial",
				PlanType: "trial",
				Price:    0.0,
				IsActive: true,
			},
			mockSetup:     func(mockRepo *utils.MockSubscriptionPlanRepository) {},
			expectedError: "trial plans must have a valid trial duration",
		},
		{
			name: "trial plan with zero duration fails",
			request: &models.CreatePlanRequest{
				Name:              "Invalid Trial",
				PlanType:          "trial",
				Price:             0.0,
				TrialDurationDays: func() *int { v := 0; return &v }(),
				IsActive:          true,
			},
			mockSetup:     func(mockRepo *utils.MockSubscriptionPlanRepository) {},
			expectedError: "trial plans must have a valid trial duration",
		},
		{
			name: "monthly plan with zero price fails",
			request: &models.CreatePlanRequest{
				Name:     "Invalid Monthly",
				PlanType: "monthly",
				Price:    0.0,
				IsActive: true,
			},
			mockSetup:     func(mockRepo *utils.MockSubscriptionPlanRepository) {},
			expectedError: "non-trial plans must have a positive price",
		},
		{
			name: "annual plan with invalid discount fails",
			request: &models.CreatePlanRequest{
				Name:                     "Invalid Annual",
				PlanType:                 "annual",
				Price:                    299.99,
				AnnualDiscountPercentage: func() *float64 { v := 150.0; return &v }(),
				IsActive:                 true,
			},
			mockSetup:     func(mockRepo *utils.MockSubscriptionPlanRepository) {},
			expectedError: "annual discount percentage must be between 0 and 100",
		},
		{
			name: "monthly plan with discount fails",
			request: &models.CreatePlanRequest{
				Name:                     "Invalid Monthly",
				PlanType:                 "monthly",
				Price:                    29.99,
				AnnualDiscountPercentage: func() *float64 { v := 10.0; return &v }(),
				IsActive:                 true,
			},
			mockSetup:     func(mockRepo *utils.MockSubscriptionPlanRepository) {},
			expectedError: "only annual plans can have discount percentage",
		},
		{
			name: "repository error during creation",
			request: &models.CreatePlanRequest{
				Name:     "Basic Monthly",
				PlanType: "monthly",
				Price:    29.99,
				IsActive: true,
			},
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				mockRepo.On("Create", mock.AnythingOfType("*models.SubscriptionPlan")).Return(errors.New("database error"))
			},
			expectedError: "failed to create plan: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := utils.SetupMockSubscriptionPlanRepository(t)
			tt.mockSetup(mockRepo)
			service := services.NewSubscriptionPlanService(mockRepo)

			// Execute
			result, err := service.CreatePlan(tt.request)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.Name, result.Name)
				assert.Equal(t, tt.request.PlanType, result.PlanType)
				assert.Equal(t, tt.request.Price, result.Price)
				assert.Equal(t, tt.request.IsActive, result.IsActive)
			}

			utils.AssertSubscriptionPlanRepositoryExpectations(t, mockRepo)
		})
	}
}

func TestSubscriptionPlanService_GetPlan(t *testing.T) {
	tests := []struct {
		name           string
		planID         int
		mockSetup      func(*utils.MockSubscriptionPlanRepository)
		expectedError  string
		expectedResult *models.PlanResponse
	}{
		{
			name:   "successful plan retrieval",
			planID: 1,
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				mockPlan := utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99)
				mockRepo.On("GetByID", 1).Return(mockPlan, nil)
			},
			expectedError: "",
		},
		{
			name:   "plan not found",
			planID: 999,
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				mockRepo.On("GetByID", 999).Return(nil, errors.New("subscription plan not found"))
			},
			expectedError: "failed to get plan: subscription plan not found",
		},
		{
			name:   "database error",
			planID: 1,
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				mockRepo.On("GetByID", 1).Return(nil, errors.New("database connection failed"))
			},
			expectedError: "failed to get plan: database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := utils.SetupMockSubscriptionPlanRepository(t)
			tt.mockSetup(mockRepo)
			service := services.NewSubscriptionPlanService(mockRepo)

			// Execute
			result, err := service.GetPlan(tt.planID)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.planID, result.ID)
			}

			utils.AssertSubscriptionPlanRepositoryExpectations(t, mockRepo)
		})
	}
}

func TestSubscriptionPlanService_GetAllPlans(t *testing.T) {
	tests := []struct {
		name          string
		limit         int
		offset        int
		isActive      *bool
		mockSetup     func(*utils.MockSubscriptionPlanRepository)
		expectedError string
		expectedCount int
		expectedTotal int
	}{
		{
			name:     "successful plans retrieval with pagination",
			limit:    10,
			offset:   0,
			isActive: nil,
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				plans := []*models.SubscriptionPlan{
					utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99),
					utils.CreateMockSubscriptionPlan(2, "Premium Annual", "annual", 299.99),
				}
				mockRepo.On("GetAll", 10, 0, (*bool)(nil)).Return(plans, nil)
				mockRepo.On("Count", (*bool)(nil)).Return(2, nil)
			},
			expectedError: "",
			expectedCount: 2,
			expectedTotal: 2,
		},
		{
			name:     "successful active plans retrieval",
			limit:    5,
			offset:   0,
			isActive: func() *bool { v := true; return &v }(),
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				plans := []*models.SubscriptionPlan{
					utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99),
				}
				active := true
				mockRepo.On("GetAll", 5, 0, &active).Return(plans, nil)
				mockRepo.On("Count", &active).Return(1, nil)
			},
			expectedError: "",
			expectedCount: 1,
			expectedTotal: 1,
		},
		{
			name:     "database error during retrieval",
			limit:    10,
			offset:   0,
			isActive: nil,
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				mockRepo.On("GetAll", 10, 0, (*bool)(nil)).Return(nil, errors.New("database error"))
			},
			expectedError: "failed to get plans: database error",
		},
		{
			name:     "database error during count",
			limit:    10,
			offset:   0,
			isActive: nil,
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				plans := []*models.SubscriptionPlan{
					utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99),
				}
				mockRepo.On("GetAll", 10, 0, (*bool)(nil)).Return(plans, nil)
				mockRepo.On("Count", (*bool)(nil)).Return(0, errors.New("count error"))
			},
			expectedError: "failed to count plans: count error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := utils.SetupMockSubscriptionPlanRepository(t)
			tt.mockSetup(mockRepo)
			service := services.NewSubscriptionPlanService(mockRepo)

			// Execute
			plans, total, err := service.GetAllPlans(tt.limit, tt.offset, tt.isActive)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, plans)
				assert.Equal(t, 0, total)
			} else {
				assert.NoError(t, err)
				assert.Len(t, plans, tt.expectedCount)
				assert.Equal(t, tt.expectedTotal, total)
			}

			utils.AssertSubscriptionPlanRepositoryExpectations(t, mockRepo)
		})
	}
}

func TestSubscriptionPlanService_UpdatePlan(t *testing.T) {
	tests := []struct {
		name           string
		planID         int
		request        *models.UpdatePlanRequest
		mockSetup      func(*utils.MockSubscriptionPlanRepository)
		expectedError  string
		expectedResult *models.PlanResponse
	}{
		{
			name:   "successful plan update",
			planID: 1,
			request: &models.UpdatePlanRequest{
				Name:  func() *string { v := "Updated Basic Monthly"; return &v }(),
				Price: func() *float64 { v := 39.99; return &v }(),
			},
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				existingPlan := utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99)
				mockRepo.On("GetByID", 1).Return(existingPlan, nil)
				mockRepo.On("Update", 1, mock.AnythingOfType("*models.SubscriptionPlan")).Return(nil)
			},
			expectedError: "",
		},
		{
			name:   "plan not found",
			planID: 999,
			request: &models.UpdatePlanRequest{
				Name: func() *string { v := "Updated Plan"; return &v }(),
			},
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				mockRepo.On("GetByID", 999).Return(nil, errors.New("subscription plan not found"))
			},
			expectedError: "failed to get existing plan: subscription plan not found",
		},
		{
			name:   "invalid update - trial plan without duration",
			planID: 1,
			request: &models.UpdatePlanRequest{
				PlanType: func() *string { v := "trial"; return &v }(),
			},
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				existingPlan := utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99)
				mockRepo.On("GetByID", 1).Return(existingPlan, nil)
			},
			expectedError: "trial plans must have a valid trial duration",
		},
		{
			name:   "invalid update - monthly plan with zero price",
			planID: 1,
			request: &models.UpdatePlanRequest{
				Price: func() *float64 { v := 0.0; return &v }(),
			},
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				existingPlan := utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99)
				mockRepo.On("GetByID", 1).Return(existingPlan, nil)
			},
			expectedError: "non-trial plans must have a positive price",
		},
		{
			name:   "database error during update",
			planID: 1,
			request: &models.UpdatePlanRequest{
				Name: func() *string { v := "Updated Plan"; return &v }(),
			},
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				existingPlan := utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99)
				mockRepo.On("GetByID", 1).Return(existingPlan, nil)
				mockRepo.On("Update", 1, mock.AnythingOfType("*models.SubscriptionPlan")).Return(errors.New("database error"))
			},
			expectedError: "failed to update plan: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := utils.SetupMockSubscriptionPlanRepository(t)
			tt.mockSetup(mockRepo)
			service := services.NewSubscriptionPlanService(mockRepo)

			// Execute
			result, err := service.UpdatePlan(tt.planID, tt.request)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.planID, result.ID)
			}

			utils.AssertSubscriptionPlanRepositoryExpectations(t, mockRepo)
		})
	}
}

func TestSubscriptionPlanService_DeletePlan(t *testing.T) {
	tests := []struct {
		name          string
		planID        int
		mockSetup     func(*utils.MockSubscriptionPlanRepository)
		expectedError string
	}{
		{
			name:   "successful plan deletion",
			planID: 1,
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				existingPlan := utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99)
				mockRepo.On("GetByID", 1).Return(existingPlan, nil)
				mockRepo.On("HasActiveSubscriptions", 1).Return(false, nil)
				mockRepo.On("Delete", 1).Return(nil)
			},
			expectedError: "",
		},
		{
			name:   "plan not found",
			planID: 999,
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				mockRepo.On("GetByID", 999).Return(nil, errors.New("subscription plan not found"))
			},
			expectedError: "failed to get plan: subscription plan not found",
		},
		{
			name:   "cannot delete plan with active subscriptions",
			planID: 1,
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				existingPlan := utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99)
				mockRepo.On("GetByID", 1).Return(existingPlan, nil)
				mockRepo.On("HasActiveSubscriptions", 1).Return(true, nil)
			},
			expectedError: "cannot delete plan: plan has active subscriptions",
		},
		{
			name:   "error checking active subscriptions",
			planID: 1,
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				existingPlan := utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99)
				mockRepo.On("GetByID", 1).Return(existingPlan, nil)
				mockRepo.On("HasActiveSubscriptions", 1).Return(false, errors.New("database error"))
			},
			expectedError: "failed to check active subscriptions: database error",
		},
		{
			name:   "database error during deletion",
			planID: 1,
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				existingPlan := utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99)
				mockRepo.On("GetByID", 1).Return(existingPlan, nil)
				mockRepo.On("HasActiveSubscriptions", 1).Return(false, nil)
				mockRepo.On("Delete", 1).Return(errors.New("database error"))
			},
			expectedError: "failed to delete plan: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := utils.SetupMockSubscriptionPlanRepository(t)
			tt.mockSetup(mockRepo)
			service := services.NewSubscriptionPlanService(mockRepo)

			// Execute
			err := service.DeletePlan(tt.planID)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			utils.AssertSubscriptionPlanRepositoryExpectations(t, mockRepo)
		})
	}
}
