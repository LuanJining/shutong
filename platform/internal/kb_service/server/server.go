package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config *config.ServerConfig
	router *gin.Engine
	server *http.Server
}

func New(cfg *config.ServerConfig, router *gin.Engine) *Server {
	return &Server{
		config: cfg,
		router: router,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)

	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Server starting on %s\n", addr)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
