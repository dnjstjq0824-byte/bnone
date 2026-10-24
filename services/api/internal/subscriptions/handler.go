package subscriptions

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

type CreateSubscriptionRequest struct {
	ServiceID          *string `json:"service_id,omitempty"`
	CustomServiceName  *string `json:"custom_service_name,omitempty"`
	PaymentAccountID   *string `json:"payment_account_id,omitempty"`
	AmountMinor        int64   `json:"amount_minor"`
	CurrencyCode       string  `json:"currency_code"`
	BillingCycle       string  `json:"billing_cycle"` // MONTHLY, YEARLY, WEEKLY
	BillingInterval    int     `json:"billing_interval"`
	StartDate          string  `json:"start_date"`
	NextRenewalAt      string  `json:"next_renewal_at"`
	ReminderEnabled    bool    `json:"reminder_enabled"`
	Notes              string  `json:"notes"`
}

func (h *Handler) ListSubscriptionServices(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, name, slug, category, website, logo_asset_id, catalog_status, created_at
	          FROM subscription_services WHERE catalog_status = 'ACTIVE' ORDER BY name ASC`
	rows, err := h.db.Query(query)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	services := make([]models.SubscriptionService, 0)
	for rows.Next() {
		var s models.SubscriptionService
		if err := rows.Scan(&s.ID, &s.Name, &s.Slug, &s.Category, &s.Website, &s.LogoAssetID, &s.CatalogStatus, &s.CreatedAt); err == nil {
			services = append(services, s)
		}
	}
	middleware.RespondJSON(w, r, http.StatusOK, services)
}

func (h *Handler) ListUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	query := `SELECT us.id, us.user_id, us.service_id, COALESCE(ss.name, us.custom_service_name, 'Subscription') as service_name,
	                 us.custom_service_name, us.payment_account_id, COALESCE(fa.name, '') as payment_account_name,
	                 us.amount_minor, us.currency_code, us.billing_cycle, us.billing_interval,
	                 us.start_date, us.next_renewal_at, us.status, us.reminder_enabled, us.notes,
	                 us.created_at, us.updated_at
	          FROM user_subscriptions us
	          LEFT JOIN subscription_services ss ON us.service_id = ss.id
	          LEFT JOIN finance_accounts fa ON us.payment_account_id = fa.id
	          WHERE us.user_id = $1 AND us.status != 'CANCELLED_RECORD'
	          ORDER BY us.next_renewal_at ASC`
	rows, err := h.db.Query(query, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	subs := make([]models.UserSubscription, 0)
	for rows.Next() {
		var s models.UserSubscription
		if err := rows.Scan(&s.ID, &s.UserID, &s.ServiceID, &s.ServiceName, &s.CustomServiceName,
			&s.PaymentAccountID, &s.PaymentAccountName, &s.AmountMinor, &s.CurrencyCode,
			&s.BillingCycle, &s.BillingInterval, &s.StartDate, &s.NextRenewalAt,
			&s.Status, &s.ReminderEnabled, &s.Notes, &s.CreatedAt, &s.UpdatedAt); err == nil {
			subs = append(subs, s)
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, subs)
}

func (h *Handler) CreateUserSubscription(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var req CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	if req.AmountMinor <= 0 {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Amount must be strictly positive", nil)
		return
	}
	if req.ServiceID == nil && (req.CustomServiceName == nil || *req.CustomServiceName == "") {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Either service_id or custom_service_name must be specified", nil)
		return
	}
	if req.BillingCycle == "" {
		req.BillingCycle = "MONTHLY"
	}
	if req.BillingInterval <= 0 {
		req.BillingInterval = 1
	}
	if req.CurrencyCode == "" {
		req.CurrencyCode = "KRW"
	}
	now := time.Now().UTC()
	if req.StartDate == "" {
		req.StartDate = now.Format("2006-01-02")
	}
	if req.NextRenewalAt == "" {
		req.NextRenewalAt = now.AddDate(0, 1, 0).Format("2006-01-02")
	}

	subID := uuid.New().String()
	query := `INSERT INTO user_subscriptions 
	          (id, user_id, service_id, custom_service_name, payment_account_id, amount_minor, currency_code, billing_cycle, billing_interval, start_date, next_renewal_at, status, reminder_enabled, notes, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 'ACTIVE', $12, $13, $14, $15)`
	_, err := h.db.Exec(query, subID, userID, req.ServiceID, req.CustomServiceName, req.PaymentAccountID, req.AmountMinor, req.CurrencyCode, req.BillingCycle, req.BillingInterval, req.StartDate, req.NextRenewalAt, req.ReminderEnabled, req.Notes, now, now)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}

	// Schedule reminder if enabled
	if req.ReminderEnabled {
		remindID := uuid.New().String()
		renewTime, _ := time.Parse("2006-01-02", req.NextRenewalAt)
		schedAt := renewTime.AddDate(0, 0, -3) // 3 days before
		h.db.Exec(`INSERT INTO subscription_reminders (id, subscription_id, scheduled_at, status, deduplication_key, created_at)
		           VALUES ($1, $2, $3, 'PENDING', $4, $5)`,
			remindID, subID, schedAt, subID+"_"+req.NextRenewalAt, now)
	}

	middleware.RespondJSON(w, r, http.StatusCreated, map[string]interface{}{
		"id":              subID,
		"message":         "Subscription saved successfully",
		"next_renewal_at": req.NextRenewalAt,
		"notice":          "BN ONE의 구독 관리는 참고 기록이며 실제 지출 거래를 자동 확정하지 않습니다.",
	})
}

func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	subID := chi.URLParam(r, "id")

	var s models.UserSubscription
	query := `SELECT us.id, us.user_id, us.service_id, COALESCE(ss.name, us.custom_service_name, 'Subscription') as service_name,
	                 us.custom_service_name, us.payment_account_id, COALESCE(fa.name, '') as payment_account_name,
	                 us.amount_minor, us.currency_code, us.billing_cycle, us.billing_interval,
	                 us.start_date, us.next_renewal_at, us.status, us.reminder_enabled, us.notes,
	                 us.created_at, us.updated_at
	          FROM user_subscriptions us
	          LEFT JOIN subscription_services ss ON us.service_id = ss.id
	          LEFT JOIN finance_accounts fa ON us.payment_account_id = fa.id
	          WHERE us.id = $1 AND us.user_id = $2`
	err := h.db.QueryRow(query, subID, userID).Scan(&s.ID, &s.UserID, &s.ServiceID, &s.ServiceName, &s.CustomServiceName,
		&s.PaymentAccountID, &s.PaymentAccountName, &s.AmountMinor, &s.CurrencyCode,
		&s.BillingCycle, &s.BillingInterval, &s.StartDate, &s.NextRenewalAt,
		&s.Status, &s.ReminderEnabled, &s.Notes, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", "Subscription not found", nil)
		return
	}

	middleware.RespondJSON(w, r, http.StatusOK, s)
}

func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	subID := chi.URLParam(r, "id")

	var body struct {
		AmountMinor     *int64  `json:"amount_minor"`
		NextRenewalAt   *string `json:"next_renewal_at"`
		ReminderEnabled *bool   `json:"reminder_enabled"`
		Notes           *string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	now := time.Now().UTC()
	query := `UPDATE user_subscriptions SET
	            amount_minor = COALESCE($1, amount_minor),
	            next_renewal_at = COALESCE($2, next_renewal_at),
	            reminder_enabled = COALESCE($3, reminder_enabled),
	            notes = COALESCE($4, notes),
	            updated_at = $5
	          WHERE id = $6 AND user_id = $7`
	_, err := h.db.Exec(query, body.AmountMinor, body.NextRenewalAt, body.ReminderEnabled, body.Notes, now, subID, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "UPDATE_ERROR", err.Error(), nil)
		return
	}

	h.GetSubscription(w, r)
}

func (h *Handler) PauseSubscription(w http.ResponseWriter, r *http.Request) {
	h.setSubscriptionStatus(w, r, "PAUSED", "Subscription paused in BN ONE")
}

func (h *Handler) ResumeSubscription(w http.ResponseWriter, r *http.Request) {
	h.setSubscriptionStatus(w, r, "ACTIVE", "Subscription resumed in BN ONE")
}

func (h *Handler) CancelRecord(w http.ResponseWriter, r *http.Request) {
	// Cancel record in BN ONE - NOT external service cancellation
	h.setSubscriptionStatus(w, r, "CANCELLED_RECORD", "BN ONE 내부 구독 기록이 해지 처리되었습니다. (실제 서비스 제공사에 자동 해지 요청이 전달되지 않으므로 서비스 웹사이트에서 별도 해지하시기 바랍니다.)")
}

func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	subID := chi.URLParam(r, "id")

	_, err := h.db.Exec("DELETE FROM user_subscriptions WHERE id = $1 AND user_id = $2", subID, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DELETE_ERROR", err.Error(), nil)
		return
	}
	middleware.RespondJSON(w, r, http.StatusOK, map[string]string{"message": "Subscription removed"})
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var monthlyTotal int64 = 0
	var count int = 0
	query := `SELECT amount_minor, billing_cycle FROM user_subscriptions WHERE user_id = $1 AND status = 'ACTIVE'`
	rows, err := h.db.Query(query, userID)
	if err == nil {
		for rows.Next() {
			var amount int64
			var cycle string
			if err := rows.Scan(&amount, &cycle); err == nil {
				count++
				if cycle == "YEARLY" {
					monthlyTotal += amount / 12
				} else if cycle == "WEEKLY" {
					monthlyTotal += amount * 4
				} else {
					monthlyTotal += amount
				}
			}
		}
		rows.Close()
	}

	middleware.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
		"active_count":           count,
		"monthly_estimate_minor": monthlyTotal,
		"yearly_estimate_minor":  monthlyTotal * 12,
		"notice":                 "동일 비용이 실제 거래에도 등록됐다면 대시보드 전체 소비 합계에 중복 반영되지 않습니다.",
	})
}

func (h *Handler) Renewals(w http.ResponseWriter, r *http.Request) {
	h.ListUserSubscriptions(w, r)
}

func (h *Handler) setSubscriptionStatus(w http.ResponseWriter, r *http.Request, status, message string) {
	userID := middleware.GetUserID(r.Context())
	subID := chi.URLParam(r, "id")

	var currentStatus string
	err := h.db.QueryRow("SELECT status FROM user_subscriptions WHERE id = $1 AND user_id = $2", subID, userID).Scan(&currentStatus)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", "Subscription not found", nil)
		return
	}

	now := time.Now().UTC()
	_, err = h.db.Exec("UPDATE user_subscriptions SET status = $1, updated_at = $2 WHERE id = $3 AND user_id = $4", status, now, subID, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "UPDATE_ERROR", err.Error(), nil)
		return
	}

	// Record history
	histID := uuid.New().String()
	h.db.Exec("INSERT INTO subscription_history (id, subscription_id, change_type, before_snapshot, after_snapshot, changed_at) VALUES ($1, $2, $3, $4, $5, $6)",
		histID, subID, "STATUS_CHANGE", currentStatus, status, now)

	middleware.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
		"id":      subID,
		"status":  status,
		"message": message,
	})
}
