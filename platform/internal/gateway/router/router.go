package router

import (
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/middleware"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/client"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/configs"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/handler"
	gw_middleware "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(cfg *configs.Config) *gin.Engine {
	gin.SetMode(cfg.Gin.Mode)
	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
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
				auth.PATCH("/change-password",
					gw_middleware.AuthRequired(iamHandler),
					iamHandler.ProxyToIamClient)
			}

			// 用户管理路由
			users := iam.Group("/users")
			users.Use(gw_middleware.AuthRequired(iamHandler))
			{
				users.GET("", iamHandler.ProxyToIamClient)
				users.GET("/by-role/:rid/space/:sid", iamHandler.ProxyToIamClient) // 根据角色和空间获取用户
				users.GET("/:id", iamHandler.ProxyToIamClient)
				users.POST("", iamHandler.ProxyToIamClient)
				users.POST("/:id/roles", iamHandler.ProxyToIamClient)
				users.PUT("/:id", iamHandler.ProxyToIamClient)
				users.DELETE("/:id", iamHandler.ProxyToIamClient)
			}

			// 角色管理路由
			roles := iam.Group("/roles")
			roles.Use(gw_middleware.AuthRequired(iamHandler))
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
			permissions.Use(gw_middleware.AuthRequired(iamHandler))
			{
				permissions.GET("", iamHandler.ProxyToIamClient)
				permissions.GET("/:id", iamHandler.ProxyToIamClient)
				permissions.POST("/check", iamHandler.ProxyToIamClient)
			}

			// 空间管理路由
			spaces := iam.Group("/spaces")
			spaces.Use(gw_middleware.AuthRequired(iamHandler))
			{
				spaces.GET("", iamHandler.ProxyToIamClient)
				spaces.GET("/:id", iamHandler.ProxyToIamClient)
				spaces.POST("", iamHandler.ProxyToIamClient)
				spaces.PUT("/:id", iamHandler.ProxyToIamClient)
				spaces.DELETE("/:id", iamHandler.ProxyToIamClient)
				spaces.GET("/:id/members", iamHandler.ProxyToIamClient)
				spaces.POST("/:id/members", iamHandler.ProxyToIamClient)
				spaces.DELETE("/:id/members/:user_id", iamHandler.ProxyToIamClient)
				spaces.PUT("/:id/members/:user_id", iamHandler.ProxyToIamClient)
				spaces.GET("/:id/members/role/:role", iamHandler.ProxyToIamClient) // 根据角色获取成员

				spaces.POST("/subspaces", iamHandler.ProxyToIamClient)

				spaces.POST("/classes", iamHandler.ProxyToIamClient)

			}

		}

		workflow := api.Group("/workflow")
		workflow.Use(gw_middleware.AuthRequired(iamHandler))
		{
			// 工作流管理
			workflows := workflow.Group("/workflows")
			{
				workflows.POST("", workflowHandler.ProxyToWorkflowClient)           // 创建工作流
				workflows.POST("/:id/start", workflowHandler.ProxyToWorkflowClient) // 启动工作流
			}

			// 任务管理
			tasks := workflow.Group("/tasks")
			{
				tasks.GET("", workflowHandler.ProxyToWorkflowClient)          // 获取任务列表
				tasks.POST("/approve", workflowHandler.ProxyToWorkflowClient) // 审批任务
			}
		}

		kb := api.Group("/kb")
		kb.Use(gw_middleware.AuthRequired(iamHandler))
		{
			kb.POST("/upload", kbHandler.ProxyToKbClient)
			kb.GET("/:id/preview", kbHandler.ProxyToKbClient)
			kb.GET("/:id/info", kbHandler.ProxyToKbClient)
			kb.GET("/:id/space", kbHandler.ProxyToKbClient)
			kb.GET("/homepage", kbHandler.ProxyToKbClient)
			kb.DELETE("/:id", kbHandler.ProxyToKbClient)
			kb.POST("/search", kbHandler.ProxyToKbClient)

			kb.POST("/:id/publish", kbHandler.ProxyToKbClient)
			kb.POST("/:id/unpublish", kbHandler.ProxyToKbClient)

			// 文档对话
			kb.POST("/:id/chat", kbHandler.ProxyToKbClient)
			kb.POST("/:id/chat/stream", kbHandler.ProxyToKbClient)
		}
	}

	return r
}
