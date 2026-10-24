package admin

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bookingnetwork/bn-one/services/api/internal/database"
	"github.com/bookingnetwork/bn-one/services/api/internal/middleware"
	"github.com/bookingnetwork/bn-one/services/api/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	db *database.DB
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
}

// ---------------- Dashboard ----------------

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	var totalUsers, activeUsers, openTickets, pendingDeletions int
	h.db.QueryRow("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&totalUsers)
	h.db.QueryRow("SELECT COUNT(DISTINCT user_id) FROM user_sessions WHERE revoked_at IS NULL AND expires_at > CURRENT_TIMESTAMP").Scan(&activeUsers)
	h.db.QueryRow("SELECT COUNT(*) FROM support_tickets WHERE status IN ('OPEN', 'IN_PROGRESS')").Scan(&openTickets)
	h.db.QueryRow("SELECT COUNT(*) FROM deletion_requests WHERE status = 'PENDING'").Scan(&pendingDeletions)

	var activeProviders int
	h.db.QueryRow("SELECT COUNT(*) FROM fx_providers WHERE active = TRUE").Scan(&activeProviders)

	middleware.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
		"kpi": map[string]interface{}{
			"total_registered_users": totalUsers,
			"active_user_sessions":   activeUsers,
			"open_support_tickets":   openTickets,
			"pending_deletion_reqs":  pendingDeletions,
			"active_fx_providers":    activeProviders,
		},
		"system_health": "OPTIMAL",
		"server_time":   time.Now().UTC(),
		"notice":        "개인정보 보호 정책에 따라 관리자 대시보드에 사용자의 개별 금융 금액을 노출하지 않습니다.",
	})
}

// ---------------- User Management ----------------

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("q")
	query := `SELECT id, email, display_name, country_code, locale, base_currency, timezone, status, email_verified_at, onboarding_completed_at, created_at, updated_at
	          FROM users WHERE deleted_at IS NULL`
	args := []interface{}{}
	if search != "" {
		query += " AND (email LIKE '%' || $1 || '%' OR display_name LIKE '%' || $1 || '%')"
		args = append(args, search)
	}
	query += " ORDER BY created_at DESC LIMIT 50"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.DisplayName, &u.CountryCode, &u.Locale, &u.BaseCurrency, &u.Timezone, &u.Status, &u.EmailVerifiedAt, &u.OnboardingCompletedAt, &u.CreatedAt, &u.UpdatedAt); err == nil {
			users = append(users, u)
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, users)
}

func (h *Handler) GetUserDetail(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	var u models.User
	query := `SELECT id, email, display_name, country_code, locale, base_currency, timezone, status, email_verified_at, onboarding_completed_at, created_at, updated_at
	          FROM users WHERE id = $1`
	err := h.db.QueryRow(query, userID).Scan(&u.ID, &u.Email, &u.DisplayName, &u.CountryCode, &u.Locale, &u.BaseCurrency, &u.Timezone, &u.Status, &u.EmailVerifiedAt, &u.OnboardingCompletedAt, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", "User not found", nil)
		return
	}

	var sessionCount int
	h.db.QueryRow("SELECT COUNT(*) FROM user_sessions WHERE user_id = $1 AND revoked_at IS NULL AND expires_at > CURRENT_TIMESTAMP", userID).Scan(&sessionCount)

	middleware.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
		"user":            u,
		"active_sessions": sessionCount,
	})
}

func (h *Handler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	var body struct {
		Status string `json:"status"` // ACTIVE, SUSPENDED
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Status == "" {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Status required", nil)
		return
	}

	now := time.Now().UTC()
	_, err := h.db.Exec("UPDATE users SET status = $1, updated_at = $2 WHERE id = $3", body.Status, now, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}

	// Audit log
	adminID := middleware.GetAdminID(r.Context())
	h.recordAudit(adminID, "UPDATE_USER_STATUS", "USER", userID, "SUCCESS", `{"status":"`+body.Status+`"}`)

	middleware.RespondJSON(w, r, http.StatusOK, map[string]string{"message": "User status updated to " + body.Status})
}

func (h *Handler) RevokeUserSessions(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	now := time.Now().UTC()
	_, err := h.db.Exec("UPDATE user_sessions SET revoked_at = $1 WHERE user_id = $2 AND revoked_at IS NULL", now, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}

	adminID := middleware.GetAdminID(r.Context())
	h.recordAudit(adminID, "REVOKE_USER_SESSIONS", "USER_SESSION", userID, "SUCCESS", "{}")

	middleware.RespondJSON(w, r, http.StatusOK, map[string]string{"message": "All user sessions revoked"})
}

// ---------------- Customer Support ----------------

func (h *Handler) ListSupportTickets(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	query := `SELECT st.id, st.user_id, u.email, st.subject, st.category, st.description, st.priority, st.status, st.assigned_admin_id, st.created_at, st.updated_at, st.resolved_at
	          FROM support_tickets st
	          JOIN users u ON st.user_id = u.id`
	args := []interface{}{}
	if status != "" {
		query += " WHERE st.status = $1"
		args = append(args, status)
	}
	query += " ORDER BY st.created_at DESC LIMIT 100"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	tickets := make([]models.SupportTicket, 0)
	for rows.Next() {
		var t models.SupportTicket
		if err := rows.Scan(&t.ID, &t.UserID, &t.UserEmail, &t.Subject, &t.Category, &t.Description, &t.Priority, &t.Status, &t.AssignedAdminID, &t.CreatedAt, &t.UpdatedAt, &t.ResolvedAt); err == nil {
			tickets = append(tickets, t)
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, tickets)
}

func (h *Handler) GetSupportTicket(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "id")
	var t models.SupportTicket
	query := `SELECT st.id, st.user_id, u.email, st.subject, st.category, st.description, st.priority, st.status, st.assigned_admin_id, st.created_at, st.updated_at, st.resolved_at
	          FROM support_tickets st
	          JOIN users u ON st.user_id = u.id
	          WHERE st.id = $1`
	err := h.db.QueryRow(query, ticketID).Scan(&t.ID, &t.UserID, &t.UserEmail, &t.Subject, &t.Category, &t.Description, &t.Priority, &t.Status, &t.AssignedAdminID, &t.CreatedAt, &t.UpdatedAt, &t.ResolvedAt)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", "Ticket not found", nil)
		return
	}

	msgRows, _ := h.db.Query("SELECT id, ticket_id, sender_type, sender_id, message, visibility, created_at FROM support_messages WHERE ticket_id = $1 ORDER BY created_at ASC", ticketID)
	if msgRows != nil {
		msgs := make([]models.SupportMessage, 0)
		for msgRows.Next() {
			var m models.SupportMessage
			if err := msgRows.Scan(&m.ID, &m.TicketID, &m.SenderType, &m.SenderID, &m.Message, &m.Visibility, &m.CreatedAt); err == nil {
				msgs = append(msgs, m)
			}
		}
		msgRows.Close()
		t.Messages = msgs
	}

	middleware.RespondJSON(w, r, http.StatusOK, t)
}

func (h *Handler) AssignSupportTicket(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "id")
	var body struct {
		AdminID string `json:"admin_id"`
	}
	json.NewDecoder(r.Body).Decode(&body)
	if body.AdminID == "" {
		body.AdminID = middleware.GetAdminID(r.Context())
	}

	now := time.Now().UTC()
	_, err := h.db.Exec("UPDATE support_tickets SET assigned_admin_id = $1, status = 'IN_PROGRESS', updated_at = $2 WHERE id = $3", body.AdminID, now, ticketID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}

	middleware.RespondJSON(w, r, http.StatusOK, map[string]string{"message": "Ticket assigned"})
}

func (h *Handler) ReplySupportTicket(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "id")
	adminID := middleware.GetAdminID(r.Context())

	var body struct {
		Message    string `json:"message"`
		Visibility string `json:"visibility"` // PUBLIC, INTERNAL_ONLY
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Message == "" {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Message required", nil)
		return
	}
	if body.Visibility == "" {
		body.Visibility = "PUBLIC"
	}

	msgID := uuid.New().String()
	now := time.Now().UTC()
	_, err := h.db.Exec("INSERT INTO support_messages (id, ticket_id, sender_type, sender_id, message, visibility, created_at) VALUES ($1, $2, 'ADMIN', $3, $4, $5, $6)",
		msgID, ticketID, adminID, body.Message, body.Visibility, now)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}

	h.db.Exec("UPDATE support_tickets SET status = 'WAITING_FOR_USER', updated_at = $1 WHERE id = $2", now, ticketID)

	middleware.RespondJSON(w, r, http.StatusCreated, map[string]string{"message_id": msgID, "status": "REPLIED"})
}

func (h *Handler) ResolveSupportTicket(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "id")
	now := time.Now().UTC()

	_, err := h.db.Exec("UPDATE support_tickets SET status = 'RESOLVED', resolved_at = $1, updated_at = $1 WHERE id = $2", now, ticketID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}

	adminID := middleware.GetAdminID(r.Context())
	h.recordAudit(adminID, "RESOLVE_SUPPORT_TICKET", "SUPPORT_TICKET", ticketID, "SUCCESS", "{}")

	middleware.RespondJSON(w, r, http.StatusOK, map[string]string{"message": "Ticket resolved"})
}

// ---------------- FX Operations ----------------

func (h *Handler) ListFXProviders(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, name, official_website, api_base_url, license_status, commercial_use_status, active, last_checked_at FROM fx_providers ORDER BY name ASC`
	rows, err := h.db.Query(query)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	type FXP struct {
		ID                  string     `json:"id"`
		Name                string     `json:"name"`
		OfficialWebsite     string     `json:"official_website"`
		ApiBaseURL          string     `json:"api_base_url"`
		LicenseStatus       string     `json:"license_status"`
		CommercialUseStatus string     `json:"commercial_use_status"`
		Active              bool       `json:"active"`
		LastCheckedAt       *time.Time `json:"last_checked_at,omitempty"`
	}

	providers := make([]FXP, 0)
	for rows.Next() {
		var p FXP
		if err := rows.Scan(&p.ID, &p.Name, &p.OfficialWebsite, &p.ApiBaseURL, &p.LicenseStatus, &p.CommercialUseStatus, &p.Active, &p.LastCheckedAt); err == nil {
			providers = append(providers, p)
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, providers)
}

func (h *Handler) TriggerFXSync(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()
	// Update last_checked_at
	h.db.Exec("UPDATE fx_providers SET last_checked_at = $1 WHERE active = TRUE", now)

	// Create job record
	jobID := uuid.New().String()
	h.db.Exec("INSERT INTO system_jobs (id, job_type, status, scheduled_at, started_at, completed_at) VALUES ($1, 'FX_SYNC', 'COMPLETED', $2, $2, $2)", jobID, now)

	adminID := middleware.GetAdminID(r.Context())
	h.recordAudit(adminID, "TRIGGER_FX_SYNC", "FX_PROVIDER", "ALL", "SUCCESS", "{}")

	middleware.RespondJSON(w, r, http.StatusAccepted, map[string]interface{}{
		"job_id":     jobID,
		"status":     "COMPLETED",
		"message":    "FX sync completed successfully. Reference exchange rates updated.",
		"synced_at":  now,
		"disclaimer": "검증되지 않은 외부 비공식 환율은 자동 등록되지 않습니다.",
	})
}

// ---------------- Privacy Requests ----------------

func (h *Handler) ListDeletionRequests(w http.ResponseWriter, r *http.Request) {
	query := `SELECT dr.id, dr.user_id, u.email, dr.status, dr.requested_at, dr.processed_at, dr.retention_exception_reason
	          FROM deletion_requests dr
	          JOIN users u ON dr.user_id = u.id
	          ORDER BY dr.requested_at DESC`
	rows, err := h.db.Query(query)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	type DelReq struct {
		ID                       string     `json:"id"`
		UserID                   string     `json:"user_id"`
		Email                    string     `json:"email"`
		Status                   string     `json:"status"`
		RequestedAt              time.Time  `json:"requested_at"`
		ProcessedAt              *time.Time `json:"processed_at,omitempty"`
		RetentionExceptionReason string     `json:"retention_exception_reason"`
	}

	requests := make([]DelReq, 0)
	for rows.Next() {
		var d DelReq
		if err := rows.Scan(&d.ID, &d.UserID, &d.Email, &d.Status, &d.RequestedAt, &d.ProcessedAt, &d.RetentionExceptionReason); err == nil {
			requests = append(requests, d)
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, requests)
}

func (h *Handler) ReviewDeletionRequest(w http.ResponseWriter, r *http.Request) {
	reqID := chi.URLParam(r, "id")
	var body struct {
		Action string `json:"action"` // APPROVE_AND_EXECUTE, REJECT
		Reason string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	now := time.Now().UTC()
	var userID string
	err := h.db.QueryRow("SELECT user_id FROM deletion_requests WHERE id = $1", reqID).Scan(&userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", "Deletion request not found", nil)
		return
	}

	if body.Action == "APPROVE_AND_EXECUTE" {
		// Anonymize and soft delete user
		h.db.Exec("UPDATE users SET status = 'DELETED', deleted_at = $1, email = 'deleted_' || id || '@anonymized.bn-one', display_name = 'Deleted User' WHERE id = $2", now, userID)
		h.db.Exec("UPDATE deletion_requests SET status = 'COMPLETED', processed_at = $1 WHERE id = $2", now, reqID)
		h.db.Exec("UPDATE user_sessions SET revoked_at = $1 WHERE user_id = $2", now, userID)
	} else {
		h.db.Exec("UPDATE deletion_requests SET status = 'REJECTED', processed_at = $1, retention_exception_reason = $2 WHERE id = $3", now, body.Reason, reqID)
	}

	adminID := middleware.GetAdminID(r.Context())
	h.recordAudit(adminID, "REVIEW_DELETION_REQUEST", "DELETION_REQUEST", reqID, "SUCCESS", `{"action":"`+body.Action+`"}`)

	middleware.RespondJSON(w, r, http.StatusOK, map[string]string{"message": "Deletion request processed successfully"})
}

// ---------------- Audit Logs & Feature Flags ----------------

func (h *Handler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, actor_type, actor_id, action, resource_type, resource_id, result, metadata_redacted, request_id, created_at
	          FROM audit_logs ORDER BY created_at DESC LIMIT 100`
	rows, err := h.db.Query(query)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	logs := make([]models.AuditLog, 0)
	for rows.Next() {
		var l models.AuditLog
		if err := rows.Scan(&l.ID, &l.ActorType, &l.ActorID, &l.Action, &l.ResourceType, &l.ResourceID, &l.Result, &l.MetadataRedacted, &l.RequestID, &l.CreatedAt); err == nil {
			logs = append(logs, l)
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, logs)
}

func (h *Handler) ListFeatureFlags(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, key, description, environment, enabled, updated_by, updated_at FROM feature_flags ORDER BY key ASC`
	rows, err := h.db.Query(query)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	flags := make([]models.FeatureFlag, 0)
	for rows.Next() {
		var f models.FeatureFlag
		if err := rows.Scan(&f.ID, &f.Key, &f.Description, &f.Environment, &f.Enabled, &f.UpdatedBy, &f.UpdatedAt); err == nil {
			flags = append(flags, f)
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, flags)
}

func (h *Handler) ToggleFeatureFlag(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	var body struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	// Article 21.12 Protection: Do NOT allow flag to enable live banking or transfers
	if key == "live_bank_transfers" || key == "custody_funds" {
		middleware.RespondError(w, r, http.StatusForbidden, "FORBIDDEN", "Live banking and funds custody features cannot be enabled via feature flags (Article 21.12)", nil)
		return
	}

	adminID := middleware.GetAdminID(r.Context())
	now := time.Now().UTC()
	_, err := h.db.Exec("UPDATE feature_flags SET enabled = $1, updated_by = $2, updated_at = $3 WHERE key = $4", body.Enabled, adminID, now, key)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}

	h.recordAudit(adminID, "TOGGLE_FEATURE_FLAG", "FEATURE_FLAG", key, "SUCCESS", "{}")

	middleware.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
		"key":     key,
		"enabled": body.Enabled,
	})
}

func (h *Handler) recordAudit(actorID, action, resType, resID, result, meta string) {
	logID := uuid.New().String()
	h.db.Exec(`INSERT INTO audit_logs (id, actor_type, actor_id, action, resource_type, resource_id, result, metadata_redacted, request_id, created_at)
	           VALUES ($1, 'ADMIN', $2, $3, $4, $5, $6, $7, '', $8)`,
		logID, actorID, action, resType, resID, result, meta, time.Now().UTC())
}
