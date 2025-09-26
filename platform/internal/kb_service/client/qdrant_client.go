package client

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/config"

	"github.com/qdrant/go-client/qdrant"
)

type QdrantClient struct {
	config *config.QdrantConfig
	client *qdrant.Client
	once   sync.Once
	err    error
}

func NewQdrantClient(config *config.QdrantConfig) *QdrantClient {
	return &QdrantClient{config: config}
}

// GetClient 获取Qdrant客户端，使用单例模式避免重复创建
func (c *QdrantClient) GetClient() (*qdrant.Client, error) {
	c.once.Do(func() {
		if err := c.validateConfig(); err != nil {
			c.err = fmt.Errorf("config validation failed: %w", err)
			return
		}

		port, err := strconv.Atoi(c.config.Port)
		if err != nil {
			c.err = fmt.Errorf("invalid Qdrant port %q: %w", c.config.Port, err)
			return
		}

		// 创建客户端
		c.client, c.err = qdrant.NewClient(&qdrant.Config{
			Host: c.config.Host,
			Port: port,
		})
		if c.err != nil {
			c.err = fmt.Errorf("failed to create Qdrant client: %w", c.err)
		}
	})
	return c.client, c.err
}

// validateConfig 验证配置参数
func (c *QdrantClient) validateConfig() error {
	if c.config.Host == "" {
		return fmt.Errorf("Qdrant host is required")
	}
	if c.config.Port == "" {
		return fmt.Errorf("Qdrant port is required")
	}
	return nil
}

// CreateCollection 创建集合
func (c *QdrantClient) CreateCollection(ctx context.Context, collectionName string, vectorSize uint64) error {
	client, err := c.GetClient()
	if err != nil {
		return err
	}
	if collectionName == "" {
		return fmt.Errorf("collection name is required")
	}
	if vectorSize == 0 {
		return fmt.Errorf("vector size must be greater than zero")
	}

	request := &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     vectorSize,
			Distance: qdrant.Distance_Cosine,
		}),
	}

	if err := client.CreateCollection(ctx, request); err != nil {
		return fmt.Errorf("failed to create collection %s: %w", collectionName, err)
	}

	return nil
}

// CollectionExists 检查集合是否存在
func (c *QdrantClient) CollectionExists(ctx context.Context, collectionName string) (bool, error) {
	client, err := c.GetClient()
	if err != nil {
		return false, err
	}
	if collectionName == "" {
		return false, fmt.Errorf("collection name is required")
	}

	exists, err := client.CollectionExists(ctx, collectionName)
	if err != nil {
		return false, fmt.Errorf("failed to check collection %s: %w", collectionName, err)
	}
	return exists, nil
}

// UpsertPoints 插入或更新向量点
func (c *QdrantClient) UpsertPoints(ctx context.Context, collectionName string, points interface{}) error {
	client, err := c.GetClient()
	if err != nil {
		return err
	}
	if collectionName == "" {
		return fmt.Errorf("collection name is required")
	}

	pointStructs, err := normalizePointStructs(points)
	if err != nil {
		return err
	}
	if len(pointStructs) == 0 {
		return fmt.Errorf("no points provided")
	}

	result, err := client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points:         pointStructs,
		Wait:           qdrant.PtrOf(true),
	})
	if err != nil {
		return fmt.Errorf("failed to upsert points into %s: %w", collectionName, err)
	}
	if status := result.GetStatus(); status != qdrant.UpdateStatus_Completed && status != qdrant.UpdateStatus_Acknowledged {
		return fmt.Errorf("upsert operation not completed: status=%v", status)
	}

	return nil
}

// SearchPoints 搜索向量点
func (c *QdrantClient) SearchPoints(ctx context.Context, collectionName string, vector []float32, limit uint64, scoreThreshold float32) ([]interface{}, error) {
	client, err := c.GetClient()
	if err != nil {
		return nil, err
	}
	if collectionName == "" {
		return nil, fmt.Errorf("collection name is required")
	}
	if len(vector) == 0 {
		return nil, fmt.Errorf("search vector is required")
	}
	if limit == 0 {
		limit = 10
	}

	req := &qdrant.SearchPoints{
		CollectionName: collectionName,
		Vector:         vector,
		Limit:          limit,
	}
	if scoreThreshold > 0 {
		req.ScoreThreshold = qdrant.PtrOf(scoreThreshold)
	}

	resp, err := client.GetPointsClient().Search(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to search points in %s: %w", collectionName, err)
	}

	result := resp.GetResult()
	points := make([]interface{}, len(result))
	for i, p := range result {
		points[i] = p
	}

	return points, nil
}

// DeletePoints 删除向量点
func (c *QdrantClient) DeletePoints(ctx context.Context, collectionName string, pointIDs []string) error {
	client, err := c.GetClient()
	if err != nil {
		return err
	}
	if collectionName == "" {
		return fmt.Errorf("collection name is required")
	}
	if len(pointIDs) == 0 {
		return fmt.Errorf("at least one point ID is required")
	}

	ids := make([]*qdrant.PointId, len(pointIDs))
	for i, id := range pointIDs {
		pointID, err := parsePointID(id)
		if err != nil {
			return fmt.Errorf("invalid point id %q: %w", id, err)
		}
		ids[i] = pointID
	}

	result, err := client.Delete(ctx, &qdrant.DeletePoints{
		CollectionName: collectionName,
		Points: &qdrant.PointsSelector{
			PointsSelectorOneOf: &qdrant.PointsSelector_Points{
				Points: &qdrant.PointsIdsList{Ids: ids},
			},
		},
		Wait: qdrant.PtrOf(true),
	})
	if err != nil {
		return fmt.Errorf("failed to delete points from %s: %w", collectionName, err)
	}
	if status := result.GetStatus(); status != qdrant.UpdateStatus_Completed && status != qdrant.UpdateStatus_Acknowledged {
		return fmt.Errorf("delete operation not completed: status=%v", status)
	}

	return nil
}

// GetCollectionInfo 获取集合信息
func (c *QdrantClient) GetCollectionInfo(ctx context.Context, collectionName string) (interface{}, error) {
	client, err := c.GetClient()
	if err != nil {
		return nil, err
	}
	if collectionName == "" {
		return nil, fmt.Errorf("collection name is required")
	}

	info, err := client.GetCollectionInfo(ctx, collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection info for %s: %w", collectionName, err)
	}
	return info, nil
}

// DeleteCollection 删除集合
func (c *QdrantClient) DeleteCollection(ctx context.Context, collectionName string) error {
	client, err := c.GetClient()
	if err != nil {
		return err
	}
	if collectionName == "" {
		return fmt.Errorf("collection name is required")
	}

	if err := client.DeleteCollection(ctx, collectionName); err != nil {
		return fmt.Errorf("failed to delete collection %s: %w", collectionName, err)
	}
	return nil
}

// ListCollections 列出所有集合
func (c *QdrantClient) ListCollections(ctx context.Context) ([]string, error) {
	client, err := c.GetClient()
	if err != nil {
		return nil, err
	}

	collections, err := client.ListCollections(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}
	return collections, nil
}

// normalizePointStructs 规范化点结构体
func normalizePointStructs(points interface{}) ([]*qdrant.PointStruct, error) {
	switch v := points.(type) {
	case nil:
		return nil, fmt.Errorf("points cannot be nil")
	case *qdrant.PointStruct:
		return []*qdrant.PointStruct{v}, nil
	case []*qdrant.PointStruct:
		return v, nil
	case []qdrant.PointStruct:
		if len(v) == 0 {
			return []*qdrant.PointStruct{}, nil
		}
		result := make([]*qdrant.PointStruct, len(v))
		for i := range v {
			result[i] = &v[i]
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported points type %T", points)
	}
}

// parsePointID 解析点ID
func parsePointID(id string) (*qdrant.PointId, error) {
	if id == "" {
		return nil, fmt.Errorf("point id is empty")
	}
	if num, err := strconv.ParseUint(id, 10, 64); err == nil {
		return &qdrant.PointId{PointIdOptions: &qdrant.PointId_Num{Num: num}}, nil
	}
	return &qdrant.PointId{PointIdOptions: &qdrant.PointId_Uuid{Uuid: id}}, nil
}
