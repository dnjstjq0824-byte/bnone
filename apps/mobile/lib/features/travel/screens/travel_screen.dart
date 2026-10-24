import 'package:flutter/material.dart';
import '../../../core/api/api_client.dart';
import '../../../core/theme/colors.dart';
import '../../../core/utils/formatters.dart';

class TravelScreen extends StatefulWidget {
  const TravelScreen({super.key});

  @override
  State<TravelScreen> createState() => _TravelScreenState();
}

class _TravelScreenState extends State<TravelScreen> {
  bool _isLoading = true;
  List<dynamic> _trips = [];

  @override
  void initState() {
    super.initState();
    _loadTrips();
  }

  Future<void> _loadTrips() async {
    setState(() => _isLoading = true);
    try {
      final list = await ApiClient().get('/trips');
      setState(() => _trips = list ?? []);
    } catch (e) {
      debugPrint('Trips error: $e');
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  void _showCreateTripDialog() {
    final nameCtrl = TextEditingController(text: '도쿄 4박 5일');
    final budgetCtrl = TextEditingController(text: '1500000');
    final cityCtrl = TextEditingController(text: 'Tokyo');

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('새 여행 등록', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(controller: nameCtrl, decoration: const InputDecoration(labelText: '여행 이름')),
            const SizedBox(height: 12),
            TextField(controller: cityCtrl, decoration: const InputDecoration(labelText: '도시/국가')),
            const SizedBox(height: 12),
            TextField(controller: budgetCtrl, keyboardType: TextInputType.number, decoration: const InputDecoration(labelText: '총 예산 (KRW)')),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
          ElevatedButton(
            onPressed: () async {
              Navigator.pop(ctx);
              await ApiClient().post('/trips', {
                'name': nameCtrl.text.trim(),
                'country_code': 'JP',
                'city_name': cityCtrl.text.trim(),
                'start_date': DateTime.now().toIso8601String().split('T')[0],
                'end_date': DateTime.now().add(const Duration(days: 5)).toIso8601String().split('T')[0],
                'budget_amount_minor': int.tryParse(budgetCtrl.text.trim()) ?? 1500000,
                'base_currency': 'KRW',
              });
              _loadTrips();
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
        title: const Text('여행 경비 & 예산'),
        actions: [
          IconButton(icon: const Icon(Icons.add), onPressed: _showCreateTripDialog),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _trips.isEmpty
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      const Icon(Icons.flight_takeoff, size: 48, color: AppColors.textSecondary),
                      const SizedBox(height: 16),
                      const Text('등록된 여행 일정이 없습니다.', style: TextStyle(fontWeight: FontWeight.bold)),
                      const SizedBox(height: 12),
                      ElevatedButton(onPressed: _showCreateTripDialog, child: const Text('여행 만들기')),
                    ],
                  ),
                )
              : ListView.builder(
                  padding: const EdgeInsets.all(16),
                  itemCount: _trips.length,
                  itemBuilder: (context, index) {
                    final t = _trips[index];
                    final budget = t['total_budget_minor'] ?? 0;
                    final spent = t['total_spent_minor'] ?? 0;
                    final remaining = budget - spent;

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
                                Text(t['name'] ?? '', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                                Container(
                                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                                  decoration: BoxDecoration(
                                    color: AppColors.brandPrimary.withOpacity(0.1),
                                    borderRadius: BorderRadius.circular(8),
                                  ),
                                  child: Text(t['status'] ?? 'PLANNED', style: const TextStyle(fontSize: 10, fontWeight: FontWeight.bold, color: AppColors.brandPrimary)),
                                ),
                              ],
                            ),
                            const SizedBox(height: 4),
                            Text('${t['city_name']} • ${t['start_date']} ~ ${t['end_date']}', style: const TextStyle(fontSize: 11, color: AppColors.textSecondary)),
                            const SizedBox(height: 16),
                            Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                Text('총 예산: ${CurrencyFormatter.format(budget)}', style: const TextStyle(fontSize: 12, color: AppColors.textSecondary)),
                                Text('사용액: ${CurrencyFormatter.format(spent)}', style: const TextStyle(fontSize: 12, fontWeight: FontWeight.bold)),
                              ],
                            ),
                            const SizedBox(height: 8),
                            Text(
                              '남은 예산: ${CurrencyFormatter.format(remaining)}',
                              style: const TextStyle(fontWeight: FontWeight.bold, color: AppColors.brandPrimary),
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

