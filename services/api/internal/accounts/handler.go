package accounts

import (
	"database/sql"
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

// ---------------- Finance Accounts ----------------

func (h *Handler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "ACTIVE"
	}

	query := `SELECT id, user_id, name, account_type, currency_code, opening_balance_minor, 
	                 opening_balance_date, status, color, icon, sort_order, created_at, updated_at, archived_at
	          FROM finance_accounts 
	          WHERE user_id = $1 AND ($2 = 'ALL' OR status = $2)
	          ORDER BY sort_order ASC, created_at ASC`
	rows, err := h.db.Query(query, userID, status)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	accounts := make([]models.FinanceAccount, 0)
	for rows.Next() {
		var a models.FinanceAccount
		if err := rows.Scan(&a.ID, &a.UserID, &a.Name, &a.AccountType, &a.CurrencyCode,
			&a.OpeningBalanceMinor, &a.OpeningBalanceDate, &a.Status, &a.Color, &a.Icon,
			&a.SortOrder, &a.CreatedAt, &a.UpdatedAt, &a.ArchivedAt); err == nil {
			
			// Calculate balance: opening + sum(account_entries)
			balance := h.calculateAccountBalance(a.ID, a.OpeningBalanceMinor)
			a.CalculatedBalanceMinor = balance
			a.SourceType = "USER_ENTERED"
			accounts = append(accounts, a)
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, accounts)
}

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var body struct {
		Name                string `json:"name"`
		AccountType         string `json:"account_type"`
		CurrencyCode        string `json:"currency_code"`
		OpeningBalanceMinor int64  `json:"opening_balance_minor"`
		OpeningBalanceDate  string `json:"opening_balance_date"`
		Color               string `json:"color"`
		Icon                string `json:"icon"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	if body.Name == "" {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Account name is required", nil)
		return
	}
	if body.AccountType == "" {
		body.AccountType = "MANUAL_BANK"
	}
	if body.CurrencyCode == "" {
		body.CurrencyCode = "KRW"
	}
	if body.OpeningBalanceDate == "" {
		body.OpeningBalanceDate = time.Now().Format("2006-01-02")
	}
	if body.Color == "" {
		body.Color = "#2563EB"
	}
	if body.Icon == "" {
		body.Icon = "wallet"
	}

	accID := uuid.New().String()
	now := time.Now().UTC()

	query := `INSERT INTO finance_accounts 
	          (id, user_id, name, account_type, currency_code, opening_balance_minor, opening_balance_date, status, color, icon, sort_order, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, 'ACTIVE', $8, $9, 0, $10, $11)`
	_, err := h.db.Exec(query, accID, userID, body.Name, body.AccountType, body.CurrencyCode, body.OpeningBalanceMinor, body.OpeningBalanceDate, body.Color, body.Icon, now, now)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}

	account := models.FinanceAccount{
		ID:                     accID,
		UserID:                 userID,
		Name:                   body.Name,
		AccountType:            body.AccountType,
		CurrencyCode:           body.CurrencyCode,
		OpeningBalanceMinor:    body.OpeningBalanceMinor,
		OpeningBalanceDate:     body.OpeningBalanceDate,
		CalculatedBalanceMinor: body.OpeningBalanceMinor,
		Status:                 "ACTIVE",
		Color:                  body.Color,
		Icon:                   body.Icon,
		CreatedAt:              now,
		UpdatedAt:              now,
		SourceType:             "USER_ENTERED",
	}

	middleware.RespondJSON(w, r, http.StatusCreated, account)
}

func (h *Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	accID := chi.URLParam(r, "id")

	var a models.FinanceAccount
	query := `SELECT id, user_id, name, account_type, currency_code, opening_balance_minor, 
	                 opening_balance_date, status, color, icon, sort_order, created_at, updated_at, archived_at
	          FROM finance_accounts WHERE id = $1 AND user_id = $2`
	err := h.db.QueryRow(query, accID, userID).Scan(&a.ID, &a.UserID, &a.Name, &a.AccountType, &a.CurrencyCode,
		&a.OpeningBalanceMinor, &a.OpeningBalanceDate, &a.Status, &a.Color, &a.Icon,
		&a.SortOrder, &a.CreatedAt, &a.UpdatedAt, &a.ArchivedAt)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "ACCOUNT_NOT_FOUND", "Account not found", nil)
		return
	}

	a.CalculatedBalanceMinor = h.calculateAccountBalance(a.ID, a.OpeningBalanceMinor)
	a.SourceType = "USER_ENTERED"

	middleware.RespondJSON(w, r, http.StatusOK, a)
}

func (h *Handler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	accID := chi.URLParam(r, "id")

	var body struct {
		Name      *string `json:"name"`
		Color     *string `json:"color"`
		Icon      *string `json:"icon"`
		SortOrder *int    `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	now := time.Now().UTC()
	query := `UPDATE finance_accounts SET 
	            name = COALESCE($1, name),
	            color = COALESCE($2, color),
	            icon = COALESCE($3, icon),
	            sort_order = COALESCE($4, sort_order),
	            updated_at = $5
	          WHERE id = $6 AND user_id = $7`
	_, err := h.db.Exec(query, body.Name, body.Color, body.Icon, body.SortOrder, now, accID, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "UPDATE_ERROR", err.Error(), nil)
		return
	}

	h.GetAccount(w, r)
}

func (h *Handler) ArchiveAccount(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	accID := chi.URLParam(r, "id")

	now := time.Now().UTC()
	query := `UPDATE finance_accounts SET status = 'ARCHIVED', archived_at = $1, updated_at = $1 
	          WHERE id = $2 AND user_id = $3`
	_, err := h.db.Exec(query, now, accID, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}

	middleware.RespondJSON(w, r, http.StatusOK, map[string]string{"message": "Account archived successfully"})
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	accID := chi.URLParam(r, "id")

	var openingBalance int64
	var currency string
	err := h.db.QueryRow("SELECT opening_balance_minor, currency_code FROM finance_accounts WHERE id = $1 AND user_id = $2", accID, userID).Scan(&openingBalance, &currency)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", "Account not found", nil)
		return
	}

	currentBalance := h.calculateAccountBalance(accID, openingBalance)
	middleware.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
		"account_id":               accID,
		"currency_code":            currency,
		"opening_balance_minor":    openingBalance,
		"calculated_balance_minor": currentBalance,
		"source_type":              "USER_ENTERED",
		"as_of":                    time.Now().UTC(),
	})
}

func (h *Handler) GetEntries(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	accID := chi.URLParam(r, "id")

	// Ensure account ownership
	var exists bool
	err := h.db.QueryRow("SELECT 1 FROM finance_accounts WHERE id = $1 AND user_id = $2", accID, userID).Scan(&exists)
	if err != nil || !exists {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", "Account not found", nil)
		return
	}

	query := `SELECT id, transaction_id, account_id, amount_minor_signed, currency_code, created_at 
	          FROM account_entries WHERE account_id = $1 ORDER BY created_at DESC LIMIT 50`
	rows, err := h.db.Query(query, accID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	entries := make([]models.AccountEntry, 0)
	for rows.Next() {
		var e models.AccountEntry
		if err := rows.Scan(&e.ID, &e.TransactionID, &e.AccountID, &e.AmountMinorSigned, &e.CurrencyCode, &e.CreatedAt); err == nil {
			entries = append(entries, e)
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, entries)
}

// ---------------- Manual Assets ----------------

func (h *Handler) ListAssets(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	query := `SELECT id, user_id, name, asset_type, currency_code, estimated_value_minor, valuation_date, valuation_source, notes, status, created_at, updated_at
	          FROM assets WHERE user_id = $1 AND status = 'ACTIVE' ORDER BY created_at DESC`
	rows, err := h.db.Query(query, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	assetList := make([]models.Asset, 0)
	for rows.Next() {
		var a models.Asset
		if err := rows.Scan(&a.ID, &a.UserID, &a.Name, &a.AssetType, &a.CurrencyCode, &a.EstimatedValueMinor, &a.ValuationDate, &a.ValuationSource, &a.Notes, &a.Status, &a.CreatedAt, &a.UpdatedAt); err == nil {
			assetList = append(assetList, a)
		}
	}
	middleware.RespondJSON(w, r, http.StatusOK, assetList)
}

func (h *Handler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var body struct {
		Name                string `json:"name"`
		AssetType           string `json:"asset_type"`
		CurrencyCode        string `json:"currency_code"`
		EstimatedValueMinor int64  `json:"estimated_value_minor"`
		ValuationDate       string `json:"valuation_date"`
		Notes               string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	if body.Name == "" {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Asset name is required", nil)
		return
	}
	if body.AssetType == "" {
		body.AssetType = "OTHER"
	}
	if body.CurrencyCode == "" {
		body.CurrencyCode = "KRW"
	}
	if body.ValuationDate == "" {
		body.ValuationDate = time.Now().Format("2006-01-02")
	}

	assetID := uuid.New().String()
	now := time.Now().UTC()

	query := `INSERT INTO assets (id, user_id, name, asset_type, currency_code, estimated_value_minor, valuation_date, valuation_source, notes, status, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, 'USER_ESTIMATE', $8, 'ACTIVE', $9, $10)`
	_, err := h.db.Exec(query, assetID, userID, body.Name, body.AssetType, body.CurrencyCode, body.EstimatedValueMinor, body.ValuationDate, body.Notes, now, now)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}

	middleware.RespondJSON(w, r, http.StatusCreated, models.Asset{
		ID:                  assetID,
		UserID:              userID,
		Name:                body.Name,
		AssetType:           body.AssetType,
		CurrencyCode:        body.CurrencyCode,
		EstimatedValueMinor: body.EstimatedValueMinor,
		ValuationDate:       body.ValuationDate,
		ValuationSource:     "USER_ESTIMATE",
		Notes:               body.Notes,
		Status:              "ACTIVE",
		CreatedAt:           now,
		UpdatedAt:           now,
	})
}

// ---------------- Manual Liabilities ----------------

func (h *Handler) ListLiabilities(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	query := `SELECT id, user_id, name, liability_type, currency_code, outstanding_amount_minor, valuation_date, interest_rate_optional, notes, status, created_at, updated_at
	          FROM liabilities WHERE user_id = $1 AND status = 'ACTIVE' ORDER BY created_at DESC`
	rows, err := h.db.Query(query, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	liabilities := make([]models.Liability, 0)
	for rows.Next() {
		var l models.Liability
		if err := rows.Scan(&l.ID, &l.UserID, &l.Name, &l.LiabilityType, &l.CurrencyCode, &l.OutstandingAmountMinor, &l.ValuationDate, &l.InterestRate, &l.Notes, &l.Status, &l.CreatedAt, &l.UpdatedAt); err == nil {
			liabilities = append(liabilities, l)
		}
	}
	middleware.RespondJSON(w, r, http.StatusOK, liabilities)
}

func (h *Handler) CreateLiability(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var body struct {
		Name                   string   `json:"name"`
		LiabilityType          string   `json:"liability_type"`
		CurrencyCode           string   `json:"currency_code"`
		OutstandingAmountMinor int64    `json:"outstanding_amount_minor"`
		ValuationDate          string   `json:"valuation_date"`
		InterestRate           *float64 `json:"interest_rate_optional"`
		Notes                  string   `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	if body.Name == "" {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Liability name is required", nil)
		return
	}
	if body.LiabilityType == "" {
		body.LiabilityType = "OTHER"
	}
	if body.CurrencyCode == "" {
		body.CurrencyCode = "KRW"
	}
	if body.ValuationDate == "" {
		body.ValuationDate = time.Now().Format("2006-01-02")
	}

	liabID := uuid.New().String()
	now := time.Now().UTC()

	query := `INSERT INTO liabilities (id, user_id, name, liability_type, currency_code, outstanding_amount_minor, valuation_date, interest_rate_optional, notes, status, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'ACTIVE', $10, $11)`
	_, err := h.db.Exec(query, liabID, userID, body.Name, body.LiabilityType, body.CurrencyCode, body.OutstandingAmountMinor, body.ValuationDate, body.InterestRate, body.Notes, now, now)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}

	middleware.RespondJSON(w, r, http.StatusCreated, models.Liability{
		ID:                     liabID,
		UserID:                 userID,
		Name:                   body.Name,
		LiabilityType:          body.LiabilityType,
		CurrencyCode:           body.CurrencyCode,
		OutstandingAmountMinor: body.OutstandingAmountMinor,
		ValuationDate:          body.ValuationDate,
		InterestRate:           body.InterestRate,
		Notes:                  body.Notes,
		Status:                 "ACTIVE",
		CreatedAt:              now,
		UpdatedAt:              now,
	})
}

// ---------------- Net Worth Calculation ----------------

func (h *Handler) GetNetWorth(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	// Fetch base currency of user
	var baseCurrency string
	err := h.db.QueryRow("SELECT base_currency FROM users WHERE id = $1", userID).Scan(&baseCurrency)
	if err != nil || baseCurrency == "" {
		baseCurrency = "KRW"
	}

	// Sum financial account balances
	var totalAccountBalance int64 = 0
	accRows, err := h.db.Query("SELECT id, opening_balance_minor, currency_code FROM finance_accounts WHERE user_id = $1 AND status = 'ACTIVE'", userID)
	if err == nil {
		for accRows.Next() {
			var id, curr string
			var openBal int64
			if err := accRows.Scan(&id, &openBal, &curr); err == nil {
				bal := h.calculateAccountBalance(id, openBal)
				// Convert to baseCurrency if different
				converted := h.convertToBase(bal, curr, baseCurrency)
				totalAccountBalance += converted
			}
		}
		accRows.Close()
	}

	// Sum manual assets
	var totalManualAssets int64 = 0
	assetRows, err := h.db.Query("SELECT estimated_value_minor, currency_code FROM assets WHERE user_id = $1 AND status = 'ACTIVE'", userID)
	if err == nil {
		for assetRows.Next() {
			var val int64
			var curr string
			if err := assetRows.Scan(&val, &curr); err == nil {
				totalManualAssets += h.convertToBase(val, curr, baseCurrency)
			}
		}
		assetRows.Close()
	}

	// Sum liabilities
	var totalLiabilities int64 = 0
	liabRows, err := h.db.Query("SELECT outstanding_amount_minor, currency_code FROM liabilities WHERE user_id = $1 AND status = 'ACTIVE'", userID)
	if err == nil {
		for liabRows.Next() {
			var out int64
			var curr string
			if err := liabRows.Scan(&out, &curr); err == nil {
				totalLiabilities += h.convertToBase(out, curr, baseCurrency)
			}
		}
		liabRows.Close()
	}

	totalAssets := totalAccountBalance + totalManualAssets
	netWorth := totalAssets - totalLiabilities

	middleware.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
		"base_currency":               baseCurrency,
		"total_assets_minor":          totalAssets,
		"financial_accounts_minor":    totalAccountBalance,
		"other_assets_minor":          totalManualAssets,
		"total_liabilities_minor":     totalLiabilities,
		"net_worth_minor":             netWorth,
		"disclaimer":                  "표시된 자산과 잔액은 사용자가 직접 입력한 기록을 기준으로 합니다.",
		"as_of":                       time.Now().UTC(),
	})
}

func (h *Handler) NetWorthHistory(w http.ResponseWriter, r *http.Request) {
	// History endpoint
	h.GetNetWorth(w, r)
}

// ---------------- Helper Methods ----------------

func (h *Handler) calculateAccountBalance(accountID string, openingBalance int64) int64 {
	var sumEntries sql.NullInt64
	query := `SELECT SUM(amount_minor_signed) FROM account_entries WHERE account_id = $1`
	err := h.db.QueryRow(query, accountID).Scan(&sumEntries)
	if err != nil || !sumEntries.Valid {
		return openingBalance
	}
	return openingBalance + sumEntries.Int64
}

func (h *Handler) convertToBase(amount int64, fromCurr, toCurr string) int64 {
	if fromCurr == toCurr || amount == 0 {
		return amount
	}
	// Look up verified rate in exchange_rates
	var rate float64
	query := `SELECT rate FROM exchange_rates 
	          WHERE base_currency = $1 AND quote_currency = $2 AND status = 'VERIFIED'
	          ORDER BY effective_at DESC LIMIT 1`
	err := h.db.QueryRow(query, fromCurr, toCurr).Scan(&rate)
	if err == nil && rate > 0 {
		return int64(float64(amount) * rate)
	}

	// Try inverse rate
	queryInverse := `SELECT rate FROM exchange_rates 
	                 WHERE base_currency = $1 AND quote_currency = $2 AND status = 'VERIFIED'
	                 ORDER BY effective_at DESC LIMIT 1`
	err = h.db.QueryRow(queryInverse, toCurr, fromCurr).Scan(&rate)
	if err == nil && rate > 0 {
		return int64(float64(amount) / rate)
	}

	// Return 1:1 fallback if unverified
	return amount
}
