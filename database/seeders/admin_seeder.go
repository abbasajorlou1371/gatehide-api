package seeders

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gatehide/gatehide-api/config"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// init registers the admin seeder
func init() {
	RegisterSeeder("admin", SeedAdmin)
}

// AdminSeeder handles seeding admin data
type AdminSeeder struct {
	db *sql.DB
}

// NewAdminSeeder creates a new admin seeder instance
func NewAdminSeeder(cfg *config.Config) (*AdminSeeder, error) {
	db, err := sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &AdminSeeder{db: db}, nil
}

// AdminData represents admin user data
type AdminData struct {
	Name     string
	Mobile   string
	Email    string
	Password string
}

// SeedAdmin is the public seeder function that can be called by the registry
func SeedAdmin(cfg *config.Config) error {
	seeder, err := NewAdminSeeder(cfg)
	if err != nil {
		return fmt.Errorf("failed to create admin seeder: %w", err)
	}
	defer seeder.Close()

	// Admin data to seed
	admin := AdminData{
		Name:     "Abbas Ajorloo",
		Mobile:   "01907529334",
		Email:    "abbas.ajorlou1371@gmail.com",
		Password: "46769732",
	}

	return seeder.seedAdmin(admin)
}

// seedAdmin seeds an admin user into the database (private method)
func (s *AdminSeeder) seedAdmin(admin AdminData) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Check if admin already exists
	var count int
	checkQuery := "SELECT COUNT(*) FROM admins WHERE email = ? OR mobile = ?"
	err = s.db.QueryRow(checkQuery, admin.Email, admin.Mobile).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing admin: %w", err)
	}

	if count > 0 {
		log.Printf("Admin with email %s or mobile %s already exists, skipping...", admin.Email, admin.Mobile)
		return nil
	}

	// Insert admin
	insertQuery := `
		INSERT INTO admins (name, mobile, email, password, created_at, updated_at) 
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	result, err := s.db.Exec(insertQuery, admin.Name, admin.Mobile, admin.Email, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to insert admin: %w", err)
	}

	adminID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	log.Printf("✅ Admin seeded successfully with ID: %d", adminID)
	return nil
}

// Close closes the database connection
func (s *AdminSeeder) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
