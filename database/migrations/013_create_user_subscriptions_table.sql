-- version: 013_create_user_subscriptions_table
-- description: Create user subscriptions table to track current subscription status

-- UP
CREATE TABLE IF NOT EXISTS user_subscriptions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    gamenet_id INT NOT NULL,
    plan_id INT NOT NULL,
    status ENUM('active', 'trial', 'expired', 'cancelled', 'grace_period') NOT NULL DEFAULT 'trial',
    started_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NULL,
    auto_renew BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (plan_id) REFERENCES subscription_plans(id) ON DELETE RESTRICT,
    FOREIGN KEY (gamenet_id) REFERENCES gamenets(id) ON DELETE CASCADE,
    INDEX idx_gamenet_id (gamenet_id),
    INDEX idx_status (status),
    INDEX idx_expires_at (expires_at),
    INDEX idx_created_at (created_at),
    UNIQUE KEY unique_active_subscription (gamenet_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- DOWN
DROP TABLE IF EXISTS user_subscriptions;
