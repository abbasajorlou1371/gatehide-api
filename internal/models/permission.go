package models

import "time"

// Permission represents a permission in the system
type Permission struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Resource    string    `json:"resource" db:"resource"`
	Action      string    `json:"action" db:"action"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Role represents a role in the system
type Role struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// RoleWithPermissions represents a role with its associated permissions
type RoleWithPermissions struct {
	Role        Role         `json:"role"`
	Permissions []Permission `json:"permissions"`
}

// PermissionString returns a formatted permission string (resource:action)
func (p *Permission) PermissionString() string {
	return p.Resource + ":" + p.Action
}

// HasPermission checks if a role has a specific permission
func (r *RoleWithPermissions) HasPermission(resource, action string) bool {
	for _, perm := range r.Permissions {
		if perm.Resource == resource && perm.Action == action {
			return true
		}
	}
	return false
}

// HasAnyPermission checks if a role has any permission for a resource
func (r *RoleWithPermissions) HasAnyPermission(resource string) bool {
	for _, perm := range r.Permissions {
		if perm.Resource == resource {
			return true
		}
	}
	return false
}

// GetPermissionStrings returns all permissions as strings
func (r *RoleWithPermissions) GetPermissionStrings() []string {
	var permissions []string
	for _, perm := range r.Permissions {
		permissions = append(permissions, perm.PermissionString())
	}
	return permissions
}

// CanAccessResource checks if a role can access a specific resource
func (r *RoleWithPermissions) CanAccessResource(resource string) bool {
	return r.HasAnyPermission(resource)
}

// Permission constants for easy reference
const (
	// Dashboard permissions
	PermissionDashboardView = "dashboard:view"

	// Gamenet permissions
	PermissionGamenetsCreate = "gamenets:create"
	PermissionGamenetsRead   = "gamenets:read"
	PermissionGamenetsUpdate = "gamenets:update"
	PermissionGamenetsDelete = "gamenets:delete"

	// User permissions
	PermissionUsersCreate = "users:create"
	PermissionUsersRead   = "users:read"
	PermissionUsersUpdate = "users:update"
	PermissionUsersDelete = "users:delete"

	// Subscription plan permissions
	PermissionSubscriptionPlansCreate = "subscription_plans:create"
	PermissionSubscriptionPlansRead   = "subscription_plans:read"
	PermissionSubscriptionPlansUpdate = "subscription_plans:update"
	PermissionSubscriptionPlansDelete = "subscription_plans:delete"

	// Analytics permissions
	PermissionAnalyticsView = "analytics:view"

	// Payment permissions
	PermissionPaymentsView = "payments:view"

	// Transaction permissions
	PermissionTransactionsView = "transactions:view"

	// Invoice permissions
	PermissionInvoicesView = "invoices:view"

	// Settings permissions
	PermissionSettingsManage = "settings:manage"

	// Support permissions
	PermissionSupportAccess = "support:access"

	// Reservation permissions
	PermissionReservationManage = "reservation:manage"

	// Wallet permissions
	PermissionWalletView = "wallet:view"
)

// Role constants
const (
	RoleAdministrator = "administrator"
	RoleGamenet       = "gamenet"
	RoleUser          = "user"
)
