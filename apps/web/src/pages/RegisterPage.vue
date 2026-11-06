<template>
  <div class="min-h-screen bg-[#08170F] flex">
    <!-- Brand panel -->
    <div class="hidden lg:flex w-[44%] flex-col justify-between p-12 text-white">
      <RouterLink to="/" class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-xl bg-[#C0A04E] flex items-center justify-center font-black text-xl text-[#08170F]">
          M
        </div>
        <div class="leading-none">
          <span class="font-extrabold text-xl tracking-tight">MOA</span>
          <span class="block text-[10px] font-semibold text-[#7E9387] tracking-[0.2em] mt-1">모아 · PRIVATE LEDGER</span>
        </div>
      </RouterLink>
      <div>
        <div class="w-12 h-1 bg-[#C0A04E] rounded-full mb-7"></div>
        <p class="text-4xl font-black tracking-tight leading-tight">첫 장부를<br />펼쳐보세요.</p>
        <div class="mt-8 space-y-4 text-sm">
          <div class="flex items-center gap-3">
            <span class="w-8 h-8 rounded-lg bg-white/10 grid place-items-center text-xs font-black text-[#E9D9A8]">1</span>
            <span class="text-[#A9BCAF]">계정을 만들고 기준 통화를 정합니다</span>
          </div>
          <div class="flex items-center gap-3">
            <span class="w-8 h-8 rounded-lg bg-white/10 grid place-items-center text-xs font-black text-[#E9D9A8]">2</span>
            <span class="text-[#A9BCAF]">자산과 예산, 구독을 차곡차곡 모읍니다</span>
          </div>
          <div class="flex items-center gap-3">
            <span class="w-8 h-8 rounded-lg bg-white/10 grid place-items-center text-xs font-black text-[#E9D9A8]">3</span>
            <span class="text-[#A9BCAF]">순자산의 흐름을 한눈에 확인합니다</span>
          </div>
        </div>
      </div>
      <p class="text-[11px] text-[#5B6960]">* MOA는 은행 연동 없이 사용자가 입력한 기록만 관리합니다.</p>
    </div>

    <!-- Form panel -->
    <div class="flex-1 bg-[#F6F8F4] flex items-center justify-center p-6 py-10 lg:rounded-l-[32px]">
      <div class="max-w-md w-full">
        <div class="lg:hidden flex items-center gap-3 mb-8">
          <div class="w-10 h-10 rounded-xl bg-[#123522] flex items-center justify-center font-black text-xl text-[#E9D9A8]">M</div>
          <span class="font-extrabold text-xl tracking-tight text-[#141F19]">MOA</span>
        </div>

        <h2 class="text-2xl font-black text-[#141F19] tracking-tight">회원가입</h2>
        <p class="text-xs text-[#8A948C] mt-1.5">MOA에서 당신의 금융 생활을 시작하세요.</p>

        <div v-if="error" class="mt-5 p-3.5 rounded-xl bg-[#FBEAEA] border border-[#F2C9C9] text-[#D64545] text-xs font-medium">
          {{ error }}
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-4 mt-6">
          <div>
            <label class="block text-xs font-bold text-[#41524A] mb-1.5">이름 또는 닉네임</label>
            <input
              v-model="displayName"
              type="text"
              required
              placeholder="홍길동"
              class="w-full px-4 py-3.5 rounded-xl border border-[#E2E7DE] text-sm focus:outline-none focus:border-[#1E5A38] focus:ring-2 focus:ring-[#1E5A38]/20 transition-all bg-white"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-[#41524A] mb-1.5">이메일</label>
            <input
              v-model="email"
              type="email"
              required
              placeholder="name@example.com"
              class="w-full px-4 py-3.5 rounded-xl border border-[#E2E7DE] text-sm focus:outline-none focus:border-[#1E5A38] focus:ring-2 focus:ring-[#1E5A38]/20 transition-all bg-white"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-[#41524A] mb-1.5">비밀번호 (8자 이상, 영문·숫자 포함)</label>
            <input
              v-model="password"
              type="password"
              required
              minlength="8"
              placeholder="••••••••"
              class="w-full px-4 py-3.5 rounded-xl border border-[#E2E7DE] text-sm focus:outline-none focus:border-[#1E5A38] focus:ring-2 focus:ring-[#1E5A38]/20 transition-all bg-white"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-[#41524A] mb-1.5">기준 통화</label>
            <select
              v-model="baseCurrency"
              class="w-full px-4 py-3.5 rounded-xl border border-[#E2E7DE] text-sm focus:outline-none focus:border-[#1E5A38] focus:ring-2 focus:ring-[#1E5A38]/20 bg-white transition-all"
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
              <input v-model="termsAccepted" type="checkbox" required class="mt-0.5 rounded text-[#1E5A38] focus:ring-[#1E5A38]" />
              <span class="text-[11px] text-[#41524A] leading-relaxed">
                [필수] <RouterLink to="/terms" class="font-bold text-[#1E5A38] hover:underline">서비스 이용약관</RouterLink> 및
                <RouterLink to="/privacy" class="font-bold text-[#1E5A38] hover:underline">개인정보 처리방침</RouterLink>에 동의합니다.
                (실제 금융기관 거래를 대행하지 않는 순수 개인정보 관리 소프트웨어임을 확인합니다)
              </span>
            </label>
          </div>

          <button
            type="submit"
            :disabled="auth.isLoading || !termsAccepted"
            class="w-full py-3.5 bg-[#123522] hover:bg-[#0C2318] active:bg-[#08170F] disabled:opacity-50 text-white rounded-xl font-bold text-sm shadow-md shadow-[#123522]/25 transition-all mt-2"
          >
            <span v-if="auth.isLoading">가입 처리 중...</span>
            <span v-else>가입하고 시작하기</span>
          </button>
        </form>

        <div class="mt-6 text-center text-xs text-[#8A948C]">
          이미 계정이 있으신가요?
          <RouterLink to="/login" class="font-bold text-[#1E5A38] hover:underline ml-1">
            로그인
          </RouterLink>
        </div>
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
