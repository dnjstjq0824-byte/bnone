<template>
  <div class="max-w-xl mx-auto space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold tracking-tight text-[#191F28]">새 거래 기록</h1>
      <p class="text-sm text-[#8B95A1] mt-1">지출, 수입 또는 계정 간 자금 이동을 기록해요.</p>
    </div>

    <!-- Mandatory Disclaimer (Article 1.5) -->
    <div class="p-4 rounded-2xl bg-[#E8F3FF] flex items-start gap-3">
      <span class="text-[#3182F6] text-base">ℹ️</span>
      <p class="text-xs text-[#1B64DA] font-semibold leading-relaxed">
        BN ONE의 모든 거래 기록은 실제 금융기관 송금이나 결제를 실행하지 않는 순수 장부 기록입니다.
      </p>
    </div>

    <!-- Type Tabs (Toss Segmented Pill Control) -->
    <div class="flex p-1.5 bg-[#F2F4F6] rounded-2xl">
      <button
        @click="activeTab = 'EXPENSE'"
        class="flex-1 py-3 text-xs font-bold rounded-xl transition-all"
        :class="activeTab === 'EXPENSE' ? 'bg-white text-[#191F28] shadow-sm' : 'text-[#8B95A1] hover:text-[#4E5968]'"
      >
        지출
      </button>
      <button
        @click="activeTab = 'INCOME'"
        class="flex-1 py-3 text-xs font-bold rounded-xl transition-all"
        :class="activeTab === 'INCOME' ? 'bg-white text-[#3182F6] shadow-sm' : 'text-[#8B95A1] hover:text-[#4E5968]'"
      >
        수입
      </button>
      <button
        @click="activeTab = 'TRANSFER'"
        class="flex-1 py-3 text-xs font-bold rounded-xl transition-all"
        :class="activeTab === 'TRANSFER' ? 'bg-white text-[#8061D9] shadow-sm' : 'text-[#8B95A1] hover:text-[#4E5968]'"
      >
        계정 간 이동
      </button>
    </div>

    <!-- Big Amount Card (Toss Keypad/Amount Style) -->
    <div class="bg-white p-6 rounded-3xl space-y-4">
      <span class="text-xs font-bold text-[#8B95A1]">금액 입력</span>
      <div class="flex items-center gap-2 border-b border-[#F2F4F6] pb-4">
        <input
          v-if="activeTab !== 'TRANSFER'"
          v-model.number="form.amount_minor"
          type="number"
          min="1"
          required
          placeholder="0"
          class="w-full text-3xl font-extrabold text-[#191F28] tracking-tight bg-transparent focus:outline-none placeholder-[#D1D6DB] tabular-nums"
        />
        <input
          v-else
          v-model.number="transferForm.amount_minor"
          type="number"
          min="1"
          required
          placeholder="0"
          class="w-full text-3xl font-extrabold text-[#191F28] tracking-tight bg-transparent focus:outline-none placeholder-[#D1D6DB] tabular-nums"
        />
        <span class="text-2xl font-bold text-[#191F28] flex-shrink-0">원</span>
      </div>

      <!-- Quick Amount Presets -->
      <div class="flex flex-wrap gap-2 pt-2">
        <button
          v-for="preset in [10000, 50000, 100000, 500000, 1000000]"
          :key="preset"
          type="button"
          @click="addPreset(preset)"
          class="px-3 py-1.5 bg-[#F2F4F6] hover:bg-[#E5E8EB] text-[#4E5968] rounded-xl text-xs font-bold transition-colors"
        >
          +{{ formatPreset(preset) }}
        </button>
      </div>
    </div>

    <!-- Form Container -->
    <div class="bg-white p-6 rounded-3xl">
      <!-- EXPENSE or INCOME FORM -->
      <form v-if="activeTab !== 'TRANSFER'" @submit.prevent="submitExpenseIncome" class="space-y-4">
        <div>
          <label class="block text-xs font-bold text-[#4E5968] mb-1.5">연결 계정</label>
          <select
            v-model="form.account_id"
            required
            class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] font-medium focus:outline-none focus:border-[#3182F6]"
          >
            <option v-for="acc in finance.accounts" :key="acc.id" :value="acc.id">
              {{ acc.name }} ({{ acc.currency_code }}) — 잔액: {{ finance.formatCurrency(acc.calculated_balance_minor, acc.currency_code) }}
            </option>
          </select>
        </div>

        <div>
          <label class="block text-xs font-bold text-[#4E5968] mb-1.5">카테고리</label>
          <select
            v-model="form.category_id"
            required
            class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] font-medium focus:outline-none focus:border-[#3182F6]"
          >
            <option v-for="cat in filteredCategories" :key="cat.id" :value="cat.id">
              {{ cat.name }}
            </option>
          </select>
        </div>

        <div>
          <label class="block text-xs font-bold text-[#4E5968] mb-1.5">설명 또는 상호명</label>
          <input
            v-model="form.description"
            type="text"
            required
            placeholder="예: 점심 식사, 스타벅스, 월급"
            class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] placeholder-[#8B95A1] focus:outline-none focus:border-[#3182F6]"
          />
        </div>

        <div>
          <label class="block text-xs font-bold text-[#4E5968] mb-1.5">거래 날짜 및 시각</label>
          <input
            v-model="form.occurred_at"
            type="datetime-local"
            class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] focus:outline-none focus:border-[#3182F6]"
          />
        </div>

        <div>
          <label class="block text-xs font-bold text-[#4E5968] mb-1.5">메모 (선택)</label>
          <textarea
            v-model="form.notes"
            rows="2"
            placeholder="추가 메모 입력"
            class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] placeholder-[#8B95A1] focus:outline-none focus:border-[#3182F6]"
          ></textarea>
        </div>

        <button
          type="submit"
          :disabled="isSubmitting"
          class="w-full py-4 bg-[#3182F6] hover:bg-[#1B64DA] text-white rounded-2xl font-bold text-base shadow-sm transition-all mt-4"
        >
          {{ isSubmitting ? '저장 중...' : '기록 완료하기' }}
        </button>
      </form>

      <!-- TRANSFER FORM -->
      <form v-else @submit.prevent="submitTransfer" class="space-y-4">
        <div>
          <label class="block text-xs font-bold text-[#4E5968] mb-1.5">출발 계정 (출금)</label>
          <select
            v-model="transferForm.source_account_id"
            required
            class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] font-medium focus:outline-none focus:border-[#8061D9]"
          >
            <option v-for="acc in finance.accounts" :key="acc.id" :value="acc.id">
              {{ acc.name }} ({{ acc.currency_code }})
            </option>
          </select>
        </div>

        <div>
          <label class="block text-xs font-bold text-[#4E5968] mb-1.5">도착 계정 (입금)</label>
          <select
            v-model="transferForm.destination_account_id"
            required
            class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] font-medium focus:outline-none focus:border-[#8061D9]"
          >
            <option
              v-for="acc in finance.accounts"
              :key="acc.id"
              :value="acc.id"
              :disabled="acc.id === transferForm.source_account_id"
            >
              {{ acc.name }} ({{ acc.currency_code }})
            </option>
          </select>
        </div>

        <div>
          <label class="block text-xs font-bold text-[#4E5968] mb-1.5">이동 사유 / 설명</label>
          <input
            v-model="transferForm.description"
            type="text"
            placeholder="예: 비상금 저축, 카드 대금 결제용 이동"
            class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] placeholder-[#8B95A1] focus:outline-none focus:border-[#8061D9]"
          />
        </div>

        <button
          type="submit"
          :disabled="isSubmitting"
          class="w-full py-4 bg-[#8061D9] hover:bg-[#6c4fc2] text-white rounded-2xl font-bold text-base shadow-sm transition-all mt-4"
        >
          {{ isSubmitting ? '저장 중...' : '계정 이동 완료하기' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { apiFetch } from '../api/client';
import { useFinanceStore } from '../stores/finance';

const router = useRouter();
const finance = useFinanceStore();

const activeTab = ref<'EXPENSE' | 'INCOME' | 'TRANSFER'>('EXPENSE');
const isSubmitting = ref(false);
const categories = ref<any[]>([]);

const form = ref({
  account_id: '',
  amount_minor: 0,
  category_id: '',
  description: '',
  merchant_name: '',
  occurred_at: new Date().toISOString().slice(0, 16),
  notes: '',
});

const transferForm = ref({
  source_account_id: '',
  destination_account_id: '',
  amount_minor: 0,
  description: '수동 계정 이동',
});

const filteredCategories = computed(() => {
  return categories.value.filter((c) => c.category_type === activeTab.value);
});

function addPreset(delta: number) {
  if (activeTab.value === 'TRANSFER') {
    transferForm.value.amount_minor = (transferForm.value.amount_minor || 0) + delta;
  } else {
    form.value.amount_minor = (form.value.amount_minor || 0) + delta;
  }
}

function formatPreset(val: number) {
  if (val >= 10000) return `${val / 10000}만`;
  return `${val}원`;
}

onMounted(async () => {
  await finance.loadDashboard();
  if (finance.accounts.length > 0) {
    form.value.account_id = finance.accounts[0].id;
    transferForm.value.source_account_id = finance.accounts[0].id;
    if (finance.accounts.length > 1) {
      transferForm.value.destination_account_id = finance.accounts[1].id;
    }
  }

  const cats = await apiFetch('/categories');
  categories.value = cats;
  if (filteredCategories.value.length > 0) {
    form.value.category_id = filteredCategories.value[0].id;
  }
});

async function submitExpenseIncome() {
  isSubmitting.value = true;
  try {
    await apiFetch('/transactions', {
      method: 'POST',
      body: JSON.stringify({
        transaction_type: activeTab.value,
        account_id: form.value.account_id,
        amount_minor: form.value.amount_minor,
        occurred_at: new Date(form.value.occurred_at).toISOString(),
        description: form.value.description,
        merchant_name: form.value.merchant_name,
        notes: form.value.notes,
        splits: [
          {
            category_id: form.value.category_id,
            amount_minor: form.value.amount_minor,
          },
        ],
      }),
    });
    await finance.loadDashboard();
    router.push('/app/transactions');
  } catch (err: any) {
    alert(err.message || '저장에 실패했습니다.');
  } finally {
    isSubmitting.value = false;
  }
}

async function submitTransfer() {
  if (transferForm.value.source_account_id === transferForm.value.destination_account_id) {
    alert('출발 계정과 도착 계정은 서로 달라야 합니다.');
    return;
  }
  isSubmitting.value = true;
  try {
    await apiFetch('/transfers/manual', {
      method: 'POST',
      body: JSON.stringify({
        source_account_id: transferForm.value.source_account_id,
        destination_account_id: transferForm.value.destination_account_id,
        source_amount_minor: transferForm.value.amount_minor,
        destination_amount_minor: transferForm.value.amount_minor,
        description: transferForm.value.description,
      }),
    });
    await finance.loadDashboard();
    router.push('/app/transactions');
  } catch (err: any) {
    alert(err.message || '이체 기록 저장에 실패했습니다.');
  } finally {
    isSubmitting.value = false;
  }
}
</script>
