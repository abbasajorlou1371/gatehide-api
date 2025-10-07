package unit

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gatehide/gatehide-api/internal/handlers"
	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/tests/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSubscriptionPlanService is a mock implementation of SubscriptionPlanServiceInterface
type MockSubscriptionPlanService struct {
	mock.Mock
}

func (m *MockSubscriptionPlanService) CreatePlan(req *models.CreatePlanRequest) (*models.PlanResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlanResponse), args.Error(1)
}

func (m *MockSubscriptionPlanService) GetPlan(id int) (*models.PlanResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlanResponse), args.Error(1)
}

func (m *MockSubscriptionPlanService) GetAllPlans(limit, offset int, isActive *bool) ([]*models.PlanResponse, int, error) {
	args := m.Called(limit, offset, isActive)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.PlanResponse), args.Int(1), args.Error(2)
}

func (m *MockSubscriptionPlanService) UpdatePlan(id int, req *models.UpdatePlanRequest) (*models.PlanResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlanResponse), args.Error(1)
}

func (m *MockSubscriptionPlanService) DeletePlan(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestSubscriptionPlanHandler_CreatePlan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockSubscriptionPlanService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful plan creation",
			requestBody: models.CreatePlanRequest{
				Name:     "Basic Monthly",
				PlanType: "monthly",
				Price:    29.99,
				IsActive: true,
			},
			mockSetup: func(mockService *MockSubscriptionPlanService) {
				response := utils.CreateMockPlanResponse(1, "Basic Monthly", "monthly", 29.99)
				mockService.On("CreatePlan", mock.AnythingOfType("*models.CreatePlanRequest")).Return(response, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid JSON request",
			requestBody: map[string]interface{}{
				"name":      "Basic Monthly",
				"plan_type": "monthly",
				"price":     "invalid_price", // Invalid type
			},
			mockSetup:      func(mockService *MockSubscriptionPlanService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request data",
		},
		{
			name: "missing required fields",
			requestBody: map[string]interface{}{
				"name": "Basic Monthly",
				// Missing plan_type and price
			},
			mockSetup:      func(mockService *MockSubscriptionPlanService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request data",
		},
		{
			name: "service error during creation",
			requestBody: models.CreatePlanRequest{
				Name:     "Basic Monthly",
				PlanType: "monthly",
				Price:    29.99,
				IsActive: true,
			},
			mockSetup: func(mockService *MockSubscriptionPlanService) {
				mockService.On("CreatePlan", mock.AnythingOfType("*models.CreatePlanRequest")).Return(nil, errors.New("validation error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to create plan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockSubscriptionPlanService)
			tt.mockSetup(mockService)
			handler := handlers.NewSubscriptionPlanHandler(mockService)

			// Create request body
			var requestBody []byte
			if tt.requestBody != nil {
				requestBody, _ = json.Marshal(tt.requestBody)
			}

			// Create HTTP request
			req := httptest.NewRequest(http.MethodPost, "/subscription-plans", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Setup Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Execute
			handler.CreatePlan(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			} else if tt.expectedStatus == http.StatusCreated {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Plan created successfully", response["message"])
				assert.NotNil(t, response["data"])
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestSubscriptionPlanHandler_GetPlan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		planID         string
		mockSetup      func(*MockSubscriptionPlanService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful plan retrieval",
			planID: "1",
			mockSetup: func(mockService *MockSubscriptionPlanService) {
				response := utils.CreateMockPlanResponse(1, "Basic Monthly", "monthly", 29.99)
				mockService.On("GetPlan", 1).Return(response, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid plan ID",
			planID:         "invalid",
			mockSetup:      func(mockService *MockSubscriptionPlanService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid plan ID",
		},
		{
			name:   "plan not found",
			planID: "999",
			mockSetup: func(mockService *MockSubscriptionPlanService) {
				mockService.On("GetPlan", 999).Return(nil, errors.New("subscription plan not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Plan not found",
		},
		{
			name:   "service error",
			planID: "1",
			mockSetup: func(mockService *MockSubscriptionPlanService) {
				mockService.On("GetPlan", 1).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Plan not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockSubscriptionPlanService)
			tt.mockSetup(mockService)
			handler := handlers.NewSubscriptionPlanHandler(mockService)

			// Create HTTP request
			req := httptest.NewRequest(http.MethodGet, "/subscription-plans/"+tt.planID, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Setup Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: tt.planID}}

			// Execute
			handler.GetPlan(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			} else if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotNil(t, response["data"])
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestSubscriptionPlanHandler_GetAllPlans(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*MockSubscriptionPlanService)
		expectedStatus int
		expectedError  string
		expectedCount  int
	}{
		{
			name:        "successful plans retrieval with default pagination",
			queryParams: "",
			mockSetup: func(mockService *MockSubscriptionPlanService) {
				plans := []*models.PlanResponse{
					utils.CreateMockPlanResponse(1, "Basic Monthly", "monthly", 29.99),
					utils.CreateMockPlanResponse(2, "Premium Annual", "annual", 299.99),
				}
				mockService.On("GetAllPlans", 10, 0, (*bool)(nil)).Return(plans, 2, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:        "successful plans retrieval with custom pagination",
			queryParams: "?limit=5&offset=10",
			mockSetup: func(mockService *MockSubscriptionPlanService) {
				plans := []*models.PlanResponse{
					utils.CreateMockPlanResponse(1, "Basic Monthly", "monthly", 29.99),
				}
				mockService.On("GetAllPlans", 5, 10, (*bool)(nil)).Return(plans, 1, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:        "successful active plans retrieval",
			queryParams: "?is_active=true",
			mockSetup: func(mockService *MockSubscriptionPlanService) {
				plans := []*models.PlanResponse{
					utils.CreateMockPlanResponse(1, "Basic Monthly", "monthly", 29.99),
				}
				active := true
				mockService.On("GetAllPlans", 10, 0, &active).Return(plans, 1, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:        "invalid limit parameter",
			queryParams: "?limit=invalid",
			mockSetup: func(mockService *MockSubscriptionPlanService) {
				plans := []*models.PlanResponse{
					utils.CreateMockPlanResponse(1, "Basic Monthly", "monthly", 29.99),
				}
				mockService.On("GetAllPlans", 10, 0, (*bool)(nil)).Return(plans, 1, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:        "invalid offset parameter",
			queryParams: "?offset=invalid",
			mockSetup: func(mockService *MockSubscriptionPlanService) {
				plans := []*models.PlanResponse{
					utils.CreateMockPlanResponse(1, "Basic Monthly", "monthly", 29.99),
				}
				mockService.On("GetAllPlans", 10, 0, (*bool)(nil)).Return(plans, 1, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:        "service error",
			queryParams: "",
			mockSetup: func(mockService *MockSubscriptionPlanService) {
				mockService.On("GetAllPlans", 10, 0, (*bool)(nil)).Return(nil, 0, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to get plans",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockSubscriptionPlanService)
			tt.mockSetup(mockService)
			handler := handlers.NewSubscriptionPlanHandler(mockService)

			// Create HTTP request
			req := httptest.NewRequest(http.MethodGet, "/subscription-plans"+tt.queryParams, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Setup Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Execute
			handler.GetAllPlans(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			} else if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotNil(t, response["data"])
				assert.NotNil(t, response["pagination"])

				if data, ok := response["data"].([]interface{}); ok {
					assert.Len(t, data, tt.expectedCount)
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestSubscriptionPlanHandler_UpdatePlan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		planID         string
		requestBody    interface{}
		mockSetup      func(*MockSubscriptionPlanService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful plan update",
			planID: "1",
			requestBody: models.UpdatePlanRequest{
				Name:  func() *string { v := "Updated Basic Monthly"; return &v }(),
				Price: func() *float64 { v := 39.99; return &v }(),
			},
			mockSetup: func(mockService *MockSubscriptionPlanService) {
				response := utils.CreateMockPlanResponse(1, "Updated Basic Monthly", "monthly", 39.99)
				mockService.On("UpdatePlan", 1, mock.AnythingOfType("*models.UpdatePlanRequest")).Return(response, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid plan ID",
			planID:         "invalid",
			requestBody:    models.UpdatePlanRequest{},
			mockSetup:      func(mockService *MockSubscriptionPlanService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid plan ID",
		},
		{
			name:   "invalid JSON request",
			planID: "1",
			requestBody: map[string]interface{}{
				"price": "invalid_price", // Invalid type
			},
			mockSetup:      func(mockService *MockSubscriptionPlanService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request data",
		},
		{
			name:   "plan not found",
			planID: "999",
			requestBody: models.UpdatePlanRequest{
				Name: func() *string { v := "Updated Plan"; return &v }(),
			},
			mockSetup: func(mockService *MockSubscriptionPlanService) {
				mockService.On("UpdatePlan", 999, mock.AnythingOfType("*models.UpdatePlanRequest")).Return(nil, errors.New("subscription plan not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to update plan",
		},
		{
			name:   "service error",
			planID: "1",
			requestBody: models.UpdatePlanRequest{
				Name: func() *string { v := "Updated Plan"; return &v }(),
			},
			mockSetup: func(mockService *MockSubscriptionPlanService) {
				mockService.On("UpdatePlan", 1, mock.AnythingOfType("*models.UpdatePlanRequest")).Return(nil, errors.New("validation error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to update plan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockSubscriptionPlanService)
			tt.mockSetup(mockService)
			handler := handlers.NewSubscriptionPlanHandler(mockService)

			// Create request body
			var requestBody []byte
			if tt.requestBody != nil {
				requestBody, _ = json.Marshal(tt.requestBody)
			}

			// Create HTTP request
			req := httptest.NewRequest(http.MethodPut, "/subscription-plans/"+tt.planID, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Setup Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: tt.planID}}

			// Execute
			handler.UpdatePlan(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			} else if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Plan updated successfully", response["message"])
				assert.NotNil(t, response["data"])
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestSubscriptionPlanHandler_DeletePlan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		planID         string
		mockSetup      func(*MockSubscriptionPlanService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful plan deletion",
			planID: "1",
			mockSetup: func(mockService *MockSubscriptionPlanService) {
				mockService.On("DeletePlan", 1).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid plan ID",
			planID:         "invalid",
			mockSetup:      func(mockService *MockSubscriptionPlanService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid plan ID",
		},
		{
			name:   "plan with active subscriptions",
			planID: "1",
			mockSetup: func(mockService *MockSubscriptionPlanService) {
				mockService.On("DeletePlan", 1).Return(errors.New("cannot delete plan: plan has active subscriptions"))
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "Cannot delete plan",
		},
		{
			name:   "plan not found",
			planID: "999",
			mockSetup: func(mockService *MockSubscriptionPlanService) {
				mockService.On("DeletePlan", 999).Return(errors.New("subscription plan not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Plan not found",
		},
		{
			name:   "service error",
			planID: "1",
			mockSetup: func(mockService *MockSubscriptionPlanService) {
				mockService.On("DeletePlan", 1).Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to delete plan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockSubscriptionPlanService)
			tt.mockSetup(mockService)
			handler := handlers.NewSubscriptionPlanHandler(mockService)

			// Create HTTP request
			req := httptest.NewRequest(http.MethodDelete, "/subscription-plans/"+tt.planID, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Setup Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: tt.planID}}

			// Execute
			handler.DeletePlan(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			} else if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Plan deleted successfully", response["message"])
			}

			mockService.AssertExpectations(t)
		})
	}
}
