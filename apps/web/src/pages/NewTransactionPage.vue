<template>
  <div class="max-w-xl mx-auto space-y-6">
    <div>
      <h1 class="text-2xl font-black text-slate-900 tracking-tight">새 거래 기록</h1>
      <p class="text-xs text-slate-500 mt-0.5">지출, 수입 또는 계정 간 자금 이동을 기록합니다.</p>
    </div>

    <!-- Mandatory Disclaimer (Article 1.5) -->
    <div class="p-3.5 rounded-2xl bg-blue-50/70 border border-blue-100 flex items-start gap-2.5">
      <span class="text-blue-600 text-sm">ℹ️</span>
      <p class="text-[11px] text-blue-900 leading-relaxed">
        BN ONE의 모든 거래 기록은 실제 금융기관 송금이나 결제를 실행하지 않는 장부 기록입니다.
      </p>
    </div>

    <!-- Type Tabs -->
    <div class="flex p-1 bg-slate-100 rounded-2xl">
      <button
        @click="activeTab = 'EXPENSE'"
        class="flex-1 py-2.5 text-xs font-bold rounded-xl transition-all"
        :class="activeTab === 'EXPENSE' ? 'bg-white text-red-600 shadow-sm' : 'text-slate-500 hover:text-slate-800'"
      >
        지출
      </button>
      <button
        @click="activeTab = 'INCOME'"
        class="flex-1 py-2.5 text-xs font-bold rounded-xl transition-all"
        :class="activeTab === 'INCOME' ? 'bg-white text-emerald-600 shadow-sm' : 'text-slate-500 hover:text-slate-800'"
      >
        수입
      </button>
      <button
        @click="activeTab = 'TRANSFER'"
        class="flex-1 py-2.5 text-xs font-bold rounded-xl transition-all"
        :class="activeTab === 'TRANSFER' ? 'bg-white text-indigo-600 shadow-sm' : 'text-slate-500 hover:text-slate-800'"
      >
        계정 간 이동
      </button>
    </div>

    <!-- Form Container -->
    <div class="bg-white p-6 rounded-3xl border border-slate-200 card-shadow">
      <!-- EXPENSE or INCOME FORM -->
      <form v-if="activeTab !== 'TRANSFER'" @submit.prevent="submitExpenseIncome" class="space-y-4">
        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1">금액</label>
          <div class="relative">
            <input
              v-model.number="form.amount_minor"
              type="number"
              min="1"
              required
              placeholder="0"
              class="w-full px-4 py-3 rounded-xl border border-slate-200 text-lg font-black focus:outline-none focus:border-brand-600"
            />
            <span class="absolute right-4 top-3.5 text-xs font-bold text-slate-400">
              {{ selectedAccountCurrency }}
            </span>
          </div>
        </div>

        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1">연결 계정</label>
          <select
            v-model="form.account_id"
            required
            class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm bg-white"
          >
            <option v-for="acc in finance.accounts" :key="acc.id" :value="acc.id">
              {{ acc.name }} ({{ acc.currency_code }}) — 잔액: {{ finance.formatCurrency(acc.calculated_balance_minor, acc.currency_code) }}
            </option>
          </select>
        </div>

        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1">카테고리</label>
          <select
            v-model="form.category_id"
            required
            class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm bg-white"
          >
            <option v-for="cat in filteredCategories" :key="cat.id" :value="cat.id">
              {{ cat.name }}
            </option>
          </select>
        </div>

        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1">설명 또는 상호명</label>
          <input
            v-model="form.description"
            type="text"
            required
            placeholder="예: 점심 식사, 스타벅스, 월급"
            class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm focus:outline-none focus:border-brand-600"
          />
        </div>

        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1">거래 날짜 및 시각</label>
          <input
            v-model="form.occurred_at"
            type="datetime-local"
            class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm bg-white"
          />
        </div>

        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1">메모 (선택)</label>
          <textarea
            v-model="form.notes"
            rows="2"
            placeholder="추가 메모 입력"
            class="w-full px-4 py-2.5 rounded-xl border border-slate-200 text-sm"
          ></textarea>
        </div>

        <button
          type="submit"
          :disabled="isSubmitting"
          class="w-full py-3.5 bg-brand-600 hover:bg-brand-700 text-white rounded-xl font-bold text-sm shadow-md shadow-brand-600/25 transition-all mt-4"
        >
          {{ isSubmitting ? '저장 중...' : '거래 저장하기' }}
        </button>
      </form>

      <!-- TRANSFER FORM -->
      <form v-else @submit.prevent="submitTransfer" class="space-y-4">
        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1">출발 계정 (출금)</label>
          <select
            v-model="transferForm.source_account_id"
            required
            class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm bg-white"
          >
            <option v-for="acc in finance.accounts" :key="acc.id" :value="acc.id">
              {{ acc.name }} ({{ acc.currency_code }})
            </option>
          </select>
        </div>

        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1">도착 계정 (입금)</label>
          <select
            v-model="transferForm.destination_account_id"
            required
            class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm bg-white"
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
          <label class="block text-xs font-bold text-slate-700 mb-1">이동 금액</label>
          <input
            v-model.number="transferForm.amount_minor"
            type="number"
            min="1"
            required
            placeholder="0"
            class="w-full px-4 py-3 rounded-xl border border-slate-200 text-lg font-black focus:outline-none focus:border-brand-600"
          />
        </div>

        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1">이동 사유 / 설명</label>
          <input
            v-model="transferForm.description"
            type="text"
            placeholder="예: 비상금 저축, 카드 대금 이체"
            class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm"
          />
        </div>

        <button
          type="submit"
          :disabled="isSubmitting"
          class="w-full py-3.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl font-bold text-sm shadow-md shadow-indigo-600/25 transition-all mt-4"
        >
          {{ isSubmitting ? '저장 중...' : '계정 이동 저장하기' }}
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

const selectedAccountCurrency = computed(() => {
  const acc = finance.accounts.find((a) => a.id === form.value.account_id);
  return acc ? acc.currency_code : 'KRW';
});

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

