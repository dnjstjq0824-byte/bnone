# 04. FLUTTER MOBILE ARCHITECTURE & SPECIFICATION

**Document ID:** BN-ONE-MOB-004  
**Version:** 1.0.0  
**Effective Date:** 2026-10-24  

---

## 1. Architectural Pattern: Feature-First with Riverpod

BN ONE Mobile is engineered with a modular, **feature-first** directory structure adhering strictly to Flutter and Dart idioms.

```text
apps/mobile/lib/
  app/
    app.dart               <-- MaterialApp.router, AppTheme, Localization
    router.dart            <-- GoRouter route table with auth redirection guard
  core/
    api/
      api_client.dart      <-- HTTP Client with Bearer token interceptor
    theme/
      colors.dart          <-- Color palette design tokens
      app_theme.dart       <-- ThemeData (Material 3, CardThemeData, Typography)
    utils/
      formatters.dart      <-- Minor-to-major currency formatters & masking
  features/
    onboarding/screens/    <-- Welcome, Login, Register, Onboarding Wizard
    home/screens/          <-- Home Dashboard (Net worth, Cashflow, Quick Actions)
    assets/screens/        <-- Asset & Liability management, manual accounts
    transactions/screens/  <-- Transaction timeline, filtering, new transaction
    budgets/screens/       <-- Category budgeting & burn-rates
    subscriptions/screens/ <-- Recurring subscriptions & renewal tracking
    travel/screens/        <-- Trip planning & multi-currency travel logs
    forex/screens/         <-- Live FX rates & reference currency converter
    notifications/screens/ <-- Push / In-app notification center
    profile/screens/       <-- User settings, data export, support tickets
```

---

## 2. State Management & Dependency Injection

- **Riverpod 2.x:** Stateless UI components consume providers via `ConsumerWidget` or `ConsumerStatefulWidget`.
- **Auth State Provider:** Monitors authentication status (token presence, user profile). When unauthenticated, GoRouter intercepts route transitions and routes the user to `/welcome`.
- **Finance State Provider:** Manages local cache of accounts, transactions, and net worth calculations.
- **Privacy Masking Provider:** Manages global state of monetary obfuscation (`bool _maskAmounts`).

---

## 3. GoRouter Route Hierarchy

```dart
final appRouter = GoRouter(
  initialLocation: '/welcome',
  routes: [
    GoRoute(path: '/welcome', builder: (c, s) => const WelcomeScreen()),
    GoRoute(path: '/login', builder: (c, s) => const LoginScreen()),
    GoRoute(path: '/register', builder: (c, s) => const RegisterScreen()),
    GoRoute(path: '/onboarding', builder: (c, s) => const OnboardingScreen()),
    
    // Bottom Tab Shell
    GoRoute(path: '/home', builder: (c, s) => const HomeDashboardScreen()),
    GoRoute(path: '/assets', builder: (c, s) => const AssetsScreen()),
    GoRoute(path: '/transactions', builder: (c, s) => const TransactionsScreen()),
    GoRoute(path: '/transactions/new', builder: (c, s) => const AddTransactionScreen()),
    GoRoute(path: '/budgets', builder: (c, s) => const BudgetsScreen()),
    GoRoute(path: '/subscriptions', builder: (c, s) => const SubscriptionsScreen()),
    GoRoute(path: '/travel', builder: (c, s) => const TravelScreen()),
    GoRoute(path: '/forex', builder: (c, s) => const ForexScreen()),
    GoRoute(path: '/notifications', builder: (c, s) => const NotificationsScreen()),
    GoRoute(path: '/settings', builder: (c, s) => const SettingsScreen()),
    GoRoute(path: '/support', builder: (c, s) => const SupportScreen()),
  ],
);
```

---

## 4. Hardware & Platform Integration

- **Android SDK:** Min SDK 21 (Android 5.0 Lollipop), Target SDK 34 (Android 14).
- **iOS:** Minimum iOS deployment target 13.0.
- **Secure Storage:** Uses Android Keystore and iOS Keychain via `flutter_secure_storage` for refresh token protection.
- **Offline Cache:** Cached snapshot of account balances and recent transactions stored locally in encrypted Key-Value storage.

