-- 000005_create_subscription_tables.up.sql
-- BN ONE v1.0 Subscription Catalog, User Subscriptions & Reminders

CREATE TABLE IF NOT EXISTS subscription_services (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    category VARCHAR(50) NOT NULL DEFAULT 'MEDIA',
    website VARCHAR(255) NOT NULL DEFAULT '',
    logo_asset_id VARCHAR(100) NOT NULL DEFAULT '',
    catalog_status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, INACTIVE
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_subscriptions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    service_id VARCHAR(36) NULL REFERENCES subscription_services(id) ON DELETE SET NULL,
    custom_service_name VARCHAR(100) NULL,
    payment_account_id VARCHAR(36) NULL REFERENCES finance_accounts(id) ON DELETE SET NULL,
    amount_minor BIGINT NOT NULL,
    currency_code VARCHAR(3) NOT NULL DEFAULT 'KRW',
    billing_cycle VARCHAR(20) NOT NULL DEFAULT 'MONTHLY', -- MONTHLY, YEARLY, WEEKLY
    billing_interval INT NOT NULL DEFAULT 1,
    start_date DATE NOT NULL,
    next_renewal_at DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, PAUSED, CANCELLED_RECORD
    reminder_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_subscriptions_user ON user_subscriptions(user_id, status);
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_renewal ON user_subscriptions(next_renewal_at);

CREATE TABLE IF NOT EXISTS subscription_history (
    id VARCHAR(36) PRIMARY KEY,
    subscription_id VARCHAR(36) NOT NULL REFERENCES user_subscriptions(id) ON DELETE CASCADE,
    change_type VARCHAR(50) NOT NULL, -- CREATED, UPDATED, PAUSED, RESUMED, CANCELLED_RECORD
    before_snapshot TEXT NOT NULL DEFAULT '',
    after_snapshot TEXT NOT NULL DEFAULT '',
    changed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS subscription_reminders (
    id VARCHAR(36) PRIMARY KEY,
    subscription_id VARCHAR(36) NOT NULL REFERENCES user_subscriptions(id) ON DELETE CASCADE,
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- PENDING, SENT, FAILED, CANCELLED
    sent_at TIMESTAMP WITH TIME ZONE NULL,
    deduplication_key VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

