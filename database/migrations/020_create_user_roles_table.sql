-- version: 020_create_user_roles_table
-- description: Create user_roles table to assign roles to users, admins, and gamenets

-- UP
-- Create user_roles table to assign roles to different user types
CREATE TABLE IF NOT EXISTS user_roles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    user_type ENUM('user', 'admin', 'gamenet') NOT NULL,
    role_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    
    UNIQUE KEY unique_user_role (user_id, user_type, role_id),
    INDEX idx_user_id_type (user_id, user_type),
    INDEX idx_role_id (role_id),
    INDEX idx_user_type (user_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Assign roles to existing records
-- Assign administrator role to all existing admins
INSERT INTO user_roles (user_id, user_type, role_id)
SELECT a.id, 'admin', r.id
FROM admins a, roles r
WHERE r.name = 'administrator'
AND NOT EXISTS (
    SELECT 1 FROM user_roles ur 
    WHERE ur.user_id = a.id 
    AND ur.user_type = 'admin' 
    AND ur.role_id = r.id
);

-- Assign user role to all existing users
INSERT INTO user_roles (user_id, user_type, role_id)
SELECT u.id, 'user', r.id
FROM users u, roles r
WHERE r.name = 'user'
AND NOT EXISTS (
    SELECT 1 FROM user_roles ur 
    WHERE ur.user_id = u.id 
    AND ur.user_type = 'user' 
    AND ur.role_id = r.id
);

-- Assign gamenet role to all existing gamenets
INSERT INTO user_roles (user_id, user_type, role_id)
SELECT g.id, 'gamenet', r.id
FROM gamenets g, roles r
WHERE r.name = 'gamenet'
AND NOT EXISTS (
    SELECT 1 FROM user_roles ur 
    WHERE ur.user_id = g.id 
    AND ur.user_type = 'gamenet' 
    AND ur.role_id = r.id
);

-- DOWN
DROP TABLE IF EXISTS user_roles;
