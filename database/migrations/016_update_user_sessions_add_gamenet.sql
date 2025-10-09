-- version: 016_update_user_sessions_add_gamenet
-- description: Update user_sessions table to support gamenet user type

-- UP
ALTER TABLE user_sessions 
MODIFY COLUMN user_type ENUM('user', 'admin', 'gamenet') NOT NULL;

-- Remove the foreign key constraint since user_id can now reference multiple tables
ALTER TABLE user_sessions 
DROP FOREIGN KEY user_sessions_ibfk_1;

-- DOWN
ALTER TABLE user_sessions 
MODIFY COLUMN user_type ENUM('user', 'admin') NOT NULL;

-- Re-add the foreign key constraint
ALTER TABLE user_sessions 
ADD CONSTRAINT user_sessions_ibfk_1 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

