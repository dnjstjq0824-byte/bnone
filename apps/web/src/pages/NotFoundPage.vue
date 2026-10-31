<template>
  <div class="min-h-screen bg-[#F9FAFB] flex items-center justify-center p-4">
    <div class="max-w-md w-full bg-white rounded-3xl p-8 text-center card-shadow border border-[#F2F4F6] space-y-4">
      <div class="w-16 h-16 rounded-3xl bg-[#F2F4F6] text-[#4E5968] flex items-center justify-center text-3xl mx-auto font-black">
        404
      </div>
      <div>
        <h1 class="text-xl font-bold text-[#191F28] tracking-tight">페이지를 찾을 수 없습니다</h1>
        <p class="text-xs text-[#8B95A1] mt-1.5 leading-relaxed">
          요청하신 페이지가 삭제되었거나 잘못된 경로입니다.<br />
          입력하신 주소가 정확한지 다시 한번 확인해주세요.
        </p>
      </div>
      <div class="pt-4 flex flex-col gap-2">
        <RouterLink
          to="/app"
          class="w-full py-3.5 bg-[#3182F6] hover:bg-[#1B64DA] text-white rounded-xl font-bold text-sm shadow-sm transition-colors block text-center"
        >
          홈 대시보드로 이동
        </RouterLink>
        <RouterLink
          to="/"
          class="w-full py-3 bg-[#F2F4F6] hover:bg-[#E5E8EB] text-[#4E5968] rounded-xl font-semibold text-xs transition-colors block text-center"
        >
          랜딩 페이지로 이동
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router';
</script>
