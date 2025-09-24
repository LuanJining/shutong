package router

import (
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/config"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/handler"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.Gin.Mode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	handler := handler.NewHandler(&cfg.Iam)
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
				auth.POST("/login", handler.ProxyToIamClient)
				auth.POST("/logout", handler.ProxyToIamClient)
				auth.POST("/refresh", handler.ProxyToIamClient)
				auth.PATCH("/change-password", handler.ProxyToIamClient)
			}

			// 用户管理路由
			users := iam.Group("/users")
			{
				users.GET("", handler.ProxyToIamClient)
				users.GET("/:id", handler.ProxyToIamClient)
				users.POST("", handler.ProxyToIamClient)
				users.PUT("/:id", handler.ProxyToIamClient)
				users.DELETE("/:id", handler.ProxyToIamClient)
			}

			// 角色管理路由
			roles := iam.Group("/roles")
			{
				roles.GET("", handler.ProxyToIamClient)
				roles.GET("/:id", handler.ProxyToIamClient)
				roles.GET("/:id/permissions", handler.ProxyToIamClient)
				roles.POST("", handler.ProxyToIamClient)
				roles.PUT("/:id", handler.ProxyToIamClient)
				roles.DELETE("/:id", handler.ProxyToIamClient)
			}

			// 权限管理路由
			permissions := iam.Group("/permissions")
			{
				permissions.GET("", handler.ProxyToIamClient)
				permissions.GET("/:id", handler.ProxyToIamClient)
				permissions.POST("/check", handler.ProxyToIamClient)
			}

			// 空间管理路由
			spaces := iam.Group("/spaces")
			{
				spaces.GET("", handler.ProxyToIamClient)
				spaces.GET("/:id", handler.ProxyToIamClient)
				spaces.POST("", handler.ProxyToIamClient)
				spaces.PUT("/:id", handler.ProxyToIamClient)
				spaces.DELETE("/:id", handler.ProxyToIamClient)
				spaces.GET("/:id/members", handler.ProxyToIamClient)
				spaces.POST("/:id/members", handler.ProxyToIamClient)
				spaces.DELETE("/:id/members/:user_id", handler.ProxyToIamClient)
			}
		}
	}

	return r
}
