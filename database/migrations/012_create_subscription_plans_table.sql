-- version: 012_create_subscription_plans_table
-- description: Create subscription plans table for managing different subscription tiers

-- UP
CREATE TABLE IF NOT EXISTS subscription_plans (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    plan_type ENUM('trial', 'monthly', 'annual') NOT NULL,
    price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    annual_discount_percentage DECIMAL(5,2) NULL DEFAULT 0.00,
    trial_duration_days INT NULL DEFAULT 30,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_plan_type (plan_type),
    INDEX idx_is_active (is_active),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- DOWN
DROP TABLE IF EXISTS subscription_plans;
