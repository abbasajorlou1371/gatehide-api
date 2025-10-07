package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// SubscriptionPlan represents a subscription plan in the system
type SubscriptionPlan struct {
	ID                       int       `json:"id" db:"id"`
	Name                     string    `json:"name" db:"name"`
	PlanType                 string    `json:"plan_type" db:"plan_type"`
	Price                    float64   `json:"price" db:"price"`
	AnnualDiscountPercentage *float64  `json:"annual_discount_percentage" db:"annual_discount_percentage"`
	TrialDurationDays        *int      `json:"trial_duration_days" db:"trial_duration_days"`
	IsActive                 bool      `json:"is_active" db:"is_active"`
	SubscriptionCount        int       `json:"subscription_count" db:"subscription_count"`
	CreatedAt                time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time `json:"updated_at" db:"updated_at"`
}

// UserSubscription represents a gamenet's current subscription
type UserSubscription struct {
	ID        int        `json:"id" db:"id"`
	GamenetID int        `json:"gamenet_id" db:"gamenet_id"`
	PlanID    int        `json:"plan_id" db:"plan_id"`
	Status    string     `json:"status" db:"status"`
	StartedAt time.Time  `json:"started_at" db:"started_at"`
	ExpiresAt *time.Time `json:"expires_at" db:"expires_at"`
	AutoRenew bool       `json:"auto_renew" db:"auto_renew"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// SubscriptionHistory represents subscription changes and payments
type SubscriptionHistory struct {
	ID               int       `json:"id" db:"id"`
	GamenetID        int       `json:"gamenet_id" db:"gamenet_id"`
	PlanID           int       `json:"plan_id" db:"plan_id"`
	Action           string    `json:"action" db:"action"`
	PreviousPlanID   *int      `json:"previous_plan_id" db:"previous_plan_id"`
	AmountPaid       *float64  `json:"amount_paid" db:"amount_paid"`
	PaymentMethod    *string   `json:"payment_method" db:"payment_method"`
	PaymentReference *string   `json:"payment_reference" db:"payment_reference"`
	Notes            *string   `json:"notes" db:"notes"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// SubscriptionPayment represents a payment transaction
type SubscriptionPayment struct {
	ID               int              `json:"id" db:"id"`
	GamenetID        int              `json:"gamenet_id" db:"gamenet_id"`
	SubscriptionID   int              `json:"subscription_id" db:"subscription_id"`
	PlanID           int              `json:"plan_id" db:"plan_id"`
	Amount           float64          `json:"amount" db:"amount"`
	Currency         string           `json:"currency" db:"currency"`
	PaymentMethod    string           `json:"payment_method" db:"payment_method"`
	PaymentReference string           `json:"payment_reference" db:"payment_reference"`
	Status           string           `json:"status" db:"status"`
	GatewayResponse  *GatewayResponse `json:"gateway_response" db:"gateway_response"`
	ProcessedAt      *time.Time       `json:"processed_at" db:"processed_at"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at" db:"updated_at"`
}

// GatewayResponse represents payment gateway response data
type GatewayResponse struct {
	TransactionID string                 `json:"transaction_id,omitempty"`
	Status        string                 `json:"status,omitempty"`
	Message       string                 `json:"message,omitempty"`
	RawResponse   map[string]interface{} `json:"raw_response,omitempty"`
}

// Value implements the driver.Valuer interface for GatewayResponse
func (gr GatewayResponse) Value() (driver.Value, error) {
	return json.Marshal(gr)
}

// Scan implements the sql.Scanner interface for GatewayResponse
func (gr *GatewayResponse) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, gr)
}

// Request/Response DTOs

// CreatePlanRequest represents a plan creation request
type CreatePlanRequest struct {
	Name                     string   `json:"name" binding:"required"`
	PlanType                 string   `json:"plan_type" binding:"required,oneof=trial monthly annual"`
	Price                    float64  `json:"price" binding:"min=0"`
	AnnualDiscountPercentage *float64 `json:"annual_discount_percentage,omitempty"`
	TrialDurationDays        *int     `json:"trial_duration_days,omitempty"`
	IsActive                 bool     `json:"is_active"`
}

// UpdatePlanRequest represents a plan update request
type UpdatePlanRequest struct {
	Name                     *string  `json:"name"`
	PlanType                 *string  `json:"plan_type,omitempty"`
	Price                    *float64 `json:"price,omitempty"`
	AnnualDiscountPercentage *float64 `json:"annual_discount_percentage,omitempty"`
	TrialDurationDays        *int     `json:"trial_duration_days,omitempty"`
	IsActive                 *bool    `json:"is_active"`
}

// PlanResponse represents a plan response
type PlanResponse struct {
	ID                       int       `json:"id"`
	Name                     string    `json:"name"`
	PlanType                 string    `json:"plan_type"`
	Price                    float64   `json:"price"`
	AnnualDiscountPercentage *float64  `json:"annual_discount_percentage"`
	TrialDurationDays        *int      `json:"trial_duration_days"`
	IsActive                 bool      `json:"is_active"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

// SubscriptionResponse represents a subscription response
type SubscriptionResponse struct {
	ID        int           `json:"id"`
	GamenetID int           `json:"gamenet_id"`
	PlanID    int           `json:"plan_id"`
	Plan      *PlanResponse `json:"plan,omitempty"`
	Status    string        `json:"status"`
	StartedAt time.Time     `json:"started_at"`
	ExpiresAt *time.Time    `json:"expires_at"`
	AutoRenew bool          `json:"auto_renew"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// CreateSubscriptionRequest represents a subscription creation request
type CreateSubscriptionRequest struct {
	GamenetID int  `json:"gamenet_id" binding:"required"`
	PlanID    int  `json:"plan_id" binding:"required"`
	AutoRenew bool `json:"auto_renew"`
}

// UpdateSubscriptionRequest represents a subscription update request
type UpdateSubscriptionRequest struct {
	PlanID    *int  `json:"plan_id"`
	AutoRenew *bool `json:"auto_renew"`
}

// PaymentRequest represents a payment request
type PaymentRequest struct {
	GamenetID        int     `json:"gamenet_id" binding:"required"`
	SubscriptionID   int     `json:"subscription_id" binding:"required"`
	PlanID           int     `json:"plan_id" binding:"required"`
	Amount           float64 `json:"amount" binding:"required,min=0"`
	Currency         string  `json:"currency" binding:"required,len=3"`
	PaymentMethod    string  `json:"payment_method" binding:"required"`
	PaymentReference string  `json:"payment_reference" binding:"required"`
}

// UpdatePaymentRequest represents a payment update request
type UpdatePaymentRequest struct {
	Status          *string          `json:"status" binding:"omitempty,oneof=pending completed failed refunded cancelled"`
	GatewayResponse *GatewayResponse `json:"gateway_response"`
	ProcessedAt     *time.Time       `json:"processed_at"`
}

// Helper methods

// ToResponse converts SubscriptionPlan to PlanResponse
func (sp *SubscriptionPlan) ToResponse() PlanResponse {
	return PlanResponse{
		ID:                       sp.ID,
		Name:                     sp.Name,
		PlanType:                 sp.PlanType,
		Price:                    sp.Price,
		AnnualDiscountPercentage: sp.AnnualDiscountPercentage,
		TrialDurationDays:        sp.TrialDurationDays,
		IsActive:                 sp.IsActive,
		CreatedAt:                sp.CreatedAt,
		UpdatedAt:                sp.UpdatedAt,
	}
}

// ToResponse converts UserSubscription to SubscriptionResponse
func (us *UserSubscription) ToResponse() SubscriptionResponse {
	return SubscriptionResponse{
		ID:        us.ID,
		GamenetID: us.GamenetID,
		PlanID:    us.PlanID,
		Status:    us.Status,
		StartedAt: us.StartedAt,
		ExpiresAt: us.ExpiresAt,
		AutoRenew: us.AutoRenew,
		CreatedAt: us.CreatedAt,
		UpdatedAt: us.UpdatedAt,
	}
}

// IsExpired checks if the subscription is expired
func (us *UserSubscription) IsExpired() bool {
	if us.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*us.ExpiresAt)
}

// IsActive checks if the subscription is currently active
func (us *UserSubscription) IsActive() bool {
	return us.Status == "active" || us.Status == "trial"
}

// GetEffectivePrice calculates the effective price considering discounts
func (sp *SubscriptionPlan) GetEffectivePrice() float64 {
	if sp.PlanType == "annual" && sp.AnnualDiscountPercentage != nil {
		discount := *sp.AnnualDiscountPercentage / 100
		return sp.Price * (1 - discount)
	}
	return sp.Price
}

// GetTrialEndDate calculates when the trial period ends
func (us *UserSubscription) GetTrialEndDate() *time.Time {
	if us.Status != "trial" {
		return nil
	}

	// This would need to be calculated based on the plan's trial duration
	// For now, return nil as we'd need the plan details
	return nil
}
