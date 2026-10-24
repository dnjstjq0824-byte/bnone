<template>
  <div class="min-h-screen bg-[#0B0F19] flex items-center justify-center p-4">
    <div class="max-w-md w-full bg-[#111827] rounded-3xl p-8 border border-slate-800 shadow-2xl">
      <div class="text-center mb-8">
        <div class="w-12 h-12 rounded-2xl bg-blue-600 mx-auto flex items-center justify-center text-white font-black text-2xl shadow-lg shadow-blue-500/30 mb-3">
          BN
        </div>
        <h2 class="text-2xl font-bold text-white tracking-tight">BN ONE 콘솔 로그인</h2>
        <p class="text-xs text-slate-400 mt-1">Master Operations & Audited Admin</p>
      </div>

      <div v-if="error" class="mb-5 p-3 rounded-xl bg-red-500/10 border border-red-500/20 text-red-400 text-xs font-medium">
        {{ error }}
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label class="block text-xs font-bold text-slate-300 mb-1.5">운영자 이메일</label>
          <input
            v-model="email"
            type="email"
            required
            placeholder="admin@bookingnetwork.com"
            class="w-full px-4 py-3 rounded-xl bg-slate-900 border border-slate-700 text-sm text-white focus:outline-none focus:border-blue-500 transition-all"
          />
        </div>

        <div>
          <label class="block text-xs font-bold text-slate-300 mb-1.5">보안 패스워드</label>
          <input
            v-model="password"
            type="password"
            required
            placeholder="••••••••••••"
            class="w-full px-4 py-3 rounded-xl bg-slate-900 border border-slate-700 text-sm text-white focus:outline-none focus:border-blue-500 transition-all"
          />
        </div>

        <button
          type="submit"
          :disabled="adminAuth.isLoading"
          class="w-full py-3.5 bg-blue-600 hover:bg-blue-500 active:bg-blue-700 disabled:opacity-50 text-white rounded-xl font-bold text-sm shadow-lg shadow-blue-600/30 transition-all mt-4"
        >
          <span v-if="adminAuth.isLoading">인증 검증 중...</span>
          <span v-else>관리자 콘솔 접속</span>
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAdminAuthStore } from '../stores/adminAuth';

const router = useRouter();
const adminAuth = useAdminAuthStore();

const email = ref('admin@bookingnetwork.com');
const password = ref('AdminSuperPass2026!');
const error = ref('');

async function handleSubmit() {
  error.value = '';
  try {
    await adminAuth.login(email.value, password.value);
    router.push('/admin/dashboard');
  } catch (err: any) {
    error.value = err.message || '인증에 실패했습니다. 권한을 확인하세요.';
  }
}
</script>
