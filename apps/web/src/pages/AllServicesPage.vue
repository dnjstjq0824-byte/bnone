<script setup lang="ts">
import { computed, ref } from 'vue';
import { Search, ChevronRight } from 'lucide-vue-next';
const query = ref('');
const groups = [
  { title: '내 돈 관리', items: [ { path: 'assets', name: '내 자산', detail: '계정 · 저축 · 자산 · 부채' }, { path: 'transactions', name: '거래 내역', detail: '수입과 지출 기록 확인' }, { path: 'transactions/new', name: '거래 기록하기', detail: '수입 · 지출 · 계정 이동' }, { path: 'analytics', name: '소비 분석', detail: '내 소비 패턴 알아보기' }, { path: 'budgets', name: '예산 관리', detail: '이번 달 소비 계획 세우기' } ] },
  { title: '생활과 여행', items: [ { path: 'subscriptions', name: '정기 구독', detail: '매달 나가는 구독료 관리' }, { path: 'trips', name: '여행 경비', detail: '여행별 예산과 지출 기록' }, { path: 'fx', name: '환율 계산', detail: '참고 환율로 외화 계산' } ] },
  { title: '이용 안내', items: [ { path: 'notifications', name: '알림', detail: '내게 도착한 소식' }, { path: 'settings', name: '설정', detail: '내 정보 · 보안 · 데이터 관리' }, { path: 'support', name: '고객센터', detail: '궁금한 점 문의하기' } ] },
];
const filtered = computed(() => groups.map(group => ({ ...group, items: group.items.filter(item => `${item.name} ${item.detail}`.includes(query.value.trim())) })).filter(group => group.items.length));
</script>
<template>
  <div class="max-w-3xl mx-auto space-y-6"><div><p class="eyebrow">BN ONE</p><h1 class="text-3xl font-bold mt-2">전체 서비스</h1></div><label class="service-search"><Search :size="21" aria-hidden="true" /><input v-model="query" type="search" placeholder="어떤 서비스가 필요하세요?" aria-label="서비스 검색" /></label><p v-if="!filtered.length" class="empty-state" role="status">검색한 서비스가 없어요. 다른 단어로 찾아보세요.</p><section v-for="group in filtered" :key="group.title" class="finance-card"><h2 class="text-sm font-semibold text-slate-500 mb-3">{{ group.title }}</h2><RouterLink v-for="item in group.items" :key="item.path" :to="`/app/${item.path}`" class="financial-row"><div class="flex-1"><h3 class="font-semibold">{{ item.name }}</h3><p class="text-sm text-slate-500 mt-1">{{ item.detail }}</p></div><ChevronRight :size="19" class="text-slate-400" /></RouterLink></section></div>
</template>
