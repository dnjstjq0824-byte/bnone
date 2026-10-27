<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { ArrowUpRight, ArrowDownLeft, ArrowLeftRight, ChevronRight, Plus, Wallet, ReceiptText, Target, Plane, Repeat, RefreshCw } from 'lucide-vue-next';
import { useAuthStore } from '../stores/auth';
import { useFinanceStore } from '../stores/finance';
const auth = useAuthStore();
const finance = useFinanceStore();
const currency = computed(() => finance.netWorth?.base_currency || auth.user?.base_currency || 'KRW');
const money = (value: number = 0) => finance.formatCurrency(value, currency.value);
const month = new Date().getMonth() + 1;
const shortcuts = [
  { path: '/app/transactions/new', label: '거래 기록', icon: Plus, color: 'blue' },
  { path: '/app/assets', label: '자산 관리', icon: Wallet, color: 'purple' },
  { path: '/app/budgets', label: '예산 관리', icon: Target, color: 'green' },
  { path: '/app/fx', label: '환율 계산', icon: ArrowLeftRight, color: 'orange' },
];
const formatDate = (iso: string) => new Intl.DateTimeFormat('ko-KR', { month: 'long', day: 'numeric' }).format(new Date(iso));
onMounted(() => finance.loadDashboard());
</script>
<template>
  <div class="home-dashboard" :aria-busy="finance.isLoading">
    <div class="dashboard-heading"><div><p class="eyebrow">나의 금융 생활</p><h1>{{ auth.user?.display_name || '사용자' }}님, 반가워요</h1><p class="text-slate-500 mt-2">오늘도 내 돈의 흐름을 가볍게 확인하세요.</p></div><button class="icon-button" @click="finance.loadDashboard" :disabled="finance.isLoading" aria-label="금융 정보 새로고침"><RefreshCw :size="20" :class="{ 'animate-spin': finance.isLoading }" /></button></div>
    <div v-if="finance.isLoading" class="dashboard-loading" role="status"><div v-for="n in 3" :key="n" class="finance-card skeleton-card"></div><span class="sr-only">금융 정보를 불러오는 중이에요</span></div>
    <div v-else-if="finance.dashboardError" class="finance-card empty-state" role="alert"><Wallet :size="36" /><h2>잠시 연결이 원활하지 않아요</h2><p>{{ finance.dashboardError }}</p><button class="primary-button" @click="finance.loadDashboard">다시 불러오기</button></div>
    <div v-else class="dashboard-columns">
      <div class="space-y-5">
        <section class="finance-card wealth-card"><RouterLink to="/app/assets" class="section-heading"><h2>내 순자산</h2><ChevronRight :size="20" /></RouterLink><p class="wealth-amount">{{ money(finance.netWorth?.net_worth_minor) }}</p><div class="wealth-detail"><span>총자산 <strong>{{ money(finance.netWorth?.total_assets_minor) }}</strong></span><span>부채 <strong>{{ money(finance.netWorth?.total_liabilities_minor) }}</strong></span></div><div class="quick-actions"><RouterLink v-for="item in shortcuts" :key="item.path" :to="item.path"><span class="shortcut-icon" :class="item.color"><component :is="item.icon" :size="23" /></span><span>{{ item.label }}</span></RouterLink></div></section>
        <section class="finance-card"><RouterLink to="/app/assets" class="section-heading"><h2>내 계정 <span class="text-blue-500">{{ finance.accounts.length }}</span></h2><ChevronRight :size="20" /></RouterLink><div v-if="!finance.accounts.length" class="empty-state"><Wallet :size="32" /><h3>첫 번째 자산을 등록해보세요</h3><p>계정과 현금, 저축을 한눈에 관리해요.</p><RouterLink to="/app/assets" class="primary-button">자산 등록하기</RouterLink></div><div v-else><RouterLink v-for="account in finance.accounts.slice(0, 5)" :key="account.id" to="/app/assets" class="financial-row"><span class="account-avatar" :style="{ backgroundColor: account.color || '#3182f6' }"><Wallet :size="21" /></span><div class="flex-1 min-w-0"><p class="text-sm text-slate-500 truncate">{{ account.name }}</p><p class="row-amount">{{ finance.formatCurrency(account.calculated_balance_minor, account.currency_code) }}</p></div><ChevronRight :size="18" class="text-slate-400" /></RouterLink><RouterLink to="/app/assets" class="soft-button mt-4">내 자산 모두 보기</RouterLink></div></section>
        <section class="finance-card"><RouterLink to="/app/transactions" class="section-heading"><h2>최근 거래</h2><ChevronRight :size="20" /></RouterLink><div v-if="!finance.recentTransactions.length" class="empty-state"><ReceiptText :size="32" /><h3>아직 기록한 거래가 없어요</h3><p>작은 소비부터 하나씩 기록해보세요.</p><RouterLink to="/app/transactions/new" class="soft-button">첫 거래 기록하기</RouterLink></div><div v-else><div v-for="tx in finance.recentTransactions.slice(0, 5)" :key="tx.id" class="financial-row"><span class="transaction-icon" :class="{ income: tx.transaction_type === 'INCOME' }"><ArrowDownLeft v-if="tx.transaction_type === 'INCOME'" :size="20" /><ArrowLeftRight v-else-if="tx.transaction_type === 'TRANSFER'" :size="20" /><ArrowUpRight v-else :size="20" /></span><div class="flex-1 min-w-0"><p class="font-semibold truncate">{{ tx.description || tx.merchant_name || '거래 기록' }}</p><p class="text-xs text-slate-500 mt-1">{{ formatDate(tx.occurred_at) }}<span v-if="tx.category_name"> · {{ tx.category_name }}</span></p></div><strong class="text-sm whitespace-nowrap" :class="{ 'text-blue-600': tx.transaction_type === 'INCOME' }">{{ finance.maskAmounts ? '' : tx.transaction_type === 'EXPENSE' ? '−' : tx.transaction_type === 'INCOME' ? '+' : '' }}{{ finance.formatCurrency(tx.amount_minor, tx.currency_code) }}</strong></div></div></section>
      </div>
      <div class="space-y-5">
        <section class="finance-card spending-card">
          <RouterLink to="/app/analytics" class="section-heading">
            <h2>{{ month }}월 소비</h2>
            <ChevronRight :size="20" />
          </RouterLink>
          <p class="text-sm text-slate-500 mt-6">이번 달 쓴 금액</p>
          <p class="spending-amount">{{ money(finance.analyticsOverview?.current_month_expense_minor) }}</p>
          <div class="monthly-income">
            <span>들어온 금액</span>
            <strong>{{ money(finance.analyticsOverview?.current_month_income_minor) }}</strong>
          </div>
          <div v-if="finance.analyticsOverview?.category_breakdown?.length" class="space-y-5 mt-6">
            <div v-for="cat in finance.analyticsOverview.category_breakdown.slice(0, 3)" :key="cat.category_id">
              <div class="flex justify-between text-sm gap-3 mb-2">
                <span class="text-slate-500">{{ cat.category_name }}</span>
                <strong>{{ money(cat.amount_minor) }}</strong>
              </div>
              <div class="spend-track" aria-hidden="true">
                <div :style="{ width: `${Math.min(100, Math.max(0, cat.percentage || 0))}%`, backgroundColor: cat.color || '#3182f6' }"></div>
              </div>
            </div>
          </div>
          <RouterLink to="/app/budgets" class="soft-button mt-6">
            이번 달 예산 관리하기 <ChevronRight :size="16" />
          </RouterLink>
        </section>
        <section class="finance-card"><h2 class="text-lg font-bold mb-3">놓치지 말고 챙겨요</h2><RouterLink to="/app/subscriptions" class="financial-row"><span class="shortcut-icon purple"><Repeat :size="22" /></span><div class="flex-1"><strong class="text-sm">매달 나가는 구독료</strong><p class="text-xs text-slate-500 mt-1">정기 결제를 한곳에서 관리해요</p></div><ChevronRight :size="18" /></RouterLink><RouterLink to="/app/trips" class="financial-row"><span class="shortcut-icon green"><Plane :size="22" /></span><div class="flex-1"><strong class="text-sm">여행도 계획적으로</strong><p class="text-xs text-slate-500 mt-1">여행 예산과 외화 지출 기록</p></div><ChevronRight :size="18" /></RouterLink></section>
        <RouterLink to="/app/all" class="service-banner"><span><strong>내게 필요한 금융 생활</strong><span class="block text-sm mt-2">BN ONE 전체 서비스 둘러보기</span></span><ArrowUpRight :size="26" /></RouterLink><p class="ledger-note">직접 입력한 기록을 기준으로 표시해요.<br />계정 이동 기록은 실제 은행 송금이 아닙니다.</p>
      </div>
    </div>
  </div>
</template>
