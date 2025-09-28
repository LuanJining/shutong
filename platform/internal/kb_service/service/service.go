package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/client"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/model"
	"gorm.io/gorm"
)

// DocumentService 文档服务
type DocumentService struct {
	db             *gorm.DB
	minioClient    *client.S3Client
	qdrantClient   *client.QdrantClient
	ocrClient      *client.OCRClient
	workflowClient *client.WorkflowClient
}

// NewDocumentService 创建文档服务
func NewDocumentService(db *gorm.DB, minioClient *client.S3Client, qdrantClient *client.QdrantClient, ocrClient *client.OCRClient, workflowClient *client.WorkflowClient) *DocumentService {
	return &DocumentService{
		db:             db,
		minioClient:    minioClient,
		qdrantClient:   qdrantClient,
		ocrClient:      ocrClient,
		workflowClient: workflowClient,
	}
}

// UploadDocumentRequest 上传文档请求
type UploadDocumentRequest struct {
	File         io.Reader
	FileName     string
	FileSize     int64
	ContentType  string
	SpaceID      uint
	Visibility   string
	Urgency      string
	Tags         string
	Summary      string
	CreatedBy    uint
	Department   string
	NeedApproval bool
}

// UploadDocumentResponse 上传文档响应
type UploadDocumentResponse struct {
	DocumentID uint                 `json:"document_id"`
	Status     model.DocumentStatus `json:"status"`
	Message    string               `json:"message"`
}

// UploadDocument 上传文档
func (s *DocumentService) UploadDocument(ctx context.Context, req *UploadDocumentRequest) (*UploadDocumentResponse, error) {
	// 设置默认值
	if req.Visibility == "" {
		req.Visibility = "private"
	}
	if req.Urgency == "" {
		req.Urgency = "normal"
	}
	if req.CreatedBy == 0 {
		req.CreatedBy = 1 // 默认系统用户ID
	}

	// 生成文件路径
	fileExt := strings.ToLower(filepath.Ext(req.FileName))
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), req.FileName)
	filePath := fmt.Sprintf("%s%s", client.PathPrefixPermanent, fileName)

	// 创建文档记录
	document := &model.Document{
		Title:        strings.TrimSuffix(req.FileName, fileExt),
		FileName:     req.FileName,
		FilePath:     filePath,
		FileSize:     req.FileSize,
		FileType:     fileExt,
		MimeType:     req.ContentType,
		Status:       model.DocumentStatusUploading,
		Visibility:   model.DocumentVisibility(req.Visibility),
		Urgency:      model.DocumentUrgency(req.Urgency),
		SpaceID:      req.SpaceID,
		CreatedBy:    req.CreatedBy,
		Department:   req.Department,
		Tags:         req.Tags,
		Summary:      req.Summary,
		NeedApproval: req.NeedApproval,
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
	}

	// 更新文档状态为处理中
	s.db.Model(document).Update("status", model.DocumentStatusPendingApproval)

	return &UploadDocumentResponse{
		DocumentID: document.ID,
		Status:     document.Status,
		Message:    "Document uploaded successfully",
	}, nil
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

func (s *DocumentService) CreateWorkflow(ctx context.Context, document *model.Document) (*model.Document, error) {
	// 直接写死空间管理员审批，后期再扩展
	workflowStep := model.WorkflowStep{
		StepName:         "文档发布审批",
		StepOrder:        1,
		ApproverType:     "space_admin",
		ApproverID:       0,
		IsRequired:       true,
		TimeoutHours:     24,
		ApprovalStrategy: "any",
	}
	workflow := model.Workflow{
		Name:        "文档发布审批流程",
		Description: "用于文档发布的审批流程",
		SpaceID:     document.SpaceID,
		Priority:    1,
		Steps:       []model.WorkflowStep{workflowStep},
	}
	workflowIDStr, err := s.workflowClient.CreateWorkflow(ctx, &workflow, document.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow: %w", err)
	}

	// 将字符串ID转换为uint
	workflowID, err := strconv.ParseUint(workflowIDStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse workflow ID: %w", err)
	}

	s.db.Model(document).Update("workflow_id", workflowID)

	// 更新workflow对象的ID，用于启动工作流
	workflow.ID = uint(workflowID)
	_, err = s.workflowClient.StartWorkflow(ctx, &workflow, document.CreatedBy, document.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to start workflow: %w", err)
	}
	return document, nil
}
