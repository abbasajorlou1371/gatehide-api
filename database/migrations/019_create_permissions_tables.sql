-- version: 019_create_permissions_tables
-- description: Create permissions, roles, and role_permissions tables for RBAC system

-- UP
-- Create permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_resource_action (resource, action),
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create role_permissions junction table
CREATE TABLE IF NOT EXISTS role_permissions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    role_id INT NOT NULL,
    permission_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    
    UNIQUE KEY unique_role_permission (role_id, permission_id),
    INDEX idx_role_id (role_id),
    INDEX idx_permission_id (permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert roles
INSERT INTO roles (name, description) VALUES
('administrator', 'System administrator with full access'),
('gamenet', 'Gaming center operator with limited access'),
('user', 'Regular user with basic access');

-- Insert permissions
INSERT INTO permissions (name, description, resource, action) VALUES
-- Dashboard permissions
('dashboard:view', 'View dashboard', 'dashboard', 'view'),

-- Gamenet permissions
('gamenets:create', 'Create gamenets', 'gamenets', 'create'),
('gamenets:read', 'View gamenets', 'gamenets', 'read'),
('gamenets:update', 'Update gamenets', 'gamenets', 'update'),
('gamenets:delete', 'Delete gamenets', 'gamenets', 'delete'),

-- User permissions
('users:create', 'Create users', 'users', 'create'),
('users:read', 'View users', 'users', 'read'),
('users:update', 'Update users', 'users', 'update'),
('users:delete', 'Delete users', 'users', 'delete'),

-- Subscription plan permissions
('subscription_plans:create', 'Create subscription plans', 'subscription_plans', 'create'),
('subscription_plans:read', 'View subscription plans', 'subscription_plans', 'read'),
('subscription_plans:update', 'Update subscription plans', 'subscription_plans', 'update'),
('subscription_plans:delete', 'Delete subscription plans', 'subscription_plans', 'delete'),

-- Analytics permissions
('analytics:view', 'View analytics', 'analytics', 'view'),

-- Payment permissions
('payments:view', 'View payments', 'payments', 'view'),

-- Transaction permissions
('transactions:view', 'View transactions', 'transactions', 'view'),

-- Invoice permissions
('invoices:view', 'View invoices', 'invoices', 'view'),

-- Settings permissions
('settings:manage', 'Manage settings', 'settings', 'manage'),

-- Support permissions
('support:access', 'Access support', 'support', 'access'),

-- Reservation permissions
('reservation:manage', 'Manage reservations', 'reservation', 'manage'),

-- Wallet permissions
('wallet:view', 'View wallet', 'wallet', 'view');

-- Assign permissions to administrator role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'administrator'
AND p.name IN (
    'dashboard:view',
    'gamenets:create', 'gamenets:read', 'gamenets:update', 'gamenets:delete',
    'subscription_plans:create', 'subscription_plans:read', 'subscription_plans:update', 'subscription_plans:delete',
    'analytics:view',
    'payments:view',
    'transactions:view',
    'invoices:view',
    'settings:manage',
    'support:access'
);

-- Assign permissions to gamenet role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'gamenet'
AND p.name IN (
    'dashboard:view',
    'users:create', 'users:read', 'users:update', 'users:delete',
    'analytics:view',
    'transactions:view',
    'payments:view',
    'support:access',
    'settings:manage'
);

-- Assign permissions to user role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'user'
AND p.name IN (
    'reservation:manage',
    'support:access',
    'settings:manage',
    'wallet:view'
);

-- DOWN
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
