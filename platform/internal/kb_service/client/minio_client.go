package client

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"sync"
	"time"

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
	if c.config.Bucket == "" {
		return fmt.Errorf("MinIO bucket name is required")
	}
	if err := c.validateBucketName(c.config.Bucket); err != nil {
		return fmt.Errorf("invalid bucket name: %w", err)
	}
	return nil
}

// validateBucketName 验证MinIO bucket名称是否符合规范
func (c *S3Client) validateBucketName(bucketName string) error {
	// MinIO bucket名称规则：
	// 1. 长度3-63字符
	// 2. 只能包含小写字母、数字、点、连字符
	// 3. 不能以点或连字符开头或结尾
	// 4. 不能包含连续的点
	// 5. 不能是IP地址格式

	if len(bucketName) < 3 || len(bucketName) > 63 {
		return fmt.Errorf("bucket name must be 3-63 characters long")
	}

	// 检查是否以点或连字符开头或结尾
	if bucketName[0] == '.' || bucketName[0] == '-' ||
		bucketName[len(bucketName)-1] == '.' || bucketName[len(bucketName)-1] == '-' {
		return fmt.Errorf("bucket name cannot start or end with dot or hyphen")
	}

	// 检查是否包含连续的点
	if regexp.MustCompile(`\.\.`).MatchString(bucketName) {
		return fmt.Errorf("bucket name cannot contain consecutive dots")
	}

	// 检查是否只包含允许的字符
	validPattern := regexp.MustCompile(`^[a-z0-9.-]+$`)
	if !validPattern.MatchString(bucketName) {
		return fmt.Errorf("bucket name can only contain lowercase letters, numbers, dots, and hyphens")
	}

	// 检查是否是IP地址格式
	ipPattern := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	if ipPattern.MatchString(bucketName) {
		return fmt.Errorf("bucket name cannot be an IP address")
	}

	return nil
}

// isSecureEndpoint 判断是否使用HTTPS
func (c *S3Client) isSecureEndpoint() bool {
	// 对于外部MinIO服务，默认使用HTTPS
	return true
}

// UploadFile 上传文件到MinIO，返回实际写入的字节数
func (c *S3Client) UploadFile(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (int64, error) {
	client, err := c.GetClient()
	if err != nil {
		return 0, err
	}

	// 确保bucket存在
	if err := c.ensureBucketExists(ctx, client, c.config.Bucket); err != nil {
		return 0, fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	// 尝试从reader获取实际大小，避免由于不准确的Size导致文件截断
	putSize := objectSize
	if seeker, ok := reader.(io.Seeker); ok {
		currentOffset, seekErr := seeker.Seek(0, io.SeekCurrent)
		if seekErr == nil {
			if endOffset, sizeErr := seeker.Seek(0, io.SeekEnd); sizeErr == nil {
				putSize = endOffset
			}
			// 恢复原始读取位置
			_, _ = seeker.Seek(currentOffset, io.SeekStart)
		}
	}
	if putSize <= 0 {
		putSize = -1
	}

	// 上传文件
	uploadInfo, err := client.PutObject(ctx, c.config.Bucket, objectName, reader, putSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to upload file: %w", err)
	}

	return uploadInfo.Size, nil
}

// DownloadFile 从MinIO下载文件
func (c *S3Client) DownloadFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	client, err := c.GetClient()
	if err != nil {
		return nil, err
	}

	// 添加超时控制 - 增加超时时间用于大文件
	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)

	object, err := client.GetObject(timeoutCtx, c.config.Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return &objectWithCancel{ReadCloser: object, cancel: cancel}, nil
}

type objectWithCancel struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (o *objectWithCancel) Close() error {
	defer o.cancel()
	return o.ReadCloser.Close()
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
		// 如果生命周期配置失败，记录错误但不影响bucket创建
		// 这可能是由于MinIO版本或配置问题
		fmt.Printf("Warning: failed to set bucket lifecycle: %v\n", err)
	}

	return nil
}

// GetBucketName 获取配置的桶名
func (c *S3Client) GetBucketName() string {
	return c.config.Bucket
}
