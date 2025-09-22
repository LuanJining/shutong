// Package main bootstraps the workflow service.
//
// @title Knowledge Base Workflow Service API
// @version 0.1.0
// @description Workflow orchestration endpoints for the knowledge base.
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "gitee.com/sichuan-shutong-zhihui-data/k-base/platform/api/openapi"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/platform/internal/common/auth"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/platform/internal/common/config"
	commondb "gitee.com/sichuan-shutong-zhihui-data/k-base/platform/internal/common/db"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/platform/internal/common/httpserver"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/platform/internal/common/logging"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/platform/internal/common/storage"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/platform/internal/workflow"
	wfhandler "gitee.com/sichuan-shutong-zhihui-data/k-base/platform/internal/workflow/handler"
	wfrepo "gitee.com/sichuan-shutong-zhihui-data/k-base/platform/internal/workflow/repository"
	wfservice "gitee.com/sichuan-shutong-zhihui-data/k-base/platform/internal/workflow/service"
)

func main() {
	cfg, err := config.Load("WORKFLOW")
	if err != nil {
		panic(err)
	}

	logger := logging.New(cfg.Env).With(slog.String("service", "workflow"))

	var (
		defRepo  wfrepo.DefinitionRepository
		instRepo wfrepo.InstanceRepository
		sqlDB    *sql.DB
	)

	if cfg.Database.DSN != "" {
		sqlDB, err = commondb.ConnectPostgres(cfg.Database.DSN)
		if err != nil {
			logger.Error("connect postgres failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
		defRepo = wfrepo.NewPostgresDefinitionRepo(sqlDB)
		instRepo = wfrepo.NewPostgresInstanceRepo(sqlDB)
		logger.Info("using postgres repositories")
	} else {
		defRepo = wfrepo.NewDefinitionRepo(storage.NewInMemory[workflow.FlowDefinition]())
		instRepo = wfrepo.NewInstanceRepo(storage.NewInMemory[workflow.FlowInstance]())
		logger.Warn("WORKFLOW_DATABASE_DSN not set, using in-memory repositories")
	}

	if sqlDB != nil {
		defer sqlDB.Close()
	}

	svc := wfservice.New(defRepo, instRepo)
	tokenManager := auth.NewManager(cfg.Auth.JWTSigningKey, cfg.Auth.JWTTTLSeconds)
	handler := wfhandler.New(svc, tokenManager, logger)

	server := httpserver.New(cfg.Server.Host, cfg.Server.Port, handler.Routes())

	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		logger.Info("starting http server", slog.String("addr", addr))
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logger.Error("server stopped unexpectedly", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown server", slog.String("error", err.Error()))
	}
	logger.Info("server stopped")
}
