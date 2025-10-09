package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	model "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/models"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/config"
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

func (c *WorkflowClient) CreateWorkflow(ctx context.Context, workflow *model.Workflow, userID uint) (uint, error) {
	targetURL := fmt.Sprintf("%s/api/v1/workflow/workflows", c.config.Url)
	jsonData, err := json.Marshal(workflow)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal workflow: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", targetURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("workflow service returned status %d", resp.StatusCode)
	}

	// 解析响应
	var response model.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	// 检查响应码
	if response.Code != 200 {
		return 0, fmt.Errorf("workflow service error: %s", response.Message)
	}

	// 从workflow对象中提取ID
	if response.Data == nil {
		return 0, fmt.Errorf("no data in response")
	}

	// 将Data转换为Workflow对象以提取ID
	workflowData, ok := response.Data.(*model.Workflow)
	if !ok {
		// 尝试转换为map[string]interface{}再提取ID
		dataMap, ok := response.Data.(map[string]any)
		if !ok {
			return 0, fmt.Errorf("invalid response data format")
		}

		idFloat, ok := dataMap["id"].(float64)
		if !ok {
			return 0, fmt.Errorf("invalid ID format in response")
		}

		if idFloat == 0 {
			return 0, fmt.Errorf("invalid ID value in response")
		}

		return uint(idFloat), nil
	}

	if workflowData.ID == 0 {
		return 0, fmt.Errorf("invalid ID value in response")
	}

	return workflowData.ID, nil
}

func (c *WorkflowClient) CheckWorkflowStatus(ctx context.Context, workflowID uint) (string, error) {
	targetURL := fmt.Sprintf("%s/api/v1/workflow/workflows/%d/status", c.config.Url, workflowID)
	resp, err := c.client.Get(targetURL)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("workflow service returned status %d", resp.StatusCode)
	}

	// 解析响应
	var response model.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// 检查响应码
	if response.Code != 200 {
		return "", fmt.Errorf("workflow service error: %s", response.Message)
	}

	return response.Data.(string), nil
}

func (c *WorkflowClient) StartWorkflow(ctx context.Context, workflowID uint, userID uint) (*model.Workflow, error) {
	targetURL := fmt.Sprintf("%s/api/v1/workflow/workflows/%d/start", c.config.Url, workflowID)

	// 构造请求体
	reqBody := model.StartWorkflowRequest{
		WorkflowID: workflowID,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", targetURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))

	// 发送请求
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("workflow service returned status %d", resp.StatusCode)
	}

	// 解析响应
	var response model.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// 检查响应码
	if response.Code != 200 {
		return nil, fmt.Errorf("workflow service error: %s", response.Message)
	}

	// 从响应中提取Workflow对象
	if response.Data == nil {
		return nil, fmt.Errorf("no data in response")
	}

	// 尝试将Data转换为Workflow对象
	dataMap, ok := response.Data.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid response data format")
	}

	// 将map转换回JSON再解析为Workflow对象
	dataJSON, err := json.Marshal(dataMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var workflow model.Workflow
	err = json.Unmarshal(dataJSON, &workflow)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal workflow: %w", err)
	}

	return &workflow, nil
}
