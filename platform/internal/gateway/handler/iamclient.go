package handler

import (
	"fmt"
	"io"
	"net/http"

	"strings"
	"time"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/config"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/models"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	config *config.IamConfig
}

func NewHandler(config *config.IamConfig) *Handler {
	return &Handler{config: config}
}

func (h *Handler) ProxyToIamClient(c *gin.Context) {
	iamClient := h.config.Url
	path := c.Request.URL.Path

	var targetURL string

	if suffix, ok := strings.CutPrefix(path, "/api/v1/iam/"); ok {
		targetURL = iamClient + "/api/v1/" + suffix
	} else {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_PATH",
				Message: "无效的请求路径",
				Details: fmt.Sprintf("路径 %s 不匹配 IAM 代理模式", path),
			},
		})
		return
	}

	// 验证目标URL格式
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "http://" + targetURL
	}

	client := &http.Client{
		Timeout: 30 * time.Second, // 减少超时时间到30秒
	}

	fullURL := targetURL
	if c.Request.URL.RawQuery != "" {
		fullURL += "?" + c.Request.URL.RawQuery
	}

	// 创建请求
	req, err := http.NewRequest(c.Request.Method, fullURL, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "PROXY_ERROR",
				Message: "创建代理请求失败",
				Details: err.Error(),
			},
		})
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
		c.JSON(http.StatusBadGateway, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "PROXY_ERROR",
				Message: "代理请求失败",
				Details: err.Error(),
			},
		})
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// 设置状态码并复制响应体
	c.Status(resp.StatusCode)

	// 对于流式响应，直接复制流
	if resp.Header.Get("Content-Type") == "text/event-stream" {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		// 流式复制响应
		buffer := make([]byte, 1024)
		for {
			n, err := resp.Body.Read(buffer)
			if n > 0 {
				c.Writer.Write(buffer[:n])
				c.Writer.Flush()
			}
			if err != nil {
				if err == io.EOF {
					// 正常结束，发送完成信号
					c.Writer.WriteString("data: [DONE]\n\n")
					c.Writer.Flush()
				} else {
					// 其他错误，记录日志
					fmt.Printf("流式响应读取错误: %v\n", err)
				}
				break
			}
		}
	} else {
		// 普通响应，直接复制
		io.Copy(c.Writer, resp.Body)
	}
}
