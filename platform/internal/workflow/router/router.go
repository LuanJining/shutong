package router

import (
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/workflow/config"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/workflow/handler"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/workflow/middleware"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/workflow/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func Setup(cfg *config.Config, db *gorm.DB) *gin.Engine {
	gin.SetMode(cfg.Gin.Mode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// r.Use(middleware.CORS())

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	// 初始化服务
	workflowService := service.NewWorkflowService(db)
	handler := handler.NewHandler(db, workflowService)

	// API路由组
	api := r.Group("/api/v1")
	{
		// 健康检查
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// 工作流路由组
		workflow := api.Group("/workflow")
		{
			// 流程定义管理
			workflows := workflow.Group("/workflows")
			{
				workflows.POST("", middleware.FetchUserFromHeader(db),
					handler.CreateWorkflow) // 创建流程
			}
		}
	}

	return r
}
