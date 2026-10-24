import { createRouter, createWebHistory } from 'vue-router';
import { useAdminAuthStore } from '../stores/adminAuth';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/admin/login',
      name: 'AdminLogin',
      component: () => import('../pages/AdminLoginPage.vue'),
    },
    {
      path: '/admin',
      component: () => import('../layouts/AdminLayout.vue'),
      meta: { requiresAdmin: true },
      children: [
        {
          path: '',
          redirect: '/admin/dashboard',
        },
        {
          path: 'dashboard',
          name: 'AdminDashboard',
          component: () => import('../pages/AdminDashboardPage.vue'),
        },
        {
          path: 'users',
          name: 'AdminUsers',
          component: () => import('../pages/AdminUsersPage.vue'),
        },
        {
          path: 'support',
          name: 'AdminSupport',
          component: () => import('../pages/AdminSupportPage.vue'),
        },
        {
          path: 'fx',
          name: 'AdminFX',
          component: () => import('../pages/AdminFXPage.vue'),
        },
        {
          path: 'privacy',
          name: 'AdminPrivacy',
          component: () => import('../pages/AdminPrivacyPage.vue'),
        },
        {
          path: 'audit-logs',
          name: 'AdminAuditLogs',
          component: () => import('../pages/AdminAuditLogsPage.vue'),
        },
        {
          path: 'feature-flags',
          name: 'AdminFeatureFlags',
          component: () => import('../pages/AdminFeatureFlagsPage.vue'),
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/admin/dashboard',
    },
  ],
});

router.beforeEach((to, from, next) => {
  const adminAuth = useAdminAuthStore();
  if (to.meta.requiresAdmin && !adminAuth.token) {
    next({ name: 'AdminLogin' });
  } else {
    next();
  }
});

export default router;
