package service

import (
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/workflow/config"
	"gorm.io/gorm"
)

type WorkflowService struct {
	db     *gorm.DB
	config *config.WorkflowConfig
}

func NewWorkflowService(db *gorm.DB, cfg *config.WorkflowConfig) *WorkflowService {
	return &WorkflowService{
		db:     db,
		config: cfg,
	}
}
