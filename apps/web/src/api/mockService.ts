// apps/web/src/api/mockService.ts
import {
  MOCK_USER,
  INITIAL_ACCOUNTS,
  INITIAL_NET_WORTH,
  INITIAL_TRANSACTIONS,
  INITIAL_ANALYTICS_OVERVIEW,
  INITIAL_ASSETS,
  INITIAL_LIABILITIES,
  INITIAL_BUDGETS,
  INITIAL_SUBSCRIPTIONS,
  INITIAL_SUBSCRIPTION_SUMMARY,
  INITIAL_TRIPS,
  INITIAL_CURRENCIES,
  INITIAL_EXCHANGE_RATES,
  INITIAL_CATEGORIES,
  INITIAL_NOTIFICATIONS,
} from './mockData';

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

  // Auth: Login
  if (path === '/auth/login' && method === 'POST') {
    const user = {
      ...MOCK_USER,
      email: body.email || 'test@example.com',
      display_name: body.email === 'test@example.com' ? '토스테스터' : (body.email.split('@')[0] || '사용자'),
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
    const user = {
      ...MOCK_USER,
      id: `u-${Date.now()}`,
      email: body.email,
      display_name: body.display_name || body.email.split('@')[0],
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

  // Net Worth
  if (path === '/net-worth') {
    const accounts = getStorage('accounts', INITIAL_ACCOUNTS);
    const assets = getStorage('assets', INITIAL_ASSETS);
    const liabilities = getStorage('liabilities', INITIAL_LIABILITIES);

    const finTotal = accounts.reduce((sum: number, a: any) => {
      // Rough conversion for demonstration
      if (a.currency_code === 'USD') return sum + (a.calculated_balance_minor / 100) * 1350;
      if (a.currency_code === 'JPY') return sum + a.calculated_balance_minor * 9;
      return sum + a.calculated_balance_minor;
    }, 0);

    const assetTotal = assets.reduce((sum: number, a: any) => sum + a.estimated_value_minor, 0);
    const liaTotal = liabilities.reduce((sum: number, l: any) => sum + l.outstanding_amount_minor, 0);
    const totalAssets = finTotal + assetTotal;

    return {
      base_currency: 'KRW',
      total_assets_minor: totalAssets,
      financial_accounts_minor: finTotal,
      other_assets_minor: assetTotal,
      total_liabilities_minor: liaTotal,
      net_worth_minor: totalAssets - liaTotal,
      disclaimer: '최신 환율(한국은행, ECB 기준)로 실시간 환산된 순자산입니다.',
    } as T;
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

  // Analytics
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

  // Budgets
  if (path === '/budgets') {
    const budgets = getStorage('budgets', INITIAL_BUDGETS);
    if (method === 'POST') {
      const newBgt = {
        id: `bgt-${Date.now()}`,
        name: body.name || '새 예산',
        category_id: body.category_id,
        category_name: body.category_name || '식비 & 카페',
        period: body.period || 'MONTHLY',
        amount_minor: Number(body.amount_minor || 0),
        spent_minor: 0,
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

  // Subscriptions
  if (path === '/me/subscriptions') {
    const subs = getStorage('subscriptions', INITIAL_SUBSCRIPTIONS);
    if (method === 'POST') {
      const newSub = {
        id: `sub-${Date.now()}`,
        name: body.name || '새 구독',
        category: body.category || 'ENTERTAINMENT',
        amount_minor: Number(body.amount_minor || 0),
        currency_code: body.currency_code || 'KRW',
        billing_cycle: body.billing_cycle || 'MONTHLY',
        next_billing_date: body.next_billing_date || '2026-11-30',
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
    return { success: true } as T;
  }

  // Trips
  if (path === '/trips') {
    const trips = getStorage('trips', INITIAL_TRIPS);
    if (method === 'POST') {
      const newTrip = {
        id: `trp-${Date.now()}`,
        name: body.name || '새 여행',
        destination: body.destination || '여행지',
        start_date: body.start_date || '2026-12-01',
        end_date: body.end_date || '2026-12-05',
        budget_minor: Number(body.budget_minor || 0),
        spent_minor: 0,
        currency_code: body.currency_code || 'KRW',
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

  // Currencies & FX
  if (path === '/currencies') {
    return INITIAL_CURRENCIES as T;
  }

  if (path === '/fx/rates') {
    return INITIAL_EXCHANGE_RATES as T;
  }

  if (path === '/fx/convert' && method === 'POST') {
    const from = body.source_currency;
    const to = body.target_currency;
    const amount = Number(body.amount_minor || 0);
    let rate = 1350.0;
    if (from === 'USD' && to === 'KRW') rate = 1350.0;
    else if (from === 'KRW' && to === 'USD') rate = 0.000741;
    else if (from === 'JPY' && to === 'KRW') rate = 9.0;
    else if (from === 'EUR' && to === 'KRW') rate = 1460.0;

    const converted = Math.round(amount * rate);
    return {
      source_currency: from,
      target_currency: to,
      original_amount_minor: amount,
      converted_amount_minor: converted,
      rate,
      effective_at: new Date().toISOString(),
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
        created_at: new Date(Date.now() - 3600000).toISOString(),
        last_active_at: new Date().toISOString(),
        is_current: true,
      },
      {
        id: 'sess-mobile',
        platform: 'MOBILE (BN ONE iOS App)',
        device_id: 'iPhone 15 Pro',
        ip_address: '118.235.8.19',
        created_at: new Date(Date.now() - 86400000 * 2).toISOString(),
        last_active_at: new Date(Date.now() - 7200000).toISOString(),
        is_current: false,
      },
    ] as T;
  }

  if (path.startsWith('/me/sessions/') && method === 'DELETE') {
    return { success: true } as T;
  }

  // Support
  if (path === '/me/support-tickets') {
    return [
      {
        id: 'tkt-001',
        title: '다통화 환전 환율 갱신 주기 문의',
        status: 'ANSWERED',
        created_at: '2026-10-28T10:00:00Z',
        messages: [
          { sender: 'USER', message: '한국은행 기준 환율은 하루에 몇 번 갱신되나요?', sent_at: '2026-10-28T10:00:00Z' },
          { sender: 'AGENT', message: '영업일 기준 매일 오전 9시 30분에 자동 동기화됩니다.', sent_at: '2026-10-28T10:15:00Z' },
        ],
      },
    ] as T;
  }

  // Data exports & deletions
  if (path === '/me/data-exports' || path === '/me/deletion-requests') {
    return { success: true, message: '요청이 정상 접수되었습니다.' } as T;
  }

  // Default fallback
  return {} as T;
}
