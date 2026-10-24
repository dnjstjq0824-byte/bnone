# 20. ADMINISTRATIVE AUDIT LOG SPECIFICATION

**Document ID:** BN-ONE-AUD-020  
**Version:** 1.0.0  
**Effective Date:** 2026-10-24  

---

## 1. Compliance Mandate & Immutability Guarantee

To fulfill ISO/IEC 27001, SOC 2 Type II, and South Korean Personal Information Protection Act (PIPA) Article 29 standards, every operational interaction with user metadata, system parameters, or security configurations is immutably recorded.

> **IMMUTABILITY INVARIANT:**  
> Audit log records are append-only. The database schema grants no `UPDATE` or `DELETE` privileges on the `admin_audit_logs` table to standard database roles. No administrative endpoint exists to alter or remove audit records.

---

## 2. Audit Event Taxonomy

| Event Code | Category | Trigger / Action Description |
|---|---|---|
| `ADMIN_LOGIN_SUCCESS` | Authentication | Operator successfully authenticated with password + 2FA |
| `ADMIN_LOGIN_FAILURE` | Authentication | Failed login attempt by operator (tracks IP and username) |
| `USER_STATUS_UPDATED` | User Ops | User account marked as `SUSPENDED`, `ACTIVE`, or `LOCKED` |
| `SUPPORT_TICKET_REPLIED` | Support | Support agent drafted and dispatched a response to user inquiry |
| `SUPPORT_TICKET_CLOSED` | Support | Support ticket marked as resolved |
| `FX_PROVIDER_SYNC` | Data Ops | Manual trigger of foreign exchange rate refresh job |
| `FEATURE_FLAG_TOGGLED` | System Ops | Dynamic feature flag state changed (e.g., enabling trip wallets) |
| `PRIVACY_EXPORT_ACCESSED` | Compliance | Operator verified user data export queue |
| `PRIVACY_PURGE_CONFIRMED` | Compliance | Operator confirmed user hard-deletion request after grace period |

---

## 3. Audit Log Schema & Payload Format

```sql
CREATE TABLE admin_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id UUID NOT NULL REFERENCES admin_users(id),
    action VARCHAR(64) NOT NULL,
    target_type VARCHAR(64) NOT NULL,
    target_id VARCHAR(128) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT NOT NULL,
    changes JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
```

### Sample Audit Payload (`changes` JSONB):
```json
{
  "before": {
    "status": "ACTIVE",
    "notes": "Standard account"
  },
  "after": {
    "status": "SUSPENDED",
    "notes": "Suspended pending fraud investigation #9281"
  }
}
```
