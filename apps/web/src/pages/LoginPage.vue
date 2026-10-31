<template>
  <div class="min-h-screen bg-[#F9FAFB] flex items-center justify-center p-4">
    <div class="max-w-md w-full bg-white rounded-3xl p-8 card-shadow border border-[#F2F4F6]">
      <!-- Header -->
      <div class="text-center mb-6">
        <div class="w-14 h-14 rounded-2xl bg-[#3182F6] mx-auto flex items-center justify-center text-white font-black text-2xl shadow-lg shadow-[#3182F6]/25 mb-3">
          BN
        </div>
        <h2 class="text-2xl font-bold text-[#191F28] tracking-tight">로그인</h2>
        <p class="text-xs text-[#8B95A1] mt-1.5">내 돈의 흐름을 한눈에. 전 세계 어디서나.</p>
      </div>

      <!-- Test Account Helper Banner (Only displayed in Demo Mode) -->
      <div v-if="isDemoMode" class="mb-6 p-4 rounded-2xl bg-[#E8F3FF] border border-[#D0E5FF] space-y-2">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-1.5">
            <span class="text-[#3182F6] text-sm font-bold">💡 데모 체험 계정 안내</span>
          </div>
          <span class="px-2 py-0.5 rounded-full bg-[#3182F6]/10 text-[#3182F6] text-[11px] font-bold">브라우저 로컬 데이터</span>
        </div>
        <p class="text-[11px] text-[#4E5968] leading-relaxed">
          현재 웹은 <strong>정적 웹 + 브라우저 내 로컬 목업 데모</strong>로 실행 중이며, 모든 데이터는 브라우저에만 저장됩니다.
        </p>
        <div class="text-xs text-[#4E5968] space-y-1 pt-1 font-medium">
          <div class="flex items-center justify-between">
            <span class="text-[#8B95A1]">데모 이메일</span>
            <code class="font-bold text-[#191F28] bg-white px-2 py-0.5 rounded border border-[#E5E8EB]">test@example.com</code>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-[#8B95A1]">데모 비밀번호</span>
            <code class="font-bold text-[#191F28] bg-white px-2 py-0.5 rounded border border-[#E5E8EB]">password123</code>
          </div>
        </div>
        <div class="pt-2 flex items-center gap-2">
          <button
            type="button"
            @click="fillTestAccount"
            class="flex-1 py-2 bg-white hover:bg-[#F2F4F6] text-[#3182F6] border border-[#D0E5FF] rounded-xl text-xs font-bold transition-colors"
          >
            계정 자동 입력
          </button>
          <button
            type="button"
            @click="handleDemoLogin"
            :disabled="auth.isLoading"
            class="flex-1 py-2 bg-[#3182F6] hover:bg-[#1B64DA] text-white rounded-xl text-xs font-bold shadow-sm transition-colors"
          >
            1초 만에 바로 시작
          </button>
        </div>
      </div>

      <!-- Error alert -->
      <div v-if="error" class="mb-5 p-3.5 rounded-xl bg-[#FEECEB] border border-[#FDD8D6] text-[#F04452] text-xs font-medium">
        {{ error }}
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label class="block text-xs font-bold text-[#4E5968] mb-1.5">이메일</label>
          <input
            v-model="email"
            type="email"
            required
            autocomplete="username"
            placeholder="name@example.com"
            class="w-full px-4 py-3.5 rounded-xl border border-[#E5E8EB] text-sm focus:outline-none focus:border-[#3182F6] focus:ring-2 focus:ring-[#3182F6]/20 transition-all bg-[#F9FAFB] focus:bg-white"
          />
        </div>

        <div>
          <label class="block text-xs font-bold text-[#4E5968] mb-1.5">비밀번호</label>
          <input
            v-model="password"
            type="password"
            required
            autocomplete="current-password"
            placeholder="••••••••"
            class="w-full px-4 py-3.5 rounded-xl border border-[#E5E8EB] text-sm focus:outline-none focus:border-[#3182F6] focus:ring-2 focus:ring-[#3182F6]/20 transition-all bg-[#F9FAFB] focus:bg-white"
          />
        </div>

        <button
          type="submit"
          :disabled="auth.isLoading"
          class="w-full py-3.5 bg-[#3182F6] hover:bg-[#1B64DA] active:bg-[#1553B7] disabled:opacity-50 text-white rounded-xl font-bold text-sm shadow-md shadow-[#3182F6]/25 transition-all mt-2"
        >
          <span v-if="auth.isLoading">로그인 중...</span>
          <span v-else>로그인</span>
        </button>
      </form>

      <div class="mt-6 text-center text-xs text-[#8B95A1]">
        아직 계정이 없으신가요?
        <RouterLink to="/register" class="font-bold text-[#3182F6] hover:underline ml-1">
          회원가입
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

// Initial form inputs must start empty (not pre-filled)
const email = ref('');
const password = ref('');
const error = ref('');

const isDemoMode = ref(import.meta.env.VITE_DEMO_MODE === 'true' || import.meta.env.VITE_DEMO_MODE === undefined);

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
