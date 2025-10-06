package migrations

import (
	"time"
)

// Migration represents a database migration
type Migration struct {
	ID          int       `json:"id" db:"id"`
	Version     string    `json:"version" db:"version"`
	Description string    `json:"description" db:"description"`
	AppliedAt   time.Time `json:"applied_at" db:"applied_at"`
}

// MigrationFile represents a migration file structure
type MigrationFile struct {
	Version     string
	Description string
	UpSQL       string
	DownSQL     string
}

// MigrationRunner interface defines methods for running migrations
type MigrationRunner interface {
	CreateMigrationTable() error
	GetAppliedMigrations() ([]Migration, error)
	ApplyMigration(version, description, upSQL string) error
	RollbackMigration(version, downSQL string) error
	CheckDatabaseExists() (bool, error)
	CreateDatabase() error
	Close() error
}
