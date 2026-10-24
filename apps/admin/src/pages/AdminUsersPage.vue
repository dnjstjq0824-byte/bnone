<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-black text-white tracking-tight">사용자 계정 운영</h1>
        <p class="text-xs text-slate-400 mt-0.5">등록된 사용자 계정 상태, 로그인 세션 및 보안 제재를 관리합니다.</p>
      </div>
    </div>

    <!-- Search Input -->
    <div class="bg-[#111827] p-4 rounded-2xl border border-slate-800 flex items-center gap-3">
      <input
        v-model="search"
        @input="loadUsers"
        type="text"
        placeholder="사용자 이메일, 이름 검색..."
        class="w-full px-4 py-2 rounded-xl bg-slate-900 border border-slate-700 text-xs text-white focus:outline-none focus:border-blue-500"
      />
    </div>

    <!-- Users Table -->
    <div class="bg-[#111827] rounded-3xl border border-slate-800 overflow-hidden">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-slate-800/50 text-[11px] font-bold text-slate-400 border-b border-slate-800">
            <th class="py-3 px-5">이메일</th>
            <th class="py-3 px-4">이름</th>
            <th class="py-3 px-4">국가 / 통화</th>
            <th class="py-3 px-4">상태</th>
            <th class="py-3 px-4">가입일</th>
            <th class="py-3 px-4 text-center">조치</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/60 text-xs">
          <tr v-for="u in users" :key="u.id" class="hover:bg-slate-800/30">
            <td class="py-3.5 px-5 font-bold text-white">{{ u.email }}</td>
            <td class="py-3.5 px-4 text-slate-300">{{ u.display_name }}</td>
            <td class="py-3.5 px-4 text-slate-400">{{ u.country_code }} • {{ u.base_currency }}</td>
            <td class="py-3.5 px-4">
              <span
                class="px-2 py-0.5 rounded-full text-[10px] font-bold"
                :class="u.status === 'ACTIVE' ? 'bg-emerald-500/10 text-emerald-400' : 'bg-red-500/10 text-red-400'"
              >
                {{ u.status }}
              </span>
            </td>
            <td class="py-3.5 px-4 text-slate-500">{{ formatDate(u.created_at) }}</td>
            <td class="py-3.5 px-4 text-center space-x-2">
              <button
                @click="toggleStatus(u)"
                class="text-[11px] font-bold"
                :class="u.status === 'ACTIVE' ? 'text-amber-400 hover:text-amber-300' : 'text-emerald-400 hover:text-emerald-300'"
              >
                {{ u.status === 'ACTIVE' ? '이용정지' : '정지해제' }}
              </button>
              <button
                @click="revokeSessions(u.id)"
                class="text-[11px] font-bold text-red-400 hover:text-red-300"
              >
                세션 강제종료
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

const users = ref<any[]>([]);
const search = ref('');

async function loadUsers() {
  users.value = await adminFetch(`/admin/users?q=${encodeURIComponent(search.value)}`);
}

onMounted(() => {
  loadUsers();
});

async function toggleStatus(u: any) {
  const nextStatus = u.status === 'ACTIVE' ? 'SUSPENDED' : 'ACTIVE';
  if (confirm(`사용자 상태를 ${nextStatus}(으)로 변경하시겠습니까?`)) {
    await adminFetch(`/admin/users/${u.id}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ status: nextStatus }),
    });
    await loadUsers();
  }
}

async function revokeSessions(id: string) {
  if (confirm('이 사용자의 모든 활성 세션을 즉시 철회하시겠습니까?')) {
    await adminFetch(`/admin/users/${id}/revoke-sessions`, { method: 'POST' });
    alert('모든 세션이 종료되었습니다.');
  }
}

function formatDate(iso: string) {
  if (!iso) return '';
  return new Date(iso).toISOString().slice(0, 10);
}
</script>

