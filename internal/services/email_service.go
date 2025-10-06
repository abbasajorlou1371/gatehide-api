package services

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"regexp"
	"strings"
	"time"

	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/internal/models"
)

// EmailService implements EmailServiceInterface for SMTP email sending
type EmailService struct {
	config *config.EmailConfig
}

// NewEmailService creates a new email service instance
func NewEmailService(cfg *config.EmailConfig) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

// SendEmail sends an email using SMTP
func (s *EmailService) SendEmail(ctx context.Context, email *models.EmailNotification) error {
	if !s.config.Enabled {
		return fmt.Errorf("email service is disabled")
	}

	// Validate email addresses
	if err := s.validateEmailAddresses(email.To); err != nil {
		return fmt.Errorf("invalid recipient addresses: %w", err)
	}
	if err := s.validateEmailAddresses(email.CC); err != nil {
		return fmt.Errorf("invalid CC addresses: %w", err)
	}
	if err := s.validateEmailAddresses(email.BCC); err != nil {
		return fmt.Errorf("invalid BCC addresses: %w", err)
	}

	// Create message
	message, err := s.createMessage(email)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	// Connect to SMTP server
	conn, err := s.connectSMTP(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Authenticate
	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPass, s.config.SMTPHost)

	// Send email
	recipients := append(email.To, email.CC...)
	recipients = append(recipients, email.BCC...)

	if err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort),
		auth,
		s.config.FromEmail,
		recipients,
		message,
	); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendBulkEmail sends multiple emails
func (s *EmailService) SendBulkEmail(ctx context.Context, emails []*models.EmailNotification) error {
	if !s.config.Enabled {
		return fmt.Errorf("email service is disabled")
	}

	var errors []string
	for i, email := range emails {
		if err := s.SendEmail(ctx, email); err != nil {
			errors = append(errors, fmt.Sprintf("email %d: %v", i+1, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("bulk email send failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

// ValidateEmailAddress validates an email address format
func (s *EmailService) ValidateEmailAddress(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// TestConnection tests the SMTP connection
func (s *EmailService) TestConnection(ctx context.Context) error {
	if !s.config.Enabled {
		return fmt.Errorf("email service is disabled")
	}

	conn, err := s.connectSMTP(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Test authentication
	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPass, s.config.SMTPHost)
	if err := conn.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	return nil
}

// connectSMTP establishes a connection to the SMTP server
func (s *EmailService) connectSMTP(ctx context.Context) (*smtp.Client, error) {
	// Create dialer with timeout
	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SMTP server: %w", err)
	}

	// Create SMTP client
	client, err := smtp.NewClient(conn, s.config.SMTPHost)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create SMTP client: %w", err)
	}

	// Configure TLS if needed
	if s.config.UseTLS || s.config.UseSSL {
		tlsConfig := &tls.Config{
			ServerName: s.config.SMTPHost,
		}

		if s.config.UseSSL {
			// For SSL, start TLS immediately
			if err := client.StartTLS(tlsConfig); err != nil {
				client.Close()
				return nil, fmt.Errorf("failed to start TLS: %w", err)
			}
		} else if s.config.UseTLS {
			// For TLS, check if server supports STARTTLS
			if ok, _ := client.Extension("STARTTLS"); ok {
				if err := client.StartTLS(tlsConfig); err != nil {
					client.Close()
					return nil, fmt.Errorf("failed to start TLS: %w", err)
				}
			}
		}
	}

	return client, nil
}

// createMessage creates the email message in RFC 2822 format
func (s *EmailService) createMessage(email *models.EmailNotification) ([]byte, error) {
	var message strings.Builder

	// Headers
	message.WriteString(fmt.Sprintf("From: %s <%s>\r\n", s.config.FromName, s.config.FromEmail))
	message.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ", ")))

	if len(email.CC) > 0 {
		message.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(email.CC, ", ")))
	}

	message.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))
	message.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	message.WriteString("MIME-Version: 1.0\r\n")

	// Gmail and Yahoo compliance headers
	message.WriteString("X-Mailer: GateHide API v1.0\r\n")
	message.WriteString("X-Report-Abuse: Please report abuse to abuse@gatehide.com\r\n")
	message.WriteString("List-Unsubscribe: <http://localhost:3000/unsubscribe>, <mailto:unsubscribe@gatehide.com>\r\n")
	message.WriteString("List-Unsubscribe-Post: List-Unsubscribe=One-Click\r\n")
	message.WriteString("Precedence: bulk\r\n")
	message.WriteString("X-Auto-Response-Suppress: All\r\n")

	// Priority header
	switch email.Priority {
	case models.NotificationPriorityHigh, models.NotificationPriorityUrgent:
		message.WriteString("X-Priority: 1\r\n")
		message.WriteString("X-MSMail-Priority: High\r\n")
	case models.NotificationPriorityLow:
		message.WriteString("X-Priority: 5\r\n")
		message.WriteString("X-MSMail-Priority: Low\r\n")
	default:
		message.WriteString("X-Priority: 3\r\n")
		message.WriteString("X-MSMail-Priority: Normal\r\n")
	}

	// Content type and body
	if email.HTMLBody != "" && email.Body != "" {
		// Multipart message with both text and HTML
		boundary := fmt.Sprintf("boundary_%d", time.Now().Unix())
		message.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", boundary))
		message.WriteString("\r\n")

		// Text part
		message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		message.WriteString("\r\n")
		message.WriteString(email.Body)
		message.WriteString("\r\n")

		// HTML part
		message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		message.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		message.WriteString("\r\n")
		message.WriteString(email.HTMLBody)
		message.WriteString("\r\n")
		message.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else if email.HTMLBody != "" {
		// HTML only
		message.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		message.WriteString("\r\n")
		message.WriteString(email.HTMLBody)
	} else {
		// Text only
		message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		message.WriteString("\r\n")
		message.WriteString(email.Body)
	}

	return []byte(message.String()), nil
}

// validateEmailAddresses validates a slice of email addresses
func (s *EmailService) validateEmailAddresses(emails []string) error {
	for _, email := range emails {
		if !s.ValidateEmailAddress(email) {
			return fmt.Errorf("invalid email address: %s", email)
		}
	}
	return nil
}
