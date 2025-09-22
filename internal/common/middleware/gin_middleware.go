package middleware

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gideonzy/knowledge-base/internal/common/auth"
	"github.com/gin-gonic/gin"
)

const ginRequestIDKey = "request_id"

// GinRequestID middleware ensures each request has an ID for tracing.
func GinRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(RequestIDHeader)
		if id == "" {
			id = generateRequestIDGin(c.Request)
		}
		c.Set(ginRequestIDKey, id)
		c.Writer.Header().Set(RequestIDHeader, id)
		c.Next()
	}
}

// GinLogging logs basic request metadata using slog.
func GinLogging(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("request completed",
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.Int("status", c.Writer.Status()),
			slog.String("request_id", RequestIDFromGinContext(c)),
			slog.Duration("elapsed", time.Since(start)),
		)
	}
}

// GinRequireAuth validates JWT tokens from the Authorization header.
func GinRequireAuth(logger *slog.Logger, validator TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		if validator == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "auth validator not configured"})
			return
		}
		token := extractBearerToken(c.Request.Header.Get("Authorization"))
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			return
		}
		claims, err := validator.Validate(token)
		if err != nil {
			logger.Warn("token validation failed", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set(authClaimsKeyGin{}, claims)
		c.Next()
	}
}

// ClaimsFromGinContext retrieves JWT claims from Gin context.
func ClaimsFromGinContext(c *gin.Context) *auth.Claims {
	if val, exists := c.Get(authClaimsKeyGin{}); exists {
		if claims, ok := val.(*auth.Claims); ok {
			return claims
		}
	}
	return nil
}

// RequestIDFromGinContext returns the request ID stored in the context.
func RequestIDFromGinContext(c *gin.Context) string {
	if val, exists := c.Get(ginRequestIDKey); exists {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return ""
}

func generateRequestIDGin(r *http.Request) string {
	h := sha1.New()
	_, _ = h.Write([]byte(fmt.Sprintf("%s-%d", r.RemoteAddr, time.Now().UnixNano())))
	return hex.EncodeToString(h.Sum(nil))
}

type authClaimsKeyGin struct{}

func extractBearerToken(header string) string {
	if header == "" {
		return ""
	}
	const prefix = "Bearer "
	if len(header) >= len(prefix) && header[:len(prefix)] == prefix {
		return header[len(prefix):]
	}
	return ""
}
