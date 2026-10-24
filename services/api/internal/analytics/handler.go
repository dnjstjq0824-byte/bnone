package analytics

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/bookingnetwork/bn-one/services/api/internal/database"
	"github.com/bookingnetwork/bn-one/services/api/internal/middleware"
)

type Handler struct {
	db *database.DB
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
}

type OverviewData struct {
	BaseCurrency       string           `json:"base_currency"`
	CurrentMonthIncome int64            `json:"current_month_income_minor"`
	CurrentMonthExpense int64           `json:"current_month_expense_minor"`
	NetCashflow        int64            `json:"net_cashflow_minor"`
	CategoryBreakdown  []CategoryStat   `json:"category_breakdown"`
	RecentSpendingTrend []MonthlyTrendStat `json:"monthly_trend"`
}

type CategoryStat struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Color        string  `json:"color"`
	Icon         string  `json:"icon"`
	AmountMinor  int64   `json:"amount_minor"`
	Percentage   float64 `json:"percentage"`
}

type MonthlyTrendStat struct {
	MonthName    string `json:"month"` // e.g. "2026-10"
	IncomeMinor  int64  `json:"income_minor"`
	ExpenseMinor int64  `json:"expense_minor"`
}

func (h *Handler) Overview(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var baseCurr string
	err := h.db.QueryRow("SELECT base_currency FROM users WHERE id = $1", userID).Scan(&baseCurr)
	if err != nil || baseCurr == "" {
		baseCurr = "KRW"
	}

	// Determine current month range in UTC
	now := time.Now().UTC()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	// Income in current month
	var currentIncome int64 = 0
	incomeQuery := `SELECT COALESCE(SUM(ae.amount_minor_signed), 0)
	                FROM transactions t
	                JOIN account_entries ae ON t.id = ae.transaction_id
	                WHERE t.user_id = $1 AND t.transaction_type = 'INCOME' AND t.status = 'POSTED' AND t.deleted_at IS NULL
	                  AND t.occurred_at >= $2 AND t.occurred_at < $3`
	h.db.QueryRow(incomeQuery, userID, startOfMonth, endOfMonth).Scan(&currentIncome)

	// Expense in current month (amount_minor_signed is negative for expense)
	var currentExpenseSigned int64 = 0
	expenseQuery := `SELECT COALESCE(SUM(ae.amount_minor_signed), 0)
	                 FROM transactions t
	                 JOIN account_entries ae ON t.id = ae.transaction_id
	                 WHERE t.user_id = $1 AND t.transaction_type = 'EXPENSE' AND t.status = 'POSTED' AND t.deleted_at IS NULL
	                   AND t.occurred_at >= $2 AND t.occurred_at < $3`
	h.db.QueryRow(expenseQuery, userID, startOfMonth, endOfMonth).Scan(&currentExpenseSigned)
	currentExpense := -currentExpenseSigned

	// Category breakdown for current month
	catQuery := `SELECT c.id, c.name, c.color, c.icon, COALESCE(SUM(ts.amount_minor), 0) as total_cat
	             FROM transaction_splits ts
	             JOIN transactions t ON ts.transaction_id = t.id
	             JOIN categories c ON ts.category_id = c.id
	             WHERE t.user_id = $1 AND t.transaction_type = 'EXPENSE' AND t.status = 'POSTED' AND t.deleted_at IS NULL
	               AND t.occurred_at >= $2 AND t.occurred_at < $3
	             GROUP BY c.id, c.name, c.color, c.icon
	             ORDER BY total_cat DESC`
	catRows, err := h.db.Query(catQuery, userID, startOfMonth, endOfMonth)
	cats := make([]CategoryStat, 0)
	if err == nil {
		for catRows.Next() {
			var cs CategoryStat
			if err := catRows.Scan(&cs.CategoryID, &cs.CategoryName, &cs.Color, &cs.Icon, &cs.AmountMinor); err == nil {
				if currentExpense > 0 {
					cs.Percentage = float64(cs.AmountMinor) / float64(currentExpense) * 100.0
				}
				cats = append(cats, cs)
			}
		}
		catRows.Close()
	}

	// Monthly trend for past 6 months
	trend := make([]MonthlyTrendStat, 0)
	for i := 5; i >= 0; i-- {
		mStart := time.Date(now.Year(), now.Month()-time.Month(i), 1, 0, 0, 0, 0, time.UTC)
		mEnd := mStart.AddDate(0, 1, 0)
		monthLabel := mStart.Format("2006-01")

		var inc, expSigned int64
		h.db.QueryRow(incomeQuery, userID, mStart, mEnd).Scan(&inc)
		h.db.QueryRow(expenseQuery, userID, mStart, mEnd).Scan(&expSigned)

		trend = append(trend, MonthlyTrendStat{
			MonthName:    monthLabel,
			IncomeMinor:  inc,
			ExpenseMinor: -expSigned,
		})
	}

	overview := OverviewData{
		BaseCurrency:        baseCurr,
		CurrentMonthIncome:  currentIncome,
		CurrentMonthExpense: currentExpense,
		NetCashflow:         currentIncome - currentExpense,
		CategoryBreakdown:   cats,
		RecentSpendingTrend: trend,
	}

	middleware.RespondJSON(w, r, http.StatusOK, overview)
}

func (h *Handler) CategorySpending(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	now := time.Now().UTC()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	if fromStr != "" {
		if t, err := time.Parse("2006-01-02", fromStr); err == nil {
			start = t
		}
	}
	if toStr != "" {
		if t, err := time.Parse("2006-01-02", toStr); err == nil {
			end = t.AddDate(0, 0, 1)
		}
	}

	catQuery := `SELECT c.id, c.name, c.color, c.icon, COALESCE(SUM(ts.amount_minor), 0) as total_cat
	             FROM transaction_splits ts
	             JOIN transactions t ON ts.transaction_id = t.id
	             JOIN categories c ON ts.category_id = c.id
	             WHERE t.user_id = $1 AND t.transaction_type = 'EXPENSE' AND t.status = 'POSTED' AND t.deleted_at IS NULL
	               AND t.occurred_at >= $2 AND t.occurred_at < $3
	             GROUP BY c.id, c.name, c.color, c.icon
	             ORDER BY total_cat DESC`
	rows, err := h.db.Query(catQuery, userID, start, end)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	var totalSpent int64 = 0
	cats := make([]CategoryStat, 0)
	for rows.Next() {
		var cs CategoryStat
		if err := rows.Scan(&cs.CategoryID, &cs.CategoryName, &cs.Color, &cs.Icon, &cs.AmountMinor); err == nil {
			totalSpent += cs.AmountMinor
			cats = append(cats, cs)
		}
	}

	for i := range cats {
		if totalSpent > 0 {
			cats[i].Percentage = float64(cats[i].AmountMinor) / float64(totalSpent) * 100.0
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
		"total_spent_minor": totalSpent,
		"from":              start,
		"to":                end,
		"categories":        cats,
	})
}

func (h *Handler) MonthlyTrend(w http.ResponseWriter, r *http.Request) {
	h.Overview(w, r)
}

func (h *Handler) IncomeExpense(w http.ResponseWriter, r *http.Request) {
	h.Overview(w, r)
}

func (h *Handler) AccountBalances(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	query := `SELECT id, name, currency_code, opening_balance_minor FROM finance_accounts WHERE user_id = $1 AND status = 'ACTIVE'`
	rows, err := h.db.Query(query, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	type AccBal struct {
		ID                     string `json:"id"`
		Name                   string `json:"name"`
		CurrencyCode           string `json:"currency_code"`
		CalculatedBalanceMinor int64  `json:"calculated_balance_minor"`
	}

	results := make([]AccBal, 0)
	for rows.Next() {
		var id, name, curr string
		var openBal int64
		if err := rows.Scan(&id, &name, &curr, &openBal); err == nil {
			var sumEntries sql.NullInt64
			h.db.QueryRow("SELECT SUM(amount_minor_signed) FROM account_entries WHERE account_id = $1", id).Scan(&sumEntries)
			bal := openBal
			if sumEntries.Valid {
				bal += sumEntries.Int64
			}
			results = append(results, AccBal{
				ID:                     id,
				Name:                   name,
				CurrencyCode:           curr,
				CalculatedBalanceMinor: bal,
			})
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, results)
}
