# 06. OPERATIONS ADMIN CONSOLE SPECIFICATION & RBAC

**Document ID:** BN-ONE-ADM-006  
**Version:** 1.0.0  
**Effective Date:** 2026-10-24  

---

## 1. Executive Mission & Security Baseline

The BN ONE Operations Admin Console (`apps/admin`) empowers authorized platform operators, support agents, compliance officers, and system administrators to oversee platform health, manage reference data, and support users without compromising user tenant privacy.

> **CRITICAL SECURITY INVARIANT:**  
> Customer financial account balances and transaction contents are **private user records**. Admin personnel cannot view, modify, or inject personal transactions into user ledgers unless explicitly requested via an authenticated support session with user consent. All admin actions are immutably logged into `admin_audit_logs`.

---

## 2. 8-Tier Role-Based Access Control (RBAC)

| Role Code | Role Name | Permitted Capabilities | Prohibited Capabilities |
|---|---|---|---|
| `PLATFORM_SUPER_ADMIN` | Super Administrator | Full platform control, RBAC assignment, system parameters, emergency controls | None |
| `OPERATIONS_ADMIN` | Operations Admin | User status management, support ticket resolution, reference data maintenance | RBAC permission edits |
| `SUPPORT_AGENT` | Customer Support | View and respond to support tickets, view basic user profile (email, status) | Feature flag alterations, FX editing |
| `CONTENT_MANAGER` | Content Manager | System announcements, localized FAQ articles, category taxonomy | User account modifications |
| `DATA_OPERATIONS` | Data Operations | Manage reference FX feeds, sync jobs, category classifications | Direct user credential edits |
| `SECURITY_ADMIN` | Security Officer | Audit log review, IP blacklists, session invalidation, privacy requests | General content updates |
| `ANALYST` | Business Analyst | Aggregated platform metrics (MAU, transaction volume counts) | PII inspection, user data modifications |
| `READ_ONLY_AUDITOR` | Compliance Auditor | Read-only inspection of audit logs, compliance reports, privacy queues | Any state modification |

---

## 3. Operations Console Architecture

```text
apps/admin/src/
  api/
    admin_client.ts        <-- Admin API client with Bearer authorization
  components/
    AdminLayout.vue        <-- Sidebar, operator profile, role badge, quick logout
  stores/
    adminAuth.ts           <-- Admin session and permission state
  views/
    AdminLoginView.vue     <-- Operator credential authentication
    AdminDashboardView.vue <-- Platform metrics (Total users, active tickets, FX health)
    AdminUsersView.vue     <-- User directory, status management (ACTIVE, SUSPENDED)
    AdminSupportView.vue   <-- Support ticket triage, reply drafting, status closure
    AdminFXView.vue        <-- Reference FX provider status, manual sync triggers
    AdminPrivacyView.vue   <-- GDPR / Privacy deletion and data export requests
    AdminAuditLogsView.vue <-- Tamper-proof filterable audit log inspection
    AdminFeatureFlagsView.vue <-- Dynamic runtime toggles (e.g., travel_wallets_enabled)
```

---

## 4. Immutable Audit Trail Integration
Every mutating API request performed by an operator executes inside a database transaction that concurrently records:
- `admin_id`: UUID of the authenticated operator.
- `action`: Specific operation code (e.g., `USER_SUSPENDED`, `FX_SYNC_TRIGGERED`, `TICKET_REPLIED`).
- `target_type` & `target_id`: Entity being impacted.
- `ip_address` & `user_agent`: Client network metadata.
- `changes`: JSON before/after state diff.
