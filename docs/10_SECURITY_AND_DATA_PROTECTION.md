# 10. SECURITY AND DATA PROTECTION SPECIFICATION

**Document ID:** BN-ONE-SEC-010  
**Version:** 1.0.0  
**Effective Date:** 2026-10-24  

---

## 1. Threat Modeling & Core Defense Principles

BN ONE holds sensitive user financial records. The platform enforces defense-in-depth across storage, transmission, and operational layers.

### 1.1 Tenant Isolation Invariant
- **Logical Data Segregation:** All user data (accounts, entries, budgets, trips, notes) is logically segregated with mandatory `user_id` ownership constraints.
- **Direct Object Reference (IDOR) Protection:** Backend services never fetch records purely by entity UUID without validating `WHERE user_id = $auth_user_id`. An attacker cannot access another user's accounts by guessing UUIDs.

---

## 2. Cryptographic Controls

- **Password Storage:** Passwords hashed with **bcrypt** (cost factor 12) using cryptographically secure random per-user salt. Plaintext passwords never hit logs or persistent storage.
- **Transport Security:** All endpoints require **TLS 1.3** (or TLS 1.2 minimum) with strict Forward Secrecy cipher suites. HSTS headers (`Strict-Transport-Security: max-age=63072000; includeSubDomains; preload`) enforced.
- **JWT Signatures:** JSON Web Tokens signed using HMAC-SHA256 (`HS256`) with a 256-bit cryptographically random secret. Token validity is 15 minutes; refresh tokens are stored hashed in database sessions.

---

## 3. OWASP Top 10 Mitigations

| Vulnerability | Mitigation Strategy |
|---|---|
| **A01: Broken Access Control** | Enforced at Go service layer before repository calls; admin RBAC matrix with route-level and action-level checks. |
| **A02: Cryptographic Failures** | Bcrypt hashing, TLS 1.3, encrypted local mobile storage (`flutter_secure_storage`). |
| **A03: Injection (SQL / XSS)** | Parameterized queries exclusively via `pgx/v5`; Vue 3 auto-escaping for HTML templates; Go JSON decoding without raw evaluation. |
| **A04: Insecure Design** | Explicit prohibition of bank account credentials or third-party scraping; manual ledger design eliminates unauthorized fund movement risk. |
| **A05: Security Misconfiguration** | Strict CORS whitelist (`AllowOrigins`), rate limiting, no debug stack traces exposed in production envelopes. |
| **A07: Identification and Authentication** | Account lockout after repeated failed attempts; bcrypt salt; JWT expiration handling. |
| **A09: Security Logging & Monitoring** | Immutable append-only audit trail for all admin actions with IP and timestamp; worker health monitors. |

---

## 4. Privacy Governance (GDPR & PIPA Compliance)

- **Right of Access & Portability:** Self-service CSV export generates structured account and transaction archives on demand.
- **Right to Erasure (Right to be Forgotten):** Users can trigger account deletion with a 30-day grace period (soft delete) or immediate hard purge cascading through all accounts and transactions.
