import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:bn_one/app/app.dart';
import 'package:bn_one/core/utils/formatters.dart';

void main() {
  test('CurrencyFormatter minor units to major units formatting', () {
    expect(CurrencyFormatter.format(1500000, currency: 'KRW'), '₩1,500,000');
    expect(CurrencyFormatter.format(125050, currency: 'USD'), '\$1,250.50');
    expect(CurrencyFormatter.format(999, currency: 'JPY'), '¥999');
    expect(CurrencyFormatter.format(5000, currency: 'EUR'), '€50.00');
    expect(CurrencyFormatter.format(1500000, currency: 'KRW', masked: true), '₩••••••');
  });

  testWidgets('BN ONE App initializes smoke test', (WidgetTester tester) async {
    await tester.pumpWidget(
      const ProviderScope(
        child: BnOneApp(),
      ),
    );
    await tester.pump(const Duration(milliseconds: 100));

    expect(find.byType(BnOneApp), findsOneWidget);
  });
}
