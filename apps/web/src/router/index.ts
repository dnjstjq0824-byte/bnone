import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Landing',
      component: () => import('../pages/LandingPage.vue'),
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('../pages/LoginPage.vue'),
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('../pages/RegisterPage.vue'),
    },
    {
      path: '/onboarding',
      name: 'Onboarding',
      component: () => import('../pages/OnboardingPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/app',
      component: () => import('../layouts/AppLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'Dashboard',
          component: () => import('../pages/DashboardPage.vue'),
        },
        {
          path: 'assets',
          name: 'Assets',
          component: () => import('../pages/AssetsPage.vue'),
        },
        {
          path: 'transactions',
          name: 'Transactions',
          component: () => import('../pages/TransactionsPage.vue'),
        },
        {
          path: 'transactions/new',
          name: 'NewTransaction',
          component: () => import('../pages/NewTransactionPage.vue'),
        },
        {
          path: 'analytics',
          name: 'Analytics',
          component: () => import('../pages/AnalyticsPage.vue'),
        },
        {
          path: 'budgets',
          name: 'Budgets',
          component: () => import('../pages/BudgetsPage.vue'),
        },
        {
          path: 'subscriptions',
          name: 'Subscriptions',
          component: () => import('../pages/SubscriptionsPage.vue'),
        },
        {
          path: 'trips',
          name: 'Trips',
          component: () => import('../pages/TripsPage.vue'),
        },
        {
          path: 'fx',
          name: 'Forex',
          component: () => import('../pages/ForexPage.vue'),
        },
        {
          path: 'notifications',
          name: 'Notifications',
          component: () => import('../pages/NotificationsPage.vue'),
        },
        {
          path: 'settings',
          name: 'Settings',
          component: () => import('../pages/SettingsPage.vue'),
        },
        {
          path: 'support',
          name: 'Support',
          component: () => import('../pages/SupportPage.vue'),
        },
      ],
    },
  ],
});

router.beforeEach(async (to, from, next) => {
  const auth = useAuthStore();
  if (to.meta.requiresAuth) {
    if (!auth.token) {
      next({ name: 'Login' });
      return;
    }
    if (!auth.user) {
      await auth.fetchMe();
      if (!auth.isAuthenticated) {
        next({ name: 'Login' });
        return;
      }
    }
  }
  next();
});

export default router;

