package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/config"
)

// PaddleOCRClient 调用 PaddleOCR 服务的客户端
type PaddleOCRClient struct {
	baseURL    string
	language   string
	httpClient *http.Client
}

type paddleOCRRequest struct {
	FileName string `json:"file_name"`
	Content  string `json:"content_base64"`
	Language string `json:"language"`
}

type paddleOCRResponse struct {
	Text  string   `json:"text"`
	Lines []string `json:"lines,omitempty"`
	Error string   `json:"error,omitempty"`
}

// NewPaddleOCRClient 创建 OCR 客户端
func NewPaddleOCRClient(cfg *config.PaddleOCRConfig) *PaddleOCRClient {
	if cfg == nil || strings.TrimSpace(cfg.BaseURL) == "" {
		return nil
	}

	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	client := &http.Client{
		Timeout: timeout,
	}

	lang := strings.TrimSpace(cfg.Language)
	if lang == "" {
		lang = "ch"
	}

	return &PaddleOCRClient{
		baseURL:    strings.TrimRight(cfg.BaseURL, "/"),
		language:   lang,
		httpClient: client,
	}
}

// Recognize 调用 OCR 服务识别文本
func (c *PaddleOCRClient) Recognize(ctx context.Context, fileName string, data []byte) (string, error) {
	if c == nil {
		return "", errors.New("paddle OCR client is not configured")
	}

	if len(data) == 0 {
		return "", errors.New("paddle OCR: empty payload")
	}

	payload := paddleOCRRequest{
		FileName: fileName,
		Content:  base64.StdEncoding.EncodeToString(data),
		Language: c.language,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("paddle OCR: failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/ocr", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("paddle OCR: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("paddle OCR: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("paddle OCR: unexpected status code %d", resp.StatusCode)
	}

	var ocrResp paddleOCRResponse
	if err := json.NewDecoder(resp.Body).Decode(&ocrResp); err != nil {
		return "", fmt.Errorf("paddle OCR: failed to decode response: %w", err)
	}

	if ocrResp.Error != "" {
		return "", fmt.Errorf("paddle OCR: service error %s", ocrResp.Error)
	}

	if ocrResp.Text != "" {
		return ocrResp.Text, nil
	}

	if len(ocrResp.Lines) > 0 {
		return strings.Join(ocrResp.Lines, "\n"), nil
	}

	return "", errors.New("paddle OCR: empty response")
}
