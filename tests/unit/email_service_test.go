package unit

import (
	"context"
	"testing"

	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmailService is a mock implementation of EmailServiceInterface
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendEmail(ctx context.Context, email *models.EmailNotification) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockEmailService) SendBulkEmail(ctx context.Context, emails []*models.EmailNotification) error {
	args := m.Called(ctx, emails)
	return args.Error(0)
}

func (m *MockEmailService) ValidateEmailAddress(email string) bool {
	args := m.Called(email)
	return args.Bool(0)
}

func (m *MockEmailService) TestConnection(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestEmailService_ValidateEmailAddress(t *testing.T) {
	// Create email service with test config
	emailConfig := &config.EmailConfig{
		Enabled:   true,
		SMTPHost:  "smtp.example.com",
		SMTPPort:  587,
		SMTPUser:  "test@example.com",
		SMTPPass:  "password",
		FromEmail: "noreply@example.com",
		FromName:  "Test App",
		UseTLS:    true,
		UseSSL:    false,
	}

	emailService := services.NewEmailService(emailConfig)

	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "valid email",
			email:    "user@example.com",
			expected: true,
		},
		{
			name:     "valid email with subdomain",
			email:    "user@mail.example.com",
			expected: true,
		},
		{
			name:     "valid email with plus sign",
			email:    "user+tag@example.com",
			expected: true,
		},
		{
			name:     "invalid email - no domain",
			email:    "user@",
			expected: false,
		},
		{
			name:     "invalid email - no @ symbol",
			email:    "userexample.com",
			expected: false,
		},
		{
			name:     "invalid email - no local part",
			email:    "@example.com",
			expected: false,
		},
		{
			name:     "invalid email - empty string",
			email:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := emailService.ValidateEmailAddress(tt.email)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEmailService_Disabled(t *testing.T) {
	// Create disabled email service
	emailConfig := &config.EmailConfig{
		Enabled: false,
	}

	emailService := services.NewEmailService(emailConfig)

	// Test that sending email returns error when disabled
	email := &models.EmailNotification{
		To:      []string{"test@example.com"},
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	err := emailService.SendEmail(context.TODO(), email)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email service is disabled")
}

func TestEmailService_InvalidEmailAddresses(t *testing.T) {
	emailConfig := &config.EmailConfig{
		Enabled:   true,
		SMTPHost:  "smtp.example.com",
		SMTPPort:  587,
		SMTPUser:  "test@example.com",
		SMTPPass:  "password",
		FromEmail: "noreply@example.com",
		FromName:  "Test App",
		UseTLS:    true,
		UseSSL:    false,
	}

	emailService := services.NewEmailService(emailConfig)

	tests := []struct {
		name  string
		email *models.EmailNotification
	}{
		{
			name: "invalid TO addresses",
			email: &models.EmailNotification{
				To:      []string{"invalid-email"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
		},
		{
			name: "invalid CC addresses",
			email: &models.EmailNotification{
				To:      []string{"valid@example.com"},
				CC:      []string{"invalid-email"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
		},
		{
			name: "invalid BCC addresses",
			email: &models.EmailNotification{
				To:      []string{"valid@example.com"},
				BCC:     []string{"invalid-email"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := emailService.SendEmail(context.TODO(), tt.email)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid email address")
		})
	}
}

func TestEmailService_PriorityHeaders(t *testing.T) {
	emailConfig := &config.EmailConfig{
		Enabled:   true,
		SMTPHost:  "smtp.example.com",
		SMTPPort:  587,
		SMTPUser:  "test@example.com",
		SMTPPass:  "password",
		FromEmail: "noreply@example.com",
		FromName:  "Test App",
		UseTLS:    true,
		UseSSL:    false,
	}

	emailService := services.NewEmailService(emailConfig)

	// Test high priority email
	highPriorityEmail := &models.EmailNotification{
		To:       []string{"user@example.com"},
		Subject:  "High Priority Test",
		Body:     "This is a high priority email",
		Priority: models.NotificationPriorityHigh,
	}

	err := emailService.SendEmail(context.Background(), highPriorityEmail)
	// This will fail because we don't have a real SMTP server, but we can test the validation
	assert.Error(t, err)
	// The error should be about connection, not validation
	assert.NotContains(t, err.Error(), "invalid email address")

	// Test low priority email
	lowPriorityEmail := &models.EmailNotification{
		To:       []string{"user@example.com"},
		Subject:  "Low Priority Test",
		Body:     "This is a low priority email",
		Priority: models.NotificationPriorityLow,
	}

	err = emailService.SendEmail(context.Background(), lowPriorityEmail)
	assert.Error(t, err)
	assert.NotContains(t, err.Error(), "invalid email address")
}
