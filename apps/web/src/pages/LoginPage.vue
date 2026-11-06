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
        <p class="text-4xl font-black tracking-tight leading-tight">다시 오신 것을<br />환영합니다.</p>
        <p class="text-sm text-[#A9BCAF] mt-4 leading-relaxed max-w-xs">로그인하면 어제의 기록부터 오늘의 잔액까지, 내 돈의 흐름이 이어집니다.</p>
      </div>
      <p class="text-[11px] text-[#5B6960]">* MOA는 은행 연동 없이 사용자가 입력한 기록만 관리합니다.</p>
    </div>

    <!-- Form panel -->
    <div class="flex-1 bg-[#F6F8F4] flex items-center justify-center p-6 lg:rounded-l-[32px]">
      <div class="max-w-md w-full">
        <div class="lg:hidden flex items-center gap-3 mb-8">
          <div class="w-10 h-10 rounded-xl bg-[#123522] flex items-center justify-center font-black text-xl text-[#E9D9A8]">M</div>
          <span class="font-extrabold text-xl tracking-tight text-[#141F19]">MOA</span>
        </div>

        <h2 class="text-2xl font-black text-[#141F19] tracking-tight">로그인</h2>
        <p class="text-xs text-[#8A948C] mt-1.5">내 돈의 흐름을 한눈에. 전 세계 어디서나.</p>

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

        <div class="mt-6 text-center text-xs text-[#8A948C]">
          아직 계정이 없으신가요?
          <RouterLink to="/register" class="font-bold text-[#1E5A38] hover:underline ml-1">
            회원가입
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
