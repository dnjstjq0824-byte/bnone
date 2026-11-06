<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-[#141F19]">예산 계획</h1>
        <p class="text-sm text-[#8A948C] mt-1">월별 지출 목표를 세우고 남은 예산을 관리해요.</p>
      </div>
      <button
        @click="showCreateModal = true"
        class="inline-flex items-center gap-1.5 px-4 py-2.5 bg-[#1E5A38] hover:bg-[#17452C] text-white rounded-xl text-sm font-bold shadow-sm transition-colors"
      >
        <span>+ 새 예산 세우기</span>
      </button>
    </div>

    <!-- Error alert -->
    <div v-if="loadError" class="p-4 rounded-2xl bg-[#FBEAEA] text-[#D64545] text-xs font-medium">
      {{ loadError }}
      <button class="underline ml-2 font-bold" @click="loadBudgets">다시 시도</button>
    </div>

    <!-- Budgets List -->
    <div v-if="budgets.length === 0" class="bg-white p-12 rounded-3xl text-center">
      <div class="w-14 h-14 mx-auto mb-4 rounded-2xl bg-[#E4EFE7] flex items-center justify-center text-[#2A7248] text-2xl">
        🎯
      </div>
      <div class="text-base font-bold text-[#141F19]">설정된 예산이 없습니다.</div>
      <p class="text-xs text-[#8A948C] mt-1 mb-5">이번 달 생활비나 식비 목표를 정해보세요.</p>
      <button
        @click="showCreateModal = true"
        class="px-5 py-2.5 bg-[#1E5A38] text-white rounded-xl text-xs font-bold"
      >
        첫 예산 만들기
      </button>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-5">
      <div
        v-for="b in budgets"
        :key="b.id"
        class="bg-white p-6 rounded-3xl space-y-4"
      >
        <div class="flex items-center justify-between">
          <div>
            <h3 class="font-bold text-base text-[#141F19]">{{ b.name }}</h3>
            <span class="text-xs text-[#8A948C] mt-0.5 block">{{ b.start_date || '-' }} ~ {{ b.end_date || '-' }}</span>
          </div>
          <span
            class="px-2.5 py-1 rounded-lg text-xs font-bold"
            :class="(b.usage_percent ?? 0) >= 100 ? 'bg-[#FBEAEA] text-[#D64545]' : (b.usage_percent ?? 0) >= 80 ? 'bg-[#FFF3E2] text-[#FF8D00]' : 'bg-[#EBF2EC] text-[#1E5A38]'"
          >
            {{ (b.usage_percent ?? 0).toFixed(0) }}% 사용
          </span>
        </div>

        <div>
          <div class="text-xl font-extrabold text-[#141F19] mb-2 tracking-tight tabular-nums">
            남은 예산 {{ finance.formatCurrency(b.remaining_minor ?? ((b.limit_minor ?? 0) - (b.spent_minor ?? 0)), b.currency_code) }}
          </div>
          <div class="w-full bg-[#ECEFE8] h-3 rounded-full overflow-hidden">
            <div
              class="h-full rounded-full transition-all"
              :class="(b.usage_percent ?? 0) >= 100 ? 'bg-[#D64545]' : (b.usage_percent ?? 0) >= 80 ? 'bg-[#FF8D00]' : 'bg-[#1E5A38]'"
              :style="{ width: `${Math.min(100, b.usage_percent ?? 0)}%` }"
            ></div>
          </div>
          <div class="flex items-center justify-between text-xs text-[#8A948C] mt-2">
            <span>쓴 돈: {{ finance.formatCurrency(b.spent_minor ?? 0, b.currency_code) }}</span>
            <span>목표 한도: {{ finance.formatCurrency(b.limit_minor ?? b.amount_minor ?? 0, b.currency_code) }}</span>
          </div>
        </div>

        <div class="flex items-center justify-end pt-3 border-t border-[#ECEFE8] text-xs">
          <button @click="deleteBudget(b.id)" class="text-[#8A948C] hover:text-[#D64545] font-semibold transition-colors">
            예산 삭제
          </button>
        </div>
      </div>
    </div>

    <!-- Create Budget Modal (Toss Rounded Dialog Style) -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="bg-white max-w-md w-full rounded-3xl p-6 space-y-5">
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-bold text-[#141F19]">새 목표 예산 세우기</h3>
          <button @click="showCreateModal = false" class="text-[#8A948C] hover:text-[#141F19]">✕</button>
        </div>
        <form @submit.prevent="submitCreateBudget" class="space-y-4">
          <div>
            <label class="block text-xs font-bold text-[#41524A] mb-1.5">예산 이름</label>
            <input
              v-model="newBudget.name"
              required
              placeholder="예: 이번 달 생활비, 식비 예산"
              class="w-full px-4 py-3 rounded-2xl bg-[#F6F8F4] border border-[#E2E7DE] text-sm text-[#141F19] focus:outline-none focus:border-[#1E5A38]"
            />
          </div>
          <div>
            <label class="block text-xs font-bold text-[#41524A] mb-1.5">한도 금액 (KRW)</label>
            <input
              v-model.number="newBudget.limit_minor"
              type="number"
              min="1000"
              required
              placeholder="1000000"
              class="w-full px-4 py-3 rounded-2xl bg-[#F6F8F4] border border-[#E2E7DE] text-sm text-[#141F19] focus:outline-none focus:border-[#1E5A38]"
            />
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-bold text-[#41524A] mb-1.5">시작일</label>
              <input
                v-model="newBudget.start_date"
                type="date"
                required
                class="w-full px-4 py-3 rounded-2xl bg-[#F6F8F4] border border-[#E2E7DE] text-xs text-[#141F19] focus:outline-none focus:border-[#1E5A38]"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-[#41524A] mb-1.5">종료일</label>
              <input
                v-model="newBudget.end_date"
                type="date"
                required
                class="w-full px-4 py-3 rounded-2xl bg-[#F6F8F4] border border-[#E2E7DE] text-xs text-[#141F19] focus:outline-none focus:border-[#1E5A38]"
              />
            </div>
          </div>
          <div class="flex gap-3 pt-2">
            <button
              type="button"
              @click="showCreateModal = false"
              class="w-1/2 py-3.5 bg-[#ECEFE8] hover:bg-[#E2E7DE] rounded-2xl text-sm font-bold text-[#41524A]"
            >
              취소
            </button>
            <button
              type="submit"
              class="w-1/2 py-3.5 bg-[#1E5A38] hover:bg-[#17452C] text-white rounded-2xl text-sm font-bold shadow-sm"
            >
              예산 만들기
            </button>
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
const loadError = ref('');

const newBudget = ref({
  name: '',
  limit_minor: 1000000,
  start_date: new Date().toISOString().slice(0, 7) + '-01',
  end_date: new Date(new Date().getFullYear(), new Date().getMonth() + 1, 0).toISOString().slice(0, 10),
});

async function loadBudgets() {
  loadError.value = '';
  try {
    const res = await apiFetch('/budgets');
    budgets.value = Array.isArray(res) ? res : [];
  } catch (err: any) {
    loadError.value = err.message || '예산 정보를 불러오지 못했습니다.';
  }
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
