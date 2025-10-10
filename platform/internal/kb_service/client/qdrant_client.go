package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/config"
)

// QdrantClient 访问 Qdrant 向量数据库的客户端
type QdrantClient struct {
	baseURL    string
	apiKey     string
	collection string
	vectorSize int
	distance   string
	httpClient *http.Client
}

// QdrantPoint 表示向量点
type QdrantPoint struct {
	ID      string         `json:"id"`
	Vector  []float64      `json:"vector"`
	Payload map[string]any `json:"payload,omitempty"`
}

type qdrantUpsertRequest struct {
	Points []QdrantPoint `json:"points"`
}

type qdrantCollectionRequest struct {
	Vectors struct {
		Size     int    `json:"size"`
		Distance string `json:"distance"`
		OnDisk   bool   `json:"on_disk"`
	} `json:"vectors"`
}

// NewQdrantClient 创建新的 Qdrant 客户端
func NewQdrantClient(cfg *config.QdrantConfig) *QdrantClient {
	if cfg == nil || strings.TrimSpace(cfg.BaseURL) == "" {
		return nil
	}

	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	vectorSize := cfg.VectorSize
	if vectorSize <= 0 {
		vectorSize = 7
	}

	distance := cfg.Distance
	if distance == "" {
		distance = "Cosine"
	}

	return &QdrantClient{
		baseURL:    strings.TrimRight(cfg.BaseURL, "/"),
		apiKey:     cfg.APIKey,
		collection: cfg.Collection,
		vectorSize: vectorSize,
		distance:   distance,
		httpClient: &http.Client{Timeout: timeout},
	}
}

// VectorSize 返回当前向量维度
func (c *QdrantClient) VectorSize() int {
	if c == nil {
		return 0
	}
	return c.vectorSize
}

// EnsureCollection 确保集合存在
func (c *QdrantClient) EnsureCollection(ctx context.Context) error {
	if c == nil {
		return errors.New("qdrant client is not configured")
	}

	reqBody := qdrantCollectionRequest{}
	reqBody.Vectors.Size = c.vectorSize
	reqBody.Vectors.Distance = c.distance
	reqBody.Vectors.OnDisk = false

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("qdrant: marshal collection request failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.baseURL+"/collections/"+c.collection, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("qdrant: create collection request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("qdrant: ensure collection request failed: %w", err)
	}
	defer resp.Body.Close()

	// 200、201、409 都视为成功
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
		return fmt.Errorf("qdrant: ensure collection unexpected status %d", resp.StatusCode)
	}

	return nil
}

// UpsertPoints 批量写入向量
func (c *QdrantClient) UpsertPoints(ctx context.Context, points []QdrantPoint) error {
	if c == nil {
		return errors.New("qdrant client is not configured")
	}

	if len(points) == 0 {
		return nil
	}

	for _, p := range points {
		if len(p.Vector) != c.vectorSize {
			return fmt.Errorf("qdrant: vector size mismatch, expect %d got %d", c.vectorSize, len(p.Vector))
		}
	}

	body, err := json.Marshal(qdrantUpsertRequest{Points: points})
	if err != nil {
		return fmt.Errorf("qdrant: marshal upsert request failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.baseURL+"/collections/"+c.collection+"/points", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("qdrant: create upsert request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("qdrant: upsert request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("qdrant: upsert unexpected status %d", resp.StatusCode)
	}

	return nil
}
