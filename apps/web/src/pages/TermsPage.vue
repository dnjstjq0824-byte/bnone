<template>
  <div class="min-h-screen bg-[#F6F8F4] py-10 px-4">
    <div class="max-w-2xl mx-auto bg-white rounded-3xl p-8 card-shadow border border-[#ECEFE8]">
      <p class="text-[11px] font-bold text-[#1E5A38]">MOA v1.0 · 정책 버전 1.0 (시행 2026-10-24)</p>
      <h1 class="text-2xl font-bold text-[#141F19] tracking-tight mt-2">서비스 이용약관</h1>
      <p class="text-xs text-[#8A948C] mt-1.5 leading-relaxed">
        본 약관은 부킹네트워크(이하 "회사")가 제공하는 MOA 서비스의 이용 조건을 정합니다.
        회원은 가입 시 본 약관 전문을 확인하고 동의해야 하며, 동의 내역은 버전과 함께 보관됩니다.
      </p>

      <div class="mt-6 space-y-6 text-[13px] leading-relaxed text-[#41524A]">
        <section>
          <h2 class="text-sm font-bold text-[#141F19] mb-1.5">제1조 (서비스의 성격과 범위)</h2>
          <p>MOA은 개인 가계·자산·구독·여행 경비를 사용자가 직접 입력·관리하는 개인 금융 기록 도구입니다. 회사는 실제 자금 이체, 오픈뱅킹 연동, 카드 발급, 투자 자문을 제공하지 않으며, 사용자의 은행 로그인 정보를 수집·보관하지 않고 고객 자금을 수탁·보관·보험하지 않습니다. 화면에 표시된 자산과 잔액은 사용자가 직접 입력한 기록을 기준으로 합니다.</p>
        </section>
        <section>
          <h2 class="text-sm font-bold text-[#141F19] mb-1.5">제2조 (계정과 보안)</h2>
          <p>회원은 이메일과 비밀번호로 가입하며, 비밀번호는 8자 이상(영문·숫자 포함)이어야 합니다. 회원은 계정과 비밀번호를 안전하게 관리할 책임이 있으며, 도용이 의심되면 즉시 비밀번호를 변경하고 회사에 알려야 합니다. 회사는 비정상 로그인 시도가 감지되면 해당 계정의 접근을 일시 제한할 수 있습니다.</p>
        </section>
        <section>
          <h2 class="text-sm font-bold text-[#141F19] mb-1.5">제3조 (참고 정보 고지)</h2>
          <p>표시 환율은 참고용이며 실제 환전·결제 환율과 다를 수 있습니다. MOA의 거래 기록은 실제 금융기관 거래를 실행하지 않습니다. 구독 갱신 알림은 참고 정보이며, 실제 결제일 및 청구 금액은 해당 서비스 제공자의 정책에 따릅니다.</p>
        </section>
        <section>
          <h2 class="text-sm font-bold text-[#141F19] mb-1.5">제4조 (데이터 내보내기와 탈퇴)</h2>
          <p>회원은 설정 화면에서 언제든지 본인 데이터를 CSV 형식으로 내려받을 수 있습니다. 탈퇴 요청 시 회사는 관련 법령이 허용하는 범위에서 회원 데이터를 영구 파기하며, 파기 완료 전까지 회원은 요청을 철회할 수 있습니다.</p>
        </section>
        <section>
          <h2 class="text-sm font-bold text-[#141F19] mb-1.5">제5조 (금지 행위)</h2>
          <p>회원은 타인의 계정 도용, 서비스의 비정상적 과부하 유발, 보안 조치 우회, 법령 위반 목적의 이용을 해서는 안 됩니다. 회사는 약관 위반이 확인되면 사전 통지 후 이용을 제한하거나 계약을 해지할 수 있습니다.</p>
        </section>
        <section>
          <h2 class="text-sm font-bold text-[#141F19] mb-1.5">제6조 (책임 제한과 분쟁 해결)</h2>
          <p>회사는 사용자 입력 기록의 정확성, 제3자 서비스의 과금·정책 변경, 천재지변 등 불가항력으로 인한 손해에 대해 책임지지 않습니다. 분쟁은 대한민국 법을 준거법으로 하며, 관할 법원은 민사소송법에 따릅니다. 문의: support@bookingnetwork.com</p>
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
