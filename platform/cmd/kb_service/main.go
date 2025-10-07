package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/database"
	model "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/models"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/server"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/client"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/config"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/router"
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

	// KB Service 需要迁移的模型
	db, err := database.Init(dbCfg, logCfg,
		&model.Document{},
	)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化minioClient
	minioClient := client.NewS3Client(&cfg.Minio)

	// 初始化workflowClient
	workflowClient := client.NewWorkflowClient(&cfg.Workflow)
	// 初始化OpenAI客户端
	openAIClient := client.NewOpenAIClient(&cfg.OpenAI)

	// 初始化路由
	r := router.Setup(cfg, db, minioClient, workflowClient, openAIClient)

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
