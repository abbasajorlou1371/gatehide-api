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
			Name:    "password_reset_email",
			Type:    models.NotificationTypeEmail,
			Subject: "بازنشانی رمز عبور - {{app_name}}",
			Content: "کاربر گرامی {{user_name}}،\n\nدرخواست بازنشانی رمز عبور برای حساب کاربری شما در {{app_name}} دریافت شده است.\n\nبرای تنظیم رمز عبور جدید، لطفاً روی لینک زیر کلیک کنید:\n{{reset_link}}\n\nاین لینک تا {{expiry_hours}} ساعت معتبر است.\n\nاگر شما این درخواست را انجام نداده‌اید، لطفاً این ایمیل را نادیده بگیرید.\n\nبا احترام،\nتیم {{app_name}}",
			HTMLContent: `<!DOCTYPE html>
<html lang="fa" dir="rtl">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>بازنشانی رمز عبور</title>
    <style>
        body {
            font-family: 'Tahoma', 'Arial', sans-serif;
            direction: rtl;
            text-align: right;
            background-color: #f9f9f9;
            margin: 0;
            padding: 0;
            line-height: 1.6;
        }
        .container {
            max-width: 600px;
            margin: 20px auto;
            background-color: #ffffff;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
        }
        .header {
            text-align: center;
            padding-bottom: 20px;
            border-bottom: 2px solid #e9ecef;
            margin-bottom: 30px;
        }
        .header h1 {
            color: #2c3e50;
            margin: 0;
            font-size: 24px;
        }
        .content {
            color: #333333;
            font-size: 16px;
        }
        .content p {
            margin-bottom: 20px;
        }
        .button-container {
            text-align: center;
            margin: 30px 0;
        }
        .button {
            display: inline-block;
            padding: 15px 30px;
            background-color: #007BFF;
            color: #ffffff;
            text-decoration: none;
            border-radius: 5px;
            font-weight: bold;
            font-size: 16px;
            transition: background-color 0.3s ease;
        }
        .button:hover {
            background-color: #0056b3;
        }
        .warning {
            background-color: #fff3cd;
            border: 1px solid #ffeaa7;
            border-radius: 5px;
            padding: 15px;
            margin: 20px 0;
            color: #856404;
        }
        .footer {
            text-align: center;
            font-size: 12px;
            color: #666666;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #e9ecef;
        }
        .footer a {
            color: #007BFF;
            text-decoration: none;
        }
        .unsubscribe {
            margin-top: 15px;
            font-size: 11px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>بازنشانی رمز عبور</h1>
        </div>
        <div class="content">
            <p>کاربر گرامی {{user_name}}،</p>
            <p>درخواست بازنشانی رمز عبور برای حساب کاربری شما در <strong>{{app_name}}</strong> دریافت شده است.</p>
            <p>برای تنظیم رمز عبور جدید، لطفاً روی دکمه زیر کلیک کنید:</p>
            
            <div class="button-container">
                <a href="{{reset_link}}" class="button">بازنشانی رمز عبور</a>
            </div>
            
            <div class="warning">
                <strong>توجه:</strong> این لینک تا {{expiry_hours}} ساعت معتبر است و فقط یک بار قابل استفاده است.
            </div>
            
            <p>اگر شما این درخواست را انجام نداده‌اید، لطفاً این ایمیل را نادیده بگیرید. رمز عبور شما تغییر نخواهد کرد.</p>
        </div>
        <div class="footer">
            <p>© 2025 {{app_name}}. تمامی حقوق محفوظ است.</p>
            <div class="unsubscribe">
                <a href="{{unsubscribe_link}}">لغو اشتراک</a> | 
                <a href="{{support_link}}">پشتیبانی</a>
            </div>
        </div>
    </div>
</body>
</html>`,
			Variables: []string{"app_name", "user_name", "reset_link", "expiry_hours", "unsubscribe_link", "support_link"},
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
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
