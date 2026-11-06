<template>
  <div class="min-h-screen bg-[#F6F8F4] flex">
    <!-- Form -->
    <div class="flex-1 flex items-center justify-center p-6 py-10 lg:p-16">
      <div class="max-w-md w-full">
        <RouterLink to="/" class="inline-flex items-center gap-3 mb-10">
          <div class="w-10 h-10 rounded-xl bg-[#123522] flex items-center justify-center font-black text-xl text-[#E9D9A8]">모</div>
          <span class="font-extrabold text-xl tracking-tight text-[#141F19]">MOA</span>
        </RouterLink>

        <h2 class="text-[32px] font-black text-[#141F19] tracking-tight leading-tight">첫 장부 펴는 데<br />1분이면 됩니다.</h2>
        <p class="text-[13px] text-[#8A948C] mt-3">계정 만들고, 어제 쓴 돈부터 적어보세요.</p>

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

        <div class="mt-6 text-xs text-[#8A948C]">
          이미 계정이 있으신가요?
          <RouterLink to="/login" class="font-bold text-[#1E5A38] hover:underline ml-1">
            로그인
          </RouterLink>
        </div>
      </div>
    </div>

    <!-- Ledger rail -->
    <div class="hidden lg:flex w-[380px] shrink-0 bg-[#0C2318] p-8 flex-col justify-end gap-4" aria-hidden="true">
      <div class="rounded-2xl bg-white/[0.06] border border-white/10 p-5 rotate-1">
        <p class="text-[11px] font-bold text-[#7E9387]">가입 첫 주에 보통 적는 것</p>
        <p class="text-[13px] text-white mt-2 leading-loose">주계좌 잔액 · 교통카드 · 점심값 3건<br />구독 2건 · 비상금 통장</p>
      </div>
      <div class="rounded-2xl bg-[#C0A04E] p-5 -rotate-1 mr-8">
        <p class="text-[11px] font-bold text-[#0C2318]/70">걱정 마세요</p>
        <p class="text-xl font-black text-[#0C2318] mt-1">나중에 다 고쳐집니다</p>
        <p class="text-[11px] text-[#0C2318]/70 mt-2">카테고리도, 계정도, 통화도</p>
      </div>
      <p class="text-[11px] text-[#5B6960] pt-2">* MOA는 은행 연동 없이 직접 적은 기록만 보여줍니다.</p>
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
