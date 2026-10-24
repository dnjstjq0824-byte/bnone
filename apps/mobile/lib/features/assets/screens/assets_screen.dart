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

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    setState(() => _isLoading = true);
    try {
      final accs = await ApiClient().get('/accounts');
      final asts = await ApiClient().get('/assets');
      final liabs = await ApiClient().get('/liabilities');
      setState(() {
        _accounts = accs ?? [];
        _assets = asts ?? [];
        _liabilities = liabs ?? [];
      });
    } catch (e) {
      debugPrint('Assets load error: $e');
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  void _showAddAccountDialog() {
    final nameCtrl = TextEditingController();
    final balCtrl = TextEditingController(text: '0');
    String type = 'MANUAL_BANK';
    String curr = 'KRW';

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('새 금융계정 추가', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(controller: nameCtrl, decoration: const InputDecoration(labelText: '계정 이름')),
            const SizedBox(height: 12),
            DropdownButtonFormField<String>(
              value: type,
              items: const [
                DropdownMenuItem(value: 'MANUAL_BANK', child: Text('은행계좌형 수동 장부')),
                DropdownMenuItem(value: 'CASH', child: Text('현금 지갑')),
                DropdownMenuItem(value: 'SAVINGS', child: Text('저축 / 적금')),
                DropdownMenuItem(value: 'E_WALLET_MANUAL', child: Text('전자지갑')),
              ],
              onChanged: (v) => type = v!,
              decoration: const InputDecoration(labelText: '유형'),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: balCtrl,
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(labelText: '시작 잔액'),
            ),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
          ElevatedButton(
            onPressed: () async {
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
            child: const Text('추가'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('수동 자산 & 부채 장부'),
        actions: [
          IconButton(icon: const Icon(Icons.add), onPressed: _showAddAccountDialog),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : RefreshIndicator(
              onRefresh: _loadData,
              child: SingleChildScrollView(
                physics: const AlwaysScrollableScrollPhysics(),
                padding: const EdgeInsets.all(16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Legal Banner
                    Container(
                      padding: const EdgeInsets.all(12),
                      decoration: BoxDecoration(
                        color: AppColors.brandPrimary.withOpacity(0.08),
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: const Text(
                        '표시된 자산과 잔액은 사용자가 직접 입력한 기록을 기준으로 합니다.',
                        style: TextStyle(fontSize: 11, color: AppColors.brandPrimary, fontWeight: FontWeight.bold),
                      ),
                    ),
                    const SizedBox(height: 20),

                    // Financial Accounts Section
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text('금융 계정 (${_accounts.length})', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                        TextButton(onPressed: _showAddAccountDialog, child: const Text('+ 계정 추가')),
                      ],
                    ),
                    const SizedBox(height: 8),

                    ListView.builder(
                      shrinkWrap: true,
                      physics: const NeverScrollableScrollPhysics(),
                      itemCount: _accounts.length,
                      itemBuilder: (context, index) {
                        final a = _accounts[index];
                        return Card(
                          margin: const EdgeInsets.only(bottom: 10),
                          child: ListTile(
                            leading: CircleAvatar(
                              backgroundColor: AppColors.brandPrimary.withOpacity(0.1),
                              child: const Icon(Icons.account_balance_wallet, color: AppColors.brandPrimary, size: 20),
                            ),
                            title: Text(a['name'] ?? '', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
                            subtitle: Text('${a['account_type']} • ${a['currency_code']}', style: const TextStyle(fontSize: 11, color: AppColors.textSecondary)),
                            trailing: Text(
                              CurrencyFormatter.format(a['calculated_balance_minor'] ?? 0, currency: a['currency_code'] ?? 'KRW'),
                              style: const TextStyle(fontWeight: FontWeight.w900, fontSize: 15),
                            ),
                          ),
                        );
                      },
                    ),
                    const SizedBox(height: 24),

                    // Other Assets & Liabilities
                    const Text('기타 자산 및 부채', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                    const SizedBox(height: 8),

                    if (_assets.isEmpty && _liabilities.isEmpty)
                      const Card(
                        child: Padding(
                          padding: EdgeInsets.all(20.0),
                          child: Center(child: Text('등록된 기타 자산 및 부채가 없습니다.', style: TextStyle(fontSize: 12, color: AppColors.textSecondary))),
                        ),
                      )
                    else ...[
                      ..._assets.map((ast) => Card(
                            margin: const EdgeInsets.only(bottom: 8),
                            child: ListTile(
                              title: Text(ast['name'] ?? ''),
                              subtitle: Text(ast['asset_type'] ?? 'ASSET'),
                              trailing: Text(CurrencyFormatter.format(ast['estimated_value_minor'] ?? 0)),
                            ),
                          )),
                      ..._liabilities.map((liab) => Card(
                            margin: const EdgeInsets.only(bottom: 8),
                            child: ListTile(
                              title: Text(liab['name'] ?? ''),
                              subtitle: Text(liab['liability_type'] ?? 'LIABILITY'),
                              trailing: Text(
                                CurrencyFormatter.format(liab['outstanding_amount_minor'] ?? 0),
                                style: const TextStyle(color: AppColors.danger, fontWeight: FontWeight.bold),
                              ),
                            ),
                          )),
                    ],
                  ],
                ),
              ),
            ),
    );
  }
}

