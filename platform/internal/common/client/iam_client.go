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

	// 创建请求并设置 User-ID header
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, errors.New("创建请求失败: " + err.Error())
	}
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", user.ID))

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("获取空间成员失败: " + resp.Status)
	}

	var response struct {
		Code    int                 `json:"code"`
		Message string              `json:"message"`
		Data    []model.SpaceMember `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.New("获取空间成员失败: " + err.Error())
	}
	var users []model.User
	for _, member := range response.Data {
		users = append(users, member.User)
	}
	return users, nil
}
