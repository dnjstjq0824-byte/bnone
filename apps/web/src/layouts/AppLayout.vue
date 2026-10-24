<template>
  <div class="min-h-screen bg-[#F8FAFC] flex flex-col md:flex-row">
    <!-- Desktop Sidebar -->
    <aside class="hidden md:flex flex-col w-64 bg-white border-r border-slate-200 p-5 shrink-0 justify-between">
      <div>
        <!-- Brand -->
        <div class="flex items-center gap-3 mb-8 px-2">
          <div class="w-10 h-10 rounded-xl bg-brand-600 flex items-center justify-center text-white font-bold text-xl shadow-md shadow-brand-500/30">
            BN
          </div>
          <div>
            <div class="font-extrabold text-slate-900 tracking-tight text-lg">BN ONE</div>
            <div class="text-[11px] text-slate-400 font-medium">Your Money. Your World.</div>
          </div>
        </div>

        <!-- Navigation Links -->
        <nav class="space-y-1">
          <RouterLink
            v-for="item in navItems"
            :key="item.path"
            :to="item.path"
            class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-semibold transition-all"
            :class="isActive(item.path) ? 'bg-brand-50 text-brand-600' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'"
          >
            <span>{{ item.icon }}</span>
            <span>{{ item.label }}</span>
          </RouterLink>
        </nav>
      </div>

      <!-- User footer -->
      <div class="pt-4 border-t border-slate-100 flex items-center justify-between">
        <div class="flex items-center gap-2 overflow-hidden">
          <div class="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center text-xs font-bold text-slate-600">
            {{ userInitial }}
          </div>
          <div class="text-xs truncate">
            <div class="font-bold text-slate-800 truncate">{{ auth.user?.display_name || 'User' }}</div>
            <div class="text-slate-400 text-[10px]">{{ auth.user?.base_currency || 'KRW' }}</div>
          </div>
        </div>
        <button @click="handleLogout" class="text-xs text-slate-400 hover:text-red-600 font-medium transition-colors">
          로그아웃
        </button>
      </div>
    </aside>

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col min-w-0 pb-16 md:pb-0">
      <!-- Top header -->
      <header class="h-16 bg-white border-b border-slate-200 px-4 md:px-8 flex items-center justify-between sticky top-0 z-20">
        <div class="flex items-center gap-3">
          <div class="md:hidden font-extrabold text-brand-600 text-lg">BN ONE</div>
          <span class="text-xs px-2.5 py-1 bg-slate-100 text-slate-600 rounded-full font-medium">
            기준통화: {{ auth.user?.base_currency || 'KRW' }}
          </span>
        </div>

        <div class="flex items-center gap-3">
          <!-- Mask Amounts toggle -->
          <button
            @click="finance.maskAmounts = !finance.maskAmounts"
            class="p-2 text-xs font-medium text-slate-600 hover:bg-slate-100 rounded-lg flex items-center gap-1 transition-colors"
            title="금액 숨기기"
          >
            <span>{{ finance.maskAmounts ? '🔒 금액표시' : '👁️ 금액숨김' }}</span>
          </button>

          <!-- Quick Transaction Button -->
          <RouterLink
            to="/app/transactions/new"
            class="px-3.5 py-2 bg-brand-600 hover:bg-brand-700 text-white rounded-xl text-xs font-bold shadow-sm transition-all flex items-center gap-1.5"
          >
            <span>+</span>
            <span>거래 입력</span>
          </RouterLink>
        </div>
      </header>

      <!-- Page View -->
      <main class="flex-1 p-4 md:p-8 max-w-[1280px] w-full mx-auto">
        <RouterView />
      </main>
    </div>

    <!-- Mobile Bottom Navigation -->
    <nav class="md:hidden fixed bottom-0 left-0 right-0 h-16 bg-white border-t border-slate-200 flex items-center justify-around z-30 px-2">
      <RouterLink
        v-for="item in mobileNavItems"
        :key="item.path"
        :to="item.path"
        class="flex flex-col items-center justify-center w-14 py-1 text-slate-500 transition-colors"
        :class="isActive(item.path) ? 'text-brand-600 font-bold' : ''"
      >
        <span class="text-lg">{{ item.icon }}</span>
        <span class="text-[10px] mt-0.5">{{ item.label }}</span>
      </RouterLink>
    </nav>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { useFinanceStore } from '../stores/finance';

const route = useRoute();
const router = useRouter();
const auth = useAuthStore();
const finance = useFinanceStore();

const userInitial = computed(() => {
  return (auth.user?.display_name || auth.user?.email || 'U').charAt(0).toUpperCase();
});

const navItems = [
  { path: '/app', label: '대시보드', icon: '📊' },
  { path: '/app/assets', label: '자산 & 부채', icon: '🏛️' },
  { path: '/app/transactions', label: '거래 내역', icon: '📝' },
  { path: '/app/analytics', label: '소비 분석', icon: '📈' },
  { path: '/app/budgets', label: '예산 관리', icon: '🎯' },
  { path: '/app/subscriptions', label: '정기 구독', icon: '🔄' },
  { path: '/app/trips', label: '여행 경비', icon: '✈️' },
  { path: '/app/fx', label: '참고 환율', icon: '💱' },
  { path: '/app/notifications', label: '알림', icon: '🔔' },
  { path: '/app/settings', label: '설정 & 보안', icon: '⚙️' },
  { path: '/app/support', label: '고객지원', icon: '💬' },
];

const mobileNavItems = [
  { path: '/app', label: '홈', icon: '🏠' },
  { path: '/app/assets', label: '자산', icon: '🏛️' },
  { path: '/app/transactions', label: '내역', icon: '📝' },
  { path: '/app/budgets', label: '예산', icon: '🎯' },
  { path: '/app/settings', label: '더보기', icon: '⚙️' },
];

function isActive(path: string): boolean {
  if (path === '/app') return route.path === '/app';
  return route.path.startsWith(path);
}

function handleLogout() {
  auth.logout();
  router.push('/login');
}
</script>

