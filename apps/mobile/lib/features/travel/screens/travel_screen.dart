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

  void _showCreateTripSheet() {
    final nameCtrl = TextEditingController(text: '도쿄 4박 5일');
    final budgetCtrl = TextEditingController(text: '1500000');
    final cityCtrl = TextEditingController(text: '도쿄');

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
                  '새 여행 일정 등록',
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
              decoration: const InputDecoration(labelText: '여행 이름', hintText: '예: 도쿄 4박 5일, 파리 배낭여행'),
            ),
            const SizedBox(height: 14),
            TextField(
              controller: cityCtrl,
              decoration: const InputDecoration(labelText: '여행 도시 / 국가', hintText: '예: 도쿄 (일본), 파리 (프랑스)'),
            ),
            const SizedBox(height: 14),
            TextField(
              controller: budgetCtrl,
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(labelText: '총 예산 한도', hintText: '1,500,000', suffixText: '원'),
            ),
            const SizedBox(height: 24),
            ElevatedButton(
              onPressed: () async {
                if (nameCtrl.text.trim().isEmpty) return;
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
              child: const Text('여행 등록하기'),
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
        title: const Text('여행 경비 & 예산'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add, color: AppColors.textPrimary),
            onPressed: _showCreateTripSheet,
            tooltip: '여행 등록',
          ),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator(color: AppColors.brandPrimary))
          : _trips.isEmpty
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
                        child: const Icon(Icons.flight_takeoff_rounded, size: 32, color: Color(0xFF00B06B)),
                      ),
                      const SizedBox(height: 16),
                      const Text(
                        '등록된 여행 일정이 없습니다.',
                        style: TextStyle(fontWeight: FontWeight.w700, fontSize: 16, color: AppColors.textPrimary),
                      ),
                      const SizedBox(height: 8),
                      const Text(
                        '여행지 화폐로 현지 지출을 기록하고 예산을 관리하세요.',
                        style: TextStyle(fontSize: 13, color: AppColors.textTertiary),
                      ),
                      const SizedBox(height: 20),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 40),
                        child: ElevatedButton(
                          onPressed: _showCreateTripSheet,
                          child: const Text('첫 여행 만들기'),
                        ),
                      ),
                    ],
                  ),
                )
              : RefreshIndicator(
                  color: AppColors.brandPrimary,
                  onRefresh: _loadTrips,
                  child: ListView.builder(
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                    itemCount: _trips.length,
                    itemBuilder: (context, index) {
                      final trip = _trips[index];
                      final spentMinor = trip['total_spent_minor'] ?? 0;
                      final budgetMinor = trip['budget_amount_minor'] ?? 0;

                      return Card(
                        child: Padding(
                          padding: const EdgeInsets.all(22.0),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Row(
                                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                children: [
                                  Row(
                                    children: [
                                      Container(
                                        width: 36,
                                        height: 36,
                                        decoration: BoxDecoration(
                                          color: const Color(0xFFE5F7F0),
                                          borderRadius: BorderRadius.circular(12),
                                        ),
                                        child: const Icon(Icons.flight_rounded, color: Color(0xFF00B06B), size: 20),
                                      ),
                                      const SizedBox(width: 10),
                                      Text(
                                        trip['name'] ?? '',
                                        style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 16, color: AppColors.textPrimary),
                                      ),
                                    ],
                                  ),
                                  Container(
                                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                                    decoration: BoxDecoration(
                                      color: AppColors.background,
                                      borderRadius: BorderRadius.circular(8),
                                    ),
                                    child: Text(
                                      trip['city_name'] ?? '해외',
                                      style: const TextStyle(fontSize: 12, fontWeight: FontWeight.w600, color: AppColors.textSecondary),
                                    ),
                                  ),
                                ],
                              ),
                              const SizedBox(height: 16),
                              Text(
                                '예산 ${CurrencyFormatter.format(budgetMinor, currency: trip['base_currency'] ?? 'KRW')}',
                                style: const TextStyle(
                                  fontSize: 22,
                                  fontWeight: FontWeight.w800,
                                  color: AppColors.textPrimary,
                                  letterSpacing: -0.5,
                                ),
                              ),
                              const SizedBox(height: 8),
                              Row(
                                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                children: [
                                  Text(
                                    '현재 사용한 경비: ${CurrencyFormatter.format(spentMinor, currency: trip['base_currency'] ?? 'KRW')}',
                                    style: const TextStyle(fontSize: 13, color: AppColors.textSecondary),
                                  ),
                                  Text(
                                    '${trip['start_date']} ~ ${trip['end_date']}',
                                    style: const TextStyle(fontSize: 12, color: AppColors.textTertiary),
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
