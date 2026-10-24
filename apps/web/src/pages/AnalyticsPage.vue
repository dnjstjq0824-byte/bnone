<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-black text-slate-900 tracking-tight">소비 & 현금흐름 분석</h1>
      <p class="text-xs text-slate-500 mt-0.5">월별 수입/지출 추이와 카테고리별 비중을 분석합니다.</p>
    </div>

    <!-- Monthly Cashflow Card -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
      <div class="p-6 rounded-3xl bg-white border border-slate-200 card-shadow">
        <span class="text-xs font-bold text-slate-400">이번 달 수입</span>
        <div class="text-2xl font-black text-emerald-600 mt-2">
          + {{ finance.formatCurrency(overview?.current_month_income_minor || 0) }}
        </div>
      </div>
      <div class="p-6 rounded-3xl bg-white border border-slate-200 card-shadow">
        <span class="text-xs font-bold text-slate-400">이번 달 지출</span>
        <div class="text-2xl font-black text-red-600 mt-2">
          − {{ finance.formatCurrency(overview?.current_month_expense_minor || 0) }}
        </div>
      </div>
      <div class="p-6 rounded-3xl bg-white border border-slate-200 card-shadow">
        <span class="text-xs font-bold text-slate-400">순 현금흐름</span>
        <div class="text-2xl font-black text-slate-900 mt-2">
          {{ finance.formatCurrency(overview?.net_cashflow_minor || 0) }}
        </div>
      </div>
    </div>

    <!-- 6-Month Trend Table & Visualization -->
    <div class="bg-white p-6 rounded-3xl border border-slate-200 card-shadow space-y-4">
      <h2 class="text-base font-bold text-slate-900">최근 6개월 현금흐름 추이</h2>
      <div class="space-y-4">
        <div v-for="t in overview?.monthly_trend || []" :key="t.month" class="space-y-1.5">
          <div class="flex items-center justify-between text-xs">
            <span class="font-bold text-slate-700">{{ t.month }}</span>
            <div class="space-x-3 text-[11px]">
              <span class="text-emerald-600 font-bold">수입 +{{ finance.formatCurrency(t.income_minor) }}</span>
              <span class="text-red-600 font-bold">지출 −{{ finance.formatCurrency(t.expense_minor) }}</span>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-2 h-3">
            <div class="bg-slate-100 rounded-full overflow-hidden flex justify-end">
              <div class="bg-emerald-500 h-full rounded-full" :style="{ width: `${getBarWidth(t.income_minor)}%` }"></div>
            </div>
            <div class="bg-slate-100 rounded-full overflow-hidden">
              <div class="bg-red-500 h-full rounded-full" :style="{ width: `${getBarWidth(t.expense_minor)}%` }"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Category Spending Table -->
    <div class="bg-white p-6 rounded-3xl border border-slate-200 card-shadow space-y-4">
      <h2 class="text-base font-bold text-slate-900">카테고리별 지출 비중</h2>
      <div v-if="!overview?.category_breakdown?.length" class="text-center py-8 text-xs text-slate-400">
        지출 내역이 없습니다.
      </div>
      <div v-else class="divide-y divide-slate-100">
        <div v-for="c in overview.category_breakdown" :key="c.category_id" class="py-3 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <span class="w-3.5 h-3.5 rounded-full" :style="{ backgroundColor: c.color || '#3B82F6' }"></span>
            <span class="text-xs font-bold text-slate-800">{{ c.category_name }}</span>
          </div>
          <div class="text-right">
            <div class="text-xs font-black text-slate-900">{{ finance.formatCurrency(c.amount_minor) }}</div>
            <div class="text-[10px] text-slate-400">{{ c.percentage.toFixed(1) }}%</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { apiFetch } from '../api/client';
import { useFinanceStore } from '../stores/finance';

const finance = useFinanceStore();
const overview = ref<any>(null);

onMounted(async () => {
  overview.value = await apiFetch('/analytics/overview');
});

function getBarWidth(val: number): number {
  if (!val) return 0;
  return Math.min(100, Math.max(5, (val / 3000000) * 100));
}
</script>

