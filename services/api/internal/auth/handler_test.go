package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/bookingnetwork/moa/services/api/internal/database"
)

func TestValidEmail(t *testing.T) {
	valid := []string{"user@example.com", "a.b+tag@mail.co.kr", "x@y.io"}
	for _, e := range valid {
		if !validEmail(e) {
			t.Errorf("expected %q to be valid", e)
		}
	}
	invalid := []string{"", "no-at-sign", "a@b", "Display <x@y.com>", "x@y.com "}
	for _, e := range invalid {
		if validEmail(e) {
			t.Errorf("expected %q to be invalid", e)
		}
	}
}

func TestValidPassword(t *testing.T) {
	if !validPassword("password123") {
		t.Error("expected letters+digits password to be valid")
	}
	weak := []string{"short1", "onlyletters", "12345678", "", "        "}
	for _, p := range weak {
		if validPassword(p) {
			t.Errorf("expected %q to be weak", p)
		}
	}
}

func setupAuthTestDB(t *testing.T) *database.DB {
	dbFile := "test_moa_auth_" + time.Now().Format("20060102150405.000") + ".db"
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
	if err := db.Migrate(migrationsPath, fixturesPath); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	return db
}

func postJSON(t *testing.T, handler http.HandlerFunc, path string, body interface{}) *httptest.ResponseRecorder {
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(raw))
	rec := httptest.NewRecorder()
	handler(rec, req)
	return rec
}

func TestAdminLogin_BootstrapDisabledByDefault(t *testing.T) {
	db := setupAuthTestDB(t)
	h := NewHandler(db, "test-jwt-secret-for-unit-tests", false, "")
	rec := postJSON(t, h.AdminLogin, "/api/v1/auth/admin/login", map[string]string{
		"email":    "admin@bookingnetwork.com",
		"password": "anything-at-all",
	})
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 with bootstrap disabled, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestAdminLogin_BootstrapHonoredOnlyWithEnvSecret(t *testing.T) {
	db := setupAuthTestDB(t)
	h := NewHandler(db, "test-jwt-secret-for-unit-tests", true, "env-provided-bootstrap-secret")
	rec := postJSON(t, h.AdminLogin, "/api/v1/auth/admin/login", map[string]string{
		"email":    "admin@bookingnetwork.com",
		"password": "env-provided-bootstrap-secret",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 with correct bootstrap secret, got %d: %s", rec.Code, rec.Body.String())
	}
	rec = postJSON(t, h.AdminLogin, "/api/v1/auth/admin/login", map[string]string{
		"email":    "admin@bookingnetwork.com",
		"password": "wrong-secret",
	})
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 with wrong bootstrap secret, got %d", rec.Code)
	}
}

func TestRegister_RejectsWeakCredentials(t *testing.T) {
	db := setupAuthTestDB(t)
	h := NewHandler(db, "test-jwt-secret-for-unit-tests", false, "")
	rec := postJSON(t, h.Register, "/api/v1/auth/register", map[string]interface{}{
		"email":                   "not-an-email",
		"password":                "password123",
		"required_terms_accepted": true,
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid email, got %d", rec.Code)
	}
	rec = postJSON(t, h.Register, "/api/v1/auth/register", map[string]interface{}{
		"email":                   "newuser@example.com",
		"password":                "onlyletters",
		"required_terms_accepted": true,
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for weak password, got %d", rec.Code)
	}
}
