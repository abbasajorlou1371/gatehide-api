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
		"DELETE FROM user_roles",
		"DELETE FROM role_permissions",
		"DELETE FROM permissions",
		"DELETE FROM roles",
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

	// Disable foreign key checks temporarily to avoid constraint issues
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// Clean up test data and reset auto-increment
	queries := []string{
		"DELETE FROM user_roles",
		"DELETE FROM role_permissions",
		"DELETE FROM permissions",
		"DELETE FROM roles",
		"DELETE FROM user_sessions",
		"DELETE FROM users",
		"DELETE FROM admins",
		"DELETE FROM gamenets",
		"DELETE FROM migrations",
		"ALTER TABLE user_roles AUTO_INCREMENT = 1",
		"ALTER TABLE role_permissions AUTO_INCREMENT = 1",
		"ALTER TABLE permissions AUTO_INCREMENT = 1",
		"ALTER TABLE roles AUTO_INCREMENT = 1",
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

	// Re-enable foreign key checks
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}

// CreateTestUser creates a test user in the database
func CreateTestUser(t *testing.T, db *sql.DB, email, password, name string) *models.User {
	hashedPassword, err := models.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Generate unique mobile number based on email hash and timestamp
	mobile := fmt.Sprintf("+1%09d", len(email)*1000+len(name)+int(time.Now().UnixNano()%10000000))

	query := `
		INSERT INTO users (name, mobile, email, password, image, balance, debt, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, 0.00, 0.00, NOW(), NOW())
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
	mobile := fmt.Sprintf("+1%09d", len(email)*2000+len(name)+int(time.Now().UnixNano()%10000000))

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
			balance DECIMAL(10, 2) DEFAULT 0.00 NOT NULL,
			debt DECIMAL(10, 2) DEFAULT 0.00 NOT NULL,
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

	// Create permissions tables for RBAC
	permissionsTable := `
		CREATE TABLE IF NOT EXISTS permissions (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL UNIQUE,
			description TEXT,
			resource VARCHAR(50) NOT NULL,
			action VARCHAR(50) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			
			INDEX idx_resource_action (resource, action),
			INDEX idx_name (name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	if _, err := db.Exec(permissionsTable); err != nil {
		return fmt.Errorf("failed to create permissions table: %w", err)
	}

	rolesTable := `
		CREATE TABLE IF NOT EXISTS roles (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(50) NOT NULL UNIQUE,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			
			INDEX idx_name (name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	if _, err := db.Exec(rolesTable); err != nil {
		return fmt.Errorf("failed to create roles table: %w", err)
	}

	rolePermissionsTable := `
		CREATE TABLE IF NOT EXISTS role_permissions (
			id INT AUTO_INCREMENT PRIMARY KEY,
			role_id INT NOT NULL,
			permission_id INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			
			FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
			FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
			
			UNIQUE KEY unique_role_permission (role_id, permission_id),
			INDEX idx_role_id (role_id),
			INDEX idx_permission_id (permission_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	if _, err := db.Exec(rolePermissionsTable); err != nil {
		return fmt.Errorf("failed to create role_permissions table: %w", err)
	}

	// Create user_roles table
	userRolesTable := `
		CREATE TABLE IF NOT EXISTS user_roles (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			user_type ENUM('user', 'admin', 'gamenet') NOT NULL,
			role_id INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			
			FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
			
			UNIQUE KEY unique_user_role (user_id, user_type, role_id),
			INDEX idx_user_id_type (user_id, user_type),
			INDEX idx_role_id (role_id),
			INDEX idx_user_type (user_type)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	if _, err := db.Exec(userRolesTable); err != nil {
		return fmt.Errorf("failed to create user_roles table: %w", err)
	}

	// Insert roles
	roles := []string{
		"INSERT IGNORE INTO roles (name, description) VALUES ('administrator', 'System administrator with full access')",
		"INSERT IGNORE INTO roles (name, description) VALUES ('gamenet', 'Gaming center operator with limited access')",
		"INSERT IGNORE INTO roles (name, description) VALUES ('user', 'Regular user with basic access')",
	}

	for _, roleQuery := range roles {
		if _, err := db.Exec(roleQuery); err != nil {
			return fmt.Errorf("failed to insert role: %w", err)
		}
	}

	// Insert permissions
	permissions := []string{
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('dashboard:view', 'View dashboard', 'dashboard', 'view')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('gamenets:create', 'Create gamenets', 'gamenets', 'create')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('gamenets:read', 'View gamenets', 'gamenets', 'read')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('gamenets:update', 'Update gamenets', 'gamenets', 'update')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('gamenets:delete', 'Delete gamenets', 'gamenets', 'delete')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('users:create', 'Create users', 'users', 'create')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('users:read', 'View users', 'users', 'read')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('users:update', 'Update users', 'users', 'update')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('users:delete', 'Delete users', 'users', 'delete')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('subscription_plans:create', 'Create subscription plans', 'subscription_plans', 'create')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('subscription_plans:read', 'View subscription plans', 'subscription_plans', 'read')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('subscription_plans:update', 'Update subscription plans', 'subscription_plans', 'update')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('subscription_plans:delete', 'Delete subscription plans', 'subscription_plans', 'delete')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('analytics:view', 'View analytics', 'analytics', 'view')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('payments:view', 'View payments', 'payments', 'view')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('transactions:view', 'View transactions', 'transactions', 'view')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('invoices:view', 'View invoices', 'invoices', 'view')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('settings:manage', 'Manage settings', 'settings', 'manage')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('support:access', 'Access support', 'support', 'access')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('reservation:manage', 'Manage reservations', 'reservation', 'manage')",
		"INSERT IGNORE INTO permissions (name, description, resource, action) VALUES ('wallet:view', 'View wallet', 'wallet', 'view')",
	}

	for _, permQuery := range permissions {
		if _, err := db.Exec(permQuery); err != nil {
			return fmt.Errorf("failed to insert permission: %w", err)
		}
	}

	// Assign permissions to roles - use individual transactions to avoid deadlocks
	rolePermissions := []string{
		// Administrator permissions
		"INSERT IGNORE INTO role_permissions (role_id, permission_id) SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'administrator' AND p.name IN ('dashboard:view', 'gamenets:create', 'gamenets:read', 'gamenets:update', 'gamenets:delete', 'users:create', 'users:read', 'users:update', 'users:delete', 'subscription_plans:create', 'subscription_plans:read', 'subscription_plans:update', 'subscription_plans:delete', 'analytics:view', 'payments:view', 'transactions:view', 'invoices:view', 'settings:manage', 'support:access')",
		// Gamenet permissions
		"INSERT IGNORE INTO role_permissions (role_id, permission_id) SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'gamenet' AND p.name IN ('dashboard:view', 'users:create', 'users:read', 'users:update', 'users:delete', 'analytics:view', 'transactions:view', 'payments:view', 'support:access', 'settings:manage')",
		// User permissions
		"INSERT IGNORE INTO role_permissions (role_id, permission_id) SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'user' AND p.name IN ('reservation:manage', 'support:access', 'settings:manage', 'wallet:view')",
	}

	for _, rpQuery := range rolePermissions {
		// Use individual transactions to avoid deadlocks
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := tx.Exec(rpQuery); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to assign role permissions: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit role permissions transaction: %w", err)
		}
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
