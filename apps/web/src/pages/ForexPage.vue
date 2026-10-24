<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-black text-slate-900 tracking-tight">공인 참고 환율</h1>
      <p class="text-xs text-slate-500 mt-0.5">글로벌 중앙은행 및 공인 기관의 참고용 환율 데이터를 제공합니다.</p>
    </div>

    <!-- Mandatory Disclaimer (Article 1.5 & Article 18) -->
    <div class="p-3.5 rounded-2xl bg-amber-50/70 border border-amber-200/60 flex items-start gap-2.5">
      <span class="text-amber-600 text-sm">⚠️</span>
      <p class="text-xs text-amber-900 leading-relaxed font-medium">
        <strong>필수 안내:</strong> 표시 환율은 참고용이며 실제 환전·결제 환율과 다를 수 있습니다. BN ONE은 사용자 간 환전 거래 및 주문을 대행하지 않습니다.
      </p>
    </div>

    <!-- Currency Converter Calculator -->
    <div class="bg-white p-6 rounded-3xl border border-slate-200 card-shadow space-y-4">
      <h2 class="text-base font-bold text-slate-900">환율 계산기</h2>

      <div class="grid grid-cols-1 md:grid-cols-5 gap-3 items-end">
        <div class="md:col-span-2">
          <label class="block text-xs font-bold text-slate-700 mb-1">기준 금액</label>
          <div class="flex gap-2">
            <input
              v-model.number="convertReq.source_amount_minor"
              type="number"
              min="1"
              placeholder="100000"
              class="w-full px-4 py-2.5 rounded-xl border border-slate-200 text-sm font-bold"
            />
            <select v-model="convertReq.source_currency" class="px-3 py-2.5 rounded-xl border border-slate-200 text-xs bg-white font-bold">
              <option v-for="c in currencies" :key="c.code" :value="c.code">{{ c.code }}</option>
            </select>
          </div>
        </div>

        <div class="text-center md:pb-2.5 text-slate-400 font-bold">
          ⇄
        </div>

        <div class="md:col-span-2">
          <label class="block text-xs font-bold text-slate-700 mb-1">대상 통화</label>
          <div class="flex gap-2">
            <div class="w-full px-4 py-2.5 rounded-xl bg-slate-50 border border-slate-200 text-sm font-black text-slate-900 flex items-center">
              {{ convertResult ? convertResult.converted_amount_minor.toLocaleString() : '-' }}
            </div>
            <select v-model="convertReq.target_currency" class="px-3 py-2.5 rounded-xl border border-slate-200 text-xs bg-white font-bold">
              <option v-for="c in currencies" :key="c.code" :value="c.code">{{ c.code }}</option>
            </select>
          </div>
        </div>
      </div>

      <div class="flex items-center justify-between pt-2">
        <div class="text-[11px] text-slate-400">
          <span v-if="convertResult?.rate">적용 참고 환율: 1 {{ convertResult.source_currency }} = {{ convertResult.rate }} {{ convertResult.target_currency }} (출처: {{ convertResult.provider || '공인 참고' }})</span>
          <span v-else-if="convertResult?.status === 'NO_VERIFIED_RATE'" class="text-red-500 font-bold">검증된 환율 정보가 없습니다.</span>
        </div>
        <button
          @click="calculateConversion"
          class="px-5 py-2.5 bg-brand-600 hover:bg-brand-700 text-white rounded-xl text-xs font-bold shadow-md shadow-brand-600/20 transition-all"
        >
          환산 계산
        </button>
      </div>
    </div>

    <!-- Reference Rates Table -->
    <div class="bg-white rounded-3xl border border-slate-200 card-shadow overflow-hidden">
      <div class="p-5 border-b border-slate-100 flex items-center justify-between">
        <h2 class="text-base font-bold text-slate-900">주요 통화 참고 환율 (기준: KRW)</h2>
        <span class="text-[11px] text-slate-400">출처: 한국은행 / 유럽중앙은행</span>
      </div>
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-slate-50/70 text-[11px] font-bold text-slate-400 border-b border-slate-100">
            <th class="py-3 px-5">통화 쌍</th>
            <th class="py-3 px-4">참고 환율</th>
            <th class="py-3 px-4">데이터 제공기관</th>
            <th class="py-3 px-4 text-center">검증 상태</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 text-xs">
          <tr v-for="r in rateList" :key="r.id" class="hover:bg-slate-50/60">
            <td class="py-3.5 px-5 font-bold text-slate-800">
              1 {{ r.base_currency }} / {{ r.quote_currency }}
            </td>
            <td class="py-3.5 px-4 font-black text-slate-900">
              {{ r.rate }}
            </td>
            <td class="py-3.5 px-4 text-slate-500">
              {{ r.provider_name }}
            </td>
            <td class="py-3.5 px-4 text-center">
              <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-600">
                {{ r.status }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
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

