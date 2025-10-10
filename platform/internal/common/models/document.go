package model

import (
	"errors"
	"io"
	"time"

	"gorm.io/gorm"
)

// DocumentStatus 文档状态
type DocumentStatus string

const (
	DocumentStatusUploading       DocumentStatus = "uploading"        // 上传中
	DocumentStatusPendingApproval DocumentStatus = "pending_approval" // 待审批
	DocumentStatusPendingPublish  DocumentStatus = "pending_publish"  // 待发布
	DocumentStatusPublished       DocumentStatus = "published"        // 已发布
	DocumentStatusFailed          DocumentStatus = "failed"           // 失败
)

// Document 文档模型
type Document struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Title        string         `json:"title" gorm:"size:200;not null"`                      // 文档标题
	Content      string         `json:"content" gorm:"type:text"`                            // 文档内容
	FileName     string         `json:"file_name" gorm:"size:255;not null"`                  // 原始文件名
	FilePath     string         `json:"file_path" gorm:"size:500;not null"`                  // 文件存储路径
	FileSize     int64          `json:"file_size" gorm:"not null"`                           // 文件大小(字节)
	FileType     string         `json:"file_type" gorm:"size:50;not null"`                   // 文件类型
	Status       DocumentStatus `json:"status" gorm:"size:20;not null;default:'uploading'"`  // 文档状态
	NeedApproval bool           `json:"need_approval" gorm:"default:false"`                  // 是否需要审批
	MimeType     string         `json:"mime_type" gorm:"size:100"`                           // 文件类型
	Version      string         `json:"version" gorm:"size:50;not null;default:'v1.0.0'"`    // 版本号
	UseType      UseType        `json:"use_type" gorm:"size:50;not null;default:'viewable'"` // 使用类型仅查看的文档不可用于对话

	// 关联字段
	SpaceID         uint   `json:"space_id" gorm:"foreignKey:SpaceID"`        // 所属空间ID
	SubSpaceID      uint   `json:"sub_space_id" gorm:"foreignKey:SubSpaceID"` // 所属空间ID
	ClassID         uint   `json:"class_id" gorm:"foreignKey:ClassID"`        // 所属分类ID
	CreatedBy       uint   `json:"created_by" gorm:"foreignKey:UserID"`       // 创建人ID (关联IAM用户)
	CreatorNickName string `json:"creator_nick_name" gorm:"size:100"`         // 创建人昵称
	Department      string `json:"department" gorm:"size:100"`                // 所属部门
	WorkflowID      uint   `json:"workflow_id"`                               // 工作流ID （关联workflow表，上传后为0表示没有，无需审批也为0，需要审批提交后为对应的workflow_id）

	// 标签和摘要
	Tags    string `json:"tags" gorm:"size:500"`     // 标签，JSON格式存储
	Summary string `json:"summary" gorm:"size:1000"` // 文档摘要

	// 时间字段
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 处理相关字段
	ParseError  string     `json:"parse_error" gorm:"type:text"`  // 解析错误信息
	ProcessedAt *time.Time `json:"processed_at"`                  // 处理完成时间
	VectorCount int        `json:"vector_count" gorm:"default:0"` // 向量数量

	// 关联实体
	Workflow Workflow `json:"workflow" gorm:"foreignKey:WorkflowID"`
}

var (
	// ErrOpenAIClientNotConfigured indicates that the OpenAI client configuration is missing.
	ErrOpenAIClientNotConfigured = errors.New("openai client is not configured")
	// ErrNoDocumentsAvailable indicates that there are no usable documents for the chat request.
	ErrNoDocumentsAvailable = errors.New("no documents available for chat")
	// ErrEmptyChatQuestion indicates an empty or whitespace-only question.
	ErrEmptyChatQuestion = errors.New("question is required")
	// ErrVectorClientNotConfigured indicates vector search is unavailable.
	ErrVectorClientNotConfigured = errors.New("vector client is not configured")
)

// ChatDocumentRequest 请求结构
type ChatDocumentRequest struct {
	Question    string `json:"question" binding:"required"`
	DocumentIDs []uint `json:"document_ids"`
	Limit       int    `json:"limit"`
}

// ChatDocumentResponse 响应结构
type ChatDocumentResponse struct {
	Answer  string               `json:"answer"`
	Sources []ChatDocumentSource `json:"sources,omitempty"`
}

// ChatDocumentSource 聊天引用的文档信息
type ChatDocumentSource struct {
	DocumentID uint   `json:"document_id"`
	Title      string `json:"title"`
	FilePath   string `json:"file_path"`
}

// UploadDocumentRequest 上传文档请求
type UploadDocumentRequest struct {
	File            io.Reader
	FileName        string
	FileSize        int64
	ContentType     string
	SpaceID         uint
	SubSpaceID      uint
	ClassID         uint
	Tags            string
	Summary         string
	CreatedBy       uint
	CreatorNickName string
	Department      string
	NeedApproval    bool
	Version         string
	UseType         UseType
}

// UploadDocumentResponse 上传文档响应
type UploadDocumentResponse struct {
	DocumentID uint           `json:"document_id"`
	Status     DocumentStatus `json:"status"`
	Message    string         `json:"message"`
}

type UseType string

const (
	UseTypeViewable   UseType = "viewable"
	UseTypeApplicable UseType = "applicable"
)

// HomepageResponse 首页响应结构
type HomepageResponse struct {
	Spaces []HomepageSpace `json:"spaces"`
}

// HomepageSpace 首页知识库结构
type HomepageSpace struct {
	ID          uint               `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	SubSpaces   []HomepageSubSpace `json:"sub_spaces"`
}

// HomepageSubSpace 首页二级知识库结构
type HomepageSubSpace struct {
	ID          uint               `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Documents   []HomepageDocument `json:"documents"`
}

// HomepageDocument 首页文档结构
type HomepageDocument struct {
	ID              uint           `json:"id"`
	Title           string         `json:"title"`
	FileName        string         `json:"file_name"`
	FileSize        int64          `json:"file_size"`
	FileType        string         `json:"file_type"`
	Status          DocumentStatus `json:"status"`
	CreatorNickName string         `json:"creator_nick_name"`
	Summary         string         `json:"summary"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// TagCloudItem 标签云项
type TagCloudItem struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

// TagCloudResponse 标签云响应
type TagCloudResponse struct {
	Items []TagCloudItem `json:"items"`
}

type KnowledgeSearchRequest struct {
	Query      string `json:"query" binding:"required"`
	Limit      int    `json:"limit"`
	SubSpaceID uint   `json:"sub_space_id"`
	ClassID    uint   `json:"class_id"`
}

type KnowledgeSearchResult struct {
	DocumentID uint    `json:"document_id"`
	ChunkID    uint    `json:"chunk_id"`
	Title      string  `json:"title"`
	Content    string  `json:"content"`
	Snippet    string  `json:"snippet"`
	Score      float64 `json:"score"`
	FileName   string  `json:"file_name"`
}

type KnowledgeSearchResponse struct {
	Items []KnowledgeSearchResult `json:"items"`
}
