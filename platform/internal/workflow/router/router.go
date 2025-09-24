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
	r.Use(middleware.CORS())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	workflowService := service.NewWorkflowService(db, &cfg.Workflow)
	handler := handler.NewHandler(db, workflowService)

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		workflow := api.Group("/workflow")
		{
			workflow.GET("/", handler.GetWorkflow)
		}
	}

	return r
}
