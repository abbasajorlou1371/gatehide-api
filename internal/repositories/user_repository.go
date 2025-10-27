package repositories

import (
	"database/sql"
	"fmt"

	"github.com/gatehide/gatehide-api/internal/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	GetAll() ([]models.User, error)
	GetAllByGamenet(gamenetID int) ([]models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByMobile(mobile string) (*models.User, error)
	GetByID(id int) (*models.User, error)
	Create(user *models.User) error
	Update(id int, user *models.UserUpdateRequest) error
	Delete(id int) error
	Search(req *models.UserSearchRequest) (*models.UserSearchResponse, error)
	SearchByGamenet(req *models.UserSearchRequest, gamenetID int) (*models.UserSearchResponse, error)
	UpdateLastLogin(id int) error
	UpdatePassword(id int, hashedPassword string) error
	UpdateProfile(id int, name, mobile, image string) error
	UpdateEmail(id int, email string) error
	LinkToGamenet(userID, gamenetID int) error
	UnlinkFromGamenet(userID, gamenetID int) error
	GetGamenetIDByUser(userID int) (*int, error)
}

// AdminRepository defines the interface for admin data operations
type AdminRepository interface {
	GetByEmail(email string) (*models.Admin, error)
	GetByID(id int) (*models.Admin, error)
	UpdateLastLogin(id int) error
	UpdatePassword(id int, hashedPassword string) error
	UpdateProfile(id int, name, mobile, image string) error
	UpdateEmail(id int, email string) error
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

// GetAll retrieves all users
func (r *userRepository) GetAll() ([]models.User, error) {
	query := `
		SELECT id, name, mobile, email, password, image, balance, debt, last_login_at, created_at, updated_at
		FROM users 
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Mobile,
			&user.Email,
			&user.Password,
			&user.Image,
			&user.Balance,
			&user.Debt,
			&user.LastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

// GetAllByGamenet retrieves all users for a specific gamenet
func (r *userRepository) GetAllByGamenet(gamenetID int) ([]models.User, error) {
	query := `
		SELECT u.id, u.name, u.mobile, u.email, u.password, u.image, u.balance, u.debt, u.last_login_at, u.created_at, u.updated_at
		FROM users u
		INNER JOIN users_gamenets ug ON u.id = ug.user_id
		WHERE ug.gamenet_id = ?
		ORDER BY u.created_at DESC
	`

	rows, err := r.db.Query(query, gamenetID)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Mobile,
			&user.Email,
			&user.Password,
			&user.Image,
			&user.Balance,
			&user.Debt,
			&user.LastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, name, mobile, email, password, image, balance, debt, last_login_at, created_at, updated_at
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
		&user.Balance,
		&user.Debt,
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

// GetByMobile retrieves a user by mobile number
func (r *userRepository) GetByMobile(mobile string) (*models.User, error) {
	query := `
		SELECT id, name, mobile, email, password, image, balance, debt, last_login_at, created_at, updated_at
		FROM users 
		WHERE mobile = ?
	`

	user := &models.User{}
	err := r.db.QueryRow(query, mobile).Scan(
		&user.ID,
		&user.Name,
		&user.Mobile,
		&user.Email,
		&user.Password,
		&user.Image,
		&user.Balance,
		&user.Debt,
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
		SELECT id, name, mobile, email, password, image, balance, debt, last_login_at, created_at, updated_at
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
		&user.Balance,
		&user.Debt,
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

// Create creates a new user
func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (name, mobile, email, password, image)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		user.Name,
		user.Mobile,
		user.Email,
		user.Password,
		user.Image,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	user.ID = int(id)
	return nil
}

// Update updates an existing user
func (r *userRepository) Update(id int, updateData *models.UserUpdateRequest) error {
	// Build dynamic query based on provided fields
	query := "UPDATE users SET "
	args := []interface{}{}
	fields := []string{}

	if updateData.Name != nil {
		fields = append(fields, "name = ?")
		args = append(args, *updateData.Name)
	}
	if updateData.Email != nil {
		fields = append(fields, "email = ?")
		args = append(args, *updateData.Email)
	}
	if updateData.Mobile != nil {
		fields = append(fields, "mobile = ?")
		args = append(args, *updateData.Mobile)
	}
	if updateData.Image != nil {
		fields = append(fields, "image = ?")
		args = append(args, *updateData.Image)
	}

	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query += fmt.Sprintf("%s", fields[0])
	for i := 1; i < len(fields); i++ {
		query += fmt.Sprintf(", %s", fields[i])
	}
	query += ", updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete deletes a user by ID
func (r *userRepository) Delete(id int) error {
	query := "DELETE FROM users WHERE id = ?"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Search searches users with pagination
func (r *userRepository) Search(req *models.UserSearchRequest) (*models.UserSearchResponse, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	offset := (req.Page - 1) * req.PageSize

	// Build search query
	var whereClause string
	var args []interface{}

	if req.Query != "" {
		whereClause = `WHERE name LIKE ? OR mobile LIKE ? OR email LIKE ?`
		searchTerm := "%" + req.Query + "%"
		args = []interface{}{searchTerm, searchTerm, searchTerm}
	}

	// Count total items
	countQuery := `SELECT COUNT(*) FROM users ` + whereClause
	var totalItems int64
	err := r.db.QueryRow(countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Calculate pagination info
	totalPages := int((totalItems + int64(req.PageSize) - 1) / int64(req.PageSize))
	hasNext := req.Page < totalPages
	hasPrev := req.Page > 1

	// Build data query
	dataQuery := `
		SELECT id, name, mobile, email, password, image, balance, debt, last_login_at, created_at, updated_at
		FROM users 
		` + whereClause + `
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	// Add limit and offset to args
	args = append(args, req.PageSize, offset)

	rows, err := r.db.Query(dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Mobile,
			&user.Email,
			&user.Password,
			&user.Image,
			&user.Balance,
			&user.Debt,
			&user.LastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	// Convert to response format
	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	return &models.UserSearchResponse{
		Data: responses,
		Pagination: models.PaginationInfo{
			CurrentPage: req.Page,
			PageSize:    req.PageSize,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			HasNext:     hasNext,
			HasPrev:     hasPrev,
		},
	}, nil
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

// UpdateProfile updates a user's profile information
func (r *userRepository) UpdateProfile(id int, name, mobile, image string) error {
	query := `UPDATE users SET name = ?, mobile = ?, image = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.Exec(query, name, mobile, image, id)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

// UpdateEmail updates a user's email
func (r *userRepository) UpdateEmail(id int, email string) error {
	query := `UPDATE users SET email = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.Exec(query, email, id)
	if err != nil {
		return fmt.Errorf("failed to update email: %w", err)
	}

	return nil
}

// SearchByGamenet searches users for a specific gamenet with pagination
func (r *userRepository) SearchByGamenet(req *models.UserSearchRequest, gamenetID int) (*models.UserSearchResponse, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	offset := (req.Page - 1) * req.PageSize

	// Build search query with gamenet join
	var whereClause string
	var args []interface{}

	baseWhere := "WHERE ug.gamenet_id = ?"
	args = append(args, gamenetID)

	if req.Query != "" {
		whereClause = baseWhere + ` AND (u.name LIKE ? OR u.mobile LIKE ? OR u.email LIKE ?)`
		searchTerm := "%" + req.Query + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	} else {
		whereClause = baseWhere
	}

	// Count total items
	countQuery := `SELECT COUNT(*) FROM users u INNER JOIN users_gamenets ug ON u.id = ug.user_id ` + whereClause
	var totalItems int64
	err := r.db.QueryRow(countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Calculate pagination info
	totalPages := int((totalItems + int64(req.PageSize) - 1) / int64(req.PageSize))
	hasNext := req.Page < totalPages
	hasPrev := req.Page > 1

	// Build data query
	dataQuery := `
		SELECT u.id, u.name, u.mobile, u.email, u.password, u.image, u.balance, u.debt, u.last_login_at, u.created_at, u.updated_at
		FROM users u
		INNER JOIN users_gamenets ug ON u.id = ug.user_id
		` + whereClause + `
		ORDER BY u.created_at DESC
		LIMIT ? OFFSET ?
	`

	// Add limit and offset to args
	args = append(args, req.PageSize, offset)

	rows, err := r.db.Query(dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Mobile,
			&user.Email,
			&user.Password,
			&user.Image,
			&user.Balance,
			&user.Debt,
			&user.LastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	// Convert to response format
	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	return &models.UserSearchResponse{
		Data: responses,
		Pagination: models.PaginationInfo{
			CurrentPage: req.Page,
			PageSize:    req.PageSize,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			HasNext:     hasNext,
			HasPrev:     hasPrev,
		},
	}, nil
}

// LinkToGamenet links a user to a gamenet
func (r *userRepository) LinkToGamenet(userID, gamenetID int) error {
	query := `INSERT INTO users_gamenets (user_id, gamenet_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP`

	_, err := r.db.Exec(query, userID, gamenetID)
	if err != nil {
		return fmt.Errorf("failed to link user to gamenet: %w", err)
	}

	return nil
}

// UnlinkFromGamenet unlinks a user from a gamenet
func (r *userRepository) UnlinkFromGamenet(userID, gamenetID int) error {
	query := `DELETE FROM users_gamenets WHERE user_id = ? AND gamenet_id = ?`

	result, err := r.db.Exec(query, userID, gamenetID)
	if err != nil {
		return fmt.Errorf("failed to unlink user from gamenet: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user is not linked to this gamenet")
	}

	return nil
}

// GetGamenetIDByUser gets the gamenet ID that created a user (first linked gamenet)
func (r *userRepository) GetGamenetIDByUser(userID int) (*int, error) {
	query := `SELECT gamenet_id FROM users_gamenets WHERE user_id = ? ORDER BY created_at ASC LIMIT 1`

	var gamenetID int
	err := r.db.QueryRow(query, userID).Scan(&gamenetID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get gamenet ID: %w", err)
	}

	return &gamenetID, nil
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

// UpdateProfile updates an admin's profile information
func (r *adminRepository) UpdateProfile(id int, name, mobile, image string) error {
	query := `UPDATE admins SET name = ?, mobile = ?, image = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.Exec(query, name, mobile, image, id)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

// UpdateEmail updates an admin's email
func (r *adminRepository) UpdateEmail(id int, email string) error {
	query := `UPDATE admins SET email = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.Exec(query, email, id)
	if err != nil {
		return fmt.Errorf("failed to update email: %w", err)
	}

	return nil
}
