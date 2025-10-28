package repositories

import (
	"github.com/gatehide/gatehide-api/internal/models"
)

// PermissionRepositoryInterface defines the interface for permission repository operations
type PermissionRepositoryInterface interface {
	GetPermissionsByRole(roleType string) ([]models.Permission, error)
	HasPermission(roleType, resource, action string) (bool, error)
	GetRoleWithPermissions(roleType string) (*models.RoleWithPermissions, error)
	GetRoleByName(roleName string) (*models.Role, error)
	GetAllRoles() ([]models.Role, error)
	GetAllPermissions() ([]models.Permission, error)
	AssignRoleToUser(userID int, userType string, roleName string) error
	GetUserRoles(userID int, userType string) ([]models.Role, error)
	GetUserPermissions(userID int, userType string) ([]models.Permission, error)
	RemoveRoleFromUser(userID int, userType string, roleName string) error
	HasUserRole(userID int, userType string, roleName string) (bool, error)
}
