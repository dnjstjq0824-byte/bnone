<template>
  <div class="min-h-screen bg-slate-900 text-white flex flex-col justify-between">
    <!-- Navbar -->
    <header class="max-w-7xl mx-auto w-full px-6 h-20 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-xl bg-brand-600 flex items-center justify-center font-black text-xl shadow-lg shadow-brand-500/30">
          BN
        </div>
        <span class="font-extrabold text-xl tracking-tight">BN ONE</span>
      </div>
      <div class="flex items-center gap-4">
        <RouterLink to="/login" class="text-sm font-semibold text-slate-300 hover:text-white px-3 py-2 transition-colors">
          로그인
        </RouterLink>
        <RouterLink to="/register" class="text-sm font-bold bg-brand-600 hover:bg-brand-500 px-4 py-2 rounded-xl transition-all shadow-md shadow-brand-600/30">
          시작하기
        </RouterLink>
      </div>
    </header>

    <!-- Hero Section -->
    <main class="max-w-5xl mx-auto px-6 py-16 text-center">
      <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-slate-800/80 border border-slate-700 text-xs font-medium text-brand-400 mb-6">
        <span>🌐 글로벌 금융 라이프스타일 슈퍼앱</span>
      </div>

      <h1 class="text-4xl md:text-6xl font-black tracking-tight leading-tight md:leading-tight mb-6">
        Your Money.<br />
        <span class="text-transparent bg-clip-text bg-gradient-to-r from-blue-400 via-teal-300 to-indigo-400">
          Your World.
        </span>
      </h1>

      <p class="text-base md:text-xl text-slate-300 max-w-2xl mx-auto mb-10 leading-relaxed">
        매일 쓰는 돈부터 세계 여행에서 쓰는 돈까지, 하나의 앱에서 관리한다.<br />
        다중 통화 자산 장부, 정기 구독, 예산과 여행 경비를 가장 정확하고 안전하게 기록하세요.
      </p>

      <div class="flex flex-col sm:flex-row items-center justify-center gap-4 mb-16">
        <RouterLink to="/register" class="w-full sm:w-auto px-8 py-3.5 bg-brand-600 hover:bg-brand-500 rounded-xl font-bold text-base shadow-lg shadow-brand-500/25 transition-all">
          무료로 자산 관리 시작하기
        </RouterLink>
        <RouterLink to="/login" class="w-full sm:w-auto px-8 py-3.5 bg-slate-800 hover:bg-slate-700 border border-slate-700 rounded-xl font-semibold text-base transition-all">
          기존 계정으로 로그인
        </RouterLink>
      </div>

      <!-- Features Grid -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 text-left">
        <div class="p-6 rounded-2xl bg-slate-800/60 border border-slate-700/60">
          <div class="text-3xl mb-3">🏛️</div>
          <h3 class="font-bold text-lg mb-2">다중 수동 금융계정</h3>
          <p class="text-slate-400 text-xs leading-relaxed">
            현금, 은행계좌형 장부, 외화 및 기타 자산을 사용자가 직접 신뢰성 있게 기록하고 순자산을 추적합니다.
          </p>
        </div>

        <div class="p-6 rounded-2xl bg-slate-800/60 border border-slate-700/60">
          <div class="text-3xl mb-3">✈️</div>
          <h3 class="font-bold text-lg mb-2">여행 예산 & 경비</h3>
          <p class="text-slate-400 text-xs leading-relaxed">
            여행 일정별 예산 설정과 원본 통화 경비 기록, 공인 참고 환율 기반의 실시간 환산 보고서를 제공합니다.
          </p>
        </div>

        <div class="p-6 rounded-2xl bg-slate-800/60 border border-slate-700/60">
          <div class="text-3xl mb-3">🔄</div>
          <h3 class="font-bold text-lg mb-2">정기 구독 & 알림</h3>
          <p class="text-slate-400 text-xs leading-relaxed">
            넷플릭스, 유튜브 등 분산된 구독 서비스의 갱신 주기와 월간/연간 예상 비용을 통합 관리합니다.
          </p>
        </div>
      </div>
    </main>

    <!-- Footer with mandatory legal notice -->
    <footer class="border-t border-slate-800 py-8 px-6 text-center text-xs text-slate-500">
      <div class="max-w-4xl mx-auto space-y-2">
        <p>Booking Network co., Ltd — BN ONE v1.0 Global Financial Lifestyle</p>
        <p class="text-[11px] text-slate-500">
          * BN ONE v1.0은 사용자가 직접 입력한 금융정보를 기반으로 자산과 지출을 관리하는 소프트웨어이며, 실제 금융기관 계좌 연결 및 금융거래를 직접 실행하지 않습니다.
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router';
</script>
