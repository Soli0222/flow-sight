-- Rollback script for initial schema

DROP TRIGGER IF EXISTS update_app_settings_updated_at ON app_settings;
DROP TRIGGER IF EXISTS update_card_monthly_totals_updated_at ON card_monthly_totals;
DROP TRIGGER IF EXISTS update_recurring_payments_updated_at ON recurring_payments;
DROP TRIGGER IF EXISTS update_monthly_income_records_updated_at ON monthly_income_records;
DROP TRIGGER IF EXISTS update_income_sources_updated_at ON income_sources;
DROP TRIGGER IF EXISTS update_credit_cards_updated_at ON credit_cards;
DROP TRIGGER IF EXISTS update_bank_accounts_updated_at ON bank_accounts;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP INDEX IF EXISTS idx_app_settings_user_id;
DROP INDEX IF EXISTS idx_card_monthly_totals_year_month;
DROP INDEX IF EXISTS idx_card_monthly_totals_credit_card_id;
DROP INDEX IF EXISTS idx_recurring_payments_active;
DROP INDEX IF EXISTS idx_recurring_payments_user_id;
DROP INDEX IF EXISTS idx_monthly_income_records_year_month;
DROP INDEX IF EXISTS idx_monthly_income_records_income_source_id;
DROP INDEX IF EXISTS idx_income_sources_active;
DROP INDEX IF EXISTS idx_income_sources_user_id;
DROP INDEX IF EXISTS idx_credit_cards_user_id;
DROP INDEX IF EXISTS idx_bank_accounts_user_id;

DROP TABLE IF EXISTS app_settings;
DROP TABLE IF EXISTS card_monthly_totals;
DROP TABLE IF EXISTS recurring_payments;
DROP TABLE IF EXISTS monthly_income_records;
DROP TABLE IF EXISTS income_sources;
DROP TABLE IF EXISTS credit_cards;
DROP TABLE IF EXISTS bank_accounts;
DROP TABLE IF EXISTS users;
