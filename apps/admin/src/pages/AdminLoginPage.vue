<template>
  <div class="min-h-screen bg-[#08170F] flex items-center justify-center p-4">
    <div class="max-w-md w-full bg-[#0C2318] rounded-3xl p-8 border border-white/10 shadow-2xl">
      <div class="text-center mb-8">
        <div class="w-12 h-12 rounded-2xl bg-[#C0A04E] mx-auto flex items-center justify-center text-[#08170F] font-black text-2xl mb-3">
          모
        </div>
        <h2 class="text-2xl font-bold text-white tracking-tight">MOA 콘솔 로그인</h2>
        <p class="text-xs text-[#7E9387] mt-1">Master Operations & Audited Admin</p>
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
            class="w-full px-4 py-3 rounded-xl bg-slate-900 border border-slate-700 text-sm text-white focus:outline-none focus:border-[#2A7248] transition-all"
          />
        </div>

        <div>
          <label class="block text-xs font-bold text-slate-300 mb-1.5">보안 패스워드</label>
          <input
            v-model="password"
            type="password"
            required
            placeholder="••••••••••••"
            class="w-full px-4 py-3 rounded-xl bg-slate-900 border border-slate-700 text-sm text-white focus:outline-none focus:border-[#2A7248] transition-all"
          />
        </div>

        <button
          type="submit"
          :disabled="adminAuth.isLoading"
          class="w-full py-3.5 bg-[#1E5A38] hover:bg-[#2A7248] active:bg-[#17452C] disabled:opacity-50 text-white rounded-xl font-bold text-sm shadow-lg shadow-[#1E5A38]/30 transition-all mt-4"
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

const email = ref('');
const password = ref('');
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

