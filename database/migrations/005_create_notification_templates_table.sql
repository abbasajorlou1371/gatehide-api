-- version: 005_create_notification_templates_table
-- description: Create notification_templates table

-- UP
CREATE TABLE IF NOT EXISTS notification_templates (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type ENUM('email', 'sms', 'database') NOT NULL,
    subject VARCHAR(500),
    content TEXT,
    html_content TEXT,
    variables JSON,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY unique_template_name_type (name, type),
    INDEX idx_templates_type (type),
    INDEX idx_templates_is_active (is_active),
    INDEX idx_templates_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- DOWN
DROP TABLE IF EXISTS notification_templates;
