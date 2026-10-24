<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-black text-slate-900 tracking-tight">1:1 고객지원</h1>
        <p class="text-xs text-slate-500 mt-0.5">이용 문의 및 오류 보고를 제출하고 처리 상태를 확인합니다.</p>
      </div>
      <button
        @click="showCreateModal = true"
        class="px-4 py-2 bg-brand-600 hover:bg-brand-700 text-white rounded-xl text-xs font-bold shadow-md shadow-brand-600/20 transition-all"
      >
        + 문의하기
      </button>
    </div>

    <!-- Mandatory Privacy Notice (Article 20) -->
    <div class="p-3.5 rounded-2xl bg-slate-100 border border-slate-200 text-slate-600 text-xs flex items-start gap-2.5">
      <span>🔒</span>
      <p class="text-[11px] leading-relaxed">
        <strong>개인정보 보호 안내:</strong> 원활한 지원 처리를 위해 계좌번호, 비밀번호 등 민감한 금융정보가 문의 내용에 불필요하게 포함되지 않도록 유의해 주시기 바랍니다.
      </p>
    </div>

    <!-- Ticket List -->
    <div v-if="tickets.length === 0" class="bg-white p-12 rounded-3xl border border-slate-200 card-shadow text-center">
      <div class="text-3xl mb-3">💬</div>
      <div class="text-sm font-bold text-slate-700">접수된 문의 내역이 없습니다.</div>
      <p class="text-xs text-slate-400 mt-1 mb-4">서비스 이용 중 궁금한 점이 있으시면 문의를 남겨주세요.</p>
      <button @click="showCreateModal = true" class="px-4 py-2 bg-brand-600 text-white rounded-xl text-xs font-bold">
        새 문의 작성하기
      </button>
    </div>

    <div v-else class="space-y-4">
      <div
        v-for="t in tickets"
        :key="t.id"
        class="bg-white p-6 rounded-3xl border border-slate-200 card-shadow space-y-3 cursor-pointer hover:border-brand-500 transition-all"
        @click="openTicket(t)"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-slate-100 text-slate-600">{{ t.category }}</span>
            <h3 class="font-bold text-sm text-slate-900">{{ t.subject }}</h3>
          </div>
          <span
            class="px-2.5 py-0.5 rounded-full text-[10px] font-bold"
            :class="t.status === 'RESOLVED' ? 'bg-emerald-50 text-emerald-600' : 'bg-blue-50 text-blue-600'"
          >
            {{ t.status }}
          </span>
        </div>
        <p class="text-xs text-slate-500 line-clamp-2">{{ t.description }}</p>
        <div class="text-[10px] text-slate-400 pt-2 border-t border-slate-100 flex items-center justify-between">
          <span>접수일: {{ t.created_at }}</span>
          <span class="text-brand-600 font-bold">상세 및 답변 보기 →</span>
        </div>
      </div>
    </div>

    <!-- Create Ticket Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="bg-white max-w-md w-full rounded-3xl p-6 card-shadow space-y-4 border border-slate-200">
        <h3 class="text-lg font-bold text-slate-900">새 고객 문의 작성</h3>
        <form @submit.prevent="submitTicket" class="space-y-3">
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">문의 유형</label>
            <select v-model="newTicket.category" class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm bg-white">
              <option value="GENERAL">일반 이용 문의</option>
              <option value="ACCOUNT">계정 및 인증</option>
              <option value="BUG">기능 오류 보고</option>
              <option value="SUGGESTION">기능 제안</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">제목</label>
            <input v-model="newTicket.subject" required placeholder="문의 제목을 입력하세요" class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:border-brand-600" />
          </div>
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">상세 내용</label>
            <textarea v-model="newTicket.description" required rows="4" placeholder="문의하실 내용을 구체적으로 작성해 주세요." class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm"></textarea>
          </div>
          <div class="flex gap-2 pt-2">
            <button type="button" @click="showCreateModal = false" class="w-1/2 py-2.5 bg-slate-100 hover:bg-slate-200 rounded-xl text-xs font-bold text-slate-700">취소</button>
            <button type="submit" class="w-1/2 py-2.5 bg-brand-600 hover:bg-brand-700 text-white rounded-xl text-xs font-bold shadow-md shadow-brand-600/25">문의 접수</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { apiFetch } from '../api/client';

const tickets = ref<any[]>([]);
const showCreateModal = ref(false);

const newTicket = ref({
  category: 'GENERAL',
  subject: '',
  description: '',
  priority: 'NORMAL',
});

async function loadTickets() {
  tickets.value = await apiFetch('/me/support-tickets');
}

onMounted(() => {
  loadTickets();
});

async function submitTicket() {
  await apiFetch('/me/support-tickets', {
    method: 'POST',
    body: JSON.stringify(newTicket.value),
  });
  showCreateModal.value = false;
  newTicket.value = { category: 'GENERAL', subject: '', description: '', priority: 'NORMAL' };
  await loadTickets();
}

async function openTicket(ticket: any) {
  const detail = await apiFetch(`/me/support-tickets/${ticket.id}`);
  let msgList = detail.messages?.map((m: any) => `[${m.sender_type}] ${m.message}`).join('\n\n') || '답변 대기 중입니다.';
  alert(`[${detail.subject}]\n상태: ${detail.status}\n\n내용:\n${detail.description}\n\n대화 내역:\n${msgList}`);
}
</script>

