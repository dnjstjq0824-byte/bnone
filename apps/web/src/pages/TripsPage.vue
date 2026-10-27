<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-[#191F28]">여행 경비 & 예산</h1>
        <p class="text-sm text-[#8B95A1] mt-1">여행별 예산을 수립하고 현지 통화 지출을 집계해요.</p>
      </div>
      <button
        @click="showCreateTripModal = true"
        class="inline-flex items-center gap-1.5 px-4 py-2.5 bg-[#3182F6] hover:bg-[#1B64DA] text-white rounded-xl text-sm font-bold shadow-sm transition-colors"
      >
        <span>+ 새 여행 만들기</span>
      </button>
    </div>

    <!-- Travel List -->
    <div v-if="trips.length === 0" class="bg-white p-12 rounded-3xl text-center">
      <div class="w-14 h-14 mx-auto mb-4 rounded-2xl bg-[#E5F7F0] flex items-center justify-center text-[#00B06B] text-2xl">
        ✈️
      </div>
      <div class="text-base font-bold text-[#191F28]">등록된 여행이 없습니다.</div>
      <p class="text-xs text-[#8B95A1] mt-1 mb-5">도쿄, 파리, 뉴욕 등 여행 일정을 만들고 경비를 관리해보세요.</p>
      <button
        @click="showCreateTripModal = true"
        class="px-5 py-2.5 bg-[#3182F6] text-white rounded-xl text-xs font-bold"
      >
        첫 여행 만들기
      </button>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-5">
      <div
        v-for="t in trips"
        :key="t.id"
        class="bg-white p-6 rounded-3xl space-y-4"
      >
        <div class="flex items-center justify-between">
          <div>
            <h3 class="font-bold text-base text-[#191F28] flex items-center gap-2">
              <span class="w-8 h-8 rounded-xl bg-[#E5F7F0] text-[#00B06B] flex items-center justify-center text-sm">✈️</span>
              <span>{{ t.name }}</span>
            </h3>
            <span class="text-xs text-[#8B95A1] mt-1 block">{{ t.city_name }} ({{ t.country_code }}) • {{ t.start_date }} ~ {{ t.end_date }}</span>
          </div>
          <span class="px-2.5 py-1 rounded-lg text-xs font-bold bg-[#E8F3FF] text-[#3182F6]">
            {{ t.status === 'PLANNING' ? '계획 중' : t.status === 'ONGOING' ? '여행 중' : '완료' }}
          </span>
        </div>

        <div class="p-4 rounded-2xl bg-[#F9FAFB] flex items-center justify-between text-xs">
          <div>
            <div class="text-[#8B95A1] text-xs">총 예산</div>
            <div class="font-extrabold text-[#191F28] text-base mt-0.5">{{ finance.formatCurrency(t.total_budget_minor, t.base_currency) }}</div>
          </div>
          <div class="text-right">
            <div class="text-[#8B95A1] text-xs">현재 사용액</div>
            <div class="font-extrabold text-[#191F28] text-base mt-0.5">{{ finance.formatCurrency(t.total_spent_minor, t.base_currency) }}</div>
          </div>
        </div>

        <div class="flex items-center justify-between pt-1 text-xs">
          <span class="text-[#4E5968]">
            남은 예산: <strong class="text-[#3182F6] font-extrabold">{{ finance.formatCurrency(t.total_budget_minor - t.total_spent_minor, t.base_currency) }}</strong>
          </span>
          <button @click="viewReport(t.id)" class="text-[#3182F6] hover:underline font-bold">
            결산 리포트 보기 →
          </button>
        </div>
      </div>
    </div>

    <!-- Create Trip Modal (Toss Rounded Dialog Style) -->
    <div v-if="showCreateTripModal" class="fixed inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="bg-white max-w-md w-full rounded-3xl p-6 space-y-5">
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-bold text-[#191F28]">새 여행 일정 등록</h3>
          <button @click="showCreateTripModal = false" class="text-[#8B95A1] hover:text-[#191F28]">✕</button>
        </div>
        <form @submit.prevent="submitCreateTrip" class="space-y-4">
          <div>
            <label class="block text-xs font-bold text-[#4E5968] mb-1.5">여행 이름</label>
            <input
              v-model="newTrip.name"
              required
              placeholder="예: 2026 도쿄 4박 5일"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] focus:outline-none focus:border-[#3182F6]"
            />
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-bold text-[#4E5968] mb-1.5">국가 코드</label>
              <input
                v-model="newTrip.country_code"
                required
                placeholder="JP"
                class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] uppercase focus:outline-none focus:border-[#3182F6]"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-[#4E5968] mb-1.5">도시명</label>
              <input
                v-model="newTrip.city_name"
                required
                placeholder="Tokyo"
                class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] focus:outline-none focus:border-[#3182F6]"
              />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-bold text-[#4E5968] mb-1.5">시작일</label>
              <input
                v-model="newTrip.start_date"
                type="date"
                required
                class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-xs text-[#191F28] focus:outline-none focus:border-[#3182F6]"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-[#4E5968] mb-1.5">종료일</label>
              <input
                v-model="newTrip.end_date"
                type="date"
                required
                class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-xs text-[#191F28] focus:outline-none focus:border-[#3182F6]"
              />
            </div>
          </div>
          <div>
            <label class="block text-xs font-bold text-[#4E5968] mb-1.5">여행 총 예산 (KRW)</label>
            <input
              v-model.number="newTrip.budget_amount_minor"
              type="number"
              min="0"
              required
              placeholder="1500000"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] focus:outline-none focus:border-[#3182F6]"
            />
          </div>
          <div class="flex gap-3 pt-2">
            <button
              type="button"
              @click="showCreateTripModal = false"
              class="w-1/2 py-3.5 bg-[#F2F4F6] hover:bg-[#E5E8EB] rounded-2xl text-sm font-bold text-[#4E5968]"
            >
              취소
            </button>
            <button
              type="submit"
              class="w-1/2 py-3.5 bg-[#3182F6] hover:bg-[#1B64DA] text-white rounded-2xl text-sm font-bold shadow-sm"
            >
              여행 만들기
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
