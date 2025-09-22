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

	var (
		userRepo   iamrepo.UserRepository
		roleRepo   iamrepo.RoleRepository
		spaceRepo  iamrepo.SpaceRepository
		policyRepo iamrepo.PolicyRepository
		sqlDB      *sql.DB
	)

	if cfg.Database.DSN != "" {
		sqlDB, err = commondb.ConnectPostgres(cfg.Database.DSN)
		if err != nil {
			logger.Error("connect postgres failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
		postgresRepos := iamrepo.NewPostgresRepositories(sqlDB)
		userRepo = postgresRepos.Users
		roleRepo = postgresRepos.Roles
		spaceRepo = postgresRepos.Spaces
		policyRepo = postgresRepos.Policies
		logger.Info("using postgres repositories")
	} else {
		repos := iamrepo.NewInMemoryRepositories()
		userRepo = iamrepo.NewUserRepo(repos.Users)
		roleRepo = iamrepo.NewRoleRepo(repos.Roles)
		spaceRepo = iamrepo.NewSpaceRepo(repos.Spaces)
		policyRepo = iamrepo.NewPolicyRepo(repos.Policies)
		logger.Warn("IAM_DATABASE_DSN not set, using in-memory repositories")
	}

	if sqlDB != nil {
		defer sqlDB.Close()
	}

	svc := iamsvc.New(userRepo, roleRepo, spaceRepo, policyRepo)

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
