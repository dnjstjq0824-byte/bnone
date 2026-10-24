package forex

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/bookingnetwork/bn-one/services/api/internal/database"
	"github.com/bookingnetwork/bn-one/services/api/internal/middleware"
	"github.com/bookingnetwork/bn-one/services/api/internal/models"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	db *database.DB
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
}

type ConvertRequest struct {
	SourceCurrency    string `json:"source_currency"`
	TargetCurrency    string `json:"target_currency"`
	SourceAmountMinor int64  `json:"source_amount_minor"`
}

func (h *Handler) ListCurrencies(w http.ResponseWriter, r *http.Request) {
	query := `SELECT code, name, symbol, decimal_digits, is_active FROM currencies WHERE is_active = TRUE ORDER BY code ASC`
	rows, err := h.db.Query(query)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	list := make([]models.Currency, 0)
	for rows.Next() {
		var c models.Currency
		if err := rows.Scan(&c.Code, &c.Name, &c.Symbol, &c.DecimalDigits, &c.IsActive); err == nil {
			list = append(list, c)
		}
	}
	middleware.RespondJSON(w, r, http.StatusOK, list)
}

func (h *Handler) GetRates(w http.ResponseWriter, r *http.Request) {
	base := r.URL.Query().Get("base")
	if base == "" {
		base = "KRW"
	}

	query := `SELECT er.id, er.provider_id, fp.name as provider_name, er.base_currency, er.quote_currency, er.rate, er.effective_at, er.retrieved_at, er.status
	          FROM exchange_rates er
	          JOIN fx_providers fp ON er.provider_id = fp.id
	          WHERE er.base_currency = $1 AND er.status = 'VERIFIED'
	          ORDER BY er.quote_currency ASC, er.effective_at DESC`
	rows, err := h.db.Query(query, base)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	rates := make([]models.ExchangeRate, 0)
	seen := make(map[string]bool)
	for rows.Next() {
		var er models.ExchangeRate
		if err := rows.Scan(&er.ID, &er.ProviderID, &er.ProviderName, &er.BaseCurrency, &er.QuoteCurrency, &er.Rate, &er.EffectiveAt, &er.RetrievedAt, &er.Status); err == nil {
			if !seen[er.QuoteCurrency] {
				rates = append(rates, er)
				seen[er.QuoteCurrency] = true
			}
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
		"base_currency": base,
		"rates":         rates,
		"disclaimer":    "표시 환율은 참고용이며 실제 환전·결제 환율과 다를 수 있습니다.",
		"retrieved_at":  time.Now().UTC(),
	})
}

func (h *Handler) GetRatePair(w http.ResponseWriter, r *http.Request) {
	base := strings.ToUpper(chi.URLParam(r, "base"))
	quote := strings.ToUpper(chi.URLParam(r, "quote"))

	if base == quote {
		middleware.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
			"base_currency":  base,
			"quote_currency": quote,
			"rate":           1.0,
			"status":         "VERIFIED",
			"effective_at":   time.Now().UTC(),
		})
		return
	}

	var rate float64
	var providerName, status string
	var effAt time.Time
	query := `SELECT er.rate, fp.name, er.effective_at, er.status
	          FROM exchange_rates er
	          JOIN fx_providers fp ON er.provider_id = fp.id
	          WHERE er.base_currency = $1 AND er.quote_currency = $2 AND er.status = 'VERIFIED'
	          ORDER BY er.effective_at DESC LIMIT 1`
	err := h.db.QueryRow(query, base, quote).Scan(&rate, &providerName, &effAt, &status)
	if err != nil {
		// Try inverse rate
		queryInv := `SELECT er.rate, fp.name, er.effective_at, er.status
		             FROM exchange_rates er
		             JOIN fx_providers fp ON er.provider_id = fp.id
		             WHERE er.base_currency = $1 AND er.quote_currency = $2 AND er.status = 'VERIFIED'
		             ORDER BY er.effective_at DESC LIMIT 1`
		errInv := h.db.QueryRow(queryInv, quote, base).Scan(&rate, &providerName, &effAt, &status)
		if errInv == nil && rate > 0 {
			rate = 1.0 / rate
		} else {
			middleware.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
				"base_currency":  base,
				"quote_currency": quote,
				"rate":           nil,
				"status":         "NO_VERIFIED_RATE",
				"message":        "No verified reference rate available for this currency pair",
			})
			return
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
		"base_currency":  base,
		"quote_currency": quote,
		"rate":           rate,
		"provider":       providerName,
		"effective_at":   effAt,
		"status":         status,
		"disclaimer":     "표시 환율은 참고용이며 실제 환전·결제 환율과 다를 수 있습니다.",
	})
}

func (h *Handler) ConvertCurrency(w http.ResponseWriter, r *http.Request) {
	var req ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	req.SourceCurrency = strings.ToUpper(strings.TrimSpace(req.SourceCurrency))
	req.TargetCurrency = strings.ToUpper(strings.TrimSpace(req.TargetCurrency))

	if req.SourceCurrency == req.TargetCurrency {
		middleware.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
			"source_currency":        req.SourceCurrency,
			"target_currency":        req.TargetCurrency,
			"source_amount_minor":    req.SourceAmountMinor,
			"converted_amount_minor": req.SourceAmountMinor,
			"rate":                   1.0,
			"status":                 "VERIFIED",
			"effective_at":           time.Now().UTC(),
		})
		return
	}

	var rate float64
	var providerName, status string
	var effAt time.Time
	query := `SELECT er.rate, fp.name, er.effective_at, er.status
	          FROM exchange_rates er
	          JOIN fx_providers fp ON er.provider_id = fp.id
	          WHERE er.base_currency = $1 AND er.quote_currency = $2 AND er.status = 'VERIFIED'
	          ORDER BY er.effective_at DESC LIMIT 1`
	err := h.db.QueryRow(query, req.SourceCurrency, req.TargetCurrency).Scan(&rate, &providerName, &effAt, &status)
	if err != nil {
		// Try inverse
		queryInv := `SELECT er.rate, fp.name, er.effective_at, er.status
		             FROM exchange_rates er
		             JOIN fx_providers fp ON er.provider_id = fp.id
		             WHERE er.base_currency = $1 AND er.quote_currency = $2 AND er.status = 'VERIFIED'
		             ORDER BY er.effective_at DESC LIMIT 1`
		errInv := h.db.QueryRow(queryInv, req.TargetCurrency, req.SourceCurrency).Scan(&rate, &providerName, &effAt, &status)
		if errInv == nil && rate > 0 {
			rate = 1.0 / rate
		} else {
			middleware.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
				"source_currency":        req.SourceCurrency,
				"target_currency":        req.TargetCurrency,
				"source_amount_minor":    req.SourceAmountMinor,
				"converted_amount_minor": 0,
				"rate":                   nil,
				"provider":               nil,
				"rate_effective_at":      nil,
				"status":                 "NO_VERIFIED_RATE",
			})
			return
		}
	}

	converted := int64(float64(req.SourceAmountMinor) * rate)

	middleware.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
		"source_currency":        req.SourceCurrency,
		"target_currency":        req.TargetCurrency,
		"source_amount_minor":    req.SourceAmountMinor,
		"converted_amount_minor": converted,
		"rate":                   rate,
		"provider":               providerName,
		"rate_effective_at":      effAt,
		"status":                 status,
		"disclaimer":             "표시 환율은 참고용이며 실제 환전·결제 환율과 다를 수 있습니다.",
	})
}
