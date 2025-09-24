package handler

import (
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/config"
	"github.com/gin-gonic/gin"
)

type WorkflowHandler struct {
	config *config.WorkflowConfig
}

func NewWorkflowHandler(config *config.WorkflowConfig) *WorkflowHandler {
	return &WorkflowHandler{config: config}
}

func (h *WorkflowHandler) ProxyToWorkflowClient(c *gin.Context) {
	c.JSON(500, gin.H{"message": "TODO"})
}
