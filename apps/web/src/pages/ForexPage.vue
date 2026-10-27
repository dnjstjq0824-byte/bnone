<template>
  <div class="space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold tracking-tight text-[#191F28]">공인 참고 환율</h1>
      <p class="text-sm text-[#8B95A1] mt-1">글로벌 중앙은행 및 공인 기관의 참고용 환율 데이터를 제공해요.</p>
    </div>

    <!-- Mandatory Disclaimer (Article 1.5 & Article 18) -->
    <div class="p-4 rounded-2xl bg-[#FFF3E2] flex items-start gap-3">
      <span class="text-[#FF8D00] text-base leading-none">⚠️</span>
      <p class="text-xs text-[#D88924] font-semibold leading-relaxed">
        <strong>필수 안내:</strong> 표시 환율은 참고용이며 실제 환전·결제 환율과 다를 수 있습니다. BN ONE은 사용자 간 환전 거래 및 주문을 대행하지 않습니다.
      </p>
    </div>

    <!-- Currency Converter Calculator (Toss Style) -->
    <div class="bg-white p-6 rounded-3xl space-y-4">
      <h2 class="text-base font-bold text-[#191F28]">환율 계산기</h2>

      <div class="grid grid-cols-1 md:grid-cols-5 gap-3 items-end">
        <div class="md:col-span-2">
          <label class="block text-xs font-bold text-[#4E5968] mb-1.5">기준 금액</label>
          <div class="flex gap-2">
            <input
              v-model.number="convertReq.source_amount_minor"
              type="number"
              min="1"
              placeholder="100000"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm font-bold text-[#191F28] focus:outline-none focus:border-[#3182F6]"
            />
            <select
              v-model="convertReq.source_currency"
              class="px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-xs font-bold text-[#191F28] focus:outline-none focus:border-[#3182F6]"
            >
              <option v-for="c in currencies" :key="c.code" :value="c.code">{{ c.code }}</option>
            </select>
          </div>
        </div>

        <div class="text-center md:pb-3.5 text-[#8B95A1] font-bold text-lg">
          ⇄
        </div>

        <div class="md:col-span-2">
          <label class="block text-xs font-bold text-[#4E5968] mb-1.5">대상 통화</label>
          <div class="flex gap-2">
            <div class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm font-extrabold text-[#3182F6] flex items-center tabular-nums">
              {{ convertResult ? convertResult.converted_amount_minor?.toLocaleString() : '-' }}
            </div>
            <select
              v-model="convertReq.target_currency"
              class="px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-xs font-bold text-[#191F28] focus:outline-none focus:border-[#3182F6]"
            >
              <option v-for="c in currencies" :key="c.code" :value="c.code">{{ c.code }}</option>
            </select>
          </div>
        </div>
      </div>

      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pt-2">
        <div class="text-xs text-[#8B95A1]">
          <span v-if="convertResult?.rate">적용 참고 환율: 1 {{ convertResult.source_currency }} = {{ convertResult.rate }} {{ convertResult.target_currency }}</span>
          <span v-else-if="convertResult?.status === 'NO_VERIFIED_RATE'" class="text-[#F04452] font-bold">검증된 환율 정보가 없습니다.</span>
        </div>
        <button
          @click="calculateConversion"
          class="px-6 py-3 bg-[#3182F6] hover:bg-[#1B64DA] text-white rounded-xl text-xs font-bold shadow-sm transition-colors"
        >
          환산 계산하기
        </button>
      </div>
    </div>

    <!-- Reference Rates Table (Toss Clean List) -->
    <div class="bg-white rounded-3xl p-6 space-y-4">
      <div class="flex items-center justify-between">
        <h2 class="text-base font-bold text-[#191F28]">주요 통화 참고 환율 (기준: USD)</h2>
        <span class="text-xs text-[#8B95A1]">한국은행 / ECB 공인 참고</span>
      </div>

      <div class="divide-y divide-[#F2F4F6]">
        <div
          v-for="r in rateList"
          :key="r.id || r.target_currency"
          class="py-3.5 flex items-center justify-between"
        >
          <div class="flex items-center gap-3">
            <span class="w-8 h-8 rounded-full bg-[#E8F3FF] text-[#3182F6] flex items-center justify-center font-bold text-xs">
              {{ r.target_currency || r.quote_currency }}
            </span>
            <div>
              <span class="text-sm font-bold text-[#191F28]">1 {{ r.base_currency || 'USD' }} / {{ r.target_currency || r.quote_currency }}</span>
              <p class="text-xs text-[#8B95A1]">{{ r.provider_name || '공인 참고기관' }}</p>
            </div>
          </div>
          <div class="text-right">
            <div class="text-base font-extrabold text-[#191F28] tabular-nums">{{ r.rate }}</div>
            <span class="text-[10px] font-bold text-[#00B06B] bg-[#E5F7F0] px-2 py-0.5 rounded-md">공인 검증 완료</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { apiFetch } from '../api/client';

const currencies = ref<any[]>([]);
const rateList = ref<any[]>([]);
const convertResult = ref<any>(null);

const convertReq = ref({
  source_currency: 'USD',
  target_currency: 'KRW',
  source_amount_minor: 100,
});

async function loadData() {
  const [currs, ratesResp] = await Promise.all([
    apiFetch('/currencies'),
    apiFetch('/fx/rates?base=USD'),
  ]);
  currencies.value = currs;
  rateList.value = ratesResp.rates || [];
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
    alert(err.message || '환산에 실패했습니다.');
  }
}
</script>
