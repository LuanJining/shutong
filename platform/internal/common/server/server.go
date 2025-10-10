package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ServerConfig 服务器配置接口
type ServerConfig interface {
	GetHost() string
	GetPort() string
}

// Server HTTP 服务器
type Server struct {
	host   string
	port   string
	router *gin.Engine
	server *http.Server
}

// New 创建新的服务器实例
func New(cfg ServerConfig, router *gin.Engine) *Server {
	return &Server{
		host:   cfg.GetHost(),
		port:   cfg.GetPort(),
		router: router,
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

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

// Stop 优雅关闭服务器
func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	fmt.Println("Stopping server...")
	return s.server.Shutdown(ctx)
}
