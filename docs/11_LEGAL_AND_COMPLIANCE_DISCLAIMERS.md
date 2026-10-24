# 11. LEGAL AND COMPLIANCE DISCLAIMERS SPECIFICATION

**Document ID:** BN-ONE-LEG-011  
**Version:** 1.0.0  
**Effective Date:** 2026-10-24  

---

## 1. Regulatory Context & Legal Boundary

BN ONE v1.0 is engineered and marketed strictly as a **Personal Financial Information & Ledger Management Tool**. It does not perform money transmission, open banking aggregation, deposit-taking, or investment advisory services.

To ensure strict compliance with Korean Financial Services Commission (FSC) guidelines, FTC consumer protection laws, and global fintech regulations, the application displays mandatory, non-dismissible disclaimers across critical surfaces.

---

## 2. Mandatory Disclaimer Matrix

| Surface | Required Disclosure Text (KO / EN) | Placement & Visual Requirement |
|---|---|---|
| **Asset & Account Management** | `표시된 자산과 잔액은 사용자가 직접 입력한 기록을 기준으로 합니다.`<br>*(Displayed assets and balances are based on records manually entered by the user.)* | Prominently displayed above the Net Worth card and on account creation modals; high-contrast amber/indigo banner with info icon. |
| **Exchange Rate (Forex)** | `표시 환율은 참고용이며 실제 환전·결제 환율과 다를 수 있습니다.`<br>*(Displayed exchange rates are for reference only and may differ from actual settlement rates.)* | Located directly above the live currency rate board and currency conversion calculator. |
| **Transaction History** | `BN ONE의 거래 기록은 실제 금융기관 거래를 실행하지 않습니다.`<br>*(BN ONE transaction records do not execute transactions with actual financial institutions.)* | Displayed on the transaction history feed and new transaction logging screen. |
| **Subscriptions** | `구독 갱신 알림은 참고 정보이며, 실제 결제일 및 청구 금액은 해당 서비스 제공자의 정책에 따릅니다.`<br>*(Subscription renewal notices are informational; actual billing dates depend on service provider policies.)* | Displayed in the subscription manager screen. |

---

## 3. Prohibited Representation Rules

1. **No Implication of Bank Synchronization:** Marketing copy and UI text must never use phrasing such as "Bank auto-sync", "Live account balance connection", "Direct wire transfer", or "Open banking partner".
2. **Clear Origin Labels:** In all exports and reports, figures are labeled as `User-Entered Ledger` rather than `Bank Verified`.
3. **No Custodial Claims:** Disclaimers explicitly specify that Booking Network Co., Ltd. does not hold, custody, or insure user funds.
