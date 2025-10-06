package models

import (
	"time"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeEmail    NotificationType = "email"
	NotificationTypeSMS      NotificationType = "sms"
	NotificationTypeDatabase NotificationType = "database"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	NotificationStatusPending   NotificationStatus = "pending"
	NotificationStatusSent      NotificationStatus = "sent"
	NotificationStatusFailed    NotificationStatus = "failed"
	NotificationStatusCancelled NotificationStatus = "cancelled"
)

// NotificationPriority represents the priority of a notification
type NotificationPriority string

const (
	NotificationPriorityLow    NotificationPriority = "low"
	NotificationPriorityNormal NotificationPriority = "normal"
	NotificationPriorityHigh   NotificationPriority = "high"
	NotificationPriorityUrgent NotificationPriority = "urgent"
)

// Notification represents a notification in the system
type Notification struct {
	ID           int                    `json:"id" db:"id"`
	Type         NotificationType       `json:"type" db:"type"`
	Status       NotificationStatus     `json:"status" db:"status"`
	Priority     NotificationPriority   `json:"priority" db:"priority"`
	Recipient    string                 `json:"recipient" db:"recipient"`
	Subject      string                 `json:"subject" db:"subject"`
	Content      string                 `json:"content" db:"content"`
	TemplateID   *int                   `json:"template_id" db:"template_id"`
	TemplateData map[string]interface{} `json:"template_data" db:"template_data"`
	Metadata     map[string]interface{} `json:"metadata" db:"metadata"`
	ScheduledAt  *time.Time             `json:"scheduled_at" db:"scheduled_at"`
	SentAt       *time.Time             `json:"sent_at" db:"sent_at"`
	ErrorMsg     *string                `json:"error_msg" db:"error_msg"`
	RetryCount   int                    `json:"retry_count" db:"retry_count"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at"`
}

// EmailNotification represents an email notification
type EmailNotification struct {
	To          []string             `json:"to"`
	CC          []string             `json:"cc,omitempty"`
	BCC         []string             `json:"bcc,omitempty"`
	Subject     string               `json:"subject"`
	Body        string               `json:"body"`
	HTMLBody    string               `json:"html_body,omitempty"`
	Attachments []string             `json:"attachments,omitempty"`
	Priority    NotificationPriority `json:"priority,omitempty"`
}

// SMSNotification represents an SMS notification
type SMSNotification struct {
	To       string               `json:"to"`
	Message  string               `json:"message"`
	Priority NotificationPriority `json:"priority,omitempty"`
}

// DatabaseNotification represents a database notification
type DatabaseNotification struct {
	UserID   int                  `json:"user_id"`
	Title    string               `json:"title"`
	Message  string               `json:"message"`
	Type     string               `json:"type"`
	Priority NotificationPriority `json:"priority,omitempty"`
}

// NotificationTemplate represents an email template
type NotificationTemplate struct {
	ID          int              `json:"id" db:"id"`
	Name        string           `json:"name" db:"name"`
	Type        NotificationType `json:"type" db:"type"`
	Subject     string           `json:"subject" db:"subject"`
	Content     string           `json:"content" db:"content"`
	HTMLContent string           `json:"html_content" db:"html_content"`
	Variables   []string         `json:"variables" db:"variables"`
	IsActive    bool             `json:"is_active" db:"is_active"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
}

// CreateNotificationRequest represents a request to create a notification
type CreateNotificationRequest struct {
	Type         NotificationType       `json:"type" binding:"required"`
	Priority     NotificationPriority   `json:"priority,omitempty"`
	Recipient    string                 `json:"recipient" binding:"required"`
	Subject      string                 `json:"subject,omitempty"`
	Content      string                 `json:"content,omitempty"`
	TemplateID   *int                   `json:"template_id,omitempty"`
	TemplateData map[string]interface{} `json:"template_data,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	ScheduledAt  *time.Time             `json:"scheduled_at,omitempty"`
}

// SendEmailRequest represents a request to send an email
type SendEmailRequest struct {
	To          []string             `json:"to" binding:"required"`
	CC          []string             `json:"cc,omitempty"`
	BCC         []string             `json:"bcc,omitempty"`
	Subject     string               `json:"subject" binding:"required"`
	Body        string               `json:"body,omitempty"`
	HTMLBody    string               `json:"html_body,omitempty"`
	Attachments []string             `json:"attachments,omitempty"`
	Priority    NotificationPriority `json:"priority,omitempty"`
}

// SendSMSRequest represents a request to send an SMS
type SendSMSRequest struct {
	To       string               `json:"to" binding:"required"`
	Message  string               `json:"message" binding:"required"`
	Priority NotificationPriority `json:"priority,omitempty"`
}

// NotificationResponse represents a notification response
type NotificationResponse struct {
	ID          int                  `json:"id"`
	Type        NotificationType     `json:"type"`
	Status      NotificationStatus   `json:"status"`
	Priority    NotificationPriority `json:"priority"`
	Recipient   string               `json:"recipient"`
	Subject     string               `json:"subject"`
	Content     string               `json:"content"`
	ScheduledAt *time.Time           `json:"scheduled_at"`
	SentAt      *time.Time           `json:"sent_at"`
	ErrorMsg    *string              `json:"error_msg"`
	RetryCount  int                  `json:"retry_count"`
	CreatedAt   time.Time            `json:"created_at"`
}

// ToResponse converts Notification to NotificationResponse
func (n *Notification) ToResponse() NotificationResponse {
	return NotificationResponse{
		ID:          n.ID,
		Type:        n.Type,
		Status:      n.Status,
		Priority:    n.Priority,
		Recipient:   n.Recipient,
		Subject:     n.Subject,
		Content:     n.Content,
		ScheduledAt: n.ScheduledAt,
		SentAt:      n.SentAt,
		ErrorMsg:    n.ErrorMsg,
		RetryCount:  n.RetryCount,
		CreatedAt:   n.CreatedAt,
	}
}
