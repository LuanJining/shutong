package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/config"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/model"
)

// TODO 存取流程打通后接入
type WorkflowClient struct {
	config *config.WorkflowConfig
	client *http.Client
}

func NewWorkflowClient(config *config.WorkflowConfig) *WorkflowClient {
	return &WorkflowClient{config: config}
}

func (c *WorkflowClient) GetClient() (*http.Client, error) {
	return c.client, nil
}

// TODO直接模版化创建文档上传流程
func (c *WorkflowClient) CreateWorkflow(ctx context.Context, workflow *model.Workflow) (string, error) {
	targetURL := fmt.Sprintf("%s/api/v1/workflows", c.config.Url)
	jsonData, err := json.Marshal(workflow)
	if err != nil {
		return "", fmt.Errorf("failed to marshal workflow: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", targetURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var response model.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	return response.Data.(map[string]any)["id"].(string), nil
}
