import 'package:flutter/material.dart';
import '../../../core/api/api_client.dart';
import '../../../core/theme/colors.dart';
import '../../../core/utils/formatters.dart';

class BudgetsScreen extends StatefulWidget {
  const BudgetsScreen({super.key});

  @override
  State<BudgetsScreen> createState() => _BudgetsScreenState();
}

class _BudgetsScreenState extends State<BudgetsScreen> {
  bool _isLoading = true;
  List<dynamic> _budgets = [];

  @override
  void initState() {
    super.initState();
    _loadBudgets();
  }

  Future<void> _loadBudgets() async {
    setState(() => _isLoading = true);
    try {
      final list = await ApiClient().get('/budgets');
      setState(() => _budgets = list ?? []);
    } catch (e) {
      debugPrint('Budgets error: $e');
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  void _showCreateBudgetSheet() {
    final nameCtrl = TextEditingController(text: '이번 달 생활비');
    final limitCtrl = TextEditingController(text: '1000000');

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.white,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
      ),
      builder: (ctx) => Padding(
        padding: EdgeInsets.only(
          left: 24,
          right: 24,
          top: 28,
          bottom: MediaQuery.of(ctx).viewInsets.bottom + 28,
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text(
                  '새 목표 예산 설정',
                  style: TextStyle(
                    fontSize: 19,
                    fontWeight: FontWeight.w800,
                    color: AppColors.textPrimary,
                  ),
                ),
                IconButton(
                  icon: const Icon(Icons.close, color: AppColors.textTertiary),
                  onPressed: () => Navigator.pop(ctx),
                ),
              ],
            ),
            const SizedBox(height: 16),
            TextField(
              controller: nameCtrl,
              decoration: const InputDecoration(labelText: '예산 이름', hintText: '예: 이번 달 생활비, 식비 예산'),
            ),
            const SizedBox(height: 14),
            TextField(
              controller: limitCtrl,
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(labelText: '예산 한도', hintText: '1,000,000', suffixText: '원'),
            ),
            const SizedBox(height: 24),
            ElevatedButton(
              onPressed: () async {
                if (nameCtrl.text.trim().isEmpty) return;
                Navigator.pop(ctx);
                await ApiClient().post('/budgets', {
                  'name': nameCtrl.text.trim(),
                  'limit_minor': int.tryParse(limitCtrl.text.trim()) ?? 1000000,
                  'period_type': 'MONTHLY',
                });
                _loadBudgets();
              },
              child: const Text('예산 만들기'),
            ),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('예산 계획'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add, color: AppColors.textPrimary),
            onPressed: _showCreateBudgetSheet,
            tooltip: '예산 만들기',
          ),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator(color: AppColors.brandPrimary))
          : _budgets.isEmpty
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Container(
                        width: 64,
                        height: 64,
                        decoration: BoxDecoration(
                          color: const Color(0xFFE5F7F0),
                          borderRadius: BorderRadius.circular(20),
                        ),
                        child: const Icon(Icons.track_changes_rounded, size: 32, color: Color(0xFF00B06B)),
                      ),
                      const SizedBox(height: 16),
                      const Text(
                        '설정된 예산이 없습니다.',
                        style: TextStyle(fontWeight: FontWeight.w700, fontSize: 16, color: AppColors.textPrimary),
                      ),
                      const SizedBox(height: 8),
                      const Text(
                        '목표 예산을 세우고 똑똑하게 지출을 관리해보세요.',
                        style: TextStyle(fontSize: 13, color: AppColors.textTertiary),
                      ),
                      const SizedBox(height: 20),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 40),
                        child: ElevatedButton(
                          onPressed: _showCreateBudgetSheet,
                          child: const Text('첫 예산 세우기'),
                        ),
                      ),
                    ],
                  ),
                )
              : RefreshIndicator(
                  color: AppColors.brandPrimary,
                  onRefresh: _loadBudgets,
                  child: ListView.builder(
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                    itemCount: _budgets.length,
                    itemBuilder: (context, index) {
                      final b = _budgets[index];
                      final double percent = (b['usage_percent'] ?? 0.0).toDouble();
                      final remaining = ((b['limit_minor'] ?? 0) - (b['spent_minor'] ?? 0)).clamp(0, 9999999999);

                      return Card(
                        child: Padding(
                          padding: const EdgeInsets.all(22),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Row(
                                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                children: [
                                  Text(
                                    b['name'] ?? '',
                                    style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 16, color: AppColors.textPrimary),
                                  ),
                                  Container(
                                    padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                                    decoration: BoxDecoration(
                                      color: percent >= 100
                                          ? const Color(0xFFFEECEB)
                                          : percent >= 80
                                              ? const Color(0xFFFFF3E2)
                                              : AppColors.primaryLight,
                                      borderRadius: BorderRadius.circular(10),
                                    ),
                                    child: Text(
                                      '${percent.toStringAsFixed(0)}% 사용',
                                      style: TextStyle(
                                        color: percent >= 100
                                            ? AppColors.danger
                                            : percent >= 80
                                                ? const Color(0xFFFF8D00)
                                                : AppColors.brandPrimary,
                                        fontWeight: FontWeight.w700,
                                        fontSize: 12,
                                      ),
                                    ),
                                  ),
                                ],
                              ),
                              const SizedBox(height: 14),
                              Text(
                                '남은 예산 ${CurrencyFormatter.format(remaining)}',
                                style: const TextStyle(
                                  fontSize: 22,
                                  fontWeight: FontWeight.w800,
                                  color: AppColors.textPrimary,
                                  letterSpacing: -0.5,
                                ),
                              ),
                              const SizedBox(height: 14),
                              ClipRRect(
                                borderRadius: BorderRadius.circular(8),
                                child: LinearProgressIndicator(
                                  value: (percent / 100).clamp(0.0, 1.0),
                                  backgroundColor: AppColors.background,
                                  color: percent >= 100
                                      ? AppColors.danger
                                      : percent >= 80
                                          ? const Color(0xFFFF8D00)
                                          : AppColors.brandPrimary,
                                  minHeight: 8,
                                ),
                              ),
                              const SizedBox(height: 14),
                              Row(
                                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                children: [
                                  Text(
                                    '쓴 돈 ${CurrencyFormatter.format(b['spent_minor'] ?? 0)}',
                                    style: const TextStyle(fontSize: 13, color: AppColors.textSecondary, fontWeight: FontWeight.w500),
                                  ),
                                  Text(
                                    '목표 ${CurrencyFormatter.format(b['limit_minor'] ?? 0)}',
                                    style: const TextStyle(fontSize: 13, fontWeight: FontWeight.w600, color: AppColors.textPrimary),
                                  ),
                                ],
                              ),
                            ],
                          ),
                        ),
                      );
                    },
                  ),
                ),
    );
  }
}
