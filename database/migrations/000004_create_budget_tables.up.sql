-- 000004_create_budget_tables.up.sql
-- BN ONE v1.0 Budgets, Category Allocations & Alert Rules

CREATE TABLE IF NOT EXISTS budgets (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    period_type VARCHAR(20) NOT NULL DEFAULT 'MONTHLY', -- MONTHLY, YEARLY, CUSTOM
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    currency_code VARCHAR(3) NOT NULL DEFAULT 'KRW',
    limit_minor BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, ARCHIVED
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_budgets_user_period ON budgets(user_id, start_date, end_date);

CREATE TABLE IF NOT EXISTS budget_categories (
    id VARCHAR(36) PRIMARY KEY,
    budget_id VARCHAR(36) NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    category_id VARCHAR(36) NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    allocated_amount_minor BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_budget_category UNIQUE (budget_id, category_id)
);

CREATE TABLE IF NOT EXISTS budget_alert_rules (
    id VARCHAR(36) PRIMARY KEY,
    budget_id VARCHAR(36) NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    threshold_percent INT NOT NULL DEFAULT 80, -- 50, 80, 100
    channel VARCHAR(20) NOT NULL DEFAULT 'IN_APP', -- IN_APP, PUSH, EMAIL
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    last_triggered_at TIMESTAMP WITH TIME ZONE NULL
);

