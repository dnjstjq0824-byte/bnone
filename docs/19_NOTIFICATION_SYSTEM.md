# 19. NOTIFICATION & ALERT SYSTEM SPECIFICATION

**Document ID:** BN-ONE-NTF-019  
**Version:** 1.0.0  
**Effective Date:** 2026-10-24  

---

## 1. Notification Categories & Triggers

BN ONE manages communications to keep users aware of financial events without spamming.

| Notification Category | Trigger Condition | Delivery Channel | Priority |
|---|---|---|---|
| **Subscription Renewal Reminder** | 3 days and 1 day prior to `next_renewal_date` for active subscriptions. | In-App Feed + Mobile Push | High |
| **Budget Threshold Warning** | Total category expenditure reaches 80% or 100% of defined monthly limit. | In-App Feed + Mobile Push | Medium |
| **Trip Expense Summary** | End date of trip workspace reached; prompt to review final settlement. | In-App Feed | Low |
| **Security & Privacy Alerts** | New device login, password change, data export ready, account deletion scheduled. | In-App Feed + Transactional Email | Critical |
| **System Announcements** | Service maintenance notices, terms of service updates. | In-App Feed Banner | Informational |

---

## 2. Notification Pipeline Architecture

```text
[Cron Worker / Event Trigger]
              │
              ▼
[Notification Service (services/api/internal/notifications)]
              │
              ├──► PostgreSQL (INSERT INTO notifications)
              │
              └──► Outbox Queue (Push / Email Delivery)
                        │
                        ▼
           [Firebase Cloud Messaging (FCM) / APNs]
                        │
                        ▼
                 User Device
```

---

## 3. Data Model & User Preferences

- **Schema:** Stored in `notifications(id, user_id, title, message, type, is_read, metadata, created_at)`.
- **Read State:** Real-time badge counter decrements when user opens the notification drawer or clicks "Mark All as Read" (`PUT /api/v1/notifications/read`).
- **User Opt-Outs:** Users can granularly disable marketing or budget alerts while keeping security-critical alerts mandatory.
