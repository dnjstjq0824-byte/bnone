<script setup lang="ts">
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Home, Wallet, ReceiptText, BarChart3, Target, Repeat, Plane, ArrowLeftRight, Bell, Settings, MessageCircle, LayoutGrid, Eye, EyeOff, LogOut } from 'lucide-vue-next';
import { useAuthStore } from '../stores/auth';
import { useFinanceStore } from '../stores/finance';
const route = useRoute();
const router = useRouter();
const auth = useAuthStore();
const finance = useFinanceStore();
const navItems = [
  { path: '/app', label: '홈', icon: Home },
  { path: '/app/assets', label: '내 자산', icon: Wallet },
  { path: '/app/transactions', label: '거래 내역', icon: ReceiptText },
  { path: '/app/analytics', label: '소비 분석', icon: BarChart3 },
  { path: '/app/budgets', label: '예산 관리', icon: Target },
  { path: '/app/subscriptions', label: '정기 구독', icon: Repeat },
  { path: '/app/trips', label: '여행 경비', icon: Plane },
  { path: '/app/fx', label: '환율 계산', icon: ArrowLeftRight },
  { path: '/app/notifications', label: '알림', icon: Bell },
  { path: '/app/settings', label: '설정', icon: Settings },
  { path: '/app/support', label: '고객센터', icon: MessageCircle },
];
const mobileItems = [navItems[0], navItems[1], navItems[3], { path: '/app/all', label: '전체', icon: LayoutGrid }];
declare const __COMMIT_HASH__: string;
const commitHash = typeof __COMMIT_HASH__ !== 'undefined' ? __COMMIT_HASH__ : 'dev';
const title = computed(() => {
  if (route.path === '/app/future') return '다가오는 기능';
  if (route.path === '/app/all') return '전체 서비스';
  const found = navItems.find(item => isActive(item.path));
  return found?.label || '금융 서비스';
});
function isActive(path: string) { return path === '/app' ? route.path === path : route.path.startsWith(path); }
function logout() { auth.logout(); router.push('/login'); }
</script>
<template>
  <div class="finance-shell">
    <a href="#main-content" class="skip-link">본문으로 바로가기</a>
    <aside class="finance-sidebar">
      <RouterLink to="/app" class="brand-lockup"><span class="brand-mark">M</span> MOA</RouterLink>
      <p class="sidebar-caption">나를 위한 금융 생활</p>
      <nav aria-label="주 메뉴" class="space-y-1">
        <RouterLink v-for="item in navItems" :key="item.path" :to="item.path" class="side-link" :class="{ selected: isActive(item.path) }" :aria-current="isActive(item.path) ? 'page' : undefined"><component :is="item.icon" :size="20" aria-hidden="true" />{{ item.label }}</RouterLink>
      </nav>
      <div class="sidebar-profile"><span class="profile-avatar">{{ (auth.user?.display_name || '나').charAt(0) }}</span><div class="min-w-0 flex-1"><strong class="block truncate text-white">{{ auth.user?.display_name || '사용자' }}님</strong><span class="text-xs text-[#7E9387]">내 금융을 한곳에서</span></div><button @click="logout" aria-label="로그아웃" class="icon-button !text-[#A9BCAF]"><LogOut :size="18" /></button></div>
      <div class="px-3 py-2 text-[10px] text-[#7E9387] flex items-center justify-between border-t border-white/10">
        <span>브라우저 목업 데모</span>
        <span class="font-mono">#{{ commitHash }}</span>
      </div>
    </aside>
    <div class="finance-main">
      <header class="finance-header"><span class="hidden md:block font-semibold text-slate-600">{{ title }}</span><RouterLink to="/app" class="md:hidden brand-lockup"><span class="brand-mark">M</span> MOA</RouterLink><div class="flex items-center gap-2"><button @click="finance.maskAmounts = !finance.maskAmounts" :aria-pressed="finance.maskAmounts" :aria-label="finance.maskAmounts ? '금액 표시' : '금액 숨기기'" class="icon-button"><EyeOff v-if="finance.maskAmounts" :size="20" /><Eye v-else :size="20" /></button><RouterLink to="/app/notifications" class="icon-button" aria-label="알림"><Bell :size="21" /></RouterLink></div></header>
      <main id="main-content" class="finance-content" tabindex="-1"><RouterView /></main>
    </div>
    <nav class="mobile-tabs" aria-label="모바일 주 메뉴"><RouterLink v-for="item in mobileItems" :key="item.path" :to="item.path" :class="{ selected: isActive(item.path) }" :aria-current="isActive(item.path) ? 'page' : undefined"><component :is="item.icon" :size="23" aria-hidden="true" /><span>{{ item.label }}</span></RouterLink></nav>
  </div>
</template>
