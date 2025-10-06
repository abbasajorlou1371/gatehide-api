package seeders

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

// init registers the notification template seeder
func init() {
	RegisterSeeder("notification_templates", SeedNotificationTemplates)
}

// NotificationTemplateSeeder seeds default notification templates
type NotificationTemplateSeeder struct {
	db *sql.DB
}

// NewNotificationTemplateSeeder creates a new notification template seeder
func NewNotificationTemplateSeeder(cfg *config.Config) (*NotificationTemplateSeeder, error) {
	db, err := sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &NotificationTemplateSeeder{db: db}, nil
}

// SeedNotificationTemplates is the public seeder function that can be called by the registry
func SeedNotificationTemplates(cfg *config.Config) error {
	seeder, err := NewNotificationTemplateSeeder(cfg)
	if err != nil {
		return fmt.Errorf("failed to create notification template seeder: %w", err)
	}
	defer seeder.Close()

	return seeder.seedTemplates()
}

// seedTemplates seeds the default notification templates
func (s *NotificationTemplateSeeder) seedTemplates() error {
	log.Println("Seeding notification templates...")

	templates := []*models.NotificationTemplate{
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

	for _, template := range templates {
		// Check if template already exists
		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM notification_templates WHERE name = ? AND type = ?",
			template.Name, template.Type).Scan(&count)
		if err != nil {
			log.Printf("Error checking template existence: %v", err)
			continue
		}

		if count > 0 {
			log.Printf("Template %s (%s) already exists, skipping...", template.Name, template.Type)
			continue
		}

		// Insert template
		variablesJSON, _ := json.Marshal(template.Variables)

		_, err = s.db.Exec(`
			INSERT INTO notification_templates (name, type, subject, content, html_content, variables, is_active, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, template.Name, template.Type, template.Subject, template.Content,
			template.HTMLContent, variablesJSON, template.IsActive,
			template.CreatedAt, template.UpdatedAt)

		if err != nil {
			log.Printf("Error seeding template %s: %v", template.Name, err)
			continue
		}

		log.Printf("✅ Seeded template: %s (%s)", template.Name, template.Type)
	}

	log.Println("Notification template seeding completed!")
	return nil
}

// Close closes the database connection
func (s *NotificationTemplateSeeder) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
