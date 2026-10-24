<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-black text-white tracking-tight">고객지원 티켓 관리</h1>
        <p class="text-xs text-slate-400 mt-0.5">사용자 문의 및 오류 접수 건에 대해 프라이버시 보호 원칙에 따라 답변을 처리합니다.</p>
      </div>
    </div>

    <!-- Privacy Banner -->
    <div class="p-3.5 rounded-2xl bg-blue-900/30 border border-blue-800 text-blue-200 text-xs flex items-center gap-2.5">
      <span>🔒</span>
      <p class="text-[11px] leading-relaxed">
        <strong>원칙 준수 (Article 21.5):</strong> 지원 담당자는 문의 내용 및 상태만을 열람하며, 사용자의 전체 자산 및 상세 거래내역 원문은 기본 노출되지 않습니다.
      </p>
    </div>

    <!-- Tickets List -->
    <div class="bg-[#111827] rounded-3xl border border-slate-800 overflow-hidden">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-slate-800/50 text-[11px] font-bold text-slate-400 border-b border-slate-800">
            <th class="py-3 px-5">문의자</th>
            <th class="py-3 px-4">유형</th>
            <th class="py-3 px-4">제목</th>
            <th class="py-3 px-4">상태</th>
            <th class="py-3 px-4">접수시각</th>
            <th class="py-3 px-4 text-center">관리</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/60 text-xs">
          <tr v-for="t in tickets" :key="t.id" class="hover:bg-slate-800/30">
            <td class="py-3.5 px-5 font-bold text-white">{{ t.user_email }}</td>
            <td class="py-3.5 px-4 text-slate-400">{{ t.category }}</td>
            <td class="py-3.5 px-4 font-bold text-slate-200">{{ t.subject }}</td>
            <td class="py-3.5 px-4">
              <span
                class="px-2 py-0.5 rounded-full text-[10px] font-bold"
                :class="t.status === 'RESOLVED' ? 'bg-emerald-500/10 text-emerald-400' : 'bg-amber-500/10 text-amber-400'"
              >
                {{ t.status }}
              </span>
            </td>
            <td class="py-3.5 px-4 text-slate-500">{{ formatDate(t.created_at) }}</td>
            <td class="py-3.5 px-4 text-center space-x-2">
              <button @click="openReply(t)" class="text-[11px] text-blue-400 hover:text-blue-300 font-bold">
                답변/열람
              </button>
              <button v-if="t.status !== 'RESOLVED'" @click="resolveTicket(t.id)" class="text-[11px] text-emerald-400 hover:text-emerald-300 font-bold">
                해결 처리
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Reply Modal -->
    <div v-if="activeTicket" class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="bg-[#111827] max-w-lg w-full rounded-3xl p-6 border border-slate-800 space-y-4 shadow-2xl">
        <div class="flex items-center justify-between border-b border-slate-800 pb-3">
          <div>
            <span class="text-[10px] font-bold text-blue-400">티켓 #{{ activeTicket.id.slice(0, 8) }}</span>
            <h3 class="text-sm font-bold text-white">{{ activeTicket.subject }}</h3>
          </div>
          <button @click="activeTicket = null" class="text-slate-400 hover:text-white">✕</button>
        </div>

        <div class="p-3 rounded-xl bg-slate-900 border border-slate-800 text-xs text-slate-300 space-y-1">
          <div class="font-bold text-slate-400">사용자 문의 내용:</div>
          <p>{{ activeTicket.description }}</p>
        </div>

        <!-- Thread -->
        <div class="space-y-2 max-h-48 overflow-y-auto pr-1">
          <div
            v-for="msg in activeTicket.messages || []"
            :key="msg.id"
            class="p-2.5 rounded-xl text-xs"
            :class="msg.sender_type === 'ADMIN' ? 'bg-blue-900/30 text-blue-200 ml-4 border border-blue-800' : 'bg-slate-800 text-slate-200 mr-4'"
          >
            <div class="text-[9px] text-slate-400 font-bold mb-1">
              {{ msg.sender_type === 'ADMIN' ? '운영 담당자' : '사용자' }}
            </div>
            <div>{{ msg.message }}</div>
          </div>
        </div>

        <!-- Reply input -->
        <div class="space-y-2 pt-2 border-t border-slate-800">
          <textarea
            v-model="replyText"
            rows="3"
            placeholder="사용자에게 전달할 답변을 작성하세요."
            class="w-full px-3.5 py-2.5 rounded-xl bg-slate-900 border border-slate-700 text-xs text-white focus:outline-none focus:border-blue-500"
          ></textarea>
          <div class="flex justify-end gap-2">
            <button @click="activeTicket = null" class="px-3 py-2 bg-slate-800 text-slate-300 text-xs rounded-xl font-bold">닫기</button>
            <button @click="submitReply" class="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white text-xs rounded-xl font-bold">답변 발송</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { adminFetch } from '../api/client';

const tickets = ref<any[]>([]);
const activeTicket = ref<any>(null);
const replyText = ref('');

async function loadTickets() {
  tickets.value = await adminFetch('/admin/support-tickets');
}

onMounted(() => {
  loadTickets();
});

async function openReply(t: any) {
  const detail = await adminFetch(`/admin/support-tickets/${t.id}`);
  activeTicket.value = detail;
  replyText.value = '';
}

async function submitReply() {
  if (!replyText.value) return;
  await adminFetch(`/admin/support-tickets/${activeTicket.value.id}/reply`, {
    method: 'POST',
    body: JSON.stringify({ message: replyText.value, visibility: 'PUBLIC' }),
  });
  activeTicket.value = null;
  await loadTickets();
}

async function resolveTicket(id: string) {
  if (confirm('이 문의 건을 해결 완료 처리하시겠습니까?')) {
    await adminFetch(`/admin/support-tickets/${id}/resolve`, { method: 'POST' });
    await loadTickets();
  }
}

function formatDate(iso: string) {
  if (!iso) return '';
  return new Date(iso).toISOString().slice(0, 16).replace('T', ' ');
}
</script>

