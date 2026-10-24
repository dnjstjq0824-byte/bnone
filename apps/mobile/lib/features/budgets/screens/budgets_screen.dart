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

  void _showCreateBudgetDialog() {
    final nameCtrl = TextEditingController(text: '이번 달 생활비');
    final limitCtrl = TextEditingController(text: '1000000');

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('새 예산 생성', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(controller: nameCtrl, decoration: const InputDecoration(labelText: '예산 이름')),
            const SizedBox(height: 12),
            TextField(
              controller: limitCtrl,
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(labelText: '예산 한도 (KRW)'),
            ),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
          ElevatedButton(
            onPressed: () async {
              Navigator.pop(ctx);
              await ApiClient().post('/budgets', {
                'name': nameCtrl.text.trim(),
                'limit_minor': int.tryParse(limitCtrl.text.trim()) ?? 1000000,
                'period_type': 'MONTHLY',
              });
              _loadBudgets();
            },
            child: const Text('생성'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('예산 관리'),
        actions: [
          IconButton(icon: const Icon(Icons.add), onPressed: _showCreateBudgetDialog),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _budgets.isEmpty
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      const Icon(Icons.track_changes, size: 48, color: AppColors.textSecondary),
                      const SizedBox(height: 16),
                      const Text('설정된 예산이 없습니다.', style: TextStyle(fontWeight: FontWeight.bold)),
                      const SizedBox(height: 12),
                      ElevatedButton(onPressed: _showCreateBudgetDialog, child: const Text('예산 만들기')),
                    ],
                  ),
                )
              : ListView.builder(
                  padding: const EdgeInsets.all(16),
                  itemCount: _budgets.length,
                  itemBuilder: (context, index) {
                    final b = _budgets[index];
                    final double percent = (b['usage_percent'] ?? 0.0).toDouble();

                    return Card(
                      margin: const EdgeInsets.only(bottom: 12),
                      child: Padding(
                        padding: const EdgeInsets.all(16),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                Text(b['name'] ?? '', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 15)),
                                Container(
                                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                                  decoration: BoxDecoration(
                                    color: percent >= 100
                                        ? AppColors.danger.withOpacity(0.1)
                                        : percent >= 80
                                            ? AppColors.warning.withOpacity(0.1)
                                            : AppColors.success.withOpacity(0.1),
                                    borderRadius: BorderRadius.circular(8),
                                  ),
                                  child: Text(
                                    '${percent.toStringAsFixed(1)}% 사용',
                                    style: TextStyle(
                                      color: percent >= 100
                                          ? AppColors.danger
                                          : percent >= 80
                                              ? AppColors.warning
                                              : AppColors.success,
                                      fontWeight: FontWeight.bold,
                                      fontSize: 11,
                                    ),
                                  ),
                                ),
                              ],
                            ),
                            const SizedBox(height: 12),
                            LinearProgressIndicator(
                              value: (percent / 100).clamp(0.0, 1.0),
                              backgroundColor: AppColors.border,
                              color: percent >= 100
                                  ? AppColors.danger
                                  : percent >= 80
                                      ? AppColors.warning
                                      : AppColors.brandPrimary,
                              minHeight: 8,
                              borderRadius: BorderRadius.circular(4),
                            ),
                            const SizedBox(height: 12),
                            Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                Text('사용: ${CurrencyFormatter.format(b['spent_minor'] ?? 0)}', style: const TextStyle(fontSize: 12, color: AppColors.textSecondary)),
                                Text('한도: ${CurrencyFormatter.format(b['limit_minor'] ?? 0)}', style: const TextStyle(fontSize: 12, fontWeight: FontWeight.bold)),
                              ],
                            ),
                          ],
                        ),
                      ),
                    );
                  },
                ),
    );
  }
}

