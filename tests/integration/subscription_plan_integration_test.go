package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gatehide/gatehide-api/internal/handlers"
	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/services"
	"github.com/gatehide/gatehide-api/tests/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestSubscriptionPlanIntegration tests the full integration of subscription plan CRUD operations
func TestSubscriptionPlanIntegration(t *testing.T) {
	// Setup test database (you'll need to implement this based on your test setup)
	// For now, we'll use mocks but in a real integration test, you'd use a test database
	gin.SetMode(gin.TestMode)

	t.Run("Complete CRUD Flow", func(t *testing.T) {
		// This test demonstrates the complete flow of creating, reading, updating, and deleting a subscription plan
		// In a real integration test, you would:
		// 1. Set up a test database
		// 2. Create a real repository instance
		// 3. Create a real service instance
		// 4. Test the actual HTTP endpoints

		// For this example, we'll use mocks but structure it like a real integration test
		mockRepo := utils.SetupMockSubscriptionPlanRepository(t)
		service := services.NewSubscriptionPlanService(mockRepo)
		handler := handlers.NewSubscriptionPlanHandler(service)

		// Test data
		createRequest := models.CreatePlanRequest{
			Name:     "Integration Test Plan",
			PlanType: "monthly",
			Price:    29.99,
			IsActive: true,
		}

		updateRequest := models.UpdatePlanRequest{
			Name:  func() *string { v := "Updated Integration Test Plan"; return &v }(),
			Price: func() *float64 { v := 39.99; return &v }(),
		}

		// Mock expectations for the complete flow
		mockRepo.On("Create", mock.AnythingOfType("*models.SubscriptionPlan")).Return(nil).Once()
		mockRepo.On("GetByID", 1).Return(utils.CreateMockSubscriptionPlan(1, "Integration Test Plan", "monthly", 29.99), nil).Times(3) // Get, Update, Delete
		mockRepo.On("Update", 1, mock.AnythingOfType("*models.SubscriptionPlan")).Return(nil).Once()
		mockRepo.On("HasActiveSubscriptions", 1).Return(false, nil).Once()
		mockRepo.On("Delete", 1).Return(nil).Once()

		// Test CREATE
		t.Run("Create Plan", func(t *testing.T) {
			requestBody, _ := json.Marshal(createRequest)
			req := httptest.NewRequest(http.MethodPost, "/subscription-plans", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			handler.CreatePlan(c)

			assert.Equal(t, http.StatusCreated, w.Code)
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, "Plan created successfully", response["message"])
		})

		// Test READ
		t.Run("Get Plan", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/subscription-plans/1", nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "1"}}

			handler.GetPlan(c)

			assert.Equal(t, http.StatusOK, w.Code)
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.NotNil(t, response["data"])
		})

		// Test UPDATE
		t.Run("Update Plan", func(t *testing.T) {
			requestBody, _ := json.Marshal(updateRequest)
			req := httptest.NewRequest(http.MethodPut, "/subscription-plans/1", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "1"}}

			handler.UpdatePlan(c)

			assert.Equal(t, http.StatusOK, w.Code)
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, "Plan updated successfully", response["message"])
		})

		// Test DELETE
		t.Run("Delete Plan", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/subscription-plans/1", nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "1"}}

			handler.DeletePlan(c)

			assert.Equal(t, http.StatusOK, w.Code)
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, "Plan deleted successfully", response["message"])
		})

		utils.AssertSubscriptionPlanRepositoryExpectations(t, mockRepo)
	})
}

// TestSubscriptionPlanValidationIntegration tests input validation in an integrated manner
func TestSubscriptionPlanValidationIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		endpoint       string
		method         string
		requestBody    interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:     "Create plan with invalid plan type",
			endpoint: "/subscription-plans",
			method:   http.MethodPost,
			requestBody: models.CreatePlanRequest{
				Name:     "Invalid Plan",
				PlanType: "invalid_type",
				Price:    29.99,
				IsActive: true,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request data",
		},
		{
			name:     "Create plan with negative price",
			endpoint: "/subscription-plans",
			method:   http.MethodPost,
			requestBody: models.CreatePlanRequest{
				Name:     "Invalid Plan",
				PlanType: "monthly",
				Price:    -10.0,
				IsActive: true,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request data",
		},
		{
			name:     "Create trial plan without duration",
			endpoint: "/subscription-plans",
			method:   http.MethodPost,
			requestBody: models.CreatePlanRequest{
				Name:     "Invalid Trial",
				PlanType: "trial",
				Price:    0.0,
				IsActive: true,
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to create plan",
		},
		{
			name:     "Update plan with invalid ID",
			endpoint: "/subscription-plans/invalid",
			method:   http.MethodPut,
			requestBody: models.UpdatePlanRequest{
				Name: func() *string { v := "Updated Plan"; return &v }(),
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid plan ID",
		},
		{
			name:           "Get plan with invalid ID",
			endpoint:       "/subscription-plans/invalid",
			method:         http.MethodGet,
			requestBody:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid plan ID",
		},
		{
			name:           "Delete plan with invalid ID",
			endpoint:       "/subscription-plans/invalid",
			method:         http.MethodDelete,
			requestBody:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid plan ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := utils.SetupMockSubscriptionPlanRepository(t)
			service := services.NewSubscriptionPlanService(mockRepo)
			handler := handlers.NewSubscriptionPlanHandler(service)

			// Create request
			var requestBody []byte
			if tt.requestBody != nil {
				requestBody, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(tt.method, tt.endpoint, bytes.NewBuffer(requestBody))
			if requestBody != nil {
				req.Header.Set("Content-Type", "application/json")
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Extract ID from endpoint if present
			if tt.endpoint != "/subscription-plans" {
				// Extract ID from endpoint like "/subscription-plans/1" -> "1"
				parts := strings.Split(tt.endpoint, "/")
				if len(parts) > 2 {
					c.Params = gin.Params{{Key: "id", Value: parts[len(parts)-1]}}
				}
			}

			// Execute based on method
			switch tt.method {
			case http.MethodPost:
				handler.CreatePlan(c)
			case http.MethodGet:
				handler.GetPlan(c)
			case http.MethodPut:
				handler.UpdatePlan(c)
			case http.MethodDelete:
				handler.DeletePlan(c)
			}

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}
}

// TestSubscriptionPlanPaginationIntegration tests pagination functionality
func TestSubscriptionPlanPaginationIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    string
		expectedLimit  int
		expectedOffset int
		expectedActive *bool
	}{
		{
			name:           "Default pagination",
			queryParams:    "",
			expectedLimit:  10,
			expectedOffset: 0,
			expectedActive: nil,
		},
		{
			name:           "Custom pagination",
			queryParams:    "?limit=5&offset=10",
			expectedLimit:  5,
			expectedOffset: 10,
			expectedActive: nil,
		},
		{
			name:           "Filter by active status",
			queryParams:    "?is_active=true",
			expectedLimit:  10,
			expectedOffset: 0,
			expectedActive: func() *bool { v := true; return &v }(),
		},
		{
			name:           "Filter by inactive status",
			queryParams:    "?is_active=false",
			expectedLimit:  10,
			expectedOffset: 0,
			expectedActive: func() *bool { v := false; return &v }(),
		},
		{
			name:           "Invalid parameters fallback to defaults",
			queryParams:    "?limit=invalid&offset=invalid",
			expectedLimit:  10,
			expectedOffset: 0,
			expectedActive: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := utils.SetupMockSubscriptionPlanRepository(t)
			service := services.NewSubscriptionPlanService(mockRepo)
			handler := handlers.NewSubscriptionPlanHandler(service)

			// Mock expectations
			plans := []*models.SubscriptionPlan{
				utils.CreateMockSubscriptionPlan(1, "Plan 1", "monthly", 29.99),
				utils.CreateMockSubscriptionPlan(2, "Plan 2", "annual", 299.99),
			}
			mockRepo.On("GetAll", tt.expectedLimit, tt.expectedOffset, tt.expectedActive).Return(plans, nil)
			mockRepo.On("Count", tt.expectedActive).Return(2, nil)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/subscription-plans"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Execute
			handler.GetAllPlans(c)

			// Assert
			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.NotNil(t, response["data"])
			assert.NotNil(t, response["pagination"])

			pagination := response["pagination"].(map[string]interface{})
			assert.Equal(t, float64(tt.expectedLimit), pagination["limit"])
			assert.Equal(t, float64(tt.expectedOffset), pagination["offset"])
			assert.Equal(t, float64(2), pagination["total"])

			utils.AssertSubscriptionPlanRepositoryExpectations(t, mockRepo)
		})
	}
}

// TestSubscriptionPlanErrorHandlingIntegration tests error handling scenarios
func TestSubscriptionPlanErrorHandlingIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		endpoint       string
		method         string
		requestBody    interface{}
		mockSetup      func(*utils.MockSubscriptionPlanRepository)
		expectedStatus int
		expectedError  string
	}{
		{
			name:     "Create plan with database error",
			endpoint: "/subscription-plans",
			method:   http.MethodPost,
			requestBody: models.CreatePlanRequest{
				Name:     "Test Plan",
				PlanType: "monthly",
				Price:    29.99,
				IsActive: true,
			},
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				mockRepo.On("Create", mock.AnythingOfType("*models.SubscriptionPlan")).Return(fmt.Errorf("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to create plan",
		},
		{
			name:        "Get non-existent plan",
			endpoint:    "/subscription-plans/999",
			method:      http.MethodGet,
			requestBody: nil,
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				mockRepo.On("GetByID", 999).Return(nil, fmt.Errorf("subscription plan not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Plan not found",
		},
		{
			name:     "Update non-existent plan",
			endpoint: "/subscription-plans/999",
			method:   http.MethodPut,
			requestBody: models.UpdatePlanRequest{
				Name: func() *string { v := "Updated Plan"; return &v }(),
			},
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				mockRepo.On("GetByID", 999).Return(nil, fmt.Errorf("subscription plan not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to update plan",
		},
		{
			name:        "Delete plan with active subscriptions",
			endpoint:    "/subscription-plans/1",
			method:      http.MethodDelete,
			requestBody: nil,
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				existingPlan := utils.CreateMockSubscriptionPlan(1, "Test Plan", "monthly", 29.99)
				mockRepo.On("GetByID", 1).Return(existingPlan, nil)
				mockRepo.On("HasActiveSubscriptions", 1).Return(true, nil)
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "Cannot delete plan",
		},
		{
			name:        "Delete non-existent plan",
			endpoint:    "/subscription-plans/999",
			method:      http.MethodDelete,
			requestBody: nil,
			mockSetup: func(mockRepo *utils.MockSubscriptionPlanRepository) {
				mockRepo.On("GetByID", 999).Return(nil, fmt.Errorf("subscription plan not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Plan not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := utils.SetupMockSubscriptionPlanRepository(t)
			tt.mockSetup(mockRepo)
			service := services.NewSubscriptionPlanService(mockRepo)
			handler := handlers.NewSubscriptionPlanHandler(service)

			// Create request
			var requestBody []byte
			if tt.requestBody != nil {
				requestBody, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(tt.method, tt.endpoint, bytes.NewBuffer(requestBody))
			if requestBody != nil {
				req.Header.Set("Content-Type", "application/json")
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Extract ID from endpoint if present
			if tt.endpoint != "/subscription-plans" {
				parts := strings.Split(tt.endpoint, "/")
				if len(parts) > 2 {
					c.Params = gin.Params{{Key: "id", Value: parts[len(parts)-1]}}
				}
			}

			// Execute based on method
			switch tt.method {
			case http.MethodPost:
				handler.CreatePlan(c)
			case http.MethodGet:
				handler.GetPlan(c)
			case http.MethodPut:
				handler.UpdatePlan(c)
			case http.MethodDelete:
				handler.DeletePlan(c)
			}

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Contains(t, response["error"], tt.expectedError)

			utils.AssertSubscriptionPlanRepositoryExpectations(t, mockRepo)
		})
	}
}
