package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/repositories"
	"github.com/gatehide/gatehide-api/internal/routes"
	testutils "github.com/gatehide/gatehide-api/tests/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserIntegrationTestSuite struct {
	suite.Suite
	db     *sql.DB
	router *gin.Engine
	cfg    *config.Config
	token  string
}

func (suite *UserIntegrationTestSuite) SetupSuite() {
	// Use testutils for proper database setup
	suite.db = testutils.SetupTestDB(suite.T())

	// Load test configuration
	suite.cfg = testutils.TestConfig()

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Setup routes
	suite.router = gin.New()
	routes.SetupRoutes(suite.router, suite.cfg, suite.db)

	// Get authentication token
	suite.token = suite.getAuthToken()
}

func (suite *UserIntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		// Clean up admin users created during tests
		_, err := suite.db.Exec("DELETE FROM admins WHERE email LIKE 'admin%@test.com'")
		if err != nil {
			suite.T().Logf("Warning: Failed to clean admin users: %v", err)
		}
		suite.db.Close()
	}
}

func (suite *UserIntegrationTestSuite) SetupTest() {
	// Clean up users and admins tables before each test
	_, err := suite.db.Exec("DELETE FROM users WHERE email LIKE 'test%@example.com'")
	if err != nil {
		suite.T().Fatalf("Failed to clean users table: %v", err)
	}

	_, err = suite.db.Exec("DELETE FROM admins WHERE email LIKE 'admin%@test.com'")
	if err != nil {
		suite.T().Fatalf("Failed to clean admins table: %v", err)
	}

	// Refresh token for each test to ensure it's valid
	suite.token = suite.getAuthToken()
}

func (suite *UserIntegrationTestSuite) TearDownTest() {
	// Clean up after each test
	_, err := suite.db.Exec("DELETE FROM users WHERE email LIKE 'test%@example.com'")
	if err != nil {
		suite.T().Logf("Warning: Failed to clean users table: %v", err)
	}
}

func (suite *UserIntegrationTestSuite) getAuthToken() string {
	// Create a test admin user first with unique email and mobile
	adminEmail := fmt.Sprintf("admin%d@test.com", time.Now().UnixNano())
	adminMobile := fmt.Sprintf("0912345%05d", time.Now().UnixNano()%1000000)
	hashedPassword, _ := models.HashPassword("adminpass")
	result, err := suite.db.Exec("INSERT INTO admins (name, email, mobile, password) VALUES (?, ?, ?, ?)",
		"Test Admin", adminEmail, adminMobile, hashedPassword)
	if err != nil {
		suite.T().Fatalf("Failed to create admin user: %v", err)
	}

	// Get the newly created admin ID
	adminID, err := result.LastInsertId()
	if err != nil {
		suite.T().Fatalf("Failed to get admin ID: %v", err)
	}
	suite.assignAdministratorRole(int(adminID))

	// Login as admin to get token
	loginData := map[string]interface{}{
		"email":    adminEmail,
		"password": "adminpass",
	}

	jsonData, _ := json.Marshal(loginData)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		suite.T().Fatalf("Failed to login for tests: status %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	return data["token"].(string)
}

// assignAdministratorRole assigns the administrator role to an admin
func (suite *UserIntegrationTestSuite) assignAdministratorRole(adminID int) {
	// First get the administrator role ID
	var roleID int
	roleQuery := "SELECT id FROM roles WHERE name = 'administrator'"
	err := suite.db.QueryRow(roleQuery).Scan(&roleID)
	if err != nil {
		suite.T().Fatalf("Failed to get administrator role: %v", err)
	}

	// Insert the role assignment
	assignQuery := `
		INSERT INTO user_roles (user_id, user_type, role_id, created_at, updated_at)
		VALUES (?, 'admin', ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE updated_at = NOW()
	`

	_, err = suite.db.Exec(assignQuery, adminID, roleID)
	if err != nil {
		suite.T().Fatalf("Failed to assign administrator role: %v", err)
	}
}

func (suite *UserIntegrationTestSuite) TestCreateUser() {
	t := suite.T()

	// Test data
	userData := map[string]interface{}{
		"name":   "Test User",
		"email":  "testuser@example.com",
		"mobile": "09123456789",
	}

	jsonData, _ := json.Marshal(userData)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	// Execute request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code, "Expected status 201 Created")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "User created successfully", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, userData["name"], data["name"])
	assert.Equal(t, userData["email"], data["email"])
	assert.Equal(t, userData["mobile"], data["mobile"])
}

func (suite *UserIntegrationTestSuite) TestCreateUserDuplicateEmail() {
	t := suite.T()

	// Create first user
	userData := map[string]interface{}{
		"name":   "Test User 1",
		"email":  "testduplicate@example.com",
		"mobile": "09123456789",
	}

	jsonData, _ := json.Marshal(userData)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Try to create second user with same email
	userData2 := map[string]interface{}{
		"name":   "Test User 2",
		"email":  "testduplicate@example.com",
		"mobile": "09987654321",
	}

	jsonData2, _ := json.Marshal(userData2)

	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewBuffer(jsonData2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+suite.token)

	w2 := httptest.NewRecorder()
	suite.router.ServeHTTP(w2, req2)

	// Should return error
	assert.Equal(t, http.StatusBadRequest, w2.Code)

	var response map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "email already exists")
}

func (suite *UserIntegrationTestSuite) TestGetAllUsers() {
	t := suite.T()

	// Create a test user first
	userRepo := repositories.NewUserRepository(suite.db)
	hashedPassword, _ := models.HashPassword("password123")

	testUser := &models.User{
		Name:     "Test User",
		Email:    "testgetall@example.com",
		Mobile:   "09123456789",
		Password: hashedPassword,
	}
	userRepo.Create(testUser)

	// Get all users
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/", nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 1, "Should have at least one user")
}

func (suite *UserIntegrationTestSuite) TestGetUserByID() {
	t := suite.T()

	// Create a test user first
	userRepo := repositories.NewUserRepository(suite.db)
	hashedPassword, _ := models.HashPassword("password123")

	testUser := &models.User{
		Name:     "Test User",
		Email:    "testgetbyid@example.com",
		Mobile:   "09123456789",
		Password: hashedPassword,
	}
	userRepo.Create(testUser)

	// Get user by ID
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%d", testUser.ID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, testUser.Name, data["name"])
	assert.Equal(t, testUser.Email, data["email"])
}

func (suite *UserIntegrationTestSuite) TestUpdateUser() {
	t := suite.T()

	// Create a test user first
	userRepo := repositories.NewUserRepository(suite.db)
	hashedPassword, _ := models.HashPassword("password123")

	testUser := &models.User{
		Name:     "Test User",
		Email:    "testupdate@example.com",
		Mobile:   "09123456789",
		Password: hashedPassword,
	}
	userRepo.Create(testUser)

	// Update user
	updateData := map[string]interface{}{
		"name": "Updated User Name",
	}

	jsonData, _ := json.Marshal(updateData)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%d", testUser.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Updated User Name", data["name"])
	assert.Equal(t, testUser.Email, data["email"]) // Email should remain unchanged
}

func (suite *UserIntegrationTestSuite) TestDeleteUser() {
	t := suite.T()

	// Create a test user first
	userRepo := repositories.NewUserRepository(suite.db)
	hashedPassword, _ := models.HashPassword("password123")

	testUser := &models.User{
		Name:     "Test User",
		Email:    "testdelete@example.com",
		Mobile:   "09123456789",
		Password: hashedPassword,
	}
	userRepo.Create(testUser)

	// Delete user
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/users/%d", testUser.ID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "User deleted successfully", response["message"])

	// Verify user is deleted
	_, err = userRepo.GetByID(testUser.ID)
	assert.Error(t, err, "User should not exist after deletion")
}

func (suite *UserIntegrationTestSuite) TestSearchUsers() {
	t := suite.T()

	// Create test users
	userRepo := repositories.NewUserRepository(suite.db)
	hashedPassword, _ := models.HashPassword("password123")

	users := []*models.User{
		{Name: "John Doe", Email: "testsearch1@example.com", Mobile: "09111111111", Password: hashedPassword},
		{Name: "Jane Doe", Email: "testsearch2@example.com", Mobile: "09222222222", Password: hashedPassword},
		{Name: "Bob Smith", Email: "testsearch3@example.com", Mobile: "09333333333", Password: hashedPassword},
	}

	for _, user := range users {
		userRepo.Create(user)
	}

	// Search for "Doe"
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/?query=Doe", nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 2, "Should find at least 2 users with 'Doe'")
}

func (suite *UserIntegrationTestSuite) TestCreateUserInvalidData() {
	t := suite.T()

	testCases := []struct {
		name         string
		userData     map[string]interface{}
		expectedCode int
	}{
		{
			name: "Missing name",
			userData: map[string]interface{}{
				"email":  "test@example.com",
				"mobile": "09123456789",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Missing email",
			userData: map[string]interface{}{
				"name":   "Test User",
				"mobile": "09123456789",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Missing mobile",
			userData: map[string]interface{}{
				"name":  "Test User",
				"email": "test@example.com",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid email format",
			userData: map[string]interface{}{
				"name":   "Test User",
				"email":  "invalid-email",
				"mobile": "09123456789",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tc.userData)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+suite.token)

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestUserIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(UserIntegrationTestSuite))
}
