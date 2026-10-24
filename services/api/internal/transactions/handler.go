package transactions

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

type CreateTransactionRequest struct {
	TransactionType string                   `json:"transaction_type"` // EXPENSE, INCOME, ADJUSTMENT
	AccountID       string                   `json:"account_id"`
	AmountMinor     int64                    `json:"amount_minor"`
	CurrencyCode    string                   `json:"currency_code"`
	OccurredAt      string                   `json:"occurred_at"`
	Description     string                   `json:"description"`
	MerchantName    string                   `json:"merchant_name"`
	Notes           string                   `json:"notes"`
	TripID          *string                  `json:"trip_id,omitempty"`
	Splits          []CreateSplitItemRequest `json:"splits"`
}

type CreateSplitItemRequest struct {
	CategoryID  string `json:"category_id"`
	AmountMinor int64  `json:"amount_minor"`
	Notes       string `json:"notes"`
}

type ManualTransferRequest struct {
	SourceAccountID        string `json:"source_account_id"`
	DestinationAccountID   string `json:"destination_account_id"`
	SourceAmountMinor      int64  `json:"source_amount_minor"`
	DestinationAmountMinor int64  `json:"destination_amount_minor"`
	SourceCurrency         string `json:"source_currency"`
	DestinationCurrency    string `json:"destination_currency"`
	OccurredAt             string `json:"occurred_at"`
	Description            string `json:"description"`
	Notes                  string `json:"notes"`
}

func (h *Handler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	accountID := r.URL.Query().Get("account_id")
	categoryID := r.URL.Query().Get("category_id")
	txType := r.URL.Query().Get("type")
	tripID := r.URL.Query().Get("trip_id")
	search := r.URL.Query().Get("q")

	query := `SELECT t.id, t.user_id, t.transaction_type, t.occurred_at, t.description, t.merchant_name, t.notes, t.trip_id, t.status, t.version, t.created_at, t.updated_at
	          FROM transactions t
	          WHERE t.user_id = $1 AND t.deleted_at IS NULL`
	args := []interface{}{userID}
	argIdx := 2

	if accountID != "" {
		query += ` AND EXISTS (SELECT 1 FROM account_entries ae WHERE ae.transaction_id = t.id AND ae.account_id = $` + string(rune('0'+argIdx)) + `)`
		args = append(args, accountID)
		argIdx++
	}
	if categoryID != "" {
		query += ` AND EXISTS (SELECT 1 FROM transaction_splits ts WHERE ts.transaction_id = t.id AND ts.category_id = $` + string(rune('0'+argIdx)) + `)`
		args = append(args, categoryID)
		argIdx++
	}
	if txType != "" {
		query += ` AND t.transaction_type = $` + string(rune('0'+argIdx))
		args = append(args, txType)
		argIdx++
	}
	if tripID != "" {
		query += ` AND t.trip_id = $` + string(rune('0'+argIdx))
		args = append(args, tripID)
		argIdx++
	}
	if search != "" {
		query += ` AND (t.description LIKE '%' || $` + string(rune('0'+argIdx)) + ` || '%' OR t.merchant_name LIKE '%' || $` + string(rune('0'+argIdx)) + ` || '%')`
		args = append(args, search)
		argIdx++
	}

	query += ` ORDER BY t.occurred_at DESC LIMIT 100`

	rows, err := h.db.Query(query, args...)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	txList := make([]models.Transaction, 0)
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.TransactionType, &t.OccurredAt, &t.Description, &t.MerchantName, &t.Notes, &t.TripID, &t.Status, &t.Version, &t.CreatedAt, &t.UpdatedAt); err == nil {
			// Populate splits and primary account/amount
			h.populateTransactionDetails(&t)
			txList = append(txList, t)
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, txList)
}

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	if req.AmountMinor <= 0 {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Transaction amount must be strictly positive", nil)
		return
	}
	if req.AccountID == "" {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Account ID is required", nil)
		return
	}
	if req.TransactionType != "EXPENSE" && req.TransactionType != "INCOME" && req.TransactionType != "ADJUSTMENT" {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid transaction type", nil)
		return
	}

	// 1. Verify Account ownership & currency
	var accCurrency string
	err := h.db.QueryRow("SELECT currency_code FROM finance_accounts WHERE id = $1 AND user_id = $2 AND status = 'ACTIVE'", req.AccountID, userID).Scan(&accCurrency)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "ACCOUNT_NOT_FOUND", "Specified account not found or access denied", nil)
		return
	}
	if req.CurrencyCode == "" {
		req.CurrencyCode = accCurrency
	}

	// 2. Validate splits
	var totalSplits int64 = 0
	for _, s := range req.Splits {
		if s.AmountMinor <= 0 {
			middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Split amount must be greater than 0", nil)
			return
		}
		totalSplits += s.AmountMinor
	}
	if len(req.Splits) > 0 && totalSplits != req.AmountMinor {
		middleware.RespondError(w, r, http.StatusBadRequest, "SPLIT_MISMATCH", "Sum of splits must exactly match total amount", nil)
		return
	}

	occurredTime := time.Now().UTC()
	if req.OccurredAt != "" {
		if t, err := time.Parse(time.RFC3339, req.OccurredAt); err == nil {
			occurredTime = t
		}
	}

	// Atomic DB Transaction
	tx, err := h.db.Begin()
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
	defer tx.Rollback()

	txID := uuid.New().String()
	now := time.Now().UTC()

	// Insert Transaction
	insertTx := `INSERT INTO transactions (id, user_id, transaction_type, occurred_at, description, merchant_name, notes, trip_id, status, version, created_at, updated_at)
	             VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'POSTED', 1, $9, $10)`
	_, err = tx.Exec(insertTx, txID, userID, req.TransactionType, occurredTime, req.Description, req.MerchantName, req.Notes, req.TripID, now, now)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INSERT_TX_ERROR", err.Error(), nil)
		return
	}

	// Calculate signed entry
	var signedAmount int64 = req.AmountMinor
	if req.TransactionType == "EXPENSE" {
		signedAmount = -req.AmountMinor
	}

	// Insert Signed Account Entry
	entryID := uuid.New().String()
	insertEntry := `INSERT INTO account_entries (id, transaction_id, account_id, amount_minor_signed, currency_code, created_at)
	                VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(insertEntry, entryID, txID, req.AccountID, signedAmount, req.CurrencyCode, now)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INSERT_ENTRY_ERROR", err.Error(), nil)
		return
	}

	// Insert splits
	for _, s := range req.Splits {
		splitID := uuid.New().String()
		insertSplit := `INSERT INTO transaction_splits (id, transaction_id, category_id, amount_minor, currency_code, notes, created_at)
		                VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err = tx.Exec(insertSplit, splitID, txID, s.CategoryID, s.AmountMinor, req.CurrencyCode, s.Notes, now)
		if err != nil {
			middleware.RespondError(w, r, http.StatusInternalServerError, "INSERT_SPLIT_ERROR", err.Error(), nil)
			return
		}
	}

	// If linked to trip, create trip_expense_links
	if req.TripID != nil && *req.TripID != "" {
		linkID := uuid.New().String()
		tx.Exec("INSERT INTO trip_expense_links (id, trip_id, transaction_id, created_at) VALUES ($1, $2, $3, $4)", linkID, *req.TripID, txID, now)
	}

	if err := tx.Commit(); err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "COMMIT_ERROR", err.Error(), nil)
		return
	}

	resTx := models.Transaction{
		ID:              txID,
		UserID:          userID,
		TransactionType: req.TransactionType,
		OccurredAt:      occurredTime,
		Description:     req.Description,
		MerchantName:    req.MerchantName,
		Notes:           req.Notes,
		TripID:          req.TripID,
		Status:          "POSTED",
		Version:         1,
		CreatedAt:       now,
		UpdatedAt:       now,
		AccountID:       req.AccountID,
		AmountMinor:     req.AmountMinor,
		CurrencyCode:    req.CurrencyCode,
	}

	middleware.RespondJSON(w, r, http.StatusCreated, resTx)
}

func (h *Handler) ManualTransfer(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var req ManualTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	if req.SourceAccountID == "" || req.DestinationAccountID == "" {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Source and destination accounts required", nil)
		return
	}
	if req.SourceAccountID == req.DestinationAccountID {
		middleware.RespondError(w, r, http.StatusBadRequest, "SAME_ACCOUNT", "Source and destination accounts cannot be identical", nil)
		return
	}
	if req.SourceAmountMinor <= 0 || req.DestinationAmountMinor <= 0 {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Amounts must be strictly positive", nil)
		return
	}

	// Verify User owns BOTH accounts
	var srcCurr, dstCurr string
	err := h.db.QueryRow("SELECT currency_code FROM finance_accounts WHERE id = $1 AND user_id = $2 AND status = 'ACTIVE'", req.SourceAccountID, userID).Scan(&srcCurr)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "SOURCE_ACCOUNT_NOT_FOUND", "Source account not found or access denied", nil)
		return
	}

	err = h.db.QueryRow("SELECT currency_code FROM finance_accounts WHERE id = $1 AND user_id = $2 AND status = 'ACTIVE'", req.DestinationAccountID, userID).Scan(&dstCurr)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "DEST_ACCOUNT_NOT_FOUND", "Destination account not found or access denied", nil)
		return
	}

	occurredTime := time.Now().UTC()
	if req.OccurredAt != "" {
		if t, err := time.Parse(time.RFC3339, req.OccurredAt); err == nil {
			occurredTime = t
		}
	}

	desc := req.Description
	if desc == "" {
		desc = "수동 계정 이동"
	}

	// Atomic Execution of Manual Transfer
	tx, err := h.db.Begin()
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
	defer tx.Rollback()

	txID := uuid.New().String()
	now := time.Now().UTC()

	// 1. Transaction Header
	insertTx := `INSERT INTO transactions (id, user_id, transaction_type, occurred_at, description, notes, status, version, created_at, updated_at)
	             VALUES ($1, $2, 'TRANSFER', $3, $4, $5, 'POSTED', 1, $6, $7)`
	_, err = tx.Exec(insertTx, txID, userID, occurredTime, desc, req.Notes, now, now)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INSERT_TX_ERROR", err.Error(), nil)
		return
	}

	// 2. Source Account Signed Entry: Negative
	srcEntryID := uuid.New().String()
	insertSrc := `INSERT INTO account_entries (id, transaction_id, account_id, amount_minor_signed, currency_code, created_at)
	              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(insertSrc, srcEntryID, txID, req.SourceAccountID, -req.SourceAmountMinor, srcCurr, now)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INSERT_SRC_ENTRY_ERROR", err.Error(), nil)
		return
	}

	// 3. Destination Account Signed Entry: Positive
	dstEntryID := uuid.New().String()
	insertDst := `INSERT INTO account_entries (id, transaction_id, account_id, amount_minor_signed, currency_code, created_at)
	              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(insertDst, dstEntryID, txID, req.DestinationAccountID, req.DestinationAmountMinor, dstCurr, now)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INSERT_DST_ENTRY_ERROR", err.Error(), nil)
		return
	}

	if err := tx.Commit(); err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "COMMIT_ERROR", err.Error(), nil)
		return
	}

	middleware.RespondJSON(w, r, http.StatusCreated, map[string]interface{}{
		"transaction_id":           txID,
		"transaction_type":         "TRANSFER",
		"source_account_id":        req.SourceAccountID,
		"destination_account_id":   req.DestinationAccountID,
		"source_amount_minor":      req.SourceAmountMinor,
		"destination_amount_minor": req.DestinationAmountMinor,
		"source_currency":          srcCurr,
		"destination_currency":     dstCurr,
		"occurred_at":              occurredTime,
		"status":                   "POSTED",
		"notice":                   "BN ONE의 거래 기록은 실제 금융기관 거래를 실행하지 않습니다.",
	})
}

func (h *Handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	txID := chi.URLParam(r, "id")

	var t models.Transaction
	query := `SELECT id, user_id, transaction_type, occurred_at, description, merchant_name, notes, trip_id, status, version, created_at, updated_at
	          FROM transactions WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`
	err := h.db.QueryRow(query, txID, userID).Scan(&t.ID, &t.UserID, &t.TransactionType, &t.OccurredAt, &t.Description, &t.MerchantName, &t.Notes, &t.TripID, &t.Status, &t.Version, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", "Transaction not found", nil)
		return
	}

	h.populateTransactionDetails(&t)
	middleware.RespondJSON(w, r, http.StatusOK, t)
}

func (h *Handler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	txID := chi.URLParam(r, "id")

	// Ensure ownership
	var exists bool
	err := h.db.QueryRow("SELECT 1 FROM transactions WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL", txID, userID).Scan(&exists)
	if err != nil || !exists {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", "Transaction not found", nil)
		return
	}

	// Atomic delete: delete entries and soft delete transaction
	now := time.Now().UTC()
	tx, err := h.db.Begin()
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
	defer tx.Rollback()

	tx.Exec("DELETE FROM account_entries WHERE transaction_id = $1", txID)
	tx.Exec("DELETE FROM transaction_splits WHERE transaction_id = $1", txID)
	tx.Exec("DELETE FROM trip_expense_links WHERE transaction_id = $1", txID)
	tx.Exec("UPDATE transactions SET deleted_at = $1, status = 'VOIDED' WHERE id = $2", now, txID)

	if err := tx.Commit(); err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "COMMIT_ERROR", err.Error(), nil)
		return
	}

	middleware.RespondJSON(w, r, http.StatusOK, map[string]string{"message": "Transaction deleted and balances recalculated"})
}

// ---------------- Categories ----------------

func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	catType := r.URL.Query().Get("type")

	query := `SELECT id, owner_user_id, parent_id, code, name, category_type, icon, color, is_system, is_active, sort_order, created_at
	          FROM categories 
	          WHERE (is_system = TRUE OR owner_user_id = $1) AND is_active = TRUE`
	args := []interface{}{userID}
	if catType != "" {
		query += " AND category_type = $2"
		args = append(args, catType)
	}
	query += " ORDER BY sort_order ASC, name ASC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	cats := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.OwnerUserID, &c.ParentID, &c.Code, &c.Name, &c.CategoryType, &c.Icon, &c.Color, &c.IsSystem, &c.IsActive, &c.SortOrder, &c.CreatedAt); err == nil {
			cats = append(cats, c)
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, cats)
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var body struct {
		Name         string  `json:"name"`
		CategoryType string  `json:"category_type"`
		Icon         string  `json:"icon"`
		Color        string  `json:"color"`
		ParentID     *string `json:"parent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	if body.Name == "" || (body.CategoryType != "EXPENSE" && body.CategoryType != "INCOME" && body.CategoryType != "TRANSFER") {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Category name and valid type required", nil)
		return
	}
	if body.Icon == "" {
		body.Icon = "tag"
	}
	if body.Color == "" {
		body.Color = "#3B82F6"
	}

	catID := uuid.New().String()
	now := time.Now().UTC()

	query := `INSERT INTO categories (id, owner_user_id, parent_id, code, name, category_type, icon, color, is_system, is_active, sort_order, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, FALSE, TRUE, 50, $9)`
	_, err := h.db.Exec(query, catID, userID, body.ParentID, "CUSTOM_"+catID[:8], body.Name, body.CategoryType, body.Icon, body.Color, now)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}

	middleware.RespondJSON(w, r, http.StatusCreated, models.Category{
		ID:           catID,
		OwnerUserID:  &userID,
		ParentID:     body.ParentID,
		Code:         "CUSTOM_" + catID[:8],
		Name:         body.Name,
		CategoryType: body.CategoryType,
		Icon:         body.Icon,
		Color:        body.Color,
		IsSystem:     false,
		IsActive:     true,
		SortOrder:    50,
		CreatedAt:    now,
	})
}

// ---------------- Helpers ----------------

func (h *Handler) populateTransactionDetails(t *models.Transaction) {
	// Account entries
	entryRows, err := h.db.Query("SELECT id, transaction_id, account_id, amount_minor_signed, currency_code, created_at FROM account_entries WHERE transaction_id = $1", t.ID)
	if err == nil {
		entries := make([]models.AccountEntry, 0)
		for entryRows.Next() {
			var e models.AccountEntry
			if err := entryRows.Scan(&e.ID, &e.TransactionID, &e.AccountID, &e.AmountMinorSigned, &e.CurrencyCode, &e.CreatedAt); err == nil {
				entries = append(entries, e)
				if t.AccountID == "" {
					t.AccountID = e.AccountID
					t.CurrencyCode = e.CurrencyCode
					if e.AmountMinorSigned < 0 {
						t.AmountMinor = -e.AmountMinorSigned
					} else {
						t.AmountMinor = e.AmountMinorSigned
					}
				}
			}
		}
		entryRows.Close()
		t.Entries = entries
	}

	// Splits
	splitRows, err := h.db.Query(`SELECT s.id, s.transaction_id, s.category_id, c.name, s.amount_minor, s.currency_code, s.notes, s.created_at
	                              FROM transaction_splits s 
	                              LEFT JOIN categories c ON s.category_id = c.id
	                              WHERE s.transaction_id = $1`, t.ID)
	if err == nil {
		splits := make([]models.TransactionSplit, 0)
		for splitRows.Next() {
			var s models.TransactionSplit
			var catName sql.NullString
			if err := splitRows.Scan(&s.ID, &s.TransactionID, &s.CategoryID, &catName, &s.AmountMinor, &s.CurrencyCode, &s.Notes, &s.CreatedAt); err == nil {
				if catName.Valid {
					s.CategoryName = catName.String
					if t.CategoryName == "" {
						t.CategoryName = catName.String
						t.CategoryID = s.CategoryID
					}
				}
				splits = append(splits, s)
			}
		}
		splitRows.Close()
		t.Splits = splits
	}
}

