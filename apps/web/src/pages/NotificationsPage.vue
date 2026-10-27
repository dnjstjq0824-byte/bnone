<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-[#191F28]">알림함</h1>
        <p class="text-sm text-[#8B95A1] mt-1">구독 갱신 예정, 예산 초과 경고 및 시스템 공지사항을 확인해요.</p>
      </div>
    </div>

    <div v-if="notifications.length === 0" class="bg-white p-12 rounded-3xl text-center">
      <div class="w-14 h-14 mx-auto mb-4 rounded-2xl bg-[#F2F4F6] flex items-center justify-center text-[#8B95A1] text-2xl">
        🔔
      </div>
      <div class="text-base font-bold text-[#191F28]">새로운 알림이 없습니다.</div>
      <p class="text-xs text-[#8B95A1] mt-1">구독 갱신이나 예산 알림이 발생하면 이곳에 표시돼요.</p>
    </div>

    <div v-else class="bg-white rounded-3xl divide-y divide-[#F2F4F6] overflow-hidden p-2">
      <div
        v-for="n in notifications"
        :key="n.id"
        class="p-4 flex items-start justify-between gap-4 rounded-2xl transition-colors"
        :class="n.status === 'READ' ? 'bg-white opacity-70' : 'bg-[#F9FAFB]'"
      >
        <div class="flex items-start gap-3.5">
          <div class="w-10 h-10 rounded-2xl bg-[#E8F3FF] flex items-center justify-center text-base shrink-0">
            {{ n.type === 'BUDGET_ALERT' ? '🎯' : n.type === 'RENEWAL_REMINDER' ? '🔄' : '📢' }}
          </div>
          <div>
            <div class="flex items-center gap-2">
              <h3 class="font-bold text-sm text-[#191F28]">{{ n.title }}</h3>
              <span v-if="n.status !== 'READ'" class="w-1.5 h-1.5 rounded-full bg-[#3182F6]"></span>
            </div>
            <p class="text-xs text-[#4E5968] mt-1 leading-relaxed">{{ n.body }}</p>
            <span class="text-xs text-[#8B95A1] mt-1.5 block">{{ formatDate(n.created_at) }}</span>
          </div>
        </div>

        <button
          v-if="n.status !== 'READ'"
          @click="markRead(n.id)"
          class="px-3 py-1.5 text-xs font-bold text-[#3182F6] hover:bg-[#E8F3FF] rounded-xl shrink-0 transition-colors"
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
