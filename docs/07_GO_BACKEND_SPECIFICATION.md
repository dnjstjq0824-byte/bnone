# 07. GO BACKEND ARCHITECTURE & ENGINE SPECIFICATION

**Document ID:** BN-ONE-BE-007  
**Version:** 1.0.0  
**Effective Date:** 2026-10-24  

---

## 1. Clean Layered Architecture Overview

The BN ONE backend server (`services/api`) is engineered in Go 1.27 with strict separation of concerns, zero circular dependencies, and high-concurrency throughput.

```text
HTTP Request
     │
     ▼
Chi Router + Global Middleware (Logger, Recoverer, CORS, RequestID, RateLimiter)
     │
     ▼
Auth Middleware (JWT Signature Verification, Claims Extraction)
     │
     ▼
HTTP Handler (JSON Decoding, Semantic Parameter Validation)
     │
     ▼
Application Service (Business Logic, Orchestration, Validation)
     │
     ▼
Domain Engine (Double-Entry Balance Updates, Net Worth Calculation)
     │
     ▼
Repository Layer (PostgreSQL pgx/v5 Driver / SQLite Pure Go in test)
     │
     ▼
PostgreSQL 16 Database (ACID Transaction Boundaries)
```

---

## 2. Directory Structure

```text
services/api/
  cmd/
    server/
      main.go              <-- Server bootstrap, dependency injection, HTTP listener
  internal/
    auth/
      service.go           <-- User registration, bcrypt verification, JWT generation
    accounts/
      service.go           <-- Manual account management (Assets, Liabilities)
    transactions/
      service.go           <-- Transaction engine, double-entry atomic updates
      service_test.go      <-- Rigorous financial integrity & tenant isolation tests
    budgets/
      service.go           <-- Monthly budget progress & spending calculations
    subscriptions/
      service.go           <-- Recurring subscriptions & next renewal projections
    forex/
      service.go           <-- Multi-currency rate cache & informational conversion
    trips/
      service.go           <-- Multi-currency travel wallets & expense settlement
    notifications/
      service.go           <-- In-app notifications & unread count tracking
    admin/
      service.go           <-- RBAC authorization, user management, audit logging
```

---

## 3. Financial Integrity & Invariant Rules

1. **Integer Minor Units:** All monetary values are strictly stored as `int64` (PostgreSQL `BIGINT`). For example:
   - ₩15,000 KRW is represented as `15000` (0 decimal places).
   - $150.25 USD is represented as `15025` (2 decimal places).
   - ¥990 JPY is represented as `990` (0 decimal places).
2. **Transfer Invariant:** An intra-account transfer from Account A to Account B must create two atomic ledger entries:
   - Debit Account A: `-amount`
   - Credit Account B: `+amount`
   - **Net Worth Impact:** Exactly zero (`Δ Net Worth = 0`).
3. **Strict Tenant Isolation:** Every repository query MUST enforce `WHERE user_id = $1`. Attempts to query or mutate an account or transaction by ID without matching `user_id` return `404 Not Found` or `403 Forbidden`.

