<template>
  <div class="max-w-3xl space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold tracking-tight text-[#191F28]">설정 & 개인정보 보호</h1>
      <p class="text-sm text-[#8B95A1] mt-1">계정 환경 설정, 보안 세션, 데이터 내보내기 및 삭제를 관리해요.</p>
    </div>

    <!-- User Profile Settings Card -->
    <div class="bg-white p-6 rounded-3xl space-y-5">
      <h2 class="text-base font-bold text-[#191F28]">사용자 프로필</h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label class="block text-xs font-bold text-[#4E5968] mb-1.5">이메일</label>
          <input
            :value="auth.user?.email"
            disabled
            class="w-full px-4 py-3 rounded-2xl bg-[#F2F4F6] border border-transparent text-sm text-[#8B95A1]"
          />
        </div>
        <div>
          <label class="block text-xs font-bold text-[#4E5968] mb-1.5">이름 / 닉네임</label>
          <input
            v-model="profileForm.display_name"
            class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] focus:outline-none focus:border-[#3182F6]"
          />
        </div>
        <div>
          <label class="block text-xs font-bold text-[#4E5968] mb-1.5">기준 통화</label>
          <select
            v-model="profileForm.base_currency"
            class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] focus:outline-none focus:border-[#3182F6]"
          >
            <option value="KRW">KRW — 대한민국 원 (₩)</option>
            <option value="USD">USD — 미국 달러 ($)</option>
            <option value="EUR">EUR — 유로 (€)</option>
            <option value="JPY">JPY — 일본 엔 (¥)</option>
          </select>
        </div>
        <div>
          <label class="block text-xs font-bold text-[#4E5968] mb-1.5">시간대</label>
          <input
            v-model="profileForm.timezone"
            class="w-full px-4 py-3 rounded-2xl bg-[#F9FAFB] border border-[#E5E8EB] text-sm text-[#191F28] focus:outline-none focus:border-[#3182F6]"
          />
        </div>
      </div>
      <div class="text-right pt-2">
        <button
          @click="saveProfile"
          class="px-5 py-2.5 bg-[#3182F6] hover:bg-[#1B64DA] text-white rounded-xl text-xs font-bold shadow-sm transition-colors"
        >
          프로필 변경 저장
        </button>
      </div>
    </div>

    <!-- Active Sessions Card -->
    <div class="bg-white p-6 rounded-3xl space-y-4">
      <h2 class="text-base font-bold text-[#191F28]">로그인 세션 관리</h2>
      <div class="divide-y divide-[#F2F4F6]">
        <div v-for="s in sessions" :key="s.id" class="py-3.5 flex items-center justify-between text-xs">
          <div>
            <div class="font-bold text-[#191F28] text-sm">{{ s.platform }} — {{ s.device_name || '웹 브라우저' }}</div>
            <div class="text-[#8B95A1] text-xs mt-0.5">최근 접속일: {{ s.created_at }}</div>
          </div>
          <button
            @click="revokeSession(s.id)"
            class="text-[#F04452] hover:underline font-bold text-xs"
          >
            세션 종료
          </button>
        </div>
      </div>
    </div>

    <!-- Data Export Card -->
    <div class="bg-white p-6 rounded-3xl space-y-3">
      <h2 class="text-base font-bold text-[#191F28]">데이터 내보내기 (GDPR / 개인정보보호법)</h2>
      <p class="text-xs text-[#8B95A1] leading-relaxed">
        사용자의 모든 거래 내역, 계정, 자산 및 예산 데이터를 CSV 포맷으로 안전하게 내보낼 수 있습니다.
      </p>
      <button
        @click="requestExport"
        class="px-5 py-2.5 bg-[#F2F4F6] hover:bg-[#E5E8EB] text-[#4E5968] rounded-xl text-xs font-bold transition-colors"
      >
        전체 데이터 내보내기 (CSV)
      </button>
    </div>

    <!-- Danger Zone Card -->
    <div class="bg-[#FEECEB] p-6 rounded-3xl space-y-3">
      <h2 class="text-base font-bold text-[#F04452]">회원 탈퇴 및 데이터 영구 파기</h2>
      <p class="text-xs text-[#F04452] leading-relaxed">
        계정 삭제 시 등록된 모든 금융 데이터가 안전하게 즉시 파기 처리됩니다.
      </p>
      <button
        @click="requestAccountDeletion"
        class="px-5 py-2.5 bg-[#F04452] hover:bg-[#dc2626] text-white rounded-xl text-xs font-bold transition-colors"
      >
        계정 영구 삭제하기
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
