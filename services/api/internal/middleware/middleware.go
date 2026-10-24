package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/bookingnetwork/bn-one/services/api/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	UserIDKey    contextKey = "user_id"
	AdminIDKey   contextKey = "admin_id"
	UserEmailKey contextKey = "user_email"
	AdminRoleKey contextKey = "admin_role"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				RespondError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "Missing or invalid authorization token", nil)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				RespondError(w, r, http.StatusUnauthorized, "INVALID_TOKEN", "Session token is expired or invalid", nil)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				RespondError(w, r, http.StatusUnauthorized, "INVALID_CLAIMS", "Could not parse token claims", nil)
				return
			}

			userID, _ := claims["sub"].(string)
			email, _ := claims["email"].(string)
			if userID == "" {
				RespondError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "Subject missing in token", nil)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserEmailKey, email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminAuthMiddleware(jwtSecret string, requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				RespondError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "Admin authorization required", nil)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				RespondError(w, r, http.StatusUnauthorized, "INVALID_TOKEN", "Admin token expired or invalid", nil)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				RespondError(w, r, http.StatusUnauthorized, "INVALID_CLAIMS", "Could not parse token claims", nil)
				return
			}

			adminID, _ := claims["sub"].(string)
			role, _ := claims["role"].(string)
			isAdmin, _ := claims["is_admin"].(bool)

			if !isAdmin || adminID == "" {
				RespondError(w, r, http.StatusForbidden, "FORBIDDEN", "Admin credentials required", nil)
				return
			}

			if len(requiredRoles) > 0 {
				authorized := false
				for _, req := range requiredRoles {
					if req == role || role == "PLATFORM_SUPER_ADMIN" {
						authorized = true
						break
					}
				}
				if !authorized {
					RespondError(w, r, http.StatusForbidden, "INSUFFICIENT_PERMISSIONS", "You do not have the required operational permissions", nil)
					return
				}
			}

			ctx := context.WithValue(r.Context(), AdminIDKey, adminID)
			ctx = context.WithValue(ctx, AdminRoleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) string {
	if val, ok := ctx.Value(UserIDKey).(string); ok {
		return val
	}
	return ""
}

func GetAdminID(ctx context.Context) string {
	if val, ok := ctx.Value(AdminIDKey).(string); ok {
		return val
	}
	return ""
}

func GetRequestID(ctx context.Context) string {
	if val, ok := ctx.Value(RequestIDKey).(string); ok {
		return val
	}
	return ""
}

func RespondJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	reqID := GetRequestID(r.Context())
	resp := models.ApiResponse{
		Data: data,
		Meta: &models.MetaInfo{
			RequestID:  reqID,
			ServerTime: time.Now().UTC(),
		},
	}
	json.NewEncoder(w).Encode(resp)
}

func RespondError(w http.ResponseWriter, r *http.Request, status int, code, message string, details []interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	reqID := GetRequestID(r.Context())
	resp := models.ApiErrorResponse{
		Error: models.ApiError{
			Code:      code,
			Message:   message,
			Details:   details,
			RequestID: reqID,
		},
	}
	json.NewEncoder(w).Encode(resp)
}
