# 05. VUE 3 USER WEB ARCHITECTURE & SPECIFICATION

**Document ID:** BN-ONE-WEB-005  
**Version:** 1.0.0  
**Effective Date:** 2026-10-24  

---

## 1. Architectural Stack: Vue 3 + TypeScript + Pinia + Tailwind CSS

BN ONE User Web (`apps/web`) delivers a fluid, responsive financial management portal for desktop and mobile web browsers.

### 1.1 Technical Stack
- **Framework:** Vue 3 Composition API (`<script setup lang="ts">`).
- **Build Tool:** Vite 5.x for sub-second HMR and optimized production bundling.
- **State Store:** Pinia stores (`authStore`, `financeStore`).
- **Routing:** Vue Router 4 with navigation guards checking active JWT tokens.
- **Styling:** Tailwind CSS with custom theme extensions matching the BN ONE design token palette.

---

## 2. Directory Structure

```text
apps/web/src/
  api/
    client.ts             <-- Axios-style Fetch client with auto Authorization header
  components/
    AppLayout.vue         <-- Responsive shell with sidebar, net worth widget, topbar
  stores/
    auth.ts               <-- User session, login, register, token persistence
    finance.ts            <-- Accounts, transactions, budgets, subscriptions state
  views/
    LandingView.vue       <-- Public SEO-friendly product showcase & feature overview
    LoginView.vue         <-- Secure email/password login
    RegisterView.vue      <-- User registration form
    OnboardingView.vue    <-- Currency, locale, and initial account setup
    DashboardView.vue     <-- Financial health dashboard with cashflow charts
    AssetsView.vue        <-- Asset & liability ledger breakdown
    TransactionsView.vue  <-- Filterable transaction ledger
    NewTransactionView.vue<-- Fast transaction entry modal / page
    AnalyticsView.vue     <-- Category spend trends and burn-rate charts
    BudgetsView.vue       <-- Monthly category budget progress bars
    SubscriptionsView.vue <-- Active subscriptions and renewal calendar
    TripsView.vue         <-- Multi-currency travel wallets and budgets
    ForexView.vue         <-- Reference FX converter & live rate board
    NotificationsView.vue <-- Activity notices and renewal reminders
    SettingsView.vue      <-- Base currency, locale, CSV export, account deletion
    SupportView.vue       <-- Help center and customer support inquiry form
```

---

## 3. Security & Anti-Leakage Boundaries

1. **Zero Financial Data in Public Pre-renders:**
   - Landing view is statically indexable for search engines.
   - All authenticated financial views (`/app/*`) require an active, validated Bearer token.
   - User financial figures are never injected into server-rendered HTML sent to search spiders.
2. **Strict Sanitization:**
   - All user-entered memo strings, category labels, and account names are automatically escaped to prevent Cross-Site Scripting (XSS).
3. **Session Expiry Handling:**
   - HTTP 401 Unauthorized responses trigger automatic session termination and redirection to `/login` with an informational notice.

