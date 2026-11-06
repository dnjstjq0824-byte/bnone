package middleware

import (
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RateLimiter is a per-client-IP fixed-window rate limiter. It answers 429
// with a Retry-After header when the window budget is exhausted, using the
// same JSON error envelope as the rest of the API.
type RateLimiter struct {
	maxRequests int
	window      time.Duration

	mu      sync.Mutex
	clients map[string]*clientWindow
}

type clientWindow struct {
	count     int
	expiresAt time.Time
}

// NewRateLimiter builds a limiter allowing maxRequests per window per IP.
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		maxRequests: maxRequests,
		window:      window,
		clients:     make(map[string]*clientWindow),
	}
	go rl.reapLoop()
	return rl
}

// Middleware adapts the limiter to chi-style middleware.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		if !rl.allow(ip) {
			retryAfter := int(rl.window.Seconds())
			if retryAfter < 1 {
				retryAfter = 1
			}
			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			RespondError(w, r, http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests. Please retry shortly.", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) allow(ip string) bool {
	now := time.Now()
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cw, ok := rl.clients[ip]
	if !ok || now.After(cw.expiresAt) {
		rl.clients[ip] = &clientWindow{count: 1, expiresAt: now.Add(rl.window)}
		return true
	}
	if cw.count >= rl.maxRequests {
		return false
	}
	cw.count++
	return true
}

func (rl *RateLimiter) reapLoop() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		rl.mu.Lock()
		for ip, cw := range rl.clients {
			if now.After(cw.expiresAt) {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// clientIP prefers the left-most X-Forwarded-For entry (set by trusted
// proxies / ingress) and falls back to the direct remote address.
func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		for i := 0; i < len(fwd); i++ {
			if fwd[i] == ',' {
				return trimSpace(fwd[:i])
			}
		}
		return trimSpace(fwd)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func trimSpace(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
