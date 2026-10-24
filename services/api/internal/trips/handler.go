package trips

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

type CreateTripRequest struct {
	Name              string `json:"name"`
	CountryCode       string `json:"country_code"`
	CityName          string `json:"city_name"`
	StartDate         string `json:"start_date"`
	EndDate           string `json:"end_date"`
	BaseCurrency      string `json:"base_currency"`
	BudgetAmountMinor int64  `json:"budget_amount_minor"`
	Notes             string `json:"notes"`
}

func (h *Handler) ListTrips(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	status := r.URL.Query().Get("status")

	query := `SELECT t.id, t.user_id, t.name, t.country_code, t.city_name, t.start_date, t.end_date, t.base_currency, t.status, t.notes, t.created_at, t.updated_at
	          FROM trips t
	          WHERE t.user_id = $1 AND ($2 = '' OR t.status = $2)
	          ORDER BY t.start_date DESC`
	rows, err := h.db.Query(query, userID, status)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	trips := make([]models.Trip, 0)
	for rows.Next() {
		var t models.Trip
		if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.CountryCode, &t.CityName, &t.StartDate, &t.EndDate, &t.BaseCurrency, &t.Status, &t.Notes, &t.CreatedAt, &t.UpdatedAt); err == nil {
			// Get budget
			var budgetMinor sql.NullInt64
			h.db.QueryRow("SELECT budget_amount_minor FROM trip_budgets WHERE trip_id = $1 LIMIT 1", t.ID).Scan(&budgetMinor)
			if budgetMinor.Valid {
				t.TotalBudget = budgetMinor.Int64
			}

			// Get spent
			t.TotalSpent = h.calculateTripSpent(t.ID)
			trips = append(trips, t)
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, trips)
}

func (h *Handler) CreateTrip(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var req CreateTripRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	if req.Name == "" || req.CountryCode == "" || req.StartDate == "" || req.EndDate == "" {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Name, country, start and end dates required", nil)
		return
	}
	if req.BaseCurrency == "" {
		req.BaseCurrency = "KRW"
	}

	tripID := uuid.New().String()
	now := time.Now().UTC()

	tx, err := h.db.Begin()
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
	defer tx.Rollback()

	insertTrip := `INSERT INTO trips (id, user_id, name, country_code, city_name, start_date, end_date, base_currency, status, notes, created_at, updated_at)
	               VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'PLANNED', $9, $10, $11)`
	_, err = tx.Exec(insertTrip, tripID, userID, req.Name, req.CountryCode, req.CityName, req.StartDate, req.EndDate, req.BaseCurrency, req.Notes, now, now)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INSERT_TRIP_ERROR", err.Error(), nil)
		return
	}

	// Insert trip budget
	if req.BudgetAmountMinor > 0 {
		budID := uuid.New().String()
		tx.Exec("INSERT INTO trip_budgets (id, trip_id, budget_amount_minor, currency_code, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
			budID, tripID, req.BudgetAmountMinor, req.BaseCurrency, now, now)
	}

	if err := tx.Commit(); err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "COMMIT_ERROR", err.Error(), nil)
		return
	}

	res := models.Trip{
		ID:           tripID,
		UserID:       userID,
		Name:         req.Name,
		CountryCode:  req.CountryCode,
		CityName:     req.CityName,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		BaseCurrency: req.BaseCurrency,
		Status:       "PLANNED",
		Notes:        req.Notes,
		TotalBudget:  req.BudgetAmountMinor,
		TotalSpent:   0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	middleware.RespondJSON(w, r, http.StatusCreated, res)
}

func (h *Handler) GetTrip(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	tripID := chi.URLParam(r, "id")

	var t models.Trip
	query := `SELECT id, user_id, name, country_code, city_name, start_date, end_date, base_currency, status, notes, created_at, updated_at
	          FROM trips WHERE id = $1 AND user_id = $2`
	err := h.db.QueryRow(query, tripID, userID).Scan(&t.ID, &t.UserID, &t.Name, &t.CountryCode, &t.CityName, &t.StartDate, &t.EndDate, &t.BaseCurrency, &t.Status, &t.Notes, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", "Trip not found", nil)
		return
	}

	var budgetMinor sql.NullInt64
	h.db.QueryRow("SELECT budget_amount_minor FROM trip_budgets WHERE trip_id = $1 LIMIT 1", t.ID).Scan(&budgetMinor)
	if budgetMinor.Valid {
		t.TotalBudget = budgetMinor.Int64
	}
	t.TotalSpent = h.calculateTripSpent(t.ID)

	middleware.RespondJSON(w, r, http.StatusOK, t)
}

func (h *Handler) UpdateTrip(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	tripID := chi.URLParam(r, "id")

	var body struct {
		Name              *string `json:"name"`
		CityName          *string `json:"city_name"`
		Status            *string `json:"status"`
		BudgetAmountMinor *int64  `json:"budget_amount_minor"`
		Notes             *string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	now := time.Now().UTC()
	query := `UPDATE trips SET 
	            name = COALESCE($1, name),
	            city_name = COALESCE($2, city_name),
	            status = COALESCE($3, status),
	            notes = COALESCE($4, notes),
	            updated_at = $5
	          WHERE id = $6 AND user_id = $7`
	_, err := h.db.Exec(query, body.Name, body.CityName, body.Status, body.Notes, now, tripID, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "UPDATE_ERROR", err.Error(), nil)
		return
	}

	if body.BudgetAmountMinor != nil {
		h.db.Exec("UPDATE trip_budgets SET budget_amount_minor = $1, updated_at = $2 WHERE trip_id = $3", *body.BudgetAmountMinor, now, tripID)
	}

	h.GetTrip(w, r)
}

func (h *Handler) DeleteTrip(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	tripID := chi.URLParam(r, "id")

	_, err := h.db.Exec("DELETE FROM trips WHERE id = $1 AND user_id = $2", tripID, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DELETE_ERROR", err.Error(), nil)
		return
	}

	middleware.RespondJSON(w, r, http.StatusOK, map[string]string{"message": "Trip deleted"})
}

func (h *Handler) GetTripExpenses(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	tripID := chi.URLParam(r, "id")

	// Ensure trip ownership
	var exists bool
	err := h.db.QueryRow("SELECT 1 FROM trips WHERE id = $1 AND user_id = $2", tripID, userID).Scan(&exists)
	if err != nil || !exists {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", "Trip not found", nil)
		return
	}

	query := `SELECT t.id, t.occurred_at, t.description, t.merchant_name, ts.amount_minor, ts.currency_code, c.name as category_name
	          FROM trip_expense_links tel
	          JOIN transactions t ON tel.transaction_id = t.id
	          LEFT JOIN transaction_splits ts ON t.id = ts.transaction_id
	          LEFT JOIN categories c ON ts.category_id = c.id
	          WHERE tel.trip_id = $1 AND t.deleted_at IS NULL
	          ORDER BY t.occurred_at DESC`
	rows, err := h.db.Query(query, tripID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	type TripExpenseItem struct {
		TransactionID string    `json:"transaction_id"`
		OccurredAt    time.Time `json:"occurred_at"`
		Description   string    `json:"description"`
		MerchantName  string    `json:"merchant_name"`
		AmountMinor   int64     `json:"amount_minor"`
		CurrencyCode  string    `json:"currency_code"`
		CategoryName  string    `json:"category_name"`
	}

	expenses := make([]TripExpenseItem, 0)
	for rows.Next() {
		var item TripExpenseItem
		var catName sql.NullString
		if err := rows.Scan(&item.TransactionID, &item.OccurredAt, &item.Description, &item.MerchantName, &item.AmountMinor, &item.CurrencyCode, &catName); err == nil {
			if catName.Valid {
				item.CategoryName = catName.String
			}
			expenses = append(expenses, item)
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, expenses)
}

func (h *Handler) LinkExpense(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	tripID := chi.URLParam(r, "id")

	var body struct {
		TransactionID string `json:"transaction_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.TransactionID == "" {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Transaction ID required", nil)
		return
	}

	// Verify ownership of both Trip and Transaction
	var tripOwner, txOwner string
	h.db.QueryRow("SELECT user_id FROM trips WHERE id = $1", tripID).Scan(&tripOwner)
	h.db.QueryRow("SELECT user_id FROM transactions WHERE id = $1 AND deleted_at IS NULL", body.TransactionID).Scan(&txOwner)

	if tripOwner != userID || txOwner != userID {
		middleware.RespondError(w, r, http.StatusForbidden, "FORBIDDEN", "Ownership verification failed", nil)
		return
	}

	linkID := uuid.New().String()
	_, err := h.db.Exec("INSERT INTO trip_expense_links (id, trip_id, transaction_id, created_at) VALUES ($1, $2, $3, $4) ON CONFLICT (trip_id, transaction_id) DO NOTHING",
		linkID, tripID, body.TransactionID, time.Now().UTC())
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}

	// Update trip_id on transaction
	h.db.Exec("UPDATE transactions SET trip_id = $1 WHERE id = $2", tripID, body.TransactionID)

	middleware.RespondJSON(w, r, http.StatusOK, map[string]string{"message": "Transaction linked to trip successfully"})
}

func (h *Handler) UnlinkExpense(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	tripID := chi.URLParam(r, "id")
	txID := chi.URLParam(r, "transaction_id")

	// Verify trip ownership
	var tripOwner string
	h.db.QueryRow("SELECT user_id FROM trips WHERE id = $1", tripID).Scan(&tripOwner)
	if tripOwner != userID {
		middleware.RespondError(w, r, http.StatusForbidden, "FORBIDDEN", "Access denied", nil)
		return
	}

	h.db.Exec("DELETE FROM trip_expense_links WHERE trip_id = $1 AND transaction_id = $2", tripID, txID)
	h.db.Exec("UPDATE transactions SET trip_id = NULL WHERE id = $1 AND trip_id = $2", txID, tripID)

	middleware.RespondJSON(w, r, http.StatusOK, map[string]string{"message": "Expense unlinked"})
}

func (h *Handler) GetTripReport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	tripID := chi.URLParam(r, "id")

	var t models.Trip
	query := `SELECT id, user_id, name, country_code, city_name, start_date, end_date, base_currency, status, notes
	          FROM trips WHERE id = $1 AND user_id = $2`
	err := h.db.QueryRow(query, tripID, userID).Scan(&t.ID, &t.UserID, &t.Name, &t.CountryCode, &t.CityName, &t.StartDate, &t.EndDate, &t.BaseCurrency, &t.Status, &t.Notes)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", "Trip not found", nil)
		return
	}

	var budgetMinor int64 = 0
	h.db.QueryRow("SELECT COALESCE(budget_amount_minor, 0) FROM trip_budgets WHERE trip_id = $1 LIMIT 1", tripID).Scan(&budgetMinor)
	totalSpent := h.calculateTripSpent(tripID)

	// Category breakdown
	catQuery := `SELECT c.name, COALESCE(SUM(ts.amount_minor), 0)
	             FROM trip_expense_links tel
	             JOIN transactions t ON tel.transaction_id = t.id
	             JOIN transaction_splits ts ON t.id = ts.transaction_id
	             JOIN categories c ON ts.category_id = c.id
	             WHERE tel.trip_id = $1 AND t.deleted_at IS NULL
	             GROUP BY c.name ORDER BY 2 DESC`
	rows, _ := h.db.Query(catQuery, tripID)
	catStats := make(map[string]int64)
	if rows != nil {
		for rows.Next() {
			var name string
			var sum int64
			if err := rows.Scan(&name, &sum); err == nil {
				catStats[name] = sum
			}
		}
		rows.Close()
	}

	middleware.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
		"trip":                t,
		"total_budget_minor":  budgetMinor,
		"total_spent_minor":   totalSpent,
		"remaining_minor":     budgetMinor - totalSpent,
		"category_breakdown":  catStats,
		"report_currency":     t.BaseCurrency,
		"fx_disclaimer":       "여행 경비의 원본 통화는 보존되며, 보고 통화 환산 시 참고 환율이 적용됩니다.",
		"generated_at":        time.Now().UTC(),
	})
}

func (h *Handler) calculateTripSpent(tripID string) int64 {
	var total sql.NullInt64
	query := `SELECT SUM(ts.amount_minor)
	          FROM trip_expense_links tel
	          JOIN transactions t ON tel.transaction_id = t.id
	          JOIN transaction_splits ts ON t.id = ts.transaction_id
	          WHERE tel.trip_id = $1 AND t.deleted_at IS NULL`
	h.db.QueryRow(query, tripID).Scan(&total)
	if total.Valid {
		return total.Int64
	}
	return 0
}
