package handler

import (
	"net/http"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/workflow/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler 处理器结构体
type Handler struct {
	db              *gorm.DB
	workflowService *service.WorkflowService
}

// NewHandler 创建新的处理器
func NewHandler(db *gorm.DB, workflowService *service.WorkflowService) *Handler {
	return &Handler{
		db:              db,
		workflowService: workflowService,
	}
}

func (h *Handler) GetWorkflow(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}
