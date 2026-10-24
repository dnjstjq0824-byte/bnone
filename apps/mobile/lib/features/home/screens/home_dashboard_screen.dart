import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../core/api/api_client.dart';
import '../../../core/theme/colors.dart';
import '../../../core/utils/formatters.dart';

class HomeDashboardScreen extends StatefulWidget {
  const HomeDashboardScreen({super.key});

  @override
  State<HomeDashboardScreen> createState() => _HomeDashboardScreenState();
}

class _HomeDashboardScreenState extends State<HomeDashboardScreen> {
  int _currentIndex = 0;
  bool _isLoading = true;
  bool _maskAmounts = false;

  Map<String, dynamic>? _netWorth;
  Map<String, dynamic>? _overview;
  List<dynamic> _recentTx = [];

  @override
  void initState() {
    super.initState();
    _loadDashboardData();
  }

  Future<void> _loadDashboardData() async {
    setState(() => _isLoading = true);
    try {
      final nw = await ApiClient().get('/net-worth');
      final ov = await ApiClient().get('/analytics/overview');
      final tx = await ApiClient().get('/transactions?limit=5');
      setState(() {
        _netWorth = nw;
        _overview = ov;
        _recentTx = tx ?? [];
      });
    } catch (e) {
      debugPrint('Dashboard load error: $e');
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  void _onBottomNavTapped(int index) {
    setState(() => _currentIndex = index);
    switch (index) {
      case 0:
        break;
      case 1:
        context.push('/assets');
        break;
      case 2:
        context.push('/transactions');
        break;
      case 3:
        context.push('/budgets');
        break;
      case 4:
        context.push('/settings');
        break;
    }
  }

  @override
  Widget build(BuildContext context) {
    final currency = _netWorth?['base_currency'] ?? 'KRW';
    final netWorthMinor = _netWorth?['net_worth_minor'] ?? 0;
    final totalAssetsMinor = _netWorth?['total_assets_minor'] ?? 0;
    final totalLiabilitiesMinor = _netWorth?['total_liabilities_minor'] ?? 0;

    final monthExpense = _overview?['current_month_expense_minor'] ?? 0;
    final monthIncome = _overview?['current_month_income_minor'] ?? 0;

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('BN ONE', style: TextStyle(fontWeight: FontWeight.w900)),
        actions: [
          IconButton(
            icon: Icon(_maskAmounts ? Icons.visibility_off : Icons.visibility),
            onPressed: () => setState(() => _maskAmounts = !_maskAmounts),
            tooltip: '금액 숨기기',
          ),
          IconButton(
            icon: const Icon(Icons.notifications_outlined),
            onPressed: () => context.push('/notifications'),
          ),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : RefreshIndicator(
              onRefresh: _loadDashboardData,
              child: SingleChildScrollView(
                physics: const AlwaysScrollableScrollPhysics(),
                padding: const EdgeInsets.all(16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Mandatory Disclaimer (Article 1.5)
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
                      margin: const EdgeInsets.only(bottom: 16),
                      decoration: BoxDecoration(
                        color: AppColors.brandPrimary.withOpacity(0.08),
                        borderRadius: BorderRadius.circular(14),
                      ),
                      child: const Row(
                        children: [
                          Icon(Icons.info_outline, size: 16, color: AppColors.brandPrimary),
                          SizedBox(width: 8),
                          Expanded(
                            child: Text(
                              '표시된 자산과 잔액은 사용자가 직접 입력한 기록을 기준으로 합니다.',
                              style: TextStyle(fontSize: 11, color: AppColors.brandPrimary, fontWeight: FontWeight.bold),
                            ),
                          ),
                        ],
                      ),
                    ),

                    // Net Worth Card
                    Card(
                      child: Padding(
                        padding: const EdgeInsets.all(20.0),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                const Text(
                                  '내 순자산',
                                  style: TextStyle(fontSize: 12, fontWeight: FontWeight.bold, color: AppColors.textSecondary),
                                ),
                                Container(
                                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                                  decoration: BoxDecoration(
                                    color: AppColors.background,
                                    borderRadius: BorderRadius.circular(8),
                                  ),
                                  child: Text(
                                    currency,
                                    style: const TextStyle(fontSize: 10, fontWeight: FontWeight.bold, color: AppColors.textSecondary),
                                  ),
                                ),
                              ],
                            ),
                            const SizedBox(height: 8),
                            Text(
                              CurrencyFormatter.format(netWorthMinor, currency: currency, masked: _maskAmounts),
                              style: const TextStyle(fontSize: 32, fontWeight: FontWeight.w900, color: AppColors.textPrimary),
                            ),
                            const SizedBox(height: 16),
                            const Divider(height: 1),
                            const SizedBox(height: 12),
                            Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                Text(
                                  '총 자산: ${CurrencyFormatter.format(totalAssetsMinor, currency: currency, masked: _maskAmounts)}',
                                  style: const TextStyle(fontSize: 12, color: AppColors.textSecondary),
                                ),
                                Text(
                                  '부채: ${CurrencyFormatter.format(totalLiabilitiesMinor, currency: currency, masked: _maskAmounts)}',
                                  style: const TextStyle(fontSize: 12, color: AppColors.textSecondary),
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),
                    ),
                    const SizedBox(height: 16),

                    // Month Summary Grid (Expense & Income)
                    Row(
                      children: [
                        Expanded(
                          child: Card(
                            child: Padding(
                              padding: const EdgeInsets.all(16.0),
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  const Text('이번 달 지출', style: TextStyle(fontSize: 11, color: AppColors.textSecondary, fontWeight: FontWeight.bold)),
                                  const SizedBox(height: 4),
                                  Text(
                                    '− ${CurrencyFormatter.format(monthExpense, currency: currency, masked: _maskAmounts)}',
                                    style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w900, color: AppColors.danger),
                                  ),
                                ],
                              ),
                            ),
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: Card(
                            child: Padding(
                              padding: const EdgeInsets.all(16.0),
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  const Text('이번 달 수입', style: TextStyle(fontSize: 11, color: AppColors.textSecondary, fontWeight: FontWeight.bold)),
                                  const SizedBox(height: 4),
                                  Text(
                                    '+ ${CurrencyFormatter.format(monthIncome, currency: currency, masked: _maskAmounts)}',
                                    style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w900, color: AppColors.success),
                                  ),
                                ],
                              ),
                            ),
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 20),

                    // Quick Actions
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceAround,
                      children: [
                        _buildQuickAction(Icons.add_circle, '거래 추가', () => context.push('/transactions/new')),
                        _buildQuickAction(Icons.account_balance, '자산 장부', () => context.push('/assets')),
                        _buildQuickAction(Icons.flight_takeoff, '여행 예산', () => context.push('/travel')),
                        _buildQuickAction(Icons.currency_exchange, '참고 환율', () => context.push('/fx')),
                      ],
                    ),
                    const SizedBox(height: 24),

                    // Recent Transactions Section
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text(
                          '최근 거래',
                          style: TextStyle(fontSize: 16, fontWeight: FontWeight.w800, color: AppColors.textPrimary),
                        ),
                        TextButton(
                          onPressed: () => context.push('/transactions'),
                          child: const Text('전체보기', style: TextStyle(fontSize: 12, fontWeight: FontWeight.bold)),
                        ),
                      ],
                    ),
                    const SizedBox(height: 8),

                    if (_recentTx.isEmpty)
                      const Center(
                        child: Padding(
                          padding: EdgeInsets.all(24.0),
                          child: Text('기록된 거래 내역이 없습니다.', style: TextStyle(color: AppColors.textSecondary, fontSize: 13)),
                        ),
                      )
                    else
                      Card(
                        child: ListView.separated(
                          shrinkWrap: true,
                          physics: const NeverScrollableScrollPhysics(),
                          itemCount: _recentTx.length,
                          separatorBuilder: (_, __) => const Divider(height: 1),
                          itemBuilder: (context, index) {
                            final tx = _recentTx[index];
                            final isExpense = tx['transaction_type'] == 'EXPENSE';
                            final isIncome = tx['transaction_type'] == 'INCOME';
                            return ListTile(
                              leading: CircleAvatar(
                                backgroundColor: isExpense
                                    ? AppColors.danger.withOpacity(0.1)
                                    : isIncome
                                        ? AppColors.success.withOpacity(0.1)
                                        : AppColors.brandPrimary.withOpacity(0.1),
                                child: Icon(
                                  isExpense ? Icons.arrow_downward : isIncome ? Icons.arrow_upward : Icons.swap_horiz,
                                  color: isExpense ? AppColors.danger : isIncome ? AppColors.success : AppColors.brandPrimary,
                                  size: 18,
                                ),
                              ),
                              title: Text(
                                tx['description'] ?? tx['merchant_name'] ?? '거래',
                                style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 13),
                              ),
                              subtitle: Text(
                                tx['category_name'] ?? tx['transaction_type'],
                                style: const TextStyle(fontSize: 11, color: AppColors.textSecondary),
                              ),
                              trailing: Text(
                                '${isExpense ? "− " : isIncome ? "+ " : ""}${CurrencyFormatter.format(tx['amount_minor'] ?? 0, currency: tx['currency_code'] ?? 'KRW', masked: _maskAmounts)}',
                                style: TextStyle(
                                  fontWeight: FontWeight.w900,
                                  fontSize: 14,
                                  color: isExpense ? AppColors.textPrimary : isIncome ? AppColors.success : AppColors.brandPrimary,
                                ),
                              ),
                            );
                          },
                        ),
                      ),
                  ],
                ),
              ),
            ),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: _currentIndex,
        onTap: _onBottomNavTapped,
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.home_outlined), activeIcon: Icon(Icons.home), label: 'HOME'),
          BottomNavigationBarItem(icon: Icon(Icons.account_balance_outlined), activeIcon: Icon(Icons.account_balance), label: 'ASSETS'),
          BottomNavigationBarItem(icon: Icon(Icons.receipt_long_outlined), activeIcon: Icon(Icons.receipt_long), label: 'ACTIVITY'),
          BottomNavigationBarItem(icon: Icon(Icons.pie_chart_outline), activeIcon: Icon(Icons.pie_chart), label: 'PLAN'),
          BottomNavigationBarItem(icon: Icon(Icons.more_horiz), activeIcon: Icon(Icons.more_horiz), label: 'MORE'),
        ],
      ),
    );
  }

  Widget _buildQuickAction(IconData icon, String label, VoidCallback onTap) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(16),
      child: Column(
        children: [
          Container(
            width: 52,
            height: 52,
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(16),
              border: Border.all(color: AppColors.border),
            ),
            child: Icon(icon, color: AppColors.brandPrimary),
          ),
          const SizedBox(height: 6),
          Text(label, style: const TextStyle(fontSize: 11, fontWeight: FontWeight.bold, color: AppColors.textPrimary)),
        ],
      ),
    );
  }
}
