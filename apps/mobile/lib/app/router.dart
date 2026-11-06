import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../core/api/api_client.dart';
import '../core/theme/colors.dart';
import '../features/onboarding/screens/welcome_screen.dart';
import '../features/onboarding/screens/login_screen.dart';
import '../features/onboarding/screens/register_screen.dart';
import '../features/onboarding/screens/onboarding_screen.dart';
import '../features/home/screens/home_dashboard_screen.dart';
import '../features/assets/screens/assets_screen.dart';
import '../features/transactions/screens/transactions_screen.dart';
import '../features/transactions/screens/add_transaction_screen.dart';
import '../features/budgets/screens/budgets_screen.dart';
import '../features/subscriptions/screens/subscriptions_screen.dart';
import '../features/travel/screens/travel_screen.dart';
import '../features/forex/screens/forex_screen.dart';
import '../features/notifications/screens/notifications_screen.dart';
import '../features/profile/screens/settings_screen.dart';
import '../features/profile/screens/support_screen.dart';

final appRouter = GoRouter(
  initialLocation: '/',
  routes: [
    GoRoute(
      path: '/',
      builder: (context, state) => const SplashScreen(),
    ),
    GoRoute(
      path: '/welcome',
      builder: (context, state) => const WelcomeScreen(),
    ),
    GoRoute(
      path: '/login',
      builder: (context, state) => const LoginScreen(),
    ),
    GoRoute(
      path: '/register',
      builder: (context, state) => const RegisterScreen(),
    ),
    GoRoute(
      path: '/onboarding',
      builder: (context, state) => const OnboardingScreen(),
    ),
    GoRoute(
      path: '/home',
      builder: (context, state) => const HomeDashboardScreen(),
    ),
    GoRoute(
      path: '/assets',
      builder: (context, state) => const AssetsScreen(),
    ),
    GoRoute(
      path: '/transactions',
      builder: (context, state) => const TransactionsScreen(),
    ),
    GoRoute(
      path: '/transactions/new',
      builder: (context, state) => const AddTransactionScreen(),
    ),
    GoRoute(
      path: '/budgets',
      builder: (context, state) => const BudgetsScreen(),
    ),
    GoRoute(
      path: '/subscriptions',
      builder: (context, state) => const SubscriptionsScreen(),
    ),
    GoRoute(
      path: '/travel',
      builder: (context, state) => const TravelScreen(),
    ),
    GoRoute(
      path: '/fx',
      builder: (context, state) => const ForexScreen(),
    ),
    GoRoute(
      path: '/notifications',
      builder: (context, state) => const NotificationsScreen(),
    ),
    GoRoute(
      path: '/settings',
      builder: (context, state) => const SettingsScreen(),
    ),
    GoRoute(
      path: '/support',
      builder: (context, state) => const SupportScreen(),
    ),
  ],
);

class SplashScreen extends StatefulWidget {
  const SplashScreen({super.key});

  @override
  State<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends State<SplashScreen> {
  @override
  void initState() {
    super.initState();
    _checkSession();
  }

  Future<void> _checkSession() async {
    await ApiClient().init();
    await Future.delayed(const Duration(milliseconds: 600));
    if (mounted) {
      if (ApiClient().isAuthenticated) {
        context.go('/home');
      } else {
        context.go('/welcome');
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.pineDeep,
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Container(
              width: 72,
              height: 72,
              decoration: BoxDecoration(
                color: AppColors.accent,
                borderRadius: BorderRadius.circular(22),
              ),
              child: const Center(
                child: Text(
                  '모',
                  style: TextStyle(
                    color: AppColors.pineDeep,
                    fontSize: 36,
                    fontWeight: FontWeight.w900,
                  ),
                ),
              ),
            ),
            const SizedBox(height: 20),
            const Text(
              'MOA',
              style: TextStyle(
                color: Colors.white,
                fontSize: 28,
                fontWeight: FontWeight.w900,
                letterSpacing: 4,
              ),
            ),
            const SizedBox(height: 8),
            const Text(
              '돈을 모으다',
              style: TextStyle(
                color: AppColors.darkTextTertiary,
                fontSize: 13,
                fontWeight: FontWeight.w600,
                letterSpacing: 1,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

