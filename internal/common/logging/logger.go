package logging

import (
	"context"
	"log/slog"
	"os"
)

// New returns a slog.Logger configured for the given environment.
func New(env string) *slog.Logger {
	var handler slog.Handler
	switch env {
	case "production", "prod":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	return slog.New(handler)
}

// WithContext attaches the logger to a context.
func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

// FromContext retrieves a slog.Logger from context, defaulting to a no-op logger.
func FromContext(ctx context.Context) *slog.Logger {
	val := ctx.Value(contextKey{})
	if logger, ok := val.(*slog.Logger); ok && logger != nil {
		return logger
	}
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
}

type contextKey struct{}
