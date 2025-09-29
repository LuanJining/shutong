package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/configs"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config *configs.ServerConfig
	router *gin.Engine
	server *http.Server
}

func New(cfg *configs.ServerConfig, router *gin.Engine) *Server {
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
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  5 * time.Minute,
	}

	fmt.Printf("Server starting on %s\n", addr)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
