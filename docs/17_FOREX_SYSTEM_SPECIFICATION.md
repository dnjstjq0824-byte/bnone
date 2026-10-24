# 17. REFERENCE FOREX SYSTEM SPECIFICATION

**Document ID:** BN-ONE-FX-017  
**Version:** 1.0.0  
**Effective Date:** 2026-10-24  

---

## 1. System Objective & Functional Scope

The BN ONE Foreign Exchange (Forex) Engine provides **informational reference exchange rates** for over 30 global fiat currencies to enable:
1. **Multi-Currency Portfolio Valuation:** Calculating a user's total Net Worth in their chosen primary base currency (e.g., converting JPY, USD, and EUR accounts into a consolidated KRW figure).
2. **Travel Expense Conversion:** Converting foreign travel expenditures into home currency for post-trip budgeting and settlement.
3. **Reference Currency Converter:** A handy on-the-go currency conversion calculator for travelers.

> **CRITICAL COMPLIANCE NOTICE:**  
> BN ONE is **NOT** a currency exchange broker, remittance provider, or forex trading platform. All displayed rates are reference mid-market rates and explicitly disclaimed as non-binding for actual retail currency transactions.

---

## 2. Reference Currency Support List

The system supports major trading pairs and regional travel currencies:
- **Major Reserve Currencies:** USD, EUR, JPY, GBP, CHF, CAD, AUD.
- **Asia-Pacific & Travel Currencies:** KRW, CNY, HKD, TWD, SGD, THB, VND, PHP, MYR, IDR.
- **Other Global Currencies:** NZD, INR, AED, SAR, TRY, BRL, MXN, ZAR, SEK, NOK, DKK, PLN, CZK, HUF.

---

## 3. Data Ingestion Architecture & Fallback Strategy

```text
External Central Bank / ECB / KEB Hana Feeds
                      │
                      ▼
        services/worker/internal/fx
         (Hourly Ingestion Worker)
                      │
                      ▼
     PostgreSQL table: exchange_rates
 (base_currency, target_currency, rate, updated_at)
                      │
                      ▼
        services/api/internal/forex
    (In-memory Read Cache with 1hr TTL)
                      │
                      ▼
             Client Mobile & Web
```

### 3.1 Worker Scheduling & Retries
- Ingestion runs on a configurable cron schedule (`@hourly`).
- If an external upstream API fails or times out, the worker implements exponential backoff retry (3 attempts).
- The cached database rates remain active, annotated with the exact `updated_at` timestamp.

