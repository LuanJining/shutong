package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/configs"
	"github.com/gin-gonic/gin"
)

type WorkflowHandler struct {
	config *configs.WorkflowConfig
}

func NewWorkflowHandler(config *configs.WorkflowConfig) *WorkflowHandler {
	return &WorkflowHandler{config: config}
}

func (h *WorkflowHandler) ProxyToWorkflowClient(c *gin.Context) {
	// 构建目标URL
	targetURL := h.config.Url + "/api/v1" + strings.TrimPrefix(c.Request.URL.Path, "/api/v1")

	// 添加查询参数
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	// 解析目标URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target URL"})
		return
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 创建请求
	req, err := http.NewRequest(c.Request.Method, parsedURL.String(), c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// 复制请求头
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to connect to workflow service"})
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// 设置状态码
	c.Status(resp.StatusCode)

	// 复制响应体
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		fmt.Printf("Error copying response body: %v\n", err)
	}
}
