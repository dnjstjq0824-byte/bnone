<template>
  <div class="min-h-screen bg-[#F6F8F4] py-10 px-4">
    <div class="max-w-2xl mx-auto bg-white rounded-3xl p-8 card-shadow border border-[#ECEFE8]">
      <p class="text-[11px] font-bold text-[#1E5A38]">MOA v1.0 · 정책 버전 1.0 (시행 2026-10-24)</p>
      <h1 class="text-2xl font-bold text-[#141F19] tracking-tight mt-2">개인정보처리방침</h1>
      <p class="text-xs text-[#8A948C] mt-1.5 leading-relaxed">
        부킹네트워크(이하 "회사")는 회원의 개인정보를 안전하게 보호하며, 본 방침에 고지한 목적과 범위에서만 처리합니다.
      </p>

      <div class="mt-6 space-y-6 text-[13px] leading-relaxed text-[#41524A]">
        <section>
          <h2 class="text-sm font-bold text-[#141F19] mb-1.5">1. 수집 항목과 목적</h2>
          <ul class="list-disc pl-5 space-y-1">
            <li>계정 식별: 이메일, 비밀번호(암호화 저장), 표시 이름, 국가·언어·통화·시간대 설정</li>
            <li>서비스 제공: 사용자가 입력한 계좌·자산·부채·거래·예산·구독·여행 기록</li>
            <li>보안·품질: 로그인 기록, 세션 정보, 오류 요청 ID, 알림 수신 토큰</li>
          </ul>
          <p class="mt-1.5">회사는 은행 로그인 정보, 카드 번호, 계좌 비밀번호를 절대 수집하지 않습니다.</p>
        </section>
        <section>
          <h2 class="text-sm font-bold text-[#141F19] mb-1.5">2. 보유와 파기</h2>
          <p>개인정보는 회원 탈퇴 시 지체 없이 파기합니다. 단, 분쟁 대응과 법령 준수를 위해 필요한 최소 기록(동의 이력, 감사 로그)은 관련 법령이 정한 기간 동안 분리 보관 후 파기합니다. 탈퇴 심사 중에는 요청일로부터 30일 이내에 파기가 완료됩니다.</p>
        </section>
        <section>
          <h2 class="text-sm font-bold text-[#141F19] mb-1.5">3. 회원의 권리</h2>
          <p>회원은 설정 화면에서 본인 정보를 조회·수정하고, CSV 내보내기로 데이터를 내려받으며, 탈퇴 요청으로 삭제를 요구할 수 있습니다. 처리 결과에 이의가 있으면 support@bookingnetwork.com 으로 문의할 수 있습니다.</p>
        </section>
        <section>
          <h2 class="text-sm font-bold text-[#141F19] mb-1.5">4. 안전 조치</h2>
          <p>비밀번호는 bcrypt로 암호화하여 저장하고, 인증 토큰은 서명된 JWT와 해시된 리프레시 토큰으로 관리합니다. 전송 구간은 TLS로 암호화하며, 관리자 접근은 역할 기반으로 통제되고 모든 조치가 감사 로그에 기록됩니다.</p>
        </section>
        <section>
          <h2 class="text-sm font-bold text-[#141F19] mb-1.5">5. 제3자 제공과 위탁</h2>
          <p>회사는 회원의 개인정보를 제3자에게 판매·제공하지 않습니다. 인프라 운영을 위해 클라우드 호스팅 사업자에게 처리를 위탁할 수 있으며, 위탁 시 본 방침과 동등한 수준의 보호 의무를 부과합니다.</p>
        </section>
      </div>

      <div class="pt-6 flex gap-2">
        <RouterLink
          to="/register"
          class="flex-1 py-3.5 bg-[#1E5A38] hover:bg-[#17452C] text-white rounded-xl font-bold text-sm shadow-sm transition-colors block text-center"
        >
          동의하고 가입하기
        </RouterLink>
        <RouterLink
          to="/"
          class="flex-1 py-3.5 bg-[#ECEFE8] hover:bg-[#E2E7DE] text-[#41524A] rounded-xl font-semibold text-sm transition-colors block text-center"
        >
          처음으로
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router';
</script>
