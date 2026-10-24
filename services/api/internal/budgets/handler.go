package budgets

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

type CreateBudgetRequest struct {
	Name         string                  `json:"name"`
	PeriodType   string                  `json:"period_type"`
	StartDate    string                  `json:"start_date"`
	EndDate      string                  `json:"end_date"`
	CurrencyCode string                  `json:"currency_code"`
	LimitMinor   int64                   `json:"limit_minor"`
	Categories   []CategoryAllocationReq `json:"categories"`
}

type CategoryAllocationReq struct {
	CategoryID           string `json:"category_id"`
	AllocatedAmountMinor int64  `json:"allocated_amount_minor"`
}

func (h *Handler) ListBudgets(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	query := `SELECT id, user_id, name, period_type, start_date, end_date, currency_code, limit_minor, status, created_at, updated_at
	          FROM budgets WHERE user_id = $1 AND status = 'ACTIVE' ORDER BY start_date DESC`
	rows, err := h.db.Query(query, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	budgetList := make([]models.Budget, 0)
	for rows.Next() {
		var b models.Budget
		if err := rows.Scan(&b.ID, &b.UserID, &b.Name, &b.PeriodType, &b.StartDate, &b.EndDate, &b.CurrencyCode, &b.LimitMinor, &b.Status, &b.CreatedAt, &b.UpdatedAt); err == nil {
			// Calculate spent during this budget's period
			b.SpentMinor = h.calculateBudgetSpent(userID, b.ID, b.StartDate, b.EndDate)
			b.RemainingMinor = b.LimitMinor - b.SpentMinor
			if b.LimitMinor > 0 {
				b.UsagePercent = float64(b.SpentMinor) / float64(b.LimitMinor) * 100.0
			}
			budgetList = append(budgetList, b)
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, budgetList)
}

func (h *Handler) CreateBudget(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var req CreateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	if req.Name == "" || req.LimitMinor <= 0 {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Name and positive limit amount required", nil)
		return
	}
	if req.PeriodType == "" {
		req.PeriodType = "MONTHLY"
	}
	if req.CurrencyCode == "" {
		req.CurrencyCode = "KRW"
	}
	now := time.Now().UTC()
	if req.StartDate == "" {
		req.StartDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	}
	if req.EndDate == "" {
		req.EndDate = time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, time.UTC).Format("2006-01-02")
	}

	budgetID := uuid.New().String()

	tx, err := h.db.Begin()
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
	defer tx.Rollback()

	insertQuery := `INSERT INTO budgets (id, user_id, name, period_type, start_date, end_date, currency_code, limit_minor, status, created_at, updated_at)
	                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'ACTIVE', $9, $10)`
	_, err = tx.Exec(insertQuery, budgetID, userID, req.Name, req.PeriodType, req.StartDate, req.EndDate, req.CurrencyCode, req.LimitMinor, now, now)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INSERT_BUDGET_ERROR", err.Error(), nil)
		return
	}

	// Insert category allocations
	for _, c := range req.Categories {
		allocID := uuid.New().String()
		tx.Exec("INSERT INTO budget_categories (id, budget_id, category_id, allocated_amount_minor, created_at) VALUES ($1, $2, $3, $4, $5)", allocID, budgetID, c.CategoryID, c.AllocatedAmountMinor, now)
	}

	// Insert default alert rule (80%)
	alertID := uuid.New().String()
	tx.Exec("INSERT INTO budget_alert_rules (id, budget_id, threshold_percent, channel, enabled) VALUES ($1, $2, 80, 'IN_APP', TRUE)", alertID, budgetID)

	if err := tx.Commit(); err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "COMMIT_ERROR", err.Error(), nil)
		return
	}

	spent := h.calculateBudgetSpent(userID, budgetID, req.StartDate, req.EndDate)
	res := models.Budget{
		ID:             budgetID,
		UserID:         userID,
		Name:           req.Name,
		PeriodType:     req.PeriodType,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		CurrencyCode:   req.CurrencyCode,
		LimitMinor:     req.LimitMinor,
		SpentMinor:     spent,
		RemainingMinor: req.LimitMinor - spent,
		UsagePercent:   0,
		Status:         "ACTIVE",
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	middleware.RespondJSON(w, r, http.StatusCreated, res)
}

func (h *Handler) GetBudget(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	budgetID := chi.URLParam(r, "id")

	var b models.Budget
	query := `SELECT id, user_id, name, period_type, start_date, end_date, currency_code, limit_minor, status, created_at, updated_at
	          FROM budgets WHERE id = $1 AND user_id = $2`
	err := h.db.QueryRow(query, budgetID, userID).Scan(&b.ID, &b.UserID, &b.Name, &b.PeriodType, &b.StartDate, &b.EndDate, &b.CurrencyCode, &b.LimitMinor, &b.Status, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", "Budget not found", nil)
		return
	}

	b.SpentMinor = h.calculateBudgetSpent(userID, b.ID, b.StartDate, b.EndDate)
	b.RemainingMinor = b.LimitMinor - b.SpentMinor
	if b.LimitMinor > 0 {
		b.UsagePercent = float64(b.SpentMinor) / float64(b.LimitMinor) * 100.0
	}

	// Fetch categories
	catRows, err := h.db.Query(`SELECT bc.id, bc.budget_id, bc.category_id, c.name, bc.allocated_amount_minor
	                            FROM budget_categories bc
	                            JOIN categories c ON bc.category_id = c.id
	                            WHERE bc.budget_id = $1`, b.ID)
	if err == nil {
		cats := make([]models.BudgetCategory, 0)
		for catRows.Next() {
			var bc models.BudgetCategory
			if err := catRows.Scan(&bc.ID, &bc.BudgetID, &bc.CategoryID, &bc.CategoryName, &bc.AllocatedAmountMinor); err == nil {
				// spent in this category
				var catSpent sql.NullInt64
				h.db.QueryRow(`SELECT SUM(ts.amount_minor) 
				               FROM transaction_splits ts
				               JOIN transactions t ON ts.transaction_id = t.id
				               WHERE t.user_id = $1 AND ts.category_id = $2 AND t.transaction_type = 'EXPENSE'
				                 AND t.occurred_at >= $3 AND t.occurred_at <= $4 AND t.deleted_at IS NULL`,
					userID, bc.CategoryID, b.StartDate, b.EndDate+" 23:59:59").Scan(&catSpent)
				if catSpent.Valid {
					bc.SpentAmountMinor = catSpent.Int64
				}
				cats = append(cats, bc)
			}
		}
		catRows.Close()
		b.Categories = cats
	}

	middleware.RespondJSON(w, r, http.StatusOK, b)
}

func (h *Handler) UpdateBudget(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	budgetID := chi.URLParam(r, "id")

	var body struct {
		Name       *string `json:"name"`
		LimitMinor *int64  `json:"limit_minor"`
		Status     *string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	now := time.Now().UTC()
	query := `UPDATE budgets SET 
	            name = COALESCE($1, name),
	            limit_minor = COALESCE($2, limit_minor),
	            status = COALESCE($3, status),
	            updated_at = $4
	          WHERE id = $5 AND user_id = $6`
	_, err := h.db.Exec(query, body.Name, body.LimitMinor, body.Status, now, budgetID, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "UPDATE_ERROR", err.Error(), nil)
		return
	}

	h.GetBudget(w, r)
}

func (h *Handler) DeleteBudget(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	budgetID := chi.URLParam(r, "id")

	_, err := h.db.Exec("UPDATE budgets SET status = 'ARCHIVED', updated_at = $1 WHERE id = $2 AND user_id = $3", time.Now().UTC(), budgetID, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DELETE_ERROR", err.Error(), nil)
		return
	}

	middleware.RespondJSON(w, r, http.StatusOK, map[string]string{"message": "Budget archived successfully"})
}

func (h *Handler) GetBudgetProgress(w http.ResponseWriter, r *http.Request) {
	h.GetBudget(w, r)
}

func (h *Handler) calculateBudgetSpent(userID, budgetID, startDate, endDate string) int64 {
	// If budget has specific categories, sum those; otherwise sum all expense in date range
	var countCats int
	h.db.QueryRow("SELECT COUNT(*) FROM budget_categories WHERE budget_id = $1", budgetID).Scan(&countCats)

	var sumSpent sql.NullInt64
	if countCats > 0 {
		query := `SELECT SUM(ts.amount_minor) 
		          FROM transaction_splits ts
		          JOIN transactions t ON ts.transaction_id = t.id
		          JOIN budget_categories bc ON ts.category_id = bc.category_id
		          WHERE t.user_id = $1 AND bc.budget_id = $2 AND t.transaction_type = 'EXPENSE'
		            AND t.occurred_at >= $3 AND t.occurred_at <= $4 AND t.deleted_at IS NULL`
		h.db.QueryRow(query, userID, budgetID, startDate, endDate+" 23:59:59").Scan(&sumSpent)
	} else {
		query := `SELECT SUM(amount_minor)
		          FROM (
		              SELECT -ae.amount_minor_signed as amount_minor
		              FROM transactions t
		              JOIN account_entries ae ON t.id = ae.transaction_id
		              WHERE t.user_id = $1 AND t.transaction_type = 'EXPENSE' AND t.status = 'POSTED'
		                AND t.occurred_at >= $2 AND t.occurred_at <= $3 AND t.deleted_at IS NULL
		          )`
		h.db.QueryRow(query, userID, startDate, endDate+" 23:59:59").Scan(&sumSpent)
	}

	if sumSpent.Valid {
		return sumSpent.Int64
	}
	return 0
}
