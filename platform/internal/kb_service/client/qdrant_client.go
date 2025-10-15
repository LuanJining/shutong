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

type QdrantMatch struct {
	Value any `json:"value"`
}

type QdrantCondition struct {
	Key   string      `json:"key"`
	Match QdrantMatch `json:"match"`
}

type QdrantFilter struct {
	Must []QdrantCondition `json:"must,omitempty"`
}

type qdrantSearchRequest struct {
	Vector      []float64     `json:"vector"`
	Top         int           `json:"top"`
	WithPayload bool          `json:"with_payload"`
	WithVector  bool          `json:"with_vector"`
	Filter      *QdrantFilter `json:"filter,omitempty"`
}

type QdrantSearchResult struct {
	ID      string         `json:"id"`
	Score   float64        `json:"score"`
	Payload map[string]any `json:"payload"`
}

type qdrantSearchResponse struct {
	Result []QdrantSearchResult `json:"result"`
	Status string               `json:"status"`
	Time   float64              `json:"time"`
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
		vectorSize = 1024 // 默认使用Jina embeddings-v3 向量维度
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

// SearchPoints 在 Qdrant 中检索相似向量
func (c *QdrantClient) SearchPoints(ctx context.Context, vector []float64, top int, filter *QdrantFilter) ([]QdrantSearchResult, error) {
	if c == nil {
		return nil, errors.New("qdrant client is not configured")
	}
	if len(vector) != c.vectorSize {
		return nil, fmt.Errorf("qdrant: vector size mismatch, expect %d got %d", c.vectorSize, len(vector))
	}
	if top <= 0 {
		top = 5
	}

	body, err := json.Marshal(qdrantSearchRequest{
		Vector:      vector,
		Top:         top,
		WithPayload: true,
		WithVector:  false,
		Filter:      filter,
	})
	if err != nil {
		return nil, fmt.Errorf("qdrant: marshal search request failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/collections/"+c.collection+"/points/search", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("qdrant: create search request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("qdrant: search request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("qdrant: search unexpected status %d", resp.StatusCode)
	}

	var searchResp qdrantSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("qdrant: failed to decode search response: %w", err)
	}

	return searchResp.Result, nil
}
