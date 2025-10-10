package model

import (
	"time"

	"gorm.io/gorm"
)

// DocumentChunk 文档分片，用于知识库存储和向量化
type DocumentChunk struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	DocumentID uint           `json:"document_id" gorm:"index;not null"`
	Index      int            `json:"index" gorm:"not null"`             // 片段序号
	Content    string         `json:"content" gorm:"type:text;not null"` // 片段内容
	TokenCount int            `json:"token_count"`                       // 简单 token 数量统计
	Metadata   string         `json:"metadata" gorm:"type:jsonb"`        // 存储空间、分类等额外信息
	VectorID   string         `json:"vector_id" gorm:"size:64;index"`    // Qdrant 对应的向量ID
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
