<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-black text-slate-900 tracking-tight">정기 구독 관리</h1>
        <p class="text-xs text-slate-500 mt-0.5">매달 결제되는 구독 서비스의 갱신일과 예상 금액을 추적합니다.</p>
      </div>
      <button
        @click="showAddModal = true"
        class="px-4 py-2 bg-brand-600 hover:bg-brand-700 text-white rounded-xl text-xs font-bold shadow-md shadow-brand-600/20 transition-all"
      >
        + 구독 추가
      </button>
    </div>

    <!-- Mandatory Disclaimer (Article 17) -->
    <div class="p-3.5 rounded-2xl bg-slate-100 border border-slate-200 text-slate-600 text-xs flex items-start gap-2.5">
      <span>ℹ️</span>
      <p class="text-[11px] leading-relaxed">
        * BN ONE의 구독 관리는 참고용 기록이며, 실제 서비스 해지 또는 자동 결제 승인을 대행하지 않습니다. 예정 구독비는 실제 거래 등록 시 별도 집계됩니다.
      </p>
    </div>

    <!-- Summary KPI -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="p-5 rounded-2xl bg-white border border-slate-200 card-shadow">
        <span class="text-xs font-bold text-slate-400">구독 중인 서비스</span>
        <div class="text-2xl font-black text-slate-900 mt-1">{{ summary?.active_count || 0 }}개</div>
      </div>
      <div class="p-5 rounded-2xl bg-white border border-slate-200 card-shadow">
        <span class="text-xs font-bold text-slate-400">월간 환산 총액</span>
        <div class="text-2xl font-black text-brand-600 mt-1">{{ finance.formatCurrency(summary?.monthly_estimate_minor || 0) }}</div>
      </div>
      <div class="p-5 rounded-2xl bg-white border border-slate-200 card-shadow">
        <span class="text-xs font-bold text-slate-400">연간 환산 예상액</span>
        <div class="text-2xl font-black text-slate-800 mt-1">{{ finance.formatCurrency(summary?.yearly_estimate_minor || 0) }}</div>
      </div>
    </div>

    <!-- Subscriptions List -->
    <div v-if="subscriptions.length === 0" class="bg-white p-12 rounded-3xl border border-slate-200 card-shadow text-center">
      <div class="text-3xl mb-3">🔄</div>
      <div class="text-sm font-bold text-slate-700">등록된 구독 서비스가 없습니다.</div>
      <p class="text-xs text-slate-400 mt-1 mb-4">넷플릭스, 유튜브, 쿠팡 등 정기 결제를 등록해보세요.</p>
      <button @click="showAddModal = true" class="px-4 py-2 bg-brand-600 text-white rounded-xl text-xs font-bold">
        첫 구독 추가하기
      </button>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
        v-for="sub in subscriptions"
        :key="sub.id"
        class="bg-white p-5 rounded-3xl border border-slate-200 card-shadow space-y-3"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2.5">
            <div class="w-10 h-10 rounded-2xl bg-slate-100 flex items-center justify-center font-bold text-sm text-slate-700">
              {{ sub.service_name.charAt(0) }}
            </div>
            <div>
              <h3 class="font-bold text-sm text-slate-900">{{ sub.service_name }}</h3>
              <span class="text-[10px] text-slate-400">{{ sub.billing_cycle }} 결제</span>
            </div>
          </div>
          <span
            class="px-2 py-0.5 rounded-full text-[10px] font-bold"
            :class="sub.status === 'ACTIVE' ? 'bg-emerald-50 text-emerald-600' : 'bg-slate-100 text-slate-500'"
          >
            {{ sub.status === 'ACTIVE' ? '이용중' : '일시정지' }}
          </span>
        </div>

        <div class="flex items-center justify-between pt-2">
          <span class="text-xs text-slate-400">결제 예정일: <strong class="text-slate-700">{{ sub.next_renewal_at }}</strong></span>
          <span class="text-base font-black text-slate-900">{{ finance.formatCurrency(sub.amount_minor, sub.currency_code) }}</span>
        </div>

        <div class="pt-3 border-t border-slate-100 flex items-center justify-between text-[11px]">
          <button
            @click="togglePause(sub)"
            class="text-slate-500 hover:text-slate-800 font-semibold"
          >
            {{ sub.status === 'ACTIVE' ? '일시정지' : '다시 시작' }}
          </button>
          <button
            @click="cancelRecord(sub.id)"
            class="text-red-500 hover:text-red-700 font-semibold"
          >
            기록 해지
          </button>
        </div>
      </div>
    </div>

    <!-- Add Subscription Modal -->
    <div v-if="showAddModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="bg-white max-w-md w-full rounded-3xl p-6 card-shadow space-y-4 border border-slate-200">
        <h3 class="text-lg font-bold text-slate-900">구독 서비스 추가</h3>
        <form @submit.prevent="submitAddSubscription" class="space-y-3">
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">서비스 선택 (또는 직접 입력)</label>
            <select v-model="newSub.service_id" class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm bg-white mb-2">
              <option value="">직접 이름 입력</option>
              <option v-for="s in catalog" :key="s.id" :value="s.id">{{ s.name }}</option>
            </select>
            <input
              v-if="!newSub.service_id"
              v-model="newSub.custom_service_name"
              required
              placeholder="서비스명 입력 (예: 챗GPT)"
              class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">금액 (KRW)</label>
            <input v-model.number="newSub.amount_minor" type="number" min="100" required placeholder="14900" class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm" />
          </div>

          <div class="grid grid-cols-2 gap-2">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">결제 주기</label>
              <select v-model="newSub.billing_cycle" class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-xs bg-white">
                <option value="MONTHLY">매월</option>
                <option value="YEARLY">매년</option>
                <option value="WEEKLY">매주</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">다음 갱신일</label>
              <input v-model="newSub.next_renewal_at" type="date" required class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-xs bg-white" />
            </div>
          </div>

          <div class="flex gap-2 pt-2">
            <button type="button" @click="showAddModal = false" class="w-1/2 py-2.5 bg-slate-100 hover:bg-slate-200 rounded-xl text-xs font-bold text-slate-700">취소</button>
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
const subscriptions = ref<any[]>([]);
const summary = ref<any>(null);
const catalog = ref<any[]>([]);
const showAddModal = ref(false);

const newSub = ref({
  service_id: '',
  custom_service_name: '',
  amount_minor: 14900,
  billing_cycle: 'MONTHLY',
  billing_interval: 1,
  currency_code: 'KRW',
  start_date: new Date().toISOString().slice(0, 10),
  next_renewal_at: new Date(new Date().setDate(new Date().getDate() + 30)).toISOString().slice(0, 10),
  reminder_enabled: true,
});

async function loadData() {
  const [subs, sum, cat] = await Promise.all([
    apiFetch('/me/subscriptions'),
    apiFetch('/me/subscriptions/summary'),
    apiFetch('/subscription-services'),
  ]);
  subscriptions.value = subs;
  summary.value = sum;
  catalog.value = cat;
}

onMounted(() => {
  loadData();
});

async function submitAddSubscription() {
  await apiFetch('/me/subscriptions', {
    method: 'POST',
    body: JSON.stringify(newSub.value),
  });
  showAddModal.value = false;
  await loadData();
}

async function togglePause(sub: any) {
  const action = sub.status === 'ACTIVE' ? 'pause' : 'resume';
  await apiFetch(`/me/subscriptions/${sub.id}/${action}`, { method: 'POST' });
  await loadData();
}

async function cancelRecord(id: string) {
  if (confirm('BN ONE 내부 구독 기록을 해지 처리하시겠습니까? (실제 서비스 웹사이트에서 별도 해지하셔야 합니다)')) {
    await apiFetch(`/me/subscriptions/${id}/cancel-record`, { method: 'POST' });
    await loadData();
  }
}
</script>

