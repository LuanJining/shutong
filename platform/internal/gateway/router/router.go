package router

import (
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/configs"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/client"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/handler"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(cfg *configs.Config) *gin.Engine {
	gin.SetMode(cfg.Gin.Mode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	iamClient := client.NewIamClient(&cfg.Iam)
	iamHandler := handler.NewIamHandler(&cfg.Iam, iamClient)
	kbHandler := handler.NewKbHandler(&cfg.Kb)
	workflowHandler := handler.NewWorkflowHandler(&cfg.Workflow)
	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		iam := api.Group("/iam")
		{
			// 认证相关路由
			auth := iam.Group("/auth")
			{
				auth.POST("/login", iamHandler.ProxyToIamClient)
				auth.POST("/logout", iamHandler.ProxyToIamClient)
				auth.POST("/refresh", iamHandler.ProxyToIamClient)
				auth.PATCH("/change-password", iamHandler.ProxyToIamClient)
			}

			// 用户管理路由
			users := iam.Group("/users")
			{
				users.GET("", iamHandler.ProxyToIamClient)
				users.GET("/:id", iamHandler.ProxyToIamClient)
				users.POST("", iamHandler.ProxyToIamClient)
				users.POST("/:id/roles", iamHandler.ProxyToIamClient)
				users.PUT("/:id", iamHandler.ProxyToIamClient)
				users.DELETE("/:id", iamHandler.ProxyToIamClient)
			}

			// 角色管理路由
			roles := iam.Group("/roles")
			{
				roles.GET("", iamHandler.ProxyToIamClient)
				roles.GET("/:id", iamHandler.ProxyToIamClient)
				roles.GET("/:id/permissions", iamHandler.ProxyToIamClient)
				roles.POST("", iamHandler.ProxyToIamClient)
				roles.PUT("/:id", iamHandler.ProxyToIamClient)
				roles.DELETE("/:id", iamHandler.ProxyToIamClient)
			}

			// 权限管理路由
			permissions := iam.Group("/permissions")
			{
				permissions.GET("", iamHandler.ProxyToIamClient)
				permissions.GET("/:id", iamHandler.ProxyToIamClient)
				permissions.POST("/check", iamHandler.ProxyToIamClient)
			}

			// 空间管理路由
			spaces := iam.Group("/spaces")
			{
				spaces.GET("", iamHandler.ProxyToIamClient)
				spaces.GET("/:id", iamHandler.ProxyToIamClient)
				spaces.POST("", iamHandler.ProxyToIamClient)
				spaces.PUT("/:id", iamHandler.ProxyToIamClient)
				spaces.DELETE("/:id", iamHandler.ProxyToIamClient)
				spaces.GET("/:id/members", iamHandler.ProxyToIamClient)
				spaces.POST("/:id/members", iamHandler.ProxyToIamClient)
				spaces.DELETE("/:id/members/:user_id", iamHandler.ProxyToIamClient)
			}
		}

		workflow := api.Group("/workflow")
		{
			// 工作流定义管理
			workflow.GET("", workflowHandler.ProxyToWorkflowClient)
			workflow.GET("/workflows/:id", workflowHandler.ProxyToWorkflowClient)
			workflow.POST("", workflowHandler.ProxyToWorkflowClient)
			workflow.PUT("/:id", workflowHandler.ProxyToWorkflowClient)
			workflow.DELETE("/:id", workflowHandler.ProxyToWorkflowClient)

			// 工作流步骤管理
			workflow.POST("/:id/steps", workflowHandler.ProxyToWorkflowClient)
			workflow.PUT("/:id/steps/:step_id", workflowHandler.ProxyToWorkflowClient)
			workflow.DELETE("/:id/steps/:step_id", workflowHandler.ProxyToWorkflowClient)

			// 工作流实例管理
			workflow.POST("/:id/instances", workflowHandler.ProxyToWorkflowClient)
			workflow.GET("/instances", workflowHandler.ProxyToWorkflowClient)
			workflow.GET("/instances/:instance_id", workflowHandler.ProxyToWorkflowClient)
			workflow.PUT("/instances/:instance_id/cancel", workflowHandler.ProxyToWorkflowClient)
			workflow.GET("/instances/user", workflowHandler.ProxyToWorkflowClient)
			// 任务管理
			workflow.GET("/tasks", workflowHandler.ProxyToWorkflowClient)
			workflow.POST("/tasks/:task_id/approve", workflowHandler.ProxyToWorkflowClient)
			workflow.PUT("/tasks/:task_id/reject", workflowHandler.ProxyToWorkflowClient)
			workflow.PUT("/tasks/:task_id/transfer", workflowHandler.ProxyToWorkflowClient)

			// 状态查询
			workflow.GET("/instances/:instance_id/status", workflowHandler.ProxyToWorkflowClient)
		}

		kb := api.Group("/kb")
		{
			kb.POST("/upload", kbHandler.ProxyToKbClient)
			kb.GET("/:id/preview", kbHandler.ProxyToKbClient)
			kb.GET("/:id/download", kbHandler.ProxyToKbClient)

			kb.GET("/:id/info", kbHandler.ProxyToKbClient)
			kb.GET("/:id/space", kbHandler.ProxyToKbClient)
		}
	}

	return r
}
