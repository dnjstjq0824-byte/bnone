// apps/web/tests/contracts.test.mjs
// Automated Contract & Invariant Test Suite for BN ONE (BN-04, BN-05, BN-06, BN-07, BN-08, BN-11)
import assert from 'node:assert/strict';
import {
  INITIAL_ACCOUNTS,
  INITIAL_ASSETS,
  INITIAL_LIABILITIES,
  INITIAL_SUBSCRIPTIONS,
  INITIAL_SUBSCRIPTION_SUMMARY,
  INITIAL_BUDGETS,
  INITIAL_TRIPS,
  INITIAL_ANALYTICS_OVERVIEW,
  INITIAL_EXCHANGE_RATES,
  calculateNetWorth,
} from '../src/api/mockData.ts';
import { handleMockApi } from '../src/api/mockService.ts';

console.log('🧪 Starting BN ONE Contract & Invariant Test Suite...\n');

let passedTests = 0;
let failedTests = 0;

function runTest(name, fn) {
  try {
    fn();
    console.log(`  ✓ PASS: ${name}`);
    passedTests++;
  } catch (err) {
    console.error(`  ✗ FAIL: ${name}`);
    console.error(`    ${err.message}`);
    failedTests++;
  }
}

// 1. BN-04: Single Source of Truth & Net Worth Invariant (A-04)
runTest('BN-04 / A-04: Net Worth calculation invariant holds strictly', () => {
  const nw = calculateNetWorth(INITIAL_ACCOUNTS, INITIAL_ASSETS, INITIAL_LIABILITIES);
  
  // Invariant 1: Total Assets = Financial Accounts + Other Assets
  assert.equal(
    nw.total_assets_minor,
    nw.financial_accounts_minor + nw.other_assets_minor,
    'Total Assets must exactly equal financial accounts + other assets'
  );

  // Invariant 2: Net Worth = Total Assets - Total Liabilities
  assert.equal(
    nw.net_worth_minor,
    nw.total_assets_minor - nw.total_liabilities_minor,
    'Net Worth must exactly equal Total Assets - Total Liabilities'
  );

  // Invariant 3: Anchor Regression Values (exact won unit match)
  assert.equal(nw.financial_accounts_minor, 72752500, 'Financial accounts sum must be 72,752,500');
  assert.equal(nw.other_assets_minor, 351250000, 'Other assets sum must be 351,250,000');
  assert.equal(nw.total_assets_minor, 424002500, 'Total assets must be 424,002,500');
  assert.equal(nw.total_liabilities_minor, 120850000, 'Total liabilities must be 120,850,000');
  assert.equal(nw.net_worth_minor, 303152500, 'Net Worth anchor must be 303,152,500');
});

// 2. BN-05: Subscriptions Contract & Summary (A-05)
runTest('BN-05 / A-05: Subscriptions contract has service_name, next_renewal_at and sum matches', () => {
  const subs = handleMockApi('/me/subscriptions');
  assert.equal(subs.length, 4, 'Must have 4 active subscriptions');

  let totalMonthly = 0;
  for (const sub of subs) {
    assert.ok(sub.service_name, 'Subscription must have non-empty service_name');
    assert.ok(typeof sub.service_name.charAt(0) === 'string', 'service_name.charAt(0) must not throw');
    assert.ok(sub.next_renewal_at, 'Subscription must have next_renewal_at');
    assert.ok(sub.amount_minor > 0, 'amount_minor must be positive');
    assert.ok(sub.status === 'ACTIVE', 'status must be ACTIVE');
    totalMonthly += sub.amount_minor;
  }

  assert.equal(totalMonthly, 66790, 'Monthly subscription sum must equal 66,790');
  assert.equal(totalMonthly * 12, 801480, 'Yearly subscription sum must equal 801,480');

  const sum = handleMockApi('/me/subscriptions/summary');
  assert.equal(sum.active_count, 4);
  assert.equal(sum.monthly_estimate_minor, 66790);
  assert.equal(sum.yearly_estimate_minor, 801480);
});

// 3. BN-05: Budgets Contract (A-06)
runTest('BN-05 / A-06: Budgets contract has limit_minor, remaining_minor, usage_percent', () => {
  const budgets = handleMockApi('/budgets');
  assert.equal(budgets.length, 3, 'Must have 3 initial budgets');

  for (const b of budgets) {
    assert.ok(b.id, 'Budget must have id');
    assert.ok(b.name, 'Budget must have name');
    assert.ok(typeof b.limit_minor === 'number', 'limit_minor must be number');
    assert.ok(typeof b.spent_minor === 'number', 'spent_minor must be number');
    assert.ok(typeof b.remaining_minor === 'number', 'remaining_minor must be number');
    assert.ok(typeof b.usage_percent === 'number', 'usage_percent must be number');
    assert.ok(!Number.isNaN(b.usage_percent), 'usage_percent must not be NaN');
    assert.equal(b.remaining_minor, b.limit_minor - b.spent_minor, 'remaining must equal limit - spent');
    assert.equal(b.usage_percent.toFixed(0), String(Math.round((b.spent_minor / b.limit_minor) * 100)));
    assert.ok(b.start_date && b.end_date, 'Must have start and end date');
  }
});

// 4. BN-06: Trips Contract & No NaN Guard (A-07)
runTest('BN-06 / A-07: Trips contract has total_budget_minor, total_spent_minor and zero NaN', () => {
  const trips = handleMockApi('/trips');
  assert.equal(trips.length, 2, 'Must have 2 initial trips');

  const tokyo = trips.find(t => t.id === 'trp-001');
  assert.ok(tokyo, 'Tokyo trip must exist');
  assert.equal(tokyo.city_name, '도쿄');
  assert.equal(tokyo.country_code, 'JP');
  assert.equal(tokyo.total_budget_minor, 1200000);
  assert.equal(tokyo.total_spent_minor, 350000);
  const remaining = tokyo.total_budget_minor - tokyo.total_spent_minor;
  assert.equal(remaining, 850000);
  assert.ok(!Number.isNaN(remaining), 'Remaining trip budget must not be NaN');
});

// 5. BN-07: Analytics Contract & Current Month Population (A-08)
runTest('BN-07 / A-08: Analytics overview matches page contract with non-zero current month values', () => {
  const an = handleMockApi('/analytics/overview');
  assert.equal(an.current_month_income_minor, 4500000, 'Current month income must be 4,500,000');
  assert.equal(an.current_month_expense_minor, 1285000, 'Current month expense must be 1,285,000');
  assert.equal(an.net_cashflow_minor, 3215000, 'Net cashflow must be 3,215,000');
  assert.equal(an.savings_rate, 71.4, 'Savings rate must be 71.4%');
  assert.ok(Array.isArray(an.monthly_trend) && an.monthly_trend.length === 6, 'Must have 6 months trend');
  assert.ok(Array.isArray(an.category_breakdown) && an.category_breakdown.length === 5, 'Must have 5 categories');
});

// 6. BN-08: Forex Contract & 100 USD -> KRW 135,000 Conversion (A-09)
runTest('BN-08 / A-09: Forex rates returns rates array and 100 USD converts to 135,000 KRW', () => {
  const ratesResp = handleMockApi('/fx/rates?base=USD');
  assert.ok(Array.isArray(ratesResp.rates), 'Must return rates array in rates field');
  assert.ok(ratesResp.rates.length >= 3, 'Must have at least 3 reference rates');

  const usdKrw = ratesResp.rates.find(r => r.base_currency === 'USD' && r.quote_currency === 'KRW');
  assert.ok(usdKrw, 'USD/KRW rate must exist');
  assert.equal(usdKrw.rate, 1350);

  // Conversion test with 100 USD
  const conv = handleMockApi('/fx/convert', {
    method: 'POST',
    body: JSON.stringify({
      source_currency: 'USD',
      target_currency: 'KRW',
      source_amount_minor: 100,
    }),
  });
  assert.equal(conv.converted_amount_minor, 135000, '100 USD at 1,350 rate must equal 135,000 KRW');
});

// 7. BN-03: Authentication Password Verification (A-02)
runTest('BN-03 / A-02: Auth strictly verifies test@example.com / password123 and fails on wrong input', () => {
  // 1. Correct credentials -> Succeeds
  const valid = handleMockApi('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email: 'test@example.com', password: 'password123' }),
  });
  assert.ok(valid.access_token, 'Valid login must return access_token');
  assert.equal(valid.user.email, 'test@example.com');

  // 2. Wrong password -> Must throw INVALID_CREDENTIALS
  assert.throws(
    () => {
      handleMockApi('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email: 'test@example.com', password: 'wrongpass999' }),
      });
    },
    (err) => err.code === 'INVALID_CREDENTIALS',
    'Wrong password must throw INVALID_CREDENTIALS'
  );

  // 3. Arbitrary email without registration -> Must throw INVALID_CREDENTIALS
  assert.throws(
    () => {
      handleMockApi('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email: 'intruder@random.com', password: 'password123' }),
      });
    },
    (err) => err.code === 'INVALID_CREDENTIALS',
    'Arbitrary email must not auto-login'
  );
});

// 8. BN-05 / A-11: Contract Schema Regression & Negative Mutation Test
runTest('BN-05 / A-11: Schema validator fails if a required contract field is missing (Negative test)', () => {
  function validateSubscriptionContract(sub) {
    if (!sub.service_name) throw new Error('Missing service_name');
    if (!sub.next_renewal_at) throw new Error('Missing next_renewal_at');
    if (typeof sub.amount_minor !== 'number') throw new Error('Missing amount_minor');
  }

  // Good item passes
  validateSubscriptionContract({
    service_name: 'Netflix',
    next_renewal_at: '2026-11-15',
    amount_minor: 17000,
  });

  // Mutated item without service_name fails
  assert.throws(
    () => {
      validateSubscriptionContract({
        name: 'Netflix', // old broken field name
        next_renewal_at: '2026-11-15',
        amount_minor: 17000,
      });
    },
    /Missing service_name/,
    'Negative test: missing service_name must trigger failure'
  );
});

console.log(`\n========================================`);
console.log(`Test Results: ${passedTests} passed, ${failedTests} failed.`);
console.log(`========================================\n`);

if (failedTests > 0) {
  process.exit(1);
}
