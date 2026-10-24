import 'package:flutter/material.dart';
import '../../../core/api/api_client.dart';
import '../../../core/theme/colors.dart';
import '../../../core/utils/formatters.dart';

class SubscriptionsScreen extends StatefulWidget {
  const SubscriptionsScreen({super.key});

  @override
  State<SubscriptionsScreen> createState() => _SubscriptionsScreenState();
}

class _SubscriptionsScreenState extends State<SubscriptionsScreen> {
  bool _isLoading = true;
  List<dynamic> _subscriptions = [];
  Map<String, dynamic>? _summary;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    setState(() => _isLoading = true);
    try {
      final subs = await ApiClient().get('/me/subscriptions');
      final sum = await ApiClient().get('/me/subscriptions/summary');
      setState(() {
        _subscriptions = subs ?? [];
        _summary = sum;
      });
    } catch (e) {
      debugPrint('Subs error: $e');
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  void _showAddSubscriptionDialog() {
    final nameCtrl = TextEditingController(text: 'Netflix');
    final amountCtrl = TextEditingController(text: '17000');

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('정기 구독 등록', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(controller: nameCtrl, decoration: const InputDecoration(labelText: '서비스명')),
            const SizedBox(height: 12),
            TextField(controller: amountCtrl, keyboardType: TextInputType.number, decoration: const InputDecoration(labelText: '월간 금액 (KRW)')),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
          ElevatedButton(
            onPressed: () async {
              Navigator.pop(ctx);
              await ApiClient().post('/me/subscriptions', {
                'custom_service_name': nameCtrl.text.trim(),
                'amount_minor': int.tryParse(amountCtrl.text.trim()) ?? 17000,
                'billing_cycle': 'MONTHLY',
                'currency_code': 'KRW',
                'start_date': DateTime.now().toIso8601String().split('T')[0],
                'next_renewal_at': DateTime.now().add(const Duration(days: 30)).toIso8601String().split('T')[0],
              });
              _loadData();
            },
            child: const Text('등록'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('정기 구독'),
        actions: [
          IconButton(icon: const Icon(Icons.add), onPressed: _showAddSubscriptionDialog),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : RefreshIndicator(
              onRefresh: _loadData,
              child: SingleChildScrollView(
                physics: const AlwaysScrollableScrollPhysics(),
                padding: const EdgeInsets.all(16),
                child: Column(
                  children: [
                    // Disclaimer
                    Container(
                      padding: const EdgeInsets.all(12),
                      margin: const EdgeInsets.only(bottom: 16),
                      decoration: BoxDecoration(
                        color: AppColors.brandPrimary.withOpacity(0.08),
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: const Text(
                        'BN ONE의 구독 관리는 참고 기록이며 실제 서비스 자동 해지를 대행하지 않습니다.',
                        style: TextStyle(fontSize: 11, color: AppColors.brandPrimary),
                      ),
                    ),

                    // Monthly Estimate Card
                    Card(
                      child: Padding(
                        padding: const EdgeInsets.all(16.0),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                const Text('월간 예상 구독비', style: TextStyle(fontSize: 12, color: AppColors.textSecondary, fontWeight: FontWeight.bold)),
                                const SizedBox(height: 4),
                                Text(
                                  CurrencyFormatter.format(_summary?['monthly_estimate_minor'] ?? 0),
                                  style: const TextStyle(fontSize: 24, fontWeight: FontWeight.w900, color: AppColors.brandPrimary),
                                ),
                              ],
                            ),
                            Text(
                              '총 ${_summary?['active_count'] ?? 0}개 서비스',
                              style: const TextStyle(fontSize: 12, fontWeight: FontWeight.bold, color: AppColors.textSecondary),
                            ),
                          ],
                        ),
                      ),
                    ),
                    const SizedBox(height: 16),

                    ListView.builder(
                      shrinkWrap: true,
                      physics: const NeverScrollableScrollPhysics(),
                      itemCount: _subscriptions.length,
                      itemBuilder: (context, index) {
                        final s = _subscriptions[index];
                        return Card(
                          margin: const EdgeInsets.only(bottom: 8),
                          child: ListTile(
                            leading: CircleAvatar(
                              backgroundColor: AppColors.brandPrimary.withOpacity(0.1),
                              child: Text(
                                (s['service_name'] ?? 'S').substring(0, 1),
                                style: const TextStyle(fontWeight: FontWeight.bold, color: AppColors.brandPrimary),
                              ),
                            ),
                            title: Text(s['service_name'] ?? '', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
                            subtitle: Text('갱신일: ${s['next_renewal_at']}', style: const TextStyle(fontSize: 11, color: AppColors.textSecondary)),
                            trailing: Text(
                              CurrencyFormatter.format(s['amount_minor'] ?? 0, currency: s['currency_code'] ?? 'KRW'),
                              style: const TextStyle(fontWeight: FontWeight.w900, fontSize: 14),
                            ),
                          ),
                        );
                      },
                    ),
                  ],
                ),
              ),
            ),
    );
  }
}
