<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-black text-white tracking-tight">개인정보 보호 & 삭제 요청 관리</h1>
      <p class="text-xs text-slate-400 mt-0.5">GDPR / 개인정보보호법에 따른 데이터 삭제 및 영구 파기 요청을 심사합니다.</p>
    </div>

    <div class="bg-[#111827] rounded-3xl border border-slate-800 overflow-hidden">
      <div v-if="requests.length === 0" class="p-12 text-center text-xs text-slate-500">
        접수된 계정 삭제 요청이 없습니다.
      </div>
      <table v-else class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-slate-800/50 text-[11px] font-bold text-slate-400 border-b border-slate-800">
            <th class="py-3 px-5">사용자 이메일</th>
            <th class="py-3 px-4">요청 시각</th>
            <th class="py-3 px-4">상태</th>
            <th class="py-3 px-4">처리 시각</th>
            <th class="py-3 px-4 text-center">심사 및 실행</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/60 text-xs">
          <tr v-for="r in requests" :key="r.id" class="hover:bg-slate-800/30">
            <td class="py-3.5 px-5 font-bold text-white">{{ r.email }}</td>
            <td class="py-3.5 px-4 text-slate-400">{{ formatDate(r.requested_at) }}</td>
            <td class="py-3.5 px-4">
              <span
                class="px-2 py-0.5 rounded-full text-[10px] font-bold"
                :class="r.status === 'COMPLETED' ? 'bg-emerald-500/10 text-emerald-400' : r.status === 'PENDING' ? 'bg-amber-500/10 text-amber-400' : 'bg-slate-800 text-slate-400'"
              >
                {{ r.status }}
              </span>
            </td>
            <td class="py-3.5 px-4 text-slate-500">{{ formatDate(r.processed_at) || '-' }}</td>
            <td class="py-3.5 px-4 text-center space-x-2">
              <button
                v-if="r.status === 'PENDING'"
                @click="reviewRequest(r.id, 'APPROVE_AND_EXECUTE')"
                class="px-2.5 py-1 bg-red-600 hover:bg-red-500 text-white rounded text-[10px] font-bold"
              >
                삭제 승인 및 영구 비식별화
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { adminFetch } from '../api/client';

const requests = ref<any[]>([]);

async function loadData() {
  requests.value = await adminFetch('/admin/privacy/deletion-requests');
}

onMounted(() => {
  loadData();
});

async function reviewRequest(id: string, action: string) {
  if (confirm('해당 사용자의 모든 개인정보를 영구 파기 및 비식별화 처리하시겠습니까? 이 작업은 되돌릴 수 없습니다.')) {
    await adminFetch(`/admin/privacy/deletion-requests/${id}/review`, {
      method: 'POST',
      body: JSON.stringify({ action }),
    });
    alert('삭제 및 비식별화 처리가 완료되었습니다.');
    await loadData();
  }
}

function formatDate(iso: string) {
  if (!iso) return '';
  return new Date(iso).toISOString().slice(0, 16).replace('T', ' ');
}
</script>

