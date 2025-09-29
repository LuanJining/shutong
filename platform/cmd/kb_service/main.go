package main

import (
	"log"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/client"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/config"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/database"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/router"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/server"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	db, err := database.Init(cfg.Database, cfg.Log)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化minioClient
	minioClient := client.NewS3Client(&cfg.Minio)
	// 初始化qdrantClient
	qdrantClient := client.NewQdrantClient(&cfg.Qdrant)
	// 初始化ocrClient
	ocrClient := client.NewOCRClient(&cfg.OCR)
	// 初始化workflowClient
	workflowClient := client.NewWorkflowClient(&cfg.Workflow)
	// 初始化OpenAI客户端
	openAIClient := client.NewOpenAIClient(&cfg.OpenAI)

	// 初始化路由
	r := router.Setup(cfg, db, minioClient, qdrantClient, ocrClient, workflowClient, openAIClient)

	// 启动服务器
	srv := server.New(&cfg.Server, r)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
