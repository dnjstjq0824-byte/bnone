import 'package:flutter/material.dart';

/// Toss Design System (TDS) Color Tokens for BN ONE
class AppColors {
  // Brand colors — Toss Blue
  static const Color brandPrimary = Color(0xFF3182F6); // Toss Blue
  static const Color primaryDark = Color(0xFF1B64DA);  // Toss Blue Dark / Pressed
  static const Color primaryLight = Color(0xFFE8F3FF); // Toss Blue Light Tint
  static const Color accent = Color(0xFF3182F6);

  // Surfaces & Backgrounds — Toss Cool Canvas
  static const Color background = Color(0xFFF2F4F6);  // Toss Background Canvas
  static const Color surface = Color(0xFFFFFFFF);     // Toss Card Surface
  static const Color surfaceElevated = Color(0xFFFFFFFF);
  static const Color grayBg = Color(0xFFF9FAFB);

  // Typography — Toss High Readability Hierarchy
  static const Color textPrimary = Color(0xFF191F28);   // Toss Black
  static const Color textSecondary = Color(0xFF4E5968); // Toss Gray 600
  static const Color textTertiary = Color(0xFF8B95A1);  // Toss Gray 400
  static const Color border = Color(0xFFE5E8EB);        // Faint border
  static const Color divider = Color(0xFFF2F4F6);

  // Semantic Status
  static const Color success = Color(0xFF3182F6); // Toss uses Blue for income/positives
  static const Color green = Color(0xFF00B06B);   // Toss Mint Green
  static const Color warning = Color(0xFFFF8D00); // Toss Warning Amber
  static const Color danger = Color(0xFFF04452);  // Toss Coral Red

  // Dark Mode colors
  static const Color darkBackground = Color(0xFF101012);
  static const Color darkSurface = Color(0xFF1C1C22);
  static const Color darkTextPrimary = Color(0xFFF9FAFB);
  static const Color darkTextSecondary = Color(0xFFA6AFB9);
  static const Color darkTextTertiary = Color(0xFF6B7684);
  static const Color darkBorder = Color(0xFF2B2B33);
}
