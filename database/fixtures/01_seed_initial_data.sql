-- 01_seed_initial_data.sql
-- BN ONE v1.0 Seed Data: Currencies, FX Providers, Rates, Categories, Subscriptions, Admin Roles

-- Currencies
INSERT INTO currencies (code, name, symbol, decimal_digits, is_active) VALUES
('KRW', 'South Korean Won', '₩', 0, TRUE),
('USD', 'United States Dollar', '$', 2, TRUE),
('EUR', 'Euro', '€', 2, TRUE),
('JPY', 'Japanese Yen', '¥', 0, TRUE),
('GBP', 'British Pound', '£', 2, TRUE),
('CNY', 'Chinese Yuan', '¥', 2, TRUE),
('SGD', 'Singapore Dollar', 'S$', 2, TRUE),
('AUD', 'Australian Dollar', 'A$', 2, TRUE),
('CAD', 'Canadian Dollar', 'C$', 2, TRUE),
('CHF', 'Swiss Franc', 'CHF', 2, TRUE),
('HKD', 'Hong Kong Dollar', 'HK$', 2, TRUE),
('THB', 'Thai Baht', '฿', 2, TRUE),
('VND', 'Vietnamese Dong', '₫', 0, TRUE),
('PHP', 'Philippine Peso', '₱', 2, TRUE)
ON CONFLICT (code) DO NOTHING;

-- FX Providers
INSERT INTO fx_providers (id, name, official_website, api_base_url, license_status, commercial_use_status, active, last_checked_at) VALUES
('p-bok-001', 'Bank of Korea Reference Rates', 'https://ecos.bok.or.kr', 'https://ecos.bok.or.kr/api', 'PUBLIC_GOV_OPEN', 'ALLOWED', TRUE, CURRENT_TIMESTAMP),
('p-ecb-002', 'European Central Bank Euro Reference', 'https://www.ecb.europa.eu', 'https://data-api.ecb.europa.eu', 'PUBLIC_OPEN', 'ALLOWED', TRUE, CURRENT_TIMESTAMP),
('p-open-003', 'Verified Global Reference FX', 'https://www.exchangerate-api.com', 'https://open.er-api.com/v6', 'PUBLIC_REFERENCE', 'ALLOWED', TRUE, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- Initial Exchange Rates (Base: KRW and USD)
INSERT INTO exchange_rates (id, provider_id, base_currency, quote_currency, rate, effective_at, retrieved_at, status) VALUES
('fx-krw-usd-1', 'p-bok-001', 'KRW', 'USD', 0.000741, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'VERIFIED'),
('fx-usd-krw-1', 'p-bok-001', 'USD', 'KRW', 1350.000000, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'VERIFIED'),
('fx-krw-jpy-1', 'p-bok-001', 'KRW', 'JPY', 0.111111, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'VERIFIED'),
('fx-jpy-krw-1', 'p-bok-001', 'JPY', 'KRW', 9.000000, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'VERIFIED'),
('fx-krw-eur-1', 'p-ecb-002', 'KRW', 'EUR', 0.000685, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'VERIFIED'),
('fx-eur-krw-1', 'p-ecb-002', 'EUR', 'KRW', 1460.000000, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'VERIFIED'),
('fx-krw-gbp-1', 'p-open-003', 'KRW', 'GBP', 0.000581, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'VERIFIED'),
('fx-gbp-krw-1', 'p-open-003', 'GBP', 'KRW', 1720.000000, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'VERIFIED'),
('fx-usd-eur-1', 'p-ecb-002', 'USD', 'EUR', 0.924000, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'VERIFIED'),
('fx-usd-jpy-1', 'p-open-003', 'USD', 'JPY', 150.000000, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'VERIFIED')
ON CONFLICT (id) DO NOTHING;

-- System Categories
INSERT INTO categories (id, owner_user_id, parent_id, code, name, category_type, icon, color, is_system, is_active, sort_order) VALUES
-- EXPENSE
('cat-exp-food', NULL, NULL, 'FOOD', '식비 & 카페', 'EXPENSE', 'utensils', '#F59E0B', TRUE, TRUE, 1),
('cat-exp-trans', NULL, NULL, 'TRANSPORT', '교통 & 차량', 'EXPENSE', 'car', '#3B82F6', TRUE, TRUE, 2),
('cat-exp-housing', NULL, NULL, 'HOUSING', '주거 & 관리비', 'EXPENSE', 'home', '#6366F1', TRUE, TRUE, 3),
('cat-exp-shop', NULL, NULL, 'SHOPPING', '쇼핑 & 마트', 'EXPENSE', 'shopping-bag', '#EC4899', TRUE, TRUE, 4),
('cat-exp-health', NULL, NULL, 'HEALTH', '의료 & 건강', 'EXPENSE', 'heart-pulse', '#EF4444', TRUE, TRUE, 5),
('cat-exp-leisure', NULL, NULL, 'LEISURE', '문화 & 여가', 'EXPENSE', 'film', '#8B5CF6', TRUE, TRUE, 6),
('cat-exp-travel', NULL, NULL, 'TRAVEL', '여행 & 숙박', 'EXPENSE', 'plane', '#06B6D4', TRUE, TRUE, 7),
('cat-exp-sub', NULL, NULL, 'SUBSCRIPTION', '구독 & 정기결제', 'EXPENSE', 'repeat', '#10B981', TRUE, TRUE, 8),
('cat-exp-edu', NULL, NULL, 'EDUCATION', '교육 & 자기계발', 'EXPENSE', 'book-open', '#14B8A6', TRUE, TRUE, 9),
('cat-exp-fee', NULL, NULL, 'FINANCIAL_FEE', '금융 수수료 & 이자', 'EXPENSE', 'banknote', '#64748B', TRUE, TRUE, 10),
('cat-exp-other', NULL, NULL, 'OTHER_EXPENSE', '기타 지출', 'EXPENSE', 'more-horizontal', '#94A3B8', TRUE, TRUE, 99),

-- INCOME
('cat-inc-salary', NULL, NULL, 'SALARY', '급여 & 정기 수입', 'INCOME', 'briefcase', '#10B981', TRUE, TRUE, 1),
('cat-inc-biz', NULL, NULL, 'BUSINESS', '사업 & 프리랜서', 'INCOME', 'building-2', '#059669', TRUE, TRUE, 2),
('cat-inc-invest', NULL, NULL, 'INVESTMENT', '금융소득 & 배당', 'INCOME', 'trending-up', '#047857', TRUE, TRUE, 3),
('cat-inc-gift', NULL, NULL, 'GIFT_ALLOWANCE', '용돈 & 상여금', 'INCOME', 'gift', '#34D399', TRUE, TRUE, 4),
('cat-inc-other', NULL, NULL, 'OTHER_INCOME', '기타 수입', 'INCOME', 'plus-circle', '#6EE7B7', TRUE, TRUE, 99),

-- TRANSFER
('cat-tra-transfer', NULL, NULL, 'ACCOUNT_TRANSFER', '계정 간 이동', 'TRANSFER', 'arrow-left-right', '#6366F1', TRUE, TRUE, 1),
('cat-tra-card', NULL, NULL, 'CARD_SETTLEMENT', '카드 대금 결제', 'TRANSFER', 'credit-card', '#4F46E5', TRUE, TRUE, 2)
ON CONFLICT (id) DO NOTHING;

-- Subscription Services Catalog
INSERT INTO subscription_services (id, name, slug, category, website, logo_asset_id, catalog_status) VALUES
('sub-net-001', 'Netflix', 'netflix', 'ENTERTAINMENT', 'https://netflix.com', 'netflix_logo', 'ACTIVE'),
('sub-spo-002', 'Spotify', 'spotify', 'MUSIC', 'https://spotify.com', 'spotify_logo', 'ACTIVE'),
('sub-ytp-003', 'YouTube Premium', 'youtube-premium', 'MEDIA', 'https://youtube.com', 'yt_premium_logo', 'ACTIVE'),
('sub-apo-004', 'Apple One', 'apple-one', 'ECOSYSTEM', 'https://apple.com', 'apple_logo', 'ACTIVE'),
('sub-gpt-005', 'ChatGPT Plus', 'chatgpt-plus', 'PRODUCTIVITY', 'https://openai.com', 'chatgpt_logo', 'ACTIVE'),
('sub-dis-006', 'Disney+', 'disney-plus', 'ENTERTAINMENT', 'https://disneyplus.com', 'disney_logo', 'ACTIVE'),
('sub-cpw-007', 'Coupang WOW', 'coupang-wow', 'COMMERCE', 'https://coupang.com', 'coupang_logo', 'ACTIVE'),
('sub-nvp-008', 'Naver Plus Membership', 'naver-plus', 'COMMERCE', 'https://naver.com', 'naver_logo', 'ACTIVE'),
('sub-cld-009', 'Google One (2TB)', 'google-one', 'STORAGE', 'https://one.google.com', 'google_one_logo', 'ACTIVE'),
('sub-adb-010', 'Adobe Creative Cloud', 'adobe-cc', 'DESIGN', 'https://adobe.com', 'adobe_logo', 'ACTIVE')
ON CONFLICT (id) DO NOTHING;

-- Admin Roles
INSERT INTO admin_roles (id, code, name, description) VALUES
('role-super-admin', 'PLATFORM_SUPER_ADMIN', 'Platform Super Admin', 'Full platform authority with strict audited data policies'),
('role-ops-admin', 'OPERATIONS_ADMIN', 'Operations Admin', 'Daily operational management, categories, and background jobs'),
('role-support-agent', 'SUPPORT_AGENT', 'Support Agent', 'Privacy-preserving customer support ticket management'),
('role-data-ops', 'DATA_OPERATIONS', 'Data Operations', 'FX Provider sync, exchange rates verification and quality inspection'),
('role-content-mgr', 'CONTENT_MANAGER', 'Content Manager', 'Localization strings, notification templates, public notices'),
('role-sec-admin', 'SECURITY_ADMIN', 'Security Admin', 'Security center, threat monitoring, session revocation, audit analysis'),
('role-analyst', 'ANALYST', 'Analyst', 'Aggregated and anonymized business metrics'),
('role-auditor', 'READ_ONLY_AUDITOR', 'Read-Only Auditor', 'Audited compliance and read-only inspection')
ON CONFLICT (id) DO NOTHING;

-- Default Feature Flags
INSERT INTO feature_flags (id, key, description, environment, enabled, updated_by, updated_at) VALUES
('ff-001', 'multi_currency_reporting', 'Enable multi-currency valuation snapshots across accounts', 'ALL', TRUE, 'SYSTEM', CURRENT_TIMESTAMP),
('ff-002', 'travel_budget_settlement', 'Enable travel module and expense link settlement', 'ALL', TRUE, 'SYSTEM', CURRENT_TIMESTAMP),
('ff-003', 'subscription_renewal_alerts', 'Send background notification for upcoming subscriptions', 'ALL', TRUE, 'SYSTEM', CURRENT_TIMESTAMP),
('ff-004', 'csv_data_exports', 'Allow users to export all transactions and accounts to CSV', 'ALL', TRUE, 'SYSTEM', CURRENT_TIMESTAMP),
('ff-005', 'account_deletion_flow', 'Allow users to initiate GDPR/PIPA compliant account deletion', 'ALL', TRUE, 'SYSTEM', CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

