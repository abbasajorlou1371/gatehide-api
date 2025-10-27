-- version: 018_add_wallet_columns_to_users
-- description: Add balance and debt columns to users table for wallet functionality

-- UP
ALTER TABLE users ADD COLUMN balance DECIMAL(10, 2) DEFAULT 0.00 NOT NULL AFTER image;
ALTER TABLE users ADD COLUMN debt DECIMAL(10, 2) DEFAULT 0.00 NOT NULL AFTER balance;

-- DOWN
ALTER TABLE users DROP COLUMN balance;
ALTER TABLE users DROP COLUMN debt;

