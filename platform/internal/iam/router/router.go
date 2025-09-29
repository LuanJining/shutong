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
	gin.SetMode(cfg.Gin.Mode)

	r := gin.New()

	// 添加中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	// r.Use(middleware.CORS())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	// 创建服务
	authService := service.NewAuthService(db, &cfg.JWT)

	// 创建处理器
	h := handler.NewHandler(db, authService)

	// API路由组
	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// 认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/login", h.Login)
			auth.POST("/logout", h.Logout)
			auth.POST("/refresh", h.RefreshToken)
			auth.PATCH("/change-password", middleware.AuthRequired(authService), h.ChangePassword)

			auth.POST("/validate-token", h.ValidateToken)
		}

		// 用户管理路由
		users := api.Group("/users")
		users.Use(middleware.AuthRequired(authService))
		{
			// 查看所有内容 - 所有认证用户都可以
			users.GET("", h.GetUsers)
			users.GET("/:id", h.GetUser)

			// 更新用户 - 先检查角色，再检查权限
			users.PUT("/:id", h.UpdateUser)

			// 添加/删除用户 - 只有超级管理员
			users.POST("", h.CreateUser)
			users.DELETE("/:id", h.DeleteUser)

			// 用户角色管理 - 先检查角色，再检查权限
			users.POST("/:id/roles", h.AssignUserRole)
			users.DELETE("/:id/roles/:role_id", h.RemoveUserRole)
		}

		// 角色管理路由
		roles := api.Group("/roles")
		roles.Use(middleware.AuthRequired(authService))
		{
			// 查看所有内容 - 所有认证用户都可以
			roles.GET("", h.GetRoles)
			roles.GET("/:id", h.GetRole)
			roles.GET("/:id/permissions", h.GetRolePermissions)

			// 角色管理 - 超级管理员、企业管理员
			roles.POST("", h.CreateRole)
			roles.PUT("/:id", h.UpdateRole)
			roles.DELETE("/:id", h.DeleteRole)

			// 角色权限管理 - 超级管理员、企业管理员
			roles.POST("/:id/permissions", h.AssignRolePermission)
			roles.DELETE("/:id/permissions/:permission_id", h.RemoveRolePermission)
		}

		// 权限管理路由
		permissions := api.Group("/permissions")
		permissions.Use(middleware.AuthRequired(authService))
		{
			permissions.GET("", h.GetPermissions)
			permissions.GET("/:id", h.GetPermission)
			permissions.POST("/check", h.CheckPermission)
		}

		// 空间管理路由
		spaces := api.Group("/spaces")
		spaces.Use(middleware.AuthRequired(authService))
		{
			// 查看所有内容 - 所有认证用户都可以
			spaces.GET("", h.GetSpaces)
			spaces.GET("/:id", h.GetSpace)

			// 创建知识空间 - 超级管理员、企业管理员、空间管理员
			spaces.POST("", h.CreateSpace)

			// 管理空间 - 先检查空间成员，再检查角色权限
			spaces.PUT("/:id", h.UpdateSpace)
			spaces.DELETE("/:id", h.DeleteSpace)

			// 空间成员管理 - 先检查空间成员，再检查角色权限
			spaces.GET("/:id/members", h.GetSpaceMembers)
			spaces.POST("/:id/members", h.AddSpaceMember)
			spaces.DELETE("/:id/members/:user_id", h.RemoveSpaceMember)
			spaces.PUT("/:id/members/:user_id", h.UpdateSpaceMemberRole)
			spaces.GET("/:id/members/:role_id", h.GetSpaceMembersByRole)
		}
	}

	return r
}
