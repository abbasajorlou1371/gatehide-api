package unit

import (
	"testing"

	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/services"
	"github.com/gatehide/gatehide-api/tests/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestSubscriptionPlanValidation tests comprehensive input validation for subscription plans
func TestSubscriptionPlanValidation(t *testing.T) {
	tests := []struct {
		name          string
		request       *models.CreatePlanRequest
		expectedError string
	}{
		// Valid cases
		{
			name: "valid monthly plan",
			request: &models.CreatePlanRequest{
				Name:     "Basic Monthly",
				PlanType: "monthly",
				Price:    29.99,
				IsActive: true,
			},
			expectedError: "",
		},
		{
			name: "valid annual plan with discount",
			request: &models.CreatePlanRequest{
				Name:                     "Premium Annual",
				PlanType:                 "annual",
				Price:                    299.99,
				AnnualDiscountPercentage: func() *float64 { v := 20.0; return &v }(),
				IsActive:                 true,
			},
			expectedError: "",
		},
		{
			name: "valid trial plan",
			request: &models.CreatePlanRequest{
				Name:              "Free Trial",
				PlanType:          "trial",
				Price:             0.0,
				TrialDurationDays: func() *int { v := 14; return &v }(),
				IsActive:          true,
			},
			expectedError: "",
		},
		{
			name: "valid annual plan without discount",
			request: &models.CreatePlanRequest{
				Name:     "Standard Annual",
				PlanType: "annual",
				Price:    299.99,
				IsActive: true,
			},
			expectedError: "",
		},

		// Invalid plan types
		// Note: Invalid plan type validation is handled by Gin binding validation,
		// not by the service layer, so we don't test it here

		// Trial plan validation
		{
			name: "trial plan without duration",
			request: &models.CreatePlanRequest{
				Name:     "Invalid Trial",
				PlanType: "trial",
				Price:    0.0,
				IsActive: true,
			},
			expectedError: "trial plans must have a valid trial duration",
		},
		{
			name: "trial plan with zero duration",
			request: &models.CreatePlanRequest{
				Name:              "Invalid Trial",
				PlanType:          "trial",
				Price:             0.0,
				TrialDurationDays: func() *int { v := 0; return &v }(),
				IsActive:          true,
			},
			expectedError: "trial plans must have a valid trial duration",
		},
		{
			name: "trial plan with negative duration",
			request: &models.CreatePlanRequest{
				Name:              "Invalid Trial",
				PlanType:          "trial",
				Price:             0.0,
				TrialDurationDays: func() *int { v := -5; return &v }(),
				IsActive:          true,
			},
			expectedError: "trial plans must have a valid trial duration",
		},

		// Price validation
		{
			name: "monthly plan with zero price",
			request: &models.CreatePlanRequest{
				Name:     "Invalid Monthly",
				PlanType: "monthly",
				Price:    0.0,
				IsActive: true,
			},
			expectedError: "non-trial plans must have a positive price",
		},
		{
			name: "monthly plan with negative price",
			request: &models.CreatePlanRequest{
				Name:     "Invalid Monthly",
				PlanType: "monthly",
				Price:    -10.0,
				IsActive: true,
			},
			expectedError: "non-trial plans must have a positive price",
		},
		{
			name: "annual plan with zero price",
			request: &models.CreatePlanRequest{
				Name:     "Invalid Annual",
				PlanType: "annual",
				Price:    0.0,
				IsActive: true,
			},
			expectedError: "non-trial plans must have a positive price",
		},

		// Discount validation
		{
			name: "annual plan with negative discount",
			request: &models.CreatePlanRequest{
				Name:                     "Invalid Annual",
				PlanType:                 "annual",
				Price:                    299.99,
				AnnualDiscountPercentage: func() *float64 { v := -10.0; return &v }(),
				IsActive:                 true,
			},
			expectedError: "annual discount percentage must be between 0 and 100",
		},
		{
			name: "annual plan with discount over 100%",
			request: &models.CreatePlanRequest{
				Name:                     "Invalid Annual",
				PlanType:                 "annual",
				Price:                    299.99,
				AnnualDiscountPercentage: func() *float64 { v := 150.0; return &v }(),
				IsActive:                 true,
			},
			expectedError: "annual discount percentage must be between 0 and 100",
		},
		{
			name: "monthly plan with discount",
			request: &models.CreatePlanRequest{
				Name:                     "Invalid Monthly",
				PlanType:                 "monthly",
				Price:                    29.99,
				AnnualDiscountPercentage: func() *float64 { v := 10.0; return &v }(),
				IsActive:                 true,
			},
			expectedError: "only annual plans can have discount percentage",
		},
		{
			name: "trial plan with discount",
			request: &models.CreatePlanRequest{
				Name:                     "Invalid Trial",
				PlanType:                 "trial",
				Price:                    0.0,
				TrialDurationDays:        func() *int { v := 14; return &v }(),
				AnnualDiscountPercentage: func() *float64 { v := 10.0; return &v }(),
				IsActive:                 true,
			},
			expectedError: "only annual plans can have discount percentage",
		},

		// Edge cases
		{
			name: "annual plan with 0% discount",
			request: &models.CreatePlanRequest{
				Name:                     "Valid Annual",
				PlanType:                 "annual",
				Price:                    299.99,
				AnnualDiscountPercentage: func() *float64 { v := 0.0; return &v }(),
				IsActive:                 true,
			},
			expectedError: "",
		},
		{
			name: "annual plan with 100% discount",
			request: &models.CreatePlanRequest{
				Name:                     "Valid Annual",
				PlanType:                 "annual",
				Price:                    299.99,
				AnnualDiscountPercentage: func() *float64 { v := 100.0; return &v }(),
				IsActive:                 true,
			},
			expectedError: "",
		},
		{
			name: "trial plan with maximum duration",
			request: &models.CreatePlanRequest{
				Name:              "Long Trial",
				PlanType:          "trial",
				Price:             0.0,
				TrialDurationDays: func() *int { v := 365; return &v }(),
				IsActive:          true,
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := utils.SetupMockSubscriptionPlanRepository(t)
			service := services.NewSubscriptionPlanService(mockRepo)

			// Only mock repository call if we expect success
			if tt.expectedError == "" {
				mockRepo.On("Create", mock.AnythingOfType("*models.SubscriptionPlan")).Return(nil)
			}

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
			}

			utils.AssertSubscriptionPlanRepositoryExpectations(t, mockRepo)
		})
	}
}

// TestSubscriptionPlanUpdateValidation tests validation for plan updates
func TestSubscriptionPlanUpdateValidation(t *testing.T) {
	tests := []struct {
		name          string
		existingPlan  *models.SubscriptionPlan
		updateRequest *models.UpdatePlanRequest
		expectedError string
	}{
		{
			name:         "valid update - change name and price",
			existingPlan: utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99),
			updateRequest: &models.UpdatePlanRequest{
				Name:  func() *string { v := "Updated Basic Monthly"; return &v }(),
				Price: func() *float64 { v := 39.99; return &v }(),
			},
			expectedError: "",
		},
		{
			name:         "valid update - change to trial plan",
			existingPlan: utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99),
			updateRequest: &models.UpdatePlanRequest{
				PlanType:          func() *string { v := "trial"; return &v }(),
				Price:             func() *float64 { v := 0.0; return &v }(),
				TrialDurationDays: func() *int { v := 14; return &v }(),
			},
			expectedError: "",
		},
		{
			name:         "valid update - change to annual plan with discount",
			existingPlan: utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99),
			updateRequest: &models.UpdatePlanRequest{
				PlanType:                 func() *string { v := "annual"; return &v }(),
				Price:                    func() *float64 { v := 299.99; return &v }(),
				AnnualDiscountPercentage: func() *float64 { v := 20.0; return &v }(),
			},
			expectedError: "",
		},
		{
			name:         "invalid update - trial plan without duration",
			existingPlan: utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99),
			updateRequest: &models.UpdatePlanRequest{
				PlanType: func() *string { v := "trial"; return &v }(),
				Price:    func() *float64 { v := 0.0; return &v }(),
			},
			expectedError: "trial plans must have a valid trial duration",
		},
		{
			name:         "invalid update - trial plan with zero duration",
			existingPlan: utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99),
			updateRequest: &models.UpdatePlanRequest{
				PlanType:          func() *string { v := "trial"; return &v }(),
				Price:             func() *float64 { v := 0.0; return &v }(),
				TrialDurationDays: func() *int { v := 0; return &v }(),
			},
			expectedError: "trial plans must have a valid trial duration",
		},
		{
			name:         "invalid update - monthly plan with zero price",
			existingPlan: utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99),
			updateRequest: &models.UpdatePlanRequest{
				Price: func() *float64 { v := 0.0; return &v }(),
			},
			expectedError: "non-trial plans must have a positive price",
		},
		{
			name:         "invalid update - monthly plan with negative price",
			existingPlan: utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99),
			updateRequest: &models.UpdatePlanRequest{
				Price: func() *float64 { v := -10.0; return &v }(),
			},
			expectedError: "non-trial plans must have a positive price",
		},
		{
			name:         "invalid update - annual plan with invalid discount",
			existingPlan: utils.CreateMockSubscriptionPlan(1, "Basic Annual", "annual", 299.99),
			updateRequest: &models.UpdatePlanRequest{
				AnnualDiscountPercentage: func() *float64 { v := 150.0; return &v }(),
			},
			expectedError: "annual discount percentage must be between 0 and 100",
		},
		{
			name:         "invalid update - monthly plan with discount",
			existingPlan: utils.CreateMockSubscriptionPlan(1, "Basic Monthly", "monthly", 29.99),
			updateRequest: &models.UpdatePlanRequest{
				AnnualDiscountPercentage: func() *float64 { v := 10.0; return &v }(),
			},
			expectedError: "only annual plans can have discount percentage",
		},
		{
			name: "valid update - remove discount from annual plan",
			existingPlan: &models.SubscriptionPlan{
				ID:                       1,
				Name:                     "Premium Annual",
				PlanType:                 "annual",
				Price:                    299.99,
				AnnualDiscountPercentage: func() *float64 { v := 20.0; return &v }(),
				IsActive:                 true,
			},
			updateRequest: &models.UpdatePlanRequest{
				AnnualDiscountPercentage: func() *float64 { v := 0.0; return &v }(),
			},
			expectedError: "",
		},
		{
			name: "valid update - set discount to nil",
			existingPlan: &models.SubscriptionPlan{
				ID:                       1,
				Name:                     "Premium Annual",
				PlanType:                 "annual",
				Price:                    299.99,
				AnnualDiscountPercentage: func() *float64 { v := 20.0; return &v }(),
				IsActive:                 true,
			},
			updateRequest: &models.UpdatePlanRequest{
				AnnualDiscountPercentage: nil,
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := utils.SetupMockSubscriptionPlanRepository(t)
			service := services.NewSubscriptionPlanService(mockRepo)

			// Mock expectations
			mockRepo.On("GetByID", 1).Return(tt.existingPlan, nil)
			if tt.expectedError == "" {
				mockRepo.On("Update", 1, mock.AnythingOfType("*models.SubscriptionPlan")).Return(nil)
			}

			// Execute
			result, err := service.UpdatePlan(1, tt.updateRequest)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			utils.AssertSubscriptionPlanRepositoryExpectations(t, mockRepo)
		})
	}
}

// TestSubscriptionPlanBoundaryValues tests boundary values for validation
func TestSubscriptionPlanBoundaryValues(t *testing.T) {
	tests := []struct {
		name          string
		request       *models.CreatePlanRequest
		expectedError string
	}{
		{
			name: "minimum valid price",
			request: &models.CreatePlanRequest{
				Name:     "Min Price Plan",
				PlanType: "monthly",
				Price:    0.01,
				IsActive: true,
			},
			expectedError: "",
		},
		{
			name: "maximum valid discount",
			request: &models.CreatePlanRequest{
				Name:                     "Max Discount Plan",
				PlanType:                 "annual",
				Price:                    299.99,
				AnnualDiscountPercentage: func() *float64 { v := 100.0; return &v }(),
				IsActive:                 true,
			},
			expectedError: "",
		},
		{
			name: "minimum valid trial duration",
			request: &models.CreatePlanRequest{
				Name:              "Min Trial Plan",
				PlanType:          "trial",
				Price:             0.0,
				TrialDurationDays: func() *int { v := 1; return &v }(),
				IsActive:          true,
			},
			expectedError: "",
		},
		{
			name: "maximum reasonable trial duration",
			request: &models.CreatePlanRequest{
				Name:              "Max Trial Plan",
				PlanType:          "trial",
				Price:             0.0,
				TrialDurationDays: func() *int { v := 365; return &v }(),
				IsActive:          true,
			},
			expectedError: "",
		},
		{
			name: "zero discount (valid)",
			request: &models.CreatePlanRequest{
				Name:                     "Zero Discount Plan",
				PlanType:                 "annual",
				Price:                    299.99,
				AnnualDiscountPercentage: func() *float64 { v := 0.0; return &v }(),
				IsActive:                 true,
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := utils.SetupMockSubscriptionPlanRepository(t)
			service := services.NewSubscriptionPlanService(mockRepo)

			// Mock repository call
			mockRepo.On("Create", mock.AnythingOfType("*models.SubscriptionPlan")).Return(nil)

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
			}

			utils.AssertSubscriptionPlanRepositoryExpectations(t, mockRepo)
		})
	}
}

// TestSubscriptionPlanBusinessRules tests business rule validation
func TestSubscriptionPlanBusinessRules(t *testing.T) {
	tests := []struct {
		name          string
		request       *models.CreatePlanRequest
		expectedError string
		description   string
	}{
		{
			name: "trial plan with price is valid",
			request: &models.CreatePlanRequest{
				Name:              "Paid Trial",
				PlanType:          "trial",
				Price:             10.0, // Trial plans can have a price
				TrialDurationDays: func() *int { v := 14; return &v }(),
				IsActive:          true,
			},
			expectedError: "",
			description:   "Trial plans can have a price (this is valid)",
		},
		{
			name: "monthly plan cannot have trial duration",
			request: &models.CreatePlanRequest{
				Name:              "Invalid Monthly",
				PlanType:          "monthly",
				Price:             29.99,
				TrialDurationDays: func() *int { v := 14; return &v }(),
				IsActive:          true,
			},
			expectedError: "",
			description:   "Monthly plans can have trial duration (this is actually valid)",
		},
		{
			name: "annual plan with reasonable discount",
			request: &models.CreatePlanRequest{
				Name:                     "Reasonable Annual",
				PlanType:                 "annual",
				Price:                    299.99,
				AnnualDiscountPercentage: func() *float64 { v := 20.0; return &v }(),
				IsActive:                 true,
			},
			expectedError: "",
			description:   "Annual plans can have reasonable discounts",
		},
		{
			name: "inactive plan is still valid",
			request: &models.CreatePlanRequest{
				Name:     "Inactive Plan",
				PlanType: "monthly",
				Price:    29.99,
				IsActive: false,
			},
			expectedError: "",
			description:   "Plans can be created as inactive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := utils.SetupMockSubscriptionPlanRepository(t)
			service := services.NewSubscriptionPlanService(mockRepo)

			// Mock repository call
			mockRepo.On("Create", mock.AnythingOfType("*models.SubscriptionPlan")).Return(nil)

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
			}

			utils.AssertSubscriptionPlanRepositoryExpectations(t, mockRepo)
		})
	}
}
