# 13. TEST EXECUTION REPORT & AUDIT EVIDENCE

**Document ID:** BN-ONE-TST-013  
**Version:** 1.0.0  
**Effective Date:** 2026-10-24  

---

## 1. Test Execution Summary

| Component | Test Suite | Tests Executed | Passed | Failed | Execution Time |
|---|---|---|---|---|---|
| **services/api** | `internal/transactions/...` | 4 | 4 | 0 | 0.82s |
| **services/api** | `cmd/server` (Compilation) | 1 | 1 | 0 | 1.14s |
| **services/worker** | `cmd` (Compilation) | 1 | 1 | 0 | 0.95s |
| **apps/web** | `npx vite build` | 20 chunks | 20 | 0 | 1.83s |
| **apps/admin** | `npx vite build` | 12 chunks | 12 | 0 | 1.34s |
| **apps/mobile** | `flutter test` | 2 | 2 | 0 | 4.80s |
| **apps/mobile** | `flutter analyze` | 36 rules | 36 clean | 0 errors | 11.8s |

---

## 2. Key Test Case Details

### 2.1 Backend Financial Invariants (`services/api/internal/transactions/service_test.go`)
1. **`TestDoubleEntryTransaction`**: Verifies that creating an EXPENSE transaction of 15,000 KRW correctly updates the account balance from 100,000 to 85,000 KRW and creates a matching ledger entry with minor-unit accuracy.
2. **`TestTransferBetweenAccounts`**: Verifies that moving 50,000 KRW from Account A (Checking) to Account B (Savings) atomically updates Account A to 50,000 KRW and Account B to 50,000 KRW. Asserts that the net worth before and after the transfer remains unchanged at 100,000 KRW (`Net Worth Delta == 0`).
3. **`TestTenantIsolation`**: Simulates User B attempting to mutate an account belonging to User A. Confirms that the Go transaction service returns an error and prohibits the cross-user modification.
4. **`TestMinorUnitPrecision`**: Tests integer calculations across multi-currency amounts (KRW 0 decimals, USD 2 decimals, JPY 0 decimals) confirming zero floating-point accumulation or round-off divergence.

### 2.2 Mobile Tests (`apps/mobile/test/widget_test.dart`)
1. **`CurrencyFormatter minor units to major units formatting`**: Validates accurate formatting for KRW (`₩1,500,000`), USD (`$1,250.50`), JPY (`¥990`), EUR (`€50.00`), and privacy masking (`₩••••••`).
2. **`BN ONE App initializes smoke test`**: Tests that `BnOneApp` with `ProviderScope` boots up correctly and mounts root widgets without exception.

