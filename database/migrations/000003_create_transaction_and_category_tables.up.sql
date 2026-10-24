-- 000003_create_transaction_and_category_tables.up.sql
-- BN ONE v1.0 Categories, Transactions, Splits & Signed Account Entries

CREATE TABLE IF NOT EXISTS categories (
    id VARCHAR(36) PRIMARY KEY,
    owner_user_id VARCHAR(36) NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id VARCHAR(36) NULL REFERENCES categories(id) ON DELETE SET NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    category_type VARCHAR(20) NOT NULL, -- INCOME, EXPENSE, TRANSFER
    icon VARCHAR(50) NOT NULL DEFAULT 'tag',
    color VARCHAR(20) NOT NULL DEFAULT '#64748B',
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_categories_owner ON categories(owner_user_id, category_type);

CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    transaction_type VARCHAR(20) NOT NULL, -- INCOME, EXPENSE, TRANSFER, ADJUSTMENT
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL,
    description VARCHAR(255) NOT NULL DEFAULT '',
    merchant_name VARCHAR(100) NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    trip_id VARCHAR(36) NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'POSTED', -- POSTED, VOIDED, PENDING
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX IF NOT EXISTS idx_transactions_user_date ON transactions(user_id, occurred_at);
CREATE INDEX IF NOT EXISTS idx_transactions_trip ON transactions(trip_id);

CREATE TABLE IF NOT EXISTS transaction_splits (
    id VARCHAR(36) PRIMARY KEY,
    transaction_id VARCHAR(36) NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    category_id VARCHAR(36) NOT NULL REFERENCES categories(id),
    amount_minor BIGINT NOT NULL,
    currency_code VARCHAR(3) NOT NULL,
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_splits_tx_id ON transaction_splits(transaction_id);

CREATE TABLE IF NOT EXISTS account_entries (
    id VARCHAR(36) PRIMARY KEY,
    transaction_id VARCHAR(36) NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    account_id VARCHAR(36) NOT NULL REFERENCES finance_accounts(id) ON DELETE CASCADE,
    amount_minor_signed BIGINT NOT NULL,
    currency_code VARCHAR(3) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_entries_account_id ON account_entries(account_id);
CREATE INDEX IF NOT EXISTS idx_entries_tx_id ON account_entries(transaction_id);

CREATE TABLE IF NOT EXISTS transaction_revisions (
    id VARCHAR(36) PRIMARY KEY,
    transaction_id VARCHAR(36) NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    revision_number INT NOT NULL,
    before_snapshot TEXT NOT NULL,
    after_snapshot TEXT NOT NULL,
    changed_by VARCHAR(36) NOT NULL,
    change_reason VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_revisions_tx_id ON transaction_revisions(transaction_id);
