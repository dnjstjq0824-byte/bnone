<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-[#191F28]">정기 구독 관리</h1>
        <p class="text-sm text-[#8B95A1] mt-1">매달 결제되는 구독 서비스의 갱신일과 예상 금액을 추적해요.</p>
      </div>
      <button
        @click="showAddModal = true"
        class="inline-flex items-center gap-1.5 px-4 py-2.5 bg-[#3182F6] hover:bg-[#1B64DA] text-white rounded-xl text-sm font-bold shadow-sm transition-colors"
      >
        <span>+ 구독 추가</span>
      </button>
    </div>

    <!-- Mandatory Disclaimer (Article 17) -->
    <div class="p-4 rounded-2xl bg-[#E8F3FF] flex items-start gap-3">
      <span class="text-[#3182F6] text-base leading-none">ℹ️</span>
      <p class="text-xs text-[#1B64DA] font-semibold leading-relaxed">
        <strong>안내:</strong> BN ONE의 구독 관리는 참고용 기록이며 실제 서비스 자동 결제 승인이나 자동 해지를 대행하지 않습니다.
      </p>
    </div>

    <!-- Summary KPI (Toss Stats Cards) -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="p-6 rounded-3xl bg-white space-y-1">
        <span class="text-xs font-bold text-[#8B95A1]">구독 중인 서비스</span>
        <div class="text-2xl font-extrabold text-[#191F28] tracking-tight tabular-nums">{{ summary?.active_count || 0 }}개</div>
      </div>
      <div class="p-6 rounded-3xl bg-white space-y-1">
        <span class="text-xs font-bold text-[#8B95A1]">매월 고정 구독료</span>
        <div class="text-2xl font-extrabold text-[#3182F6] tracking-tight tabular-nums">{{ finance.formatCurrency(summary?.monthly_estimate_minor || 0) }}</div>
      </div>
      <div class="p-6 rounded-3xl bg-white space-y-1">
        <span class="text-xs font-bold text-[#8B95A1]">연간 환산 예상액</span>
        <div class="text-2xl font-extrabold text-[#191F28] tracking-tight tabular-nums">{{ finance.formatCurrency(summary?.yearly_estimate_minor || 0) }}</div>
      </div>
    </div>

    <!-- Subscriptions List -->
    <div v-if="subscriptions.length === 0" class="bg-white p-12 rounded-3xl text-center">
      <div class="w-14 h-14 mx-auto mb-4 rounded-2xl bg-[#F1EDFF] flex items-center justify-center text-[#8061D9] text-2xl">
        🔄
      </div>
      <div class="text-base font-bold text-[#191F28]">등록된 구독 서비스가 없습니다.</div>
      <p class="text-xs text-[#8B95A1] mt-1 mb-5">넷플릭스, 유튜브, 쿠팡 등 정기 결제를 등록해보세요.</p>
      <button
        @click="showAddModal = true"
        class="px-5 py-2.5 bg-[#3182F6] text-white rounded-xl text-xs font-bold"
      >
        첫 구독 추가하기
      </button>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
        v-for="sub in subscriptions"
        :key="sub.id"
        class="bg-white p-6 rounded-3xl space-y-4"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div class="w-11 h-11 rounded-2xl bg-[#F1EDFF] text-[#8061D9] flex items-center justify-center font-extrabold text-base">
              {{ sub.service_name.charAt(0) }}
            </div>
            <div>
              <h3 class="font-bold text-base text-[#191F28]">{{ sub.service_name }}</h3>
              <span class="text-xs text-[#8B95A1]">{{ sub.billing_cycle === 'MONTHLY' ? '매월' : '매년' }} 결제</span>
            </div>
          </div>
          <span
            class="px-2.5 py-1 rounded-lg text-xs font-bold"
            :class="sub.status === 'ACTIVE' ? 'bg-[#E5F7F0] text-[#00B06B]' : 'bg-[#F2F4F6] text-[#8B95A1]'"
          >
            {{ sub.status === 'ACTIVE' ? '이용중' : '일시정지' }}
          </span>
        </div>

        <div class="flex items-center justify-between pt-1">
          <span class="text-xs text-[#8B95A1]">다음 결제일: <strong class="text-[#191F28] font-bold">{{ sub.next_renewal_at }}</strong></span>
          <span class="text-lg font-extrabold text-[#191F28] tabular-nums">{{ finance.formatCurrency(sub.amount_minor, sub.currency_code) }}</span>
        </div>

        <div class="pt-3 border-t border-[#F2F4F6] flex items-center justify-between text-xs">
          <button
            @click="togglePause(sub)"
            class="text-[#4E5968] hover:text-[#191F28] font-bold transition-colors"
          >
            {{ sub.status === 'ACTIVE' ? '일시정지' : '다시 시작' }}
          </button>
          <button
            @click="cancelRecord(sub.id)"
            class="text-[#8B95A1] hover:text-[#F04452] font-semibold transition-colors"
          >
            기록 해지
          </button>
        </div>
      </div>
    </div>

    <!-- Add Subscription Modal (Toss Rounded Dialog Style) -->
    <div v-if="showAddModal" class="fixed inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="bg-white max-w-md w-full rounded-3xl p-6 space-y-5">
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-bold text-[#191F28]">정기 구독 추가</h3>
          <button @click="showAddModal = false" class="text-[#8B95A1] hover:text-[#191F28]">✕</button>
        </div>
        <form @submit.prevent="submitAddSubscription" class="space-y-4">
          <div>
            <label class="block text-xs font-bold text-[#4E5968] mb-1.5">서비스 선택 (또는 직접 입력)</label>
            <select
              v-model="newSub.service_id"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] focus:outline-none focus:border-[#3182F6] mb-2"
            >
              <option value="">직접 이름 입력</option>
              <option v-for="s in catalog" :key="s.id" :value="s.id">{{ s.name }}</option>
            </select>
            <input
              v-if="!newSub.service_id"
              v-model="newSub.custom_service_name"
              required
              placeholder="서비스명 입력 (예: 챗GPT, 유튜브 프리미엄)"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] focus:outline-none focus:border-[#3182F6]"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-[#4E5968] mb-1.5">월 결제 금액 (KRW)</label>
            <input
              v-model.number="newSub.amount_minor"
              type="number"
              min="100"
              required
              placeholder="14900"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] focus:outline-none focus:border-[#3182F6]"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-bold text-[#4E5968] mb-1.5">결제 주기</label>
              <select
                v-model="newSub.billing_cycle"
                class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-xs text-[#191F28] focus:outline-none focus:border-[#3182F6]"
              >
                <option value="MONTHLY">매월</option>
                <option value="YEARLY">매년</option>
                <option value="WEEKLY">매주</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-bold text-[#4E5968] mb-1.5">다음 갱신일</label>
              <input
                v-model="newSub.next_renewal_at"
                type="date"
                required
                class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-xs text-[#191F28] focus:outline-none focus:border-[#3182F6]"
              />
            </div>
          </div>

          <div class="flex gap-3 pt-2">
            <button
              type="button"
              @click="showAddModal = false"
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
