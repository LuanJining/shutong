package router

import (
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/client"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/middleware"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/workflow/config"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/workflow/handler"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/workflow/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func Setup(cfg *config.Config, db *gorm.DB, iamClient *client.IamClient) *gin.Engine {
	gin.SetMode(cfg.Gin.Mode)
	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	// 初始化服务
	workflowService := service.NewWorkflowService(db, iamClient)
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

				workflows.POST("/:id/start", middleware.FetchUserFromHeader(db),
					handler.StartWorkflow) // 启动流程
			}

			tasks := workflow.Group("/tasks")
			{
				tasks.GET("", middleware.FetchUserFromHeader(db),
					handler.GetTasks) // 获取任务
				tasks.POST("/:id/approve", middleware.FetchUserFromHeader(db),
					handler.ApproveTask) // 审批任务
			}
		}
	}

	return r
}
