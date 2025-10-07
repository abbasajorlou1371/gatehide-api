package repositories

import (
	"database/sql"
	"fmt"

	"github.com/gatehide/gatehide-api/internal/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	GetByEmail(email string) (*models.User, error)
	GetByID(id int) (*models.User, error)
	UpdateLastLogin(id int) error
	UpdatePassword(id int, hashedPassword string) error
}

// AdminRepository defines the interface for admin data operations
type AdminRepository interface {
	GetByEmail(email string) (*models.Admin, error)
	GetByID(id int) (*models.Admin, error)
	UpdateLastLogin(id int) error
	UpdatePassword(id int, hashedPassword string) error
}

// userRepository implements UserRepository interface
type userRepository struct {
	db *sql.DB
}

// adminRepository implements AdminRepository interface
type adminRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// NewAdminRepository creates a new admin repository
func NewAdminRepository(db *sql.DB) AdminRepository {
	return &adminRepository{db: db}
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, name, mobile, email, password, image, last_login_at, created_at, updated_at
		FROM users 
		WHERE email = ?
	`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Mobile,
		&user.Email,
		&user.Password,
		&user.Image,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id int) (*models.User, error) {
	query := `
		SELECT id, name, mobile, email, password, image, last_login_at, created_at, updated_at
		FROM users 
		WHERE id = ?
	`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Mobile,
		&user.Email,
		&user.Password,
		&user.Image,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateLastLogin updates the last login timestamp for a user
func (r *userRepository) UpdateLastLogin(id int) error {
	query := `UPDATE users SET last_login_at = NOW() WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// UpdatePassword updates the password for a user
func (r *userRepository) UpdatePassword(id int, hashedPassword string) error {
	query := `UPDATE users SET password = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.Exec(query, hashedPassword, id)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// GetByEmail retrieves an admin by email
func (r *adminRepository) GetByEmail(email string) (*models.Admin, error) {
	query := `
		SELECT id, name, mobile, email, password, image, last_login_at, created_at, updated_at
		FROM admins 
		WHERE email = ?
	`

	admin := &models.Admin{}
	err := r.db.QueryRow(query, email).Scan(
		&admin.ID,
		&admin.Name,
		&admin.Mobile,
		&admin.Email,
		&admin.Password,
		&admin.Image,
		&admin.LastLoginAt,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("admin not found")
		}
		return nil, fmt.Errorf("failed to get admin: %w", err)
	}

	return admin, nil
}

// GetByID retrieves an admin by ID
func (r *adminRepository) GetByID(id int) (*models.Admin, error) {
	query := `
		SELECT id, name, mobile, email, password, image, last_login_at, created_at, updated_at
		FROM admins 
		WHERE id = ?
	`

	admin := &models.Admin{}
	err := r.db.QueryRow(query, id).Scan(
		&admin.ID,
		&admin.Name,
		&admin.Mobile,
		&admin.Email,
		&admin.Password,
		&admin.Image,
		&admin.LastLoginAt,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("admin not found")
		}
		return nil, fmt.Errorf("failed to get admin: %w", err)
	}

	return admin, nil
}

// UpdateLastLogin updates the last login timestamp for an admin
func (r *adminRepository) UpdateLastLogin(id int) error {
	query := `UPDATE admins SET last_login_at = NOW() WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// UpdatePassword updates the password for an admin
func (r *adminRepository) UpdatePassword(id int, hashedPassword string) error {
	query := `UPDATE admins SET password = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.Exec(query, hashedPassword, id)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
