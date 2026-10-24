import 'package:flutter/material.dart';
import '../core/theme/app_theme.dart';
import 'router.dart';

class BnOneApp extends StatelessWidget {
  const BnOneApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      title: 'BN ONE',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.lightTheme,
      darkTheme: AppTheme.darkTheme,
      themeMode: ThemeMode.system,
      routerConfig: appRouter,
    );
  }
}

