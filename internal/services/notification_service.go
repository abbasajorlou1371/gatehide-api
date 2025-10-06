package services

import (
	"context"
	"fmt"
	"time"

	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/repositories"
)

// NotificationService implements NotificationServiceInterface
type NotificationService struct {
	emailService          EmailServiceInterface
	smsService            SMSServiceInterface
	dbNotificationService DatabaseNotificationServiceInterface
	templateService       TemplateServiceInterface
	notificationRepo      repositories.NotificationRepository
	config                *config.Config
}

// NewNotificationService creates a new notification service instance
func NewNotificationService(
	emailService EmailServiceInterface,
	smsService SMSServiceInterface,
	dbNotificationService DatabaseNotificationServiceInterface,
	templateService TemplateServiceInterface,
	notificationRepo repositories.NotificationRepository,
	cfg *config.Config,
) *NotificationService {
	return &NotificationService{
		emailService:          emailService,
		smsService:            smsService,
		dbNotificationService: dbNotificationService,
		templateService:       templateService,
		notificationRepo:      notificationRepo,
		config:                cfg,
	}
}

// SendNotification sends a notification of any type
func (s *NotificationService) SendNotification(ctx context.Context, notification *models.CreateNotificationRequest) error {
	// Create notification record
	notificationRecord := &models.Notification{
		Type:         notification.Type,
		Status:       models.NotificationStatusPending,
		Priority:     notification.Priority,
		Recipient:    notification.Recipient,
		Subject:      notification.Subject,
		Content:      notification.Content,
		TemplateID:   notification.TemplateID,
		TemplateData: notification.TemplateData,
		Metadata:     notification.Metadata,
		ScheduledAt:  notification.ScheduledAt,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Set default priority if not specified
	if notificationRecord.Priority == "" {
		notificationRecord.Priority = models.NotificationPriorityNormal
	}

	// Save notification record
	if err := s.notificationRepo.Create(notificationRecord); err != nil {
		return fmt.Errorf("failed to create notification record: %w", err)
	}

	// Process the notification based on type
	var err error
	switch notification.Type {
	case models.NotificationTypeEmail:
		err = s.processEmailNotification(ctx, notificationRecord)
	case models.NotificationTypeSMS:
		err = s.processSMSNotification(ctx, notificationRecord)
	case models.NotificationTypeDatabase:
		err = s.processDatabaseNotification(ctx, notificationRecord)
	default:
		err = fmt.Errorf("unsupported notification type: %s", notification.Type)
	}

	// Update notification status
	if err != nil {
		errorMsg := err.Error()
		notificationRecord.Status = models.NotificationStatusFailed
		notificationRecord.ErrorMsg = &errorMsg
		notificationRecord.RetryCount++
	} else {
		notificationRecord.Status = models.NotificationStatusSent
		now := time.Now()
		notificationRecord.SentAt = &now
	}

	notificationRecord.UpdatedAt = time.Now()
	if updateErr := s.notificationRepo.Update(notificationRecord); updateErr != nil {
		fmt.Printf("Warning: failed to update notification status: %v\n", updateErr)
	}

	return err
}

// SendEmail sends an email notification
func (s *NotificationService) SendEmail(ctx context.Context, email *models.SendEmailRequest) error {
	// Convert to EmailNotification
	emailNotification := &models.EmailNotification{
		To:          email.To,
		CC:          email.CC,
		BCC:         email.BCC,
		Subject:     email.Subject,
		Body:        email.Body,
		HTMLBody:    email.HTMLBody,
		Attachments: email.Attachments,
		Priority:    email.Priority,
	}

	// Set default priority if not specified
	if emailNotification.Priority == "" {
		emailNotification.Priority = models.NotificationPriorityNormal
	}

	return s.emailService.SendEmail(ctx, emailNotification)
}

// SendSMS sends an SMS notification
func (s *NotificationService) SendSMS(ctx context.Context, sms *models.SendSMSRequest) error {
	// Convert to SMSNotification
	smsNotification := &models.SMSNotification{
		To:       sms.To,
		Message:  sms.Message,
		Priority: sms.Priority,
	}

	// Set default priority if not specified
	if smsNotification.Priority == "" {
		smsNotification.Priority = models.NotificationPriorityNormal
	}

	return s.smsService.SendSMS(ctx, smsNotification)
}

// SendDatabaseNotification sends a database notification
func (s *NotificationService) SendDatabaseNotification(ctx context.Context, dbNotification *models.DatabaseNotification) error {
	// Set default priority if not specified
	if dbNotification.Priority == "" {
		dbNotification.Priority = models.NotificationPriorityNormal
	}

	return s.dbNotificationService.CreateNotification(ctx, dbNotification)
}

// GetNotification retrieves a notification by ID
func (s *NotificationService) GetNotification(ctx context.Context, id int) (*models.Notification, error) {
	return s.notificationRepo.GetByID(id)
}

// GetNotifications retrieves notifications with filtering
func (s *NotificationService) GetNotifications(ctx context.Context, filters map[string]interface{}) ([]*models.Notification, error) {
	return s.notificationRepo.GetWithFilters(filters)
}

// UpdateNotificationStatus updates the status of a notification
func (s *NotificationService) UpdateNotificationStatus(ctx context.Context, id int, status models.NotificationStatus, errorMsg *string) error {
	notification, err := s.notificationRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get notification: %w", err)
	}

	notification.Status = status
	notification.ErrorMsg = errorMsg
	notification.UpdatedAt = time.Now()

	if status == models.NotificationStatusSent {
		now := time.Now()
		notification.SentAt = &now
	}

	return s.notificationRepo.Update(notification)
}

// RetryFailedNotification retries a failed notification
func (s *NotificationService) RetryFailedNotification(ctx context.Context, id int) error {
	notification, err := s.notificationRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get notification: %w", err)
	}

	if notification.Status != models.NotificationStatusFailed {
		return fmt.Errorf("notification is not in failed status")
	}

	// Reset status and retry
	notification.Status = models.NotificationStatusPending
	notification.ErrorMsg = nil
	notification.RetryCount++
	notification.UpdatedAt = time.Now()

	if err := s.notificationRepo.Update(notification); err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	// Process the notification again
	var processErr error
	switch notification.Type {
	case models.NotificationTypeEmail:
		processErr = s.processEmailNotification(ctx, notification)
	case models.NotificationTypeSMS:
		processErr = s.processSMSNotification(ctx, notification)
	case models.NotificationTypeDatabase:
		processErr = s.processDatabaseNotification(ctx, notification)
	default:
		processErr = fmt.Errorf("unsupported notification type: %s", notification.Type)
	}

	// Update status based on result
	if processErr != nil {
		errorMsg := processErr.Error()
		notification.Status = models.NotificationStatusFailed
		notification.ErrorMsg = &errorMsg
	} else {
		notification.Status = models.NotificationStatusSent
		now := time.Now()
		notification.SentAt = &now
	}

	notification.UpdatedAt = time.Now()
	return s.notificationRepo.Update(notification)
}

// processEmailNotification processes an email notification
func (s *NotificationService) processEmailNotification(ctx context.Context, notification *models.Notification) error {
	var emailNotification *models.EmailNotification

	// Check if we need to use a template
	if notification.TemplateID != nil {
		template, err := s.templateService.GetTemplate(ctx, *notification.TemplateID)
		if err != nil {
			return fmt.Errorf("failed to get template: %w", err)
		}

		// Render template
		subject, content, err := s.templateService.RenderTemplate(ctx, template, notification.TemplateData)
		if err != nil {
			return fmt.Errorf("failed to render template: %w", err)
		}

		emailNotification = &models.EmailNotification{
			To:       []string{notification.Recipient},
			Subject:  subject,
			HTMLBody: content, // Assuming template returns HTML content
			Priority: notification.Priority,
		}
	} else {
		// Use direct content
		emailNotification = &models.EmailNotification{
			To:       []string{notification.Recipient},
			Subject:  notification.Subject,
			Body:     notification.Content,
			Priority: notification.Priority,
		}
	}

	return s.emailService.SendEmail(ctx, emailNotification)
}

// processSMSNotification processes an SMS notification
func (s *NotificationService) processSMSNotification(ctx context.Context, notification *models.Notification) error {
	var smsNotification *models.SMSNotification

	// Check if we need to use a template
	if notification.TemplateID != nil {
		template, err := s.templateService.GetTemplate(ctx, *notification.TemplateID)
		if err != nil {
			return fmt.Errorf("failed to get template: %w", err)
		}

		// Render template
		_, content, err := s.templateService.RenderTemplate(ctx, template, notification.TemplateData)
		if err != nil {
			return fmt.Errorf("failed to render template: %w", err)
		}

		smsNotification = &models.SMSNotification{
			To:       notification.Recipient,
			Message:  content,
			Priority: notification.Priority,
		}
	} else {
		// Use direct content
		smsNotification = &models.SMSNotification{
			To:       notification.Recipient,
			Message:  notification.Content,
			Priority: notification.Priority,
		}
	}

	return s.smsService.SendSMS(ctx, smsNotification)
}

// processDatabaseNotification processes a database notification
func (s *NotificationService) processDatabaseNotification(ctx context.Context, notification *models.Notification) error {
	var dbNotification *models.DatabaseNotification

	// Check if we need to use a template
	if notification.TemplateID != nil {
		template, err := s.templateService.GetTemplate(ctx, *notification.TemplateID)
		if err != nil {
			return fmt.Errorf("failed to get template: %w", err)
		}

		// Render template
		subject, content, err := s.templateService.RenderTemplate(ctx, template, notification.TemplateData)
		if err != nil {
			return fmt.Errorf("failed to render template: %w", err)
		}

		// Extract user ID from metadata or template data
		userID, ok := notification.TemplateData["user_id"].(int)
		if !ok {
			return fmt.Errorf("user_id not found in template data")
		}

		dbNotification = &models.DatabaseNotification{
			UserID:   userID,
			Title:    subject,
			Message:  content,
			Type:     "system",
			Priority: notification.Priority,
		}
	} else {
		// Use direct content
		userID, ok := notification.Metadata["user_id"].(int)
		if !ok {
			return fmt.Errorf("user_id not found in metadata")
		}

		dbNotification = &models.DatabaseNotification{
			UserID:   userID,
			Title:    notification.Subject,
			Message:  notification.Content,
			Type:     "system",
			Priority: notification.Priority,
		}
	}

	return s.dbNotificationService.CreateNotification(ctx, dbNotification)
}
