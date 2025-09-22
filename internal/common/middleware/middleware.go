package middleware

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gideonzy/knowledge-base/internal/common/auth"
)

// RequestIDHeader is the header used to propagate request identifiers.
const RequestIDHeader = "X-Request-ID"

// TokenValidator validates tokens and returns the embedded claims.
type TokenValidator interface {
	Validate(token string) (*auth.Claims, error)
}

// RequestID ensures each request has an ID for tracing.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(RequestIDHeader)
		if id == "" {
			id = generateRequestID(r)
		}
		ctx := context.WithValue(r.Context(), requestIDKey{}, id)
		w.Header().Set(RequestIDHeader, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logging logs basic request metadata.
func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("request completed",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("request_id", RequestIDFromContext(r.Context())),
				slog.Duration("elapsed", time.Since(started)),
			)
		})
	}
}

// RequireAuth validates JWT tokens from the Authorization header.
func RequireAuth(logger *slog.Logger, validator TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if validator == nil {
				http.Error(w, "auth validator not configured", http.StatusInternalServerError)
				return
			}
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}
			claims, err := validator.Validate(token)
			if err != nil {
				logger.Warn("token validation failed",
					slog.String("error", err.Error()),
					slog.String("request_id", RequestIDFromContext(r.Context())))
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), authClaimsKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequestIDFromContext retrieves the request ID.
func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		return id
	}
	return ""
}

// ClaimsFromContext returns JWT claims stored by the auth middleware.
func ClaimsFromContext(ctx context.Context) *auth.Claims {
	if ctx == nil {
		return nil
	}
	if claims, ok := ctx.Value(authClaimsKey{}).(*auth.Claims); ok {
		return claims
	}
	return nil
}

func generateRequestID(r *http.Request) string {
	h := sha1.New()
	_, _ = h.Write([]byte(fmt.Sprintf("%s-%d", r.RemoteAddr, time.Now().UnixNano())))
	return hex.EncodeToString(h.Sum(nil))
}

type requestIDKey struct{}
type authClaimsKey struct{}
