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
	iamhandler "github.com/gideonzy/knowledge-base/internal/iam/handler"
	iamrepo "github.com/gideonzy/knowledge-base/internal/iam/repository"
	iamsvc "github.com/gideonzy/knowledge-base/internal/iam/service"
)

func main() {
	cfg, err := config.Load("IAM")
	if err != nil {
		panic(err)
	}

	logger := logging.New(cfg.Env).With(slog.String("service", "iam"))

	repos := iamrepo.NewInMemoryRepositories()
	svc := iamsvc.New(
		iamrepo.NewUserRepo(repos.Users),
		iamrepo.NewRoleRepo(repos.Roles),
		iamrepo.NewSpaceRepo(repos.Spaces),
		iamrepo.NewPolicyRepo(repos.Policies),
	)

	handler := iamhandler.New(svc, logger)
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
