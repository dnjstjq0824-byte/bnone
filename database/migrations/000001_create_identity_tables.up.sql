-- 000001_create_identity_tables.up.sql
-- BN ONE v1.0 Identity & User Tables

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    country_code VARCHAR(2) NOT NULL DEFAULT 'KR',
    locale VARCHAR(10) NOT NULL DEFAULT 'ko-KR',
    base_currency VARCHAR(3) NOT NULL DEFAULT 'KRW',
    timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Seoul',
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, SUSPENDED, DELETION_REQUESTED, DELETED
    email_verified_at TIMESTAMP WITH TIME ZONE NULL,
    onboarding_completed_at TIMESTAMP WITH TIME ZONE NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

CREATE TABLE IF NOT EXISTS user_sessions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash VARCHAR(255) NOT NULL,
    device_id VARCHAR(100) NOT NULL DEFAULT '',
    device_name VARCHAR(100) NOT NULL DEFAULT '',
    platform VARCHAR(20) NOT NULL DEFAULT 'WEB', -- ANDROID, IOS, WEB, ADMIN
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_refresh_token ON user_sessions(refresh_token_hash);

CREATE TABLE IF NOT EXISTS user_preferences (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    mask_amounts BOOLEAN NOT NULL DEFAULT FALSE,
    dark_mode BOOLEAN NOT NULL DEFAULT FALSE,
    start_of_week VARCHAR(10) NOT NULL DEFAULT 'SUNDAY', -- SUNDAY, MONDAY
    notification_settings TEXT NOT NULL DEFAULT '{"push":true,"email":true,"budget_alerts":true,"renewal_reminders":true}',
    analytics_consent BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS consent_records (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    consent_type VARCHAR(50) NOT NULL, -- TERMS_OF_SERVICE, PRIVACY_POLICY, MARKETING
    policy_version VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'GRANTED', -- GRANTED, WITHDRAWN
    consented_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    withdrawn_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX IF NOT EXISTS idx_consent_records_user_id ON consent_records(user_id);

