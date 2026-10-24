<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-black text-white tracking-tight">FX 데이터 공급자 및 환율 운영</h1>
        <p class="text-xs text-slate-400 mt-0.5">공식 참고 환율 Provider 등록 상태 및 자동 수집 배치를 관리합니다.</p>
      </div>
      <button
        @click="triggerSync"
        class="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-xl text-xs font-bold transition-all"
      >
        수집 동기화 즉시 실행
      </button>
    </div>

    <!-- Providers List -->
    <div class="bg-[#111827] rounded-3xl border border-slate-800 p-6 space-y-4">
      <h2 class="text-base font-bold text-white">등록된 환율 Provider</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div v-for="p in providers" :key="p.id" class="p-4 rounded-2xl bg-slate-900 border border-slate-800 space-y-2">
          <div class="flex items-center justify-between">
            <span class="font-bold text-xs text-white">{{ p.name }}</span>
            <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-500/10 text-emerald-400">
              {{ p.active ? '활성' : '비활성' }}
            </span>
          </div>
          <div class="text-[11px] text-slate-400 truncate">{{ p.official_website }}</div>
          <div class="text-[10px] text-slate-500 pt-2 border-t border-slate-800 flex justify-between">
            <span>라이선스: {{ p.license_status }}</span>
            <span>최종 점검: {{ formatDate(p.last_checked_at) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { adminFetch } from '../api/client';

const providers = ref<any[]>([]);

async function loadData() {
  providers.value = await adminFetch('/admin/fx/providers');
}

onMounted(() => {
  loadData();
});

async function triggerSync() {
  const res = await adminFetch('/admin/fx/sync-jobs', { method: 'POST' });
  alert(res.message);
  await loadData();
}

function formatDate(iso: string) {
  if (!iso) return '미확인';
  return new Date(iso).toISOString().slice(0, 10);
}
</script>
