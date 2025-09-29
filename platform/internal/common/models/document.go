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

// DocumentVisibility 文档可见性
type DocumentVisibility string

const (
	DocumentVisibilityPublic    DocumentVisibility = "public"    // 公开
	DocumentVisibilityInternal  DocumentVisibility = "internal"  // 内部
	DocumentVisibilityPrivate   DocumentVisibility = "private"   // 私有
	DocumentVisibilityProtected DocumentVisibility = "protected" // 受保护
)

// DocumentUrgency 紧急程度
type DocumentUrgency string

const (
	DocumentUrgencyNormal DocumentUrgency = "normal" // 一般
	DocumentUrgencyUrgent DocumentUrgency = "urgent" // 紧急
)

// Document 文档模型
type Document struct {
	ID           uint               `json:"id" gorm:"primaryKey"`
	Title        string             `json:"title" gorm:"size:200;not null"`                       // 文档标题
	Content      string             `json:"content" gorm:"type:text"`                             // 文档内容
	OriginalText string             `json:"original_text" gorm:"type:text"`                       // 原始解析文本
	FileName     string             `json:"file_name" gorm:"size:255;not null"`                   // 原始文件名
	FilePath     string             `json:"file_path" gorm:"size:500;not null"`                   // 文件存储路径
	FileSize     int64              `json:"file_size" gorm:"not null"`                            // 文件大小(字节)
	FileType     string             `json:"file_type" gorm:"size:50;not null"`                    // 文件类型
	MimeType     string             `json:"mime_type" gorm:"size:100"`                            // MIME类型
	Status       DocumentStatus     `json:"status" gorm:"size:20;not null;default:'uploading'"`   // 文档状态
	Visibility   DocumentVisibility `json:"visibility" gorm:"size:20;not null;default:'private'"` // 可见性
	Urgency      DocumentUrgency    `json:"urgency" gorm:"size:20;not null;default:'normal'"`     // 紧急程度
	NeedApproval bool               `json:"need_approval" gorm:"default:false"`                   // 是否需要审批
	// 关联字段
	SpaceID         uint   `json:"space_id" gorm:"not null"`          // 所属空间ID
	CreatedBy       uint   `json:"created_by" gorm:"not null"`        // 创建人ID (关联IAM用户)
	CreatorNickName string `json:"creator_nick_name" gorm:"size:100"` // 创建人昵称
	Approver        uint   `json:"approver"`                          // 审批人ID (关联IAM用户)
	Department      string `json:"department" gorm:"size:100"`        // 所属部门
	WorkflowID      uint   `json:"workflow_id"`                       // 工作流ID （关联workflow表，上传后为0表示没有，无需审批也为0，需要审批提交后为对应的workflow_id）

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
}

// DocumentChunk 文档分块模型
type DocumentChunk struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	DocumentID uint      `json:"document_id" gorm:"not null;index"` // 关联文档ID
	ChunkIndex int       `json:"chunk_index" gorm:"not null"`       // 分块索引
	Content    string    `json:"content" gorm:"type:text;not null"` // 分块内容
	VectorID   string    `json:"vector_id" gorm:"size:100;unique"`  // 向量ID
	TokenCount int       `json:"token_count" gorm:"default:0"`      // Token数量
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Workflow 工作流模型
type Workflow struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description"`
	SpaceID     uint           `json:"space_id" binding:"required"`
	Priority    int            `json:"priority"`
	Steps       []WorkflowStep `json:"steps" binding:"required"`
}

var (
	// ErrOpenAIClientNotConfigured indicates that the OpenAI client configuration is missing.
	ErrOpenAIClientNotConfigured = errors.New("openai client is not configured")
	// ErrNoDocumentsAvailable indicates that there are no usable documents for the chat request.
	ErrNoDocumentsAvailable = errors.New("no documents available for chat")
	// ErrEmptyChatQuestion indicates an empty or whitespace-only question.
	ErrEmptyChatQuestion = errors.New("question is required")
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
	DocumentID uint           `json:"document_id"`
	Status     DocumentStatus `json:"status"`
	Message    string         `json:"message"`
}
