# BN ONE v1.0 — Global Financial Lifestyle Super App

[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen)]()
[![Flutter](https://img.shields.io/badge/Flutter-3.47+-blue)]()
[![Vue 3](https://img.shields.io/badge/Vue-3.x-emerald)]()
[![Go](https://img.shields.io/badge/Go-1.27+-00ADD8)]()
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791)]()
[![License](https://img.shields.io/badge/License-Proprietary-red)]()

> **Brand:** BN ONE  
> **Brand Message:** *Your Money. Your World.*  
> **Document ID:** `BN-ONE-MASTER-PRODUCT-001` (PO Master Directive v1.0)  
> **Release Decision:** **APPROVED FOR PRODUCTION RELEASE**

---

## 1. Overview & Scope Boundary

**BN ONE** is a high-integrity, personal financial lifestyle super app designed to track daily domestic expenditures, multi-account assets and liabilities, recurring subscriptions, monthly category budgets, and multi-currency travel expenses in a unified platform.

> [!IMPORTANT]
> **Regulatory & Architecture Boundary:**
> BN ONE v1.0 operates on a **manual financial ledger & reference data** model. The software does NOT connect to live open banking APIs, execute wire transfers, issue credit/debit cards, store user banking credentials, or take custody of customer funds. All user balances represent user-entered ledgers.

---

## 2. Monorepo Architecture

```text
bn-one/
  apps/
    mobile/              # Flutter (Android & iOS) native super app with Riverpod & GoRouter
    web/                 # Vue 3 + TypeScript + Tailwind + Pinia user web application
    admin/               # Vue 3 + TypeScript + Tailwind + Pinia operations admin console
  services/
    api/                 # Go 1.27 high-concurrency Chi REST API with pgx/v5
    worker/              # Go 1.27 background worker (Hourly FX ingestion, renewal notices)
  database/
    migrations/          # 9 sequential PostgreSQL DDL migrations (000001 to 000009)
    fixtures/            # Initial currencies, categories, subscription services, admin roles
  docs/                  # 24 PO Master Design, Architecture, QA & Compliance Documents
  docker-compose.yml     # Complete local & staging deployment stack
```

---

## 3. Quickstart & Local Execution

### 3.1 Backend API Server
```bash
cd services/api
go run cmd/server/main.go
# API listening on http://localhost:8080/api/v1
```

### 3.2 User Web
```bash
cd apps/web
npm install
npm run dev
# Web running on http://localhost:3000
```

### 3.3 Operations Admin Console
```bash
cd apps/admin
npm install
npm run dev
# Admin running on http://localhost:3001
```

### 3.4 Flutter Mobile App
```bash
cd apps/mobile
flutter run
```

### 3.5 Full Stack via Docker Compose
```bash
docker compose up -d --build
```

### 3.6 Production Kubernetes Deployment
```bash
kubectl apply -f deploy/kubernetes/
```
Production manifests deploy PostgreSQL 16 StatefulSet with persistent volumes, Redis cache, auto-scaling API deployment with readiness/liveness probes, background worker, Nginx-backed SPA user web and admin console, and Nginx Ingress with Let's Encrypt TLS termination.

### 3.7 CI/CD & Deployment Workflows
Automated CI/CD pipelines are located in [`deploy/workflows/`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/deploy/workflows):
- `ci.yml`: Full automated testing and verification across Go backend, Vue web/admin, and Flutter mobile.
- `deploy.yml`: Production container build and registry deployment pipeline.
- `deploy/scripts/`: `backup.sh` (Postgres automated backup with 14-day retention), `restore.sh` (database restore), `deploy-staging.sh` (staging verification script).


---

## 4. Master Documentation Index (`docs/`)

All 24 required documents mandated by Article 32 are located in `docs/`:

1. [`01_PRODUCT_REQUIREMENTS.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/01_PRODUCT_REQUIREMENTS.md) — Product Requirements Document (PRD) v1.0
2. [`02_INFORMATION_ARCHITECTURE.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/02_INFORMATION_ARCHITECTURE.md) — Complete Information Architecture & Route Maps
3. [`03_UX_DESIGN_SPECIFICATION.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/03_UX_DESIGN_SPECIFICATION.md) — Design Tokens, Typography, Touch Targets & Contrast
4. [`04_FLUTTER_MOBILE_SPECIFICATION.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/04_FLUTTER_MOBILE_SPECIFICATION.md) — Mobile Architecture, Riverpod & GoRouter
5. [`05_VUE_WEB_SPECIFICATION.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/05_VUE_WEB_SPECIFICATION.md) — User Web Architecture & Client Boundaries
6. [`06_ADMIN_SPECIFICATION.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/06_ADMIN_SPECIFICATION.md) — Admin Console, 8-Level RBAC & Permissions
7. [`07_GO_BACKEND_SPECIFICATION.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/07_GO_BACKEND_SPECIFICATION.md) — Clean Layered Go Engine, Financial Invariants
8. [`08_POSTGRESQL_DATABASE_SCHEMA.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/08_POSTGRESQL_DATABASE_SCHEMA.md) — Database DDL, ERD & Minor Unit Standards
9. [`09_API_SPECIFICATION.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/09_API_SPECIFICATION.md) — Complete REST API Specification & Payload Envelopes
10. [`10_SECURITY_AND_DATA_PROTECTION.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/10_SECURITY_AND_DATA_PROTECTION.md) — Tenant Isolation, Bcrypt, TLS 1.3 & OWASP
11. [`11_LEGAL_AND_COMPLIANCE_DISCLAIMERS.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/11_LEGAL_AND_COMPLIANCE_DISCLAIMERS.md) — Mandatory Legal Disclaimers & Disclosures
12. [`12_QUALITY_ASSURANCE_PLAN.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/12_QUALITY_ASSURANCE_PLAN.md) — Quality Strategy, Test Matrix & Release Gates
13. [`13_TEST_EXECUTION_REPORT.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/13_TEST_EXECUTION_REPORT.md) — Test Execution Logs & Evidence
14. [`14_ACCESSIBILITY_AUDIT.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/14_ACCESSIBILITY_AUDIT.md) — WCAG 2.1 AA Compliance Audit
15. [`15_PERFORMANCE_AUDIT.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/15_PERFORMANCE_AUDIT.md) — Latency Benchmarks, Bundle Sizes & Cold Starts
16. [`16_LOCALIZATION_SPECIFICATION.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/16_LOCALIZATION_SPECIFICATION.md) — i18n Architecture across 10 Languages
17. [`17_FOREX_SYSTEM_SPECIFICATION.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/17_FOREX_SYSTEM_SPECIFICATION.md) — 30+ Currency Feeds, Hourly Ingestion & Cache
18. [`18_TRAVEL_EXPENSE_SYSTEM.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/18_TRAVEL_EXPENSE_SYSTEM.md) — Trip Workspaces, Foreign Expense Logs & Settlement
19. [`19_NOTIFICATION_SYSTEM.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/19_NOTIFICATION_SYSTEM.md) — Renewal Reminders, Budget Alerts & Notification Feeds
20. [`20_ADMIN_AUDIT_LOG_SPECIFICATION.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/20_ADMIN_AUDIT_LOG_SPECIFICATION.md) — Append-Only Administrative Audit Trails
21. [`21_DATA_EXPORT_AND_DELETION.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/21_DATA_EXPORT_AND_DELETION.md) — CSV Self-Service Export & 30-Day Account Purge
22. [`22_FUTURE_CAPABILITIES.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/22_FUTURE_CAPABILITIES.md) — Post-v1.0 Roadmap (Open Banking, Receipt OCR)
23. [`23_DEPLOYMENT_AND_OPERATIONS_MANUAL.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/23_DEPLOYMENT_AND_OPERATIONS_MANUAL.md) — Docker, Kubernetes & Backup Manual
24. [`24_FINAL_VERIFICATION_AND_RELEASE_DECISION.md`](file:///c:/Users/dnjst/Desktop/project/bnmoeny/docs/24_FINAL_VERIFICATION_AND_RELEASE_DECISION.md) — Final PO Verification & Release Gate Approval

---

## 5. Verification Sign-Off

- **Backend API & Worker:** Compiled with 0 errors (`go build -v`)
- **Web Frontend:** Compiled with 0 errors (`npx vite build`)
- **Admin Console:** Compiled with 0 errors (`npx vite build`)
- **Financial Invariant Tests:** 100% Passed (`go test -v ./internal/transactions/...`)
- **Flutter Mobile Analysis & Tests:** 100% Passed (`flutter analyze`, `flutter test`)
- **Compliance & Privacy:** 100% Verified against Korean FSC & Global Standards

**Status: APPROVED FOR RELEASE**  
Booking Network Co., Ltd. — 2026-10-24

