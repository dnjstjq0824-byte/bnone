<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-black text-white tracking-tight">기능 플래그 통제 (Feature Flags)</h1>
      <p class="text-xs text-slate-400 mt-0.5">환경별 기능 릴리스 및 모듈 활성화 여부를 동적으로 통제합니다.</p>
    </div>

    <!-- Article 21.12 Warning -->
    <div class="p-3.5 rounded-2xl bg-amber-500/10 border border-amber-500/20 text-amber-300 text-xs flex items-center gap-2.5">
      <span>⚠️</span>
      <p class="text-[11px] leading-relaxed">
        <strong>보안 정책 (Article 21.12):</strong> 실제 고객 송금, 결제 집행 및 금융 자금 보관 기능은 단순 Feature Flag 변경만으로 활성화할 수 없도록 원천 차단되어 있습니다.
      </p>
    </div>

    <!-- Feature Flags Table -->
    <div class="bg-[#111827] rounded-3xl border border-slate-800 overflow-hidden">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-slate-800/50 text-[11px] font-bold text-slate-400 border-b border-slate-800">
            <th class="py-3 px-5">플래그 키 (Key)</th>
            <th class="py-3 px-4">설명</th>
            <th class="py-3 px-4">대상 환경</th>
            <th class="py-3 px-4 text-center">상태</th>
            <th class="py-3 px-4 text-center">토글</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/60 text-xs">
          <tr v-for="f in flags" :key="f.key" class="hover:bg-slate-800/30">
            <td class="py-3.5 px-5 font-bold text-blue-400 font-mono text-[11px]">{{ f.key }}</td>
            <td class="py-3.5 px-4 text-slate-300">{{ f.description }}</td>
            <td class="py-3.5 px-4 text-slate-500">{{ f.environment }}</td>
            <td class="py-3.5 px-4 text-center">
              <span
                class="px-2 py-0.5 rounded-full text-[10px] font-bold"
                :class="f.enabled ? 'bg-emerald-500/10 text-emerald-400' : 'bg-slate-800 text-slate-500'"
              >
                {{ f.enabled ? '활성화' : '비활성' }}
              </span>
            </td>
            <td class="py-3.5 px-4 text-center">
              <button
                @click="toggleFlag(f)"
                class="px-3 py-1 rounded-xl text-xs font-bold transition-all"
                :class="f.enabled ? 'bg-slate-800 hover:bg-slate-700 text-slate-300' : 'bg-blue-600 hover:bg-blue-500 text-white'"
              >
                {{ f.enabled ? '끄기' : '켜기' }}
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

const flags = ref<any[]>([]);

async function loadFlags() {
  flags.value = await adminFetch('/admin/feature-flags');
}

onMounted(() => {
  loadFlags();
});

async function toggleFlag(f: any) {
  try {
    await adminFetch(`/admin/feature-flags/${f.key}`, {
      method: 'PATCH',
      body: JSON.stringify({ enabled: !f.enabled }),
    });
    await loadFlags();
  } catch (err: any) {
    alert(err.message || '플래그 토글에 실패했습니다.');
  }
}
</script>
