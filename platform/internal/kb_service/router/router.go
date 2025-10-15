package router

import (
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/middleware"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/client"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/config"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/handler"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func Setup(cfg *config.Config, db *gorm.DB, minioClient *client.S3Client, workflowClient *client.WorkflowClient, openaiClient *client.OpenAIClient, ocrClient *client.PaddleOCRClient, vectorClient *client.QdrantClient) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(cfg.Gin.Mode)

	r := gin.New()

	// 添加中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	// 创建服务层
	documentService := service.NewDocumentService(db, minioClient, workflowClient, openaiClient, ocrClient, vectorClient)

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
			documents.POST("/upload", middleware.FetchUserFromHeader(db), documentHandler.UploadDocument)
			documents.GET("/tag-cloud", documentHandler.GetTagCloud)
			documents.POST("/search", documentHandler.SearchKnowledge)

			documents.GET("/:id/preview", documentHandler.PreviewDocument)
			documents.GET("/:id/info", documentHandler.GetDocument)
			documents.GET("/:id/space", documentHandler.GetDocumentsBySpaceId)
			documents.GET("/homepage", documentHandler.GetHomepageDocuments) // 展示5个知识库，3个二级知识库，每个二级知识库展示6个文档

			documents.DELETE("/:id", middleware.FetchUserFromHeader(db), documentHandler.DeleteDocument)

			documents.POST("/:id/publish", middleware.FetchUserFromHeader(db), documentHandler.PublishDocument)
			documents.POST("/:id/unpublish", middleware.FetchUserFromHeader(db), documentHandler.UnpublishDocument)

			// 重试处理
			documents.POST("/retry-process", middleware.FetchUserFromHeader(db), documentHandler.RetryProcessDocument)

			// 文档对话
			documents.POST("/:id/chat", documentHandler.ChatDocument)
			documents.POST("/:id/chat/stream", documentHandler.ChatDocumentStream)
		}
	}

	return r
}
