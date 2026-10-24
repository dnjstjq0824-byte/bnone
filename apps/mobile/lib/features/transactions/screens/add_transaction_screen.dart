import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../core/api/api_client.dart';
import '../../../core/theme/colors.dart';

class AddTransactionScreen extends StatefulWidget {
  const AddTransactionScreen({super.key});

  @override
  State<AddTransactionScreen> createState() => _AddTransactionScreenState();
}

class _AddTransactionScreenState extends State<AddTransactionScreen> with SingleTickerProviderStateMixin {
  late TabController _tabController;
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
    _tabController = TabController(length: 3, vsync: this);
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

  Future<void> _submitTransaction() async {
    final amount = int.tryParse(_amountController.text.trim()) ?? 0;
    if (amount <= 0) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('금액을 올바르게 입력하세요.')));
      return;
    }

    setState(() => _isLoading = true);
    try {
      if (_tabController.index == 2) {
        // Manual Transfer
        if (_selectedAccount == _destAccount) {
          ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('출발 계정과 도착 계정이 같을 수 없습니다.')));
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
        final type = _tabController.index == 0 ? 'EXPENSE' : 'INCOME';
        await ApiClient().post('/transactions', {
          'transaction_type': type,
          'account_id': _selectedAccount,
          'amount_minor': amount,
          'description': _descController.text.trim(),
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
      appBar: AppBar(
        title: const Text('거래 등록'),
        bottom: TabBar(
          controller: _tabController,
          labelColor: AppColors.brandPrimary,
          indicatorColor: AppColors.brandPrimary,
          tabs: const [
            Tab(text: '지출'),
            Tab(text: '수입'),
            Tab(text: '계정 이동'),
          ],
        ),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Notice
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: AppColors.brandPrimary.withOpacity(0.08),
                borderRadius: BorderRadius.circular(12),
              ),
              child: const Text(
                'BN ONE의 거래 기록은 실제 금융기관 송금을 실행하지 않습니다.',
                style: TextStyle(fontSize: 11, color: AppColors.brandPrimary),
              ),
            ),
            const SizedBox(height: 24),

            TextField(
              controller: _amountController,
              keyboardType: TextInputType.number,
              style: const TextStyle(fontSize: 28, fontWeight: FontWeight.w900),
              decoration: InputDecoration(
                labelText: '금액 (KRW)',
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(14)),
              ),
            ),
            const SizedBox(height: 16),

            // Account Selector
            DropdownButtonFormField<String>(
              value: _selectedAccount,
              decoration: InputDecoration(
                labelText: _tabController.index == 2 ? '출발 계정' : '출금/입금 계정',
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(14)),
              ),
              items: _accounts.map((a) {
                return DropdownMenuItem<String>(
                  value: a['id'],
                  child: Text('${a['name']} (${a['currency_code']})'),
                );
              }).toList(),
              onChanged: (v) => setState(() => _selectedAccount = v),
            ),
            const SizedBox(height: 16),

            // If Transfer: Destination Account
            if (_tabController.index == 2) ...[
              DropdownButtonFormField<String>(
                value: _destAccount,
                decoration: InputDecoration(
                  labelText: '도착 계정',
                  border: OutlineInputBorder(borderRadius: BorderRadius.circular(14)),
                ),
                items: _accounts.map((a) {
                  return DropdownMenuItem<String>(
                    value: a['id'],
                    child: Text('${a['name']} (${a['currency_code']})'),
                  );
                }).toList(),
                onChanged: (v) => setState(() => _destAccount = v),
              ),
              const SizedBox(height: 16),
            ] else ...[
              // Category
              DropdownButtonFormField<String>(
                value: _selectedCategory,
                decoration: InputDecoration(
                  labelText: '카테고리',
                  border: OutlineInputBorder(borderRadius: BorderRadius.circular(14)),
                ),
                items: _categories.map((c) {
                  return DropdownMenuItem<String>(
                    value: c['id'],
                    child: Text(c['name'] ?? ''),
                  );
                }).toList(),
                onChanged: (v) => setState(() => _selectedCategory = v),
              ),
              const SizedBox(height: 16),
            ],

            TextField(
              controller: _descController,
              decoration: InputDecoration(
                labelText: '설명',
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(14)),
              ),
            ),
            const SizedBox(height: 32),

            SizedBox(
              width: double.infinity,
              height: 52,
              child: ElevatedButton(
                onPressed: _isLoading ? null : _submitTransaction,
                style: ElevatedButton.styleFrom(
                  backgroundColor: AppColors.brandPrimary,
                  foregroundColor: Colors.white,
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
                ),
                child: _isLoading
                    ? const CircularProgressIndicator(color: Colors.white)
                    : const Text('기록 저장하기', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
