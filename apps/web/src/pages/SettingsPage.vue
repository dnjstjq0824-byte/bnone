<template>
  <div class="max-w-3xl space-y-6">
    <div>
      <h1 class="text-2xl font-black text-slate-900 tracking-tight">설정 & 개인정보 보호</h1>
      <p class="text-xs text-slate-500 mt-0.5">계정 환경 설정, 보안 세션, 데이터 내보내기 및 삭제를 관리합니다.</p>
    </div>

    <!-- User Profile Settings -->
    <div class="bg-white p-6 rounded-3xl border border-slate-200 card-shadow space-y-4">
      <h2 class="text-base font-bold text-slate-900">사용자 프로필</h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1">이메일</label>
          <input :value="auth.user?.email" disabled class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-xs bg-slate-50 text-slate-500" />
        </div>
        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1">이름 / 닉네임</label>
          <input v-model="profileForm.display_name" class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-xs focus:outline-none focus:border-brand-600" />
        </div>
        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1">기준 통화</label>
          <select v-model="profileForm.base_currency" class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-xs bg-white">
            <option value="KRW">KRW — 대한민국 원 (₩)</option>
            <option value="USD">USD — 미국 달러 ($)</option>
            <option value="EUR">EUR — 유로 (€)</option>
            <option value="JPY">JPY — 일본 엔 (¥)</option>
          </select>
        </div>
        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1">시간대</label>
          <input v-model="profileForm.timezone" class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-xs focus:outline-none focus:border-brand-600" />
        </div>
      </div>
      <div class="text-right pt-2">
        <button @click="saveProfile" class="px-4 py-2 bg-brand-600 hover:bg-brand-700 text-white rounded-xl text-xs font-bold shadow-sm transition-all">
          프로필 변경 저장
        </button>
      </div>
    </div>

    <!-- Active Sessions -->
    <div class="bg-white p-6 rounded-3xl border border-slate-200 card-shadow space-y-4">
      <h2 class="text-base font-bold text-slate-900">로그인 세션 관리</h2>
      <div class="divide-y divide-slate-100">
        <div v-for="s in sessions" :key="s.id" class="py-3 flex items-center justify-between text-xs">
          <div>
            <div class="font-bold text-slate-800">{{ s.platform }} — {{ s.device_name || '웹 브라우저' }}</div>
            <div class="text-[10px] text-slate-400">생성일: {{ s.created_at }}</div>
          </div>
          <button @click="revokeSession(s.id)" class="text-red-500 hover:text-red-700 font-bold text-[11px]">
            세션 종료
          </button>
        </div>
      </div>
    </div>

    <!-- Data Export Section -->
    <div class="bg-white p-6 rounded-3xl border border-slate-200 card-shadow space-y-3">
      <h2 class="text-base font-bold text-slate-900">데이터 내보내기 (GDPR / PIPA)</h2>
      <p class="text-xs text-slate-500 leading-relaxed">
        사용자의 모든 거래 내역, 계정, 자산 및 예산 데이터를 CSV 포맷으로 안전하게 내보낼 수 있습니다.
      </p>
      <button @click="requestExport" class="px-4 py-2.5 bg-slate-800 hover:bg-slate-700 text-white rounded-xl text-xs font-bold transition-all">
        전체 데이터 내보내기 요청 (CSV)
      </button>
    </div>

    <!-- Danger Zone: Account Deletion -->
    <div class="bg-red-50/50 p-6 rounded-3xl border border-red-200 card-shadow space-y-3">
      <h2 class="text-base font-bold text-red-900">회원 탈퇴 및 데이터 삭제</h2>
      <p class="text-xs text-red-700 leading-relaxed">
        계정 삭제 시 등록된 모든 금융 데이터가 안전하게 파기 또는 비식별화 처리됩니다. 법정 의무 보존 기간이 존재하는 경우 별도 분리 보관됩니다.
      </p>
      <button @click="requestAccountDeletion" class="px-4 py-2.5 bg-red-600 hover:bg-red-700 text-white rounded-xl text-xs font-bold transition-all">
        계정 영구 삭제 요청
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { apiFetch } from '../api/client';
import { useAuthStore } from '../stores/auth';

const router = useRouter();
const auth = useAuthStore();

const sessions = ref<any[]>([]);
const profileForm = ref({
  display_name: auth.user?.display_name || '',
  base_currency: auth.user?.base_currency || 'KRW',
  timezone: auth.user?.timezone || 'Asia/Seoul',
});

async function loadData() {
  sessions.value = await apiFetch('/me/sessions');
}

onMounted(() => {
  loadData();
});

async function saveProfile() {
  await apiFetch('/me', {
    method: 'PATCH',
    body: JSON.stringify(profileForm.value),
  });
  await auth.fetchMe();
  alert('프로필 설정이 저장되었습니다.');
}

async function revokeSession(id: string) {
  if (confirm('이 세션을 종료하시겠습니까?')) {
    await apiFetch(`/me/sessions/${id}`, { method: 'DELETE' });
    await loadData();
  }
}

async function requestExport() {
  await apiFetch('/me/data-exports', {
    method: 'POST',
    body: JSON.stringify({ format: 'CSV' }),
  });
  alert('데이터 내보내기가 완료되었습니다.');
}

async function requestAccountDeletion() {
  if (confirm('정말로 계정을 영구 삭제하시겠습니까? 모든 개인 금융 데이터가 영구 파기됩니다.')) {
    await apiFetch('/me/deletion-requests', { method: 'POST' });
    alert('계정 삭제 요청이 접수되었습니다.');
    auth.logout();
    router.push('/login');
  }
}
</script>

