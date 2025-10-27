-- version: 017_create_users_gamenets_table
-- description: Create users_gamenets junction table for many-to-many relationship between users and gamenets

-- UP
CREATE TABLE IF NOT EXISTS users_gamenets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    gamenet_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Foreign key constraints
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (gamenet_id) REFERENCES gamenets(id) ON DELETE CASCADE,
    
    -- Unique constraint to prevent duplicate relationships
    UNIQUE KEY unique_user_gamenet (user_id, gamenet_id),
    
    -- Indexes for better query performance
    INDEX idx_user_id (user_id),
    INDEX idx_gamenet_id (gamenet_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- DOWN
DROP TABLE IF EXISTS users_gamenets;

