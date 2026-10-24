# 09. REST API SPECIFICATION

**Document ID:** BN-ONE-API-009  
**Version:** 1.0.0  
**Effective Date:** 2026-10-24  

---

## 1. Global API Conventions

- **Base URL:** `https://api.bnone.bookingnetwork.com/api/v1`
- **Protocol:** HTTPS / TLS 1.3
- **Content-Type:** `application/json; charset=utf-8`
- **Authentication:** `Authorization: Bearer <JWT_ACCESS_TOKEN>`
- **Standard Error Envelope:**
  ```json
  {
    "error": {
      "code": "VALIDATION_FAILED",
      "message": "Amount must be a positive integer in minor units.",
      "details": null
    }
  }
  ```

---

## 2. API Endpoint Matrix

### 2.1 Authentication & Profile (`/auth`, `/users`)
- `POST /api/v1/auth/register` — Register new user account with base currency & locale.
- `POST /api/v1/auth/login` — Authenticate credentials; return JWT access & refresh tokens.
- `GET /api/v1/users/me` — Retrieve current authenticated user profile.
- `PUT /api/v1/users/me` — Update base currency, locale, and notification preferences.
- `GET /api/v1/users/export` — Trigger full user data export (CSV format).
- `DELETE /api/v1/users/me` — Schedule account deletion with 30-day grace period or immediate hard purge.

### 2.2 Accounts & Net Worth (`/accounts`)
- `GET /api/v1/accounts` — List all user financial accounts (Assets & Liabilities).
- `POST /api/v1/accounts` — Create manual ledger account (Cash, Bank, Card, Loan, Asset).
- `GET /api/v1/accounts/net-worth` — Return aggregated Net Worth, Total Assets, and Total Liabilities in base currency.
- `DELETE /api/v1/accounts/{id}` — Archive or delete an account.

### 2.3 Transactions Engine (`/transactions`)
- `GET /api/v1/transactions` — Query paginated transactions with category/date filters.
- `POST /api/v1/transactions` — Record new transaction (`INCOME`, `EXPENSE`, `TRANSFER`).
- `GET /api/v1/transactions/{id}` — Retrieve transaction details and split entries.
- `DELETE /api/v1/transactions/{id}` — Delete transaction and atomically revert account balances.

### 2.4 Budgets & Analytics (`/budgets`, `/analytics`)
- `GET /api/v1/budgets` — List monthly category budgets with actual spending vs limit.
- `POST /api/v1/budgets` — Set or update monthly budget for a category.
- `GET /api/v1/analytics/spending` — Retrieve category-wise expense distribution and cashflow totals.

### 2.5 Subscriptions (`/subscriptions`)
- `GET /api/v1/subscriptions` — List active recurring subscriptions and billing cycles.
- `POST /api/v1/subscriptions` — Register recurring subscription (service, amount, frequency, next date).
- `DELETE /api/v1/subscriptions/{id}` — Cancel or delete subscription tracker.

### 2.6 Reference Forex (`/forex`)
- `GET /api/v1/forex/rates` — Retrieve latest reference exchange rates.
- `POST /api/v1/forex/convert` — Perform reference currency conversion.

### 2.7 Travel Expense Management (`/trips`)
- `GET /api/v1/trips` — List travel workspaces.
- `POST /api/v1/trips` — Create travel workspace with destination, dates, and budget.
- `GET /api/v1/trips/{id}/expenses` — Retrieve trip-specific expenses and settlement breakdown.
- `POST /api/v1/trips/{id}/expenses` — Log trip expense in local foreign currency.

### 2.8 Operations Admin (`/admin`)
- `POST /api/v1/admin/login` — Operator authentication.
- `GET /api/v1/admin/dashboard` — Platform operational metrics.
- `GET /api/v1/admin/users` — User administration list.
- `GET /api/v1/admin/support/tickets` — Triage customer support tickets.
- `POST /api/v1/admin/support/tickets/{id}/reply` — Post response to support inquiry.
- `GET /api/v1/admin/audit-logs` — Query immutable administrative audit logs.
- `GET /api/v1/admin/feature-flags` — Query and toggle runtime feature flags.
