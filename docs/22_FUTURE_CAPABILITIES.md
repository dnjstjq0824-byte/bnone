# 22. FUTURE CAPABILITIES & ROADMAP (v2.0+)

**Document ID:** BN-ONE-FTR-022  
**Version:** 1.0.0  
**Effective Date:** 2026-10-24  

---

## 1. Scope Boundary Isolation Rule

> **MANDATORY INVARIANT:**  
> The features outlined in this document are strictly scheduled for future versions (v2.0 and beyond). They are **NOT** active in BN ONE v1.0. The v1.0 codebase contains zero mock payment buttons, zero simulated bank logins, and zero fake transfer screens.

---

## 2. Future Capability Horizon (v2.0+)

### 2.1 Automated Open Banking Aggregation (v2.0)
- **Scope:** Read-only aggregation of bank account balances and credit card transactions via licensed Open Banking APIs (Korea Financial Telecommunications & Clearings Institute / KFTC, Plaid, or Yapily).
- **Prerequisites:** Acquisition of MyData / Account Information Service Provider (AISP) regulatory licensing, mutual TLS (mTLS) integration, user explicit consent management screens.

### 2.2 Camera-Based Receipt OCR & Automatic Parsing (v2.1)
- **Scope:** On-device camera capture of paper receipts (especially foreign language receipts in Japan, Europe, Southeast Asia) using localized OCR models to extract merchant name, date, currency, tax, and itemized lines.
- **Workflow:** User reviews OCR-extracted fields before confirming addition into the ledger.

### 2.3 Shared Group Travel Wallets (Split-the-Bill) (v2.2)
- **Scope:** Real-time multi-user synchronization for travel workspaces allowing groups of travelers to log joint expenses, calculate pairwise debt settlements (e.g., "Alice owes Bob ¥4,200"), and generate settlement summaries.

### 2.4 Multi-Asset Portfolio Tracker (v2.3)
- **Scope:** Real-time quotation tracking for public equities, ETFs, mutual funds, and digital assets using financial market data feeds, with automated unrealized gain/loss calculation.
