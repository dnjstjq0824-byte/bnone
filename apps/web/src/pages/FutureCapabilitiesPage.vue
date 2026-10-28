<script setup lang="ts">
import { Rocket, Sparkles, Camera, Users, Briefcase, ShieldCheck } from 'lucide-vue-next';

const items = [
  {
    icon: Briefcase,
    title: '오픈뱅킹 자동 연동',
    version: 'v2.0',
    desc: '금융기관 계좌와 카드 거래를 읽기 전용으로 자동 수집합니다. MyData/AISP 라이선스와 명시적 동의 절차가 선행됩니다.',
  },
  {
    icon: Camera,
    title: '영수증 OCR 자동 입력',
    version: 'v2.1',
    desc: '카메라로 종이 영수증을 촬영하면 상호명·날짜·통화·세금·품목을 자동 추출합니다. 사용자가 확인 후 가계부에 반영합니다.',
  },
  {
    icon: Users,
    title: '공동 여행 지갑',
    version: 'v2.2',
    desc: '여행 워크스페이스에서 동행자와 지출을 공동 기록하고 정산 금액을 계산합니다.',
  },
  {
    icon: Sparkles,
    title: '멀티자산 포트폴리오',
    version: 'v2.3',
    desc: '주식·ETF·펀드·디지털 자산의 실시간 시세와 미실현 손익을 추적합니다.',
  },
];
</script>

<template>
  <div class="max-w-3xl mx-auto space-y-6">
    <div>
      <p class="eyebrow">BN ONE</p>
      <h1 class="text-3xl font-bold mt-2">다가오는 기능</h1>
      <p class="text-sm text-slate-500 mt-2">
        v1.0에서는 수동 가계부와 참고 환율만 제공합니다. 아래 기능은 v2.0 이후에 순차적으로 제공될 예정입니다.
      </p>
    </div>

    <section class="finance-card">
      <div class="flex items-center gap-2 mb-4">
        <Rocket :size="18" class="text-brand-600" />
        <h2 class="text-sm font-semibold text-slate-500">로드맵</h2>
      </div>
      <ul class="space-y-4">
        <li v-for="item in items" :key="item.title" class="flex gap-4">
          <div class="w-10 h-10 rounded-2xl bg-brand-50 text-brand-600 grid place-items-center shrink-0">
            <component :is="item.icon" :size="20" />
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 flex-wrap">
              <h3 class="font-semibold text-slate-800">{{ item.title }}</h3>
              <span class="text-[11px] font-bold px-2 py-0.5 rounded-full bg-slate-100 text-slate-500">{{ item.version }} 예정</span>
            </div>
            <p class="text-sm text-slate-500 mt-1 leading-relaxed">{{ item.desc }}</p>
          </div>
        </li>
      </ul>
    </section>

    <section class="finance-card">
      <div class="flex items-center gap-2 mb-2">
        <ShieldCheck :size="18" class="text-slate-500" />
        <h2 class="text-sm font-semibold text-slate-500">v1.0 범위 안내</h2>
      </div>
      <p class="text-sm text-slate-600 leading-relaxed">
        BN ONE v1.0은 수동 금융 가계부와 참고 환율 정보만 제공합니다. 자동 은행 연동, 카드 결제, 송금, 자금 보관 기능은 제공하지 않습니다.
      </p>
    </section>
  </div>
</template>