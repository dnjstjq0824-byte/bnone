<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-[#141F19]">여행 경비 & 예산</h1>
        <p class="text-sm text-[#8A948C] mt-1">여행별 예산을 수립하고 현지 통화 지출을 집계해요.</p>
      </div>
      <button
        @click="showCreateTripModal = true"
        class="inline-flex items-center gap-1.5 px-4 py-2.5 bg-[#1E5A38] hover:bg-[#17452C] text-white rounded-xl text-sm font-bold shadow-sm transition-colors"
      >
        <span>+ 새 여행 만들기</span>
      </button>
    </div>

    <!-- Error alert -->
    <div v-if="loadError" class="p-4 rounded-2xl bg-[#FBEAEA] text-[#D64545] text-xs font-medium">
      {{ loadError }}
      <button class="underline ml-2 font-bold" @click="loadTrips">다시 시도</button>
    </div>

    <!-- Travel List -->
    <div v-if="trips.length === 0" class="bg-white p-12 rounded-3xl text-center">
      <div class="w-14 h-14 mx-auto mb-4 rounded-2xl bg-[#E4EFE7] flex items-center justify-center text-[#2A7248] text-2xl">
        ✈️
      </div>
      <div class="text-base font-bold text-[#141F19]">등록된 여행이 없습니다.</div>
      <p class="text-xs text-[#8A948C] mt-1 mb-5">도쿄, 파리, 뉴욕 등 여행 일정을 만들고 경비를 관리해보세요.</p>
      <button
        @click="showCreateTripModal = true"
        class="px-5 py-2.5 bg-[#1E5A38] text-white rounded-xl text-xs font-bold"
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
            <h3 class="font-bold text-base text-[#141F19] flex items-center gap-2">
              <span class="w-8 h-8 rounded-xl bg-[#E4EFE7] text-[#2A7248] flex items-center justify-center text-sm">✈️</span>
              <span>{{ t.name }}</span>
            </h3>
            <span class="text-xs text-[#8A948C] mt-1 block">
              {{ t.city_name || (t.destination ? t.destination.split(',')[0] : '여행지') }}
              ({{ t.country_code || (t.destination ? t.destination.split(',')[1]?.trim() : 'KR') }}) • {{ t.start_date }} ~ {{ t.end_date }}
            </span>
          </div>
          <span class="px-2.5 py-1 rounded-lg text-xs font-bold bg-[#EBF2EC] text-[#1E5A38]">
            {{ t.status === 'PLANNING' ? '계획 중' : t.status === 'ONGOING' ? '여행 중' : '완료' }}
          </span>
        </div>

        <div class="p-4 rounded-2xl bg-[#F6F8F4] flex items-center justify-between text-xs">
          <div>
            <div class="text-[#8A948C] text-xs">총 예산</div>
            <div class="font-extrabold text-[#141F19] text-base mt-0.5">
              {{ finance.formatCurrency(t.total_budget_minor ?? t.budget_minor ?? 0, t.base_currency || t.currency_code || 'KRW') }}
            </div>
          </div>
          <div class="text-right">
            <div class="text-[#8A948C] text-xs">현재 사용액</div>
            <div class="font-extrabold text-[#141F19] text-base mt-0.5">
              {{ finance.formatCurrency(t.total_spent_minor ?? t.spent_minor ?? 0, t.base_currency || t.currency_code || 'KRW') }}
            </div>
          </div>
        </div>

        <div class="flex items-center justify-between pt-1 text-xs">
          <span class="text-[#41524A]">
            남은 예산: <strong class="text-[#1E5A38] font-extrabold">{{
              finance.formatCurrency(
                (t.total_budget_minor ?? t.budget_minor ?? 0) - (t.total_spent_minor ?? t.spent_minor ?? 0),
                t.base_currency || t.currency_code || 'KRW'
              )
            }}</strong>
          </span>
          <button @click="viewReport(t.id)" class="text-[#1E5A38] hover:underline font-bold">
            결산 리포트 보기 →
          </button>
        </div>
      </div>
    </div>

    <!-- Create Trip Modal (Toss Rounded Dialog Style) -->
    <div v-if="showCreateTripModal" class="fixed inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="bg-white max-w-md w-full rounded-3xl p-6 space-y-5">
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-bold text-[#141F19]">새 여행 일정 등록</h3>
          <button @click="showCreateTripModal = false" class="text-[#8A948C] hover:text-[#141F19]">✕</button>
        </div>
        <form @submit.prevent="submitCreateTrip" class="space-y-4">
          <div>
            <label class="block text-xs font-bold text-[#41524A] mb-1.5">여행 이름</label>
            <input
              v-model="newTrip.name"
              required
              placeholder="예: 2026 도쿄 4박 5일"
              class="w-full px-4 py-3 rounded-2xl bg-[#F6F8F4] border border-[#E2E7DE] text-sm text-[#141F19] focus:outline-none focus:border-[#1E5A38]"
            />
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-bold text-[#41524A] mb-1.5">국가 코드</label>
              <input
                v-model="newTrip.country_code"
                required
                placeholder="JP"
                class="w-full px-4 py-3 rounded-2xl bg-[#F6F8F4] border border-[#E2E7DE] text-sm text-[#141F19] uppercase focus:outline-none focus:border-[#1E5A38]"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-[#41524A] mb-1.5">도시명</label>
              <input
                v-model="newTrip.city_name"
                required
                placeholder="Tokyo"
                class="w-full px-4 py-3 rounded-2xl bg-[#F6F8F4] border border-[#E2E7DE] text-sm text-[#141F19] focus:outline-none focus:border-[#1E5A38]"
              />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-bold text-[#41524A] mb-1.5">시작일</label>
              <input
                v-model="newTrip.start_date"
                type="date"
                required
                class="w-full px-4 py-3 rounded-2xl bg-[#F6F8F4] border border-[#E2E7DE] text-xs text-[#141F19] focus:outline-none focus:border-[#1E5A38]"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-[#41524A] mb-1.5">종료일</label>
              <input
                v-model="newTrip.end_date"
                type="date"
                required
                class="w-full px-4 py-3 rounded-2xl bg-[#F6F8F4] border border-[#E2E7DE] text-xs text-[#141F19] focus:outline-none focus:border-[#1E5A38]"
              />
            </div>
          </div>
          <div>
            <label class="block text-xs font-bold text-[#41524A] mb-1.5">여행 총 예산 (KRW)</label>
            <input
              v-model.number="newTrip.budget_amount_minor"
              type="number"
              min="0"
              required
              placeholder="1500000"
              class="w-full px-4 py-3 rounded-2xl bg-[#F6F8F4] border border-[#E2E7DE] text-sm text-[#141F19] focus:outline-none focus:border-[#1E5A38]"
            />
          </div>
          <div class="flex gap-3 pt-2">
            <button
              type="button"
              @click="showCreateTripModal = false"
              class="w-1/2 py-3.5 bg-[#ECEFE8] hover:bg-[#E2E7DE] rounded-2xl text-sm font-bold text-[#41524A]"
            >
              취소
            </button>
            <button
              type="submit"
              class="w-1/2 py-3.5 bg-[#1E5A38] hover:bg-[#17452C] text-white rounded-2xl text-sm font-bold shadow-sm"
            >
              여행 만들기
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Trip Report Modal -->
    <div v-if="reportDetail" class="fixed inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="bg-white max-w-md w-full rounded-3xl p-6 space-y-5">
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-bold text-[#141F19]">{{ reportDetail.trip?.name }} 결산 리포트</h3>
          <button @click="reportDetail = null" class="text-[#8A948C] hover:text-[#141F19]">✕</button>
        </div>
        <div class="p-4 rounded-2xl bg-[#F6F8F4] space-y-3 text-xs">
          <div class="flex items-center justify-between">
            <span class="text-[#8A948C]">총 예산</span>
            <strong class="text-[#141F19] text-base">{{ finance.formatCurrency(reportDetail.total_budget_minor ?? 0, reportDetail.trip?.base_currency || 'KRW') }}</strong>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-[#8A948C]">총 지출</span>
            <strong class="text-[#141F19] text-base">{{ finance.formatCurrency(reportDetail.total_spent_minor ?? 0, reportDetail.trip?.base_currency || 'KRW') }}</strong>
          </div>
          <div class="flex items-center justify-between pt-3 border-t border-[#E2E7DE]">
            <span class="text-[#8A948C]">남은 금액</span>
            <strong class="text-[#1E5A38] text-base">{{ finance.formatCurrency(reportDetail.remaining_minor ?? 0, reportDetail.trip?.base_currency || 'KRW') }}</strong>
          </div>
        </div>
        <p v-if="reportDetail.fx_disclaimer" class="text-[11px] text-[#8A948C] leading-relaxed">* {{ reportDetail.fx_disclaimer }}</p>
        <button
          @click="reportDetail = null"
          class="w-full py-3.5 bg-[#ECEFE8] hover:bg-[#E2E7DE] rounded-2xl text-sm font-bold text-[#41524A]"
        >
          닫기
        </button>
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
const loadError = ref('');
const reportDetail = ref<any | null>(null);

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
  loadError.value = '';
  try {
    const res = await apiFetch('/trips');
    trips.value = Array.isArray(res) ? res : [];
  } catch (err: any) {
    loadError.value = err.message || '여행 일정을 불러오지 못했습니다.';
  }
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
  reportDetail.value = await apiFetch(`/trips/${tripId}/report`);
}
</script>
