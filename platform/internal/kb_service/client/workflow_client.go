package client

import (
	"net/http"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/config"
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

func (c *WorkflowClient) validateConfig() error {
	return nil
}

// TODO直接模版化创建文档上传流程
