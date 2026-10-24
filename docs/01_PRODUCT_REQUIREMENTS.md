# 01. PRODUCT REQUIREMENTS DOCUMENT (PRD) — BN ONE v1.0

**Document ID:** BN-ONE-PRD-001  
**Version:** 1.0.0  
**Effective Date:** 2026-10-24  
**Author:** PO / Lead Architecture Team  
**Scope:** Global Financial Lifestyle Super App v1.0  
**Target Platforms:** Android, iOS, PC Web, Mobile Web, Admin Console  

---

## 1. Executive Summary & Brand Identity

### 1.1 Brand Mission
- **Brand Name:** BN ONE
- **Tagline:** *Your Money. Your World.*
- **Mission:** Empower global citizens to track and manage their financial lifestyle seamlessly across domestic daily expenditures and international travels within a unified, high-integrity platform.

### 1.2 Core Thesis & Scope Boundary
BN ONE v1.0 is a **manual financial lifestyle and personal wealth tracking software**. It strictly operates as an immutable ledger and planner based on user-entered transactions, manual account balances, and verified third-party reference exchange rate feeds.

> **CRITICAL LEGAL INVARIANT:**  
> BN ONE v1.0 does **NOT** connect to live banking APIs, execute wire transfers, issue debit/credit cards, collect user banking credentials, or take custody of customer funds. All user balances represent user-recorded ledgers.

---

## 2. In-Scope Deliverables (v1.0 In-Scope)

1. **Identity & Authentication:**
   - Email/password registration and cryptographically secure login (bcrypt salted hashes, JWT Bearer tokens).
   - Secure onboarding wizard capturing primary base currency (e.g., KRW, USD, EUR, JPY) and UI locale.
2. **Multi-Account & Ledger Management:**
   - User-defined accounts: Cash, Bank Accounts (Manual Ledger), Credit Cards, Loans, Investments (Manual), and Custom Assets.
   - Classification of Accounts into **Assets** and **Liabilities**.
3. **Transaction Engine (Double-Entry Minor Units):**
   - Income, Expense, and Intra-Account Transfer transactions.
   - Strict `BIGINT` minor-unit representation with zero floating-point roundoff issues.
   - Atomic balance updates: transfers debit the source account and credit the destination account with zero net-worth alteration.
4. **Budgeting & Category Management:**
   - Hierarchical transaction categories (Food & Dining, Housing, Transport, Shopping, Travel, Health, etc.).
   - Monthly category spending budgets with live progress tracking and threshold notifications.
5. **Subscription Tracking & Renewal Management:**
   - Tracking recurring bills, SaaS, streaming services, and memberships.
   - Next billing date calculation, monthly/annual cycle normalization, and proactive renewal alerts.
6. **Reference Foreign Exchange (Forex) System:**
   - Live tracking of 30+ major and regional fiat currencies against KRW and USD.
   - Multi-currency portfolio valuation converted to user's primary base currency.
   - Reference FX converter with clear disclaimers that rates are informational.
7. **Travel Expense & Multi-Currency Trip Wallets:**
   - Dedicated trip workspaces with destination country, dates, and trip budget.
   - Multi-currency expense logging tagged to trips.
   - Real-time travel budget burn-rate and post-trip settlement breakdown.
8. **Privacy, Compliance & Data Governance:**
   - Self-service data export: User transactions and accounts exported as structured CSV/JSON.
   - Self-service account deletion: 30-day grace period soft deletion or immediate hard purge.
9. **Operations Admin Console:**
   - Comprehensive back-office console with 8-level Role-Based Access Control (RBAC).
   - Real-time audit logs recording all staff operations with immutable actor signatures.
   - Feature flags, reference FX feed monitoring, user management, and support ticket triage.

---

## 3. Out-of-Scope (v1.0 Strictly Prohibited)

- Direct Open Banking scraping or bank credential collection.
- Real-time account balance auto-sync without user entry.
- Domestic wire transfers or automated clearinghouse (ACH/KFTC) execution.
- Cross-border remittances or commercial foreign exchange settlement.
- Credit/debit card issuance or physical card manufacturing.
- Custody or escrow of client money or electronic stored-value instruments.
- Execution of stock trades, cryptocurrency swaps, or loan origination.
- "Mock payment" buttons that deceive users into thinking a real payment was processed.

---

## 4. User Personas

| Persona | Description | Primary Needs |
|---|---|---|
| **Global Nomad / Frequent Traveler** | Works remotely or travels frequently across Japan, US, Europe, and Southeast Asia. | Tracking multi-currency expenses without manual exchange rate math; dedicated trip budgets. |
| **Disciplined Budgeter** | Strives to save 30% of income; manages multiple accounts (checking, savings, credit cards). | Category-based budgeting, net worth tracking, visibility into recurring subscriptions. |
| **Privacy-Conscious Planner** | Distrusts third-party screen scraping and open banking permissions; prefers manual control. | Complete data isolation, zero bank password sharing, clear audit trail, exportable CSV. |
| **Back-Office Auditor / Operator** | Operations and compliance personnel ensuring platform integrity and user support. | RBAC permissioning, tamper-proof audit trails, rapid customer support ticket handling. |

---

## 5. Non-Functional Requirements Summary

- **Security:** Zero plaintext secrets, bcrypt cost factor 12, short-lived JWT access tokens (15 min), secure refresh tokens, complete tenant isolation.
- **Performance:** Sub-200ms API response time at p95 for all read/write endpoints; 60fps Flutter rendering on mid-range devices; Web first contentful paint (FCP) < 1.2s.
- **Accessibility:** WCAG 2.1 Level AA compliance, minimum 48x48dp touch targets, semantic headings, full screen-reader labelling.
- **Reliability:** 99.9% uptime target; automated hourly FX refresh background workers; database transactional ACID guarantees.
