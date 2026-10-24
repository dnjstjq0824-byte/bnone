package auth

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/bookingnetwork/bn-one/services/api/internal/database"
	"github.com/bookingnetwork/bn-one/services/api/internal/middleware"
	"github.com/bookingnetwork/bn-one/services/api/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db        *database.DB
	jwtSecret string
}

func NewHandler(db *database.DB, jwtSecret string) *Handler {
	return &Handler{db: db, jwtSecret: jwtSecret}
}

type RegisterRequest struct {
	Email                 string `json:"email"`
	Password              string `json:"password"`
	DisplayName           string `json:"display_name"`
	CountryCode           string `json:"country_code"`
	Locale                string `json:"locale"`
	BaseCurrency          string `json:"base_currency"`
	Timezone              string `json:"timezone"`
	RequiredTermsAccepted bool   `json:"required_terms_accepted"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceID string `json:"device_id"`
	Platform string `json:"platform"`
}

type TokenResponse struct {
	AccessToken              string       `json:"access_token"`
	RefreshToken             string       `json:"refresh_token"`
	ExpiresIn                int64        `json:"expires_in"`
	User                     *models.User `json:"user"`
	EmailVerificationRequired bool        `json:"email_verification_required"`
	OnboardingRequired       bool         `json:"onboarding_required"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_EMAIL", "A valid email address is required", nil)
		return
	}
	if len(req.Password) < 8 {
		middleware.RespondError(w, r, http.StatusBadRequest, "WEAK_PASSWORD", "Password must be at least 8 characters", nil)
		return
	}
	if !req.RequiredTermsAccepted {
		middleware.RespondError(w, r, http.StatusBadRequest, "TERMS_NOT_ACCEPTED", "You must accept the terms of service", nil)
		return
	}

	// Check if email already exists
	var existingID string
	err := h.db.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingID)
	if err == nil {
		middleware.RespondError(w, r, http.StatusConflict, "EMAIL_EXISTS", "An account with this email already exists", nil)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Password processing failed", nil)
		return
	}

	userID := uuid.New().String()
	displayName := req.DisplayName
	if displayName == "" {
		displayName = strings.Split(req.Email, "@")[0]
	}
	countryCode := req.CountryCode
	if countryCode == "" {
		countryCode = "KR"
	}
	locale := req.Locale
	if locale == "" {
		locale = "ko-KR"
	}
	baseCurrency := req.BaseCurrency
	if baseCurrency == "" {
		baseCurrency = "KRW"
	}
	tz := req.Timezone
	if tz == "" {
		tz = "Asia/Seoul"
	}

	now := time.Now().UTC()
	tx, err := h.db.Begin()
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer tx.Rollback()

	// Insert User
	queryUser := `INSERT INTO users (id, email, password_hash, display_name, country_code, locale, base_currency, timezone, status, created_at, updated_at) 
	              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'ACTIVE', $9, $10)`
	if _, err := tx.Exec(queryUser, userID, req.Email, string(hashedPassword), displayName, countryCode, locale, baseCurrency, tz, now, now); err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INSERT_USER_ERROR", err.Error(), nil)
		return
	}

	// Insert User Preferences
	prefID := uuid.New().String()
	queryPref := `INSERT INTO user_preferences (id, user_id, mask_amounts, dark_mode, start_of_week, notification_settings, analytics_consent, created_at, updated_at)
	              VALUES ($1, $2, FALSE, FALSE, 'SUNDAY', '{"push":true,"email":true,"budget_alerts":true,"renewal_reminders":true}', TRUE, $3, $4)`
	if _, err := tx.Exec(queryPref, prefID, userID, now, now); err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INSERT_PREF_ERROR", err.Error(), nil)
		return
	}

	// Insert Consent Record
	consentID := uuid.New().String()
	queryConsent := `INSERT INTO consent_records (id, user_id, consent_type, policy_version, status, consented_at)
	                VALUES ($1, $2, 'TERMS_OF_SERVICE', '1.0', 'GRANTED', $3)`
	if _, err := tx.Exec(queryConsent, consentID, userID, now); err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INSERT_CONSENT_ERROR", err.Error(), nil)
		return
	}

	if err := tx.Commit(); err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "COMMIT_ERROR", err.Error(), nil)
		return
	}

	// Generate Tokens
	tokenResp, err := h.generateSessionTokens(userID, req.Email, "WEB", "Initial Web Session")
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "TOKEN_ERROR", err.Error(), nil)
		return
	}

	tokenResp.User = &models.User{
		ID:           userID,
		Email:        req.Email,
		DisplayName:  displayName,
		CountryCode:  countryCode,
		Locale:       locale,
		BaseCurrency: baseCurrency,
		Timezone:     tz,
		Status:       "ACTIVE",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	tokenResp.EmailVerificationRequired = false
	tokenResp.OnboardingRequired = true

	middleware.RespondJSON(w, r, http.StatusCreated, tokenResp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	var user models.User
	var passwordHash string

	query := `SELECT id, email, password_hash, display_name, country_code, locale, base_currency, timezone, status, email_verified_at, onboarding_completed_at, created_at, updated_at 
	          FROM users WHERE email = $1 AND deleted_at IS NULL`
	row := h.db.QueryRow(query, req.Email)
	err := row.Scan(&user.ID, &user.Email, &passwordHash, &user.DisplayName, &user.CountryCode, &user.Locale, &user.BaseCurrency, &user.Timezone, &user.Status, &user.EmailVerifiedAt, &user.OnboardingCompletedAt, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			middleware.RespondError(w, r, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password", nil)
			return
		}
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		middleware.RespondError(w, r, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password", nil)
		return
	}

	platform := req.Platform
	if platform == "" {
		platform = "WEB"
	}

	tokenResp, err := h.generateSessionTokens(user.ID, user.Email, platform, req.DeviceID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "TOKEN_ERROR", err.Error(), nil)
		return
	}

	tokenResp.User = &user
	tokenResp.EmailVerificationRequired = user.EmailVerifiedAt == nil
	tokenResp.OnboardingRequired = user.OnboardingCompletedAt == nil

	middleware.RespondJSON(w, r, http.StatusOK, tokenResp)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RefreshToken == "" {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Refresh token required", nil)
		return
	}

	hash := sha256.Sum256([]byte(body.RefreshToken))
	hashStr := hex.EncodeToString(hash[:])

	var sessionID, userID, platform, deviceName string
	var expiresAt time.Time
	query := `SELECT id, user_id, platform, device_name, expires_at FROM user_sessions WHERE refresh_token_hash = $1 AND revoked_at IS NULL`
	err := h.db.QueryRow(query, hashStr).Scan(&sessionID, &userID, &platform, &deviceName, &expiresAt)
	if err != nil || time.Now().UTC().After(expiresAt) {
		middleware.RespondError(w, r, http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", "Session expired or revoked", nil)
		return
	}

	// Revoke old session and issue new pair
	now := time.Now().UTC()
	h.db.Exec("UPDATE user_sessions SET revoked_at = $1 WHERE id = $2", now, sessionID)

	var email string
	h.db.QueryRow("SELECT email FROM users WHERE id = $1", userID).Scan(&email)

	tokenResp, err := h.generateSessionTokens(userID, email, platform, deviceName)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "TOKEN_ERROR", err.Error(), nil)
		return
	}

	middleware.RespondJSON(w, r, http.StatusOK, tokenResp)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID != "" {
		h.db.Exec("UPDATE user_sessions SET revoked_at = $1 WHERE user_id = $2 AND revoked_at IS NULL", time.Now().UTC(), userID)
	}
	middleware.RespondJSON(w, r, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func (h *Handler) AdminLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Check if this is super admin initial login or DB admin
	var adminID, fullName, passHash, roleCode string
	query := `SELECT a.id, a.full_name, a.password_hash, r.code 
	          FROM admin_users a
	          JOIN admin_user_roles aur ON a.id = aur.admin_user_id
	          JOIN admin_roles r ON aur.role_id = r.id
	          WHERE a.email = $1 AND a.status = 'ACTIVE' LIMIT 1`
	err := h.db.QueryRow(query, req.Email).Scan(&adminID, &fullName, &passHash, &roleCode)
	if err != nil {
		// If superadmin not yet seeded in admin_users, allow initial bootstrap credential
		if req.Email == "admin@bookingnetwork.com" && req.Password == "AdminSuperPass2026!" {
			adminID = "admin-bootstrap-001"
			fullName = "System Super Admin"
			roleCode = "PLATFORM_SUPER_ADMIN"
		} else {
			middleware.RespondError(w, r, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid admin credentials", nil)
			return
		}
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(passHash), []byte(req.Password)); err != nil {
			middleware.RespondError(w, r, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid admin credentials", nil)
			return
		}
	}

	// Create JWT for Admin
	claims := jwt.MapClaims{
		"sub":       adminID,
		"email":     req.Email,
		"full_name": fullName,
		"role":      roleCode,
		"is_admin":  true,
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(12 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "TOKEN_ERROR", err.Error(), nil)
		return
	}

	middleware.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
		"access_token": tokenStr,
		"role":         roleCode,
		"full_name":    fullName,
		"email":        req.Email,
		"expires_in":   43200,
	})
}

func (h *Handler) generateSessionTokens(userID, email, platform, deviceName string) (*TokenResponse, error) {
	// 24 hour Access Token
	now := time.Now().UTC()
	accessClaims := jwt.MapClaims{
		"sub":      userID,
		"email":    email,
		"platform": platform,
		"iat":      now.Unix(),
		"exp":      now.Add(24 * time.Hour).Unix(),
	}
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessTokenObj.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return nil, err
	}

	// 30 days Refresh Token
	refreshToken := uuid.New().String() + "-" + uuid.New().String()
	hash := sha256.Sum256([]byte(refreshToken))
	refreshTokenHash := hex.EncodeToString(hash[:])

	sessionID := uuid.New().String()
	expiresAt := now.Add(30 * 24 * time.Hour)

	query := `INSERT INTO user_sessions (id, user_id, refresh_token_hash, device_id, device_name, platform, expires_at, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = h.db.Exec(query, sessionID, userID, refreshTokenHash, deviceName, deviceName, platform, expiresAt, now)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    86400,
	}, nil
}

