-- version: 007_create_gamenets_table
-- description: Create gamenets table for managing gaming centers

-- UP
CREATE TABLE IF NOT EXISTS gamenets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner_name VARCHAR(255) NOT NULL,
    owner_mobile VARCHAR(20) NOT NULL,
    address TEXT NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    license_attachment VARCHAR(500) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_name (name),
    INDEX idx_owner_mobile (owner_mobile),
    INDEX idx_email (email),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- DOWN
DROP TABLE IF EXISTS gamenets;
