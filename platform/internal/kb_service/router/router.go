package router

import (
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/client"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/config"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func Setup(cfg *config.Config, db *gorm.DB, minioClient *client.S3Client, qdrantClient *client.QdrantClient) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(cfg.Gin.Mode)

	r := gin.New()

	// 添加中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	// API路由组
	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	return r
}
