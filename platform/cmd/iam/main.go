package main

import (
	"log"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam/config"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam/database"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam/router"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam/server"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	db, err := database.Init(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化路由
	r := router.Setup(cfg, db)

	// 启动服务器
	srv := server.New(&cfg.Server, r)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
