// apps/web/src/api/mockData.ts
// Comprehensive realistic mock data aligned with Toss Design System (TDS) and BN ONE PRD

export const MOCK_USER = {
  id: 'u-tester-001',
  email: 'test@example.com',
  display_name: '토스테스터',
  country_code: 'KR',
  locale: 'ko-KR',
  base_currency: 'KRW',
  timezone: 'Asia/Seoul',
  status: 'ACTIVE',
  onboarding_completed_at: '2026-01-01T00:00:00Z',
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-10-31T00:00:00Z',
};

export const INITIAL_ACCOUNTS = [
  {
    id: 'acc-001',
    name: '토스뱅크 주거래통장',
    account_type: 'CHECKING',
    currency_code: 'KRW',
    calculated_balance_minor: 14850000,
    opening_balance_minor: 10000000,
    color: '#3182F6',
    icon: 'wallet',
    source_type: 'BANK',
  },
  {
    id: 'acc-002',
    name: 'KB국민 비상금통장',
    account_type: 'SAVINGS',
    currency_code: 'KRW',
    calculated_balance_minor: 32500000,
    opening_balance_minor: 30000000,
    color: '#F59E0B',
    icon: 'bank',
    source_type: 'BANK',
  },
  {
    id: 'acc-003',
    name: '미국 배당성장주 (USD)',
    account_type: 'INVESTMENT',
    currency_code: 'USD',
    calculated_balance_minor: 1825000, // $18,250.00
    opening_balance_minor: 1500000,
    color: '#10B981',
    icon: 'trending-up',
    source_type: 'BROKERAGE',
  },
  {
    id: 'acc-004',
    name: '일본 엔화 환전 (JPY)',
    account_type: 'CASH',
    currency_code: 'JPY',
    calculated_balance_minor: 85000, // ¥85,000
    opening_balance_minor: 50000,
    color: '#EC4899',
    icon: 'circle-dollar-sign',
    source_type: 'MANUAL',
  },
];

export const INITIAL_NET_WORTH = {
  base_currency: 'KRW',
  total_assets_minor: 72425000,
  financial_accounts_minor: 62425000,
  other_assets_minor: 10000000,
  total_liabilities_minor: 1500000,
  net_worth_minor: 70925000,
  disclaimer: '최신 환율(한국은행, ECB 기준)로 실시간 환산된 순자산입니다.',
};

export const INITIAL_TRANSACTIONS = [
  {
    id: 'tx-001',
    transaction_type: 'EXPENSE',
    occurred_at: new Date(Date.now() - 3600000 * 4).toISOString(),
    description: '스타벅스 강남점 (아메리카노 2잔)',
    merchant_name: '스타벅스',
    amount_minor: 9000,
    currency_code: 'KRW',
    category_name: '식비 & 카페',
    account_id: 'acc-001',
  },
  {
    id: 'tx-002',
    transaction_type: 'EXPENSE',
    occurred_at: new Date(Date.now() - 86400000).toISOString(),
    description: '쿠팡 로켓프레시 장보기',
    merchant_name: '쿠팡',
    amount_minor: 42800,
    currency_code: 'KRW',
    category_name: '쇼핑 & 마트',
    account_id: 'acc-001',
  },
  {
    id: 'tx-003',
    transaction_type: 'INCOME',
    occurred_at: new Date(Date.now() - 86400000 * 5).toISOString(),
    description: '10월 정기 급여',
    merchant_name: '주식회사 비엔',
    amount_minor: 4500000,
    currency_code: 'KRW',
    category_name: '급여 & 정기 수입',
    account_id: 'acc-002',
  },
  {
    id: 'tx-004',
    transaction_type: 'EXPENSE',
    occurred_at: new Date(Date.now() - 86400000 * 7).toISOString(),
    description: 'Netflix 4K 프리미엄',
    merchant_name: 'Netflix',
    amount_minor: 17000,
    currency_code: 'KRW',
    category_name: '구독 & 정기결제',
    account_id: 'acc-001',
  },
  {
    id: 'tx-005',
    transaction_type: 'EXPENSE',
    occurred_at: new Date(Date.now() - 86400000 * 8).toISOString(),
    description: '신분당선 교통비',
    merchant_name: '서울교통공사',
    amount_minor: 3100,
    currency_code: 'KRW',
    category_name: '교통 & 차량',
    account_id: 'acc-001',
  },
  {
    id: 'tx-006',
    transaction_type: 'EXPENSE',
    occurred_at: new Date(Date.now() - 86400000 * 10).toISOString(),
    description: 'ChatGPT Plus 정기 결제 (USD 20.00)',
    merchant_name: 'OpenAI',
    amount_minor: 27000,
    currency_code: 'KRW',
    category_name: '구독 & 정기결제',
    account_id: 'acc-003',
  },
];

export const INITIAL_ANALYTICS_OVERVIEW = {
  total_spent_minor: 1285000,
  total_income_minor: 4500000,
  savings_rate_pct: 71.4,
  top_categories: [
    { name: '식비 & 카페', amount_minor: 485000, percentage: 37.7 },
    { name: '쇼핑 & 마트', amount_minor: 320000, percentage: 24.9 },
    { name: '주거 & 관리비', amount_minor: 280000, percentage: 21.8 },
    { name: '교통 & 차량', amount_minor: 110000, percentage: 8.5 },
    { name: '구독 & 정기결제', amount_minor: 90000, percentage: 7.1 },
  ],
};

export const INITIAL_ASSETS = [
  {
    id: 'ast-001',
    name: '순금 골드바 10g',
    asset_type: 'PRECIOUS_METAL',
    estimated_value_minor: 1250000,
    currency_code: 'KRW',
    notes: 'KRX 금시장 현물 보관',
    updated_at: '2026-10-25T00:00:00Z',
  },
  {
    id: 'ast-002',
    name: '아파트 전세보증금',
    asset_type: 'REAL_ESTATE',
    estimated_value_minor: 350000000,
    currency_code: 'KRW',
    notes: '2027년 4월 만기 예정',
    updated_at: '2026-10-01T00:00:00Z',
  },
];

export const INITIAL_LIABILITIES = [
  {
    id: 'lia-001',
    name: '버팀목 전세자금대출',
    liability_type: 'MORTGAGE',
    outstanding_amount_minor: 120000000,
    currency_code: 'KRW',
    notes: '금리 2.1% 고정',
    updated_at: '2026-10-15T00:00:00Z',
  },
  {
    id: 'lia-002',
    name: '현대카드 10월 결제예정금',
    liability_type: 'CREDIT_CARD',
    outstanding_amount_minor: 850000,
    currency_code: 'KRW',
    notes: '25일 자동이체',
    updated_at: '2026-10-28T00:00:00Z',
  },
];

export const INITIAL_BUDGETS = [
  {
    id: 'bgt-001',
    name: '10월 식비 예산',
    category_id: 'cat-exp-food',
    category_name: '식비 & 카페',
    period: 'MONTHLY',
    amount_minor: 600000,
    spent_minor: 485000,
    currency_code: 'KRW',
  },
  {
    id: 'bgt-002',
    name: '쇼핑 & 마트 예산',
    category_id: 'cat-exp-shop',
    category_name: '쇼핑 & 마트',
    period: 'MONTHLY',
    amount_minor: 400000,
    spent_minor: 320000,
    currency_code: 'KRW',
  },
  {
    id: 'bgt-003',
    name: '문화/여가 생활',
    category_id: 'cat-exp-leisure',
    category_name: '문화 & 여가',
    period: 'MONTHLY',
    amount_minor: 200000,
    spent_minor: 75000,
    currency_code: 'KRW',
  },
];

export const INITIAL_SUBSCRIPTIONS = [
  {
    id: 'sub-001',
    name: 'Netflix Premium (4K)',
    category: 'ENTERTAINMENT',
    amount_minor: 17000,
    currency_code: 'KRW',
    billing_cycle: 'MONTHLY',
    next_billing_date: '2026-11-15',
    status: 'ACTIVE',
    service_id: 'sub-net-001',
  },
  {
    id: 'sub-002',
    name: 'YouTube Premium',
    category: 'MEDIA',
    amount_minor: 14900,
    currency_code: 'KRW',
    billing_cycle: 'MONTHLY',
    next_billing_date: '2026-11-20',
    status: 'ACTIVE',
    service_id: 'sub-ytp-003',
  },
  {
    id: 'sub-003',
    name: 'ChatGPT Plus (OpenAI)',
    category: 'PRODUCTIVITY',
    amount_minor: 27000,
    currency_code: 'KRW',
    billing_cycle: 'MONTHLY',
    next_billing_date: '2026-11-10',
    status: 'ACTIVE',
    service_id: 'sub-gpt-005',
  },
  {
    id: 'sub-004',
    name: '쿠팡 와우 멤버십',
    category: 'COMMERCE',
    amount_minor: 7890,
    currency_code: 'KRW',
    billing_cycle: 'MONTHLY',
    next_billing_date: '2026-11-04',
    status: 'ACTIVE',
    service_id: 'sub-cpw-007',
  },
];

export const INITIAL_SUBSCRIPTION_SUMMARY = {
  active_count: 4,
  monthly_estimate_minor: 66790,
  yearly_estimate_minor: 801480,
  base_currency: 'KRW',
};

export const INITIAL_TRIPS = [
  {
    id: 'trp-001',
    name: '도쿄 3박 4일 미식 & 커피 투어',
    destination: '도쿄, 일본',
    start_date: '2026-11-20',
    end_date: '2026-11-23',
    budget_minor: 1200000,
    spent_minor: 350000,
    currency_code: 'KRW',
    status: 'PLANNING',
  },
  {
    id: 'trp-002',
    name: '싱가포르 테크 컨퍼런스 출장',
    destination: '싱가포르',
    start_date: '2026-12-05',
    end_date: '2026-12-08',
    budget_minor: 2500000,
    spent_minor: 0,
    currency_code: 'KRW',
    status: 'PLANNING',
  },
];

export const INITIAL_CURRENCIES = [
  { code: 'KRW', name: 'South Korean Won', symbol: '₩', decimal_digits: 0 },
  { code: 'USD', name: 'United States Dollar', symbol: '$', decimal_digits: 2 },
  { code: 'EUR', name: 'Euro', symbol: '€', decimal_digits: 2 },
  { code: 'JPY', name: 'Japanese Yen', symbol: '¥', decimal_digits: 0 },
  { code: 'GBP', name: 'British Pound', symbol: '£', decimal_digits: 2 },
  { code: 'SGD', name: 'Singapore Dollar', symbol: 'S$', decimal_digits: 2 },
];

export const INITIAL_EXCHANGE_RATES = [
  { id: 'fx-1', base_currency: 'USD', quote_currency: 'KRW', rate: 1350.0, provider_name: '한국은행(BOK)', status: 'VERIFIED', effective_at: new Date().toISOString() },
  { id: 'fx-2', base_currency: 'USD', quote_currency: 'JPY', rate: 150.0, provider_name: '글로벌 레퍼런스', status: 'VERIFIED', effective_at: new Date().toISOString() },
  { id: 'fx-3', base_currency: 'USD', quote_currency: 'EUR', rate: 0.924, provider_name: '유럽중앙은행(ECB)', status: 'VERIFIED', effective_at: new Date().toISOString() },
  { id: 'fx-4', base_currency: 'KRW', quote_currency: 'USD', rate: 0.000741, provider_name: '한국은행(BOK)', status: 'VERIFIED', effective_at: new Date().toISOString() },
  { id: 'fx-5', base_currency: 'KRW', quote_currency: 'JPY', rate: 0.111, provider_name: '한국은행(BOK)', status: 'VERIFIED', effective_at: new Date().toISOString() },
];

export const INITIAL_CATEGORIES = [
  { id: 'cat-exp-food', name: '식비 & 카페', category_type: 'EXPENSE', icon: 'utensils', color: '#F59E0B' },
  { id: 'cat-exp-trans', name: '교통 & 차량', category_type: 'EXPENSE', icon: 'car', color: '#3B82F6' },
  { id: 'cat-exp-housing', name: '주거 & 관리비', category_type: 'EXPENSE', icon: 'home', color: '#6366F1' },
  { id: 'cat-exp-shop', name: '쇼핑 & 마트', category_type: 'EXPENSE', icon: 'shopping-bag', color: '#EC4899' },
  { id: 'cat-exp-health', name: '의료 & 건강', category_type: 'EXPENSE', icon: 'heart-pulse', color: '#EF4444' },
  { id: 'cat-exp-leisure', name: '문화 & 여가', category_type: 'EXPENSE', icon: 'film', color: '#8B5CF6' },
  { id: 'cat-exp-sub', name: '구독 & 정기결제', category_type: 'EXPENSE', icon: 'repeat', color: '#10B981' },
  { id: 'cat-inc-salary', name: '급여 & 정기 수입', category_type: 'INCOME', icon: 'briefcase', color: '#10B981' },
  { id: 'cat-inc-biz', name: '사업 & 프리랜서', category_type: 'INCOME', icon: 'building-2', color: '#059669' },
  { id: 'cat-tra-transfer', name: '계정 간 이동', category_type: 'TRANSFER', icon: 'arrow-left-right', color: '#6366F1' },
];

export const INITIAL_NOTIFICATIONS = [
  {
    id: 'noti-001',
    title: '🎉 BN ONE 슈퍼앱에 오신 것을 환영합니다!',
    body: '토스 스타일의 직관적인 다통화 자산 관리 서비스를 경험해보세요.',
    is_read: false,
    created_at: new Date(Date.now() - 1800000).toISOString(),
    notification_type: 'ANNOUNCEMENT',
  },
  {
    id: 'noti-002',
    title: '🔔 쿠팡 와우 멤버십 결제 예정 (D-3)',
    body: '11월 4일에 7,890원이 결제될 예정입니다.',
    is_read: false,
    created_at: new Date(Date.now() - 3600000 * 5).toISOString(),
    notification_type: 'SUBSCRIPTION_RENEWAL',
  },
  {
    id: 'noti-003',
    title: '📈 환율 알림: 원/달러 1,350원 기준선 도달',
    body: '한국은행 공시 기준환율이 업데이트되었습니다.',
    is_read: true,
    created_at: new Date(Date.now() - 86400000).toISOString(),
    notification_type: 'FX_ALERT',
  },
];

