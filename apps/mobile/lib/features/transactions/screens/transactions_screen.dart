import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../core/api/api_client.dart';
import '../../../core/theme/colors.dart';
import '../../../core/utils/formatters.dart';

class TransactionsScreen extends StatefulWidget {
  const TransactionsScreen({super.key});

  @override
  State<TransactionsScreen> createState() => _TransactionsScreenState();
}

class _TransactionsScreenState extends State<TransactionsScreen> {
  bool _isLoading = true;
  List<dynamic> _transactions = [];
  String _filterType = '';

  @override
  void initState() {
    super.initState();
    _loadTransactions();
  }

  Future<void> _loadTransactions() async {
    setState(() => _isLoading = true);
    try {
      final endpoint = _filterType.isNotEmpty ? '/transactions?type=$_filterType' : '/transactions';
      final list = await ApiClient().get(endpoint);
      setState(() => _transactions = list ?? []);
    } catch (e) {
      debugPrint('Transactions error: $e');
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  Future<void> _deleteTx(String id) async {
    if (await showDialog<bool>(
          context: context,
          builder: (ctx) => AlertDialog(
            title: const Text('거래 삭제'),
            content: const Text('해당 거래를 삭제하시겠습니까? 잔액과 예산이 자동 재계산됩니다.'),
            actions: [
              TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('취소')),
              ElevatedButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('삭제')),
            ],
          ),
        ) ==
        true) {
      await ApiClient().delete('/transactions/$id');
      _loadTransactions();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('거래 내역'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () => context.push('/transactions/new'),
          ),
        ],
      ),
      body: Column(
        children: [
          // Filter Chips
          SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            child: Row(
              children: [
                _buildFilterChip('전체', ''),
                const SizedBox(width: 8),
                _buildFilterChip('지출', 'EXPENSE'),
                const SizedBox(width: 8),
                _buildFilterChip('수입', 'INCOME'),
                const SizedBox(width: 8),
                _buildFilterChip('계정 이동', 'TRANSFER'),
              ],
            ),
          ),

          Expanded(
            child: _isLoading
                ? const Center(child: CircularProgressIndicator())
                : _transactions.isEmpty
                    ? const Center(child: Text('거래 내역이 없습니다.', style: TextStyle(color: AppColors.textSecondary)))
                    : ListView.builder(
                        padding: const EdgeInsets.all(16),
                        itemCount: _transactions.length,
                        itemBuilder: (context, index) {
                          final tx = _transactions[index];
                          final isExpense = tx['transaction_type'] == 'EXPENSE';
                          final isIncome = tx['transaction_type'] == 'INCOME';

                          return Card(
                            margin: const EdgeInsets.only(bottom: 8),
                            child: ListTile(
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
                                '${tx['occurred_at']?.toString().split("T")[0]} • ${tx['category_name'] ?? tx['transaction_type']}',
                                style: const TextStyle(fontSize: 11, color: AppColors.textSecondary),
                              ),
                              trailing: Text(
                                '${isExpense ? "− " : isIncome ? "+ " : ""}${CurrencyFormatter.format(tx['amount_minor'] ?? 0, currency: tx['currency_code'] ?? 'KRW')}',
                                style: TextStyle(
                                  fontWeight: FontWeight.w900,
                                  fontSize: 14,
                                  color: isExpense ? AppColors.textPrimary : isIncome ? AppColors.success : AppColors.brandPrimary,
                                ),
                              ),
                              onLongPress: () => _deleteTx(tx['id']),
                            ),
                          );
                        },
                      ),
          ),
        ],
      ),
    );
  }

  Widget _buildFilterChip(String label, String value) {
    final selected = _filterType == value;
    return ChoiceChip(
      label: Text(label),
      selected: selected,
      onSelected: (_) {
        setState(() => _filterType = value);
        _loadTransactions();
      },
      selectedColor: AppColors.brandPrimary.withOpacity(0.15),
      labelStyle: TextStyle(
        color: selected ? AppColors.brandPrimary : AppColors.textSecondary,
        fontWeight: selected ? FontWeight.bold : FontWeight.normal,
        fontSize: 12,
      ),
    );
  }
}
