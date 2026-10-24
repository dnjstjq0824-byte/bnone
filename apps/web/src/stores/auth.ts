import { defineStore } from 'pinia';
import { ref } from 'vue';
import { apiFetch } from '../api/client';

export interface User {
  id: string;
  email: string;
  display_name: string;
  country_code: string;
  locale: string;
  base_currency: string;
  timezone: string;
  status: string;
  onboarding_completed_at?: string;
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null);
  const token = ref<string | null>(localStorage.getItem('bn_token'));
  const refreshToken = ref<string | null>(localStorage.getItem('bn_refresh_token'));
  const isAuthenticated = ref<boolean>(!!token.value);
  const isLoading = ref<boolean>(false);

  async function login(email: string, pass: string) {
    isLoading.value = true;
    try {
      const data = await apiFetch('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password: pass }),
      });
      token.value = data.access_token;
      refreshToken.value = data.refresh_token;
      user.value = data.user;
      isAuthenticated.value = true;

      localStorage.setItem('bn_token', data.access_token);
      localStorage.setItem('bn_refresh_token', data.refresh_token);
      return data;
    } finally {
      isLoading.value = false;
    }
  }

  async function register(params: {
    email: string;
    password: string;
    display_name?: string;
    base_currency?: string;
    locale?: string;
    country_code?: string;
  }) {
    isLoading.value = true;
    try {
      const data = await apiFetch('/auth/register', {
        method: 'POST',
        body: JSON.stringify({ ...params, required_terms_accepted: true }),
      });
      token.value = data.access_token;
      refreshToken.value = data.refresh_token;
      user.value = data.user;
      isAuthenticated.value = true;

      localStorage.setItem('bn_token', data.access_token);
      localStorage.setItem('bn_refresh_token', data.refresh_token);
      return data;
    } finally {
      isLoading.value = false;
    }
  }

  async function fetchMe() {
    if (!token.value) return null;
    try {
      const data = await apiFetch<User>('/me');
      user.value = data;
      isAuthenticated.value = true;
      return data;
    } catch (e) {
      logout();
      return null;
    }
  }

  function logout() {
    if (token.value) {
      apiFetch('/auth/logout', { method: 'POST' }).catch(() => {});
    }
    user.value = null;
    token.value = null;
    refreshToken.value = null;
    isAuthenticated.value = false;
    localStorage.removeItem('bn_token');
    localStorage.removeItem('bn_refresh_token');
  }

  return {
    user,
    token,
    isAuthenticated,
    isLoading,
    login,
    register,
    fetchMe,
    logout,
  };
});

