<template>
  <div class="space-y-6">
    <!-- Welcome Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-slate-900 tracking-tight">
          안녕하세요, {{ auth.user?.display_name }}님 👋
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">내 돈의 흐름을 한눈에 확인하세요.</p>
      </div>
      <div class="flex items-center gap-2">
        <RouterLink
          to="/app/transactions/new"
          class="px-4 py-2.5 bg-brand-600 hover:bg-brand-700 text-white rounded-xl text-xs font-bold shadow-md shadow-brand-600/20 transition-all flex items-center gap-1.5"
        >
          <span>+</span>
          <span>새 거래 기록</span>
        </RouterLink>
      </div>
    </div>

    <!-- Mandatory Legal Banner (Article 1.5) -->
    <div class="p-3.5 rounded-2xl bg-blue-50/70 border border-blue-100 flex items-start gap-3">
      <span class="text-blue-600 text-base leading-none">ℹ️</span>
      <p class="text-xs text-blue-900 font-medium leading-relaxed">
        <strong>알림:</strong> 표시된 자산과 잔액은 사용자가 직접 입력한 기록을 기준으로 합니다. BN ONE의 거래 기록은 실제 금융기관 거래를 실행하지 않습니다.
      </p>
    </div>

    <!-- KPI Grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
      <!-- Net Worth Card -->
      <div class="p-6 rounded-3xl bg-white border border-slate-200 card-shadow flex flex-col justify-between">
        <div class="flex items-center justify-between mb-4">
          <span class="text-xs font-bold text-slate-400">내 순자산</span>
          <span class="text-xs px-2 py-0.5 rounded-full bg-slate-100 text-slate-600 font-bold">
            {{ finance.netWorth?.base_currency || 'KRW' }}
          </span>
        </div>
        <div>
          <div class="text-3xl font-black text-slate-900 tracking-tight">
            {{ finance.formatCurrency(finance.netWorth?.net_worth_minor || 0, finance.netWorth?.base_currency) }}
          </div>
          <div class="text-[11px] text-slate-400 mt-2 flex items-center justify-between">
            <span>총자산: {{ finance.formatCurrency(finance.netWorth?.total_assets_minor || 0, finance.netWorth?.base_currency) }}</span>
            <span>부채: {{ finance.formatCurrency(finance.netWorth?.total_liabilities_minor || 0, finance.netWorth?.base_currency) }}</span>
          </div>
        </div>
      </div>

      <!-- Month Spending -->
      <div class="p-6 rounded-3xl bg-white border border-slate-200 card-shadow flex flex-col justify-between">
        <div class="flex items-center justify-between mb-4">
          <span class="text-xs font-bold text-slate-400">이번 달 총 지출</span>
          <span class="text-xs text-red-600 bg-red-50 font-bold px-2 py-0.5 rounded-full">지출</span>
        </div>
        <div>
          <div class="text-3xl font-black text-slate-900 tracking-tight text-red-600">
            − {{ finance.formatCurrency(finance.analyticsOverview?.current_month_expense_minor || 0, finance.netWorth?.base_currency) }}
          </div>
          <div class="text-[11px] text-slate-400 mt-2">
            지출 내역 {{ finance.analyticsOverview?.category_breakdown?.length || 0 }}개 카테고리
          </div>
        </div>
      </div>

      <!-- Month Income -->
      <div class="p-6 rounded-3xl bg-white border border-slate-200 card-shadow flex flex-col justify-between">
        <div class="flex items-center justify-between mb-4">
          <span class="text-xs font-bold text-slate-400">이번 달 총 수입</span>
          <span class="text-xs text-emerald-600 bg-emerald-50 font-bold px-2 py-0.5 rounded-full">수입</span>
        </div>
        <div>
          <div class="text-3xl font-black text-emerald-600 tracking-tight">
            + {{ finance.formatCurrency(finance.analyticsOverview?.current_month_income_minor || 0, finance.netWorth?.base_currency) }}
          </div>
          <div class="text-[11px] text-slate-400 mt-2">
            순현금흐름: {{ finance.formatCurrency(finance.analyticsOverview?.net_cashflow_minor || 0, finance.netWorth?.base_currency) }}
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content Grid (Two Columns) -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Left 2 Cols: Accounts & Category Breakdown -->
      <div class="lg:col-span-2 space-y-6">
        <!-- Accounts Card List -->
        <div class="bg-white p-6 rounded-3xl border border-slate-200 card-shadow">
          <div class="flex items-center justify-between mb-5">
            <h2 class="text-base font-bold text-slate-900">내 금융 계정</h2>
            <RouterLink to="/app/assets" class="text-xs font-bold text-brand-600 hover:underline">
              계정 전체보기 →
            </RouterLink>
          </div>

          <div v-if="finance.accounts.length === 0" class="text-center py-8 text-xs text-slate-400">
            등록된 계정이 없습니다. 계정을 추가해보세요.
          </div>
          <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
            <div
              v-for="acc in finance.accounts"
              :key="acc.id"
              class="p-4 rounded-2xl border border-slate-100 bg-slate-50/50 hover:bg-slate-50 transition-colors"
            >
              <div class="flex items-center justify-between mb-2">
                <div class="flex items-center gap-2">
                  <span class="w-3 h-3 rounded-full" :style="{ backgroundColor: acc.color }"></span>
                  <span class="font-bold text-xs text-slate-800">{{ acc.name }}</span>
                </div>
                <span class="text-[10px] text-slate-400 font-semibold px-2 py-0.5 rounded bg-white border border-slate-200">
                  {{ acc.currency_code }}
                </span>
              </div>
              <div class="text-lg font-black text-slate-900 mt-1">
                {{ finance.formatCurrency(acc.calculated_balance_minor, acc.currency_code) }}
              </div>
            </div>
          </div>
        </div>

        <!-- Recent Transactions -->
        <div class="bg-white p-6 rounded-3xl border border-slate-200 card-shadow">
          <div class="flex items-center justify-between mb-5">
            <h2 class="text-base font-bold text-slate-900">최근 거래 내역</h2>
            <RouterLink to="/app/transactions" class="text-xs font-bold text-brand-600 hover:underline">
              전체 내역보기 →
            </RouterLink>
          </div>

          <div v-if="finance.recentTransactions.length === 0" class="text-center py-8 text-xs text-slate-400">
            아직 기록된 거래가 없습니다. '+ 새 거래 기록' 버튼을 눌러보세요.
          </div>
          <div v-else class="divide-y divide-slate-100">
            <div
              v-for="tx in finance.recentTransactions"
              :key="tx.id"
              class="py-3 flex items-center justify-between"
            >
              <div class="flex items-center gap-3">
                <div
                  class="w-9 h-9 rounded-xl flex items-center justify-center text-sm font-bold"
                  :class="tx.transaction_type === 'EXPENSE' ? 'bg-red-50 text-red-600' : tx.transaction_type === 'INCOME' ? 'bg-emerald-50 text-emerald-600' : 'bg-indigo-50 text-indigo-600'"
                >
                  {{ tx.transaction_type === 'EXPENSE' ? '↓' : tx.transaction_type === 'INCOME' ? '↑' : '⇄' }}
                </div>
                <div>
                  <div class="text-xs font-bold text-slate-800">
                    {{ tx.description || tx.merchant_name || '거래' }}
                  </div>
                  <div class="text-[10px] text-slate-400">
                    {{ formatDate(tx.occurred_at) }} • {{ tx.category_name || tx.transaction_type }}
                  </div>
                </div>
              </div>
              <div
                class="text-xs font-black tracking-tight"
                :class="tx.transaction_type === 'EXPENSE' ? 'text-slate-900' : tx.transaction_type === 'INCOME' ? 'text-emerald-600' : 'text-indigo-600'"
              >
                {{ tx.transaction_type === 'EXPENSE' ? '− ' : tx.transaction_type === 'INCOME' ? '+ ' : '' }}
                {{ finance.formatCurrency(tx.amount_minor, tx.currency_code) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right 1 Col: Category Breakdown & Quick Info -->
      <div class="space-y-6">
        <!-- Category Breakdown -->
        <div class="bg-white p-6 rounded-3xl border border-slate-200 card-shadow">
          <h2 class="text-base font-bold text-slate-900 mb-4">카테고리별 지출</h2>
          
          <div v-if="!finance.analyticsOverview?.category_breakdown?.length" class="text-center py-6 text-xs text-slate-400">
            이번 달 지출 내역이 없습니다.
          </div>
          <div v-else class="space-y-3">
            <div
              v-for="cat in finance.analyticsOverview.category_breakdown.slice(0, 5)"
              :key="cat.category_id"
              class="space-y-1"
            >
              <div class="flex items-center justify-between text-xs">
                <span class="font-bold text-slate-700">{{ cat.category_name }}</span>
                <span class="font-bold text-slate-900">
                  {{ finance.formatCurrency(cat.amount_minor, finance.netWorth?.base_currency) }}
                  <span class="text-[10px] text-slate-400 font-normal">({{ cat.percentage.toFixed(1) }}%)</span>
                </span>
              </div>
              <div class="w-full bg-slate-100 h-2 rounded-full overflow-hidden">
                <div class="h-full rounded-full transition-all" :style="{ width: `${cat.percentage}%`, backgroundColor: cat.color || '#3B82F6' }"></div>
              </div>
            </div>
          </div>
        </div>

        <!-- Quick Travel & Subs Widget -->
        <div class="bg-gradient-to-br from-indigo-900 to-slate-900 p-6 rounded-3xl text-white shadow-xl">
          <div class="flex items-center justify-between mb-4">
            <span class="text-xs font-bold text-indigo-300">여행 & 환율 도우미</span>
            <span class="text-base">✈️</span>
          </div>
          <p class="text-xs text-slate-300 mb-4 leading-relaxed">
            해외 여행 경비를 원본 통화로 기록하고, 공인 참고 환율로 실시간 예산 초과 여부를 확인하세요.
          </p>
          <div class="grid grid-cols-2 gap-2">
            <RouterLink to="/app/trips" class="py-2.5 px-3 bg-white/10 hover:bg-white/20 text-center rounded-xl text-xs font-bold transition-colors">
              여행 플랜 관리
            </RouterLink>
            <RouterLink to="/app/fx" class="py-2.5 px-3 bg-white/10 hover:bg-white/20 text-center rounded-xl text-xs font-bold transition-colors">
              참고 환율 계산기
            </RouterLink>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { RouterLink } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { useFinanceStore } from '../stores/finance';

const auth = useAuthStore();
const finance = useFinanceStore();

onMounted(async () => {
  await finance.loadDashboard();
});

function formatDate(iso: string) {
  if (!iso) return '';
  const d = new Date(iso);
  return `${d.getMonth() + 1}월 ${d.getDate()}일`;
}
</script>
