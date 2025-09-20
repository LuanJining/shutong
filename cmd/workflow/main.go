package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gideonzy/knowledge-base/internal/common/config"
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

	defRepo := wfrepo.NewDefinitionRepo(storage.NewInMemory[workflow.FlowDefinition]())
	instRepo := wfrepo.NewInstanceRepo(storage.NewInMemory[workflow.FlowInstance]())

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
