# 12. QUALITY ASSURANCE PLAN & RELEASE GATES

**Document ID:** BN-ONE-QA-012  
**Version:** 1.0.0  
**Effective Date:** 2026-10-24  

---

## 1. Quality Engineering Strategy

The BN ONE QA plan is structured across three verification horizons:
1. **Unit & Financial Invariant Testing:** Testing calculation correctness, minor-unit conversions, double-entry arithmetic, and tenant isolation at the service level.
2. **Integration & API Verification:** Testing database transaction rollbacks, foreign exchange synchronization workers, and auth session flows.
3. **End-to-End & UI/UX Validation:** Verifying that Flutter mobile apps and Vue web frontends properly render layouts, honor accessibility requirements, and handle offline transitions.

---

## 2. Release Gates Matrix

| Release Gate | Verification Criteria | Status |
|---|---|---|
| **Gate 1: Compilation & Static Analysis** | Backend (`go build`), Worker (`go build`), Web (`npx vite build`), Admin (`npx vite build`), Flutter (`flutter analyze --no-fatal-infos`) must compile with 0 fatal errors. | **PASSED** |
| **Gate 2: Financial Integrity Tests** | Atomic transfers must preserve net worth; minor units must never produce rounding errors; tenant queries must strictly reject cross-user access. | **PASSED** |
| **Gate 3: UI Smoke Tests** | Flutter widget tests initialize cleanly; Vue stores mount without runtime errors; responsive layout scales to 320px mobile width. | **PASSED** |
| **Gate 4: Security & Privacy Check** | Bcrypt cost factor >= 12; zero plain-text passwords in DB; audit logging fires on admin operations; disclaimers visible on all 4 core screens. | **PASSED** |
| **Gate 5: PO Master Review** | 100% feature coverage of PRD v1.0 specifications; all 24 master documents documented. | **PASSED** |

