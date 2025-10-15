package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
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
	ocrClient      *client.PaddleOCRClient
	vectorClient   *client.QdrantClient
}

var errUnsupportedFileType = errors.New("unsupported file type for text extraction")

type chunkSearchResult struct {
	Document *model.Document
	ChunkID  uint
	Content  string
	Score    float64
	FileName string
}

type ChatDocumentStreamResult struct {
	Stream  *ssestream.Stream[openai.ChatCompletionChunk]
	Sources []model.ChatDocumentSource
}

// NewDocumentService 创建文档服务
func NewDocumentService(
	db *gorm.DB,
	minioClient *client.S3Client,
	workflowClient *client.WorkflowClient,
	openaiClient *client.OpenAIClient,
	ocrClient *client.PaddleOCRClient,
	vectorClient *client.QdrantClient,
) *DocumentService {
	return &DocumentService{
		db:             db,
		minioClient:    minioClient,
		workflowClient: workflowClient,
		openaiClient:   openaiClient,
		ocrClient:      ocrClient,
		vectorClient:   vectorClient,
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
		document, err = s.StartWorkflow(ctx, document)
		if err != nil {
			return nil, fmt.Errorf("failed to start workflow: %w", err)
		}
	} else {
		// 不需要审批，直接设置为已发布
		s.db.Model(document).Update("status", model.DocumentStatusPendingPublish)
		document.Status = model.DocumentStatusPendingPublish
	}

	// 异步处理文档（OCR/向量化）
	go func(docID uint) {
		if err := s.ProcessDocument(context.Background(), docID); err != nil {
			log.Printf("failed to process document %d: %v", docID, err)
		}
	}(document.ID)

	return document, nil
}

// GetDocument 获取文档详情
func (s *DocumentService) GetDocument(ctx context.Context, documentID uint) (*model.Document, error) {
	var document model.Document
	if err := s.db.Preload("Workflow").First(&document, documentID).Error; err != nil {
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
	if err := query.Preload("Workflow").Order("created_at DESC").
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

	if len(req.DocumentIDs) > 0 || s.vectorClient == nil {
		return s.prepareChatDocumentFromFiles(ctx, spaceID, question, req, limit)
	}

	chunks, err := s.searchChunks(ctx, spaceID, question, 0, 0, limit)
	if err != nil {
		if errors.Is(err, model.ErrVectorClientNotConfigured) {
			return s.prepareChatDocumentFromFiles(ctx, spaceID, question, req, limit)
		}
		return "", nil, nil, err
	}

	// 如果向量搜索没有结果，降级：获取该space所有chunks
	if len(chunks) == 0 {
		log.Printf("Vector search returned no results for space %d, falling back to all chunks", spaceID)
		chunks, err = s.getAllSpaceChunks(ctx, spaceID, limit)
		if err != nil {
			return "", nil, nil, err
		}
	}

	fileContents := make([]string, 0, len(chunks))
	sources := make([]model.ChatDocumentSource, 0)
	seenDocuments := make(map[uint]struct{})

	for _, chunk := range chunks {
		if chunk.Document == nil {
			continue
		}
		content := strings.TrimSpace(chunk.Content)
		if content == "" {
			continue
		}
		formatted := fmt.Sprintf("标题: %s\n文件: %s\n相关片段:\n%s", chunk.Document.Title, chunk.FileName, content)
		fileContents = append(fileContents, formatted)

		if _, exists := seenDocuments[chunk.Document.ID]; !exists {
			sources = append(sources, model.ChatDocumentSource{
				DocumentID: chunk.Document.ID,
				Title:      chunk.Document.Title,
				FilePath:   chunk.Document.FilePath,
			})
			seenDocuments[chunk.Document.ID] = struct{}{}
		}
	}

	// 如果chunks内容都为空，说明该space真的没有可用内容
	if len(fileContents) == 0 {
		return "", nil, nil, model.ErrNoDocumentsAvailable
	}

	return question, fileContents, sources, nil
}

func (s *DocumentService) prepareChatDocumentFromFiles(ctx context.Context, spaceID uint, question string, req *model.ChatDocumentRequest, limit int) (string, []string, []model.ChatDocumentSource, error) {
	query := s.db.WithContext(ctx).
		Where("space_id = ?", spaceID).
		Where("status IN ?", []model.DocumentStatus{
			model.DocumentStatusPublished,
			model.DocumentStatusPendingPublish,
		})

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

	fileContents, err := s.openaiClient.ExtractMinioFileContents(ctx, s.minioClient, objectNames)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to load file contents: %w", err)
	}

	return question, fileContents, sources, nil
}

func (s *DocumentService) CreateWorkflow(ctx context.Context, document *model.Document) (*model.Document, error) {
	step := model.Step{
		StepName:     "文档发布审批",
		StepOrder:    1,
		StepRole:     string(model.SpaceMemberRoleApprover),
		IsRequired:   true,
		TimeoutHours: 24 * 7,
		Status:       model.StepStatusProcessing,
	}
	workflow := model.Workflow{
		Name:         "文档发布审批流程",
		Description:  "用于文档发布的审批流程",
		SpaceID:      document.SpaceID,
		Status:       model.WorkflowStatusProcessing,
		Steps:        []model.Step{step},
		ResourceType: "document",
		ResourceID:   document.ID,
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

func (s *DocumentService) StartWorkflow(ctx context.Context, document *model.Document) (*model.Document, error) {
	_, err := s.workflowClient.StartWorkflow(ctx, document.WorkflowID, document.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to start workflow: %w", err)
	}

	s.db.Model(document).Updates(map[string]any{
		"status":     model.DocumentStatusPendingApproval,
		"updated_at": time.Now(),
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

// GetHomepageDocuments 获取首页展示的文档
// 返回5个知识库，每个知识库包含3个二级知识库，每个二级知识库包含6个文档
func (s *DocumentService) GetHomepageDocuments(ctx context.Context) (*model.HomepageResponse, error) {
	// 获取5个知识库（按创建时间倒序）
	var spaces []model.Space
	if err := s.db.WithContext(ctx).
		Where("status = ?", 1).
		Order("created_at DESC").
		Limit(5).
		Find(&spaces).Error; err != nil {
		return nil, fmt.Errorf("failed to get spaces: %w", err)
	}

	response := &model.HomepageResponse{
		Spaces: make([]model.HomepageSpace, 0, len(spaces)),
	}

	// 遍历每个知识库
	for _, space := range spaces {
		homepageSpace := model.HomepageSpace{
			ID:          space.ID,
			Name:        space.Name,
			Description: space.Description,
			SubSpaces:   make([]model.HomepageSubSpace, 0, 3),
		}

		// 获取该知识库下的3个二级知识库
		var subSpaces []model.SubSpace
		if err := s.db.WithContext(ctx).
			Where("space_id = ? AND status = ?", space.ID, 1).
			Order("created_at DESC").
			Limit(3).
			Find(&subSpaces).Error; err != nil {
			log.Printf("failed to get subspaces for space %d: %v", space.ID, err)
			continue
		}

		// 遍历每个二级知识库
		for _, subSpace := range subSpaces {
			homepageSubSpace := model.HomepageSubSpace{
				ID:          subSpace.ID,
				Name:        subSpace.Name,
				Description: subSpace.Description,
				Documents:   make([]model.HomepageDocument, 0, 6),
			}

			// 获取该二级知识库下的6个文档（只获取已发布的文档）
			var documents []model.Document
			if err := s.db.WithContext(ctx).
				Where("sub_space_id = ? AND status IN ?", subSpace.ID, []model.DocumentStatus{
					model.DocumentStatusPublished,
				}).
				Order("created_at DESC").
				Limit(6).
				Find(&documents).Error; err != nil {
				log.Printf("failed to get documents for subspace %d: %v", subSpace.ID, err)
				continue
			}

			// 将文档转换为首页文档结构
			for _, doc := range documents {
				homepageSubSpace.Documents = append(homepageSubSpace.Documents, model.HomepageDocument{
					ID:              doc.ID,
					Title:           doc.Title,
					FileName:        doc.FileName,
					FileSize:        doc.FileSize,
					FileType:        doc.FileType,
					Status:          doc.Status,
					CreatorNickName: doc.CreatorNickName,
					Summary:         doc.Summary,
					CreatedAt:       doc.CreatedAt,
					UpdatedAt:       doc.UpdatedAt,
				})
			}

			homepageSpace.SubSpaces = append(homepageSpace.SubSpaces, homepageSubSpace)
		}

		response.Spaces = append(response.Spaces, homepageSpace)
	}

	return response, nil
}

// GetTagCloud 聚合标签云
func (s *DocumentService) GetTagCloud(ctx context.Context, spaceID, subSpaceID uint, limit int) ([]model.TagCloudItem, error) {
	if limit <= 0 {
		limit = 50
	}

	query := s.db.WithContext(ctx).
		Model(&model.Document{}).
		Where("tags IS NOT NULL AND tags <> ''").
		Where("status IN ?", []model.DocumentStatus{
			model.DocumentStatusPublished,
		})

	if spaceID > 0 {
		query = query.Where("space_id = ?", spaceID)
	}

	if subSpaceID > 0 {
		query = query.Where("sub_space_id = ?", subSpaceID)
	}

	var tagStrings []string
	if err := query.Pluck("tags", &tagStrings).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch document tags: %w", err)
	}

	if len(tagStrings) == 0 {
		return []model.TagCloudItem{}, nil
	}

	tagCounts := make(map[string]int)
	for _, raw := range tagStrings {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}

		tags := make([]string, 0)
		if err := json.Unmarshal([]byte(raw), &tags); err != nil {
			// 尝试以逗号分隔的字符串解析
			for tag := range strings.SplitSeq(raw, ",") {
				tag = strings.TrimSpace(tag)
				tag = strings.Trim(tag, "\"'")
				if tag != "" {
					tags = append(tags, tag)
				}
			}
		}

		if len(tags) == 0 {
			continue
		}

		seen := make(map[string]struct{}, len(tags))
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			if tag == "" {
				continue
			}
			if _, exists := seen[tag]; exists {
				continue
			}
			seen[tag] = struct{}{}
			tagCounts[tag]++
		}
	}

	items := make([]model.TagCloudItem, 0, len(tagCounts))
	for tag, count := range tagCounts {
		items = append(items, model.TagCloudItem{
			Tag:   tag,
			Count: count,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Count == items[j].Count {
			return items[i].Tag < items[j].Tag
		}
		return items[i].Count > items[j].Count
	})

	if limit < len(items) {
		items = items[:limit]
	}

	return items, nil
}

// SearchKnowledge 基于向量的知识搜索
func (s *DocumentService) SearchKnowledge(ctx context.Context, req *model.KnowledgeSearchRequest) (*model.KnowledgeSearchResponse, error) {
	if req == nil {
		return nil, errors.New("search request is nil")
	}

	query := strings.TrimSpace(req.Query)
	if query == "" {
		return nil, fmt.Errorf("query is required")
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}

	chunks, err := s.searchChunks(ctx, req.SpaceID, query, req.SubSpaceID, req.ClassID, limit)
	if err != nil {
		return nil, err
	}

	results := make([]model.KnowledgeSearchResult, 0, len(chunks))
	for _, chunk := range chunks {
		if chunk.Document == nil {
			continue
		}
		content := strings.TrimSpace(chunk.Content)
		if content == "" {
			continue
		}

		results = append(results, model.KnowledgeSearchResult{
			DocumentID: chunk.Document.ID,
			ChunkID:    chunk.ChunkID,
			Title:      chunk.Document.Title,
			Content:    content,
			Snippet:    buildSnippet(content, 200),
			Score:      chunk.Score,
			FileName:   chunk.FileName,
		})
	}

	return &model.KnowledgeSearchResponse{
		Items: results,
	}, nil
}

// getAllSpaceChunks 获取space下所有chunks（降级策略）
func (s *DocumentService) getAllSpaceChunks(ctx context.Context, spaceID uint, limit int) ([]chunkSearchResult, error) {
	// 查询该space下已发布的文档
	var documents []model.Document
	if err := s.db.WithContext(ctx).
		Where("space_id = ?", spaceID).
		Where("status IN ?", []model.DocumentStatus{
			model.DocumentStatusPublished,
			model.DocumentStatusPendingPublish,
		}).
		Order("updated_at DESC").
		Find(&documents).Error; err != nil {
		return nil, fmt.Errorf("failed to load documents: %w", err)
	}

	if len(documents) == 0 {
		return nil, model.ErrNoDocumentsAvailable
	}

	// 提取document IDs
	docIDs := make([]uint, len(documents))
	docMap := make(map[uint]*model.Document)
	for i, doc := range documents {
		docIDs[i] = doc.ID
		docCopy := doc
		docMap[doc.ID] = &docCopy
	}

	// 查询这些文档的chunks
	var dbChunks []model.DocumentChunk
	if err := s.db.WithContext(ctx).
		Where("document_id IN ?", docIDs).
		Order("document_id DESC, index ASC").
		Limit(limit).
		Find(&dbChunks).Error; err != nil {
		return nil, fmt.Errorf("failed to load chunks: %w", err)
	}

	if len(dbChunks) == 0 {
		return nil, model.ErrNoDocumentsAvailable
	}

	// 转换为chunkSearchResult格式
	results := make([]chunkSearchResult, 0, len(dbChunks))
	for _, chunk := range dbChunks {
		doc, exists := docMap[chunk.DocumentID]
		if !exists {
			continue
		}
		results = append(results, chunkSearchResult{
			Document: doc,
			ChunkID:  chunk.ID,
			Content:  chunk.Content,
			Score:    1.0, // 降级策略，没有相关性分数
			FileName: doc.FileName,
		})
	}

	return results, nil
}

func (s *DocumentService) searchChunks(ctx context.Context, spaceID uint, query string, subSpaceID, classID uint, limit int) ([]chunkSearchResult, error) {
	if s.vectorClient == nil {
		return nil, model.ErrVectorClientNotConfigured
	}

	// 使用OpenAI生成查询向量
	vector, err := s.generateEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}
	vector = ensureVectorSize(vector, s.vectorClient.VectorSize())

	var must []client.QdrantCondition
	if spaceID > 0 {
		must = append(must, client.QdrantCondition{
			Key:   "space_id",
			Match: client.QdrantMatch{Value: spaceID},
		})
	}
	if subSpaceID > 0 {
		must = append(must, client.QdrantCondition{
			Key:   "sub_space_id",
			Match: client.QdrantMatch{Value: subSpaceID},
		})
	}
	if classID > 0 {
		must = append(must, client.QdrantCondition{
			Key:   "class_id",
			Match: client.QdrantMatch{Value: classID},
		})
	}

	var filter *client.QdrantFilter
	if len(must) > 0 {
		filter = &client.QdrantFilter{Must: must}
	}

	results, err := s.vectorClient.SearchPoints(ctx, vector, limit, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to search vector database: %w", err)
	}
	if len(results) == 0 {
		return []chunkSearchResult{}, nil
	}

	type rawChunk struct {
		docID    uint
		chunkID  uint
		content  string
		fileName string
		score    float64
	}

	rawChunks := make([]rawChunk, 0, len(results))
	docIDSet := make(map[uint]struct{})

	for _, res := range results {
		docID, ok := payloadUint(res.Payload["document_id"])
		if !ok || docID == 0 {
			continue
		}
		chunkID, ok := payloadUint(res.Payload["chunk_id"])
		if !ok {
			continue
		}
		content := payloadString(res.Payload["content"])
		if strings.TrimSpace(content) == "" {
			continue
		}
		fileName := payloadString(res.Payload["file_name"])

		rawChunks = append(rawChunks, rawChunk{
			docID:    docID,
			chunkID:  chunkID,
			content:  content,
			fileName: fileName,
			score:    res.Score,
		})
		docIDSet[docID] = struct{}{}
	}

	if len(rawChunks) == 0 {
		return []chunkSearchResult{}, nil
	}

	docIDs := make([]uint, 0, len(docIDSet))
	for id := range docIDSet {
		docIDs = append(docIDs, id)
	}

	var documents []model.Document
	dbQuery := s.db.WithContext(ctx).
		Where("id IN ?", docIDs).
		Where("status IN ?", []model.DocumentStatus{
			model.DocumentStatusPublished,
			model.DocumentStatusPendingPublish,
		})

	if spaceID > 0 {
		dbQuery = dbQuery.Where("space_id = ?", spaceID)
	}

	if err := dbQuery.Find(&documents).Error; err != nil {
		return nil, fmt.Errorf("failed to load documents: %w", err)
	}

	docMap := make(map[uint]*model.Document, len(documents))
	for i := range documents {
		doc := &documents[i]
		docMap[doc.ID] = doc
	}

	resultsChunks := make([]chunkSearchResult, 0, len(rawChunks))
	for _, chunk := range rawChunks {
		doc, ok := docMap[chunk.docID]
		if !ok {
			continue
		}
		resultsChunks = append(resultsChunks, chunkSearchResult{
			Document: doc,
			ChunkID:  chunk.chunkID,
			Content:  chunk.content,
			Score:    chunk.score,
			FileName: chunk.fileName,
		})
	}

	return resultsChunks, nil
}

// ProcessDocument 执行文档OCR与向量化处理
func (s *DocumentService) ProcessDocument(ctx context.Context, documentID uint) error {
	if s.minioClient == nil {
		return errors.New("minio client is not configured")
	}

	var document model.Document
	if err := s.db.WithContext(ctx).First(&document, documentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("document %d not found", documentID)
		}
		return fmt.Errorf("failed to load document %d: %w", documentID, err)
	}

	reader, err := s.minioClient.DownloadFile(ctx, document.FilePath)
	if err != nil {
		s.markDocumentProcessingError(document.ID, fmt.Errorf("failed to download file: %w", err))
		return err
	}
	defer reader.Close()

	fileBytes, err := io.ReadAll(reader)
	if err != nil {
		s.markDocumentProcessingError(document.ID, fmt.Errorf("failed to read file: %w", err))
		return err
	}

	text, err := extractPlainText(document.FileType, fileBytes)
	if err != nil {
		if errors.Is(err, errUnsupportedFileType) {
			if s.ocrClient == nil {
				err = fmt.Errorf("unsupported file type %s and OCR client not configured", document.FileType)
				s.markDocumentProcessingError(document.ID, err)
				return err
			}
			text, err = s.ocrClient.Recognize(ctx, document.FileName, fileBytes)
			if err != nil {
				err = fmt.Errorf("ocr recognition failed: %w", err)
				s.markDocumentProcessingError(document.ID, err)
				return err
			}
		} else {
			s.markDocumentProcessingError(document.ID, err)
			return err
		}
	}

	if strings.TrimSpace(text) == "" {
		err = errors.New("empty text extracted from document")
		s.markDocumentProcessingError(document.ID, err)
		return err
	}

	chunks := splitIntoChunks(text, 800, 120)

	if err := s.storeChunks(ctx, &document, chunks); err != nil {
		s.markDocumentProcessingError(document.ID, err)
		return err
	}

	return nil
}

func (s *DocumentService) storeChunks(ctx context.Context, document *model.Document, chunks []string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("document_id = ?", document.ID).Delete(&model.DocumentChunk{}).Error; err != nil {
			return fmt.Errorf("failed to clear existing chunks: %w", err)
		}

		metadataJSON, err := json.Marshal(map[string]any{
			"space_id":     document.SpaceID,
			"sub_space_id": document.SubSpaceID,
			"class_id":     document.ClassID,
		})
		if err != nil {
			return fmt.Errorf("failed to marshal chunk metadata: %w", err)
		}

		// 批量生成所有chunk的embeddings
		validChunks := make([]string, 0, len(chunks))
		validIndices := make([]int, 0, len(chunks))
		for idx, content := range chunks {
			if strings.TrimSpace(content) != "" {
				validChunks = append(validChunks, content)
				validIndices = append(validIndices, idx)
			}
		}

		if len(validChunks) == 0 {
			return errors.New("no valid chunks to store")
		}

		embeddings, err := s.generateEmbeddingBatch(ctx, validChunks)
		if err != nil {
			return fmt.Errorf("failed to generate embeddings: %w", err)
		}

		points := make([]client.QdrantPoint, 0, len(validChunks))
		createdChunks := 0
		for i, content := range validChunks {
			idx := validIndices[i]
			embeddingVector := embeddings[i]
			vectorID := uuid.NewString()

			chunk := model.DocumentChunk{
				DocumentID: document.ID,
				Index:      idx,
				Content:    content,
				TokenCount: countTokens(content),
				Metadata:   string(metadataJSON),
				VectorID:   vectorID,
			}

			if err := tx.Create(&chunk).Error; err != nil {
				return fmt.Errorf("failed to create chunk: %w", err)
			}
			createdChunks++

			if s.vectorClient != nil {
				payload := map[string]any{
					"document_id":  document.ID,
					"chunk_id":     chunk.ID,
					"space_id":     document.SpaceID,
					"sub_space_id": document.SubSpaceID,
					"class_id":     document.ClassID,
					"created_at":   chunk.CreatedAt,
					"title":        document.Title,
					"file_name":    document.FileName,
					"content":      content,
				}

				points = append(points, client.QdrantPoint{
					ID:      vectorID,
					Vector:  embeddingVector,
					Payload: payload,
				})
			}
		}

		if s.vectorClient != nil && len(points) > 0 {
			if err := s.vectorClient.UpsertPoints(ctx, points); err != nil {
				return fmt.Errorf("failed to upsert document vectors: %w", err)
			}
		}

		processedAt := time.Now()
		updateData := map[string]any{
			"processed_at": &processedAt,
			"parse_error":  "",
			"vector_count": createdChunks,
		}

		if err := tx.Model(&model.Document{}).Where("id = ?", document.ID).Updates(updateData).Error; err != nil {
			return fmt.Errorf("failed to update document meta: %w", err)
		}

		return nil
	})
}

func (s *DocumentService) markDocumentProcessingError(documentID uint, procErr error) {
	if procErr == nil {
		return
	}

	updateData := map[string]any{
		"parse_error":  procErr.Error(),
		"vector_count": 0,
		"processed_at": nil,
	}

	if err := s.db.Model(&model.Document{}).Where("id = ?", documentID).Updates(updateData).Error; err != nil {
		log.Printf("failed to update document %d after processing error: %v", documentID, err)
	}
}

func extractPlainText(fileType string, data []byte) (string, error) {
	switch strings.ToLower(fileType) {
	case ".txt", ".md", ".csv", ".log":
		return string(data), nil
	case ".json":
		var buf bytes.Buffer
		if err := json.Indent(&buf, data, "", "  "); err == nil {
			return buf.String(), nil
		}
		return string(data), nil
	case ".html", ".htm":
		result := stripHTMLTags(string(data))
		if strings.TrimSpace(result) == "" {
			return "", fmt.Errorf("html content is empty after stripping tags")
		}
		return result, nil
	default:
		return "", fmt.Errorf("%w: %s", errUnsupportedFileType, fileType)
	}
}

func splitIntoChunks(text string, chunkSize int, overlap int) []string {
	clean := strings.TrimSpace(text)
	if clean == "" {
		return []string{}
	}

	if chunkSize <= 0 {
		chunkSize = 800
	}
	if overlap < 0 {
		overlap = 0
	}
	if overlap >= chunkSize {
		overlap = chunkSize / 2
	}

	runes := []rune(clean)
	var chunks []string
	start := 0

	for start < len(runes) {
		end := min(start+chunkSize, len(runes))

		chunk := strings.TrimSpace(string(runes[start:end]))
		if chunk != "" {
			chunks = append(chunks, chunk)
		}

		if end == len(runes) {
			break
		}

		start = max(end-overlap, 0)
	}

	return chunks
}

func countTokens(content string) int {
	return len(strings.Fields(content))
}

// generateEmbedding 生成单个文本的向量，优先使用OpenAI，降级到简单统计特征
func (s *DocumentService) generateEmbedding(ctx context.Context, text string) ([]float64, error) {
	if s.openaiClient != nil {
		embedding, err := s.openaiClient.CreateEmbedding(ctx, text)
		if err != nil {
			log.Printf("OpenAI embedding failed, falling back to simple embedding: %v", err)
			return simpleEmbedding(text), nil
		}
		return embedding, nil
	}
	// 降级策略：使用简单统计特征
	return simpleEmbedding(text), nil
}

// generateEmbeddingBatch 批量生成向量
func (s *DocumentService) generateEmbeddingBatch(ctx context.Context, texts []string) ([][]float64, error) {
	if s.openaiClient != nil {
		embeddings, err := s.openaiClient.CreateEmbeddingBatch(ctx, texts)
		if err != nil {
			log.Printf("OpenAI batch embedding failed, falling back to simple embedding: %v", err)
			// 降级策略：逐个生成简单向量
			result := make([][]float64, len(texts))
			for i, text := range texts {
				result[i] = simpleEmbedding(text)
			}
			return result, nil
		}
		return embeddings, nil
	}
	// 降级策略：使用简单统计特征
	result := make([][]float64, len(texts))
	for i, text := range texts {
		result[i] = simpleEmbedding(text)
	}
	return result, nil
}

// simpleEmbedding 简单统计特征向量（降级方案）
func simpleEmbedding(content string) []float64 {
	words := strings.Fields(content)
	wordCount := len(words)

	// Count runes (characters) for better Chinese support
	runes := []rune(content)
	runeCount := len(runes)

	uniqueWords := make(map[string]struct{}, wordCount)
	var totalWordLen int
	var digitCount, upperCount, chineseCount int

	for _, w := range words {
		uniqueWords[strings.ToLower(w)] = struct{}{}
		totalWordLen += len(w)
	}

	for _, r := range runes {
		if unicode.IsDigit(r) {
			digitCount++
		}
		if unicode.IsUpper(r) {
			upperCount++
		}
		// Check if character is Chinese (CJK Unified Ideographs)
		if unicode.Is(unicode.Han, r) {
			chineseCount++
		}
	}

	// For Chinese-dominant text, use character count; otherwise use word count
	var effectiveTokenCount float64
	var uniqueTokenCount float64
	var avgTokenLen float64

	chineseRatio := float64(chineseCount) / float64(runeCount)
	if chineseRatio > 0.3 {
		// Chinese-dominant text: use character-based metrics
		effectiveTokenCount = float64(runeCount)
		// For Chinese, estimate unique characters
		uniqueChars := make(map[rune]struct{})
		for _, r := range runes {
			uniqueChars[r] = struct{}{}
		}
		uniqueTokenCount = float64(len(uniqueChars))
		avgTokenLen = 1.0 // Each Chinese character is one unit
	} else {
		// English/Latin-dominant text: use word-based metrics
		effectiveTokenCount = float64(wordCount)
		uniqueTokenCount = float64(len(uniqueWords))
		if wordCount > 0 {
			avgTokenLen = float64(totalWordLen) / float64(wordCount)
		}
	}

	return []float64{
		float64(runeCount),               // 字符数（更准确的文本长度）
		effectiveTokenCount,              // 有效token数（中文用字符数，英文用单词数）
		uniqueTokenCount,                 // 去重token数
		float64(digitCount),              // 数字字符数
		float64(upperCount),              // 大写字母数
		avgTokenLen,                      // 平均token长度
		float64(countSentences(content)), // 句子数
	}
}

func countSentences(content string) int {
	count := 0
	for _, ch := range content {
		if ch == '.' || ch == '!' || ch == '?' || ch == '。' || ch == '！' || ch == '？' {
			count++
		}
	}
	if count == 0 && strings.TrimSpace(content) != "" {
		return 1
	}
	return count
}

func stripHTMLTags(input string) string {
	var output strings.Builder
	inTag := false
	for _, r := range input {
		switch r {
		case '<':
			inTag = true
		case '>':
			inTag = false
		default:
			if !inTag {
				output.WriteRune(r)
			}
		}
	}
	return output.String()
}

func ensureVectorSize(vec []float64, size int) []float64 {
	if size <= 0 {
		return vec
	}
	if len(vec) == size {
		return vec
	}
	result := make([]float64, size)
	copyCount := min(len(vec), size)
	copy(result, vec[:copyCount])
	return result
}

func payloadUint(value any) (uint, bool) {
	switch v := value.(type) {
	case float64:
		if v < 0 {
			return 0, false
		}
		return uint(v + 0.5), true
	case int:
		if v < 0 {
			return 0, false
		}
		return uint(v), true
	case int32:
		if v < 0 {
			return 0, false
		}
		return uint(v), true
	case int64:
		if v < 0 {
			return 0, false
		}
		return uint(v), true
	case uint:
		return v, true
	case uint32:
		return uint(v), true
	case uint64:
		return uint(v), true
	case json.Number:
		if i, err := v.Int64(); err == nil && i >= 0 {
			return uint(i), true
		}
	case string:
		if strings.TrimSpace(v) == "" {
			return 0, false
		}
		if num, err := strconv.ParseUint(v, 10, 64); err == nil {
			return uint(num), true
		}
	}
	return 0, false
}

func payloadString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func buildSnippet(content string, maxRunes int) string {
	trimmed := strings.TrimSpace(content)
	if maxRunes <= 0 {
		return trimmed
	}
	runes := []rune(trimmed)
	if len(runes) <= maxRunes {
		return trimmed
	}
	return strings.TrimSpace(string(runes[:maxRunes])) + "..."
}

func (s *DocumentService) PublishDocument(ctx context.Context, document *model.Document) error {
	return s.db.WithContext(ctx).Model(document).Update("status", document.Status).Error
}

func (s *DocumentService) UnpublishDocument(ctx context.Context, document *model.Document) error {
	return s.db.WithContext(ctx).Model(document).Update("status", document.Status).Error
}
