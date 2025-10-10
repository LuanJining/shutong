package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	model "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/models"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/configs"
	"github.com/gin-gonic/gin"
)

type KbHandler struct {
	config *configs.KbConfig
	client *http.Client
}

func NewKbHandler(config *configs.KbConfig) *KbHandler {
	return &KbHandler{
		config: config,
		client: &http.Client{
			Timeout: 5 * time.Minute, // 设置5分钟超时，适合大文件上传
		},
	}
}

func (h *KbHandler) ProxyToKbClient(c *gin.Context) {
	// 构建目标URL
	targetURL := h.buildTargetURL(c)

	// 读取请求体
	var body io.Reader
	if c.Request.Body != nil {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
			return
		}
		body = bytes.NewReader(bodyBytes)
	}

	// 创建代理请求
	proxyReq, err := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, targetURL, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request"})
		return
	}

	// 复制请求头
	h.copyHeaders(c.Request.Header, proxyReq.Header)

	// 添加用户ID到请求头（从 Gateway 的 AuthRequired 中间件获取）
	user, _ := c.Get("user")

	if user != nil {
		proxyReq.Header.Add("X-User-ID", fmt.Sprintf("%d", user.(*model.User).ID))
	}

	// 发送请求
	resp, err := h.client.Do(proxyReq)
	if err != nil {
		fmt.Printf("Gateway proxy error: %v\n", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to proxy request: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	h.copyHeaders(resp.Header, c.Writer.Header())

	// 设置状态码
	c.Status(resp.StatusCode)

	// 复制响应体
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		fmt.Printf("Gateway response copy error: %v\n", err)
		// 不要返回错误，因为可能已经部分传输了
	}
}

func (h *KbHandler) buildTargetURL(c *gin.Context) string {
	// 从请求路径中提取kb相关的路径
	path := c.Request.URL.Path
	// 移除 /api/v1/kb 前缀
	kbPath := strings.TrimPrefix(path, "/api/v1/kb")

	// 构建目标URL
	targetURL := strings.TrimSuffix(h.config.Url, "/") + "/api/v1/documents" + kbPath

	// 添加查询参数
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	return targetURL
}

func (h *KbHandler) copyHeaders(src, dst http.Header) {
	for key, values := range src {
		// 跳过一些不需要转发的头
		if key == "Host" || key == "Content-Length" {
			continue
		}
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}
