-- Migration: Convert to single-user schema (remove users and user_id columns)
-- Up: Drop users table and user_id dependencies, adjust unique constraints

-- 1) bank_accounts: drop FK, index, and user_id
ALTER TABLE IF EXISTS bank_accounts DROP CONSTRAINT IF EXISTS bank_accounts_user_id_fkey;
DROP INDEX IF EXISTS idx_bank_accounts_user_id;
ALTER TABLE IF EXISTS bank_accounts DROP COLUMN IF EXISTS user_id;

-- 2) credit_cards: drop FK, index, and user_id
ALTER TABLE IF EXISTS credit_cards DROP CONSTRAINT IF EXISTS credit_cards_user_id_fkey;
DROP INDEX IF EXISTS idx_credit_cards_user_id;
ALTER TABLE IF EXISTS credit_cards DROP COLUMN IF EXISTS user_id;

-- 3) income_sources: drop FK, index, and user_id
ALTER TABLE IF EXISTS income_sources DROP CONSTRAINT IF EXISTS income_sources_user_id_fkey;
DROP INDEX IF EXISTS idx_income_sources_user_id;
ALTER TABLE IF EXISTS income_sources DROP COLUMN IF EXISTS user_id;

-- 4) recurring_payments: drop FK, index, and user_id
ALTER TABLE IF EXISTS recurring_payments DROP CONSTRAINT IF EXISTS recurring_payments_user_id_fkey;
DROP INDEX IF EXISTS idx_recurring_payments_user_id;
ALTER TABLE IF EXISTS recurring_payments DROP COLUMN IF EXISTS user_id;

-- 5) app_settings: drop FK/index, change unique constraint to UNIQUE(key), drop user_id
-- Ensure new uniqueness on key
ALTER TABLE IF EXISTS app_settings ADD CONSTRAINT IF NOT EXISTS app_settings_key_unique UNIQUE (key);
ALTER TABLE IF EXISTS app_settings DROP CONSTRAINT IF EXISTS app_settings_user_id_fkey;
DROP INDEX IF EXISTS idx_app_settings_user_id;
-- Dropping column will remove the old UNIQUE(user_id, key) constraint implicitly
ALTER TABLE IF EXISTS app_settings DROP COLUMN IF EXISTS user_id;

-- 6) drop users triggers and table
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TABLE IF EXISTS users;

-- Deprecated: Consolidated into 000001_initial_schema.up.sql as single-user baseline.
-- Intentionally left blank for fresh deployments.
