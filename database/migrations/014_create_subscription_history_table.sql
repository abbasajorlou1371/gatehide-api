-- version: 014_create_subscription_history_table
-- description: Create subscription history table to track subscription changes and payments

-- UP
CREATE TABLE IF NOT EXISTS subscription_history (
    id INT AUTO_INCREMENT PRIMARY KEY,
    gamenet_id INT NOT NULL,
    plan_id INT NOT NULL,
    action ENUM('created', 'renewed', 'upgraded', 'downgraded', 'cancelled', 'expired', 'grace_period_started', 'grace_period_ended') NOT NULL,
    previous_plan_id INT NULL,
    amount_paid DECIMAL(10,2) NULL,
    payment_method VARCHAR(100) NULL,
    payment_reference VARCHAR(255) NULL,
    notes TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (plan_id) REFERENCES subscription_plans(id) ON DELETE RESTRICT,
    FOREIGN KEY (previous_plan_id) REFERENCES subscription_plans(id) ON DELETE SET NULL,
    FOREIGN KEY (gamenet_id) REFERENCES gamenets(id) ON DELETE CASCADE,
    INDEX idx_gamenet_id (gamenet_id),
    INDEX idx_action (action),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- DOWN
DROP TABLE IF EXISTS subscription_history;
