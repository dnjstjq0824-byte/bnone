import 'package:flutter/material.dart';
import '../../../core/api/api_client.dart';
import '../../../core/theme/colors.dart';

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
    final amount = int.tryParse(_amountCtrl.text.trim()) ?? 0;
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
      appBar: AppBar(title: const Text('공인 참고 환율')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Mandatory Disclaimer (Article 1.5 & Article 18)
            Container(
              padding: const EdgeInsets.all(12),
              margin: const EdgeInsets.only(bottom: 16),
              decoration: BoxDecoration(
                color: AppColors.warning.withOpacity(0.1),
                borderRadius: BorderRadius.circular(12),
              ),
              child: const Text(
                '표시 환율은 참고용이며 실제 환전·결제 환율과 다를 수 있습니다.',
                style: TextStyle(fontSize: 12, color: AppColors.warning, fontWeight: FontWeight.bold),
              ),
            ),

            // Calculator Card
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text('환산 계산기', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                    const SizedBox(height: 12),
                    Row(
                      children: [
                        Expanded(
                          child: TextField(
                            controller: _amountCtrl,
                            keyboardType: TextInputType.number,
                            decoration: const InputDecoration(labelText: '금액'),
                          ),
                        ),
                        const SizedBox(width: 8),
                        DropdownButton<String>(
                          value: _sourceCurrency,
                          items: ['USD', 'KRW', 'EUR', 'JPY'].map((c) => DropdownMenuItem(value: c, child: Text(c))).toList(),
                          onChanged: (v) => setState(() => _sourceCurrency = v!),
                        ),
                        const Padding(padding: EdgeInsets.symmetric(horizontal: 8), child: Text('→')),
                        DropdownButton<String>(
                          value: _targetCurrency,
                          items: ['KRW', 'USD', 'EUR', 'JPY'].map((c) => DropdownMenuItem(value: c, child: Text(c))).toList(),
                          onChanged: (v) => setState(() => _targetCurrency = v!),
                        ),
                      ],
                    ),
                    const SizedBox(height: 16),
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton(onPressed: _calculateConversion, child: const Text('환산하기')),
                    ),
                    if (_convertResult != null) ...[
                      const SizedBox(height: 12),
                      Container(
                        padding: const EdgeInsets.all(12),
                        decoration: BoxDecoration(color: AppColors.background, borderRadius: BorderRadius.circular(8)),
                        child: Text(
                          '결과: ${_convertResult!['converted_amount_minor']} $_targetCurrency (적용 환율: ${_convertResult!['rate']})',
                          style: const TextStyle(fontWeight: FontWeight.bold),
                        ),
                      ),
                    ],
                  ],
                ),
              ),
            ),
            const SizedBox(height: 24),

            const Text('주요 통화 참고 환율 (기준: USD)', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
            const SizedBox(height: 8),

            _isLoading
                ? const Center(child: CircularProgressIndicator())
                : ListView.builder(
                    shrinkWrap: true,
                    physics: const NeverScrollableScrollPhysics(),
                    itemCount: _rates.length,
                    itemBuilder: (context, index) {
                      final r = _rates[index];
                      return Card(
                        margin: const EdgeInsets.only(bottom: 8),
                        child: ListTile(
                          title: Text('1 ${r['base_currency']} = ${r['rate']} ${r['quote_currency']}'),
                          subtitle: Text('출처: ${r['provider_name']}'),
                          trailing: const Text('VERIFIED', style: TextStyle(color: AppColors.success, fontWeight: FontWeight.bold, fontSize: 11)),
                        ),
                      );
                    },
                  ),
          ],
        ),
      ),
    );
  }
}
