<template>
  <div class="min-h-screen bg-slate-50 flex items-center justify-center p-4">
    <div class="max-w-lg w-full bg-white rounded-3xl p-8 card-shadow border border-slate-100">
      <!-- Progress Bar -->
      <div class="flex items-center justify-between mb-8">
        <div class="flex items-center gap-2">
          <div class="w-8 h-8 rounded-full bg-brand-600 text-white flex items-center justify-center text-xs font-bold">
            {{ step }}
          </div>
          <span class="text-xs font-bold text-slate-700">단계 {{ step }} / 3</span>
        </div>
        <span class="text-xs text-slate-400 font-medium">초기 설정 온보딩</span>
      </div>

      <!-- Step 1: Language & Currency -->
      <div v-if="step === 1" class="space-y-6">
        <div>
          <h3 class="text-xl font-bold text-slate-900 tracking-tight">기준 통화 및 국가 설정</h3>
          <p class="text-xs text-slate-500 mt-1">
            BN ONE에서 전체 자산과 순자산을 집계할 기준 통화를 선택하세요.
          </p>
        </div>

        <div class="space-y-4">
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1.5">거주 국가</label>
            <select v-model="countryCode" class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm bg-white">
              <option value="KR">대한민국 (South Korea)</option>
              <option value="US">미국 (United States)</option>
              <option value="JP">일본 (Japan)</option>
              <option value="SG">싱가포르 (Singapore)</option>
              <option value="GB">영국 (United Kingdom)</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1.5">기준 통화</label>
            <select v-model="baseCurrency" class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm bg-white">
              <option value="KRW">KRW — 대한민국 원 (₩)</option>
              <option value="USD">USD — 미국 달러 ($)</option>
              <option value="EUR">EUR — 유로 (€)</option>
              <option value="JPY">JPY — 일본 엔 (¥)</option>
              <option value="GBP">GBP — 영국 파운드 (£)</option>
            </select>
          </div>
        </div>

        <button
          @click="step = 2"
          class="w-full py-3.5 bg-brand-600 hover:bg-brand-700 text-white rounded-xl font-bold text-sm shadow-md shadow-brand-600/25 transition-all"
        >
          다음: 첫 금융계좌 등록
        </button>
      </div>

      <!-- Step 2: First Financial Account -->
      <div v-else-if="step === 2" class="space-y-6">
        <div>
          <h3 class="text-xl font-bold text-slate-900 tracking-tight">첫 번째 금융계정 등록</h3>
          <p class="text-xs text-slate-500 mt-1">
            자주 사용하는 주거래 수동 계좌 또는 현금 지갑을 등록하세요.
          </p>
        </div>

        <div class="space-y-4">
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1.5">계정 이름</label>
            <input
              v-model="accountName"
              type="text"
              placeholder="예: 주거래 급여계좌, 비상금 지갑"
              class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm focus:outline-none focus:border-brand-600"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1.5">계정 유형</label>
            <select v-model="accountType" class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm bg-white">
              <option value="MANUAL_BANK">은행계좌형 수동 장부</option>
              <option value="CASH">현금 지갑</option>
              <option value="SAVINGS">저축 / 적금</option>
              <option value="E_WALLET_MANUAL">전자지갑 (페이)</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1.5">시작 잔액 ({{ baseCurrency }})</label>
            <input
              v-model.number="openingBalance"
              type="number"
              min="0"
              placeholder="0"
              class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm focus:outline-none focus:border-brand-600"
            />
            <p class="text-[11px] text-slate-400 mt-1">
              * 표시된 자산과 잔액은 사용자가 직접 입력한 기록을 기준으로 계산됩니다.
            </p>
          </div>
        </div>

        <div class="flex gap-3">
          <button @click="step = 1" class="w-1/3 py-3.5 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl font-bold text-sm">
            이전
          </button>
          <button
            @click="submitFirstAccount"
            :disabled="!accountName || isSubmitting"
            class="w-2/3 py-3.5 bg-brand-600 hover:bg-brand-700 disabled:opacity-50 text-white rounded-xl font-bold text-sm shadow-md shadow-brand-600/25 transition-all"
          >
            {{ isSubmitting ? '저장 중...' : '다음: 설정 완료' }}
          </button>
        </div>
      </div>

      <!-- Step 3: Complete -->
      <div v-else class="text-center py-6 space-y-6">
        <div class="w-16 h-16 rounded-full bg-emerald-100 text-emerald-600 mx-auto flex items-center justify-center text-3xl">
          ✓
        </div>

        <div>
          <h3 class="text-2xl font-black text-slate-900 tracking-tight">준비가 완료되었습니다!</h3>
          <p class="text-xs text-slate-500 mt-2 max-w-sm mx-auto leading-relaxed">
            이제 대시보드에서 매일의 수입과 지출, 정기 구독, 예산과 여행 경비를 편리하게 관리하세요.
          </p>
        </div>

        <button
          @click="finishOnboarding"
          class="w-full py-4 bg-brand-600 hover:bg-brand-700 text-white rounded-xl font-bold text-base shadow-lg shadow-brand-600/30 transition-all"
        >
          대시보드로 이동하기
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { apiFetch } from '../api/client';
import { useAuthStore } from '../stores/auth';

const router = useRouter();
const auth = useAuthStore();

const step = ref(1);
const countryCode = ref(auth.user?.country_code || 'KR');
const baseCurrency = ref(auth.user?.base_currency || 'KRW');
const accountName = ref('주거래 계좌');
const accountType = ref('MANUAL_BANK');
const openingBalance = ref<number>(1000000);
const isSubmitting = ref(false);

async function submitFirstAccount() {
  isSubmitting.value = true;
  try {
    // 1. Update user currency/country
    await apiFetch('/me', {
      method: 'PATCH',
      body: JSON.stringify({
        country_code: countryCode.value,
        base_currency: baseCurrency.value,
      }),
    });

    // 2. Create first manual account
    await apiFetch('/accounts', {
      method: 'POST',
      body: JSON.stringify({
        name: accountName.value,
        account_type: accountType.value,
        currency_code: baseCurrency.value,
        opening_balance_minor: openingBalance.value,
        opening_balance_date: new Date().toISOString().split('T')[0],
      }),
    });

    step.value = 3;
  } catch (err) {
    console.error(err);
  } finally {
    isSubmitting.value = false;
  }
}

async function finishOnboarding() {
  await apiFetch('/me', {
    method: 'PATCH',
    body: JSON.stringify({ complete_onboarding: true }),
  });
  await auth.fetchMe();
  router.push('/app');
}
</script>

