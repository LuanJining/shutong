package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Server wraps http.Server with sane defaults for graceful shutdown.
type Server struct {
	httpServer *http.Server
}

// New creates a new HTTP server configured with timeouts.
func New(host string, port int, handler http.Handler) *Server {
	addr := fmt.Sprintf("%s:%d", host, port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	return &Server{httpServer: srv}
}

// Start begins serving HTTP requests.
func (s *Server) Start() error {
	if s.httpServer == nil {
		return fmt.Errorf("http server is not initialized")
	}
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}
