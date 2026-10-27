<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-wrap items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-[#191F28]">자산 & 부채 관리</h1>
        <p class="text-sm text-[#8B95A1] mt-1">금융 계정 및 보유 자산과 부채를 한눈에 관리해요.</p>
      </div>
      <div class="flex flex-wrap items-center gap-2">
        <button
          @click="showAddAccountModal = true"
          class="inline-flex items-center gap-1.5 px-4 py-2.5 bg-[#3182F6] hover:bg-[#1B64DA] text-white rounded-xl text-sm font-bold shadow-sm transition-colors"
        >
          <span>+ 계정 추가</span>
        </button>
        <button
          @click="showAddAssetModal = true"
          class="inline-flex items-center gap-1.5 px-4 py-2.5 bg-[#F2F4F6] hover:bg-[#E5E8EB] text-[#4E5968] rounded-xl text-sm font-semibold transition-colors"
        >
          <span>+ 기타 자산</span>
        </button>
        <button
          @click="showAddLiabilityModal = true"
          class="inline-flex items-center gap-1.5 px-4 py-2.5 bg-[#FEECEB] hover:bg-[#FDD8D6] text-[#F04452] rounded-xl text-sm font-semibold transition-colors"
        >
          <span>+ 부채 등록</span>
        </button>
      </div>
    </div>

    <!-- Mandatory Disclaimer (Article 1.5) -->
    <div class="p-4 rounded-2xl bg-[#E8F3FF] flex items-start gap-3">
      <span class="text-[#3182F6] text-base leading-none">ℹ️</span>
      <p class="text-xs text-[#1B64DA] font-semibold leading-relaxed">
        <strong>안내:</strong> 표시된 자산과 잔액은 사용자가 직접 입력한 장부 기록을 기준으로 합니다. 실제 금융기관 연동 조회가 아니며 잔액 변경은 거래 기록 또는 잔액 수정을 통해 반영됩니다.
      </p>
    </div>

    <div v-if="loadError" role="alert" class="p-4 rounded-2xl bg-[#FEECEB] text-[#F04452] text-sm">
      {{ loadError }}
      <button class="underline ml-2 font-bold" @click="loadData">다시 시도</button>
    </div>

    <p v-if="loading" role="status" class="text-sm text-[#8B95A1]">자산 정보를 불러오는 중이에요...</p>

    <AssetEntryDialog v-if="showAddAssetModal" kind="asset" @close="showAddAssetModal = false" @saved="entrySaved" />
    <AssetEntryDialog v-if="showAddLiabilityModal" kind="liability" @close="showAddLiabilityModal = false" @saved="entrySaved" />

    <!-- Accounts Section (Toss Account Cards) -->
    <div class="bg-white p-6 rounded-3xl space-y-4">
      <div class="flex items-center justify-between">
        <h2 class="text-base font-bold text-[#191F28]">금융 계정 장부 ({{ accounts.length }})</h2>
        <button @click="showAddAccountModal = true" class="text-xs font-bold text-[#3182F6] hover:underline">
          + 새 계정 등록
        </button>
      </div>

      <div v-if="accounts.length === 0" class="text-center py-12 text-sm text-[#8B95A1]">
        등록된 계정이 없습니다. '+ 계정 추가'를 눌러 시작해보세요.
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="acc in accounts"
          :key="acc.id"
          class="p-5 rounded-2xl bg-[#F9FAFB] hover:bg-white hover:shadow-sm border border-[#F2F4F6] hover:border-[#3182F6] transition-all flex flex-col justify-between"
        >
          <div>
            <div class="flex items-center justify-between mb-3">
              <span class="text-xs font-bold text-[#191F28] flex items-center gap-2">
                <span class="w-2.5 h-2.5 rounded-full" :style="{ backgroundColor: acc.color || '#3182F6' }"></span>
                {{ acc.name }}
              </span>
              <span class="text-[10px] font-bold px-2 py-0.5 rounded-md bg-white text-[#8B95A1] border border-[#E5E8EB]">
                {{ acc.account_type === 'MANUAL_BANK' ? '수동 계좌' : acc.account_type === 'CASH' ? '현금 지갑' : acc.account_type === 'SAVINGS' ? '저축/적금' : '전자지갑' }}
              </span>
            </div>
            <div class="text-2xl font-extrabold text-[#191F28] tracking-tight tabular-nums">
              {{ finance.formatCurrency(acc.calculated_balance_minor, acc.currency_code) }}
            </div>
          </div>
          <div class="mt-4 pt-3 border-t border-[#E5E8EB] flex items-center justify-between text-xs text-[#8B95A1]">
            <span>시작: {{ finance.formatCurrency(acc.opening_balance_minor, acc.currency_code) }}</span>
            <button class="text-[#8B95A1] hover:text-[#F04452] font-semibold transition-colors" @click="archiveAccount(acc.id)">
              보관하기
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Manual Assets & Liabilities (Toss Side-by-Side Cards) -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Manual Assets -->
      <div class="bg-white p-6 rounded-3xl space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="text-base font-bold text-[#191F28]">기타 보유자산 (부동산, 차량, 투자)</h2>
          <button @click="showAddAssetModal = true" class="text-xs font-bold text-[#3182F6] hover:underline">+ 추가</button>
        </div>
        <div v-if="assets.length === 0" class="text-center py-8 text-xs text-[#8B95A1]">
          등록된 기타 자산이 없습니다.
        </div>
        <div v-else class="space-y-3">
          <div
            v-for="ast in assets"
            :key="ast.id"
            class="p-4 rounded-2xl bg-[#F9FAFB] flex items-center justify-between"
          >
            <div>
              <div class="text-sm font-bold text-[#191F28]">{{ ast.name }}</div>
              <div class="text-xs text-[#8B95A1] mt-0.5">{{ ast.asset_type }} • {{ ast.valuation_date }}</div>
            </div>
            <div class="text-base font-extrabold text-[#191F28] tabular-nums">
              {{ finance.formatCurrency(ast.estimated_value_minor, ast.currency_code) }}
            </div>
          </div>
        </div>
      </div>

      <!-- Liabilities -->
      <div class="bg-white p-6 rounded-3xl space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="text-base font-bold text-[#191F28]">부채 내역 (대출, 할부금)</h2>
          <button @click="showAddLiabilityModal = true" class="text-xs font-bold text-[#F04452] hover:underline">+ 등록</button>
        </div>
        <div v-if="liabilities.length === 0" class="text-center py-8 text-xs text-[#8B95A1]">
          등록된 부채 내역이 없습니다.
        </div>
        <div v-else class="space-y-3">
          <div
            v-for="liab in liabilities"
            :key="liab.id"
            class="p-4 rounded-2xl bg-[#F9FAFB] flex items-center justify-between"
          >
            <div>
              <div class="text-sm font-bold text-[#191F28]">{{ liab.name }}</div>
              <div class="text-xs text-[#8B95A1] mt-0.5">{{ liab.liability_type }} • {{ liab.valuation_date }}</div>
            </div>
            <div class="text-base font-extrabold text-[#F04452] tabular-nums">
              {{ finance.formatCurrency(liab.outstanding_amount_minor, liab.currency_code) }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Account Modal (Toss Rounded Dialog Style) -->
    <div v-if="showAddAccountModal" class="fixed inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="bg-white max-w-md w-full rounded-3xl p-6 space-y-5">
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-bold text-[#191F28]">새 금융 계정 등록</h3>
          <button @click="showAddAccountModal = false" class="text-[#8B95A1] hover:text-[#191F28]">✕</button>
        </div>
        <form @submit.prevent="submitAddAccount" class="space-y-4">
          <div>
            <label class="block text-xs font-bold text-[#4E5968] mb-1.5">계정 이름</label>
            <input
              v-model="newAcc.name"
              required
              placeholder="예: 급여통장, 비상금 현금"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] focus:outline-none focus:border-[#3182F6]"
            />
          </div>
          <div>
            <label class="block text-xs font-bold text-[#4E5968] mb-1.5">계정 유형</label>
            <select
              v-model="newAcc.account_type"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] focus:outline-none focus:border-[#3182F6]"
            >
              <option value="MANUAL_BANK">은행계좌형 수동 장부</option>
              <option value="CASH">현금 지갑</option>
              <option value="SAVINGS">저축 / 적금</option>
              <option value="E_WALLET_MANUAL">전자지갑</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-bold text-[#4E5968] mb-1.5">통화</label>
            <select
              v-model="newAcc.currency_code"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] focus:outline-none focus:border-[#3182F6]"
            >
              <option value="KRW">KRW (대한민국 원)</option>
              <option value="USD">USD (미국 달러)</option>
              <option value="EUR">EUR (유로)</option>
              <option value="JPY">JPY (일본 엔)</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-bold text-[#4E5968] mb-1.5">시작 잔액</label>
            <input
              v-model.number="newAcc.opening_balance_minor"
              type="number"
              min="0"
              placeholder="0"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] focus:outline-none focus:border-[#3182F6]"
            />
          </div>
          <div class="flex gap-3 pt-2">
            <button
              type="button"
              @click="showAddAccountModal = false"
              class="w-1/2 py-3.5 bg-[#F2F4F6] hover:bg-[#E5E8EB] rounded-2xl text-sm font-bold text-[#4E5968]"
            >
              취소
            </button>
            <button
              type="submit"
              class="w-1/2 py-3.5 bg-[#3182F6] hover:bg-[#1B64DA] text-white rounded-2xl text-sm font-bold shadow-sm"
            >
              등록하기
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
import AssetEntryDialog from '../components/AssetEntryDialog.vue';

const finance = useFinanceStore();

const accounts = ref<any[]>([]);
const assets = ref<any[]>([]);
const liabilities = ref<any[]>([]);
const loading = ref(false);
const loadError = ref('');

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
  loading.value = true;
  loadError.value = '';
  try {
    const [accs, asts, liabs] = await Promise.all([
      apiFetch('/accounts'),
      apiFetch('/assets'),
      apiFetch('/liabilities'),
    ]);
    accounts.value = accs;
    assets.value = asts;
    liabilities.value = liabs;
  } catch {
    loadError.value = '자산 정보를 불러오지 못했어요. 표시된 정보가 최신이 아닐 수 있어요.';
  } finally {
    loading.value = false;
  }
}

async function entrySaved() {
  showAddAssetModal.value = false;
  showAddLiabilityModal.value = false;
  await loadData();
  await finance.loadDashboard();
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
