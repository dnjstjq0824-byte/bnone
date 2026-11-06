<template>
  <div class="min-h-screen bg-[#08170F] text-slate-100 flex flex-col md:flex-row">
    <!-- Left Sidebar -->
    <aside class="w-full md:w-64 bg-[#0C2318] border-r border-white/10 p-5 flex flex-col justify-between shrink-0">
      <div>
        <!-- Admin Console Header -->
        <div class="flex items-center gap-3 mb-8 px-2">
          <div class="w-9 h-9 rounded-xl bg-[#C0A04E] flex items-center justify-center font-black text-[#08170F] text-lg">
            M
          </div>
          <div>
            <div class="font-extrabold text-white tracking-tight text-sm">MOA OPERATIONS</div>
            <div class="text-[10px] text-[#C0A04E] font-semibold uppercase tracking-wider">Master Admin</div>
          </div>
        </div>

        <!-- Navigation Menu -->
        <nav class="space-y-1">
          <RouterLink
            v-for="item in navItems"
            :key="item.path"
            :to="item.path"
            class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-xs font-bold transition-all"
            :class="isActive(item.path) ? 'bg-[#C0A04E]/15 text-[#E9D9A8]' : 'text-[#A9BCAF] hover:bg-white/5 hover:text-white'"
          >
            <component :is="item.icon" :size="17" />
            <span>{{ item.label }}</span>
          </RouterLink>
        </nav>
      </div>

      <!-- Admin Identity & Logout -->
      <div class="pt-4 border-t border-white/10 flex items-center justify-between">
        <div class="text-xs truncate">
          <div class="font-bold text-white truncate">{{ adminAuth.fullName || 'Operator' }}</div>
          <span class="inline-block px-1.5 py-0.5 rounded bg-[#C0A04E]/15 text-[#E9D9A8] text-[9px] font-bold mt-0.5">
            {{ adminAuth.role || 'ADMIN' }}
          </span>
        </div>
        <button @click="handleLogout" class="text-xs text-[#7E9387] hover:text-red-400 font-semibold transition-colors">
          로그아웃
        </button>
      </div>
    </aside>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col min-w-0">
      <!-- Top Bar -->
      <header class="h-16 bg-[#0C2318] border-b border-white/10 px-6 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <span class="px-2.5 py-1 rounded-full bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 text-[10px] font-extrabold uppercase tracking-wider flex items-center gap-1.5">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse"></span>
            MOA v1.0 LIVE
          </span>
        </div>

        <div class="text-xs text-[#7E9387] font-medium">
          {{ adminAuth.adminEmail }}
        </div>
      </header>

      <!-- View Area -->
      <main class="flex-1 p-6 md:p-8 max-w-[1400px] w-full mx-auto overflow-auto">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter, RouterLink, RouterView } from 'vue-router';
import { LayoutDashboard, Users, MessagesSquare, ArrowLeftRight, ShieldCheck, ScrollText, Flag } from 'lucide-vue-next';
import { useAdminAuthStore } from '../stores/adminAuth';

const route = useRoute();
const router = useRouter();
const adminAuth = useAdminAuthStore();

const navItems = [
  { path: '/admin/dashboard', label: '운영 대시보드', icon: LayoutDashboard },
  { path: '/admin/users', label: '사용자 관리', icon: Users },
  { path: '/admin/support', label: '고객지원 티켓', icon: MessagesSquare },
  { path: '/admin/fx', label: 'FX 공급자 및 환율', icon: ArrowLeftRight },
  { path: '/admin/privacy', label: '개인정보·삭제 요청', icon: ShieldCheck },
  { path: '/admin/audit-logs', label: '불변 감사 로그', icon: ScrollText },
  { path: '/admin/feature-flags', label: '기능 플래그 통제', icon: Flag },
];

function isActive(path: string): boolean {
  return route.path === path;
}

function handleLogout() {
  adminAuth.logout();
  router.push('/admin/login');
}
</script>
