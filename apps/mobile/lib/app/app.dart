import 'package:flutter/material.dart';
import '../core/theme/app_theme.dart';
import 'router.dart';

class MoaApp extends StatelessWidget {
  const MoaApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      title: 'MOA',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.lightTheme,
      darkTheme: AppTheme.darkTheme,
      themeMode: ThemeMode.system,
      routerConfig: appRouter,
    );
  }
}

