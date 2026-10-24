<template>
  <div class="min-h-screen bg-slate-50 flex items-center justify-center p-4">
    <div class="max-w-md w-full bg-white rounded-3xl p-8 card-shadow border border-slate-100">
      <!-- Header -->
      <div class="text-center mb-8">
        <div class="w-12 h-12 rounded-2xl bg-brand-600 mx-auto flex items-center justify-center text-white font-black text-2xl shadow-lg shadow-brand-500/30 mb-3">
          BN
        </div>
        <h2 class="text-2xl font-bold text-slate-900 tracking-tight">로그인</h2>
        <p class="text-xs text-slate-400 mt-1">내 돈의 흐름을 한눈에. 전 세계 어디서나.</p>
      </div>

      <!-- Error alert -->
      <div v-if="error" class="mb-5 p-3 rounded-xl bg-red-50 border border-red-200 text-red-600 text-xs font-medium">
        {{ error }}
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-4">
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
          <label class="block text-xs font-bold text-slate-700 mb-1.5">비밀번호</label>
          <input
            v-model="password"
            type="password"
            required
            placeholder="••••••••"
            class="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm focus:outline-none focus:border-brand-600 focus:ring-2 focus:ring-brand-500/20 transition-all"
          />
        </div>

        <button
          type="submit"
          :disabled="auth.isLoading"
          class="w-full py-3.5 bg-brand-600 hover:bg-brand-700 active:bg-brand-800 disabled:opacity-50 text-white rounded-xl font-bold text-sm shadow-md shadow-brand-600/25 transition-all mt-2"
        >
          <span v-if="auth.isLoading">로그인 중...</span>
          <span v-else>로그인</span>
        </button>
      </form>

      <div class="mt-6 text-center text-xs text-slate-500">
        아직 계정이 없으신가요?
        <RouterLink to="/register" class="font-bold text-brand-600 hover:underline ml-1">
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

const email = ref('');
const password = ref('');
const error = ref('');

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

