// Shared toast state. Import useToast() in any page, render <ToastHost /> once in App.vue.
import { reactive } from 'vue';

export interface ToastItem {
  id: number;
  kind: 'success' | 'error' | 'info';
  message: string;
}

const toasts = reactive<ToastItem[]>([]);
let nextId = 1;

function push(kind: ToastItem['kind'], message: string) {
  const id = nextId++;
  toasts.push({ id, kind, message });
  window.setTimeout(() => {
    const idx = toasts.findIndex((t) => t.id === id);
    if (idx >= 0) toasts.splice(idx, 1);
  }, 4200);
}

export function useToast() {
  return {
    toasts,
    success: (message: string) => push('success', message),
    error: (message: string) => push('error', message),
    info: (message: string) => push('info', message),
  };
}
