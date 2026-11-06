<template>
  <div class="fixed bottom-6 left-1/2 -translate-x-1/2 z-50 flex flex-col items-center gap-2 pointer-events-none px-4 w-full max-w-md">
    <div
      v-for="toast in toasts"
      :key="toast.id"
      role="status"
      class="pointer-events-auto w-full px-4 py-3 rounded-2xl text-[13px] font-semibold shadow-lg border transition-all"
      :class="{
        'bg-[#141F19] text-white border-[#141F19]': toast.kind === 'success',
        'bg-[#FBEAEA] text-[#BE3A3A] border-[#F2C9C9]': toast.kind === 'error',
        'bg-white text-[#41524A] border-[#E2E7DE]': toast.kind === 'info',
      }"
    >
      {{ toast.message }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { useToast } from '../composables/useToast';

const { toasts } = useToast();
</script>
