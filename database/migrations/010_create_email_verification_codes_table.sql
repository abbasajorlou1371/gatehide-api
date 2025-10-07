-- version: 010_create_email_verification_codes_table
-- description: Create table to store email verification codes with expiration

-- UP
CREATE TABLE IF NOT EXISTS email_verification_codes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    user_type ENUM('user', 'admin') NOT NULL,
    email VARCHAR(255) NOT NULL,
    code VARCHAR(64) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_email (user_id, user_type),
    INDEX idx_code (code),
    INDEX idx_expires_at (expires_at)
);

-- DOWN
DROP TABLE IF EXISTS email_verification_codes;
