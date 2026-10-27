import 'package:flutter/material.dart';
import '../../../core/api/api_client.dart';
import '../../../core/theme/colors.dart';
import '../../../core/utils/formatters.dart';

class AssetsScreen extends StatefulWidget {
  const AssetsScreen({super.key});

  @override
  State<AssetsScreen> createState() => _AssetsScreenState();
}

class _AssetsScreenState extends State<AssetsScreen> {
  bool _isLoading = true;
  List<dynamic> _accounts = [];
  List<dynamic> _assets = [];
  List<dynamic> _liabilities = [];
  int _totalNetWorth = 0;
  int _totalAssets = 0;
  int _totalLiabilities = 0;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    setState(() => _isLoading = true);
    try {
      final nw = await ApiClient().get('/net-worth');
      final accs = await ApiClient().get('/accounts');
      final asts = await ApiClient().get('/assets');
      final liabs = await ApiClient().get('/liabilities');
      setState(() {
        _accounts = accs ?? [];
        _assets = asts ?? [];
        _liabilities = liabs ?? [];
        _totalNetWorth = nw?['net_worth_minor'] ?? 0;
        _totalAssets = nw?['total_assets_minor'] ?? 0;
        _totalLiabilities = nw?['total_liabilities_minor'] ?? 0;
      });
    } catch (e) {
      debugPrint('Assets load error: $e');
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  void _showAddAccountSheet() {
    final nameCtrl = TextEditingController();
    final balCtrl = TextEditingController(text: '0');
    String type = 'MANUAL_BANK';
    String curr = 'KRW';

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.white,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
      ),
      builder: (ctx) => StatefulBuilder(
        builder: (context, setModalState) => Padding(
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
                    '새 금융 계정 등록',
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
                decoration: const InputDecoration(
                  labelText: '계정 이름',
                  hintText: '예: 급여통장, 카카오뱅크, 비상금 현금',
                ),
              ),
              const SizedBox(height: 14),
              DropdownButtonFormField<String>(
                value: type,
                items: const [
                  DropdownMenuItem(value: 'MANUAL_BANK', child: Text('수동 은행 계좌')),
                  DropdownMenuItem(value: 'CASH', child: Text('현금 지갑')),
                  DropdownMenuItem(value: 'SAVINGS', child: Text('저축 / 적금')),
                  DropdownMenuItem(value: 'E_WALLET_MANUAL', child: Text('간편결제 / 전자지갑')),
                ],
                onChanged: (v) => setModalState(() => type = v!),
                decoration: const InputDecoration(labelText: '계정 유형'),
              ),
              const SizedBox(height: 14),
              TextField(
                controller: balCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: '시작 잔액',
                  hintText: '0',
                  suffixText: '원',
                ),
              ),
              const SizedBox(height: 24),
              ElevatedButton(
                onPressed: () async {
                  if (nameCtrl.text.trim().isEmpty) return;
                  Navigator.pop(ctx);
                  await ApiClient().post('/accounts', {
                    'name': nameCtrl.text.trim(),
                    'account_type': type,
                    'currency_code': curr,
                    'opening_balance_minor': int.tryParse(balCtrl.text.trim()) ?? 0,
                    'opening_balance_date': DateTime.now().toIso8601String().split('T')[0],
                  });
                  _loadData();
                },
                child: const Text('등록하기'),
              ),
            ],
          ),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('내 자산'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add, color: AppColors.textPrimary),
            onPressed: _showAddAccountSheet,
            tooltip: '계정 추가',
          ),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator(color: AppColors.brandPrimary))
          : RefreshIndicator(
              color: AppColors.brandPrimary,
              onRefresh: _loadData,
              child: SingleChildScrollView(
                physics: const AlwaysScrollableScrollPhysics(),
                padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Legal Banner (Article 1.5)
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
                              '표시된 자산과 잔액은 사용자가 직접 입력한 기록을 기준으로 합니다.',
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

                    // Net Worth Summary Card
                    Card(
                      child: Padding(
                        padding: const EdgeInsets.all(24.0),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text(
                              '순자산 총액',
                              style: TextStyle(
                                fontSize: 14,
                                fontWeight: FontWeight.w600,
                                color: AppColors.textSecondary,
                              ),
                            ),
                            const SizedBox(height: 8),
                            Text(
                              CurrencyFormatter.format(_totalNetWorth),
                              style: const TextStyle(
                                fontSize: 30,
                                fontWeight: FontWeight.w800,
                                color: AppColors.textPrimary,
                                letterSpacing: -0.8,
                              ),
                            ),
                            const SizedBox(height: 14),
                            Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                Text(
                                  '총 자산 ${CurrencyFormatter.format(_totalAssets)}',
                                  style: const TextStyle(fontSize: 13, color: AppColors.textSecondary, fontWeight: FontWeight.w600),
                                ),
                                Text(
                                  '부채 ${CurrencyFormatter.format(_totalLiabilities)}',
                                  style: const TextStyle(fontSize: 13, color: AppColors.danger, fontWeight: FontWeight.w600),
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),
                    ),

                    const SizedBox(height: 12),

                    // Financial Accounts Card (Toss Account List)
                    Card(
                      child: Padding(
                        padding: const EdgeInsets.all(20.0),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                Text(
                                  '금융 계정 (${_accounts.length})',
                                  style: const TextStyle(
                                    fontWeight: FontWeight.w700,
                                    fontSize: 16,
                                    color: AppColors.textPrimary,
                                  ),
                                ),
                                InkWell(
                                  onTap: _showAddAccountSheet,
                                  borderRadius: BorderRadius.circular(8),
                                  child: const Padding(
                                    padding: EdgeInsets.symmetric(horizontal: 6, vertical: 4),
                                    child: Text(
                                      '+ 계정 추가',
                                      style: TextStyle(
                                        fontSize: 13,
                                        fontWeight: FontWeight.w700,
                                        color: AppColors.brandPrimary,
                                      ),
                                    ),
                                  ),
                                ),
                              ],
                            ),
                            const SizedBox(height: 14),
                            if (_accounts.isEmpty)
                              const Padding(
                                padding: EdgeInsets.symmetric(vertical: 24.0),
                                child: Center(
                                  child: Text('등록된 계정이 없습니다.', style: TextStyle(color: AppColors.textTertiary, fontSize: 13)),
                                ),
                              )
                            else
                              ListView.separated(
                                shrinkWrap: true,
                                physics: const NeverScrollableScrollPhysics(),
                                itemCount: _accounts.length,
                                separatorBuilder: (_, __) => const Divider(height: 18),
                                itemBuilder: (context, index) {
                                  final a = _accounts[index];
                                  final typeName = a['account_type'] == 'MANUAL_BANK'
                                      ? '수동 계좌'
                                      : a['account_type'] == 'CASH'
                                          ? '현금 지갑'
                                          : a['account_type'] == 'SAVINGS'
                                              ? '저축/적금'
                                              : '간편결제';
                                  return Row(
                                    children: [
                                      Container(
                                        width: 42,
                                        height: 42,
                                        decoration: BoxDecoration(
                                          color: AppColors.primaryLight,
                                          borderRadius: BorderRadius.circular(15),
                                        ),
                                        child: const Icon(
                                          Icons.account_balance_wallet_rounded,
                                          color: AppColors.brandPrimary,
                                          size: 22,
                                        ),
                                      ),
                                      const SizedBox(width: 14),
                                      Expanded(
                                        child: Column(
                                          crossAxisAlignment: CrossAxisAlignment.start,
                                          children: [
                                            Text(
                                              a['name'] ?? '',
                                              style: const TextStyle(
                                                fontWeight: FontWeight.w700,
                                                fontSize: 15,
                                                color: AppColors.textPrimary,
                                              ),
                                            ),
                                            const SizedBox(height: 2),
                                            Text(
                                              '$typeName · ${a['currency_code']}',
                                              style: const TextStyle(
                                                fontSize: 12,
                                                color: AppColors.textTertiary,
                                              ),
                                            ),
                                          ],
                                        ),
                                      ),
                                      Text(
                                        CurrencyFormatter.format(a['calculated_balance_minor'] ?? 0, currency: a['currency_code'] ?? 'KRW'),
                                        style: const TextStyle(
                                          fontWeight: FontWeight.w800,
                                          fontSize: 16,
                                          color: AppColors.textPrimary,
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

                    const SizedBox(height: 12),

                    // Other Assets & Liabilities Card
                    Card(
                      child: Padding(
                        padding: const EdgeInsets.all(20.0),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text(
                              '기타 자산 및 부채',
                              style: TextStyle(
                                fontWeight: FontWeight.w700,
                                fontSize: 16,
                                color: AppColors.textPrimary,
                              ),
                            ),
                            const SizedBox(height: 14),
                            if (_assets.isEmpty && _liabilities.isEmpty)
                              const Padding(
                                padding: EdgeInsets.symmetric(vertical: 20.0),
                                child: Center(
                                  child: Text('등록된 기타 자산 및 부채가 없습니다.', style: TextStyle(color: AppColors.textTertiary, fontSize: 13)),
                                ),
                              )
                            else ...[
                              ..._assets.map((ast) => Padding(
                                    padding: const EdgeInsets.symmetric(vertical: 8.0),
                                    child: Row(
                                      children: [
                                        Container(
                                          width: 38,
                                          height: 38,
                                          decoration: BoxDecoration(
                                            color: const Color(0xFFF1EDFF),
                                            borderRadius: BorderRadius.circular(13),
                                          ),
                                          child: const Icon(Icons.home_work_rounded, color: Color(0xFF8061D9), size: 20),
                                        ),
                                        const SizedBox(width: 12),
                                        Expanded(
                                          child: Column(
                                            crossAxisAlignment: CrossAxisAlignment.start,
                                            children: [
                                              Text(ast['name'] ?? '', style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 14)),
                                              Text(ast['asset_type'] ?? '기타 자산', style: const TextStyle(fontSize: 12, color: AppColors.textTertiary)),
                                            ],
                                          ),
                                        ),
                                        Text(
                                          CurrencyFormatter.format(ast['estimated_value_minor'] ?? 0),
                                          style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 15),
                                        ),
                                      ],
                                    ),
                                  )),
                              ..._liabilities.map((liab) => Padding(
                                    padding: const EdgeInsets.symmetric(vertical: 8.0),
                                    child: Row(
                                      children: [
                                        Container(
                                          width: 38,
                                          height: 38,
                                          decoration: BoxDecoration(
                                            color: const Color(0xFFFEECEB),
                                            borderRadius: BorderRadius.circular(13),
                                          ),
                                          child: const Icon(Icons.credit_card_rounded, color: AppColors.danger, size: 20),
                                        ),
                                        const SizedBox(width: 12),
                                        Expanded(
                                          child: Column(
                                            crossAxisAlignment: CrossAxisAlignment.start,
                                            children: [
                                              Text(liab['name'] ?? '', style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 14)),
                                              Text(liab['liability_type'] ?? '부채', style: const TextStyle(fontSize: 12, color: AppColors.textTertiary)),
                                            ],
                                          ),
                                        ),
                                        Text(
                                          CurrencyFormatter.format(liab['outstanding_amount_minor'] ?? 0),
                                          style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 15, color: AppColors.danger),
                                        ),
                                      ],
                                    ),
                                  )),
                            ],
                          ],
                        ),
                      ),
                    ),

                    const SizedBox(height: 24),
                  ],
                ),
              ),
            ),
    );
  }
}
