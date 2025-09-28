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
		api.GET("/health", handler.Health)

		// 工作流路由组
		workflow := api.Group("/workflow")
		{
			// 流程定义管理
			workflows := workflow.Group("/workflows")
			{
				workflows.POST("", middleware.FetchUserIdFromHeader(),
					handler.CreateWorkflow) // 创建流程
				workflows.GET("", handler.GetWorkflows)    // 获取流程列表
				workflows.GET("/:id", handler.GetWorkflow) // 获取流程详情
				workflows.PUT("/:id", middleware.FetchUserIdFromHeader(),
					handler.UpdateWorkflow) // 更新流程
				workflows.DELETE("/:id", middleware.FetchUserIdFromHeader(),
					handler.DeleteWorkflow) // 删除流程
			}

			// 流程实例管理
			instances := workflow.Group("/instances")
			{
				instances.POST("", middleware.FetchUserIdFromHeader(),
					handler.StartWorkflow) // 启动流程实例
				instances.GET("", handler.GetInstances)    // 获取实例列表
				instances.GET("/:id", handler.GetInstance) // 获取实例详情
				instances.PUT("/:id/cancel", middleware.FetchUserIdFromHeader(),
					handler.CancelInstance) // 取消实例
				instances.GET("/user", middleware.FetchUserIdFromHeader(),
					handler.GetInstanceByUserID) // 获取用户创建的实例详情
			}

			// 审批任务管理
			tasks := workflow.Group("/tasks")
			{
				tasks.GET("", middleware.FetchUserIdFromHeader(),
					handler.GetMyTasks) // 获取我的待办任务
				tasks.POST("/:id/approve", middleware.FetchUserIdFromHeader(),
					handler.ApproveTask) // 审批通过
				tasks.POST("/:id/reject", middleware.FetchUserIdFromHeader(),
					handler.RejectTask) // 审批拒绝
				tasks.POST("/:id/transfer", middleware.FetchUserIdFromHeader(),
					handler.TransferTask) // 转交任务
			}

			// 通知管理
			// notifications := workflow.Group("/notifications")
			{
				// notifications.GET("", handler.GetNotifications)     // 获取通知列表
				// notifications.PUT("/:id/read", handler.MarkAsRead)  // 标记已读
			}
		}
	}

	return r
}
