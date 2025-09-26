package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"sync"
	"time"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/config"
)

type OCRClient struct {
	config     *config.OCRConfig
	httpClient *http.Client
	once       sync.Once
	err        error
}

// OCRResponse OCR解析响应
type OCRResponse struct {
	Success bool   `json:"success"`
	Text    string `json:"text"`
	Error   string `json:"error,omitempty"`
}

// OCRRequest OCR解析请求
type OCRRequest struct {
	File     io.Reader
	FileName string
	Language string // 语言代码，如 "chi_sim+eng" 表示中文简体+英文
}

func NewOCRClient(config *config.OCRConfig) *OCRClient {
	return &OCRClient{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetClient 获取OCR客户端，使用单例模式避免重复创建
func (c *OCRClient) GetClient() (*http.Client, error) {
	c.once.Do(func() {
		if err := c.validateConfig(); err != nil {
			c.err = fmt.Errorf("config validation failed: %w", err)
			return
		}
	})
	return c.httpClient, c.err
}

// validateConfig 验证配置参数
func (c *OCRClient) validateConfig() error {
	if c.config.Url == "" {
		return fmt.Errorf("OCR service URL is required")
	}
	return nil
}

// ParseImage 解析图片文件
func (c *OCRClient) ParseImage(ctx context.Context, req *OCRRequest) (*OCRResponse, error) {
	client, err := c.GetClient()
	if err != nil {
		return nil, err
	}

	// 创建multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加文件字段
	fileWriter, err := writer.CreateFormFile("file", req.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(fileWriter, req.File); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// 添加语言参数
	if req.Language != "" {
		writer.WriteField("language", req.Language)
	}

	writer.Close()

	// 创建请求
	url := fmt.Sprintf("%s/api/v1/ocr/parse/image", c.config.Url)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var ocrResp OCRResponse
	if err := json.NewDecoder(resp.Body).Decode(&ocrResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !ocrResp.Success {
		return nil, fmt.Errorf("OCR parsing failed: %s", ocrResp.Error)
	}

	return &ocrResp, nil
}

// ParsePDF 解析PDF文件
func (c *OCRClient) ParsePDF(ctx context.Context, req *OCRRequest) (*OCRResponse, error) {
	client, err := c.GetClient()
	if err != nil {
		return nil, err
	}

	// 创建multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加文件字段
	fileWriter, err := writer.CreateFormFile("file", req.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(fileWriter, req.File); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// 添加语言参数
	if req.Language != "" {
		writer.WriteField("language", req.Language)
	}

	writer.Close()

	// 创建请求
	url := fmt.Sprintf("%s/api/v1/ocr/parse/pdf", c.config.Url)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var ocrResp OCRResponse
	if err := json.NewDecoder(resp.Body).Decode(&ocrResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !ocrResp.Success {
		return nil, fmt.Errorf("OCR parsing failed: %s", ocrResp.Error)
	}

	return &ocrResp, nil
}

// ParseWord 解析Word文档
func (c *OCRClient) ParseWord(ctx context.Context, req *OCRRequest) (*OCRResponse, error) {
	client, err := c.GetClient()
	if err != nil {
		return nil, err
	}

	// 创建multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加文件字段
	fileWriter, err := writer.CreateFormFile("file", req.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(fileWriter, req.File); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// 添加语言参数
	if req.Language != "" {
		writer.WriteField("language", req.Language)
	}

	writer.Close()

	// 创建请求
	url := fmt.Sprintf("%s/api/v1/ocr/parse/word", c.config.Url)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var ocrResp OCRResponse
	if err := json.NewDecoder(resp.Body).Decode(&ocrResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !ocrResp.Success {
		return nil, fmt.Errorf("OCR parsing failed: %s", ocrResp.Error)
	}

	return &ocrResp, nil
}

// HealthCheck 检查OCR服务健康状态
func (c *OCRClient) HealthCheck(ctx context.Context) error {
	client, err := c.GetClient()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v1/ocr/health", c.config.Url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send health check request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("OCR service is not healthy, status: %d", resp.StatusCode)
	}

	return nil
}
