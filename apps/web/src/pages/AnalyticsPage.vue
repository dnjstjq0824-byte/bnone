<template>
  <div class="space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold tracking-tight text-[#191F28]">소비 & 현금흐름 분석</h1>
      <p class="text-sm text-[#8B95A1] mt-1">월별 수입/지출 추이와 카테고리별 비중을 분석해요.</p>
    </div>

    <!-- Monthly Cashflow Cards (Toss Stats Cards) -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
      <div class="p-6 rounded-3xl bg-white space-y-2">
        <span class="text-xs font-bold text-[#8B95A1]">이번 달 들어온 돈</span>
        <div class="text-2xl font-extrabold text-[#3182F6] tracking-tight tabular-nums">
          + {{ finance.formatCurrency(overview?.current_month_income_minor || 0) }}
        </div>
      </div>
      <div class="p-6 rounded-3xl bg-white space-y-2">
        <span class="text-xs font-bold text-[#8B95A1]">이번 달 쓴 금액</span>
        <div class="text-2xl font-extrabold text-[#191F28] tracking-tight tabular-nums">
          − {{ finance.formatCurrency(overview?.current_month_expense_minor || 0) }}
        </div>
      </div>
      <div class="p-6 rounded-3xl bg-white space-y-2">
        <span class="text-xs font-bold text-[#8B95A1]">이번 달 남은 순현금흐름</span>
        <div
          class="text-2xl font-extrabold tracking-tight tabular-nums"
          :class="(overview?.net_cashflow_minor || 0) >= 0 ? 'text-[#3182F6]' : 'text-[#F04452]'"
        >
          {{ finance.formatCurrency(overview?.net_cashflow_minor || 0) }}
        </div>
      </div>
    </div>

    <!-- 6-Month Trend Card -->
    <div class="bg-white p-6 rounded-3xl space-y-5">
      <h2 class="text-base font-bold text-[#191F28]">최근 6개월 수입 & 소비 추이</h2>
      <div class="space-y-5">
        <div v-for="t in overview?.monthly_trend || []" :key="t.month" class="space-y-2">
          <div class="flex items-center justify-between text-xs">
            <span class="font-bold text-[#191F28]">{{ t.month }}</span>
            <div class="space-x-3 text-xs">
              <span class="text-[#3182F6] font-bold">수입 +{{ finance.formatCurrency(t.income_minor) }}</span>
              <span class="text-[#191F28] font-bold">지출 −{{ finance.formatCurrency(t.expense_minor) }}</span>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-2 h-3.5">
            <div class="bg-[#F2F4F6] rounded-full overflow-hidden flex justify-end">
              <div class="bg-[#3182F6] h-full rounded-full transition-all" :style="{ width: `${getBarWidth(t.income_minor)}%` }"></div>
            </div>
            <div class="bg-[#F2F4F6] rounded-full overflow-hidden">
              <div class="bg-[#F04452] h-full rounded-full transition-all" :style="{ width: `${getBarWidth(t.expense_minor)}%` }"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Category Spending Card -->
    <div class="bg-white p-6 rounded-3xl space-y-4">
      <h2 class="text-base font-bold text-[#191F28]">카테고리별 지출 비중</h2>
      <div v-if="!overview?.category_breakdown?.length" class="text-center py-8 text-xs text-[#8B95A1]">
        지출 내역이 없습니다.
      </div>
      <div v-else class="divide-y divide-[#F2F4F6]">
        <div v-for="c in overview.category_breakdown" :key="c.category_id" class="py-3.5 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <span class="w-3.5 h-3.5 rounded-full flex-shrink-0" :style="{ backgroundColor: c.color || '#3182F6' }"></span>
            <span class="text-sm font-bold text-[#191F28]">{{ c.category_name }}</span>
          </div>
          <div class="text-right">
            <div class="text-sm font-extrabold text-[#191F28] tabular-nums">{{ finance.formatCurrency(c.amount_minor) }}</div>
            <div class="text-xs text-[#8B95A1]">{{ c.percentage.toFixed(1) }}%</div>
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
