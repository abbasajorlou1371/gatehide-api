package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gatehide/gatehide-api/internal/models"
)

// NotificationRepository handles notification data operations
type NotificationRepository interface {
	Create(notification *models.Notification) error
	GetByID(id int) (*models.Notification, error)
	GetWithFilters(filters map[string]interface{}) ([]*models.Notification, error)
	Update(notification *models.Notification) error
	Delete(id int) error
	GetPendingNotifications(limit int) ([]*models.Notification, error)
	GetFailedNotifications(limit int) ([]*models.Notification, error)
}

// MySQLNotificationRepository implements NotificationRepository for MySQL
type MySQLNotificationRepository struct {
	db *sql.DB
}

// NewMySQLNotificationRepository creates a new MySQL notification repository
func NewMySQLNotificationRepository(db *sql.DB) *MySQLNotificationRepository {
	return &MySQLNotificationRepository{db: db}
}

// Create creates a new notification
func (r *MySQLNotificationRepository) Create(notification *models.Notification) error {
	query := `
		INSERT INTO notifications (
			type, status, priority, recipient, subject, content, 
			template_id, template_data, metadata, scheduled_at, 
			retry_count, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	templateDataJSON, _ := json.Marshal(notification.TemplateData)
	metadataJSON, _ := json.Marshal(notification.Metadata)

	result, err := r.db.Exec(
		query,
		notification.Type,
		notification.Status,
		notification.Priority,
		notification.Recipient,
		notification.Subject,
		notification.Content,
		notification.TemplateID,
		templateDataJSON,
		metadataJSON,
		notification.ScheduledAt,
		notification.RetryCount,
		notification.CreatedAt,
		notification.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	notification.ID = int(id)
	return nil
}

// GetByID retrieves a notification by ID
func (r *MySQLNotificationRepository) GetByID(id int) (*models.Notification, error) {
	query := `
		SELECT id, type, status, priority, recipient, subject, content,
			   template_id, template_data, metadata, scheduled_at, sent_at,
			   error_msg, retry_count, created_at, updated_at
		FROM notifications WHERE id = ?
	`

	var notification models.Notification
	var templateDataJSON, metadataJSON sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&notification.ID,
		&notification.Type,
		&notification.Status,
		&notification.Priority,
		&notification.Recipient,
		&notification.Subject,
		&notification.Content,
		&notification.TemplateID,
		&templateDataJSON,
		&metadataJSON,
		&notification.ScheduledAt,
		&notification.SentAt,
		&notification.ErrorMsg,
		&notification.RetryCount,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("notification not found")
		}
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	// Parse JSON fields
	if templateDataJSON.Valid && templateDataJSON.String != "" {
		if err := json.Unmarshal([]byte(templateDataJSON.String), &notification.TemplateData); err != nil {
			return nil, fmt.Errorf("failed to parse template data: %w", err)
		}
	}

	if metadataJSON.Valid && metadataJSON.String != "" {
		if err := json.Unmarshal([]byte(metadataJSON.String), &notification.Metadata); err != nil {
			return nil, fmt.Errorf("failed to parse metadata: %w", err)
		}
	}

	return &notification, nil
}

// GetWithFilters retrieves notifications with optional filters
func (r *MySQLNotificationRepository) GetWithFilters(filters map[string]interface{}) ([]*models.Notification, error) {
	query := `
		SELECT id, type, status, priority, recipient, subject, content,
			   template_id, template_data, metadata, scheduled_at, sent_at,
			   error_msg, retry_count, created_at, updated_at
		FROM notifications
	`

	args := []interface{}{}
	whereClauses := []string{}

	// Add filters
	if status, ok := filters["status"]; ok {
		whereClauses = append(whereClauses, "status = ?")
		args = append(args, status)
	}

	if notificationType, ok := filters["type"]; ok {
		whereClauses = append(whereClauses, "type = ?")
		args = append(args, notificationType)
	}

	if recipient, ok := filters["recipient"]; ok {
		whereClauses = append(whereClauses, "recipient = ?")
		args = append(args, recipient)
	}

	if priority, ok := filters["priority"]; ok {
		whereClauses = append(whereClauses, "priority = ?")
		args = append(args, priority)
	}

	// Add WHERE clause if filters exist
	if len(whereClauses) > 0 {
		query += " WHERE " + fmt.Sprintf("%s", whereClauses[0])
		for i := 1; i < len(whereClauses); i++ {
			query += " AND " + whereClauses[i]
		}
	}

	// Add ordering and limit
	query += " ORDER BY created_at DESC"

	if limit, ok := filters["limit"]; ok {
		if limitInt, ok := limit.(int); ok && limitInt > 0 {
			query += " LIMIT ?"
			args = append(args, limitInt)
		}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		var notification models.Notification
		var templateDataJSON, metadataJSON sql.NullString

		err := rows.Scan(
			&notification.ID,
			&notification.Type,
			&notification.Status,
			&notification.Priority,
			&notification.Recipient,
			&notification.Subject,
			&notification.Content,
			&notification.TemplateID,
			&templateDataJSON,
			&metadataJSON,
			&notification.ScheduledAt,
			&notification.SentAt,
			&notification.ErrorMsg,
			&notification.RetryCount,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}

		// Parse JSON fields
		if templateDataJSON.Valid && templateDataJSON.String != "" {
			if err := json.Unmarshal([]byte(templateDataJSON.String), &notification.TemplateData); err != nil {
				return nil, fmt.Errorf("failed to parse template data: %w", err)
			}
		}

		if metadataJSON.Valid && metadataJSON.String != "" {
			if err := json.Unmarshal([]byte(metadataJSON.String), &notification.Metadata); err != nil {
				return nil, fmt.Errorf("failed to parse metadata: %w", err)
			}
		}

		notifications = append(notifications, &notification)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notifications: %w", err)
	}

	return notifications, nil
}

// Update updates an existing notification
func (r *MySQLNotificationRepository) Update(notification *models.Notification) error {
	query := `
		UPDATE notifications SET
			type = ?, status = ?, priority = ?, recipient = ?, subject = ?, content = ?,
			template_id = ?, template_data = ?, metadata = ?, scheduled_at = ?, sent_at = ?,
			error_msg = ?, retry_count = ?, updated_at = ?
		WHERE id = ?
	`

	templateDataJSON, _ := json.Marshal(notification.TemplateData)
	metadataJSON, _ := json.Marshal(notification.Metadata)

	_, err := r.db.Exec(
		query,
		notification.Type,
		notification.Status,
		notification.Priority,
		notification.Recipient,
		notification.Subject,
		notification.Content,
		notification.TemplateID,
		templateDataJSON,
		metadataJSON,
		notification.ScheduledAt,
		notification.SentAt,
		notification.ErrorMsg,
		notification.RetryCount,
		notification.UpdatedAt,
		notification.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	return nil
}

// Delete deletes a notification
func (r *MySQLNotificationRepository) Delete(id int) error {
	query := "DELETE FROM notifications WHERE id = ?"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

// GetPendingNotifications retrieves pending notifications
func (r *MySQLNotificationRepository) GetPendingNotifications(limit int) ([]*models.Notification, error) {
	filters := map[string]interface{}{
		"status": models.NotificationStatusPending,
		"limit":  limit,
	}
	return r.GetWithFilters(filters)
}

// GetFailedNotifications retrieves failed notifications
func (r *MySQLNotificationRepository) GetFailedNotifications(limit int) ([]*models.Notification, error) {
	filters := map[string]interface{}{
		"status": models.NotificationStatusFailed,
		"limit":  limit,
	}
	return r.GetWithFilters(filters)
}
