package repositories

import (
	"database/sql"
	"fmt"

	"github.com/gatehide/gatehide-api/internal/models"
)

// SubscriptionPlanRepository handles subscription plan database operations
type SubscriptionPlanRepository struct {
	db *sql.DB
}

// NewSubscriptionPlanRepository creates a new subscription plan repository
func NewSubscriptionPlanRepository(db *sql.DB) *SubscriptionPlanRepository {
	return &SubscriptionPlanRepository{db: db}
}

// Create creates a new subscription plan
func (r *SubscriptionPlanRepository) Create(plan *models.SubscriptionPlan) error {
	query := `
		INSERT INTO subscription_plans (
			name, plan_type, price, annual_discount_percentage, 
			trial_duration_days, is_active
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		plan.Name,
		plan.PlanType,
		plan.Price,
		plan.AnnualDiscountPercentage,
		plan.TrialDurationDays,
		plan.IsActive,
	)

	if err != nil {
		return fmt.Errorf("failed to create subscription plan: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	plan.ID = int(id)
	return nil
}

// GetByID retrieves a subscription plan by ID
func (r *SubscriptionPlanRepository) GetByID(id int) (*models.SubscriptionPlan, error) {
	query := `
		SELECT sp.id, sp.name, sp.plan_type, sp.price, sp.annual_discount_percentage, 
		       sp.trial_duration_days, sp.is_active, sp.created_at, sp.updated_at,
		       COALESCE(COUNT(us.id), 0) as subscription_count
		FROM subscription_plans sp
		LEFT JOIN user_subscriptions us ON sp.id = us.plan_id AND us.status IN ('active', 'trial')
		WHERE sp.id = ?
		GROUP BY sp.id
	`

	plan := &models.SubscriptionPlan{}
	err := r.db.QueryRow(query, id).Scan(
		&plan.ID,
		&plan.Name,
		&plan.PlanType,
		&plan.Price,
		&plan.AnnualDiscountPercentage,
		&plan.TrialDurationDays,
		&plan.IsActive,
		&plan.CreatedAt,
		&plan.UpdatedAt,
		&plan.SubscriptionCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subscription plan not found")
		}
		return nil, fmt.Errorf("failed to get subscription plan: %w", err)
	}

	return plan, nil
}

// GetAll retrieves all subscription plans with optional filters
func (r *SubscriptionPlanRepository) GetAll(limit, offset int, isActive *bool) ([]*models.SubscriptionPlan, error) {
	query := `
		SELECT sp.id, sp.name, sp.plan_type, sp.price, sp.annual_discount_percentage, 
		       sp.trial_duration_days, sp.is_active, sp.created_at, sp.updated_at,
		       COALESCE(COUNT(us.id), 0) as subscription_count
		FROM subscription_plans sp
		LEFT JOIN user_subscriptions us ON sp.id = us.plan_id AND us.status IN ('active', 'trial')
	`
	args := []interface{}{}

	if isActive != nil {
		query += " WHERE sp.is_active = ?"
		args = append(args, *isActive)
	}

	query += " GROUP BY sp.id ORDER BY sp.created_at DESC"

	if limit > 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription plans: %w", err)
	}
	defer rows.Close()

	var plans []*models.SubscriptionPlan
	for rows.Next() {
		plan := &models.SubscriptionPlan{}
		err := rows.Scan(
			&plan.ID,
			&plan.Name,
			&plan.PlanType,
			&plan.Price,
			&plan.AnnualDiscountPercentage,
			&plan.TrialDurationDays,
			&plan.IsActive,
			&plan.CreatedAt,
			&plan.UpdatedAt,
			&plan.SubscriptionCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription plan: %w", err)
		}
		plans = append(plans, plan)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating subscription plans: %w", err)
	}

	return plans, nil
}

// Update updates an existing subscription plan
func (r *SubscriptionPlanRepository) Update(id int, plan *models.SubscriptionPlan) error {
	query := `
		UPDATE subscription_plans 
		SET name = ?, plan_type = ?, price = ?, annual_discount_percentage = ?, 
		    trial_duration_days = ?, is_active = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := r.db.Exec(query,
		plan.Name,
		plan.PlanType,
		plan.Price,
		plan.AnnualDiscountPercentage,
		plan.TrialDurationDays,
		plan.IsActive,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update subscription plan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription plan not found")
	}

	return nil
}

// Delete deletes a subscription plan
func (r *SubscriptionPlanRepository) Delete(id int) error {
	query := `DELETE FROM subscription_plans WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription plan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription plan not found")
	}

	return nil
}

// Count returns the total number of subscription plans
func (r *SubscriptionPlanRepository) Count(isActive *bool) (int, error) {
	query := `SELECT COUNT(*) FROM subscription_plans`
	args := []interface{}{}

	if isActive != nil {
		query += " WHERE is_active = ?"
		args = append(args, *isActive)
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count subscription plans: %w", err)
	}

	return count, nil
}

// HasActiveSubscriptions checks if a plan has any active subscriptions
func (r *SubscriptionPlanRepository) HasActiveSubscriptions(planID int) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM user_subscriptions 
		WHERE plan_id = ? AND status IN ('active', 'trial')
	`

	var count int
	err := r.db.QueryRow(query, planID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check active subscriptions: %w", err)
	}

	return count > 0, nil
}
