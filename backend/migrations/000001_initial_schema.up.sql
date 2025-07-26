-- Initial schema for Flow Sight financial management application

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bank accounts table
CREATE TABLE IF NOT EXISTS bank_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0, -- Amount in cents
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Assets table (credit cards and loans)
CREATE TABLE IF NOT EXISTS assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    asset_type VARCHAR(50) NOT NULL CHECK (asset_type IN ('card', 'loan')),
    closing_day INTEGER, -- For credit cards (1-31)
    payment_day INTEGER NOT NULL, -- Payment day (1-31)
    bank_account UUID NOT NULL REFERENCES bank_accounts(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_closing_day CHECK (closing_day IS NULL OR (closing_day >= 1 AND closing_day <= 31)),
    CONSTRAINT check_payment_day CHECK (payment_day >= 1 AND payment_day <= 31)
);

-- Income sources table
CREATE TABLE IF NOT EXISTS income_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    income_type VARCHAR(50) NOT NULL CHECK (income_type IN ('monthly_fixed', 'one_time')),
    base_amount BIGINT NOT NULL, -- Amount in cents
    bank_account UUID NOT NULL REFERENCES bank_accounts(id),
    scheduled_year_month VARCHAR(7), -- Format: "2024-01" for one-time income
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    amount BIGINT NOT NULL, -- Amount in cents
    payment_day INTEGER NOT NULL, -- Payment day (1-31)
    start_year_month VARCHAR(7) NOT NULL, -- Format: "2024-01"
    total_payments INTEGER, -- For loans
    remaining_payments INTEGER, -- For loans
    bank_account UUID NOT NULL REFERENCES bank_accounts(id),
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
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    year_month VARCHAR(7) NOT NULL, -- Format: "2024-01"
    total_amount BIGINT NOT NULL, -- Amount in cents
    is_confirmed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(asset_id, year_month)
);

-- App settings table
CREATE TABLE IF NOT EXISTS app_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key VARCHAR(255) NOT NULL,
    value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, key)
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_bank_accounts_user_id ON bank_accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_assets_user_id ON assets(user_id);
CREATE INDEX IF NOT EXISTS idx_assets_type ON assets(asset_type);
CREATE INDEX IF NOT EXISTS idx_income_sources_user_id ON income_sources(user_id);
CREATE INDEX IF NOT EXISTS idx_income_sources_active ON income_sources(is_active);
CREATE INDEX IF NOT EXISTS idx_monthly_income_records_income_source_id ON monthly_income_records(income_source_id);
CREATE INDEX IF NOT EXISTS idx_monthly_income_records_year_month ON monthly_income_records(year_month);
CREATE INDEX IF NOT EXISTS idx_recurring_payments_user_id ON recurring_payments(user_id);
CREATE INDEX IF NOT EXISTS idx_recurring_payments_active ON recurring_payments(is_active);
CREATE INDEX IF NOT EXISTS idx_card_monthly_totals_asset_id ON card_monthly_totals(asset_id);
CREATE INDEX IF NOT EXISTS idx_card_monthly_totals_year_month ON card_monthly_totals(year_month);
CREATE INDEX IF NOT EXISTS idx_app_settings_user_id ON app_settings(user_id);

-- Update triggers for updated_at columns
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bank_accounts_updated_at BEFORE UPDATE ON bank_accounts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_assets_updated_at BEFORE UPDATE ON assets
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
