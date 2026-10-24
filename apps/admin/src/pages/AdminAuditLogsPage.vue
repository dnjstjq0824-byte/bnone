<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-black text-white tracking-tight">불변 감사 로그 (Immutable Audit Trail)</h1>
      <p class="text-xs text-slate-400 mt-0.5">운영자 활동, 권한 변경, 세션 철회 및 프라이버시 처리 이력을 엄격하게 추적합니다.</p>
    </div>

    <div class="bg-[#111827] rounded-3xl border border-slate-800 overflow-hidden">
      <div v-if="logs.length === 0" class="p-12 text-center text-xs text-slate-500">
        기록된 감사 로그가 없습니다.
      </div>
      <table v-else class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-slate-800/50 text-[11px] font-bold text-slate-400 border-b border-slate-800">
            <th class="py-3 px-5">시각</th>
            <th class="py-3 px-4">행위자</th>
            <th class="py-3 px-4">작업 (Action)</th>
            <th class="py-3 px-4">대상 리소스</th>
            <th class="py-3 px-4">결과</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/60 text-xs">
          <tr v-for="l in logs" :key="l.id" class="hover:bg-slate-800/30">
            <td class="py-3.5 px-5 text-slate-400">{{ formatDate(l.created_at) }}</td>
            <td class="py-3.5 px-4 font-bold text-white">{{ l.actor_type }} #{{ l.actor_id.slice(0, 8) }}</td>
            <td class="py-3.5 px-4 font-bold text-blue-400">{{ l.action }}</td>
            <td class="py-3.5 px-4 text-slate-300">{{ l.resource_type }} ({{ l.resource_id }})</td>
            <td class="py-3.5 px-4">
              <span
                class="px-2 py-0.5 rounded-full text-[10px] font-bold"
                :class="l.result === 'SUCCESS' ? 'bg-emerald-500/10 text-emerald-400' : 'bg-red-500/10 text-red-400'"
              >
                {{ l.result }}
              </span>
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

const logs = ref<any[]>([]);

onMounted(async () => {
  logs.value = await adminFetch('/admin/audit-logs');
});

function formatDate(iso: string) {
  if (!iso) return '';
  return new Date(iso).toISOString().slice(0, 19).replace('T', ' ');
}
</script>
