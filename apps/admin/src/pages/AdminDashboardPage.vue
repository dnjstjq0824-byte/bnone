<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-black text-white tracking-tight">마스터 운영 대시보드</h1>
        <p class="text-xs text-slate-400 mt-0.5">BN ONE v1.0 글로벌 서비스 상태 및 운영 지표를 모니터링합니다.</p>
      </div>
      <div class="flex items-center gap-2">
        <button
          @click="loadData"
          class="px-3.5 py-2 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded-xl text-xs font-bold transition-all"
        >
          새로고침
        </button>
      </div>
    </div>

    <!-- Mandatory Privacy Banner (Article 21.3) -->
    <div class="p-4 rounded-2xl bg-slate-800/80 border border-slate-700 flex items-start gap-3">
      <span class="text-blue-400 text-lg leading-none">🛡️</span>
      <p class="text-xs text-slate-300 leading-relaxed">
        <strong>개인정보 보호 정책 준수:</strong> 일반 관리자 콘솔 대시보드에서는 사용자의 개별 금융 잔액, 상세 거래내역 및 민감 개인정보를 표시하지 않으며, 집계 및 익명화된 운영 지표만을 표출합니다.
      </p>
    </div>

    <!-- Operational KPI Grid -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4">
      <div class="p-5 rounded-2xl bg-[#111827] border border-slate-800">
        <span class="text-xs font-bold text-slate-400">가입 사용자</span>
        <div class="text-3xl font-black text-white mt-1">{{ dashboard?.kpi?.total_registered_users || 0 }}</div>
        <span class="text-[10px] text-emerald-400 font-semibold mt-1 block">정상 등록 계정</span>
      </div>

      <div class="p-5 rounded-2xl bg-[#111827] border border-slate-800">
        <span class="text-xs font-bold text-slate-400">활성 세션</span>
        <div class="text-3xl font-black text-blue-400 mt-1">{{ dashboard?.kpi?.active_user_sessions || 0 }}</div>
        <span class="text-[10px] text-slate-400 font-semibold mt-1 block">실시간 토큰 유효 세션</span>
      </div>

      <div class="p-5 rounded-2xl bg-[#111827] border border-slate-800">
        <span class="text-xs font-bold text-slate-400">미처리 고객지원</span>
        <div class="text-3xl font-black text-amber-400 mt-1">{{ dashboard?.kpi?.open_support_tickets || 0 }}</div>
        <span class="text-[10px] text-slate-400 font-semibold mt-1 block">답변 대기 중 티켓</span>
      </div>

      <div class="p-5 rounded-2xl bg-[#111827] border border-slate-800">
        <span class="text-xs font-bold text-slate-400">데이터 삭제 요청</span>
        <div class="text-3xl font-black text-red-400 mt-1">{{ dashboard?.kpi?.pending_deletion_reqs || 0 }}</div>
        <span class="text-[10px] text-slate-400 font-semibold mt-1 block">GDPR/PIPA 심사 대기</span>
      </div>

      <div class="p-5 rounded-2xl bg-[#111827] border border-slate-800">
        <span class="text-xs font-bold text-slate-400">활성 FX Provider</span>
        <div class="text-3xl font-black text-emerald-400 mt-1">{{ dashboard?.kpi?.active_fx_providers || 0 }}</div>
        <span class="text-[10px] text-slate-400 font-semibold mt-1 block">한국은행 / ECB 등</span>
      </div>
    </div>

    <!-- Quick Operations Tools -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="p-6 rounded-3xl bg-[#111827] border border-slate-800 space-y-3">
        <h3 class="text-sm font-bold text-white flex items-center gap-2">
          <span>💱</span>
          <span>FX 참고 환율 즉시 동기화</span>
        </h3>
        <p class="text-xs text-slate-400 leading-relaxed">
          공인 기관(BOK, ECB)의 최신 기준 환율을 즉시 수집하고 유효성을 재검증합니다.
        </p>
        <button @click="triggerFXSync" class="w-full py-2.5 bg-blue-600 hover:bg-blue-500 text-white rounded-xl text-xs font-bold transition-all">
          환율 수집 배치 실행
        </button>
      </div>

      <div class="p-6 rounded-3xl bg-[#111827] border border-slate-800 space-y-3">
        <h3 class="text-sm font-bold text-white flex items-center gap-2">
          <span>👥</span>
          <span>사용자 계정 상태 관리</span>
        </h3>
        <p class="text-xs text-slate-400 leading-relaxed">
          이상 활동 사용자 계정을 검토하거나 긴급 세션 철회 조치를 수행합니다.
        </p>
        <RouterLink to="/admin/users" class="block w-full py-2.5 bg-slate-800 hover:bg-slate-700 text-center text-white rounded-xl text-xs font-bold transition-all">
          사용자 목록 조회
        </RouterLink>
      </div>

      <div class="p-6 rounded-3xl bg-[#111827] border border-slate-800 space-y-3">
        <h3 class="text-sm font-bold text-white flex items-center gap-2">
          <span>📜</span>
          <span>보안 및 감사 로그</span>
        </h3>
        <p class="text-xs text-slate-400 leading-relaxed">
          운영자 권한 변경, 세션 철회, 데이터 삭제 요청 등 모든 감사 이력을 확인합니다.
        </p>
        <RouterLink to="/admin/audit-logs" class="block w-full py-2.5 bg-slate-800 hover:bg-slate-700 text-center text-white rounded-xl text-xs font-bold transition-all">
          감사 로그 뷰어
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { RouterLink } from 'vue-router';
import { adminFetch } from '../api/client';

const dashboard = ref<any>(null);

async function loadData() {
  dashboard.value = await adminFetch('/admin/dashboard');
}

onMounted(() => {
  loadData();
});

async function triggerFXSync() {
  const res = await adminFetch('/admin/fx/sync-jobs', { method: 'POST' });
  alert(res.message);
  await loadData();
}
</script>

