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
    final confirmed = await showModalBottomSheet<bool>(
      context: context,
      backgroundColor: Colors.white,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (ctx) => Padding(
        padding: const EdgeInsets.all(24.0),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              '거래 내역 삭제',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.w800, color: AppColors.textPrimary),
            ),
            const SizedBox(height: 10),
            const Text(
              '이 거래 기록을 삭제하시겠습니까?\n연결된 계정 잔액과 예산이 자동으로 다시 계산됩니다.',
              style: TextStyle(fontSize: 14, color: AppColors.textSecondary, height: 1.5),
            ),
            const SizedBox(height: 24),
            Row(
              children: [
                Expanded(
                  child: OutlinedButton(
                    onPressed: () => Navigator.pop(ctx, false),
                    child: const Text('취소'),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: ElevatedButton(
                    style: ElevatedButton.styleFrom(backgroundColor: AppColors.danger),
                    onPressed: () => Navigator.pop(ctx, true),
                    child: const Text('삭제하기'),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );

    if (confirmed == true) {
      await ApiClient().delete('/transactions/$id');
      _loadTransactions();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('소비 내역'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add, color: AppColors.textPrimary),
            onPressed: () => context.push('/transactions/new'),
            tooltip: '새 거래 기록',
          ),
        ],
      ),
      body: Column(
        children: [
          // Filter Chips (Toss Pill Filter Style)
          Container(
            color: AppColors.background,
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
                ? const Center(child: CircularProgressIndicator(color: AppColors.brandPrimary))
                : _transactions.isEmpty
                    ? Center(
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Container(
                              width: 64,
                              height: 64,
                              decoration: BoxDecoration(
                                color: AppColors.primaryLight,
                                borderRadius: BorderRadius.circular(20),
                              ),
                              child: const Icon(Icons.receipt_long_rounded, color: AppColors.brandPrimary, size: 32),
                            ),
                            const SizedBox(height: 16),
                            const Text(
                              '조회된 거래 내역이 없습니다.',
                              style: TextStyle(color: AppColors.textSecondary, fontSize: 15, fontWeight: FontWeight.w600),
                            ),
                            const SizedBox(height: 12),
                            TextButton(
                              onPressed: () => context.push('/transactions/new'),
                              child: const Text('+ 새 거래 기록하기'),
                            ),
                          ],
                        ),
                      )
                    : RefreshIndicator(
                        color: AppColors.brandPrimary,
                        onRefresh: _loadTransactions,
                        child: ListView.builder(
                          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                          itemCount: _transactions.length,
                          itemBuilder: (context, index) {
                            final tx = _transactions[index];
                            final isExpense = tx['transaction_type'] == 'EXPENSE';
                            final isIncome = tx['transaction_type'] == 'INCOME';
                            final dateStr = tx['occurred_at']?.toString().split("T")[0] ?? '';

                            return InkWell(
                              onLongPress: () => _deleteTx(tx['id']),
                              borderRadius: BorderRadius.circular(20),
                              child: Card(
                                child: Padding(
                                  padding: const EdgeInsets.all(18.0),
                                  child: Row(
                                    children: [
                                      Container(
                                        width: 44,
                                        height: 44,
                                        decoration: BoxDecoration(
                                          color: isIncome
                                              ? AppColors.primaryLight
                                              : isExpense
                                                  ? AppColors.background
                                                  : const Color(0xFFF1EDFF),
                                          borderRadius: BorderRadius.circular(16),
                                        ),
                                        child: Icon(
                                          isIncome
                                              ? Icons.arrow_downward_rounded
                                              : isExpense
                                                  ? Icons.arrow_upward_rounded
                                                  : Icons.swap_horiz_rounded,
                                          color: isIncome
                                              ? AppColors.brandPrimary
                                              : isExpense
                                                  ? AppColors.textSecondary
                                                  : const Color(0xFF8061D9),
                                          size: 22,
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
                                                fontSize: 15,
                                                color: AppColors.textPrimary,
                                              ),
                                              maxLines: 1,
                                              overflow: TextOverflow.ellipsis,
                                            ),
                                            const SizedBox(height: 3),
                                            Text(
                                              '$dateStr · ${tx['category_name'] ?? (isIncome ? "수입" : isExpense ? "지출" : "계정 이동")}',
                                              style: const TextStyle(
                                                fontSize: 12,
                                                color: AppColors.textTertiary,
                                              ),
                                            ),
                                          ],
                                        ),
                                      ),
                                      Text(
                                        '${isExpense ? "− " : isIncome ? "+ " : ""}${CurrencyFormatter.format(tx['amount_minor'] ?? 0, currency: tx['currency_code'] ?? 'KRW')}',
                                        style: TextStyle(
                                          fontWeight: FontWeight.w800,
                                          fontSize: 16,
                                          color: isIncome
                                              ? AppColors.brandPrimary
                                              : AppColors.textPrimary,
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                              ),
                            );
                          },
                        ),
                      ),
          ),
        ],
      ),
    );
  }

  Widget _buildFilterChip(String label, String value) {
    final selected = _filterType == value;
    return InkWell(
      onTap: () {
        setState(() => _filterType = value);
        _loadTransactions();
      },
      borderRadius: BorderRadius.circular(20),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 150),
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
        decoration: BoxDecoration(
          color: selected ? AppColors.textPrimary : Colors.white,
          borderRadius: BorderRadius.circular(20),
        ),
        child: Text(
          label,
          style: TextStyle(
            color: selected ? Colors.white : AppColors.textSecondary,
            fontWeight: selected ? FontWeight.w700 : FontWeight.w500,
            fontSize: 13,
          ),
        ),
      ),
    );
  }
}
