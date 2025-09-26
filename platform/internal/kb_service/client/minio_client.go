package client

import (
	"context"
	"fmt"
	"io"
	"sync"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
)

// PathPrefix 路径前缀 固定变量
const (
	PathPrefixTemp      = "/temp/"
	PathPrefixPermanent = "/permanent/"
)

type S3Client struct {
	config *config.MinioConfig
	client *minio.Client
	once   sync.Once
	err    error
}

func NewS3Client(config *config.MinioConfig) *S3Client {
	return &S3Client{config: config}
}

// GetClient 获取MinIO客户端，使用单例模式避免重复创建
func (c *S3Client) GetClient() (*minio.Client, error) {
	c.once.Do(func() {
		if err := c.validateConfig(); err != nil {
			c.err = fmt.Errorf("config validation failed: %w", err)
			return
		}

		c.client, c.err = minio.New(c.config.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(c.config.AccessKey, c.config.SecretKey, ""),
			Secure: c.isSecureEndpoint(),
		})
		if c.err != nil {
			c.err = fmt.Errorf("failed to create MinIO client: %w", c.err)
		}
	})
	return c.client, c.err
}

// validateConfig 验证配置参数
func (c *S3Client) validateConfig() error {
	if c.config.Endpoint == "" {
		return fmt.Errorf("MinIO endpoint is required")
	}
	if c.config.AccessKey == "" {
		return fmt.Errorf("MinIO access key is required")
	}
	if c.config.SecretKey == "" {
		return fmt.Errorf("MinIO secret key is required")
	}
	return nil
}

// isSecureEndpoint 判断是否使用HTTPS
func (c *S3Client) isSecureEndpoint() bool {
	return c.config.Endpoint[:8] == "https://"
}

// UploadFile 上传文件到MinIO
func (c *S3Client) UploadFile(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) error {
	client, err := c.GetClient()
	if err != nil {
		return err
	}

	// 确保bucket存在
	if err := c.ensureBucketExists(ctx, client, c.config.Bucket); err != nil {
		return fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	// 上传文件
	_, err = client.PutObject(ctx, c.config.Bucket, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

// DownloadFile 从MinIO下载文件
func (c *S3Client) DownloadFile(ctx context.Context, objectName string) (io.Reader, error) {
	client, err := c.GetClient()
	if err != nil {
		return nil, err
	}

	object, err := client.GetObject(ctx, c.config.Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return object, nil
}

// DeleteFile 从MinIO删除文件
func (c *S3Client) DeleteFile(ctx context.Context, objectName string) error {
	client, err := c.GetClient()
	if err != nil {
		return err
	}

	err = client.RemoveObject(ctx, c.config.Bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// ensureBucketExists 确保bucket存在
func (c *S3Client) ensureBucketExists(ctx context.Context, client *minio.Client, bucketName string) error {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = c.createBucket(ctx, client)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *S3Client) createBucket(ctx context.Context, client *minio.Client) error {
	err := client.MakeBucket(ctx, c.config.Bucket, minio.MakeBucketOptions{
		Region: c.config.Region,
	})
	if err != nil {
		return err
	}

	// temp路径配置1天 生命周期
	lifecycleConfig := &lifecycle.Configuration{
		Rules: []lifecycle.Rule{
			{
				ID:     "temp-cleanup",
				Status: "Enabled",
				Prefix: "temp/",
				Expiration: lifecycle.Expiration{
					Days: 1,
				},
			},
		},
	}

	err = client.SetBucketLifecycle(ctx, c.config.Bucket, lifecycleConfig)
	if err != nil {
		return err
	}

	return nil
}

// GetBucketName 获取配置的桶名
func (c *S3Client) GetBucketName() string {
	return c.config.Bucket
}
