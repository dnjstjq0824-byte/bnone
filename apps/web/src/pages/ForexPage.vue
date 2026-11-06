<template>
  <div class="space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold tracking-tight text-[#141F19]">공인 참고 환율</h1>
      <p class="text-sm text-[#8A948C] mt-1">글로벌 중앙은행 및 공인 기관의 참고용 환율 데이터를 제공해요.</p>
    </div>

    <!-- Mandatory Disclaimer (Article 1.5 & Article 18) -->
    <div class="p-4 rounded-2xl bg-[#FFF3E2] flex items-start gap-3">
      <span class="text-[#FF8D00] text-base leading-none">⚠️</span>
      <p class="text-xs text-[#D88924] font-semibold leading-relaxed">
        <strong>필수 안내:</strong> 표시 환율은 참고용이며 실제 환전·결제 환율과 다를 수 있습니다. MOA은 사용자 간 환전 거래 및 주문을 대행하지 않습니다.
      </p>
    </div>

    <!-- Currency Converter Calculator (Toss Style) -->
    <div class="bg-white p-6 rounded-3xl space-y-4">
      <h2 class="text-base font-bold text-[#141F19]">환율 계산기</h2>

      <div class="grid grid-cols-1 md:grid-cols-5 gap-3 items-end">
        <div class="md:col-span-2">
          <label class="block text-xs font-bold text-[#41524A] mb-1.5">기준 금액</label>
          <div class="flex gap-2">
            <input
              v-model.number="convertReq.source_amount_minor"
              type="number"
              min="1"
              placeholder="100000"
              class="w-full px-4 py-3 rounded-2xl bg-[#F6F8F4] border border-[#E2E7DE] text-sm font-bold text-[#141F19] focus:outline-none focus:border-[#1E5A38]"
            />
            <select
              v-model="convertReq.source_currency"
              class="px-4 py-3 rounded-2xl bg-[#F6F8F4] border border-[#E2E7DE] text-xs font-bold text-[#141F19] focus:outline-none focus:border-[#1E5A38]"
            >
              <option v-for="c in currencies" :key="c.code" :value="c.code">{{ c.code }}</option>
            </select>
          </div>
        </div>

        <div class="text-center md:pb-3.5 text-[#8A948C] font-bold text-lg">
          ⇄
        </div>

        <div class="md:col-span-2">
          <label class="block text-xs font-bold text-[#41524A] mb-1.5">대상 통화</label>
          <div class="flex gap-2">
            <div class="w-full px-4 py-3 rounded-2xl bg-[#F6F8F4] border border-[#E2E7DE] text-sm font-extrabold text-[#1E5A38] flex items-center tabular-nums">
              {{ convertResult ? convertResult.converted_amount_minor?.toLocaleString() : '-' }}
            </div>
            <select
              v-model="convertReq.target_currency"
              class="px-4 py-3 rounded-2xl bg-[#F6F8F4] border border-[#E2E7DE] text-xs font-bold text-[#141F19] focus:outline-none focus:border-[#1E5A38]"
            >
              <option v-for="c in currencies" :key="c.code" :value="c.code">{{ c.code }}</option>
            </select>
          </div>
        </div>
      </div>

      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pt-2">
        <div class="text-xs text-[#8A948C]">
          <span v-if="convertResult?.rate">적용 참고 환율: 1 {{ convertResult.source_currency }} = {{ convertResult.rate }} {{ convertResult.target_currency }}</span>
          <span v-else-if="convertResult?.status === 'NO_VERIFIED_RATE'" class="text-[#D64545] font-bold">검증된 환율 정보가 없습니다.</span>
        </div>
        <button
          @click="calculateConversion"
          class="px-6 py-3 bg-[#1E5A38] hover:bg-[#17452C] text-white rounded-xl text-xs font-bold shadow-sm transition-colors"
        >
          환산 계산하기
        </button>
      </div>
    </div>

    <!-- Reference Rates Table (Toss Clean List) -->
    <div class="bg-white rounded-3xl p-6 space-y-4">
      <div class="flex items-center justify-between">
        <h2 class="text-base font-bold text-[#141F19]">주요 통화 참고 환율 (기준: USD)</h2>
        <span class="text-xs text-[#8A948C]">예시 환율 (기준일: 2026-10-24 고정)</span>
      </div>

      <div class="divide-y divide-[#ECEFE8]">
        <div
          v-for="r in rateList"
          :key="r.id || r.target_currency || r.quote_currency"
          class="py-3.5 flex items-center justify-between"
        >
          <div class="flex items-center gap-3">
            <span class="w-8 h-8 rounded-full bg-[#EBF2EC] text-[#1E5A38] flex items-center justify-center font-bold text-xs">
              {{ r.target_currency || r.quote_currency }}
            </span>
            <div>
              <span class="text-sm font-bold text-[#141F19]">1 {{ r.base_currency || 'USD' }} / {{ r.target_currency || r.quote_currency }}</span>
              <p class="text-xs text-[#8A948C]">{{ r.provider_name || '예시 참고기관' }}</p>
            </div>
          </div>
          <div class="text-right">
            <div class="text-base font-extrabold text-[#141F19] tabular-nums">{{ r.rate }}</div>
            <span class="text-[10px] font-bold text-[#41524A] bg-[#ECEFE8] px-2 py-0.5 rounded-md">정적 예시 환율</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { apiFetch } from '../api/client';
import { useToast } from '../composables/useToast';

const toast = useToast();

const currencies = ref<any[]>([]);
const rateList = ref<any[]>([]);
const convertResult = ref<any>(null);

const convertReq = ref({
  source_currency: 'USD',
  target_currency: 'KRW',
  source_amount_minor: 100,
});

async function loadData() {
  try {
    const [currs, ratesResp] = await Promise.all([
      apiFetch('/currencies'),
      apiFetch('/fx/rates?base=USD'),
    ]);
    currencies.value = Array.isArray(currs) ? currs : [];
    rateList.value = ratesResp?.rates || (Array.isArray(ratesResp) ? ratesResp : []);
  } catch (err: any) {
    console.warn('[ForexPage] Failed to load data:', err);
  }
}

onMounted(() => {
  loadData();
});

async function calculateConversion() {
  try {
    convertResult.value = await apiFetch('/fx/convert', {
      method: 'POST',
      body: JSON.stringify(convertReq.value),
    });
  } catch (err: any) {
    toast.error(err.message || '환산에 실패했습니다.');
  }
}
</script>
