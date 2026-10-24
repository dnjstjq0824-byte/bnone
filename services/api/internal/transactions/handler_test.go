package transactions

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/bookingnetwork/bn-one/services/api/internal/database"
	"github.com/bookingnetwork/bn-one/services/api/internal/middleware"
)

func setupTestDB(t *testing.T) *database.DB {
	dbFile := "test_bn_one_" + time.Now().Format("20060102150405.000") + ".db"
	db, err := database.Connect("sqlite", dbFile)
	if err != nil {
		t.Fatalf("Failed to open test db: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
		os.Remove(dbFile)
	})

	migrationsPath := "../../../../database/migrations"
	fixturesPath := "../../../../database/fixtures"
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		migrationsPath = "../../../database/migrations"
		fixturesPath = "../../../database/fixtures"
	}
	err = db.Migrate(migrationsPath, fixturesPath)
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func TestFinancialIntegrity_ExpenseAndBalance(t *testing.T) {
	db := setupTestDB(t)
	handler := NewHandler(db)

	userID := "user-test-001"
	accountID := "acc-test-001"

	// Insert test user and account
	now := time.Now().UTC()
	db.Exec(`INSERT INTO users (id, email, password_hash, display_name, country_code, locale, base_currency, timezone, status, created_at, updated_at)
	         VALUES ($1, 'test@example.com', 'hash', 'Tester', 'KR', 'ko-KR', 'KRW', 'Asia/Seoul', 'ACTIVE', $2, $2)`, userID, now)

	// Opening balance: 1,000,000 KRW (1000000 minor units)
	db.Exec(`INSERT INTO finance_accounts (id, user_id, name, account_type, currency_code, opening_balance_minor, opening_balance_date, status, color, icon, sort_order, created_at, updated_at)
	         VALUES ($1, $2, '급여계좌', 'MANUAL_BANK', 'KRW', 1000000, '2026-10-01', 'ACTIVE', '#2563EB', 'wallet', 0, $3, $3)`, accountID, userID, now)

	// Create Expense: 24,000 KRW
	reqBody := CreateTransactionRequest{
		TransactionType: "EXPENSE",
		AccountID:       accountID,
		AmountMinor:     24000,
		CurrencyCode:    "KRW",
		OccurredAt:      now.Format(time.RFC3339),
		Description:     "식사",
		Splits: []CreateSplitItemRequest{
			{
				CategoryID:  "cat-exp-food",
				AmountMinor: 24000,
			},
		},
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/transactions", bytes.NewReader(bodyBytes))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.CreateTransaction(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	// Verify Account Entry amount is negative (-24000)
	var signedEntry int64
	err := db.QueryRow("SELECT amount_minor_signed FROM account_entries WHERE account_id = $1", accountID).Scan(&signedEntry)
	if err != nil {
		t.Fatalf("Failed to query entry: %v", err)
	}
	if signedEntry != -24000 {
		t.Fatalf("Expected signed entry to be -24000, got %d", signedEntry)
	}

	// Verify calculated balance is 1000000 - 24000 = 976000
	var sumEntries int64
	db.QueryRow("SELECT SUM(amount_minor_signed) FROM account_entries WHERE account_id = $1", accountID).Scan(&sumEntries)
	expectedBalance := int64(1000000) + sumEntries
	if expectedBalance != 976000 {
		t.Fatalf("Expected calculated balance 976000, got %d", expectedBalance)
	}
}

func TestFinancialIntegrity_ManualTransfer(t *testing.T) {
	db := setupTestDB(t)
	handler := NewHandler(db)

	userID := "user-transfer-001"
	accSrc := "acc-src-001"
	accDst := "acc-dst-001"

	now := time.Now().UTC()
	db.Exec(`INSERT INTO users (id, email, password_hash, display_name, country_code, locale, base_currency, timezone, status, created_at, updated_at)
	         VALUES ($1, 'transfer@example.com', 'hash', 'TransferUser', 'KR', 'ko-KR', 'KRW', 'Asia/Seoul', 'ACTIVE', $2, $2)`, userID, now)

	// Source account: 500,000 KRW
	db.Exec(`INSERT INTO finance_accounts (id, user_id, name, account_type, currency_code, opening_balance_minor, opening_balance_date, status, color, icon, sort_order, created_at, updated_at)
	         VALUES ($1, $2, '현금', 'CASH', 'KRW', 500000, '2026-10-01', 'ACTIVE', '#2563EB', 'wallet', 0, $3, $3)`, accSrc, userID, now)
	// Destination account: 100,000 KRW
	db.Exec(`INSERT INTO finance_accounts (id, user_id, name, account_type, currency_code, opening_balance_minor, opening_balance_date, status, color, icon, sort_order, created_at, updated_at)
	         VALUES ($1, $2, '저축', 'SAVINGS', 'KRW', 100000, '2026-10-01', 'ACTIVE', '#2563EB', 'wallet', 1, $3, $3)`, accDst, userID, now)

	// Transfer 150,000 KRW from cash to savings
	transferReq := ManualTransferRequest{
		SourceAccountID:        accSrc,
		DestinationAccountID:   accDst,
		SourceAmountMinor:      150000,
		DestinationAmountMinor: 150000,
		SourceCurrency:         "KRW",
		DestinationCurrency:    "KRW",
		OccurredAt:             now.Format(time.RFC3339),
		Description:            "저축 이동",
	}
	bodyBytes, _ := json.Marshal(transferReq)
	req := httptest.NewRequest("POST", "/api/v1/transfers/manual", bytes.NewReader(bodyBytes))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ManualTransfer(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	// Verify source balance: 500000 - 150000 = 350000
	var srcEntries int64
	db.QueryRow("SELECT SUM(amount_minor_signed) FROM account_entries WHERE account_id = $1", accSrc).Scan(&srcEntries)
	if 500000+srcEntries != 350000 {
		t.Fatalf("Expected source balance 350000, got %d", 500000+srcEntries)
	}

	// Verify destination balance: 100000 + 150000 = 250000
	var dstEntries int64
	db.QueryRow("SELECT SUM(amount_minor_signed) FROM account_entries WHERE account_id = $1", accDst).Scan(&dstEntries)
	if 100000+dstEntries != 250000 {
		t.Fatalf("Expected destination balance 250000, got %d", 100000+dstEntries)
	}

	// Net Worth across both accounts invariant: (350000 + 250000) == (500000 + 100000) = 600000
	totalBalanceAfter := (500000 + srcEntries) + (100000 + dstEntries)
	if totalBalanceAfter != 600000 {
		t.Fatalf("Net Worth invariant violated! Expected 600000, got %d", totalBalanceAfter)
	}
}

func TestFinancialIntegrity_RejectSameAccountTransfer(t *testing.T) {
	db := setupTestDB(t)
	handler := NewHandler(db)

	transferReq := ManualTransferRequest{
		SourceAccountID:        "acc-same-001",
		DestinationAccountID:   "acc-same-001",
		SourceAmountMinor:      50000,
		DestinationAmountMinor: 50000,
	}
	bodyBytes, _ := json.Marshal(transferReq)
	req := httptest.NewRequest("POST", "/api/v1/transfers/manual", bytes.NewReader(bodyBytes))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-001")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ManualTransfer(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 Bad Request for same account transfer, got %d", w.Code)
	}
}

func TestTenantIsolation_CrossUserAccessForbidden(t *testing.T) {
	db := setupTestDB(t)
	handler := NewHandler(db)

	now := time.Now().UTC()
	// User A
	db.Exec(`INSERT INTO users (id, email, password_hash, display_name, country_code, locale, base_currency, timezone, status, created_at, updated_at)
	         VALUES ('user-a', 'usera@example.com', 'h', 'User A', 'KR', 'ko-KR', 'KRW', 'Asia/Seoul', 'ACTIVE', $1, $1)`, now)
	db.Exec(`INSERT INTO finance_accounts (id, user_id, name, account_type, currency_code, opening_balance_minor, opening_balance_date, status, created_at, updated_at)
	         VALUES ('acc-a', 'user-a', 'A Account', 'MANUAL_BANK', 'KRW', 1000000, '2026-10-01', 'ACTIVE', $1, $1)`, now)

	// User B tries to create an expense charged to User A's account
	reqBody := CreateTransactionRequest{
		TransactionType: "EXPENSE",
		AccountID:       "acc-a",
		AmountMinor:     10000,
		CurrencyCode:    "KRW",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/transactions", bytes.NewReader(bodyBytes))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-b") // Authenticated as User B
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.CreateTransaction(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected 404/Access Denied when user-b tries to spend from user-a's account, got %d: %s", w.Code, w.Body.String())
	}
}
