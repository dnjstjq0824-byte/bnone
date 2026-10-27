<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-[#191F28]">소비 & 거래 내역</h1>
        <p class="text-sm text-[#8B95A1] mt-1">수입, 지출 및 계정 간 자금 이동 기록을 확인해요.</p>
      </div>
      <div class="flex items-center gap-2">
        <button
          @click="exportCSV"
          class="inline-flex items-center gap-2 px-4 py-2.5 bg-white hover:bg-[#F2F4F6] text-[#4E5968] rounded-xl text-sm font-semibold border border-[#E5E8EB] transition-colors"
        >
          <Download :size="16" />
          <span>CSV 내보내기</span>
        </button>
        <RouterLink
          to="/app/transactions/new"
          class="inline-flex items-center gap-1.5 px-4 py-2.5 bg-[#3182F6] hover:bg-[#1B64DA] text-white rounded-xl text-sm font-bold shadow-sm transition-colors"
        >
          <Plus :size="18" />
          <span>거래 기록하기</span>
        </RouterLink>
      </div>
    </div>

    <!-- Filters & Search (Toss Pill Filter Style) -->
    <div class="bg-white p-4 rounded-2xl flex flex-wrap items-center justify-between gap-4">
      <div class="flex items-center gap-1.5 p-1 bg-[#F2F4F6] rounded-xl">
        <button
          v-for="tab in filterTabs"
          :key="tab.value"
          @click="selectFilter(tab.value)"
          class="px-3.5 py-1.5 rounded-lg text-xs font-bold transition-all"
          :class="filterType === tab.value ? 'bg-white text-[#191F28] shadow-sm' : 'text-[#8B95A1] hover:text-[#4E5968]'"
        >
          {{ tab.label }}
        </button>
      </div>

      <div class="relative flex-1 min-w-[220px] max-w-sm">
        <Search :size="16" class="absolute left-3.5 top-1/2 -translate-y-1/2 text-[#8B95A1]" />
        <input
          v-model="searchQuery"
          @input="fetchTransactions"
          type="text"
          placeholder="상호명, 메모 검색"
          class="w-full pl-9 pr-3.5 py-2 bg-[#F9FAFB] rounded-xl text-xs text-[#191F28] placeholder-[#8B95A1] border border-transparent focus:border-[#3182F6] focus:bg-white focus:outline-none transition-colors"
        />
      </div>
    </div>

    <!-- Mandatory Notice (Article 1.5) -->
    <div class="p-3.5 rounded-2xl bg-[#E8F3FF] flex items-center gap-2.5">
      <span class="text-[#3182F6] text-sm">ℹ️</span>
      <p class="text-xs text-[#1B64DA] font-semibold">
        BN ONE의 거래 기록은 사용자가 직접 입력한 장부이며 실제 금융기관 거래를 실행하지 않습니다.
      </p>
    </div>

    <!-- Transaction List Card (Toss Timeline Style) -->
    <div class="bg-white rounded-3xl p-6">
      <div v-if="loading" class="text-center py-16 text-sm text-[#8B95A1]">
        거래 내역을 불러오는 중이에요...
      </div>
      <div v-else-if="transactions.length === 0" class="text-center py-16">
        <div class="w-14 h-14 mx-auto mb-4 rounded-2xl bg-[#F2F4F6] flex items-center justify-center text-[#8B95A1]">
          <ReceiptText :size="28" />
        </div>
        <h3 class="text-base font-bold text-[#191F28]">기록된 거래가 없어요</h3>
        <p class="text-xs text-[#8B95A1] mt-1">오늘 쓴 커피 한 잔부터 가볍게 기록해보세요.</p>
        <RouterLink
          to="/app/transactions/new"
          class="inline-block mt-4 px-4 py-2 bg-[#3182F6] text-white rounded-xl text-xs font-bold"
        >
          첫 거래 기록하기
        </RouterLink>
      </div>
      <div v-else class="divide-y divide-[#F2F4F6]">
        <div
          v-for="tx in transactions"
          :key="tx.id"
          class="py-4 first:pt-0 last:pb-0 flex items-center justify-between gap-4 group hover:bg-[#F9FAFB] px-2 rounded-2xl transition-colors"
        >
          <div class="flex items-center gap-3.5 min-w-0">
            <div
              class="w-11 h-11 rounded-2xl flex items-center justify-center flex-shrink-0"
              :class="tx.transaction_type === 'INCOME' ? 'bg-[#E8F3FF] text-[#3182F6]' : tx.transaction_type === 'TRANSFER' ? 'bg-[#F1EDFF] text-[#8061D9]' : 'bg-[#F2F4F6] text-[#4E5968]'"
            >
              <ArrowDownLeft v-if="tx.transaction_type === 'INCOME'" :size="20" />
              <ArrowLeftRight v-else-if="tx.transaction_type === 'TRANSFER'" :size="20" />
              <ArrowUpRight v-else :size="20" />
            </div>
            <div class="min-w-0">
              <div class="flex items-center gap-2">
                <p class="text-sm font-bold text-[#191F28] truncate">{{ tx.description || tx.merchant_name || '거래 기록' }}</p>
                <span
                  class="px-2 py-0.5 rounded-md text-[10px] font-bold"
                  :class="tx.transaction_type === 'INCOME' ? 'bg-[#E8F3FF] text-[#3182F6]' : tx.transaction_type === 'TRANSFER' ? 'bg-[#F1EDFF] text-[#8061D9]' : 'bg-[#F2F4F6] text-[#6B7684]'"
                >
                  {{ tx.transaction_type === 'INCOME' ? '수입' : tx.transaction_type === 'TRANSFER' ? '이동' : '지출' }}
                </span>
              </div>
              <p class="text-xs text-[#8B95A1] mt-0.5 truncate">
                {{ formatDate(tx.occurred_at) }}
                <span v-if="tx.category_name"> · {{ tx.category_name }}</span>
                <span v-if="tx.merchant_name && tx.description !== tx.merchant_name"> · {{ tx.merchant_name }}</span>
              </p>
            </div>
          </div>

          <div class="flex items-center gap-3 flex-shrink-0">
            <span
              class="text-base font-extrabold tracking-tight tabular-nums"
              :class="tx.transaction_type === 'INCOME' ? 'text-[#3182F6]' : 'text-[#191F28]'"
            >
              {{ tx.transaction_type === 'EXPENSE' ? '− ' : tx.transaction_type === 'INCOME' ? '+ ' : '' }}
              {{ finance.formatCurrency(tx.amount_minor, tx.currency_code) }}
            </span>
            <button
              @click="deleteTx(tx.id)"
              class="opacity-0 group-hover:opacity-100 p-1.5 text-[#8B95A1] hover:text-[#F04452] rounded-lg transition-all"
              title="거래 삭제"
            >
              <Trash2 :size="16" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { RouterLink } from 'vue-router';
import { Download, Plus, Search, ReceiptText, ArrowDownLeft, ArrowUpRight, ArrowLeftRight, Trash2 } from 'lucide-vue-next';
import { apiFetch } from '../api/client';
import { useFinanceStore } from '../stores/finance';

const finance = useFinanceStore();

const transactions = ref<any[]>([]);
const loading = ref(false);
const searchQuery = ref('');
const filterType = ref('');

const filterTabs = [
  { label: '전체', value: '' },
  { label: '지출', value: 'EXPENSE' },
  { label: '수입', value: 'INCOME' },
  { label: '계정 이동', value: 'TRANSFER' },
];

function selectFilter(val: string) {
  filterType.value = val;
  fetchTransactions();
}

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
  if (confirm('이 거래를 삭제하시겠습니까? 관련 계정 잔액과 예산이 자동으로 다시 계산됩니다.')) {
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
  return `${d.getMonth() + 1}월 ${d.getDate()}일`;
}
</script>
