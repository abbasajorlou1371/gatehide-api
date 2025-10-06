package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/internal/migrations"
)

const (
	migrationsDir = "database/migrations"
)

func main() {
	var (
		command = flag.String("command", "status", "Migration command: status, up, down, create")
		name    = flag.String("name", "", "Migration name (for create command)")
		steps   = flag.Int("steps", 1, "Number of migrations to run (for up/down commands)")
	)
	flag.Parse()

	// Load configuration
	cfg := config.Load()

	// Get migrations directory path
	migrationsPath, err := getMigrationsPath()
	if err != nil {
		log.Fatalf("Failed to get migrations path: %v", err)
	}

	// Create migration runner
	runner, err := migrations.NewMySQLRunner(cfg)
	if err != nil {
		log.Fatalf("Failed to create migration runner: %v", err)
	}
	defer runner.Close()

	// Execute command
	switch *command {
	case "status":
		if err := runStatus(runner, migrationsPath); err != nil {
			log.Fatalf("Status command failed: %v", err)
		}
	case "up":
		if err := runUp(runner, migrationsPath, *steps); err != nil {
			log.Fatalf("Up command failed: %v", err)
		}
	case "down":
		if err := runDown(runner, migrationsPath, *steps); err != nil {
			log.Fatalf("Down command failed: %v", err)
		}
	case "create":
		if *name == "" {
			log.Fatal("Migration name is required for create command")
		}
		if err := runCreate(*name, migrationsPath); err != nil {
			log.Fatalf("Create command failed: %v", err)
		}
	default:
		log.Fatalf("Unknown command: %s", *command)
	}
}

func runStatus(runner migrations.MigrationRunner, migrationsPath string) error {
	// Create migration table if it doesn't exist
	if err := runner.CreateMigrationTable(); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	// Get applied migrations
	applied, err := runner.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Get available migration files
	available, err := migrations.LoadMigrationFiles(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to load migration files: %w", err)
	}

	// Create a map of applied migrations for quick lookup
	appliedMap := make(map[string]bool)
	for _, m := range applied {
		appliedMap[m.Version] = true
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")

	if len(available) == 0 {
		fmt.Println("No migration files found.")
		return nil
	}

	fmt.Printf("%-20s %-30s %-15s\n", "Version", "Description", "Status")
	fmt.Println(strings.Repeat("-", 65))

	for _, migration := range available {
		status := "PENDING"
		if appliedMap[migration.Version] {
			status = "APPLIED"
		}
		fmt.Printf("%-20s %-30s %-15s\n", migration.Version, migration.Description, status)
	}

	return nil
}

func runUp(runner migrations.MigrationRunner, migrationsPath string, steps int) error {
	// Create migration table if it doesn't exist
	if err := runner.CreateMigrationTable(); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	// Get applied migrations
	applied, err := runner.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Create a map of applied migrations for quick lookup
	appliedMap := make(map[string]bool)
	for _, m := range applied {
		appliedMap[m.Version] = true
	}

	// Get available migration files
	available, err := migrations.LoadMigrationFiles(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to load migration files: %w", err)
	}

	// Find pending migrations
	var pending []migrations.MigrationFile
	for _, migration := range available {
		if !appliedMap[migration.Version] {
			pending = append(pending, migration)
		}
	}

	if len(pending) == 0 {
		fmt.Println("No pending migrations.")
		return nil
	}

	// Limit by steps
	if steps > len(pending) {
		steps = len(pending)
	}

	fmt.Printf("Applying %d migration(s)...\n", steps)

	// Apply migrations
	for i := 0; i < steps; i++ {
		migration := pending[i]
		fmt.Printf("Applying migration %s: %s\n", migration.Version, migration.Description)

		if err := runner.ApplyMigration(migration.Version, migration.Description, migration.UpSQL); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}

		fmt.Printf("✅ Migration %s applied successfully\n", migration.Version)
	}

	return nil
}

func runDown(runner migrations.MigrationRunner, migrationsPath string, steps int) error {
	// Get applied migrations
	applied, err := runner.GetAppliedMigrations()
	if err != nil {
		// If migrations table doesn't exist, there are no migrations to rollback
		if strings.Contains(err.Error(), "doesn't exist") || strings.Contains(err.Error(), "Table") {
			fmt.Println("No applied migrations to rollback.")
			return nil
		}
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(applied) == 0 {
		fmt.Println("No applied migrations to rollback.")
		return nil
	}

	// Limit by steps
	if steps > len(applied) {
		steps = len(applied)
	}

	// Get available migration files
	available, err := migrations.LoadMigrationFiles(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to load migration files: %w", err)
	}

	// Create a map of available migrations for quick lookup
	availableMap := make(map[string]migrations.MigrationFile)
	for _, migration := range available {
		availableMap[migration.Version] = migration
	}

	fmt.Printf("Rolling back %d migration(s)...\n", steps)

	// Rollback migrations (from latest to oldest)
	for i := len(applied) - 1; i >= len(applied)-steps; i-- {
		migration := applied[i]

		// Get migration file for rollback SQL
		migrationFile, exists := availableMap[migration.Version]
		if !exists {
			return fmt.Errorf("migration file not found for version %s", migration.Version)
		}

		fmt.Printf("Rolling back migration %s: %s\n", migration.Version, migration.Description)

		if err := runner.RollbackMigration(migration.Version, migrationFile.DownSQL); err != nil {
			return fmt.Errorf("failed to rollback migration %s: %w", migration.Version, err)
		}

		fmt.Printf("✅ Migration %s rolled back successfully\n", migration.Version)
	}

	return nil
}

func runCreate(name, migrationsPath string) error {
	// Create migrations directory if it doesn't exist
	if err := os.MkdirAll(migrationsPath, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Generate timestamp for version
	timestamp := fmt.Sprintf("%d", getCurrentTimestamp())
	version := fmt.Sprintf("%s_%s", timestamp, sanitizeName(name))

	// Create migration file
	filename := fmt.Sprintf("%s.sql", version)
	filepath := filepath.Join(migrationsPath, filename)

	template := fmt.Sprintf(`-- version: %s
-- description: %s

-- UP
-- Add your migration SQL here
-- Example:
-- CREATE TABLE example (
--     id INT AUTO_INCREMENT PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- DOWN
-- Add your rollback SQL here
-- Example:
-- DROP TABLE IF EXISTS example;
`, version, name)

	if err := os.WriteFile(filepath, []byte(template), 0644); err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}

	fmt.Printf("✅ Created migration file: %s\n", filepath)
	return nil
}

func getMigrationsPath() (string, error) {
	// Try to find the project root
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up the directory tree to find go.mod
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, migrationsDir), nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// Fallback to current directory
	return migrationsDir, nil
}

func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}

func sanitizeName(name string) string {
	// Replace spaces and special characters with underscores
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")

	// Remove any non-alphanumeric characters except underscores
	var result strings.Builder
	for _, char := range name {
		if (char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' {
			result.WriteRune(char)
		}
	}

	return strings.ToLower(result.String())
}
