<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-black text-slate-900 tracking-tight">여행 경비 & 예산</h1>
        <p class="text-xs text-slate-500 mt-0.5">국내외 여행별 예산을 수립하고 현지 통화 지출을 집계합니다.</p>
      </div>
      <button
        @click="showCreateTripModal = true"
        class="px-4 py-2 bg-brand-600 hover:bg-brand-700 text-white rounded-xl text-xs font-bold shadow-md shadow-brand-600/20 transition-all flex items-center gap-1.5"
      >
        <span>+</span>
        <span>새 여행 만들기</span>
      </button>
    </div>

    <!-- Travel List -->
    <div v-if="trips.length === 0" class="bg-white p-12 rounded-3xl border border-slate-200 card-shadow text-center">
      <div class="text-3xl mb-3">✈️</div>
      <div class="text-sm font-bold text-slate-700">등록된 여행이 없습니다.</div>
      <p class="text-xs text-slate-400 mt-1 mb-4">도쿄, 파리, 뉴욕 등 여행 일정을 만들고 경비를 관리해보세요.</p>
      <button @click="showCreateTripModal = true" class="px-4 py-2 bg-brand-600 text-white rounded-xl text-xs font-bold">
        첫 여행 만들기
      </button>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-5">
      <div
        v-for="t in trips"
        :key="t.id"
        class="bg-white p-6 rounded-3xl border border-slate-200 card-shadow space-y-4"
      >
        <div class="flex items-center justify-between">
          <div>
            <h3 class="font-bold text-base text-slate-900 flex items-center gap-1.5">
              <span>✈️</span>
              <span>{{ t.name }}</span>
            </h3>
            <span class="text-[10px] text-slate-400">{{ t.city_name }} ({{ t.country_code }}) • {{ t.start_date }} ~ {{ t.end_date }}</span>
          </div>
          <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-blue-50 text-blue-600">
            {{ t.status }}
          </span>
        </div>

        <div class="p-4 rounded-2xl bg-slate-50 border border-slate-100 flex items-center justify-between text-xs">
          <div>
            <div class="text-slate-400 text-[10px]">총 예산</div>
            <div class="font-bold text-slate-800">{{ finance.formatCurrency(t.total_budget_minor, t.base_currency) }}</div>
          </div>
          <div class="text-right">
            <div class="text-slate-400 text-[10px]">현재 사용액</div>
            <div class="font-black text-slate-900">{{ finance.formatCurrency(t.total_spent_minor, t.base_currency) }}</div>
          </div>
        </div>

        <div class="flex items-center justify-between pt-2 text-[11px]">
          <span class="text-slate-500">
            남은 예산: <strong class="text-brand-600">{{ finance.formatCurrency(t.total_budget_minor - t.total_spent_minor, t.base_currency) }}</strong>
          </span>
          <button @click="viewReport(t.id)" class="text-brand-600 hover:underline font-bold">
            결산 리포트 보기 →
          </button>
        </div>
      </div>
    </div>

    <!-- Create Trip Modal -->
    <div v-if="showCreateTripModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="bg-white max-w-md w-full rounded-3xl p-6 card-shadow space-y-4 border border-slate-200">
        <h3 class="text-lg font-bold text-slate-900">새 여행 생성</h3>
        <form @submit.prevent="submitCreateTrip" class="space-y-3">
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">여행 이름</label>
            <input v-model="newTrip.name" required placeholder="예: 2026 도쿄 4박 5일" class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:border-brand-600" />
          </div>
          <div class="grid grid-cols-2 gap-2">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">국가 코드</label>
              <input v-model="newTrip.country_code" required placeholder="JP" class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm uppercase" />
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">도시</label>
              <input v-model="newTrip.city_name" required placeholder="Tokyo" class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm" />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-2">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">시작일</label>
              <input v-model="newTrip.start_date" type="date" required class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-xs bg-white" />
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">종료일</label>
              <input v-model="newTrip.end_date" type="date" required class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-xs bg-white" />
            </div>
          </div>
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">여행 총 예산</label>
            <input v-model.number="newTrip.budget_amount_minor" type="number" min="0" required placeholder="1500000" class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm" />
          </div>
          <div class="flex gap-2 pt-2">
            <button type="button" @click="showCreateTripModal = false" class="w-1/2 py-2.5 bg-slate-100 hover:bg-slate-200 rounded-xl text-xs font-bold text-slate-700">취소</button>
            <button type="submit" class="w-1/2 py-2.5 bg-brand-600 hover:bg-brand-700 text-white rounded-xl text-xs font-bold shadow-md shadow-brand-600/25">여행 만들기</button>
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
const trips = ref<any[]>([]);
const showCreateTripModal = ref(false);

const newTrip = ref({
  name: '',
  country_code: 'JP',
  city_name: 'Tokyo',
  start_date: new Date().toISOString().slice(0, 10),
  end_date: new Date(new Date().setDate(new Date().getDate() + 5)).toISOString().slice(0, 10),
  base_currency: 'KRW',
  budget_amount_minor: 1500000,
});

async function loadTrips() {
  trips.value = await apiFetch('/trips');
}

onMounted(() => {
  loadTrips();
});

async function submitCreateTrip() {
  await apiFetch('/trips', {
    method: 'POST',
    body: JSON.stringify(newTrip.value),
  });
  showCreateTripModal.value = false;
  newTrip.value.name = '';
  await loadTrips();
}

async function viewReport(tripId: string) {
  const report = await apiFetch(`/trips/${tripId}/report`);
  alert(`[${report.trip.name} 결산 리포트]\n총 예산: ${report.total_budget_minor} KRW\n총 지출: ${report.total_spent_minor} KRW\n남은 금액: ${report.remaining_minor} KRW\n\n* ${report.fx_disclaimer}`);
}
</script>
