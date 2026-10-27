import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../core/api/api_client.dart';
import '../../../core/theme/colors.dart';

class AddTransactionScreen extends StatefulWidget {
  const AddTransactionScreen({super.key});

  @override
  State<AddTransactionScreen> createState() => _AddTransactionScreenState();
}

class _AddTransactionScreenState extends State<AddTransactionScreen> {
  int _selectedTypeIndex = 0; // 0: 지출, 1: 수입, 2: 계정 이동
  final _amountController = TextEditingController();
  final _descController = TextEditingController();
  final _merchantController = TextEditingController();

  List<dynamic> _accounts = [];
  List<dynamic> _categories = [];
  String? _selectedAccount;
  String? _selectedCategory;
  String? _destAccount;
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _loadDependencies();
  }

  Future<void> _loadDependencies() async {
    try {
      final accs = await ApiClient().get('/accounts');
      final cats = await ApiClient().get('/categories');
      setState(() {
        _accounts = accs ?? [];
        _categories = cats ?? [];
        if (_accounts.isNotEmpty) {
          _selectedAccount = _accounts[0]['id'];
          if (_accounts.length > 1) {
            _destAccount = _accounts[1]['id'];
          }
        }
        if (_categories.isNotEmpty) {
          _selectedCategory = _categories[0]['id'];
        }
      });
    } catch (e) {
      debugPrint('Load deps error: $e');
    }
  }

  void _addQuickAmount(int delta) {
    final current = int.tryParse(_amountController.text.replaceAll(',', '')) ?? 0;
    final next = current + delta;
    _amountController.text = next.toString();
  }

  Future<void> _submitTransaction() async {
    final amount = int.tryParse(_amountController.text.replaceAll(',', '').trim()) ?? 0;
    if (amount <= 0) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('금액을 올바르게 입력해주세요.')),
      );
      return;
    }

    setState(() => _isLoading = true);
    try {
      if (_selectedTypeIndex == 2) {
        // Manual Transfer
        if (_selectedAccount == _destAccount) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('출발 계정과 도착 계정이 같을 수 없습니다.')),
          );
          return;
        }
        await ApiClient().post('/transfers/manual', {
          'source_account_id': _selectedAccount,
          'destination_account_id': _destAccount,
          'source_amount_minor': amount,
          'destination_amount_minor': amount,
          'description': _descController.text.trim().isEmpty ? '수동 계정 이동' : _descController.text.trim(),
        });
      } else {
        // Expense or Income
        final type = _selectedTypeIndex == 0 ? 'EXPENSE' : 'INCOME';
        await ApiClient().post('/transactions', {
          'transaction_type': type,
          'account_id': _selectedAccount,
          'amount_minor': amount,
          'description': _descController.text.trim().isEmpty
              ? (_selectedTypeIndex == 0 ? '지출 기록' : '수입 기록')
              : _descController.text.trim(),
          'merchant_name': _merchantController.text.trim(),
          'splits': [
            {
              'category_id': _selectedCategory,
              'amount_minor': amount,
            }
          ],
        });
      }

      if (mounted) {
        context.pop();
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(e.toString())));
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('거래 기록'),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Type Segment Control (Toss Pill Selector)
            Container(
              padding: const EdgeInsets.all(4),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(16),
              ),
              child: Row(
                children: [
                  _buildSegmentButton('지출', 0),
                  _buildSegmentButton('수입', 1),
                  _buildSegmentButton('계정 이동', 2),
                ],
              ),
            ),

            const SizedBox(height: 16),

            // Amount Input Card (Toss Big Amount Display)
            Card(
              child: Padding(
                padding: const EdgeInsets.all(24.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      '금액 입력',
                      style: TextStyle(fontSize: 13, fontWeight: FontWeight.w600, color: AppColors.textTertiary),
                    ),
                    const SizedBox(height: 10),
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        Expanded(
                          child: TextField(
                            controller: _amountController,
                            keyboardType: TextInputType.number,
                            style: const TextStyle(
                              fontSize: 34,
                              fontWeight: FontWeight.w800,
                              color: AppColors.textPrimary,
                              letterSpacing: -1,
                            ),
                            decoration: const InputDecoration(
                              hintText: '0',
                              border: InputBorder.none,
                              enabledBorder: InputBorder.none,
                              focusedBorder: InputBorder.none,
                              filled: false,
                              contentPadding: EdgeInsets.zero,
                            ),
                          ),
                        ),
                        const Text(
                          '원',
                          style: TextStyle(
                            fontSize: 26,
                            fontWeight: FontWeight.w700,
                            color: AppColors.textPrimary,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 16),
                    // Quick Amount Presets
                    SingleChildScrollView(
                      scrollDirection: Axis.horizontal,
                      child: Row(
                        children: [
                          _buildPresetChip('+1만', 10000),
                          const SizedBox(width: 8),
                          _buildPresetChip('+5만', 50000),
                          const SizedBox(width: 8),
                          _buildPresetChip('+10만', 100000),
                          const SizedBox(width: 8),
                          _buildPresetChip('+50만', 500000),
                          const SizedBox(width: 8),
                          _buildPresetChip('+100만', 1000000),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),

            const SizedBox(height: 12),

            // Account & Category Fields Card
            Card(
              child: Padding(
                padding: const EdgeInsets.all(20.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    if (_selectedTypeIndex == 2) ...[
                      // Transfer Fields
                      DropdownButtonFormField<String>(
                        value: _selectedAccount,
                        items: _accounts
                            .map((a) => DropdownMenuItem<String>(
                                  value: a['id'] as String,
                                  child: Text('출발: ${a['name']} (${a['currency_code']})'),
                                ))
                            .toList(),
                        onChanged: (v) => setState(() => _selectedAccount = v),
                        decoration: const InputDecoration(labelText: '보내는 계정'),
                      ),
                      const SizedBox(height: 14),
                      DropdownButtonFormField<String>(
                        value: _destAccount,
                        items: _accounts
                            .map((a) => DropdownMenuItem<String>(
                                  value: a['id'] as String,
                                  child: Text('도착: ${a['name']} (${a['currency_code']})'),
                                ))
                            .toList(),
                        onChanged: (v) => setState(() => _destAccount = v),
                        decoration: const InputDecoration(labelText: '받는 계정'),
                      ),
                    ] else ...[
                      // Expense / Income Fields
                      DropdownButtonFormField<String>(
                        value: _selectedAccount,
                        items: _accounts
                            .map((a) => DropdownMenuItem<String>(
                                  value: a['id'] as String,
                                  child: Text('${a['name']} (${a['currency_code']})'),
                                ))
                            .toList(),
                        onChanged: (v) => setState(() => _selectedAccount = v),
                        decoration: const InputDecoration(labelText: '연결 계정'),
                      ),
                      const SizedBox(height: 14),
                      DropdownButtonFormField<String>(
                        value: _selectedCategory,
                        items: _categories
                            .map((c) => DropdownMenuItem<String>(
                                  value: c['id'] as String,
                                  child: Text(c['name'] as String),
                                ))
                            .toList(),
                        onChanged: (v) => setState(() => _selectedCategory = v),
                        decoration: const InputDecoration(labelText: '카테고리'),
                      ),
                    ],
                    const SizedBox(height: 14),
                    TextField(
                      controller: _descController,
                      decoration: const InputDecoration(
                        labelText: '내용 / 메모',
                        hintText: '예: 점심 식사, 스타벅스, 월급',
                      ),
                    ),
                    const SizedBox(height: 14),
                    TextField(
                      controller: _merchantController,
                      decoration: const InputDecoration(
                        labelText: '상호명 (선택)',
                        hintText: '예: 스타벅스 강남점',
                      ),
                    ),
                  ],
                ),
              ),
            ),

            const SizedBox(height: 20),

            // Mandatory Legal Disclaimer Note
            const Padding(
              padding: EdgeInsets.symmetric(horizontal: 4.0),
              child: Text(
                '※ BN ONE의 모든 기록은 사용자 직접 입력 장부이며 실제 금융기관 송금이나 결제를 실행하지 않습니다.',
                style: TextStyle(fontSize: 11, color: AppColors.textTertiary, height: 1.4),
              ),
            ),

            const SizedBox(height: 24),

            // Submit Button
            ElevatedButton(
              onPressed: _isLoading ? null : _submitTransaction,
              child: _isLoading
                  ? const SizedBox(
                      width: 20,
                      height: 20,
                      child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2),
                    )
                  : const Text('기록 완료'),
            ),

            const SizedBox(height: 32),
          ],
        ),
      ),
    );
  }

  Widget _buildSegmentButton(String label, int index) {
    final isSelected = _selectedTypeIndex == index;
    return Expanded(
      child: InkWell(
        onTap: () => setState(() => _selectedTypeIndex = index),
        borderRadius: BorderRadius.circular(12),
        child: AnimatedContainer(
          duration: const Duration(milliseconds: 150),
          padding: const EdgeInsets.symmetric(vertical: 12),
          decoration: BoxDecoration(
            color: isSelected ? AppColors.textPrimary : Colors.transparent,
            borderRadius: BorderRadius.circular(12),
          ),
          child: Center(
            child: Text(
              label,
              style: TextStyle(
                color: isSelected ? Colors.white : AppColors.textSecondary,
                fontWeight: isSelected ? FontWeight.w700 : FontWeight.w500,
                fontSize: 14,
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildPresetChip(String label, int delta) {
    return InkWell(
      onTap: () => _addQuickAmount(delta),
      borderRadius: BorderRadius.circular(12),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 8),
        decoration: BoxDecoration(
          color: AppColors.background,
          borderRadius: BorderRadius.circular(12),
        ),
        child: Text(
          label,
          style: const TextStyle(
            fontSize: 12,
            fontWeight: FontWeight.w700,
            color: AppColors.textSecondary,
          ),
        ),
      ),
    );
  }
}
