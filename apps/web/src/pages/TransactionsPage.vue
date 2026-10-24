<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-slate-900 tracking-tight">거래 내역 장부</h1>
        <p class="text-xs text-slate-500 mt-0.5">수입, 지출 및 계정 간 이동 기록을 검색하고 관리합니다.</p>
      </div>
      <div class="flex items-center gap-2">
        <button
          @click="exportCSV"
          class="px-3.5 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl text-xs font-bold transition-all flex items-center gap-1.5"
        >
          <span>📥</span>
          <span>CSV 내보내기</span>
        </button>
        <RouterLink
          to="/app/transactions/new"
          class="px-4 py-2 bg-brand-600 hover:bg-brand-700 text-white rounded-xl text-xs font-bold shadow-md shadow-brand-600/20 transition-all flex items-center gap-1.5"
        >
          <span>+</span>
          <span>거래 추가</span>
        </RouterLink>
      </div>
    </div>

    <!-- Filters & Search -->
    <div class="bg-white p-4 rounded-2xl border border-slate-200 card-shadow flex flex-wrap items-center gap-3">
      <div class="flex-1 min-w-[200px]">
        <input
          v-model="searchQuery"
          @input="fetchTransactions"
          type="text"
          placeholder="상호명, 설명 검색..."
          class="w-full px-3.5 py-2 rounded-xl border border-slate-200 text-xs focus:outline-none focus:border-brand-600"
        />
      </div>

      <div class="flex items-center gap-2">
        <select
          v-model="filterType"
          @change="fetchTransactions"
          class="px-3 py-2 rounded-xl border border-slate-200 text-xs bg-white font-medium text-slate-700"
        >
          <option value="">전체 종류</option>
          <option value="EXPENSE">지출</option>
          <option value="INCOME">수입</option>
          <option value="TRANSFER">계정 이동</option>
        </select>
      </div>
    </div>

    <!-- Mandatory Notice -->
    <div class="text-[11px] text-slate-400 px-1">
      * BN ONE의 거래 기록은 실제 금융기관 거래를 실행하지 않습니다.
    </div>

    <!-- Transaction Table -->
    <div class="bg-white rounded-3xl border border-slate-200 card-shadow overflow-hidden">
      <div v-if="loading" class="text-center py-12 text-xs text-slate-400">
        거래 내역을 불러오는 중...
      </div>
      <div v-else-if="transactions.length === 0" class="text-center py-12 text-xs text-slate-400">
        조회된 거래 내역이 없습니다.
      </div>
      <table v-else class="w-full text-left border-collapse">
        <thead>
          <tr class="border-b border-slate-100 bg-slate-50/70 text-[11px] font-bold text-slate-400">
            <th class="py-3 px-5">날짜</th>
            <th class="py-3 px-4">종류</th>
            <th class="py-3 px-4">설명 / 상호명</th>
            <th class="py-3 px-4">카테고리</th>
            <th class="py-3 px-4 text-right">금액</th>
            <th class="py-3 px-4 text-center">관리</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 text-xs">
          <tr v-for="tx in transactions" :key="tx.id" class="hover:bg-slate-50/60 transition-colors">
            <td class="py-3.5 px-5 text-slate-500 font-medium">
              {{ formatDate(tx.occurred_at) }}
            </td>
            <td class="py-3.5 px-4">
              <span
                class="px-2 py-0.5 rounded-full text-[10px] font-bold"
                :class="tx.transaction_type === 'EXPENSE' ? 'bg-red-50 text-red-600' : tx.transaction_type === 'INCOME' ? 'bg-emerald-50 text-emerald-600' : 'bg-indigo-50 text-indigo-600'"
              >
                {{ tx.transaction_type === 'EXPENSE' ? '지출' : tx.transaction_type === 'INCOME' ? '수입' : '계정 이동' }}
              </span>
            </td>
            <td class="py-3.5 px-4 font-bold text-slate-800">
              <div>{{ tx.description || '거래' }}</div>
              <div v-if="tx.merchant_name" class="text-[10px] text-slate-400 font-normal">{{ tx.merchant_name }}</div>
            </td>
            <td class="py-3.5 px-4 text-slate-600">
              {{ tx.category_name || '-' }}
            </td>
            <td class="py-3.5 px-4 text-right font-black tracking-tight" :class="tx.transaction_type === 'EXPENSE' ? 'text-slate-900' : tx.transaction_type === 'INCOME' ? 'text-emerald-600' : 'text-indigo-600'">
              {{ tx.transaction_type === 'EXPENSE' ? '− ' : tx.transaction_type === 'INCOME' ? '+ ' : '' }}
              {{ finance.formatCurrency(tx.amount_minor, tx.currency_code) }}
            </td>
            <td class="py-3.5 px-4 text-center">
              <button @click="deleteTx(tx.id)" class="text-[11px] text-slate-400 hover:text-red-600 font-semibold transition-colors">
                삭제
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { RouterLink } from 'vue-router';
import { apiFetch } from '../api/client';
import { useFinanceStore } from '../stores/finance';

const finance = useFinanceStore();

const transactions = ref<any[]>([]);
const loading = ref(false);
const searchQuery = ref('');
const filterType = ref('');

async function fetchTransactions() {
  loading.value = true;
  try {
    let url = '/transactions?';
    if (filterType.value) url += `type=${filterType.value}&`;
    if (searchQuery.value) url += `q=${encodeURIComponent(searchQuery.value)}&`;
    transactions.value = await apiFetch(url);
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  fetchTransactions();
});

async function deleteTx(id: string) {
  if (confirm('이 거래를 삭제하시겠습니까? 관련 잔액과 예산이 자동 재계산됩니다.')) {
    await apiFetch(`/transactions/${id}`, { method: 'DELETE' });
    await fetchTransactions();
    await finance.loadDashboard();
  }
}

async function exportCSV() {
  await apiFetch('/me/data-exports', {
    method: 'POST',
    body: JSON.stringify({ format: 'CSV' }),
  });
  alert('CSV 내보내기 작업이 생성되었습니다. 설정 > 데이터 내보내기에서 다운로드할 수 있습니다.');
}

function formatDate(iso: string) {
  if (!iso) return '';
  const d = new Date(iso);
  return `${d.getFullYear()}.${String(d.getMonth() + 1).padStart(2, '0')}.${String(d.getDate()).padStart(2, '0')}`;
}
</script>
