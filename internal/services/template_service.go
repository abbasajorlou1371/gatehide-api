package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/repositories"
)

// TemplateService implements TemplateServiceInterface for managing notification templates
type TemplateService struct {
	templateRepo repositories.NotificationTemplateRepository
}

// NewTemplateService creates a new template service instance
func NewTemplateService(templateRepo repositories.NotificationTemplateRepository) *TemplateService {
	return &TemplateService{
		templateRepo: templateRepo,
	}
}

// GetTemplate retrieves a template by ID
func (s *TemplateService) GetTemplate(ctx context.Context, id int) (*models.NotificationTemplate, error) {
	return s.templateRepo.GetByID(id)
}

// GetTemplateByName retrieves a template by name and type
func (s *TemplateService) GetTemplateByName(ctx context.Context, name string, templateType models.NotificationType) (*models.NotificationTemplate, error) {
	return s.templateRepo.GetByNameAndType(name, templateType)
}

// CreateTemplate creates a new template
func (s *TemplateService) CreateTemplate(ctx context.Context, template *models.NotificationTemplate) error {
	// Validate template
	if err := s.validateTemplate(template); err != nil {
		return fmt.Errorf("template validation failed: %w", err)
	}

	// Extract variables from template content
	template.Variables = s.extractVariables(template.Subject, template.Content, template.HTMLContent)

	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()

	return s.templateRepo.Create(template)
}

// UpdateTemplate updates an existing template
func (s *TemplateService) UpdateTemplate(ctx context.Context, template *models.NotificationTemplate) error {
	// Validate template
	if err := s.validateTemplate(template); err != nil {
		return fmt.Errorf("template validation failed: %w", err)
	}

	// Extract variables from template content
	template.Variables = s.extractVariables(template.Subject, template.Content, template.HTMLContent)

	template.UpdatedAt = time.Now()

	return s.templateRepo.Update(template)
}

// DeleteTemplate deletes a template
func (s *TemplateService) DeleteTemplate(ctx context.Context, id int) error {
	return s.templateRepo.Delete(id)
}

// ListTemplates lists all templates with optional filtering
func (s *TemplateService) ListTemplates(ctx context.Context, templateType *models.NotificationType) ([]*models.NotificationTemplate, error) {
	if templateType != nil {
		return s.templateRepo.GetByType(*templateType)
	}
	return s.templateRepo.GetAll()
}

// RenderTemplate renders a template with provided data
func (s *TemplateService) RenderTemplate(ctx context.Context, template *models.NotificationTemplate, data map[string]interface{}) (string, string, error) {
	// Render subject
	subject, err := s.renderString(template.Subject, data)
	if err != nil {
		return "", "", fmt.Errorf("failed to render subject: %w", err)
	}

	// Render content
	var content string
	if template.HTMLContent != "" {
		content, err = s.renderString(template.HTMLContent, data)
		if err != nil {
			return "", "", fmt.Errorf("failed to render HTML content: %w", err)
		}
	} else {
		content, err = s.renderString(template.Content, data)
		if err != nil {
			return "", "", fmt.Errorf("failed to render content: %w", err)
		}
	}

	return subject, content, nil
}

// validateTemplate validates a template
func (s *TemplateService) validateTemplate(template *models.NotificationTemplate) error {
	if template.Name == "" {
		return fmt.Errorf("template name is required")
	}

	if template.Type == "" {
		return fmt.Errorf("template type is required")
	}

	if template.Subject == "" {
		return fmt.Errorf("template subject is required")
	}

	if template.Content == "" && template.HTMLContent == "" {
		return fmt.Errorf("template content or HTML content is required")
	}

	// Validate template type
	validTypes := []models.NotificationType{
		models.NotificationTypeEmail,
		models.NotificationTypeSMS,
		models.NotificationTypeDatabase,
	}

	validType := false
	for _, validT := range validTypes {
		if template.Type == validT {
			validType = true
			break
		}
	}

	if !validType {
		return fmt.Errorf("invalid template type: %s", template.Type)
	}

	return nil
}

// extractVariables extracts variable placeholders from template strings
func (s *TemplateService) extractVariables(strs ...string) []string {
	variableRegex := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	variables := make(map[string]bool)

	for _, str := range strs {
		if str == "" {
			continue
		}
		matches := variableRegex.FindAllStringSubmatch(str, -1)
		for _, match := range matches {
			if len(match) > 1 {
				variables[match[1]] = true
			}
		}
	}

	var result []string
	for variable := range variables {
		result = append(result, variable)
	}

	return result
}

// renderString renders a template string with provided data
func (s *TemplateService) renderString(template string, data map[string]interface{}) (string, error) {
	variableRegex := regexp.MustCompile(`\{\{([^}]+)\}\}`)

	result := variableRegex.ReplaceAllStringFunc(template, func(match string) string {
		// Extract variable name
		variableName := strings.TrimSpace(match[2 : len(match)-2])

		// Get value from data
		if value, exists := data[variableName]; exists {
			return fmt.Sprintf("%v", value)
		}

		// Return original match if variable not found
		return match
	})

	return result, nil
}

// GetDefaultTemplates returns a list of default templates
func (s *TemplateService) GetDefaultTemplates() []*models.NotificationTemplate {
	return []*models.NotificationTemplate{
		{
			Name:        "welcome_email",
			Type:        models.NotificationTypeEmail,
			Subject:     "Welcome to {{app_name}}, {{user_name}}!",
			Content:     "Hi {{user_name}},\n\nWelcome to {{app_name}}! We're excited to have you on board.\n\nBest regards,\nThe {{app_name}} Team",
			HTMLContent: "<h2>Welcome to {{app_name}}, {{user_name}}!</h2><p>Hi {{user_name}},</p><p>Welcome to {{app_name}}! We're excited to have you on board.</p><p>Best regards,<br>The {{app_name}} Team</p>",
			Variables:   []string{"app_name", "user_name"},
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "password_reset_email",
			Type:        models.NotificationTypeEmail,
			Subject:     "Password Reset - {{app_name}}",
			Content:     "Hi {{user_name}},\n\nYou requested a password reset for your {{app_name}} account.\n\nClick the link below to reset your password:\n{{reset_link}}\n\nThis link will expire in {{expiry_hours}} hours.\n\nIf you didn't request this, please ignore this email.\n\nBest regards,\nThe {{app_name}} Team",
			HTMLContent: "<h2>Password Reset - {{app_name}}</h2><p>Hi {{user_name}},</p><p>You requested a password reset for your {{app_name}} account.</p><p><a href=\"{{reset_link}}\">Click here to reset your password</a></p><p>This link will expire in {{expiry_hours}} hours.</p><p>If you didn't request this, please ignore this email.</p><p>Best regards,<br>The {{app_name}} Team</p>",
			Variables:   []string{"app_name", "user_name", "reset_link", "expiry_hours"},
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "login_notification_email",
			Type:        models.NotificationTypeEmail,
			Subject:     "New Login - {{app_name}}",
			Content:     "Hi {{user_name}},\n\nThere was a new login to your {{app_name}} account.\n\nDetails:\n- Time: {{login_time}}\n- IP Address: {{ip_address}}\n- Device: {{device}}\n\nIf this wasn't you, please secure your account immediately.\n\nBest regards,\nThe {{app_name}} Team",
			HTMLContent: "<h2>New Login - {{app_name}}</h2><p>Hi {{user_name}},</p><p>There was a new login to your {{app_name}} account.</p><p><strong>Details:</strong></p><ul><li>Time: {{login_time}}</li><li>IP Address: {{ip_address}}</li><li>Device: {{device}}</li></ul><p>If this wasn't you, please secure your account immediately.</p><p>Best regards,<br>The {{app_name}} Team</p>",
			Variables:   []string{"app_name", "user_name", "login_time", "ip_address", "device"},
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:      "welcome_sms",
			Type:      models.NotificationTypeSMS,
			Subject:   "Welcome to {{app_name}}",
			Content:   "Hi {{user_name}}, welcome to {{app_name}}! We're excited to have you on board.",
			Variables: []string{"app_name", "user_name"},
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "password_reset_sms",
			Type:      models.NotificationTypeSMS,
			Subject:   "Password Reset Code",
			Content:   "Your {{app_name}} password reset code is: {{reset_code}}. Valid for {{expiry_minutes}} minutes.",
			Variables: []string{"app_name", "reset_code", "expiry_minutes"},
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "system_notification",
			Type:      models.NotificationTypeDatabase,
			Subject:   "System Notification",
			Content:   "{{title}}\n\n{{message}}",
			Variables: []string{"title", "message"},
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}
