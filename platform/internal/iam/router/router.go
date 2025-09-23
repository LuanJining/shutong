package router

import (
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam/config"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam/handler"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(cfg *config.Config, db *gorm.DB) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// 添加中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API路由组
	api := r.Group("/api/v1")
	{
		// 认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/login", handler.Login)
			auth.POST("/logout", handler.Logout)
			auth.POST("/refresh", handler.RefreshToken)
		}

		// 用户管理路由
		users := api.Group("/users")
		users.Use(middleware.AuthRequired())
		{
			users.GET("", handler.GetUsers)
			users.GET("/:id", handler.GetUser)
			users.POST("", handler.CreateUser)
			users.PUT("/:id", handler.UpdateUser)
			users.DELETE("/:id", handler.DeleteUser)
		}

		// 角色管理路由
		roles := api.Group("/roles")
		roles.Use(middleware.AuthRequired())
		{
			roles.GET("", handler.GetRoles)
			roles.GET("/:id", handler.GetRole)
			roles.POST("", handler.CreateRole)
			roles.PUT("/:id", handler.UpdateRole)
			roles.DELETE("/:id", handler.DeleteRole)
		}

		// 权限管理路由
		permissions := api.Group("/permissions")
		permissions.Use(middleware.AuthRequired())
		{
			permissions.GET("", handler.GetPermissions)
			permissions.GET("/:id", handler.GetPermission)
		}
	}

	return r
}
