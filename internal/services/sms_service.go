// Package services provides SMS notification functionality using Kavenegar API.
//
// Example usage:
//
//	cfg := &config.SMSConfig{
//	    Enabled:    true,
//	    APIKey:     "your-kavenegar-api-key",
//	    Sender:     "10008663",
//	    TestMode:   true,
//	    MaxRetries: 3,
//	}
//	smsService := NewSMSService(cfg)
//
//	sms := &models.SMSNotification{
//	    To:      "09123456789",
//	    Message: "Hello from GateHide!",
//	}
//	err := smsService.SendSMS(context.Background(), sms)
package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gatehide/gatehide-api/config"
	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/kavenegar/kavenegar-go"
)

// SMSService implements SMSServiceInterface using Kavenegar
type SMSService struct {
	client *kavenegar.Kavenegar
	config *config.SMSConfig
}

// NewSMSService creates a new SMS service instance
func NewSMSService(cfg *config.SMSConfig) *SMSService {
	if !cfg.Enabled || cfg.APIKey == "" {
		return &SMSService{
			client: nil,
			config: cfg,
		}
	}

	client := kavenegar.New(cfg.APIKey)
	return &SMSService{
		client: client,
		config: cfg,
	}
}

// SendSMS sends an SMS message using Kavenegar
func (s *SMSService) SendSMS(ctx context.Context, sms *models.SMSNotification) error {
	if !s.config.Enabled {
		return fmt.Errorf("SMS service is disabled")
	}

	if s.client == nil {
		return fmt.Errorf("SMS service not properly configured")
	}

	// Validate phone number
	if !s.ValidatePhoneNumber(sms.To) {
		return fmt.Errorf("invalid phone number: %s", sms.To)
	}

	// Normalize phone number
	phoneNumber := s.normalizePhoneNumber(sms.To)

	// Prepare message
	message := strings.TrimSpace(sms.Message)
	if message == "" {
		return fmt.Errorf("message cannot be empty")
	}

	// Set timeout for the request
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Send SMS with retry logic
	var lastErr error
	for attempt := 1; attempt <= s.config.MaxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Send the SMS
		receptor := []string{phoneNumber}
		sender := s.config.Sender
		if s.config.TestMode {
			// In test mode, we might want to use a different sender or add test prefix
			message = fmt.Sprintf("[TEST] %s", message)
		}

		res, err := s.client.Message.Send(sender, receptor, message, nil)
		if err != nil {
			lastErr = err
			if attempt < s.config.MaxRetries {
				// Wait before retry (exponential backoff)
				time.Sleep(time.Duration(attempt) * time.Second)
				continue
			}
			return s.handleKavenegarError(err)
		}

		// Check if the response indicates success
		if len(res) > 0 && res[0].Status == 1 {
			return nil
		}

		// If we get here, the response indicates failure
		lastErr = fmt.Errorf("SMS sending failed with status: %d", res[0].Status)
		if attempt < s.config.MaxRetries {
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}
	}

	return fmt.Errorf("SMS sending failed after %d attempts: %w", s.config.MaxRetries, lastErr)
}

// SendBulkSMS sends multiple SMS messages
func (s *SMSService) SendBulkSMS(ctx context.Context, smsMessages []*models.SMSNotification) error {
	if !s.config.Enabled {
		return fmt.Errorf("SMS service is disabled")
	}

	if s.client == nil {
		return fmt.Errorf("SMS service not properly configured")
	}

	if len(smsMessages) == 0 {
		return fmt.Errorf("no SMS messages to send")
	}

	// Validate all phone numbers and messages
	for i, sms := range smsMessages {
		if !s.ValidatePhoneNumber(sms.To) {
			return fmt.Errorf("invalid phone number at index %d: %s", i, sms.To)
		}
		if strings.TrimSpace(sms.Message) == "" {
			return fmt.Errorf("empty message at index %d", i)
		}
	}

	// Set timeout for the request
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Send each SMS individually with retry logic
	var lastErr error
	for i, sms := range smsMessages {
		phoneNumber := s.normalizePhoneNumber(sms.To)
		message := strings.TrimSpace(sms.Message)

		if s.config.TestMode {
			message = fmt.Sprintf("[TEST] %s", message)
		}

		// Send individual SMS with retry logic
		for attempt := 1; attempt <= s.config.MaxRetries; attempt++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// Send the SMS
			res, err := s.client.Message.Send(s.config.Sender, []string{phoneNumber}, message, nil)
			if err != nil {
				lastErr = err
				if attempt < s.config.MaxRetries {
					time.Sleep(time.Duration(attempt) * time.Second)
					continue
				}
				return s.handleKavenegarError(err)
			}

			// Check if the response indicates success
			if len(res) > 0 && res[0].Status == 1 {
				break // Success, move to next SMS
			}

			// If we get here, the response indicates failure
			lastErr = fmt.Errorf("SMS sending failed for message %d with status: %d", i+1, res[0].Status)
			if attempt < s.config.MaxRetries {
				time.Sleep(time.Duration(attempt) * time.Second)
				continue
			}
		}

		// If we couldn't send this SMS after all retries, return error
		if lastErr != nil {
			return fmt.Errorf("failed to send SMS %d after %d attempts: %w", i+1, s.config.MaxRetries, lastErr)
		}
	}

	return nil
}

// ValidatePhoneNumber validates a phone number format
func (s *SMSService) ValidatePhoneNumber(phone string) bool {
	if phone == "" {
		return false
	}

	// Remove all non-digit characters
	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	// Check if it's a valid Iranian mobile number (09xxxxxxxxx)
	if len(cleaned) == 11 && strings.HasPrefix(cleaned, "09") {
		return true
	}

	// Check if it's a valid international number (starts with +98)
	if strings.HasPrefix(phone, "+98") {
		cleaned = strings.TrimPrefix(cleaned, "98")
		if len(cleaned) == 11 && strings.HasPrefix(cleaned, "09") {
			return true
		}
	}

	// Check if it's already in international format without +
	if len(cleaned) == 12 && strings.HasPrefix(cleaned, "98") {
		cleaned = strings.TrimPrefix(cleaned, "98")
		if len(cleaned) == 11 && strings.HasPrefix(cleaned, "09") {
			return true
		}
	}

	return false
}

// TestConnection tests the SMS service connection
func (s *SMSService) TestConnection(ctx context.Context) error {
	if !s.config.Enabled {
		return fmt.Errorf("SMS service is disabled")
	}

	if s.client == nil {
		return fmt.Errorf("SMS service not properly configured")
	}

	// Set timeout for the request
	_, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Try to get account info to test the connection
	_, err := s.client.Account.Info()
	if err != nil {
		return s.handleKavenegarError(err)
	}

	return nil
}

// normalizePhoneNumber normalizes a phone number to the format expected by Kavenegar
func (s *SMSService) normalizePhoneNumber(phone string) string {
	// Remove all non-digit characters
	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	// If it starts with +98, remove the +
	if strings.HasPrefix(phone, "+98") {
		return cleaned
	}

	// If it starts with 98, keep as is
	if strings.HasPrefix(cleaned, "98") {
		return cleaned
	}

	// If it starts with 09, add 98 prefix
	if strings.HasPrefix(cleaned, "09") {
		return "98" + cleaned
	}

	// If it's 11 digits and starts with 9, add 98 prefix
	if len(cleaned) == 11 && strings.HasPrefix(cleaned, "9") {
		return "98" + cleaned
	}

	// Return as is if it doesn't match any pattern
	return cleaned
}

// SendGamenetCredentials sends gamenet credentials using Kavenegar Verify Lookup or regular SMS as fallback
func (s *SMSService) SendGamenetCredentials(ctx context.Context, mobile, email, password string) error {
	if !s.config.Enabled {
		return fmt.Errorf("SMS service is disabled")
	}

	if s.client == nil {
		return fmt.Errorf("SMS service not properly configured")
	}

	// Validate phone number
	if !s.ValidatePhoneNumber(mobile) {
		return fmt.Errorf("invalid phone number: %s", mobile)
	}

	// Normalize phone number
	phoneNumber := s.normalizePhoneNumber(mobile)

	// Set timeout for the request
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Template name from Kavenegar panel
	template := "gamenet-credentials"

	// Try Verify Lookup first (preferred method with template)
	params := &kavenegar.VerifyLookupParam{
		Token2: email,
		Token3: password,
	}

	res, err := s.client.Verify.Lookup(phoneNumber, template, email, params)
	if err != nil {
		// Check if it's a template not found error (424)
		if apiErr, ok := err.(*kavenegar.APIError); ok && apiErr.Status == 424 {
			// Template doesn't exist, use fallback to regular SMS
			fmt.Printf("Template not found, using regular SMS fallback\n")
			return s.sendCredentialsViaSMS(ctx, phoneNumber, email, password)
		}
		return s.handleKavenegarError(err)
	}

	// Check if the response indicates success
	if res.Status == 200 {
		fmt.Printf("Successfully sent credentials SMS via Verify Lookup to %s\n", mobile)
		return nil
	}

	// If Verify Lookup failed, try regular SMS as fallback
	fmt.Printf("Verify Lookup returned status %d, using regular SMS fallback\n", res.Status)
	return s.sendCredentialsViaSMS(ctx, phoneNumber, email, password)
}

// sendCredentialsViaSMS sends credentials using regular SMS (fallback method)
func (s *SMSService) sendCredentialsViaSMS(ctx context.Context, phoneNumber, email, password string) error {
	// ctx parameter is required by interface but not used in this implementation
	// Construct message
	message := fmt.Sprintf("اطلاعات ورود به سیستم گیت نت:\nایمیل: %s\nرمز عبور: %s", email, password)

	// Send the SMS (NO RETRIES to avoid duplicate credentials)
	receptor := []string{phoneNumber}

	// Use empty sender to use account's default sender line (as per Kavenegar SDK docs)
	// Reference: https://github.com/kavenegar/kavenegar-go
	sender := ""

	fmt.Printf("📤 Sending credentials SMS to %s...\n", phoneNumber)
	res, err := s.client.Message.Send(sender, receptor, message, nil)
	if err != nil {
		fmt.Printf("❌ SMS Error: %v\n", err)
		return s.handleKavenegarError(err)
	}

	// Check if the response indicates message was accepted
	if len(res) > 0 {
		fmt.Printf("📱 SMS Response: Status = %d, MessageID = %d\n", res[0].Status, res[0].MessageID)

		// Kavenegar Status Codes:
		// 1 = Queued (normal success)
		// 5 = Sent with sender warning (still delivered)
		// Accept both status 1 and 5 as success since message is delivered
		if res[0].Status == 1 || res[0].Status == 5 {
			fmt.Printf("✅ Credentials SMS sent successfully to %s (MessageID: %d)\n", phoneNumber, res[0].MessageID)
			return nil
		}

		return fmt.Errorf("SMS sending failed with status: %d", res[0].Status)
	}

	return fmt.Errorf("SMS sending failed: no response from Kavenegar")
}

// SendUserCredentials sends user credentials using Kavenegar Verify Lookup or regular SMS as fallback
func (s *SMSService) SendUserCredentials(ctx context.Context, mobile, email, password string) error {
	if !s.config.Enabled {
		return fmt.Errorf("SMS service is disabled")
	}

	if s.client == nil {
		return fmt.Errorf("SMS service not properly configured")
	}

	// Validate phone number
	if !s.ValidatePhoneNumber(mobile) {
		return fmt.Errorf("invalid phone number: %s", mobile)
	}

	// Normalize phone number
	phoneNumber := s.normalizePhoneNumber(mobile)

	// Set timeout for the request
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Template name from Kavenegar panel
	template := "user-credentials"

	// Try Verify Lookup first (preferred method with template)
	params := &kavenegar.VerifyLookupParam{
		Token2: email,
		Token3: password,
	}

	res, err := s.client.Verify.Lookup(phoneNumber, template, email, params)
	if err != nil {
		// Check if it's a template not found error (424)
		if apiErr, ok := err.(*kavenegar.APIError); ok && apiErr.Status == 424 {
			// Template doesn't exist, use fallback to regular SMS
			fmt.Printf("Template not found, using regular SMS fallback\n")
			return s.sendUserCredentialsViaSMS(ctx, phoneNumber, email, password)
		}
		return s.handleKavenegarError(err)
	}

	// Check if the response indicates success
	if res.Status == 200 {
		fmt.Printf("Successfully sent user credentials SMS via Verify Lookup to %s\n", mobile)
		return nil
	}

	// If Verify Lookup failed, try regular SMS as fallback
	fmt.Printf("Verify Lookup returned status %d, using regular SMS fallback\n", res.Status)
	return s.sendUserCredentialsViaSMS(ctx, phoneNumber, email, password)
}

// sendUserCredentialsViaSMS sends user credentials using regular SMS (fallback method)
func (s *SMSService) sendUserCredentialsViaSMS(ctx context.Context, phoneNumber, email, password string) error {
	// ctx parameter is required by interface but not used in this implementation
	// Construct message
	message := fmt.Sprintf("اطلاعات ورود به سیستم:\nایمیل: %s\nرمز عبور: %s", email, password)

	// Send the SMS (NO RETRIES to avoid duplicate credentials)
	receptor := []string{phoneNumber}

	// Use empty sender to use account's default sender line
	sender := ""

	fmt.Printf("📤 Sending user credentials SMS to %s...\n", phoneNumber)
	res, err := s.client.Message.Send(sender, receptor, message, nil)
	if err != nil {
		fmt.Printf("❌ SMS Error: %v\n", err)
		return s.handleKavenegarError(err)
	}

	// Check if the response indicates message was accepted
	if len(res) > 0 {
		fmt.Printf("📱 SMS Response: Status = %d, MessageID = %d\n", res[0].Status, res[0].MessageID)

		// Accept both status 1 and 5 as success
		if res[0].Status == 1 || res[0].Status == 5 {
			fmt.Printf("✅ User credentials SMS sent successfully to %s (MessageID: %d)\n", phoneNumber, res[0].MessageID)
			return nil
		}

		return fmt.Errorf("SMS sending failed with status: %d", res[0].Status)
	}

	return fmt.Errorf("SMS sending failed: no response from Kavenegar")
}

// handleKavenegarError handles Kavenegar-specific errors
func (s *SMSService) handleKavenegarError(err error) error {
	switch e := err.(type) {
	case *kavenegar.APIError:
		return fmt.Errorf("kavenegar API error: %s (Status: %d)", e.Error(), e.Status)
	case *kavenegar.HTTPError:
		return fmt.Errorf("kavenegar HTTP error: %s", e.Error())
	default:
		return fmt.Errorf("SMS service error: %w", err)
	}
}
