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

	"github.com/gideonzy/knowledge-base/internal/common/config"
	commondb "github.com/gideonzy/knowledge-base/internal/common/db"
	"github.com/gideonzy/knowledge-base/internal/common/httpserver"
	"github.com/gideonzy/knowledge-base/internal/common/logging"
	"github.com/gideonzy/knowledge-base/internal/common/storage"
	"github.com/gideonzy/knowledge-base/internal/workflow"
	wfhandler "github.com/gideonzy/knowledge-base/internal/workflow/handler"
	wfrepo "github.com/gideonzy/knowledge-base/internal/workflow/repository"
	wfservice "github.com/gideonzy/knowledge-base/internal/workflow/service"
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
	handler := wfhandler.New(svc, logger)

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
