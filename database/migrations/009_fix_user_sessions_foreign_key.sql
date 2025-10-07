-- version: 009_fix_user_sessions_foreign_key
-- description: Fix foreign key constraint in user_sessions table to support both users and admins

-- UP
-- Remove the foreign key constraint that only references users table
ALTER TABLE user_sessions DROP FOREIGN KEY user_sessions_ibfk_1;

-- DOWN
-- Add back the foreign key constraint (this might fail if there are admin sessions)
-- ALTER TABLE user_sessions ADD CONSTRAINT user_sessions_ibfk_1 FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
