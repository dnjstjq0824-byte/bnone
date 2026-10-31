// apps/web/src/api/mockService.ts
import {
  MOCK_USER,
  INITIAL_ACCOUNTS,
  INITIAL_ASSETS,
  INITIAL_LIABILITIES,
  INITIAL_TRANSACTIONS,
  INITIAL_ANALYTICS_OVERVIEW,
  INITIAL_BUDGETS,
  INITIAL_SUBSCRIPTIONS,
  INITIAL_SUBSCRIPTION_SUMMARY,
  INITIAL_TRIPS,
  INITIAL_CURRENCIES,
  INITIAL_EXCHANGE_RATES,
  INITIAL_CATEGORIES,
  INITIAL_NOTIFICATIONS,
  calculateNetWorth,
} from './mockData';
import { ApiError } from './client';

function getStorage<T>(key: string, defaultValue: T): T {
  try {
    const item = localStorage.getItem(`bn_mock_${key}`);
    return item ? JSON.parse(item) : defaultValue;
  } catch {
    return defaultValue;
  }
}

function setStorage<T>(key: string, value: T): void {
  try {
    localStorage.setItem(`bn_mock_${key}`, JSON.stringify(value));
  } catch {}
}

export function handleMockApi<T = any>(endpoint: string, options: RequestInit = {}): T {
  const method = (options.method || 'GET').toUpperCase();
  const body = options.body ? JSON.parse(options.body as string) : {};
  const [path, queryString] = endpoint.split('?');
  const params = new URLSearchParams(queryString || '');

  // Fail-closed check: If demo mode is not explicitly enabled and no backend is set, reject all logins
  const isDemoModeEnabled = typeof import.meta !== 'undefined' && import.meta.env ? (import.meta.env.VITE_DEMO_MODE === 'true' || import.meta.env.VITE_DEMO_MODE === undefined) : true;
  const hasBackend = typeof import.meta !== 'undefined' && import.meta.env ? !!import.meta.env.VITE_API_URL : false;
  if (!isDemoModeEnabled && !hasBackend) {
    throw new ApiError('BACKEND_UNAVAILABLE', '백엔드 API 서버가 연결되지 않았습니다. 관리자에게 문의하세요.');
  }

  // Auth: Login - Strictly validates credentials even in demo mode
  if (path === '/auth/login' && method === 'POST') {
    const email = (body.email || '').trim().toLowerCase();
    const password = body.password || '';

    // Only test@example.com / password123 is valid for demo login
    if (email !== 'test@example.com' || password !== 'password123') {
      throw new ApiError('INVALID_CREDENTIALS', '이메일 또는 비밀번호가 올바르지 않습니다.');
    }

    const user = {
      ...MOCK_USER,
      email: 'test@example.com',
      display_name: '토스테스터',
    };
    setStorage('user', user);

    return {
      access_token: 'bn-mock-jwt-token-verified',
      refresh_token: 'bn-mock-refresh-token-verified',
      expires_in: 86400,
      user,
      email_verification_required: false,
      onboarding_required: false,
    } as T;
  }

  // Auth: Register
  if (path === '/auth/register' && method === 'POST') {
    const email = (body.email || '').trim().toLowerCase();
    if (!email || !email.includes('@')) {
      throw new ApiError('INVALID_EMAIL', '올바른 이메일 주소를 입력해주세요.');
    }
    if ((body.password || '').length < 8) {
      throw new ApiError('WEAK_PASSWORD', '비밀번호는 최소 8자 이상이어야 합니다.');
    }

    const user = {
      ...MOCK_USER,
      id: `u-${Date.now()}`,
      email,
      display_name: body.display_name || email.split('@')[0],
      base_currency: body.base_currency || 'KRW',
    };
    setStorage('user', user);

    return {
      access_token: 'bn-mock-jwt-token-verified',
      refresh_token: 'bn-mock-refresh-token-verified',
      expires_in: 86400,
      user,
      email_verification_required: false,
      onboarding_required: false,
    } as T;
  }

  // Auth: Logout
  if (path === '/auth/logout' && method === 'POST') {
    return { success: true } as T;
  }

  // User: /me
  if (path === '/me') {
    if (method === 'PATCH' || method === 'PUT') {
      const user = { ...getStorage('user', MOCK_USER), ...body };
      setStorage('user', user);
      return user as T;
    }
    return getStorage('user', MOCK_USER) as T;
  }

  // Accounts
  if (path === '/accounts') {
    const accounts = getStorage('accounts', INITIAL_ACCOUNTS);
    if (method === 'POST') {
      const newAcc = {
        id: `acc-${Date.now()}`,
        name: body.name || '새 계정',
        account_type: body.account_type || 'CHECKING',
        currency_code: body.currency_code || 'KRW',
        calculated_balance_minor: Number(body.opening_balance_minor || 0),
        opening_balance_minor: Number(body.opening_balance_minor || 0),
        color: body.color || '#3182F6',
        icon: body.icon || 'wallet',
        source_type: body.source_type || 'MANUAL',
      };
      const updated = [...accounts, newAcc];
      setStorage('accounts', updated);
      return newAcc as T;
    }
    return accounts as T;
  }

  if (path.startsWith('/accounts/') && path.endsWith('/archive') && method === 'POST') {
    const id = path.split('/')[2];
    const accounts = getStorage('accounts', INITIAL_ACCOUNTS).filter((a: any) => a.id !== id);
    setStorage('accounts', accounts);
    return { success: true } as T;
  }

  // Net Worth (Single Source of Truth: dynamically computed)
  if (path === '/net-worth') {
    const accounts = getStorage('accounts', INITIAL_ACCOUNTS);
    const assets = getStorage('assets', INITIAL_ASSETS);
    const liabilities = getStorage('liabilities', INITIAL_LIABILITIES);
    const calculated = calculateNetWorth(accounts, assets, liabilities);
    return calculated as T;
  }

  // Transactions
  if (path === '/transactions') {
    const txs = getStorage('transactions', INITIAL_TRANSACTIONS);
    if (method === 'POST') {
      const newTx = {
        id: `tx-${Date.now()}`,
        transaction_type: body.transaction_type || 'EXPENSE',
        occurred_at: body.occurred_at || new Date().toISOString(),
        description: body.description || '새 지출',
        merchant_name: body.merchant_name || body.description || '',
        amount_minor: Number(body.amount_minor || 0),
        currency_code: body.currency_code || 'KRW',
        category_name: body.category_name || '기타',
        account_id: body.account_id || 'acc-001',
      };
      const updated = [newTx, ...txs];
      setStorage('transactions', updated);

      // Update account balance
      const accounts = getStorage('accounts', INITIAL_ACCOUNTS);
      const acc = accounts.find((a: any) => a.id === newTx.account_id);
      if (acc) {
        if (newTx.transaction_type === 'EXPENSE') acc.calculated_balance_minor -= newTx.amount_minor;
        else if (newTx.transaction_type === 'INCOME') acc.calculated_balance_minor += newTx.amount_minor;
        setStorage('accounts', accounts);
      }

      return newTx as T;
    }

    const limit = Number(params.get('limit') || 50);
    const accId = params.get('account_id');
    let filtered = txs;
    if (accId) filtered = filtered.filter((t: any) => t.account_id === accId);
    return filtered.slice(0, limit) as T;
  }

  if (path.startsWith('/transactions/') && method === 'DELETE') {
    const id = path.split('/')[2];
    const txs = getStorage('transactions', INITIAL_TRANSACTIONS).filter((t: any) => t.id !== id);
    setStorage('transactions', txs);
    return { success: true } as T;
  }

  // Manual Transfer
  if (path === '/transfers/manual' && method === 'POST') {
    const accounts = getStorage('accounts', INITIAL_ACCOUNTS);
    const fromAcc = accounts.find((a: any) => a.id === body.source_account_id);
    const toAcc = accounts.find((a: any) => a.id === body.destination_account_id);
    const amount = Number(body.amount_minor || 0);

    if (fromAcc) fromAcc.calculated_balance_minor -= amount;
    if (toAcc) toAcc.calculated_balance_minor += amount;
    setStorage('accounts', accounts);

    const txs = getStorage('transactions', INITIAL_TRANSACTIONS);
    const transferTx = {
      id: `tx-tf-${Date.now()}`,
      transaction_type: 'TRANSFER',
      occurred_at: new Date().toISOString(),
      description: `${fromAcc?.name || '출금'} → ${toAcc?.name || '입금'} 이체`,
      merchant_name: '계정 간 이체',
      amount_minor: amount,
      currency_code: fromAcc?.currency_code || 'KRW',
      category_name: '계정 간 이동',
      account_id: fromAcc?.id || 'acc-001',
    };
    setStorage('transactions', [transferTx, ...txs]);
    return { success: true, transfer_id: `tr-${Date.now()}` } as T;
  }

  // Analytics Overview (Aligned with AnalyticsPage contract)
  if (path === '/analytics/overview') {
    return INITIAL_ANALYTICS_OVERVIEW as T;
  }

  // Assets & Liabilities
  if (path === '/assets') {
    const assets = getStorage('assets', INITIAL_ASSETS);
    if (method === 'POST') {
      const newAsset = {
        id: `ast-${Date.now()}`,
        name: body.name,
        asset_type: body.asset_type || 'OTHER',
        estimated_value_minor: Number(body.estimated_value_minor || 0),
        currency_code: body.currency_code || 'KRW',
        notes: body.notes || '',
        updated_at: new Date().toISOString(),
      };
      setStorage('assets', [...assets, newAsset]);
      return newAsset as T;
    }
    return assets as T;
  }

  if (path === '/liabilities') {
    const liabilities = getStorage('liabilities', INITIAL_LIABILITIES);
    if (method === 'POST') {
      const newLia = {
        id: `lia-${Date.now()}`,
        name: body.name,
        liability_type: body.liability_type || 'OTHER',
        outstanding_amount_minor: Number(body.outstanding_amount_minor || 0),
        currency_code: body.currency_code || 'KRW',
        notes: body.notes || '',
        updated_at: new Date().toISOString(),
      };
      setStorage('liabilities', [...liabilities, newLia]);
      return newLia as T;
    }
    return liabilities as T;
  }

  // Budgets (Aligned with BudgetsPage contract: limit_minor, remaining_minor, usage_percent, start_date, end_date)
  if (path === '/budgets') {
    const budgets = getStorage('budgets', INITIAL_BUDGETS);
    if (method === 'POST') {
      const limit = Number(body.limit_minor || body.amount_minor || 0);
      const spent = 0;
      const newBgt = {
        id: `bgt-${Date.now()}`,
        name: body.name || '새 예산',
        category_id: body.category_id || 'cat-exp-food',
        category_name: body.category_name || '식비 & 카페',
        period: body.period || 'MONTHLY',
        limit_minor: limit,
        amount_minor: limit,
        spent_minor: spent,
        remaining_minor: limit - spent,
        usage_percent: limit > 0 ? (spent / limit) * 100 : 0,
        start_date: '2026-11-01',
        end_date: '2026-11-30',
        currency_code: body.currency_code || 'KRW',
      };
      setStorage('budgets', [...budgets, newBgt]);
      return newBgt as T;
    }
    return budgets as T;
  }

  if (path.startsWith('/budgets/') && method === 'DELETE') {
    const id = path.split('/')[2];
    const budgets = getStorage('budgets', INITIAL_BUDGETS).filter((b: any) => b.id !== id);
    setStorage('budgets', budgets);
    return { success: true } as T;
  }

  // Subscriptions (Aligned with SubscriptionsPage: service_name, next_renewal_at)
  if (path === '/me/subscriptions') {
    const subs = getStorage('subscriptions', INITIAL_SUBSCRIPTIONS);
    if (method === 'POST') {
      const serviceName = body.service_name || body.name || body.custom_service_name || '새 구독';
      const newSub = {
        id: `sub-${Date.now()}`,
        service_name: serviceName,
        name: serviceName,
        category: body.category || 'ENTERTAINMENT',
        amount_minor: Number(body.amount_minor || 0),
        currency_code: body.currency_code || 'KRW',
        billing_cycle: body.billing_cycle || 'MONTHLY',
        next_renewal_at: body.next_renewal_at || body.next_billing_date || '2026-11-30',
        next_billing_date: body.next_renewal_at || body.next_billing_date || '2026-11-30',
        status: 'ACTIVE',
        service_id: body.service_id || null,
      };
      setStorage('subscriptions', [...subs, newSub]);
      return newSub as T;
    }
    return subs as T;
  }

  if (path === '/me/subscriptions/summary') {
    const subs = getStorage('subscriptions', INITIAL_SUBSCRIPTIONS);
    const active = subs.filter((s: any) => s.status === 'ACTIVE');
    const monthly = active.reduce((sum: number, s: any) => sum + s.amount_minor, 0);
    return {
      active_count: active.length,
      monthly_estimate_minor: monthly,
      yearly_estimate_minor: monthly * 12,
      base_currency: 'KRW',
    } as T;
  }

  if (path === '/subscription-services') {
    return [
      { id: 'sub-net-001', name: 'Netflix', category: 'ENTERTAINMENT', website: 'https://netflix.com' },
      { id: 'sub-ytp-003', name: 'YouTube Premium', category: 'MEDIA', website: 'https://youtube.com' },
      { id: 'sub-gpt-005', name: 'ChatGPT Plus', category: 'PRODUCTIVITY', website: 'https://openai.com' },
      { id: 'sub-cpw-007', name: 'Coupang WOW', category: 'COMMERCE', website: 'https://coupang.com' },
      { id: 'sub-spo-002', name: 'Spotify', category: 'MUSIC', website: 'https://spotify.com' },
    ] as T;
  }

  if (path.startsWith('/me/subscriptions/') && method === 'POST') {
    const parts = path.split('/');
    const id = parts[3];
    const action = parts[4];
    const subs = getStorage('subscriptions', INITIAL_SUBSCRIPTIONS);
    if (action === 'cancel-record') {
      const filtered = subs.filter((s: any) => s.id !== id);
      setStorage('subscriptions', filtered);
    } else {
      const updated = subs.map((s: any) => {
        if (s.id === id) {
          return { ...s, status: s.status === 'ACTIVE' ? 'PAUSED' : 'ACTIVE' };
        }
        return s;
      });
      setStorage('subscriptions', updated);
    }
    return { success: true } as T;
  }

  // Trips (Aligned with TripsPage: total_budget_minor, total_spent_minor, base_currency, city_name, country_code)
  if (path === '/trips') {
    const trips = getStorage('trips', INITIAL_TRIPS);
    if (method === 'POST') {
      const budget = Number(body.total_budget_minor || body.budget_minor || 0);
      const newTrip = {
        id: `trp-${Date.now()}`,
        name: body.name || '새 여행',
        city_name: body.city_name || '여행지',
        country_code: body.country_code || 'KR',
        destination: `${body.city_name || '여행지'}, ${body.country_code || 'KR'}`,
        start_date: body.start_date || '2026-12-01',
        end_date: body.end_date || '2026-12-05',
        total_budget_minor: budget,
        total_spent_minor: 0,
        budget_minor: budget,
        spent_minor: 0,
        base_currency: body.base_currency || body.currency_code || 'KRW',
        currency_code: body.base_currency || body.currency_code || 'KRW',
        status: 'PLANNING',
      };
      setStorage('trips', [...trips, newTrip]);
      return newTrip as T;
    }
    return trips as T;
  }

  if (path.includes('/report')) {
    return {
      trip_id: 'trp-001',
      total_budget_minor: 1200000,
      total_spent_minor: 350000,
      base_currency: 'KRW',
      settlements: [],
    } as T;
  }

  // Currencies & Reference FX
  if (path === '/currencies') {
    return INITIAL_CURRENCIES as T;
  }

  // Aligned with ForexPage: returns { rates: [...] } as well as array compatibility
  if (path === '/fx/rates') {
    const ratesData = {
      base_currency: params.get('base') || 'USD',
      effective_at: '2026-10-24T00:00:00Z',
      is_example: true,
      rates: INITIAL_EXCHANGE_RATES,
    };
    return ratesData as T;
  }

  if (path === '/fx/convert' && method === 'POST') {
    const from = body.source_currency || 'USD';
    const to = body.target_currency || 'KRW';
    const amount = Number(body.source_amount_minor || body.amount_minor || 0);

    let rate = 1350.0;
    if (from === 'USD' && to === 'KRW') rate = 1350.0;
    else if (from === 'KRW' && to === 'USD') rate = 0.000741;
    else if (from === 'JPY' && to === 'KRW') rate = 9.0;
    else if (from === 'EUR' && to === 'KRW') rate = 1460.0;

    const converted = Math.round(amount * rate);
    return {
      source_currency: from,
      target_currency: to,
      source_amount_minor: amount,
      original_amount_minor: amount,
      converted_amount_minor: converted,
      rate,
      effective_at: '2026-10-24T00:00:00Z',
      status: 'EXAMPLE_RATE',
    } as T;
  }

  // Categories
  if (path === '/categories') {
    return INITIAL_CATEGORIES as T;
  }

  // Notifications
  if (path === '/me/notifications') {
    return getStorage('notifications', INITIAL_NOTIFICATIONS) as T;
  }

  if (path.startsWith('/me/notifications/') && path.endsWith('/read')) {
    const id = path.split('/')[3];
    const notis = getStorage('notifications', INITIAL_NOTIFICATIONS).map((n: any) =>
      n.id === id ? { ...n, is_read: true } : n
    );
    setStorage('notifications', notis);
    return { success: true } as T;
  }

  // Sessions
  if (path === '/me/sessions') {
    return [
      {
        id: 'sess-current',
        platform: 'WEB (Chrome / macOS & Windows)',
        device_id: 'browser-active-session',
        ip_address: '211.234.120.45',
        created_at: '2026-10-31T08:00:00Z',
        last_active_at: new Date().toISOString(),
        is_current: true,
      },
    ] as T;
  }

  if (path.startsWith('/me/sessions/') && method === 'DELETE') {
    return { success: true } as T;
  }

  // Support Tickets
  if (path === '/me/support-tickets') {
    return [
      {
        id: 'tkt-001',
        title: '다통화 환전 환율 갱신 주기 문의',
        status: 'ANSWERED',
        created_at: '2026-10-28T10:00:00Z',
        messages: [
          { sender: 'USER', message: '한국은행 기준 환율은 하루에 몇 번 갱신되나요?', sent_at: '2026-10-28T10:00:00Z' },
          { sender: 'AGENT', message: '데모 모드에서는 2026-10-24 기준 예시 환율이 고정 표기됩니다.', sent_at: '2026-10-28T10:15:00Z' },
        ],
      },
    ] as T;
  }

  // Data exports & deletions
  if (path === '/me/data-exports' || path === '/me/deletion-requests') {
    return { success: true, message: '요청이 정상 접수되었습니다.' } as T;
  }

  return {} as T;
}
