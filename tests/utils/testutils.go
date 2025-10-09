package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

// TestConfig creates a test configuration
func TestConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			Port:    "8081",
			GinMode: "test",
		},
		App: config.AppConfig{
			Name:    "GateHide API Test",
			Version: "1.0.0-test",
		},
		Security: config.SecurityConfig{
			APISecret:     "test-api-secret",
			JWTSecret:     "test-jwt-secret-key-for-testing-only",
			JWTExpiration: 1, // 1 hour for tests
		},
		Database: config.DatabaseConfig{
			Host:     getEnv("TEST_DB_HOST", "localhost"),
			Port:     getEnv("TEST_DB_PORT", "3306"),
			User:     getEnv("TEST_DB_USER", "root"),
			Password: getEnv("TEST_DB_PASSWORD", ""),
			DBName:   getEnv("TEST_DB_NAME", "gatehide_test"),
			SSLMode:  "false",
			Driver:   "mysql",
		},
	}
}

// TestDB creates a test database connection
func TestDB(t *testing.T) *sql.DB {
	cfg := TestConfig()

	db, err := sql.Open(cfg.Database.Driver, cfg.GetDSN())
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	return db
}

// SetupTestDB sets up a test database with migrations
func SetupTestDB(t *testing.T) *sql.DB {
	db := TestDB(t)

	// Run migrations for test database
	if err := runTestMigrations(db); err != nil {
		t.Fatalf("Failed to run test migrations: %v", err)
	}

	return db
}

// CleanupTestDB cleans up test data
func CleanupTestDB(t *testing.T, db *sql.DB) {
	// Check if database is still open
	if err := db.Ping(); err != nil {
		log.Printf("Warning: database connection is closed, skipping cleanup: %v", err)
		return
	}

	// Clean up test data in reverse order to avoid foreign key constraints
	queries := []string{
		"DELETE FROM user_sessions",
		"DELETE FROM users",
		"DELETE FROM admins",
		"DELETE FROM gamenets",
		"DELETE FROM migrations",
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			log.Printf("Warning: failed to clean up test data: %v", err)
		}
	}
}

// CleanupTestDBForce cleans up test data and resets auto-increment
func CleanupTestDBForce(t *testing.T, db *sql.DB) {
	// Check if database is still open
	if err := db.Ping(); err != nil {
		log.Printf("Warning: database connection is closed, skipping cleanup: %v", err)
		return
	}

	// Clean up test data and reset auto-increment
	queries := []string{
		"DELETE FROM user_sessions",
		"DELETE FROM users",
		"DELETE FROM admins",
		"DELETE FROM gamenets",
		"DELETE FROM migrations",
		"ALTER TABLE user_sessions AUTO_INCREMENT = 1",
		"ALTER TABLE users AUTO_INCREMENT = 1",
		"ALTER TABLE admins AUTO_INCREMENT = 1",
		"ALTER TABLE gamenets AUTO_INCREMENT = 1",
		"ALTER TABLE migrations AUTO_INCREMENT = 1",
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			log.Printf("Warning: failed to clean up test data: %v", err)
		}
	}
}

// CreateTestUser creates a test user in the database
func CreateTestUser(t *testing.T, db *sql.DB, email, password, name string) *models.User {
	hashedPassword, err := models.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Generate unique mobile number based on email hash and timestamp
	mobile := fmt.Sprintf("+1%09d", len(email)*1000+len(name)+int(time.Now().UnixNano()%1000000))

	query := `
		INSERT INTO users (name, mobile, email, password, image, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := db.Exec(query, name, mobile, email, hashedPassword, nil)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	return &models.User{
		ID:       int(id),
		Name:     name,
		Mobile:   mobile,
		Email:    email,
		Password: hashedPassword,
	}
}

// CreateTestAdmin creates a test admin in the database
func CreateTestAdmin(t *testing.T, db *sql.DB, email, password, name string) *models.Admin {
	hashedPassword, err := models.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Generate unique mobile number based on email hash and timestamp
	mobile := fmt.Sprintf("+1%09d", len(email)*2000+len(name)+int(time.Now().UnixNano()%1000000))

	query := `
		INSERT INTO admins (name, mobile, email, password, image, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := db.Exec(query, name, mobile, email, hashedPassword, nil)
	if err != nil {
		t.Fatalf("Failed to create test admin: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get admin ID: %v", err)
	}

	return &models.Admin{
		ID:       int(id),
		Name:     name,
		Mobile:   mobile,
		Email:    email,
		Password: hashedPassword,
	}
}

// runTestMigrations runs basic migrations for testing
func runTestMigrations(db *sql.DB) error {
	// Create migrations table
	migrationsTable := `
		CREATE TABLE IF NOT EXISTS migrations (
			id INT AUTO_INCREMENT PRIMARY KEY,
			version VARCHAR(255) NOT NULL UNIQUE,
			description TEXT,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	if _, err := db.Exec(migrationsTable); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Create users table
	usersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			mobile VARCHAR(20) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			image VARCHAR(500) NULL,
			last_login_at TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			
			INDEX idx_email (email),
			INDEX idx_mobile (mobile),
			INDEX idx_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	if _, err := db.Exec(usersTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create admins table
	adminsTable := `
		CREATE TABLE IF NOT EXISTS admins (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			mobile VARCHAR(20) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			image VARCHAR(500) NULL,
			last_login_at TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			
			INDEX idx_email (email),
			INDEX idx_mobile (mobile),
			INDEX idx_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	if _, err := db.Exec(adminsTable); err != nil {
		return fmt.Errorf("failed to create admins table: %w", err)
	}

	// Create gamenets table
	gamenetsTable := `
		CREATE TABLE IF NOT EXISTS gamenets (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			owner_name VARCHAR(255) NOT NULL,
			owner_mobile VARCHAR(20) NOT NULL,
			address TEXT NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			license_attachment VARCHAR(500) NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			
			INDEX idx_email (email),
			INDEX idx_owner_mobile (owner_mobile),
			INDEX idx_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	if _, err := db.Exec(gamenetsTable); err != nil {
		return fmt.Errorf("failed to create gamenets table: %w", err)
	}

	// Create user_sessions table
	userSessionsTable := `
		CREATE TABLE IF NOT EXISTS user_sessions (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			user_type ENUM('user', 'admin', 'gamenet') NOT NULL,
			session_token VARCHAR(500) NOT NULL UNIQUE,
			device_info TEXT NULL,
			ip_address VARCHAR(45) NULL,
			user_agent TEXT NULL,
			is_active BOOLEAN DEFAULT TRUE,
			last_activity_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP NOT NULL,
			
			INDEX idx_user_id (user_id),
			INDEX idx_user_type (user_type),
			INDEX idx_session_token (session_token),
			INDEX idx_is_active (is_active),
			INDEX idx_expires_at (expires_at),
			INDEX idx_last_activity (last_activity_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	if _, err := db.Exec(userSessionsTable); err != nil {
		return fmt.Errorf("failed to create user_sessions table: %w", err)
	}

	return nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SkipIfNoDB skips the test if no database is available
func SkipIfNoDB(t *testing.T) {
	if os.Getenv("SKIP_DB_TESTS") == "true" {
		t.Skip("Skipping database tests")
	}
}
