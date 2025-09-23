package router

import (
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam/config"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam/handler"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam/middleware"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	// 创建服务
	authService := service.NewAuthService(db, &cfg.JWT)

	// 创建处理器
	h := handler.NewHandler(db, authService)

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
			auth.POST("/login", h.Login)
			auth.POST("/logout", h.Logout)
			auth.POST("/refresh", h.RefreshToken)
			auth.PATCH("/change-password", middleware.AuthRequired(authService), h.ChangePassword)
		}

		// 用户管理路由
		users := api.Group("/users")
		users.Use(middleware.AuthRequired(authService))
		{
			users.GET("", h.GetUsers)
			users.GET("/:id", h.GetUser)
			users.PUT("/:id", h.UpdateUser)

			// 只有超级管理员才能创建和删除用户
			users.POST("", middleware.RequireRole([]string{"super_admin"}), h.CreateUser)
			users.DELETE("/:id", middleware.RequireRole([]string{"super_admin"}), h.DeleteUser)
		}

		// 角色管理路由
		roles := api.Group("/roles")
		roles.Use(middleware.AuthRequired(authService))
		{
			roles.GET("", h.GetRoles)
			roles.GET("/:id", h.GetRole)
			roles.POST("", h.CreateRole)
			roles.PUT("/:id", h.UpdateRole)
			roles.DELETE("/:id", h.DeleteRole)
		}

		// 权限管理路由
		permissions := api.Group("/permissions")
		permissions.Use(middleware.AuthRequired(authService))
		{
			permissions.GET("", h.GetPermissions)
			permissions.GET("/:id", h.GetPermission)
			permissions.POST("/check", h.CheckPermission)
		}
	}

	return r
}
