import { defineStore } from 'pinia';
import { ref } from 'vue';
import { apiFetch } from '../api/client';

export interface FinanceAccount {
  id: string;
  name: string;
  account_type: string;
  currency_code: string;
  calculated_balance_minor: number;
  opening_balance_minor: number;
  color: string;
  icon: string;
  source_type: string;
}

export interface NetWorthData {
  base_currency: string;
  total_assets_minor: number;
  financial_accounts_minor: number;
  other_assets_minor: number;
  total_liabilities_minor: number;
  net_worth_minor: number;
  disclaimer: string;
}

export interface TransactionItem {
  id: string;
  transaction_type: string;
  occurred_at: string;
  description: string;
  merchant_name: string;
  amount_minor: number;
  currency_code: string;
  category_name?: string;
  account_id?: string;
}

export const useFinanceStore = defineStore('finance', () => {
  const accounts = ref<FinanceAccount[]>([]);
  const netWorth = ref<NetWorthData | null>(null);
  const recentTransactions = ref<TransactionItem[]>([]);
  const analyticsOverview = ref<any>(null);
  const isLoading = ref<boolean>(false);
  const dashboardError = ref('');
  const maskAmounts = ref<boolean>(false);

  async function loadDashboard() {
    isLoading.value = true;
    dashboardError.value = '';
    try {
      const [accs, nw, txs, overview] = await Promise.all([
        apiFetch<FinanceAccount[]>('/accounts'),
        apiFetch<NetWorthData>('/net-worth'),
        apiFetch<TransactionItem[]>('/transactions?limit=10'),
        apiFetch<any>('/analytics/overview'),
      ]);
      accounts.value = accs;
      netWorth.value = nw;
      recentTransactions.value = txs;
      analyticsOverview.value = overview;
    } catch {
      accounts.value = [];
      netWorth.value = null;
      recentTransactions.value = [];
      analyticsOverview.value = null;
      dashboardError.value = '금융 정보를 불러오지 못했어요. 연결을 확인하고 다시 시도해주세요.';
    } finally {
      isLoading.value = false;
    }
  }

  function formatCurrency(minor: number, currency: string = 'KRW'): string {
    if (maskAmounts.value) return '••••••';
    const amount = currency === 'KRW' || currency === 'JPY' || currency === 'VND' ? minor : minor / 100;
    return new Intl.NumberFormat('ko-KR', {
      style: 'currency',
      currency: currency,
      maximumFractionDigits: (currency === 'KRW' || currency === 'JPY') ? 0 : 2,
    }).format(amount);
  }

  return {
    accounts,
    netWorth,
    recentTransactions,
    analyticsOverview,
    isLoading,
    dashboardError,
    maskAmounts,
    loadDashboard,
    formatCurrency,
  };
});

