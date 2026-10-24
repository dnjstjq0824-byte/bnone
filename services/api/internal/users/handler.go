package users

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

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var user models.User

	query := `SELECT id, email, display_name, country_code, locale, base_currency, timezone, status, email_verified_at, onboarding_completed_at, created_at, updated_at 
	          FROM users WHERE id = $1 AND deleted_at IS NULL`
	row := h.db.QueryRow(query, userID)
	err := row.Scan(&user.ID, &user.Email, &user.DisplayName, &user.CountryCode, &user.Locale, &user.BaseCurrency, &user.Timezone, &user.Status, &user.EmailVerifiedAt, &user.OnboardingCompletedAt, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "USER_NOT_FOUND", "User profile not found", nil)
		return
	}

	middleware.RespondJSON(w, r, http.StatusOK, user)
}

func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var body struct {
		DisplayName        *string `json:"display_name"`
		CountryCode        *string `json:"country_code"`
		Locale             *string `json:"locale"`
		BaseCurrency       *string `json:"base_currency"`
		Timezone           *string `json:"timezone"`
		CompleteOnboarding *bool   `json:"complete_onboarding"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	now := time.Now().UTC()
	var onboardingTime *time.Time
	if body.CompleteOnboarding != nil && *body.CompleteOnboarding {
		onboardingTime = &now
	}

	// Dynamic update or COALESCE update
	query := `UPDATE users SET 
	            display_name = COALESCE($1, display_name),
	            country_code = COALESCE($2, country_code),
	            locale = COALESCE($3, locale),
	            base_currency = COALESCE($4, base_currency),
	            timezone = COALESCE($5, timezone),
	            onboarding_completed_at = COALESCE($6, onboarding_completed_at),
	            updated_at = $7
	          WHERE id = $8 AND deleted_at IS NULL`
	_, err := h.db.Exec(query, body.DisplayName, body.CountryCode, body.Locale, body.BaseCurrency, body.Timezone, onboardingTime, now, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "UPDATE_ERROR", err.Error(), nil)
		return
	}

	h.GetMe(w, r)
}

func (h *Handler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var p models.UserPreferences

	query := `SELECT id, user_id, mask_amounts, dark_mode, start_of_week, notification_settings, analytics_consent, updated_at 
	          FROM user_preferences WHERE user_id = $1`
	err := h.db.QueryRow(query, userID).Scan(&p.ID, &p.UserID, &p.MaskAmounts, &p.DarkMode, &p.StartOfWeek, &p.NotificationSettings, &p.AnalyticsConsent, &p.UpdatedAt)
	if err != nil {
		// Return default preferences
		p = models.UserPreferences{
			UserID:               userID,
			MaskAmounts:          false,
			DarkMode:             false,
			StartOfWeek:          "SUNDAY",
			NotificationSettings: `{"push":true,"email":true,"budget_alerts":true,"renewal_reminders":true}`,
			AnalyticsConsent:     true,
			UpdatedAt:            time.Now().UTC(),
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, p)
}

func (h *Handler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var body struct {
		MaskAmounts          *bool   `json:"mask_amounts"`
		DarkMode             *bool   `json:"dark_mode"`
		StartOfWeek          *string `json:"start_of_week"`
		NotificationSettings *string `json:"notification_settings"`
		AnalyticsConsent     *bool   `json:"analytics_consent"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	now := time.Now().UTC()
	query := `UPDATE user_preferences SET
	            mask_amounts = COALESCE($1, mask_amounts),
	            dark_mode = COALESCE($2, dark_mode),
	            start_of_week = COALESCE($3, start_of_week),
	            notification_settings = COALESCE($4, notification_settings),
	            analytics_consent = COALESCE($5, analytics_consent),
	            updated_at = $6
	          WHERE user_id = $7`
	res, err := h.db.Exec(query, body.MaskAmounts, body.DarkMode, body.StartOfWeek, body.NotificationSettings, body.AnalyticsConsent, now, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "UPDATE_ERROR", err.Error(), nil)
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		// Insert if not exists
		prefID := uuid.New().String()
		insQuery := `INSERT INTO user_preferences (id, user_id, mask_amounts, dark_mode, start_of_week, notification_settings, analytics_consent, updated_at)
		             VALUES ($1, $2, COALESCE($3, FALSE), COALESCE($4, FALSE), COALESCE($5, 'SUNDAY'), COALESCE($6, '{}'), COALESCE($7, TRUE), $8)`
		h.db.Exec(insQuery, prefID, userID, body.MaskAmounts, body.DarkMode, body.StartOfWeek, body.NotificationSettings, body.AnalyticsConsent, now)
	}

	h.GetPreferences(w, r)
}

func (h *Handler) GetSessions(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	query := `SELECT id, user_id, device_id, device_name, platform, expires_at, revoked_at, created_at 
	          FROM user_sessions WHERE user_id = $1 ORDER BY created_at DESC LIMIT 20`
	rows, err := h.db.Query(query, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	sessions := make([]models.UserSession, 0)
	for rows.Next() {
		var s models.UserSession
		if err := rows.Scan(&s.ID, &s.UserID, &s.DeviceID, &s.DeviceName, &s.Platform, &s.ExpiresAt, &s.RevokedAt, &s.CreatedAt); err == nil {
			sessions = append(sessions, s)
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, sessions)
}

func (h *Handler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	sessionID := chi.URLParam(r, "id")

	now := time.Now().UTC()
	_, err := h.db.Exec("UPDATE user_sessions SET revoked_at = $1 WHERE id = $2 AND user_id = $3", now, sessionID, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}

	middleware.RespondJSON(w, r, http.StatusOK, map[string]string{"message": "Session revoked"})
}

func (h *Handler) CreateDataExport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var body struct {
		Format string `json:"format"`
	}
	json.NewDecoder(r.Body).Decode(&body)
	if body.Format == "" {
		body.Format = "CSV"
	}

	jobID := uuid.New().String()
	now := time.Now().UTC()
	expiresAt := now.Add(7 * 24 * time.Hour)

	query := `INSERT INTO data_export_jobs (id, user_id, format, status, requested_at, expires_at)
	          VALUES ($1, $2, $3, 'COMPLETED', $4, $5)`
	_, err := h.db.Exec(query, jobID, userID, body.Format, now, expiresAt)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "EXPORT_ERROR", err.Error(), nil)
		return
	}

	middleware.RespondJSON(w, r, http.StatusAccepted, map[string]interface{}{
		"export_id":    jobID,
		"status":       "COMPLETED",
		"download_url": "/api/v1/me/data-exports/" + jobID + "/download",
		"format":       body.Format,
		"expires_at":   expiresAt,
	})
}

func (h *Handler) GetDataExport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	exportID := chi.URLParam(r, "id")

	var status, format string
	var reqAt time.Time
	query := `SELECT status, format, requested_at FROM data_export_jobs WHERE id = $1 AND user_id = $2`
	err := h.db.QueryRow(query, exportID, userID).Scan(&status, &format, &reqAt)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", "Export job not found", nil)
		return
	}

	middleware.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
		"id":           exportID,
		"status":       status,
		"format":       format,
		"download_url": "/api/v1/me/data-exports/" + exportID + "/download",
		"requested_at": reqAt,
	})
}

func (h *Handler) CreateDeletionRequest(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	reqID := uuid.New().String()
	now := time.Now().UTC()

	query := `INSERT INTO deletion_requests (id, user_id, status, requested_at)
	          VALUES ($1, $2, 'PENDING', $3)`
	_, err := h.db.Exec(query, reqID, userID, now)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DELETION_ERROR", err.Error(), nil)
		return
	}

	// Update user status
	h.db.Exec("UPDATE users SET status = 'DELETION_REQUESTED', updated_at = $1 WHERE id = $2", now, userID)

	middleware.RespondJSON(w, r, http.StatusAccepted, map[string]interface{}{
		"deletion_request_id": reqID,
		"status":              "PENDING",
		"requested_at":        now,
		"message":             "Your account deletion request has been registered and will be processed within the statutory review period.",
	})
}
