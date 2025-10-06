package migrations

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gatehide/gatehide-api/config"
	_ "github.com/go-sql-driver/mysql"
)

// MySQLRunner implements MigrationRunner for MySQL
type MySQLRunner struct {
	db     *sql.DB
	config *config.Config
}

// NewMySQLRunner creates a new MySQL migration runner
func NewMySQLRunner(cfg *config.Config) (*MySQLRunner, error) {
	// First check if database exists
	dbExists, err := checkDatabaseExists(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to check database existence: %w", err)
	}

	if !dbExists {
		// Check if we should auto-create database (for non-interactive environments)
		autoCreate := os.Getenv("DB_AUTO_CREATE") == "true"

		if !autoCreate && !promptCreateDatabase(cfg.Database.DBName) {
			return nil, fmt.Errorf("database creation cancelled by user")
		}

		// Create database
		if err := createDatabase(cfg); err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
	}

	// Connect to the specific database
	db, err := sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &MySQLRunner{
		db:     db,
		config: cfg,
	}, nil
}

// CreateMigrationTable creates the migrations tracking table
func (r *MySQLRunner) CreateMigrationTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id INT AUTO_INCREMENT PRIMARY KEY,
		version VARCHAR(255) NOT NULL UNIQUE,
		description VARCHAR(500) NOT NULL,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	_, err := r.db.ExecContext(context.Background(), query)
	return err
}

// GetAppliedMigrations returns all applied migrations
func (r *MySQLRunner) GetAppliedMigrations() ([]Migration, error) {
	query := "SELECT id, version, description, applied_at FROM migrations ORDER BY version"
	rows, err := r.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []Migration
	for rows.Next() {
		var m Migration
		if err := rows.Scan(&m.ID, &m.Version, &m.Description, &m.AppliedAt); err != nil {
			return nil, err
		}
		migrations = append(migrations, m)
	}

	return migrations, nil
}

// ApplyMigration applies a migration
func (r *MySQLRunner) ApplyMigration(version, description, upSQL string) error {
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration SQL
	if _, err := tx.ExecContext(context.Background(), upSQL); err != nil {
		return fmt.Errorf("failed to execute migration %s: %w", version, err)
	}

	// Record migration
	insertQuery := "INSERT INTO migrations (version, description) VALUES (?, ?)"
	if _, err := tx.ExecContext(context.Background(), insertQuery, version, description); err != nil {
		return fmt.Errorf("failed to record migration %s: %w", version, err)
	}

	return tx.Commit()
}

// RollbackMigration rolls back a migration
func (r *MySQLRunner) RollbackMigration(version, downSQL string) error {
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute rollback SQL
	if _, err := tx.ExecContext(context.Background(), downSQL); err != nil {
		return fmt.Errorf("failed to execute rollback for %s: %w", version, err)
	}

	// Remove migration record
	deleteQuery := "DELETE FROM migrations WHERE version = ?"
	if _, err := tx.ExecContext(context.Background(), deleteQuery, version); err != nil {
		return fmt.Errorf("failed to remove migration record %s: %w", version, err)
	}

	return tx.Commit()
}

// CheckDatabaseExists checks if the database exists
func (r *MySQLRunner) CheckDatabaseExists() (bool, error) {
	return checkDatabaseExists(r.config)
}

// CreateDatabase creates the database
func (r *MySQLRunner) CreateDatabase() error {
	return createDatabase(r.config)
}

// Close closes the database connection
func (r *MySQLRunner) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// LoadMigrationFiles loads migration files from the migrations directory
func LoadMigrationFiles(migrationsDir string) ([]MigrationFile, error) {
	var migrationFiles []MigrationFile

	// Check if migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return migrationFiles, nil
	}

	// Read directory
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}

	// Parse migration files
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migration, err := parseMigrationFile(filepath.Join(migrationsDir, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s: %w", file.Name(), err)
			}
			migrationFiles = append(migrationFiles, migration)
		}
	}

	// Sort by version
	sort.Slice(migrationFiles, func(i, j int) bool {
		return migrationFiles[i].Version < migrationFiles[j].Version
	})

	return migrationFiles, nil
}

// parseMigrationFile parses a migration file and extracts version, description, and SQL
func parseMigrationFile(filePath string) (MigrationFile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return MigrationFile{}, err
	}
	defer file.Close()

	var migration MigrationFile
	var currentSection string
	var upSQL, downSQL strings.Builder

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Parse sections first (before header comments)
		if strings.EqualFold(line, "-- UP") {
			currentSection = "up"
			continue
		}
		if strings.EqualFold(line, "-- DOWN") {
			currentSection = "down"
			continue
		}

		// Parse header comments (only if not in a section)
		if strings.HasPrefix(line, "--") && currentSection == "" {
			content := strings.TrimSpace(line[2:])
			if strings.HasPrefix(content, "version:") {
				migration.Version = strings.TrimSpace(content[8:])
			} else if strings.HasPrefix(content, "description:") {
				migration.Description = strings.TrimSpace(content[12:])
			}
			continue
		}

		// Add SQL to appropriate section (only if we have content and are in a section)
		if line != "" && currentSection != "" {
			if currentSection == "up" {
				upSQL.WriteString(line)
				upSQL.WriteString("\n")
			} else if currentSection == "down" {
				downSQL.WriteString(line)
				downSQL.WriteString("\n")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return MigrationFile{}, err
	}

	migration.UpSQL = strings.TrimSpace(upSQL.String())
	migration.DownSQL = strings.TrimSpace(downSQL.String())

	if migration.Version == "" {
		return MigrationFile{}, fmt.Errorf("migration file must contain version comment")
	}

	return migration, nil
}

// checkDatabaseExists checks if the database exists
func checkDatabaseExists(cfg *config.Config) (bool, error) {
	// Connect to server without database
	db, err := sql.Open("mysql", cfg.GetServerDSN())
	if err != nil {
		return false, err
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return false, err
	}

	// Check if database exists
	query := "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?"
	var dbName string
	err = db.QueryRowContext(context.Background(), query, cfg.Database.DBName).Scan(&dbName)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// createDatabase creates the database
func createDatabase(cfg *config.Config) error {
	// Connect to server without database
	db, err := sql.Open("mysql", cfg.GetServerDSN())
	if err != nil {
		return err
	}
	defer db.Close()

	// Create database
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.Database.DBName)
	_, err = db.ExecContext(context.Background(), query)
	return err
}

// promptCreateDatabase prompts user to create database
func promptCreateDatabase(dbName string) bool {
	fmt.Printf("Database '%s' does not exist. Do you want to create it? (y/N): ", dbName)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading input: %v", err)
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
