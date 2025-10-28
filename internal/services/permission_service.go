package services

import (
	"database/sql"
	"fmt"

	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/repositories"
)

// PermissionServiceInterface defines the interface for permission business logic
type PermissionServiceInterface interface {
	CheckPermission(userType, resource, action string) error
	CheckUserPermission(userID int, userType, resource, action string) error
	GetUserPermissions(userType string) ([]string, error)
	GetUserPermissionsByID(userID int, userType string) ([]string, error)
	CanAccessResource(userType string, resourceType string, resourceID int, userID int) (bool, error)
	GetRoleWithPermissions(roleType string) (*models.RoleWithPermissions, error)
	HasPermission(userType, resource, action string) (bool, error)
}

// PermissionService handles permission business logic
type PermissionService struct {
	permissionRepo *repositories.PermissionRepository
	db             *sql.DB
}

// NewPermissionService creates a new permission service
func NewPermissionService(permissionRepo *repositories.PermissionRepository, db *sql.DB) *PermissionService {
	return &PermissionService{
		permissionRepo: permissionRepo,
		db:             db,
	}
}

// CheckPermission checks if a user type has a specific permission
func (s *PermissionService) CheckPermission(userType, resource, action string) error {
	// Map user types to role names
	roleName := s.mapUserTypeToRoleName(userType)

	hasPermission, err := s.permissionRepo.HasPermission(roleName, resource, action)
	if err != nil {
		return fmt.Errorf("failed to check permission: %w", err)
	}

	if !hasPermission {
		return fmt.Errorf("permission denied: %s:%s for role %s", resource, action, userType)
	}

	return nil
}

// CheckUserPermission checks if a specific user has a specific permission
func (s *PermissionService) CheckUserPermission(userID int, userType, resource, action string) error {
	hasPermission, err := s.permissionRepo.HasUserPermission(userID, userType, resource, action)
	if err != nil {
		return fmt.Errorf("failed to check user permission: %w", err)
	}

	if !hasPermission {
		return fmt.Errorf("permission denied: %s:%s for user %d", resource, action, userID)
	}

	return nil
}

// GetUserPermissions retrieves all permissions for a user type
func (s *PermissionService) GetUserPermissions(userType string) ([]string, error) {
	// Map user types to role names
	roleName := s.mapUserTypeToRoleName(userType)

	permissions, err := s.permissionRepo.GetPermissionsByRole(roleName)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	var permissionStrings []string
	for _, perm := range permissions {
		permissionStrings = append(permissionStrings, perm.PermissionString())
	}

	return permissionStrings, nil
}

// GetUserPermissionsByID retrieves all permissions for a specific user
func (s *PermissionService) GetUserPermissionsByID(userID int, userType string) ([]string, error) {
	permissions, err := s.permissionRepo.GetUserPermissions(userID, userType)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	var permissionStrings []string
	for _, perm := range permissions {
		permissionStrings = append(permissionStrings, perm.PermissionString())
	}

	return permissionStrings, nil
}

// mapUserTypeToRoleName maps user types to role names
func (s *PermissionService) mapUserTypeToRoleName(userType string) string {
	switch userType {
	case "admin":
		return "administrator"
	case "user":
		return "user"
	case "gamenet":
		return "gamenet"
	default:
		return userType // fallback to the original value
	}
}

// CanAccessResource checks if a user can access a specific resource with ownership validation
func (s *PermissionService) CanAccessResource(userType string, resourceType string, resourceID int, userID int) (bool, error) {
	// Check if user has administrator role
	hasAdminRole, err := s.permissionRepo.HasUserRole(userID, userType, "administrator")
	if err != nil {
		return false, fmt.Errorf("failed to check administrator role: %w", err)
	}

	// Administrators can access all resources
	if hasAdminRole {
		return true, nil
	}

	// Check if user has permission for the resource type
	roleName := s.mapUserTypeToRoleName(userType)
	hasPermission, err := s.permissionRepo.HasPermission(roleName, resourceType, "read")
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}

	if !hasPermission {
		return false, nil
	}

	// For gamenets accessing users, check ownership via users_gamenets table
	if userType == models.RoleGamenet && resourceType == "users" {
		return s.checkGamenetUserOwnership(userID, resourceID)
	}

	// For users accessing their own profile/settings
	if userType == models.RoleUser && (resourceType == "profile" || resourceType == "settings") {
		return userID == resourceID, nil
	}

	// Default: if they have permission, they can access
	return true, nil
}

// checkGamenetUserOwnership checks if a gamenet owns/manages a specific user
func (s *PermissionService) checkGamenetUserOwnership(gamenetID, userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM users_gamenets WHERE gamenet_id = ? AND user_id = ?`
	var count int
	err := s.db.QueryRow(query, gamenetID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check gamenet-user ownership: %w", err)
	}

	return count > 0, nil
}

// GetRoleWithPermissions retrieves a role with all its permissions
func (s *PermissionService) GetRoleWithPermissions(roleType string) (*models.RoleWithPermissions, error) {
	return s.permissionRepo.GetRoleWithPermissions(roleType)
}

// HasPermission checks if a user type has a specific permission
func (s *PermissionService) HasPermission(userType, resource, action string) (bool, error) {
	// Map user types to role names
	roleName := s.mapUserTypeToRoleName(userType)

	return s.permissionRepo.HasPermission(roleName, resource, action)
}

// CheckResourceOwnership is a helper method to check if a user owns a resource
func (s *PermissionService) CheckResourceOwnership(userType string, userID int, resourceType string, resourceID int) (bool, error) {
	// Check if user has administrator role
	hasAdminRole, err := s.permissionRepo.HasUserRole(userID, userType, "administrator")
	if err != nil {
		return false, fmt.Errorf("failed to check administrator role: %w", err)
	}

	// Administrators can access all resources
	if hasAdminRole {
		return true, nil
	}

	switch resourceType {
	case "users":
		if userType == models.RoleGamenet {
			return s.checkGamenetUserOwnership(userID, resourceID)
		}
		if userType == models.RoleUser {
			return userID == resourceID, nil
		}
		return true, nil // Admin can access all
	case "profile", "settings":
		return userID == resourceID, nil
	default:
		return true, nil // Default to allowing access
	}
}

// ValidateUserAccess validates if a user can perform an action on a resource
func (s *PermissionService) ValidateUserAccess(userType string, userID int, resourceType string, resourceID int, action string) error {
	// Check permission first
	err := s.CheckPermission(userType, resourceType, action)
	if err != nil {
		return err
	}

	// Check resource ownership
	canAccess, err := s.CanAccessResource(userType, resourceType, resourceID, userID)
	if err != nil {
		return fmt.Errorf("failed to check resource access: %w", err)
	}

	if !canAccess {
		return fmt.Errorf("access denied: insufficient ownership for resource %s:%d", resourceType, resourceID)
	}

	return nil
}
