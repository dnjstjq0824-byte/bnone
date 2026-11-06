<template>
  <div class="min-h-screen bg-[#F6F8F4] flex">
    <!-- Form -->
    <div class="flex-1 flex items-center justify-center p-6 lg:p-16">
      <div class="max-w-md w-full">
        <RouterLink to="/" class="inline-flex items-center gap-3 mb-10">
          <div class="w-10 h-10 rounded-xl bg-[#123522] flex items-center justify-center font-black text-xl text-[#E9D9A8]">모</div>
          <span class="font-extrabold text-xl tracking-tight text-[#141F19]">MOA</span>
        </RouterLink>

        <h2 class="text-[32px] font-black text-[#141F19] tracking-tight leading-tight">어제 쓴 돈,<br />기억나세요?</h2>
        <p class="text-[13px] text-[#8A948C] mt-3">로그인하면 어제 적어둔 장부부터 이어집니다.</p>

        <!-- Demo banner (explicit demo builds only) -->
        <div v-if="isDemoMode" class="mt-6 p-4 rounded-2xl bg-[#EBF2EC] border border-[#CBDACF] space-y-2">
          <div class="flex items-center justify-between">
            <span class="text-[#1E5A38] text-sm font-bold">데모 체험 계정 안내</span>
            <span class="px-2 py-0.5 rounded-full bg-[#1E5A38]/10 text-[#1E5A38] text-[11px] font-bold">브라우저 로컬 데이터</span>
          </div>
          <p class="text-[11px] text-[#41524A] leading-relaxed">
            현재 웹은 <strong>정적 웹 + 브라우저 내 로컬 목업 데모</strong>로 실행 중이며, 모든 데이터는 브라우저에만 저장됩니다.
          </p>
          <div class="text-xs text-[#41524A] space-y-1 pt-1 font-medium">
            <div class="flex items-center justify-between">
              <span class="text-[#8A948C]">데모 이메일</span>
              <code class="font-bold text-[#141F19] bg-white px-2 py-0.5 rounded border border-[#E2E7DE]">test@example.com</code>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-[#8A948C]">데모 비밀번호</span>
              <code class="font-bold text-[#141F19] bg-white px-2 py-0.5 rounded border border-[#E2E7DE]">password123</code>
            </div>
          </div>
          <div class="pt-2 flex items-center gap-2">
            <button
              type="button"
              @click="fillTestAccount"
              class="flex-1 py-2 bg-white hover:bg-[#ECEFE8] text-[#1E5A38] border border-[#CBDACF] rounded-xl text-xs font-bold transition-colors"
            >
              계정 자동 입력
            </button>
            <button
              type="button"
              @click="handleDemoLogin"
              :disabled="auth.isLoading"
              class="flex-1 py-2 bg-[#123522] hover:bg-[#0C2318] text-white rounded-xl text-xs font-bold shadow-sm transition-colors"
            >
              1초 만에 바로 시작
            </button>
          </div>
        </div>

        <!-- Error alert -->
        <div v-if="error" class="mt-5 p-3.5 rounded-xl bg-[#FBEAEA] border border-[#F2C9C9] text-[#D64545] text-xs font-medium">
          {{ error }}
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-4 mt-6">
          <div>
            <label class="block text-xs font-bold text-[#41524A] mb-1.5">이메일</label>
            <input
              v-model="email"
              type="email"
              required
              autocomplete="username"
              placeholder="name@example.com"
              class="w-full px-4 py-3.5 rounded-xl border border-[#E2E7DE] text-sm focus:outline-none focus:border-[#1E5A38] focus:ring-2 focus:ring-[#1E5A38]/20 transition-all bg-white"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-[#41524A] mb-1.5">비밀번호</label>
            <input
              v-model="password"
              type="password"
              required
              autocomplete="current-password"
              placeholder="••••••••"
              class="w-full px-4 py-3.5 rounded-xl border border-[#E2E7DE] text-sm focus:outline-none focus:border-[#1E5A38] focus:ring-2 focus:ring-[#1E5A38]/20 transition-all bg-white"
            />
          </div>

          <button
            type="submit"
            :disabled="auth.isLoading"
            class="w-full py-3.5 bg-[#123522] hover:bg-[#0C2318] active:bg-[#08170F] disabled:opacity-50 text-white rounded-xl font-bold text-sm shadow-md shadow-[#123522]/25 transition-all mt-2"
          >
            <span v-if="auth.isLoading">로그인 중...</span>
            <span v-else>로그인</span>
          </button>
        </form>

        <div class="mt-6 text-xs text-[#8A948C]">
          아직 계정이 없으신가요?
          <RouterLink to="/register" class="font-bold text-[#1E5A38] hover:underline ml-1">
            회원가입
          </RouterLink>
        </div>
      </div>
    </div>

    <!-- Ledger rail -->
    <div class="hidden lg:flex w-[380px] shrink-0 bg-[#0C2318] p-8 flex-col justify-end gap-4" aria-hidden="true">
      <div class="rounded-2xl bg-white/[0.06] border border-white/10 p-5 -rotate-1">
        <p class="text-[11px] font-bold text-[#7E9387]">어제 장부 · 4건</p>
        <p class="text-xl font-black text-white mt-1">−₩50,300</p>
        <p class="text-[11px] text-[#A9BCAF] mt-2 leading-relaxed">택시 14,800 · 점심 12,000 · 서점 23,500 · 커피 0원(얻어먹음)</p>
      </div>
      <div class="rounded-2xl bg-[#C0A04E] p-5 rotate-1 ml-8">
        <p class="text-[11px] font-bold text-[#0C2318]/70">11월 교통비</p>
        <p class="text-xl font-black text-[#0C2318] mt-1">예산의 62%</p>
        <p class="text-[11px] text-[#0C2318]/70 mt-2">남은 달력 24일 · 페이스 양호</p>
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

// Initial form inputs must start empty (not pre-filled)
const email = ref('');
const password = ref('');
const error = ref('');

const isDemoMode = ref(import.meta.env.VITE_MOCK_ENABLED === 'true');

function fillTestAccount() {
  email.value = 'test@example.com';
  password.value = 'password123';
}

async function handleDemoLogin() {
  fillTestAccount();
  await handleSubmit();
}

async function handleSubmit() {
  error.value = '';
  try {
    const res = await auth.login(email.value, password.value);
    if (res.onboarding_required) {
      router.push('/onboarding');
    } else {
      router.push('/app');
    }
  } catch (err: any) {
    error.value = err.message || '로그인에 실패했습니다. 이메일과 비밀번호를 확인하세요.';
  }
}
</script>
