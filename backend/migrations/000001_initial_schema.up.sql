-- Initial schema (single-user baseline)
-- No authentication, no users table, no user_id columns.

-- Enable extension for gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Bank accounts table
CREATE TABLE IF NOT EXISTS bank_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0, -- Amount in cents
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Credit cards table
CREATE TABLE IF NOT EXISTS credit_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    closing_day INTEGER, -- Closing day of the month (1-31)
    payment_day INTEGER NOT NULL, -- Payment day (1-31)
    bank_account UUID NOT NULL REFERENCES bank_accounts(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_closing_day CHECK (closing_day IS NULL OR (closing_day >= 1 AND closing_day <= 31)),
    CONSTRAINT check_payment_day CHECK (payment_day >= 1 AND payment_day <= 31)
);

-- Income sources table
CREATE TABLE IF NOT EXISTS income_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    income_type VARCHAR(50) NOT NULL CHECK (income_type IN ('monthly_fixed', 'one_time')),
    base_amount BIGINT NOT NULL, -- Amount in cents
    bank_account UUID NOT NULL REFERENCES bank_accounts(id) ON DELETE CASCADE,
    scheduled_year_month VARCHAR(7), -- Format: "2024-01" for one-time income
    payment_day INTEGER, -- Payment day for monthly_fixed income (1-31)
    scheduled_date DATE, -- Scheduled date for one_time income
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_payment_day CHECK (payment_day IS NULL OR (payment_day >= 1 AND payment_day <= 31))
);

-- Monthly income records table
CREATE TABLE IF NOT EXISTS monthly_income_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    income_source_id UUID NOT NULL REFERENCES income_sources(id) ON DELETE CASCADE,
    year_month VARCHAR(7) NOT NULL, -- Format: "2024-01"
    actual_amount BIGINT NOT NULL, -- Amount in cents
    is_confirmed BOOLEAN NOT NULL DEFAULT false,
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(income_source_id, year_month)
);

-- Recurring payments table
CREATE TABLE IF NOT EXISTS recurring_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    amount BIGINT NOT NULL, -- Amount in cents
    payment_day INTEGER NOT NULL, -- Payment day (1-31)
    start_year_month VARCHAR(7) NOT NULL, -- Format: "2024-01"
    total_payments INTEGER, -- For loans
    remaining_payments INTEGER, -- For loans
    bank_account UUID NOT NULL REFERENCES bank_accounts(id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_payment_day CHECK (payment_day >= 1 AND payment_day <= 31),
    CONSTRAINT check_total_payments CHECK (total_payments IS NULL OR total_payments > 0),
    CONSTRAINT check_remaining_payments CHECK (remaining_payments IS NULL OR remaining_payments >= 0)
);

-- Card monthly totals table
CREATE TABLE IF NOT EXISTS card_monthly_totals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    credit_card_id UUID NOT NULL REFERENCES credit_cards(id) ON DELETE CASCADE,
    year_month VARCHAR(7) NOT NULL, -- Format: "2024-01"
    total_amount BIGINT NOT NULL, -- Amount in cents
    is_confirmed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(credit_card_id, year_month)
);

-- App settings table
CREATE TABLE IF NOT EXISTS app_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(255) NOT NULL,
    value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT app_settings_key_unique UNIQUE(key)
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_income_sources_active ON income_sources(is_active);
CREATE INDEX IF NOT EXISTS idx_monthly_income_records_income_source_id ON monthly_income_records(income_source_id);
CREATE INDEX IF NOT EXISTS idx_monthly_income_records_year_month ON monthly_income_records(year_month);
CREATE INDEX IF NOT EXISTS idx_recurring_payments_active ON recurring_payments(is_active);
CREATE INDEX IF NOT EXISTS idx_card_monthly_totals_credit_card_id ON card_monthly_totals(credit_card_id);
CREATE INDEX IF NOT EXISTS idx_card_monthly_totals_year_month ON card_monthly_totals(year_month);

-- Update triggers for updated_at columns
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_bank_accounts_updated_at BEFORE UPDATE ON bank_accounts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_credit_cards_updated_at BEFORE UPDATE ON credit_cards
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_income_sources_updated_at BEFORE UPDATE ON income_sources
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_monthly_income_records_updated_at BEFORE UPDATE ON monthly_income_records
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_recurring_payments_updated_at BEFORE UPDATE ON recurring_payments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_card_monthly_totals_updated_at BEFORE UPDATE ON card_monthly_totals
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_app_settings_updated_at BEFORE UPDATE ON app_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
