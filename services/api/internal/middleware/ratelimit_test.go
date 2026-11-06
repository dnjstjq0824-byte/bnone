package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter_AllowsBurstThenBlocks(t *testing.T) {
	rl := NewRateLimiter(3, time.Minute)
	ok := 0
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "10.0.0.1:1234"
		rec := httptest.NewRecorder()
		rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(rec, req)
		if rec.Code == http.StatusOK {
			ok++
		}
	}
	if ok != 3 {
		t.Fatalf("expected 3 allowed requests, got %d", ok)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	rec := httptest.NewRecorder()
	rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 after budget exhausted, got %d", rec.Code)
	}
	if rec.Header().Get("Retry-After") == "" {
		t.Fatal("expected Retry-After header on 429")
	}
}

func TestRateLimiter_IsolatesClients(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute)
	serve := func(ip string) int {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = ip + ":1234"
		rec := httptest.NewRecorder()
		rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(rec, req)
		return rec.Code
	}
	if serve("10.0.0.2") != http.StatusOK {
		t.Fatal("first client should pass")
	}
	if serve("10.0.0.2") != http.StatusTooManyRequests {
		t.Fatal("first client should now be limited")
	}
	if serve("10.0.0.3") != http.StatusOK {
		t.Fatal("second client must have an independent budget")
	}
}

func TestRateLimiter_PrefersForwardedFor(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute)
	serve := func(fwd string) int {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "10.9.9.9:1234"
		if fwd != "" {
			req.Header.Set("X-Forwarded-For", fwd)
		}
		rec := httptest.NewRecorder()
		rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(rec, req)
		return rec.Code
	}
	if serve("203.0.113.7") != http.StatusOK {
		t.Fatal("forwarded client should pass once")
	}
	if serve("203.0.113.7, 70.41.3.18") != http.StatusTooManyRequests {
		t.Fatal("same left-most forwarded IP must share the budget")
	}
	if serve("203.0.113.8") != http.StatusOK {
		t.Fatal("different forwarded IP must have its own budget")
	}
}

func TestSecurityHeaders_IncludesCSP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	SecurityHeadersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rec, req)
	if got := rec.Header().Get("Content-Security-Policy"); got == "" {
		t.Fatal("expected Content-Security-Policy header")
	}
	if got := rec.Header().Get("Permissions-Policy"); got == "" {
		t.Fatal("expected Permissions-Policy header")
	}
}
