package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gatehide/gatehide-api/internal/models"
)

// NotificationTemplateRepository handles notification template data operations
type NotificationTemplateRepository interface {
	Create(template *models.NotificationTemplate) error
	GetByID(id int) (*models.NotificationTemplate, error)
	GetByNameAndType(name string, templateType models.NotificationType) (*models.NotificationTemplate, error)
	GetAll() ([]*models.NotificationTemplate, error)
	GetByType(templateType models.NotificationType) ([]*models.NotificationTemplate, error)
	Update(template *models.NotificationTemplate) error
	Delete(id int) error
}

// MySQLTemplateRepository implements NotificationTemplateRepository for MySQL
type MySQLTemplateRepository struct {
	db *sql.DB
}

// NewMySQLTemplateRepository creates a new MySQL template repository
func NewMySQLTemplateRepository(db *sql.DB) *MySQLTemplateRepository {
	return &MySQLTemplateRepository{db: db}
}

// Create creates a new notification template
func (r *MySQLTemplateRepository) Create(template *models.NotificationTemplate) error {
	query := `
		INSERT INTO notification_templates (
			name, type, subject, content, html_content, variables, is_active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	variablesJSON, _ := json.Marshal(template.Variables)

	result, err := r.db.Exec(
		query,
		template.Name,
		template.Type,
		template.Subject,
		template.Content,
		template.HTMLContent,
		variablesJSON,
		template.IsActive,
		template.CreatedAt,
		template.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	template.ID = int(id)
	return nil
}

// GetByID retrieves a template by ID
func (r *MySQLTemplateRepository) GetByID(id int) (*models.NotificationTemplate, error) {
	query := `
		SELECT id, name, type, subject, content, html_content, variables, is_active, created_at, updated_at
		FROM notification_templates WHERE id = ?
	`

	var template models.NotificationTemplate
	var variablesJSON string

	err := r.db.QueryRow(query, id).Scan(
		&template.ID,
		&template.Name,
		&template.Type,
		&template.Subject,
		&template.Content,
		&template.HTMLContent,
		&variablesJSON,
		&template.IsActive,
		&template.CreatedAt,
		&template.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("template not found")
		}
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	// Parse variables JSON
	if err := json.Unmarshal([]byte(variablesJSON), &template.Variables); err != nil {
		return nil, fmt.Errorf("failed to parse variables: %w", err)
	}

	return &template, nil
}

// GetByNameAndType retrieves a template by name and type
func (r *MySQLTemplateRepository) GetByNameAndType(name string, templateType models.NotificationType) (*models.NotificationTemplate, error) {
	query := `
		SELECT id, name, type, subject, content, html_content, variables, is_active, created_at, updated_at
		FROM notification_templates WHERE name = ? AND type = ? AND is_active = 1
	`

	var template models.NotificationTemplate
	var variablesJSON string

	err := r.db.QueryRow(query, name, templateType).Scan(
		&template.ID,
		&template.Name,
		&template.Type,
		&template.Subject,
		&template.Content,
		&template.HTMLContent,
		&variablesJSON,
		&template.IsActive,
		&template.CreatedAt,
		&template.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("template not found")
		}
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	// Parse variables JSON
	if err := json.Unmarshal([]byte(variablesJSON), &template.Variables); err != nil {
		return nil, fmt.Errorf("failed to parse variables: %w", err)
	}

	return &template, nil
}

// GetAll retrieves all templates
func (r *MySQLTemplateRepository) GetAll() ([]*models.NotificationTemplate, error) {
	query := `
		SELECT id, name, type, subject, content, html_content, variables, is_active, created_at, updated_at
		FROM notification_templates ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query templates: %w", err)
	}
	defer rows.Close()

	var templates []*models.NotificationTemplate
	for rows.Next() {
		var template models.NotificationTemplate
		var variablesJSON string

		err := rows.Scan(
			&template.ID,
			&template.Name,
			&template.Type,
			&template.Subject,
			&template.Content,
			&template.HTMLContent,
			&variablesJSON,
			&template.IsActive,
			&template.CreatedAt,
			&template.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}

		// Parse variables JSON
		if err := json.Unmarshal([]byte(variablesJSON), &template.Variables); err != nil {
			return nil, fmt.Errorf("failed to parse variables: %w", err)
		}

		templates = append(templates, &template)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating templates: %w", err)
	}

	return templates, nil
}

// GetByType retrieves templates by type
func (r *MySQLTemplateRepository) GetByType(templateType models.NotificationType) ([]*models.NotificationTemplate, error) {
	query := `
		SELECT id, name, type, subject, content, html_content, variables, is_active, created_at, updated_at
		FROM notification_templates WHERE type = ? ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, templateType)
	if err != nil {
		return nil, fmt.Errorf("failed to query templates: %w", err)
	}
	defer rows.Close()

	var templates []*models.NotificationTemplate
	for rows.Next() {
		var template models.NotificationTemplate
		var variablesJSON string

		err := rows.Scan(
			&template.ID,
			&template.Name,
			&template.Type,
			&template.Subject,
			&template.Content,
			&template.HTMLContent,
			&variablesJSON,
			&template.IsActive,
			&template.CreatedAt,
			&template.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}

		// Parse variables JSON
		if err := json.Unmarshal([]byte(variablesJSON), &template.Variables); err != nil {
			return nil, fmt.Errorf("failed to parse variables: %w", err)
		}

		templates = append(templates, &template)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating templates: %w", err)
	}

	return templates, nil
}

// Update updates an existing template
func (r *MySQLTemplateRepository) Update(template *models.NotificationTemplate) error {
	query := `
		UPDATE notification_templates SET
			name = ?, type = ?, subject = ?, content = ?, html_content = ?,
			variables = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`

	variablesJSON, _ := json.Marshal(template.Variables)

	_, err := r.db.Exec(
		query,
		template.Name,
		template.Type,
		template.Subject,
		template.Content,
		template.HTMLContent,
		variablesJSON,
		template.IsActive,
		template.UpdatedAt,
		template.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}

	return nil
}

// Delete deletes a template
func (r *MySQLTemplateRepository) Delete(id int) error {
	query := "DELETE FROM notification_templates WHERE id = ?"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}
	return nil
}
