<template>
  <div class="min-h-screen bg-slate-50 flex items-center justify-center p-4">
    <div class="max-w-md w-full bg-white rounded-3xl p-8 card-shadow border border-slate-100">
      <!-- Header -->
      <div class="text-center mb-6">
        <div class="w-12 h-12 rounded-2xl bg-brand-600 mx-auto flex items-center justify-center text-white font-black text-2xl shadow-lg shadow-brand-500/30 mb-3">
          BN
        </div>
        <h2 class="text-2xl font-bold text-slate-900 tracking-tight">회원가입</h2>
        <p class="text-xs text-slate-400 mt-1">BN ONE에서 당신의 글로벌 금융 라이프를 시작하세요.</p>
      </div>

      <div v-if="error" class="mb-5 p-3 rounded-xl bg-red-50 border border-red-200 text-red-600 text-xs font-medium">
        {{ error }}
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1.5">이름 또는 닉네임</label>
          <input
            v-model="displayName"
            type="text"
            required
            placeholder="홍길동"
            class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm focus:outline-none focus:border-brand-600 focus:ring-2 focus:ring-brand-500/20 transition-all"
          />
        </div>

        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1.5">이메일</label>
          <input
            v-model="email"
            type="email"
            required
            placeholder="name@example.com"
            class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm focus:outline-none focus:border-brand-600 focus:ring-2 focus:ring-brand-500/20 transition-all"
          />
        </div>

        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1.5">비밀번호 (8자 이상)</label>
          <input
            v-model="password"
            type="password"
            required
            minlength="8"
            placeholder="••••••••"
            class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm focus:outline-none focus:border-brand-600 focus:ring-2 focus:ring-brand-500/20 transition-all"
          />
        </div>

        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1.5">기준 통화</label>
          <select
            v-model="baseCurrency"
            class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm focus:outline-none focus:border-brand-600 focus:ring-2 focus:ring-brand-500/20 bg-white transition-all"
          >
            <option value="KRW">대한민국 원 (₩ KRW)</option>
            <option value="USD">미국 달러 ($ USD)</option>
            <option value="EUR">유로 (€ EUR)</option>
            <option value="JPY">일본 엔 (¥ JPY)</option>
            <option value="GBP">영국 파운드 (£ GBP)</option>
          </select>
        </div>

        <div class="pt-2">
          <label class="flex items-start gap-2.5 cursor-pointer">
            <input v-model="termsAccepted" type="checkbox" required class="mt-0.5 rounded text-brand-600 focus:ring-brand-500" />
            <span class="text-[11px] text-slate-600 leading-tight">
              [필수] 서비스 이용약관 및 개인정보 처리방침에 동의합니다. (실제 금융기관 거래를 대행하지 않는 순수 개인정보 관리 소프트웨어임을 확인합니다)
            </span>
          </label>
        </div>

        <button
          type="submit"
          :disabled="auth.isLoading || !termsAccepted"
          class="w-full py-3.5 bg-brand-600 hover:bg-brand-700 active:bg-brand-800 disabled:opacity-50 text-white rounded-xl font-bold text-sm shadow-md shadow-brand-600/25 transition-all mt-2"
        >
          <span v-if="auth.isLoading">가입 처리 중...</span>
          <span v-else>가입하고 시작하기</span>
        </button>
      </form>

      <div class="mt-6 text-center text-xs text-slate-500">
        이미 계정이 있으신가요?
        <RouterLink to="/login" class="font-bold text-brand-600 hover:underline ml-1">
          로그인
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter, RouterLink } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const router = useRouter();
const auth = useAuthStore();

const displayName = ref('');
const email = ref('');
const password = ref('');
const baseCurrency = ref('KRW');
const termsAccepted = ref(false);
const error = ref('');

async function handleSubmit() {
  error.value = '';
  try {
    await auth.register({
      email: email.value,
      password: password.value,
      display_name: displayName.value,
      base_currency: baseCurrency.value,
      locale: 'ko-KR',
      country_code: 'KR',
    });
    router.push('/onboarding');
  } catch (err: any) {
    error.value = err.message || '회원가입에 실패했습니다.';
  }
}
</script>
