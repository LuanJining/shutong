package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/server"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/configs"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/router"
)

func main() {
	// 加载配置
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化路由
	r := router.Setup(cfg)

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

	log.Println("Server exited")
}
