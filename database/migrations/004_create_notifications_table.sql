-- version: 004_create_notifications_table
-- description: Create notifications table

-- UP
CREATE TABLE IF NOT EXISTS notifications (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type ENUM('email', 'sms', 'database') NOT NULL,
    status ENUM('pending', 'sent', 'failed', 'cancelled') NOT NULL DEFAULT 'pending',
    priority ENUM('low', 'normal', 'high', 'urgent') NOT NULL DEFAULT 'normal',
    recipient VARCHAR(255) NOT NULL,
    subject VARCHAR(500),
    content TEXT,
    template_id INT,
    template_data JSON,
    metadata JSON,
    scheduled_at TIMESTAMP NULL,
    sent_at TIMESTAMP NULL,
    error_msg TEXT,
    retry_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_notifications_type (type),
    INDEX idx_notifications_status (status),
    INDEX idx_notifications_recipient (recipient),
    INDEX idx_notifications_priority (priority),
    INDEX idx_notifications_created_at (created_at),
    INDEX idx_notifications_scheduled_at (scheduled_at),
    INDEX idx_notifications_template_id (template_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- DOWN
DROP TABLE IF EXISTS notifications;
