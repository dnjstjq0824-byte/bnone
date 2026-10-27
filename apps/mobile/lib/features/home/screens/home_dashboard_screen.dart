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
    final currentMonth = DateTime.now().month;

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: Row(
          children: [
            Container(
              width: 28,
              height: 28,
              decoration: BoxDecoration(
                color: AppColors.brandPrimary,
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Center(
                child: Text(
                  '1',
                  style: TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.w900,
                    fontSize: 16,
                  ),
                ),
              ),
            ),
            const SizedBox(width: 8),
            const Text(
              'BN ONE',
              style: TextStyle(
                fontWeight: FontWeight.w900,
                fontSize: 20,
                letterSpacing: -0.5,
              ),
            ),
          ],
        ),
        actions: [
          IconButton(
            icon: Icon(
              _maskAmounts ? Icons.visibility_off_outlined : Icons.visibility_outlined,
              color: AppColors.textSecondary,
            ),
            onPressed: () => setState(() => _maskAmounts = !_maskAmounts),
            tooltip: '금액 숨기기',
          ),
          IconButton(
            icon: const Icon(
              Icons.notifications_none_rounded,
              color: AppColors.textSecondary,
            ),
            onPressed: () => context.push('/notifications'),
            tooltip: '알림',
          ),
          const SizedBox(width: 8),
        ],
      ),
      body: _isLoading
          ? const Center(
              child: CircularProgressIndicator(color: AppColors.brandPrimary),
            )
          : RefreshIndicator(
              color: AppColors.brandPrimary,
              onRefresh: _loadDashboardData,
              child: SingleChildScrollView(
                physics: const AlwaysScrollableScrollPhysics(),
                padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Legal Disclaimer Banner (Article 1.5)
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                      margin: const EdgeInsets.only(bottom: 12),
                      decoration: BoxDecoration(
                        color: AppColors.primaryLight,
                        borderRadius: BorderRadius.circular(16),
                      ),
                      child: Row(
                        children: [
                          const Icon(Icons.info_outline, size: 16, color: AppColors.brandPrimary),
                          const SizedBox(width: 10),
                          Expanded(
                            child: Text(
                              '표시된 자산과 잔액은 직접 입력한 장부 기록을 기준으로 합니다.',
                              style: TextStyle(
                                fontSize: 12,
                                color: AppColors.primaryDark,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),

                    // Net Worth Card (Toss Bank Card Style)
                    InkWell(
                      onTap: () => context.push('/assets'),
                      borderRadius: BorderRadius.circular(24),
                      child: Card(
                        child: Padding(
                          padding: const EdgeInsets.all(24.0),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Row(
                                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                children: [
                                  const Row(
                                    children: [
                                      Text(
                                        '내 순자산',
                                        style: TextStyle(
                                          fontSize: 15,
                                          fontWeight: FontWeight.w700,
                                          color: AppColors.textSecondary,
                                        ),
                                      ),
                                      SizedBox(width: 4),
                                      Icon(Icons.chevron_right, size: 18, color: AppColors.textTertiary),
                                    ],
                                  ),
                                  Container(
                                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                                    decoration: BoxDecoration(
                                      color: AppColors.background,
                                      borderRadius: BorderRadius.circular(8),
                                    ),
                                    child: Text(
                                      currency,
                                      style: const TextStyle(
                                        fontSize: 11,
                                        fontWeight: FontWeight.w700,
                                        color: AppColors.textSecondary,
                                      ),
                                    ),
                                  ),
                                ],
                              ),
                              const SizedBox(height: 12),
                              Text(
                                CurrencyFormatter.format(netWorthMinor, currency: currency, masked: _maskAmounts),
                                style: const TextStyle(
                                  fontSize: 32,
                                  fontWeight: FontWeight.w800,
                                  color: AppColors.textPrimary,
                                  letterSpacing: -1.0,
                                ),
                              ),
                              const SizedBox(height: 12),
                              Row(
                                children: [
                                  Text(
                                    '총자산 ${CurrencyFormatter.format(totalAssetsMinor, currency: currency, masked: _maskAmounts)}',
                                    style: const TextStyle(
                                      fontSize: 13,
                                      color: AppColors.textSecondary,
                                      fontWeight: FontWeight.w500,
                                    ),
                                  ),
                                  const SizedBox(width: 10),
                                  const Text('·', style: TextStyle(color: AppColors.textTertiary)),
                                  const SizedBox(width: 10),
                                  Text(
                                    '부채 ${CurrencyFormatter.format(totalLiabilitiesMinor, currency: currency, masked: _maskAmounts)}',
                                    style: const TextStyle(
                                      fontSize: 13,
                                      color: AppColors.danger,
                                      fontWeight: FontWeight.w500,
                                    ),
                                  ),
                                ],
                              ),
                              const SizedBox(height: 20),
                              const Divider(),
                              const SizedBox(height: 16),
                              // Toss Quick Actions Row
                              Row(
                                mainAxisAlignment: MainAxisAlignment.spaceAround,
                                children: [
                                  _buildTossQuickAction(
                                    icon: Icons.add,
                                    label: '거래 기록',
                                    color: AppColors.brandPrimary,
                                    bg: AppColors.primaryLight,
                                    onTap: () => context.push('/transactions/new'),
                                  ),
                                  _buildTossQuickAction(
                                    icon: Icons.account_balance_wallet_outlined,
                                    label: '자산 관리',
                                    color: const Color(0xFF8061D9),
                                    bg: const Color(0xFFF1EDFF),
                                    onTap: () => context.push('/assets'),
                                  ),
                                  _buildTossQuickAction(
                                    icon: Icons.track_changes_outlined,
                                    label: '예산 관리',
                                    color: const Color(0xFF00B06B),
                                    bg: const Color(0xFFE5F7F0),
                                    onTap: () => context.push('/budgets'),
                                  ),
                                  _buildTossQuickAction(
                                    icon: Icons.currency_exchange_rounded,
                                    label: '환율 계산',
                                    color: const Color(0xFFFF8D00),
                                    bg: const Color(0xFFFFF3E2),
                                    onTap: () => context.push('/fx'),
                                  ),
                                ],
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),

                    const SizedBox(height: 12),

                    // Monthly Spending Card (Toss "이번 달 쓴 금액" Style)
                    InkWell(
                      onTap: () => context.push('/transactions'),
                      borderRadius: BorderRadius.circular(24),
                      child: Card(
                        child: Padding(
                          padding: const EdgeInsets.all(24.0),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Row(
                                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                children: [
                                  Text(
                                    '$currentMonth월 소비',
                                    style: const TextStyle(
                                      fontSize: 15,
                                      fontWeight: FontWeight.w700,
                                      color: AppColors.textSecondary,
                                    ),
                                  ),
                                  const Icon(Icons.chevron_right, size: 18, color: AppColors.textTertiary),
                                ],
                              ),
                              const SizedBox(height: 10),
                              const Text(
                                '이번 달 쓴 금액',
                                style: TextStyle(
                                  fontSize: 13,
                                  color: AppColors.textTertiary,
                                ),
                              ),
                              const SizedBox(height: 4),
                              Text(
                                CurrencyFormatter.format(monthExpense, currency: currency, masked: _maskAmounts),
                                style: const TextStyle(
                                  fontSize: 28,
                                  fontWeight: FontWeight.w800,
                                  color: AppColors.textPrimary,
                                  letterSpacing: -0.8,
                                ),
                              ),
                              const SizedBox(height: 16),
                              Container(
                                padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
                                decoration: BoxDecoration(
                                  color: AppColors.background,
                                  borderRadius: BorderRadius.circular(12),
                                ),
                                child: Row(
                                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                  children: [
                                    const Text(
                                      '이번 달 들어온 돈',
                                      style: TextStyle(
                                        fontSize: 13,
                                        color: AppColors.textSecondary,
                                        fontWeight: FontWeight.w600,
                                      ),
                                    ),
                                    Text(
                                      '+ ${CurrencyFormatter.format(monthIncome, currency: currency, masked: _maskAmounts)}',
                                      style: const TextStyle(
                                        fontSize: 14,
                                        fontWeight: FontWeight.w700,
                                        color: AppColors.brandPrimary,
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),

                    const SizedBox(height: 12),

                    // Lifestyle Services Card (Toss Feature Banner Style)
                    Card(
                      child: Padding(
                        padding: const EdgeInsets.symmetric(vertical: 8),
                        child: Column(
                          children: [
                            ListTile(
                              leading: Container(
                                width: 42,
                                height: 42,
                                decoration: BoxDecoration(
                                  color: const Color(0xFFF1EDFF),
                                  borderRadius: BorderRadius.circular(14),
                                ),
                                child: const Icon(Icons.repeat_rounded, color: Color(0xFF8061D9), size: 22),
                              ),
                              title: const Text(
                                '매달 나가는 구독료',
                                style: TextStyle(fontWeight: FontWeight.w700, fontSize: 14),
                              ),
                              subtitle: const Text(
                                '정기 결제를 한곳에서 모아 관리해요',
                                style: TextStyle(fontSize: 12, color: AppColors.textTertiary),
                              ),
                              trailing: const Icon(Icons.chevron_right, size: 18, color: AppColors.textTertiary),
                              onTap: () => context.push('/subscriptions'),
                            ),
                            const Padding(
                              padding: EdgeInsets.symmetric(horizontal: 16),
                              child: Divider(),
                            ),
                            ListTile(
                              leading: Container(
                                width: 42,
                                height: 42,
                                decoration: BoxDecoration(
                                  color: const Color(0xFFE5F7F0),
                                  borderRadius: BorderRadius.circular(14),
                                ),
                                child: const Icon(Icons.flight_takeoff_rounded, color: Color(0xFF00B06B), size: 22),
                              ),
                              title: const Text(
                                '여행도 계획적으로',
                                style: TextStyle(fontWeight: FontWeight.w700, fontSize: 14),
                              ),
                              subtitle: const Text(
                                '여행 예산과 외화 지출을 현지 통화로 기록해요',
                                style: TextStyle(fontSize: 12, color: AppColors.textTertiary),
                              ),
                              trailing: const Icon(Icons.chevron_right, size: 18, color: AppColors.textTertiary),
                              onTap: () => context.push('/travel'),
                            ),
                          ],
                        ),
                      ),
                    ),

                    const SizedBox(height: 12),

                    // Recent Transactions Section
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
                                  '최근 거래',
                                  style: TextStyle(
                                    fontSize: 17,
                                    fontWeight: FontWeight.w700,
                                    color: AppColors.textPrimary,
                                  ),
                                ),
                                InkWell(
                                  onTap: () => context.push('/transactions'),
                                  child: const Padding(
                                    padding: EdgeInsets.symmetric(vertical: 4, horizontal: 4),
                                    child: Row(
                                      children: [
                                        Text(
                                          '전체보기',
                                          style: TextStyle(
                                            fontSize: 13,
                                            fontWeight: FontWeight.w600,
                                            color: AppColors.textTertiary,
                                          ),
                                        ),
                                        Icon(Icons.chevron_right, size: 16, color: AppColors.textTertiary),
                                      ],
                                    ),
                                  ),
                                ),
                              ],
                            ),
                            const SizedBox(height: 12),
                            if (_recentTx.isEmpty)
                              const Padding(
                                padding: EdgeInsets.symmetric(vertical: 24.0),
                                child: Center(
                                  child: Text(
                                    '아직 기록된 거래가 없어요.',
                                    style: TextStyle(color: AppColors.textTertiary, fontSize: 13),
                                  ),
                                ),
                              )
                            else
                              ListView.separated(
                                shrinkWrap: true,
                                physics: const NeverScrollableScrollPhysics(),
                                itemCount: _recentTx.length,
                                separatorBuilder: (_, __) => const Divider(height: 16),
                                itemBuilder: (context, index) {
                                  final tx = _recentTx[index];
                                  final isExpense = tx['transaction_type'] == 'EXPENSE';
                                  final isIncome = tx['transaction_type'] == 'INCOME';
                                  return Row(
                                    children: [
                                      Container(
                                        width: 40,
                                        height: 40,
                                        decoration: BoxDecoration(
                                          color: isIncome
                                              ? AppColors.primaryLight
                                              : AppColors.background,
                                          borderRadius: BorderRadius.circular(14),
                                        ),
                                        child: Icon(
                                          isIncome
                                              ? Icons.arrow_downward_rounded
                                              : isExpense
                                                  ? Icons.arrow_upward_rounded
                                                  : Icons.swap_horiz_rounded,
                                          color: isIncome
                                              ? AppColors.brandPrimary
                                              : AppColors.textSecondary,
                                          size: 20,
                                        ),
                                      ),
                                      const SizedBox(width: 14),
                                      Expanded(
                                        child: Column(
                                          crossAxisAlignment: CrossAxisAlignment.start,
                                          children: [
                                            Text(
                                              tx['description'] ?? tx['merchant_name'] ?? '거래 기록',
                                              style: const TextStyle(
                                                fontWeight: FontWeight.w700,
                                                fontSize: 14,
                                                color: AppColors.textPrimary,
                                              ),
                                              maxLines: 1,
                                              overflow: TextOverflow.ellipsis,
                                            ),
                                            const SizedBox(height: 2),
                                            Text(
                                              tx['category_name'] ?? (isIncome ? '수입' : isExpense ? '지출' : '계정 이동'),
                                              style: const TextStyle(
                                                fontSize: 12,
                                                color: AppColors.textTertiary,
                                              ),
                                            ),
                                          ],
                                        ),
                                      ),
                                      Text(
                                        '${isExpense ? "− " : isIncome ? "+ " : ""}${CurrencyFormatter.format(tx['amount_minor'] ?? 0, currency: tx['currency_code'] ?? 'KRW', masked: _maskAmounts)}',
                                        style: TextStyle(
                                          fontWeight: FontWeight.w800,
                                          fontSize: 15,
                                          color: isIncome
                                              ? AppColors.brandPrimary
                                              : AppColors.textPrimary,
                                        ),
                                      ),
                                    ],
                                  );
                                },
                              ),
                          ],
                        ),
                      ),
                    ),

                    const SizedBox(height: 24),
                  ],
                ),
              ),
            ),
      bottomNavigationBar: Container(
        decoration: const BoxDecoration(
          color: Colors.white,
          border: Border(top: BorderSide(color: AppColors.divider)),
        ),
        child: BottomNavigationBar(
          currentIndex: _currentIndex,
          onTap: _onBottomNavTapped,
          items: const [
            BottomNavigationBarItem(
              icon: Icon(Icons.home_outlined),
              activeIcon: Icon(Icons.home_filled),
              label: '홈',
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.account_balance_wallet_outlined),
              activeIcon: Icon(Icons.account_balance_wallet),
              label: '내 자산',
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.receipt_long_outlined),
              activeIcon: Icon(Icons.receipt_long),
              label: '소비',
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.pie_chart_outline_rounded),
              activeIcon: Icon(Icons.pie_chart_rounded),
              label: '계획',
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.menu_rounded),
              activeIcon: Icon(Icons.menu_open_rounded),
              label: '전체',
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildTossQuickAction({
    required IconData icon,
    required String label,
    required Color color,
    required Color bg,
    required VoidCallback onTap,
  }) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(16),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
        child: Column(
          children: [
            Container(
              width: 48,
              height: 48,
              decoration: BoxDecoration(
                color: bg,
                borderRadius: BorderRadius.circular(18),
              ),
              child: Icon(icon, color: color, size: 22),
            ),
            const SizedBox(height: 8),
            Text(
              label,
              style: const TextStyle(
                fontSize: 12,
                fontWeight: FontWeight.w600,
                color: AppColors.textSecondary,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
