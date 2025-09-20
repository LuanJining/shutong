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
)

// RequestIDHeader is the header used to propagate request identifiers.
const RequestIDHeader = "X-Request-ID"

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

// RequireAuth is a placeholder middleware that validates the Authorization header is present.
// Services should replace this with proper JWT validation.
func RequireAuth(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}
			// TODO: replace with JWT verification.
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), authTokenKey{}, token)
			logger.Debug("authorization header accepted", slog.String("request_id", RequestIDFromContext(r.Context())))
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

// AuthTokenFromContext retrieves the raw auth token from context.
func AuthTokenFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if token, ok := ctx.Value(authTokenKey{}).(string); ok {
		return token
	}
	return ""
}

func generateRequestID(r *http.Request) string {
	h := sha1.New()
	_, _ = h.Write([]byte(fmt.Sprintf("%s-%d", r.RemoteAddr, time.Now().UnixNano())))
	return hex.EncodeToString(h.Sum(nil))
}

type requestIDKey struct{}
type authTokenKey struct{}
