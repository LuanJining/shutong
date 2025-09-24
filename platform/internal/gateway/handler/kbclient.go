package handler

import (
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/config"
	"github.com/gin-gonic/gin"
)

type KbHandler struct {
	config *config.KbConfig
}

func NewKbHandler(config *config.KbConfig) *KbHandler {
	return &KbHandler{config: config}
}

func (h *KbHandler) ProxyToKbClient(c *gin.Context) {
	c.JSON(500, gin.H{"message": "TODO"})
}
