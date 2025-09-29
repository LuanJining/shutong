package router

import (
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/client"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/config"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/handler"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func Setup(cfg *config.Config, db *gorm.DB, minioClient *client.S3Client, workflowClient *client.WorkflowClient, openaiClient *client.OpenAIClient) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(cfg.Gin.Mode)

	r := gin.New()

	// 添加中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// r.Use(middleware.CORS())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	// 创建服务层
	documentService := service.NewDocumentService(db, minioClient, workflowClient, openaiClient)

	// 创建处理器
	documentHandler := handler.NewDocumentHandler(documentService)

	// API路由组
	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// 文档相关路由
		documents := api.Group("/documents")
		{
			documents.POST("/upload", documentHandler.UploadDocument)
			// 文档预览和下载
			documents.GET("/:id/preview", documentHandler.PreviewDocument)
			documents.GET("/:id/download", documentHandler.DownloadDocument)

			documents.POST("/:id/submit", documentHandler.SubmitDocument)

			documents.GET("/:id/info", documentHandler.GetDocument)
			documents.GET(":id/space", documentHandler.GetDocumentsBySpaceId)

			documents.POST("/:id/publish", documentHandler.PublishDocument)

			documents.POST("/:id/chat", documentHandler.ChatDocument)              // space_id
			documents.POST("/:id/chat/stream", documentHandler.ChatDocumentStream) // space_id
		}
	}

	return r
}
