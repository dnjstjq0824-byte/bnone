<template>
  <div class="min-h-screen bg-[#08170F] text-white flex flex-col">
    <!-- Navbar -->
    <header class="max-w-6xl mx-auto w-full px-6 h-20 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-xl bg-[#C0A04E] flex items-center justify-center font-black text-xl text-[#08170F]">
          모
        </div>
        <span class="font-extrabold text-xl tracking-tight">MOA</span>
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
    <main class="max-w-6xl mx-auto px-6 w-full grid grid-cols-1 lg:grid-cols-[1.05fr_1fr] gap-12 items-start py-14 lg:py-20">
      <div class="pt-2">
        <h1 class="text-5xl md:text-[64px] font-black tracking-tight leading-[1.06]">
          쓸 때마다 적고,<br />
          월말에 웃는다.
        </h1>
        <p class="text-[15px] text-[#A9BCAF] max-w-md mt-6 leading-relaxed">
          MOA는 은행에 연결하지 않습니다. 당신이 직접 적는 장부라서 더 정확하고,
          통장 잔고가 아닌 <span class="text-white font-semibold">내 돈의 흐름</span>이 보입니다.
          하루 10초, 쓸 때 바로 적으세요.
        </p>
        <div class="flex items-center gap-4 mt-9">
          <RouterLink to="/register" class="px-8 py-3.5 bg-[#E9D9A8] hover:bg-white text-[#08170F] rounded-xl font-bold text-base transition-all">
            무료로 시작하기
          </RouterLink>
          <RouterLink to="/login" class="text-sm font-semibold text-[#A9BCAF] hover:text-white transition-colors">
            기존 계정으로 로그인 →
          </RouterLink>
        </div>
        <p class="text-xs text-[#5B6960] mt-8">은행 연동 없음 · 카드 긁힘 없음 · CSV 내보내기 지원</p>
      </div>

      <!-- Ledger preview -->
      <div class="rounded-3xl bg-white text-[#141F19] p-7 shadow-2xl shadow-black/40 lg:rotate-[0.6deg]" aria-hidden="true">
        <div class="flex items-baseline justify-between">
          <p class="text-xs font-bold text-[#5B6960]">11월 쓴 돈</p>
          <p class="text-[11px] text-[#8A948C]">지난달보다 <span class="font-bold text-[#1E5A38]">₩128,400 덜 씀</span></p>
        </div>
        <p class="text-4xl font-black tracking-tight mt-1">₩1,842,300</p>
        <div class="mt-5 border-t border-[#ECEFE8]">
          <div class="flex items-center gap-3 py-3 border-b border-[#F0F3EE]">
            <div class="flex-1 min-w-0">
              <p class="text-[13px] font-bold truncate">택시 · 합정→성수</p>
              <p class="text-[11px] text-[#8A948C]">어제 21:42 · 교통</p>
            </div>
            <strong class="text-[13px]">−₩14,800</strong>
          </div>
          <div class="flex items-center gap-3 py-3 border-b border-[#F0F3EE]">
            <div class="flex-1 min-w-0">
              <p class="text-[13px] font-bold truncate">급여 입금</p>
              <p class="text-[11px] text-[#8A948C]">어제 08:15 · 주계좌</p>
            </div>
            <strong class="text-[13px] text-[#1E5A38]">+₩4,200,000</strong>
          </div>
          <div class="flex items-center gap-3 py-3 border-b border-[#F0F3EE]">
            <div class="flex-1 min-w-0">
              <p class="text-[13px] font-bold truncate">대한항공 도쿄 왕복</p>
              <p class="text-[11px] text-[#8A948C]">11월 3일 · 여행 적립</p>
            </div>
            <strong class="text-[13px]">−₩486,000</strong>
          </div>
          <div class="flex items-center gap-3 py-3">
            <div class="flex-1 min-w-0">
              <p class="text-[13px] font-bold truncate">스트리밍 4건 묶음</p>
              <p class="text-[11px] text-[#8A948C]">다음 갱신 11월 12일</p>
            </div>
            <strong class="text-[13px]">₩43,600/월</strong>
          </div>
        </div>
      </div>
    </main>

    <!-- How it works -->
    <section class="max-w-6xl mx-auto px-6 w-full pb-16">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-10 border-t border-white/10 pt-10">
        <div>
          <h2 class="font-black text-2xl tracking-tight">기록은 가볍게,<br />결산은 자동으로.</h2>
          <p class="text-sm text-[#A9BCAF] leading-relaxed mt-4 max-w-sm">
            자판기 커피 1,500원도 그냥 적으세요. 카테고리 합산, 월별 추이,
            여행별 환산까지 MOA가 알아서 묶어줍니다. 엑셀 정리하던 일요일 오후는 이제 없습니다.
          </p>
          <RouterLink to="/register" class="inline-block mt-6 text-sm font-bold text-[#E9D9A8] hover:text-white transition-colors">
            첫 기록 남기기 →
          </RouterLink>
        </div>
        <div class="rounded-2xl bg-white/[0.04] border border-white/10 p-6 font-mono text-[12px] leading-loose">
          <p class="text-[#7E9387]">11월 6일 · 목요일</p>
          <p><span class="text-[#7E9387]">08:15</span>　급여 +4,200,000</p>
          <p><span class="text-[#7E9387]">12:30</span>　점심 −12,000 <span class="text-[#7E9387]">식비</span></p>
          <p><span class="text-[#7E9387]">19:02</span>　서점 −23,500 <span class="text-[#7E9387]">취미</span></p>
          <p><span class="text-[#7E9387]">21:42</span>　택시 −14,800 <span class="text-[#7E9387]">교통</span></p>
          <p class="text-[#E9D9A8] font-bold pt-2 border-t border-white/10 mt-2">합계 −50,300 · 남음 1,792,000</p>
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
