package repositories

import (
	"database/sql"
	"fmt"

	"github.com/gatehide/gatehide-api/internal/models"
)

// GamenetRepository defines the interface for gamenet data operations
type GamenetRepository interface {
	GetAll() ([]models.Gamenet, error)
	GetByID(id int) (*models.Gamenet, error)
	Create(gamenet *models.Gamenet) error
	Update(id int, gamenet *models.GamenetUpdateRequest) error
	Delete(id int) error
	Search(req *models.GamenetSearchRequest) (*models.GamenetSearchResponse, error)
}

// gamenetRepository implements GamenetRepository interface
type gamenetRepository struct {
	db *sql.DB
}

// NewGamenetRepository creates a new gamenet repository
func NewGamenetRepository(db *sql.DB) GamenetRepository {
	return &gamenetRepository{db: db}
}

// GetAll retrieves all gamenets
func (r *gamenetRepository) GetAll() ([]models.Gamenet, error) {
	query := `
		SELECT id, name, owner_name, owner_mobile, address, email, password, license_attachment, 
		       created_at, updated_at
		FROM gamenets 
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query gamenets: %w", err)
	}
	defer rows.Close()

	var gamenets []models.Gamenet
	for rows.Next() {
		var gamenet models.Gamenet
		err := rows.Scan(
			&gamenet.ID,
			&gamenet.Name,
			&gamenet.OwnerName,
			&gamenet.OwnerMobile,
			&gamenet.Address,
			&gamenet.Email,
			&gamenet.Password,
			&gamenet.LicenseAttachment,
			&gamenet.CreatedAt,
			&gamenet.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan gamenet: %w", err)
		}
		gamenets = append(gamenets, gamenet)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating gamenets: %w", err)
	}

	return gamenets, nil
}

// GetByID retrieves a gamenet by ID
func (r *gamenetRepository) GetByID(id int) (*models.Gamenet, error) {
	query := `
		SELECT id, name, owner_name, owner_mobile, address, email, password, license_attachment, 
		       created_at, updated_at
		FROM gamenets 
		WHERE id = ?
	`

	var gamenet models.Gamenet
	err := r.db.QueryRow(query, id).Scan(
		&gamenet.ID,
		&gamenet.Name,
		&gamenet.OwnerName,
		&gamenet.OwnerMobile,
		&gamenet.Address,
		&gamenet.Email,
		&gamenet.Password,
		&gamenet.LicenseAttachment,
		&gamenet.CreatedAt,
		&gamenet.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("gamenet not found")
		}
		return nil, fmt.Errorf("failed to get gamenet: %w", err)
	}

	return &gamenet, nil
}

// Create creates a new gamenet
func (r *gamenetRepository) Create(gamenet *models.Gamenet) error {
	query := `
		INSERT INTO gamenets (name, owner_name, owner_mobile, address, email, password, license_attachment)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		gamenet.Name,
		gamenet.OwnerName,
		gamenet.OwnerMobile,
		gamenet.Address,
		gamenet.Email,
		gamenet.Password,
		gamenet.LicenseAttachment,
	)

	if err != nil {
		return fmt.Errorf("failed to create gamenet: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	gamenet.ID = int(id)
	return nil
}

// Update updates an existing gamenet
func (r *gamenetRepository) Update(id int, updateData *models.GamenetUpdateRequest) error {
	// Build dynamic query based on provided fields
	query := "UPDATE gamenets SET "
	args := []interface{}{}
	fields := []string{}

	if updateData.Name != nil {
		fields = append(fields, "name = ?")
		args = append(args, *updateData.Name)
	}
	if updateData.OwnerName != nil {
		fields = append(fields, "owner_name = ?")
		args = append(args, *updateData.OwnerName)
	}
	if updateData.OwnerMobile != nil {
		fields = append(fields, "owner_mobile = ?")
		args = append(args, *updateData.OwnerMobile)
	}
	if updateData.Address != nil {
		fields = append(fields, "address = ?")
		args = append(args, *updateData.Address)
	}
	if updateData.Email != nil {
		fields = append(fields, "email = ?")
		args = append(args, *updateData.Email)
	}
	if updateData.Password != nil {
		fields = append(fields, "password = ?")
		args = append(args, *updateData.Password)
	}
	if updateData.LicenseAttachment != nil {
		fields = append(fields, "license_attachment = ?")
		args = append(args, *updateData.LicenseAttachment)
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
		return fmt.Errorf("failed to update gamenet: %w", err)
	}

	return nil
}

// Delete deletes a gamenet by ID
func (r *gamenetRepository) Delete(id int) error {
	query := "DELETE FROM gamenets WHERE id = ?"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete gamenet: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("gamenet not found")
	}

	return nil
}

// Search searches gamenets with pagination
func (r *gamenetRepository) Search(req *models.GamenetSearchRequest) (*models.GamenetSearchResponse, error) {
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
		whereClause = `WHERE name LIKE ? OR owner_name LIKE ? OR owner_mobile LIKE ? OR address LIKE ? OR email LIKE ?`
		searchTerm := "%" + req.Query + "%"
		args = []interface{}{searchTerm, searchTerm, searchTerm, searchTerm, searchTerm}
	}

	// Count total items
	countQuery := `SELECT COUNT(*) FROM gamenets ` + whereClause
	var totalItems int64
	err := r.db.QueryRow(countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("failed to count gamenets: %w", err)
	}

	// Calculate pagination info
	totalPages := int((totalItems + int64(req.PageSize) - 1) / int64(req.PageSize))
	hasNext := req.Page < totalPages
	hasPrev := req.Page > 1

	// Build data query
	dataQuery := `
		SELECT id, name, owner_name, owner_mobile, address, email, password, license_attachment, 
		       created_at, updated_at
		FROM gamenets 
		` + whereClause + `
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	// Add limit and offset to args
	args = append(args, req.PageSize, offset)

	rows, err := r.db.Query(dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query gamenets: %w", err)
	}
	defer rows.Close()

	var gamenets []models.Gamenet
	for rows.Next() {
		var gamenet models.Gamenet
		err := rows.Scan(
			&gamenet.ID,
			&gamenet.Name,
			&gamenet.OwnerName,
			&gamenet.OwnerMobile,
			&gamenet.Address,
			&gamenet.Email,
			&gamenet.Password,
			&gamenet.LicenseAttachment,
			&gamenet.CreatedAt,
			&gamenet.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan gamenet: %w", err)
		}
		gamenets = append(gamenets, gamenet)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating gamenets: %w", err)
	}

	// Convert to response format
	var responses []models.GamenetResponse
	for _, gamenet := range gamenets {
		responses = append(responses, gamenet.ToResponse())
	}

	return &models.GamenetSearchResponse{
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
