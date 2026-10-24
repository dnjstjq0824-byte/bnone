<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-black text-slate-900 tracking-tight">예산 관리</h1>
        <p class="text-xs text-slate-500 mt-0.5">월별 지출 목표를 설정하고 초과 지출을 방지합니다.</p>
      </div>
      <button
        @click="showCreateModal = true"
        class="px-4 py-2 bg-brand-600 hover:bg-brand-700 text-white rounded-xl text-xs font-bold shadow-md shadow-brand-600/20 transition-all"
      >
        + 새 예산 생성
      </button>
    </div>

    <!-- Budgets List -->
    <div v-if="budgets.length === 0" class="bg-white p-12 rounded-3xl border border-slate-200 card-shadow text-center">
      <div class="text-3xl mb-3">🎯</div>
      <div class="text-sm font-bold text-slate-700">설정된 예산이 없습니다.</div>
      <p class="text-xs text-slate-400 mt-1 mb-4">이번 달 식비나 생활비 예산을 설정해보세요.</p>
      <button @click="showCreateModal = true" class="px-4 py-2 bg-brand-600 text-white rounded-xl text-xs font-bold">
        첫 예산 만들기
      </button>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-5">
      <div
        v-for="b in budgets"
        :key="b.id"
        class="bg-white p-6 rounded-3xl border border-slate-200 card-shadow space-y-4"
      >
        <div class="flex items-center justify-between">
          <div>
            <h3 class="font-bold text-sm text-slate-900">{{ b.name }}</h3>
            <span class="text-[10px] text-slate-400">{{ b.start_date }} ~ {{ b.end_date }}</span>
          </div>
          <span
            class="px-2.5 py-0.5 rounded-full text-[10px] font-bold"
            :class="b.usage_percent >= 100 ? 'bg-red-50 text-red-600' : b.usage_percent >= 80 ? 'bg-amber-50 text-amber-600' : 'bg-emerald-50 text-emerald-600'"
          >
            {{ b.usage_percent.toFixed(1) }}% 사용
          </span>
        </div>

        <div>
          <div class="flex items-center justify-between text-xs mb-1.5 font-bold">
            <span class="text-slate-500">사용: {{ finance.formatCurrency(b.spent_minor, b.currency_code) }}</span>
            <span class="text-slate-900">한도: {{ finance.formatCurrency(b.limit_minor, b.currency_code) }}</span>
          </div>
          <div class="w-full bg-slate-100 h-2.5 rounded-full overflow-hidden">
            <div
              class="h-full rounded-full transition-all"
              :class="b.usage_percent >= 100 ? 'bg-red-500' : b.usage_percent >= 80 ? 'bg-amber-500' : 'bg-brand-600'"
              :style="{ width: `${Math.min(100, b.usage_percent)}%` }"
            ></div>
          </div>
        </div>

        <div class="flex items-center justify-between pt-3 border-t border-slate-100 text-[11px]">
          <span class="text-slate-400">
            남은 예산: <strong>{{ finance.formatCurrency(b.remaining_minor, b.currency_code) }}</strong>
          </span>
          <button @click="deleteBudget(b.id)" class="text-slate-400 hover:text-red-600 font-bold transition-colors">
            삭제
          </button>
        </div>
      </div>
    </div>

    <!-- Create Budget Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="bg-white max-w-md w-full rounded-3xl p-6 card-shadow space-y-4 border border-slate-200">
        <h3 class="text-lg font-bold text-slate-900">새 예산 생성</h3>
        <form @submit.prevent="submitCreateBudget" class="space-y-3">
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">예산 이름</label>
            <input v-model="newBudget.name" required placeholder="예: 11월 총 생활비" class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:border-brand-600" />
          </div>
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">한도 금액 (KRW)</label>
            <input v-model.number="newBudget.limit_minor" type="number" min="1000" required placeholder="1000000" class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm" />
          </div>
          <div class="grid grid-cols-2 gap-2">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">시작일</label>
              <input v-model="newBudget.start_date" type="date" required class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-xs bg-white" />
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">종료일</label>
              <input v-model="newBudget.end_date" type="date" required class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-xs bg-white" />
            </div>
          </div>
          <div class="flex gap-2 pt-2">
            <button type="button" @click="showCreateModal = false" class="w-1/2 py-2.5 bg-slate-100 hover:bg-slate-200 rounded-xl text-xs font-bold text-slate-700">취소</button>
            <button type="submit" class="w-1/2 py-2.5 bg-brand-600 hover:bg-brand-700 text-white rounded-xl text-xs font-bold shadow-md shadow-brand-600/25">생성하기</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { apiFetch } from '../api/client';
import { useFinanceStore } from '../stores/finance';

const finance = useFinanceStore();
const budgets = ref<any[]>([]);
const showCreateModal = ref(false);

const newBudget = ref({
  name: '',
  limit_minor: 1000000,
  start_date: new Date().toISOString().slice(0, 7) + '-01',
  end_date: new Date(new Date().getFullYear(), new Date().getMonth() + 1, 0).toISOString().slice(0, 10),
});

async function loadBudgets() {
  budgets.value = await apiFetch('/budgets');
}

onMounted(() => {
  loadBudgets();
});

async function submitCreateBudget() {
  await apiFetch('/budgets', {
    method: 'POST',
    body: JSON.stringify(newBudget.value),
  });
  showCreateModal.value = false;
  newBudget.value.name = '';
  await loadBudgets();
}

async function deleteBudget(id: string) {
  if (confirm('이 예산을 삭제하시겠습니까?')) {
    await apiFetch(`/budgets/${id}`, { method: 'DELETE' });
    await loadBudgets();
  }
}
</script>

