import 'package:intl/intl.dart';

class CurrencyFormatter {
  static String format(int minor, {String currency = 'KRW', bool masked = false}) {
    if (masked) return '${_getSymbol(currency)}••••••';
    final amount = (currency == 'KRW' || currency == 'JPY' || currency == 'VND')
        ? minor
        : minor / 100.0;

    final format = NumberFormat.currency(
      locale: 'ko_KR',
      symbol: _getSymbol(currency),
      decimalDigits: (currency == 'KRW' || currency == 'JPY') ? 0 : 2,
    );
    return format.format(amount);
  }

  static String _getSymbol(String currency) {
    switch (currency) {
      case 'KRW':
        return '₩';
      case 'USD':
        return '\$';
      case 'EUR':
        return '€';
      case 'JPY':
        return '¥';
      case 'GBP':
        return '£';
      default:
        return currency;
    }
  }
}
