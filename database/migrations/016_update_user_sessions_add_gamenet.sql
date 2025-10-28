-- version: 016_update_user_sessions_add_gamenet
-- description: Update user_sessions table to support gamenet user type

-- UP
-- Modify the user_type column to include gamenet
-- Note: The foreign key was already removed in migration 009
ALTER TABLE user_sessions MODIFY COLUMN user_type ENUM('user', 'admin', 'gamenet') NOT NULL;

-- DOWN
-- Modify the user_type column back to original values
-- Note: Do not re-add foreign key as it was removed in migration 009
-- First, update any 'gamenet' records to 'admin' to avoid data truncation
UPDATE user_sessions SET user_type = 'admin' WHERE user_type = 'gamenet';
-- Then modify the column to only allow 'user' and 'admin'
ALTER TABLE user_sessions MODIFY COLUMN user_type ENUM('user', 'admin') NOT NULL;

