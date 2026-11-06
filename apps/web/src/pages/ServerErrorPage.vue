<template>
  <div class="min-h-screen bg-[#F6F8F4] flex items-center justify-center p-4">
    <div class="max-w-md w-full bg-white rounded-3xl p-8 text-center card-shadow border border-[#ECEFE8] space-y-4">
      <div class="w-16 h-16 rounded-3xl bg-[#FBEAEA] text-[#D64545] flex items-center justify-center text-3xl mx-auto font-black">
        500
      </div>
      <div>
        <h1 class="text-xl font-bold text-[#141F19] tracking-tight">일시적인 오류가 발생했습니다</h1>
        <p class="text-xs text-[#8A948C] mt-1.5 leading-relaxed">
          서버 처리 중 문제가 발생했습니다. 잠시 후 다시 시도해주세요.<br />
          문제가 계속되면 아래 요청 ID와 함께 고객센터에 문의해주세요.
        </p>
        <p v-if="requestId" class="mt-3 inline-block px-3 py-1.5 rounded-lg bg-[#ECEFE8] text-[11px] font-mono text-[#41524A]">
          요청 ID: {{ requestId }}
        </p>
      </div>
      <div class="pt-4 flex flex-col gap-2">
        <button
          @click="retry"
          class="w-full py-3.5 bg-[#1E5A38] hover:bg-[#17452C] text-white rounded-xl font-bold text-sm shadow-sm transition-colors"
        >
          다시 시도
        </button>
        <RouterLink
          to="/app"
          class="w-full py-3 bg-[#ECEFE8] hover:bg-[#E2E7DE] text-[#41524A] rounded-xl font-semibold text-xs transition-colors block text-center"
        >
          홈 대시보드로 이동
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { RouterLink, useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();
const requestId = ref(typeof route.query.request_id === 'string' ? route.query.request_id : '');

function retry() {
  const from = typeof route.query.from === 'string' && route.query.from ? route.query.from : '/app';
  router.push(from);
}
</script>
