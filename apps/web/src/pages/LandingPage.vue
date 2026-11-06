<template>
  <div class="min-h-screen bg-[#08170F] text-white flex flex-col">
    <!-- Navbar -->
    <header class="max-w-6xl mx-auto w-full px-6 h-20 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-xl bg-[#C0A04E] flex items-center justify-center font-black text-xl text-[#08170F]">
          M
        </div>
        <div class="leading-none">
          <span class="font-extrabold text-xl tracking-tight">MOA</span>
          <span class="block text-[10px] font-semibold text-[#7E9387] tracking-[0.2em] mt-1">모아 · PRIVATE LEDGER</span>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <RouterLink to="/login" class="text-sm font-semibold text-[#A9BCAF] hover:text-white px-4 py-2 transition-colors">
          로그인
        </RouterLink>
        <RouterLink to="/register" class="text-sm font-bold bg-[#E9D9A8] hover:bg-white text-[#08170F] px-5 py-2.5 rounded-xl transition-all">
          시작하기
        </RouterLink>
      </div>
    </header>

    <!-- Hero -->
    <main class="max-w-6xl mx-auto px-6 w-full grid grid-cols-1 lg:grid-cols-2 gap-12 items-center py-14 lg:py-20">
      <div>
        <div class="w-12 h-1 bg-[#C0A04E] rounded-full mb-7"></div>
        <h1 class="text-5xl md:text-6xl font-black tracking-tight leading-[1.08]">
          돈을 모으다,<br />
          <span class="text-[#E9D9A8]">삶을 모으다.</span>
        </h1>
        <p class="text-base text-[#A9BCAF] max-w-md mt-6 leading-relaxed">
          매일 쓰는 돈부터 세계 여행에서 쓰는 돈까지. 자산과 예산, 구독과 환율을
          하나의 장부에 정확하고 안전하게 기록하세요.
        </p>
        <div class="flex items-center gap-4 mt-9">
          <RouterLink to="/register" class="px-8 py-3.5 bg-[#E9D9A8] hover:bg-white text-[#08170F] rounded-xl font-bold text-base transition-all">
            무료로 시작하기
          </RouterLink>
          <RouterLink to="/login" class="text-sm font-semibold text-[#A9BCAF] hover:text-white transition-colors">
            기존 계정으로 로그인 →
          </RouterLink>
        </div>
        <div class="flex items-center gap-6 mt-10 text-xs text-[#7E9387]">
          <span><strong class="text-white text-sm font-extrabold">14</strong>개 통화</span>
          <span class="w-px h-4 bg-white/15"></span>
          <span><strong class="text-white text-sm font-extrabold">9</strong>개 자산 영역</span>
          <span class="w-px h-4 bg-white/15"></span>
          <span>은행 연동 없는 <strong class="text-white font-bold">수기 장부</strong></span>
        </div>
      </div>

      <!-- Product preview card -->
      <div class="rounded-3xl bg-white text-[#141F19] p-7 shadow-2xl shadow-black/40" aria-hidden="true">
        <div class="flex items-center justify-between">
          <p class="text-xs font-bold text-[#5B6960]">내 순자산</p>
          <span class="px-2.5 py-1 rounded-full bg-[#EBF2EC] text-[#1E5A38] text-[10px] font-extrabold">이번 달 +2.4%</span>
        </div>
        <p class="text-4xl font-black tracking-tight mt-2">₩303,152,500</p>
        <div class="flex gap-5 mt-3 text-xs text-[#5B6960]">
          <span>총자산 <strong class="text-[#141F19]">₩341,900,000</strong></span>
          <span>부채 <strong class="text-[#141F19]">₩38,747,500</strong></span>
        </div>
        <div class="mt-6 space-y-3">
          <div class="flex items-center gap-3 p-3 rounded-2xl bg-[#F6F8F4]">
            <span class="w-9 h-9 rounded-xl bg-[#123522] text-[#E9D9A8] grid place-items-center text-sm font-black">급</span>
            <div class="flex-1">
              <p class="text-xs font-bold">급여 입금</p>
              <p class="text-[10px] text-[#8A948C]">11월 5일 · 주계좌</p>
            </div>
            <strong class="text-sm text-[#1E5A38]">+₩4,200,000</strong>
          </div>
          <div class="flex items-center gap-3 p-3 rounded-2xl bg-[#F6F8F4]">
            <span class="w-9 h-9 rounded-xl bg-[#ECEFE8] text-[#41524A] grid place-items-center text-sm font-black">항</span>
            <div class="flex-1">
              <p class="text-xs font-bold">대한항공 도쿄 왕복</p>
              <p class="text-[10px] text-[#8A948C]">11월 3일 · 여행</p>
            </div>
            <strong class="text-sm">−₩486,000</strong>
          </div>
          <div class="flex items-center gap-3 p-3 rounded-2xl bg-[#F6F8F4]">
            <span class="w-9 h-9 rounded-xl bg-[#F2EAD3] text-[#9A7A33] grid place-items-center text-sm font-black">구</span>
            <div class="flex-1">
              <p class="text-xs font-bold">스트리밍 구독 4건</p>
              <p class="text-[10px] text-[#8A948C]">다음 갱신 11월 12일</p>
            </div>
            <strong class="text-sm">₩43,600/월</strong>
          </div>
        </div>
      </div>
    </main>

    <!-- Editorial feature rows -->
    <section class="max-w-6xl mx-auto px-6 w-full pb-16">
      <div class="border-t border-white/10">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8 py-8 border-b border-white/10">
          <div class="flex gap-4">
            <span class="text-[#C0A04E] font-black text-sm">01</span>
            <div>
              <h3 class="font-bold text-base">다중 통화 수기 장부</h3>
              <p class="text-[#7E9387] text-xs leading-relaxed mt-1.5">현금과 계좌, 외화 자산을 직접 기록하고 순자산을 추적합니다.</p>
            </div>
          </div>
          <div class="flex gap-4">
            <span class="text-[#C0A04E] font-black text-sm">02</span>
            <div>
              <h3 class="font-bold text-base">여행 예산과 경비</h3>
              <p class="text-[#7E9387] text-xs leading-relaxed mt-1.5">일정별 예산과 원화 환산 보고서를 자동으로 정리합니다.</p>
            </div>
          </div>
          <div class="flex gap-4">
            <span class="text-[#C0A04E] font-black text-sm">03</span>
            <div>
              <h3 class="font-bold text-base">구독과 갱신 알림</h3>
              <p class="text-[#7E9387] text-xs leading-relaxed mt-1.5">분산된 구독의 갱신일과 월 부담액을 한눈에 봅니다.</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="border-t border-white/10 py-8 px-6 text-center text-xs text-[#7E9387]">
      <div class="max-w-4xl mx-auto space-y-2">
        <p>Booking Network Co., Ltd — MOA v1.0 Premium Financial Lifestyle</p>
        <p class="text-[11px]">
          * MOA v1.0은 사용자가 직접 입력한 금융정보를 기반으로 자산과 지출을 관리하는 소프트웨어이며, 실제 금융기관 계좌 연결 및 금융거래를 직접 실행하지 않습니다.
        </p>
        <div class="pt-2 flex items-center justify-center gap-4">
          <RouterLink to="/terms" class="hover:text-white transition-colors">이용약관</RouterLink>
          <RouterLink to="/privacy" class="hover:text-white transition-colors">개인정보처리방침</RouterLink>
          <span class="font-mono">Commit: #{{ commitHash }}</span>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router';

declare const __COMMIT_HASH__: string;
const commitHash = typeof __COMMIT_HASH__ !== 'undefined' ? __COMMIT_HASH__ : 'dev';
</script>
