package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	model "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/models"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/configs"
)

type IamClient struct {
	config *configs.IamConfig
	client *http.Client
}

func NewIamClient(config *configs.IamConfig) *IamClient {
	return &IamClient{
		config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *IamClient) ValidateToken(token string) (*model.User, error) {
	req, err := http.NewRequest("POST", c.config.Url+"/api/v1/auth/validate-token", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", token)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("无效的token: " + resp.Status)
	}

	var response struct {
		Message string     `json:"message"`
		Data    model.User `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response.Data, nil
}
