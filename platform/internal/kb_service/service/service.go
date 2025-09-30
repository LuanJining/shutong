package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"

	openai "github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/packages/ssestream"

	model "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/models"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/client"
	"gorm.io/gorm"
)

// DocumentService 文档服务
type DocumentService struct {
	db             *gorm.DB
	minioClient    *client.S3Client
	workflowClient *client.WorkflowClient
	openaiClient   *client.OpenAIClient
}

type ChatDocumentStreamResult struct {
	Stream  *ssestream.Stream[openai.ChatCompletionChunk]
	Sources []model.ChatDocumentSource
}

// NewDocumentService 创建文档服务
func NewDocumentService(db *gorm.DB, minioClient *client.S3Client, workflowClient *client.WorkflowClient, openaiClient *client.OpenAIClient) *DocumentService {
	return &DocumentService{
		db:             db,
		minioClient:    minioClient,
		workflowClient: workflowClient,
		openaiClient:   openaiClient,
	}
}

// UploadDocument 上传文档
func (s *DocumentService) UploadDocument(ctx context.Context, req *model.UploadDocumentRequest) (*model.Document, error) {
	// 设置默认值
	if req.CreatedBy == 0 {
		req.CreatedBy = 1 // 默认系统用户ID
	}

	// 生成文件路径
	fileExt := strings.ToLower(filepath.Ext(req.FileName))
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), req.FileName)
	filePath := fmt.Sprintf("%s%s", client.PathPrefixPermanent, fileName)

	// 创建文档记录
	document := &model.Document{
		Title:           strings.TrimSuffix(req.FileName, fileExt),
		FileName:        req.FileName,
		FilePath:        filePath,
		FileSize:        req.FileSize,
		MimeType:        req.ContentType,
		FileType:        fileExt,
		Status:          model.DocumentStatusUploading,
		SpaceID:         req.SpaceID,
		SubSpaceID:      req.SubSpaceID,
		ClassID:         req.ClassID,
		CreatedBy:       req.CreatedBy,
		CreatorNickName: req.CreatorNickName,
		Department:      req.Department,
		Tags:            req.Tags,
		Summary:         req.Summary,
		NeedApproval:    req.NeedApproval,
		Version:         req.Version,
		UseType:         req.UseType,
	}
	// 上传文件到MinIO
	uploadedSize, err := s.minioClient.UploadFile(ctx, filePath, req.File, req.FileSize, req.ContentType)
	if err != nil {
		// 更新文档状态为失败
		s.db.Model(document).Updates(map[string]any{
			"status":      model.DocumentStatusFailed,
			"parse_error": err.Error(),
		})
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	if uploadedSize > 0 {
		document.FileSize = uploadedSize
	}

	// 保存到数据库
	if err := s.db.Create(document).Error; err != nil {
		return nil, fmt.Errorf("failed to create document record: %w", err)
	}

	// 如果需要审批，则创建审批流程
	if document.NeedApproval {
		document, err = s.CreateWorkflow(ctx, document)
		if err != nil {
			return nil, fmt.Errorf("failed to create workflow: %w", err)
		}
		// 设置状态为待审批
		s.db.Model(document).Update("status", model.DocumentStatusPendingApproval)
		document.Status = model.DocumentStatusPendingApproval
	} else {
		// 不需要审批，直接设置为已发布
		s.db.Model(document).Update("status", model.DocumentStatusPendingPublish)
		document.Status = model.DocumentStatusPendingPublish
	}

	return document, nil
}

// GetDocument 获取文档详情
func (s *DocumentService) GetDocument(ctx context.Context, documentID uint) (*model.Document, error) {
	var document model.Document
	if err := s.db.First(&document, documentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("document not found")
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	return &document, nil
}

// GetDocumentsBySpaceId 获取空间下的文档
func (s *DocumentService) GetDocumentsBySpaceId(ctx context.Context, spaceID uint, page int, pageSize int) (*model.PaginationResponse, error) {
	var documents []model.Document
	var total int64

	// 构建查询条件
	query := s.db.Where("space_id = ?", spaceID)

	// 获取总数
	if err := query.Model(&model.Document{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count documents: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&documents).Error; err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}

	// 计算总页数
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &model.PaginationResponse{
		Items:      documents,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// DownloadDocument 下载文档 给前端返回文件流
func (s *DocumentService) DownloadDocument(ctx context.Context, documentID uint) (io.ReadCloser, error) {
	var document model.Document
	if err := s.db.First(&document, documentID).Error; err != nil {
		return nil, fmt.Errorf("document not found")
	}

	// 从MinIO下载文件
	fileReader, err := s.minioClient.DownloadFile(ctx, document.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return fileReader, nil
}

// ChatDocument 基于空间内的文档进行对话
func (s *DocumentService) ChatDocument(ctx context.Context, spaceID uint, req *model.ChatDocumentRequest) (*model.ChatDocumentResponse, error) {
	if s.openaiClient == nil {
		return nil, model.ErrOpenAIClientNotConfigured
	}

	question, fileContents, sources, err := s.prepareChatDocument(ctx, spaceID, req)
	if err != nil {
		return nil, err
	}

	answer, err := s.openaiClient.ChatWithFiles(ctx, question, fileContents)
	if err != nil {
		return nil, fmt.Errorf("failed to chat with documents: %w", err)
	}

	return &model.ChatDocumentResponse{
		Answer:  answer,
		Sources: sources,
	}, nil
}

// ChatDocumentStream 基于空间内的文档进行流式对话
func (s *DocumentService) ChatDocumentStream(ctx context.Context, spaceID uint, req *model.ChatDocumentRequest) (*ChatDocumentStreamResult, error) {
	if s.openaiClient == nil {
		return nil, model.ErrOpenAIClientNotConfigured
	}

	question, fileContents, sources, err := s.prepareChatDocument(ctx, spaceID, req)
	if err != nil {
		return nil, err
	}

	stream, err := s.openaiClient.ChatWithFilesStream(ctx, question, fileContents)
	if err != nil {
		return nil, fmt.Errorf("failed to stream chat with documents: %w", err)
	}

	return &ChatDocumentStreamResult{
		Stream:  stream,
		Sources: sources,
	}, nil
}

func (s *DocumentService) prepareChatDocument(ctx context.Context, spaceID uint, req *model.ChatDocumentRequest) (string, []string, []model.ChatDocumentSource, error) {
	if req == nil {
		return "", nil, nil, model.ErrEmptyChatQuestion
	}

	question := strings.TrimSpace(req.Question)
	if question == "" {
		return "", nil, nil, model.ErrEmptyChatQuestion
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 3
	}
	if limit > 10 {
		limit = 10
	}

	query := s.db.WithContext(ctx).
		Where("space_id = ?", spaceID)

	if len(req.DocumentIDs) > 0 {
		query = query.Where("id IN ?", req.DocumentIDs)
	}

	var documents []model.Document
	if err := query.
		Order("updated_at DESC").
		Limit(limit).
		Find(&documents).Error; err != nil {
		return "", nil, nil, fmt.Errorf("failed to load documents: %w", err)
	}

	log.Println("documents", documents)

	if len(documents) == 0 {
		return "", nil, nil, model.ErrNoDocumentsAvailable
	}

	objectNames := make([]string, 0, len(documents))
	sources := make([]model.ChatDocumentSource, 0, len(documents))
	for _, doc := range documents {
		filePath := strings.TrimSpace(doc.FilePath)
		if filePath == "" {
			continue
		}
		objectNames = append(objectNames, filePath)
		sources = append(sources, model.ChatDocumentSource{
			DocumentID: doc.ID,
			Title:      doc.Title,
			FilePath:   filePath,
		})
	}

	if len(objectNames) == 0 {
		return "", nil, nil, model.ErrNoDocumentsAvailable
	}

	log.Println("objectNames", objectNames)

	fileContents, err := s.openaiClient.ExtractMinioFileContents(ctx, s.minioClient, objectNames)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to load file contents: %w", err)
	}

	log.Println("fileContents loaded", len(fileContents))

	return question, fileContents, sources, nil
}

func (s *DocumentService) CreateWorkflow(ctx context.Context, document *model.Document) (*model.Document, error) {
	step := model.Step{
		StepName:     "文档发布审批",
		StepOrder:    1,
		StepRole:     "content_viewer",
		IsRequired:   true,
		TimeoutHours: 24 * 7,
		Status:       model.StepStatusProcessing,
	}
	workflow := model.Workflow{
		Name:        "文档发布审批流程",
		Description: "用于文档发布的审批流程",
		SpaceID:     document.SpaceID,
		Steps:       []model.Step{step},
	}

	workflowID, err := s.workflowClient.CreateWorkflow(ctx, &workflow, document.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow: %w", err)
	}

	s.db.Model(document).Updates(map[string]any{
		"workflow_id": workflowID,
	})

	return document, nil
}

func (s *DocumentService) CheckWorkflowStatus(ctx context.Context, workflowID uint) (string, error) {
	status, err := s.workflowClient.CheckWorkflowStatus(ctx, workflowID)
	if err != nil {
		return "", fmt.Errorf("failed to check workflow status: %w", err)
	}
	return status, nil
}

func (s *DocumentService) DeleteDocument(ctx context.Context, documentID uint) error {
	var document model.Document
	if err := s.db.First(&document, documentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("document not found")
		}
		return fmt.Errorf("failed to get document: %w", err)
	}

	// minio先删
	err := s.minioClient.DeleteFile(ctx, document.FilePath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	err = s.db.Delete(&model.Document{}, documentID).Error
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}
