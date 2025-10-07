-- version: 011_update_email_verification_codes_code_field
-- description: Update email_verification_codes table to support hashed codes (64 characters)

-- UP
ALTER TABLE email_verification_codes MODIFY COLUMN code VARCHAR(64) NOT NULL;

-- DOWN
ALTER TABLE email_verification_codes MODIFY COLUMN code VARCHAR(10) NOT NULL;
