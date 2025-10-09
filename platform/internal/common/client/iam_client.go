package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	commonConfig "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/config"
	model "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/models"
)

type IamClient struct {
	config *commonConfig.IamConfig
	client *http.Client
}

func NewIamClient(config *commonConfig.IamConfig) *IamClient {
	return &IamClient{
		config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *IamClient) GetSpaceMemebersByRole(user *model.User, spaceID uint, role string) ([]model.User, error) {
	targetURL := fmt.Sprintf("%s/api/v1/spaces/%d/members/role/%s", c.config.Url, spaceID, role)
	resp, err := c.client.Get(targetURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("获取空间成员失败: " + resp.Status)
	}

	var response struct {
		Code    int          `json:"code"`
		Message string       `json:"message"`
		Data    []model.User `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.New("获取空间成员失败: " + err.Error())
	}
	return response.Data, nil
}
