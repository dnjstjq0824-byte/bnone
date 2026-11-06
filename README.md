# MOA v1.0 — Global Financial Lifestyle Super App

[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen)]()
[![Flutter](https://img.shields.io/badge/Flutter-3.47+-blue)]()
[![Vue 3](https://img.shields.io/badge/Vue-3.x-emerald)]()
[![Go](https://img.shields.io/badge/Go-1.27+-00ADD8)]()
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791)]()
[![Deployment](https://img.shields.io/badge/Live%20Demo-Static%20Web%20Workers-orange)]()

> **Brand:** MOA  
> **Brand Message:** *Gather all your money.*  
> **Document ID:** `MOA-MASTER-PRODUCT-001` (PO Master Directive v1.0)  
> **Deployment Status:** **정적 웹 + 브라우저 내 로컬 목업 데모 (Static Web + Local Mock Demo)**  
> **Canonical Live URL:** **https://moa.dnjstjqppk.workers.dev**

---

## 0. Cloudflare Workers 배포 및 라이브 도메인 안내

> [!IMPORTANT]
> **라이브 환경 배포 범위 안내:**
> - 현재 Cloudflare Workers(`https://moa.dnjstjqppk.workers.dev`)에 배포된 아티팩트는 **`apps/web` 정적 빌드 SPA(Single Page Application)**입니다.
> - Go 백엔드 API(`services/api`), PostgreSQL 데이터베이스, 관리자 콘솔(`apps/admin`), 모바일 앱(`apps/mobile`)은 본 Cloudflare Workers 단일 주소에 배포되어 있지 않습니다.
> - 웹 화면의 데이터는 브라우저 로컬 스토리지(`localStorage` 키 `moa_mock_*`)에 저장·반영되는 **로컬 목업 데모**입니다.
> 
> **도메인 구성 안내:**
> - **정본 라이브 도메인:** `https://moa.dnjstjqppk.workers.dev` (Cloudflare 계정 `dnjstjqppk@gmail.com`의 정본 워커)
> - **서브도메인 참고:** `dnjstjq0824.workers.dev` 또한 Cloudflare 상에 존재하며 동일한 빌드가 호스팅될 수 있으나, 공식 기준 정본 도메인은 `dnjstjqppk.workers.dev` 1개입니다.
> - **Cloudflare Builds 설정:**
>   - Root config: `wrangler.jsonc` (assets: `./apps/web/dist`, not_found_handling: `single-page-application`)
>   - Build command: `npm --prefix apps/web install && npm --prefix apps/web run build`
>   - Target branch: `main` (정본)
>   - **필수 빌드 환경변수:** `VITE_MOCK_ENABLED=true` (이 워커는 백엔드 없이 동작하는 공개 데모이므로 데모 모드를 명시해야 하며, 미설정 시 앱이 백엔드 미연결 오류를 표시합니다)

---

## 1. Overview & Scope Boundary

**MOA** is a high-integrity, personal financial lifestyle super app designed to track daily domestic expenditures, multi-account assets and liabilities, recurring subscriptions, monthly category budgets, and multi-currency travel expenses in a unified platform.

> [!IMPORTANT]
> **Regulatory & Architecture Boundary:**
> MOA v1.0 operates on a **manual financial ledger & reference data** model. The software does NOT connect to live open banking APIs, execute wire transfers, issue credit/debit cards, store user banking credentials, or take custody of customer funds. All user balances represent user-entered ledgers.

---

## 2. Monorepo Architecture

```text
moa/
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
- **Web Frontend:** Compiled with 0 errors (`npm --prefix apps/web run build`)
- **Admin Console:** Compiled with 0 errors (`npm --prefix apps/admin run build`)
- **Contract & Invariant Tests:** 100% Passed (`npm test`)
- **Financial Invariant Tests:** 100% Passed (`go test -v ./internal/transactions/...`)
- **Flutter Mobile Analysis & Tests:** 100% Passed (`flutter analyze`, `flutter test`)
- **Compliance & Privacy:** 100% Verified against Korean FSC & Global Standards

**Status: Monorepo Codebase Complete & Static Web Demo Deployed**  
Booking Network Co., Ltd. — 2026-11-01

