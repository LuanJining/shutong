package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "gitee.com/sichuan-shutong-zhihui-data/k-base/api/workflow/openapi"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/client"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/database"
	model "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/models"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/server"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/workflow/config"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/workflow/router"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	dbCfg := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}
	logCfg := database.LogConfig{
		DBLogLevel: cfg.Log.DBLogLevel,
	}

	// Workflow 服务需要迁移的模型
	db, err := database.Init(dbCfg, logCfg,
		&model.Workflow{},
		&model.Step{},
		&model.Task{},
	)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化IAM客户端
	iamClient := client.NewIamClient(&cfg.Iam)

	// 初始化路由
	r := router.Setup(cfg, db, iamClient)

	// 启动服务器
	srv := server.New(&cfg.Server, r)

	// 在 goroutine 中启动服务器
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 设置 5 秒的超时时间用于优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭服务器
	if err := srv.Stop(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// 关闭数据库连接
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	log.Println("Server exited")
}
