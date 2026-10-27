<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { apiFetch } from '../api/client';
const props = defineProps<{ kind: 'asset' | 'liability' }>();
const emit = defineEmits<{ close: []; saved: [] }>();
const dialog = ref<HTMLDialogElement>();
const name = ref('');
const type = ref('OTHER');
const currency = ref('KRW');
const amount = ref('');
const notes = ref('');
const busy = ref(false);
const error = ref('');
const isAsset = computed(() => props.kind === 'asset');
const types = computed(() => isAsset.value ? [['REAL_ESTATE', '부동산'], ['VEHICLE', '차량'], ['INVESTMENT_MANUAL', '투자 자산'], ['OTHER', '기타 자산']] : [['MORTGAGE', '주택담보대출'], ['PERSONAL_LOAN', '신용대출'], ['CREDIT_CARD_DEBT', '카드 대금'], ['STUDENT_LOAN', '학자금'], ['OTHER', '기타 부채']]);
onMounted(() => dialog.value?.showModal());
async function save() {
  if (busy.value) return;
  error.value = '';
  const multiplier = ['KRW', 'JPY'].includes(currency.value) ? 1 : 100;
  const value = Number(amount.value);
  const minor = Math.round(value * multiplier);
  if (!name.value.trim() || !amount.value.trim() || !Number.isFinite(value) || value < 0 || !Number.isSafeInteger(minor) || Math.abs(value * multiplier - minor) > .000001) {
    error.value = '이름과 올바른 금액을 입력해주세요. 원·엔은 정수, 달러·유로는 소수 둘째 자리까지 입력할 수 있어요.';
    return;
  }
  busy.value = true;
  try {
    await apiFetch(isAsset.value ? '/assets' : '/liabilities', { method: 'POST', body: JSON.stringify({ name: name.value.trim(), currency_code: currency.value, notes: notes.value, ...(isAsset.value ? { asset_type: type.value, estimated_value_minor: minor } : { liability_type: type.value, outstanding_amount_minor: minor }) }) });
    emit('saved');
  } catch { error.value = '저장하지 못했어요. 연결을 확인하고 다시 시도해주세요.'; }
  finally { busy.value = false; }
}
</script>
<template>
  <dialog ref="dialog" class="asset-dialog" aria-labelledby="entry-title" @cancel.prevent="!busy && emit('close')">
    <form @submit.prevent="save" class="space-y-5">
      <h2 id="entry-title" class="text-xl font-bold">{{ isAsset ? '보유 자산 등록' : '부채 등록' }}</h2>
      <p class="text-sm text-slate-500">{{ isAsset ? '현재 예상 가치를 기록해주세요.' : '현재 남아 있는 금액을 기록해주세요.' }}</p>
      <label class="entry-field">이름<input v-model="name" required maxlength="100" autofocus :disabled="busy" placeholder="예: 내 자동차, 주택담보대출" /></label>
      <label class="entry-field">종류<select v-model="type" :disabled="busy"><option v-for="option in types" :key="option[0]" :value="option[0]">{{ option[1] }}</option></select></label>
      <label class="entry-field">통화<select v-model="currency" :disabled="busy"><option value="KRW">대한민국 원 (KRW)</option><option value="USD">미국 달러 (USD)</option><option value="EUR">유로 (EUR)</option><option value="JPY">일본 엔 (JPY)</option></select></label>
      <label class="entry-field">{{ isAsset ? '예상 가치' : '남은 금액' }} ({{ currency }})<input v-model="amount" type="number" min="0" :step="['KRW', 'JPY'].includes(currency) ? 1 : 0.01" required :disabled="busy" placeholder="0" /></label>
      <label class="entry-field">메모 (선택)<input v-model="notes" maxlength="500" :disabled="busy" /></label>
      <p v-if="error" role="alert" class="text-sm text-red-600">{{ error }}</p>
      <div class="flex gap-3"><button type="button" class="soft-button flex-1" :disabled="busy" @click="emit('close')">취소</button><button class="primary-button flex-1" :disabled="busy">{{ busy ? '저장 중…' : '등록하기' }}</button></div>
    </form>
  </dialog>
</template>
<style scoped>
.asset-dialog { border: 0; border-radius: 24px; padding: 28px; width: min(460px, calc(100vw - 32px)); max-height: calc(100dvh - 40px); overflow-y: auto; color: #191f28; }
.asset-dialog::backdrop { background: rgb(25 31 40 / 45%); }
.entry-field { display: grid; gap: 8px; font-size: 14px; font-weight: 600; }
.entry-field input,.entry-field select { border: 1px solid #d1d6db; border-radius: 12px; padding: 12px; width: 100%; background: white; font-weight: 400; }
</style>
