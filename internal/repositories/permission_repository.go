package repositories

import (
	"database/sql"
	"fmt"

	"github.com/gatehide/gatehide-api/internal/models"
)

// PermissionRepository handles permission-related database operations
type PermissionRepository struct {
	db *sql.DB
}

// NewPermissionRepository creates a new permission repository
func NewPermissionRepository(db *sql.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// GetPermissionsByRole retrieves all permissions for a specific role
func (r *PermissionRepository) GetPermissionsByRole(roleType string) ([]models.Permission, error) {
	query := `
		SELECT p.id, p.name, p.description, p.resource, p.action, p.created_at, p.updated_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN roles r ON rp.role_id = r.id
		WHERE r.name = ?
		ORDER BY p.resource, p.action
	`

	rows, err := r.db.Query(query, roleType)
	if err != nil {
		return nil, fmt.Errorf("failed to query permissions: %w", err)
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var perm models.Permission
		err := rows.Scan(
			&perm.ID,
			&perm.Name,
			&perm.Description,
			&perm.Resource,
			&perm.Action,
			&perm.CreatedAt,
			&perm.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, perm)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating permissions: %w", err)
	}

	return permissions, nil
}

// HasPermission checks if a role has a specific permission
func (r *PermissionRepository) HasPermission(roleType, resource, action string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN roles r ON rp.role_id = r.id
		WHERE r.name = ? AND p.resource = ? AND p.action = ?
	`

	var count int
	err := r.db.QueryRow(query, roleType, resource, action).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}

	return count > 0, nil
}

// HasUserPermission checks if a specific user has a specific permission
func (r *PermissionRepository) HasUserPermission(userID int, userType, resource, action string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN roles r ON rp.role_id = r.id
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = ? AND ur.user_type = ? AND p.resource = ? AND p.action = ?
	`

	var count int
	err := r.db.QueryRow(query, userID, userType, resource, action).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user permission: %w", err)
	}

	return count > 0, nil
}

// GetRoleWithPermissions retrieves a role with all its permissions
func (r *PermissionRepository) GetRoleWithPermissions(roleType string) (*models.RoleWithPermissions, error) {
	// First get the role
	roleQuery := `SELECT id, name, description, created_at, updated_at FROM roles WHERE name = ?`
	var role models.Role
	err := r.db.QueryRow(roleQuery, roleType).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found: %s", roleType)
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	// Then get permissions for this role
	permissions, err := r.GetPermissionsByRole(roleType)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	return &models.RoleWithPermissions{
		Role:        role,
		Permissions: permissions,
	}, nil
}

// GetRoleByName retrieves a role by name
func (r *PermissionRepository) GetRoleByName(roleName string) (*models.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE name = ?`
	var role models.Role
	err := r.db.QueryRow(query, roleName).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found: %s", roleName)
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &role, nil
}

// GetAllRoles retrieves all roles
func (r *PermissionRepository) GetAllRoles() ([]models.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles ORDER BY name`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %w", err)
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating roles: %w", err)
	}

	return roles, nil
}

// GetAllPermissions retrieves all permissions
func (r *PermissionRepository) GetAllPermissions() ([]models.Permission, error) {
	query := `SELECT id, name, description, resource, action, created_at, updated_at FROM permissions ORDER BY resource, action`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query permissions: %w", err)
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var perm models.Permission
		err := rows.Scan(
			&perm.ID,
			&perm.Name,
			&perm.Description,
			&perm.Resource,
			&perm.Action,
			&perm.CreatedAt,
			&perm.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, perm)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating permissions: %w", err)
	}

	return permissions, nil
}

// AssignRoleToUser assigns a role to a user (user, admin, or gamenet)
func (r *PermissionRepository) AssignRoleToUser(userID int, userType string, roleName string) error {
	// First get the role ID
	role, err := r.GetRoleByName(roleName)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	// Insert the role assignment
	query := `
		INSERT INTO user_roles (user_id, user_type, role_id, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE updated_at = NOW()
	`

	_, err = r.db.Exec(query, userID, userType, role.ID)
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}

// GetUserRoles retrieves all roles assigned to a specific user
func (r *PermissionRepository) GetUserRoles(userID int, userType string) ([]models.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = ? AND ur.user_type = ?
		ORDER BY r.name
	`

	rows, err := r.db.Query(query, userID, userType)
	if err != nil {
		return nil, fmt.Errorf("failed to query user roles: %w", err)
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating roles: %w", err)
	}

	return roles, nil
}

// GetUserPermissions retrieves all permissions for a specific user based on their roles
func (r *PermissionRepository) GetUserPermissions(userID int, userType string) ([]models.Permission, error) {
	query := `
		SELECT DISTINCT p.id, p.name, p.description, p.resource, p.action, p.created_at, p.updated_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = ? AND ur.user_type = ?
		ORDER BY p.resource, p.action
	`

	rows, err := r.db.Query(query, userID, userType)
	if err != nil {
		return nil, fmt.Errorf("failed to query user permissions: %w", err)
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var perm models.Permission
		err := rows.Scan(
			&perm.ID,
			&perm.Name,
			&perm.Description,
			&perm.Resource,
			&perm.Action,
			&perm.CreatedAt,
			&perm.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, perm)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating permissions: %w", err)
	}

	return permissions, nil
}

// RemoveRoleFromUser removes a role from a user
func (r *PermissionRepository) RemoveRoleFromUser(userID int, userType string, roleName string) error {
	// First get the role ID
	role, err := r.GetRoleByName(roleName)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	// Delete the role assignment
	query := `DELETE FROM user_roles WHERE user_id = ? AND user_type = ? AND role_id = ?`

	result, err := r.db.Exec(query, userID, userType, role.ID)
	if err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role assignment not found")
	}

	return nil
}

// HasUserRole checks if a user has a specific role
func (r *PermissionRepository) HasUserRole(userID int, userType string, roleName string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM user_roles ur
		INNER JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = ? AND ur.user_type = ? AND r.name = ?
	`

	var count int
	err := r.db.QueryRow(query, userID, userType, roleName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user role: %w", err)
	}

	return count > 0, nil
}
