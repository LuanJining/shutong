package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/config"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/model"
)

// TODO 存取流程打通后接入
type WorkflowClient struct {
	config *config.WorkflowConfig
	client *http.Client
}

func NewWorkflowClient(config *config.WorkflowConfig) *WorkflowClient {
	return &WorkflowClient{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *WorkflowClient) GetClient() (*http.Client, error) {
	return c.client, nil
}

// TODO直接模版化创建文档上传流程
func (c *WorkflowClient) CreateWorkflow(ctx context.Context, workflow *model.Workflow, userID uint) (string, error) {
	targetURL := fmt.Sprintf("%s/api/v1/workflow/workflows", c.config.Url)
	jsonData, err := json.Marshal(workflow)
	if err != nil {
		return "", fmt.Errorf("failed to marshal workflow: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", targetURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
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

	// 从workflow对象中提取ID
	if response.Data == nil {
		return "", fmt.Errorf("no data in response")
	}

	// 将Data转换为map以提取ID
	dataMap, ok := response.Data.(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid response data format")
	}

	idFloat, ok := dataMap["id"].(float64) // JSON数字默认解析为float64
	if !ok {
		return "", fmt.Errorf("invalid ID format in response")
	}

	// 将float64转换为string
	id := fmt.Sprintf("%.0f", idFloat)
	return id, nil
}

func (c *WorkflowClient) StartWorkflow(ctx context.Context, workflow *model.Workflow, userID uint, documentID uint) (string, error) {
	targetURL := fmt.Sprintf("%s/api/v1/workflow/instances", c.config.Url)

	workflowRequest := model.StartWorkflowRequest{
		WorkflowID:   workflow.ID,
		Title:        "文档上传审批",
		Description:  "文档上传审批",
		ResourceType: "document",
		ResourceID:   documentID,
		SpaceID:      workflow.SpaceID,
		Priority:     "normal",
	}
	jsonData, err := json.Marshal(workflowRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal workflow: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", targetURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
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

	// 从workflow对象中提取ID
	if response.Data == nil {
		return "", fmt.Errorf("no data in response")
	}

	// 将Data转换为map以提取ID
	dataMap, ok := response.Data.(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid response data format")
	}

	idFloat, ok := dataMap["id"].(float64) // JSON数字默认解析为float64
	if !ok {
		return "", fmt.Errorf("invalid ID format in response")
	}

	// 将float64转换为string
	id := fmt.Sprintf("%.0f", idFloat)
	return id, nil
}
