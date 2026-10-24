<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-black text-slate-900 tracking-tight">자산 & 부채 관리</h1>
        <p class="text-xs text-slate-500 mt-0.5">수동 금융계정 및 보유 자산·부채를 체계적으로 관리합니다.</p>
      </div>
      <div class="flex items-center gap-2">
        <button
          @click="showAddAccountModal = true"
          class="px-3.5 py-2 bg-brand-600 hover:bg-brand-700 text-white rounded-xl text-xs font-bold shadow-sm transition-all"
        >
          + 계정 추가
        </button>
        <button
          @click="showAddAssetModal = true"
          class="px-3.5 py-2 bg-slate-800 hover:bg-slate-700 text-white rounded-xl text-xs font-bold transition-all"
        >
          + 기타 자산
        </button>
        <button
          @click="showAddLiabilityModal = true"
          class="px-3.5 py-2 bg-slate-100 hover:bg-slate-200 text-slate-800 rounded-xl text-xs font-bold transition-all"
        >
          + 부채 등록
        </button>
      </div>
    </div>

    <!-- Mandatory Disclaimer (Article 1.5) -->
    <div class="p-3.5 rounded-2xl bg-amber-50/70 border border-amber-200/60 flex items-start gap-3">
      <span class="text-amber-600 text-base leading-none">⚠️</span>
      <p class="text-xs text-amber-900 font-medium">
        <strong>필수 안내:</strong> 표시된 자산과 잔액은 사용자가 직접 입력한 기록을 기준으로 합니다. 금융기관 실제 연동 기능이 아니며 잔액 변경은 거래 또는 잔액 조정으로 반영됩니다.
      </p>
    </div>

    <!-- Accounts Section -->
    <div class="bg-white p-6 rounded-3xl border border-slate-200 card-shadow space-y-4">
      <h2 class="text-base font-bold text-slate-900">금융 계정 장부 ({{ accounts.length }})</h2>
      <div v-if="accounts.length === 0" class="text-center py-8 text-xs text-slate-400">
        등록된 계정이 없습니다. '+ 계정 추가'를 눌러 시작해보세요.
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="acc in accounts"
          :key="acc.id"
          class="p-5 rounded-2xl border border-slate-200/80 bg-slate-50/50 hover:bg-white hover:border-brand-500 transition-all card-shadow flex flex-col justify-between"
        >
          <div>
            <div class="flex items-center justify-between mb-3">
              <span class="text-xs font-bold text-slate-800 flex items-center gap-2">
                <span class="w-3 h-3 rounded-full" :style="{ backgroundColor: acc.color }"></span>
                {{ acc.name }}
              </span>
              <span class="text-[10px] font-bold px-2 py-0.5 rounded bg-white border border-slate-200 text-slate-500">
                {{ acc.account_type }}
              </span>
            </div>
            <div class="text-2xl font-black text-slate-900 tracking-tight">
              {{ finance.formatCurrency(acc.calculated_balance_minor, acc.currency_code) }}
            </div>
          </div>
          <div class="mt-4 pt-3 border-t border-slate-100 flex items-center justify-between text-[11px] text-slate-400">
            <span>시작: {{ finance.formatCurrency(acc.opening_balance_minor, acc.currency_code) }}</span>
            <span class="text-brand-600 font-bold hover:underline cursor-pointer" @click="archiveAccount(acc.id)">
              보관하기
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Manual Assets & Liabilities -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Manual Assets -->
      <div class="bg-white p-6 rounded-3xl border border-slate-200 card-shadow space-y-4">
        <h2 class="text-base font-bold text-slate-900">기타 보유자산 (부동산, 차량, 수동투자)</h2>
        <div v-if="assets.length === 0" class="text-center py-8 text-xs text-slate-400">
          등록된 기타 자산이 없습니다.
        </div>
        <div v-else class="space-y-3">
          <div v-for="ast in assets" :key="ast.id" class="p-4 rounded-2xl border border-slate-100 bg-slate-50 flex items-center justify-between">
            <div>
              <div class="text-xs font-bold text-slate-800">{{ ast.name }}</div>
              <div class="text-[10px] text-slate-400">{{ ast.asset_type }} • {{ ast.valuation_date }}</div>
            </div>
            <div class="text-sm font-black text-slate-900">
              {{ finance.formatCurrency(ast.estimated_value_minor, ast.currency_code) }}
            </div>
          </div>
        </div>
      </div>

      <!-- Liabilities -->
      <div class="bg-white p-6 rounded-3xl border border-slate-200 card-shadow space-y-4">
        <h2 class="text-base font-bold text-slate-900">부채 내역 (대출, 할부금)</h2>
        <div v-if="liabilities.length === 0" class="text-center py-8 text-xs text-slate-400">
          등록된 부채 내역이 없습니다.
        </div>
        <div v-else class="space-y-3">
          <div v-for="liab in liabilities" :key="liab.id" class="p-4 rounded-2xl border border-slate-100 bg-slate-50 flex items-center justify-between">
            <div>
              <div class="text-xs font-bold text-slate-800">{{ liab.name }}</div>
              <div class="text-[10px] text-slate-400">{{ liab.liability_type }} • {{ liab.valuation_date }}</div>
            </div>
            <div class="text-sm font-black text-red-600">
              {{ finance.formatCurrency(liab.outstanding_amount_minor, liab.currency_code) }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Account Modal -->
    <div v-if="showAddAccountModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="bg-white max-w-md w-full rounded-3xl p-6 card-shadow space-y-4 border border-slate-200">
        <h3 class="text-lg font-bold text-slate-900">새 금융계정 추가</h3>
        <form @submit.prevent="submitAddAccount" class="space-y-3">
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">계정 이름</label>
            <input v-model="newAcc.name" required placeholder="예: 비상금 입출금" class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:border-brand-600" />
          </div>
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">계정 유형</label>
            <select v-model="newAcc.account_type" class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm bg-white">
              <option value="MANUAL_BANK">은행계좌형 수동 장부</option>
              <option value="CASH">현금 지갑</option>
              <option value="SAVINGS">저축 / 적금</option>
              <option value="E_WALLET_MANUAL">전자지갑</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">통화</label>
            <select v-model="newAcc.currency_code" class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm bg-white">
              <option value="KRW">KRW</option>
              <option value="USD">USD</option>
              <option value="EUR">EUR</option>
              <option value="JPY">JPY</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">시작 잔액</label>
            <input v-model.number="newAcc.opening_balance_minor" type="number" min="0" placeholder="0" class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm" />
          </div>
          <div class="flex gap-2 pt-2">
            <button type="button" @click="showAddAccountModal = false" class="w-1/2 py-2.5 bg-slate-100 hover:bg-slate-200 rounded-xl text-xs font-bold text-slate-700">취소</button>
            <button type="submit" class="w-1/2 py-2.5 bg-brand-600 hover:bg-brand-700 text-white rounded-xl text-xs font-bold shadow-md shadow-brand-600/25">등록하기</button>
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

const accounts = ref<any[]>([]);
const assets = ref<any[]>([]);
const liabilities = ref<any[]>([]);

const showAddAccountModal = ref(false);
const showAddAssetModal = ref(false);
const showAddLiabilityModal = ref(false);

const newAcc = ref({
  name: '',
  account_type: 'MANUAL_BANK',
  currency_code: 'KRW',
  opening_balance_minor: 0,
});

async function loadData() {
  const [accs, asts, liabs] = await Promise.all([
    apiFetch('/accounts'),
    apiFetch('/assets'),
    apiFetch('/liabilities'),
  ]);
  accounts.value = accs;
  assets.value = asts;
  liabilities.value = liabs;
}

onMounted(() => {
  loadData();
});

async function submitAddAccount() {
  await apiFetch('/accounts', {
    method: 'POST',
    body: JSON.stringify(newAcc.value),
  });
  showAddAccountModal.value = false;
  newAcc.value = { name: '', account_type: 'MANUAL_BANK', currency_code: 'KRW', opening_balance_minor: 0 };
  await loadData();
  await finance.loadDashboard();
}

async function archiveAccount(id: string) {
  if (confirm('이 계정을 보관하시겠습니까?')) {
    await apiFetch(`/accounts/${id}/archive`, { method: 'POST' });
    await loadData();
    await finance.loadDashboard();
  }
}
</script>
