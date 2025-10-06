package services

import (
	"context"

	"github.com/gatehide/gatehide-api/internal/models"
)

// NotificationServiceInterface defines the contract for notification services
type NotificationServiceInterface interface {
	// SendNotification sends a notification of any type
	SendNotification(ctx context.Context, notification *models.CreateNotificationRequest) error

	// SendEmail sends an email notification
	SendEmail(ctx context.Context, email *models.SendEmailRequest) error

	// SendSMS sends an SMS notification
	SendSMS(ctx context.Context, sms *models.SendSMSRequest) error

	// SendDatabaseNotification sends a database notification
	SendDatabaseNotification(ctx context.Context, dbNotification *models.DatabaseNotification) error

	// GetNotification retrieves a notification by ID
	GetNotification(ctx context.Context, id int) (*models.Notification, error)

	// GetNotifications retrieves notifications with filtering
	GetNotifications(ctx context.Context, filters map[string]interface{}) ([]*models.Notification, error)

	// UpdateNotificationStatus updates the status of a notification
	UpdateNotificationStatus(ctx context.Context, id int, status models.NotificationStatus, errorMsg *string) error

	// RetryFailedNotification retries a failed notification
	RetryFailedNotification(ctx context.Context, id int) error
}

// EmailServiceInterface defines the contract for email services
type EmailServiceInterface interface {
	// SendEmail sends an email using SMTP
	SendEmail(ctx context.Context, email *models.EmailNotification) error

	// SendBulkEmail sends multiple emails
	SendBulkEmail(ctx context.Context, emails []*models.EmailNotification) error

	// ValidateEmailAddress validates an email address
	ValidateEmailAddress(email string) bool

	// TestConnection tests the SMTP connection
	TestConnection(ctx context.Context) error
}

// SMSServiceInterface defines the contract for SMS services
type SMSServiceInterface interface {
	// SendSMS sends an SMS message
	SendSMS(ctx context.Context, sms *models.SMSNotification) error

	// SendBulkSMS sends multiple SMS messages
	SendBulkSMS(ctx context.Context, smsMessages []*models.SMSNotification) error

	// ValidatePhoneNumber validates a phone number
	ValidatePhoneNumber(phone string) bool

	// TestConnection tests the SMS service connection
	TestConnection(ctx context.Context) error
}

// DatabaseNotificationServiceInterface defines the contract for database notification services
type DatabaseNotificationServiceInterface interface {
	// CreateNotification creates a database notification
	CreateNotification(ctx context.Context, notification *models.DatabaseNotification) error

	// GetUserNotifications retrieves notifications for a specific user
	GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]*models.DatabaseNotification, error)

	// MarkAsRead marks a notification as read
	MarkAsRead(ctx context.Context, notificationID int, userID int) error

	// MarkAllAsRead marks all notifications as read for a user
	MarkAllAsRead(ctx context.Context, userID int) error

	// DeleteNotification deletes a notification
	DeleteNotification(ctx context.Context, notificationID int, userID int) error
}

// TemplateServiceInterface defines the contract for template services
type TemplateServiceInterface interface {
	// GetTemplate retrieves a template by ID
	GetTemplate(ctx context.Context, id int) (*models.NotificationTemplate, error)

	// GetTemplateByName retrieves a template by name
	GetTemplateByName(ctx context.Context, name string, templateType models.NotificationType) (*models.NotificationTemplate, error)

	// CreateTemplate creates a new template
	CreateTemplate(ctx context.Context, template *models.NotificationTemplate) error

	// UpdateTemplate updates an existing template
	UpdateTemplate(ctx context.Context, template *models.NotificationTemplate) error

	// DeleteTemplate deletes a template
	DeleteTemplate(ctx context.Context, id int) error

	// ListTemplates lists all templates with optional filtering
	ListTemplates(ctx context.Context, templateType *models.NotificationType) ([]*models.NotificationTemplate, error)

	// RenderTemplate renders a template with provided data
	RenderTemplate(ctx context.Context, template *models.NotificationTemplate, data map[string]interface{}) (string, string, error) // returns subject, content, error
}
