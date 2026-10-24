<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-black text-slate-900 tracking-tight">사용자 알림함</h1>
        <p class="text-xs text-slate-500 mt-0.5">구독 갱신 예정, 예산 초과 경고 및 시스템 공지사항을 확인합니다.</p>
      </div>
    </div>

    <div v-if="notifications.length === 0" class="bg-white p-12 rounded-3xl border border-slate-200 card-shadow text-center">
      <div class="text-3xl mb-3">🔔</div>
      <div class="text-sm font-bold text-slate-700">새로운 알림이 없습니다.</div>
      <p class="text-xs text-slate-400 mt-1">구독 갱신이나 예산 알림이 발생하면 이곳에 표시됩니다.</p>
    </div>

    <div v-else class="bg-white rounded-3xl border border-slate-200 card-shadow divide-y divide-slate-100 overflow-hidden">
      <div
        v-for="n in notifications"
        :key="n.id"
        class="p-5 flex items-start justify-between gap-4 transition-colors"
        :class="n.status === 'READ' ? 'bg-white opacity-70' : 'bg-blue-50/20'"
      >
        <div class="flex items-start gap-3.5">
          <div class="w-9 h-9 rounded-xl bg-slate-100 flex items-center justify-center text-sm shrink-0">
            {{ n.type === 'BUDGET_ALERT' ? '🎯' : n.type === 'RENEWAL_REMINDER' ? '🔄' : '📢' }}
          </div>
          <div>
            <div class="flex items-center gap-2">
              <h3 class="font-bold text-xs text-slate-900">{{ n.title }}</h3>
              <span v-if="n.status !== 'READ'" class="w-1.5 h-1.5 rounded-full bg-brand-600"></span>
            </div>
            <p class="text-xs text-slate-600 mt-1 leading-relaxed">{{ n.body }}</p>
            <span class="text-[10px] text-slate-400 mt-2 block">{{ formatDate(n.created_at) }}</span>
          </div>
        </div>

        <button
          v-if="n.status !== 'READ'"
          @click="markRead(n.id)"
          class="px-2.5 py-1 text-[11px] font-bold text-slate-500 hover:text-brand-600 hover:bg-slate-100 rounded-lg shrink-0 transition-colors"
        >
          읽음 표시
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { apiFetch } from '../api/client';

const notifications = ref<any[]>([]);

async function loadNotifications() {
  notifications.value = await apiFetch('/me/notifications');
}

onMounted(() => {
  loadNotifications();
});

async function markRead(id: string) {
  await apiFetch(`/me/notifications/${id}/read`, { method: 'PATCH' });
  await loadNotifications();
}

function formatDate(iso: string) {
  if (!iso) return '';
  const d = new Date(iso);
  return `${d.getMonth() + 1}월 ${d.getDate()}일 ${d.getHours()}:${String(d.getMinutes()).padStart(2, '0')}`;
}
</script>

