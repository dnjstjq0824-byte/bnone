-- 000006_create_fx_tables.up.sql
-- BN ONE v1.0 Currencies, FX Providers, Exchange Rates & Snapshots

CREATE TABLE IF NOT EXISTS currencies (
    code VARCHAR(3) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    decimal_digits INT NOT NULL DEFAULT 2,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS fx_providers (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    official_website VARCHAR(255) NOT NULL DEFAULT '',
    api_base_url VARCHAR(255) NOT NULL DEFAULT '',
    license_status VARCHAR(50) NOT NULL DEFAULT 'PUBLIC_REFERENCE',
    commercial_use_status VARCHAR(50) NOT NULL DEFAULT 'ALLOWED',
    active BOOLEAN NOT NULL DEFAULT TRUE,
    last_checked_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE TABLE IF NOT EXISTS exchange_rates (
    id VARCHAR(36) PRIMARY KEY,
    provider_id VARCHAR(36) NOT NULL REFERENCES fx_providers(id) ON DELETE CASCADE,
    base_currency VARCHAR(3) NOT NULL,
    quote_currency VARCHAR(3) NOT NULL,
    rate NUMERIC(18, 6) NOT NULL,
    effective_at TIMESTAMP WITH TIME ZONE NOT NULL,
    retrieved_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'VERIFIED', -- VERIFIED, ESTIMATED, STALE
    CONSTRAINT uq_rate_pair_time UNIQUE (provider_id, base_currency, quote_currency, effective_at)
);

CREATE INDEX IF NOT EXISTS idx_exchange_rates_pair ON exchange_rates(base_currency, quote_currency, effective_at DESC);

CREATE TABLE IF NOT EXISTS fx_conversion_snapshots (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    base_currency VARCHAR(3) NOT NULL,
    quote_currency VARCHAR(3) NOT NULL,
    source_amount_minor BIGINT NOT NULL,
    converted_amount_minor BIGINT NOT NULL,
    rate NUMERIC(18, 6) NOT NULL,
    rate_effective_at TIMESTAMP WITH TIME ZONE NOT NULL,
    provider_id VARCHAR(36) NOT NULL REFERENCES fx_providers(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

