import 'package:flutter/material.dart';
import '../../../core/api/api_client.dart';
import '../../../core/theme/colors.dart';
import '../../../core/utils/formatters.dart';

class ForexScreen extends StatefulWidget {
  const ForexScreen({super.key});

  @override
  State<ForexScreen> createState() => _ForexScreenState();
}

class _ForexScreenState extends State<ForexScreen> {
  bool _isLoading = true;
  List<dynamic> _rates = [];
  final _amountCtrl = TextEditingController(text: '100');
  String _sourceCurrency = 'USD';
  String _targetCurrency = 'KRW';
  Map<String, dynamic>? _convertResult;

  @override
  void initState() {
    super.initState();
    _loadRates();
  }

  Future<void> _loadRates() async {
    setState(() => _isLoading = true);
    try {
      final res = await ApiClient().get('/fx/rates?base=USD');
      setState(() => _rates = res?['rates'] ?? []);
    } catch (e) {
      debugPrint('FX load error: $e');
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  Future<void> _calculateConversion() async {
    final amount = int.tryParse(_amountCtrl.text.replaceAll(',', '').trim()) ?? 0;
    try {
      final res = await ApiClient().post('/fx/convert', {
        'source_currency': _sourceCurrency,
        'target_currency': _targetCurrency,
        'source_amount_minor': amount,
      });
      setState(() => _convertResult = res);
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(e.toString())));
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(title: const Text('공인 참고 환율')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Mandatory Disclaimer (Article 1.5 & Article 18)
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
              margin: const EdgeInsets.only(bottom: 12),
              decoration: BoxDecoration(
                color: const Color(0xFFFFF3E2),
                borderRadius: BorderRadius.circular(16),
              ),
              child: const Row(
                children: [
                  Icon(Icons.info_outline, size: 16, color: Color(0xFFFF8D00)),
                  SizedBox(width: 10),
                  Expanded(
                    child: Text(
                      '표시 환율은 참고용이며 실제 환전·결제 환율과 다를 수 있습니다.',
                      style: TextStyle(fontSize: 12, color: Color(0xFFD88924), fontWeight: FontWeight.bold),
                    ),
                  ),
                ],
              ),
            ),

            // Currency Converter Card (Toss Style)
            Card(
              child: Padding(
                padding: const EdgeInsets.all(22.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      '환율 계산기',
                      style: TextStyle(fontWeight: FontWeight.w800, fontSize: 17, color: AppColors.textPrimary),
                    ),
                    const SizedBox(height: 16),
                    Row(
                      children: [
                        Expanded(
                          child: TextField(
                            controller: _amountCtrl,
                            keyboardType: TextInputType.number,
                            decoration: const InputDecoration(labelText: '환산할 금액'),
                          ),
                        ),
                        const SizedBox(width: 12),
                        Container(
                          padding: const EdgeInsets.symmetric(horizontal: 12),
                          decoration: BoxDecoration(
                            color: Colors.white,
                            borderRadius: BorderRadius.circular(16),
                            border: Border.all(color: AppColors.border),
                          ),
                          child: DropdownButtonHideUnderline(
                            child: DropdownButton<String>(
                              value: _sourceCurrency,
                              items: ['USD', 'KRW', 'EUR', 'JPY'].map((c) => DropdownMenuItem(value: c, child: Text(c, style: const TextStyle(fontWeight: FontWeight.w700)))).toList(),
                              onChanged: (v) => setState(() => _sourceCurrency = v!),
                            ),
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 12),
                    Row(
                      children: [
                        const Text('환산 대상 통화: ', style: TextStyle(fontSize: 13, color: AppColors.textSecondary)),
                        Container(
                          padding: const EdgeInsets.symmetric(horizontal: 12),
                          decoration: BoxDecoration(
                            color: Colors.white,
                            borderRadius: BorderRadius.circular(16),
                            border: Border.all(color: AppColors.border),
                          ),
                          child: DropdownButtonHideUnderline(
                            child: DropdownButton<String>(
                              value: _targetCurrency,
                              items: ['KRW', 'USD', 'EUR', 'JPY'].map((c) => DropdownMenuItem(value: c, child: Text(c, style: const TextStyle(fontWeight: FontWeight.w700)))).toList(),
                              onChanged: (v) => setState(() => _targetCurrency = v!),
                            ),
                          ),
                        ),
                        const Spacer(),
                        ElevatedButton(
                          style: ElevatedButton.styleFrom(
                            minimumSize: const Size(100, 48),
                            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
                          ),
                          onPressed: _calculateConversion,
                          child: const Text('계산'),
                        ),
                      ],
                    ),
                    if (_convertResult != null) ...[
                      const SizedBox(height: 18),
                      Container(
                        width: double.infinity,
                        padding: const EdgeInsets.all(16),
                        decoration: BoxDecoration(
                          color: AppColors.primaryLight,
                          borderRadius: BorderRadius.circular(16),
                        ),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text('환산 결과', style: TextStyle(fontSize: 12, color: AppColors.textSecondary, fontWeight: FontWeight.bold)),
                            const SizedBox(height: 4),
                            Text(
                              '${CurrencyFormatter.format(_convertResult!['target_amount_minor'] ?? 0, currency: _targetCurrency)}',
                              style: const TextStyle(fontSize: 24, fontWeight: FontWeight.w800, color: AppColors.brandPrimary),
                            ),
                            const SizedBox(height: 4),
                            Text(
                              '적용 환율: 1 $_sourceCurrency = ${_convertResult!['exchange_rate']} $_targetCurrency',
                              style: const TextStyle(fontSize: 12, color: AppColors.textTertiary),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ],
                ),
              ),
            ),

            const SizedBox(height: 12),

            // Live Rates Table Card
            Card(
              child: Padding(
                padding: const EdgeInsets.all(20.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      '주요 통화 참고 환율 (기준: USD)',
                      style: TextStyle(fontWeight: FontWeight.w700, fontSize: 16, color: AppColors.textPrimary),
                    ),
                    const SizedBox(height: 14),
                    if (_isLoading)
                      const Center(child: Padding(padding: EdgeInsets.all(24), child: CircularProgressIndicator(color: AppColors.brandPrimary)))
                    else
                      ListView.separated(
                        shrinkWrap: true,
                        physics: const NeverScrollableScrollPhysics(),
                        itemCount: _rates.length,
                        separatorBuilder: (_, __) => const Divider(height: 18),
                        itemBuilder: (context, index) {
                          final r = _rates[index];
                          return Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Text(
                                '${r['target_currency']}',
                                style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 15, color: AppColors.textPrimary),
                              ),
                              Text(
                                '${r['rate']}',
                                style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 15, color: AppColors.brandPrimary),
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
    );
  }
}
