import { defineStore } from 'pinia';
import { ref } from 'vue';
import { adminFetch } from '../api/client';

export const useAdminAuthStore = defineStore('adminAuth', () => {
  const token = ref<string | null>(localStorage.getItem('bn_admin_token'));
  const role = ref<string | null>(localStorage.getItem('bn_admin_role'));
  const adminEmail = ref<string | null>(localStorage.getItem('bn_admin_email'));
  const fullName = ref<string | null>(localStorage.getItem('bn_admin_name'));
  const isAuthenticated = ref<boolean>(!!token.value);
  const isLoading = ref<boolean>(false);

  async function login(email: string, pass: string) {
    isLoading.value = true;
    try {
      const data = await adminFetch('/auth/admin/login', {
        method: 'POST',
        body: JSON.stringify({ email, password: pass }),
      });

      token.value = data.access_token;
      role.value = data.role;
      adminEmail.value = data.email;
      fullName.value = data.full_name;
      isAuthenticated.value = true;

      localStorage.setItem('bn_admin_token', data.access_token);
      localStorage.setItem('bn_admin_role', data.role);
      localStorage.setItem('bn_admin_email', data.email);
      localStorage.setItem('bn_admin_name', data.full_name);

      return data;
    } finally {
      isLoading.value = false;
    }
  }

  function logout() {
    token.value = null;
    role.value = null;
    adminEmail.value = null;
    fullName.value = null;
    isAuthenticated.value = false;

    localStorage.removeItem('bn_admin_token');
    localStorage.removeItem('bn_admin_role');
    localStorage.removeItem('bn_admin_email');
    localStorage.removeItem('bn_admin_name');
  }

  return {
    token,
    role,
    adminEmail,
    fullName,
    isAuthenticated,
    isLoading,
    login,
    logout,
  };
});

