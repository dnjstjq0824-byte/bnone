import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../core/api/api_client.dart';
import '../../../core/theme/colors.dart';

class OnboardingScreen extends StatefulWidget {
  const OnboardingScreen({super.key});

  @override
  State<OnboardingScreen> createState() => _OnboardingScreenState();
}

class _OnboardingScreenState extends State<OnboardingScreen> {
  int _step = 1;
  String _currency = 'KRW';
  final _accountNameController = TextEditingController(text: '주거래 급여계좌');
  final _balanceController = TextEditingController(text: '1000000');
  bool _isLoading = false;

  Future<void> _submitFirstAccount() async {
    setState(() => _isLoading = true);
    try {
      // 1. Update user currency
      await ApiClient().patch('/me', {'base_currency': _currency});

      // 2. Create first manual account
      await ApiClient().post('/accounts', {
        'name': _accountNameController.text.trim(),
        'account_type': 'MANUAL_BANK',
        'currency_code': _currency,
        'opening_balance_minor': int.tryParse(_balanceController.text.trim()) ?? 0,
        'opening_balance_date': DateTime.now().toIso8601String().split('T')[0],
      });

      setState(() => _step = 2);
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(e.toString())));
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  Future<void> _finishOnboarding() async {
    await ApiClient().patch('/me', {'complete_onboarding': true});
    if (mounted) {
      context.go('/home');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(24.0),
          child: _step == 1
              ? Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      '초기 설정 (1/2)',
                      style: TextStyle(color: AppColors.brandPrimary, fontWeight: FontWeight.bold),
                    ),
                    const SizedBox(height: 8),
                    const Text(
                      '첫 번째 금융계정 등록',
                      style: TextStyle(fontSize: 26, fontWeight: FontWeight.w900),
                    ),
                    const SizedBox(height: 8),
                    const Text(
                      '자주 쓰는 수동 장부 계좌와 기준 통화를 설정하세요.',
                      style: TextStyle(color: AppColors.textSecondary, fontSize: 13),
                    ),
                    const SizedBox(height: 32),
                    TextField(
                      controller: _accountNameController,
                      decoration: InputDecoration(
                        labelText: '계정 이름',
                        border: OutlineInputBorder(borderRadius: BorderRadius.circular(14)),
                      ),
                    ),
                    const SizedBox(height: 16),
                    DropdownButtonFormField<String>(
                      value: _currency,
                      decoration: InputDecoration(
                        labelText: '기준 통화',
                        border: OutlineInputBorder(borderRadius: BorderRadius.circular(14)),
                      ),
                      items: const [
                        DropdownMenuItem(value: 'KRW', child: Text('KRW — 대한민국 원 (₩)')),
                        DropdownMenuItem(value: 'USD', child: Text('USD — 미국 달러 (\$)')),
                        DropdownMenuItem(value: 'EUR', child: Text('EUR — 유로 (€)')),
                        DropdownMenuItem(value: 'JPY', child: Text('JPY — 일본 엔 (¥)')),
                      ],
                      onChanged: (v) => setState(() => _currency = v!),
                    ),
                    const SizedBox(height: 16),
                    TextField(
                      controller: _balanceController,
                      keyboardType: TextInputType.number,
                      decoration: InputDecoration(
                        labelText: '시작 잔액',
                        border: OutlineInputBorder(borderRadius: BorderRadius.circular(14)),
                      ),
                    ),
                    const Spacer(),
                    SizedBox(
                      width: double.infinity,
                      height: 52,
                      child: ElevatedButton(
                        onPressed: _isLoading ? null : _submitFirstAccount,
                        style: ElevatedButton.styleFrom(
                          backgroundColor: AppColors.brandPrimary,
                          foregroundColor: Colors.white,
                          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
                        ),
                        child: _isLoading
                            ? const CircularProgressIndicator(color: Colors.white)
                            : const Text('다음', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                      ),
                    ),
                  ],
                )
              : Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Container(
                        width: 80,
                        height: 80,
                        decoration: BoxDecoration(
                          color: AppColors.success.withOpacity(0.1),
                          shape: BoxShape.circle,
                        ),
                        child: const Icon(Icons.check, color: AppColors.success, size: 48),
                      ),
                      const SizedBox(height: 24),
                      const Text(
                        '설정 완료!',
                        style: TextStyle(fontSize: 26, fontWeight: FontWeight.w900),
                      ),
                      const SizedBox(height: 8),
                      const Padding(
                        padding: EdgeInsets.symmetric(horizontal: 32.0),
                        child: Text(
                          '모든 준비가 끝났습니다. 이제 홈 대시보드로 이동합니다.',
                          textAlign: TextAlign.center,
                          style: TextStyle(color: AppColors.textSecondary, fontSize: 14),
                        ),
                      ),
                      const SizedBox(height: 40),
                      SizedBox(
                        width: 200,
                        height: 50,
                        child: ElevatedButton(
                          onPressed: _finishOnboarding,
                          style: ElevatedButton.styleFrom(
                            backgroundColor: AppColors.brandPrimary,
                            foregroundColor: Colors.white,
                            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
                          ),
                          child: const Text('대시보드 시작', style: TextStyle(fontWeight: FontWeight.bold)),
                        ),
                      ),
                    ],
                  ),
                ),
        ),
      ),
    );
  }
}

