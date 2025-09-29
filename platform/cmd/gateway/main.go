package main

import (
	"log"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/configs"
	serverCfg "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/configs"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/router"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/server"
)

func main() {
	// 加载配置
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	serverCfg, err := serverCfg.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化路由
	r := router.Setup(cfg)

	// 启动服务器
	srv := server.New(&serverCfg.Server, r)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
